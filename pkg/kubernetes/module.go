package kubernetes

import (
	"fmt"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	k8sExternal "github.com/dash-ops/dash-ops/pkg/kubernetes/adapters/external"
	k8sAdaptersHttp "github.com/dash-ops/dash-ops/pkg/kubernetes/adapters/http"
	kubernetes "github.com/dash-ops/dash-ops/pkg/kubernetes/controllers"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/handlers"
	k8sLogic "github.com/dash-ops/dash-ops/pkg/kubernetes/logic"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Module represents the Kubernetes module with all its components
type Module struct {
	// Core components
	Controller *kubernetes.KubernetesController
	Handler    *handlers.HTTPHandler

	// Logic components
	HealthCalculator *k8sLogic.HealthCalculator

	// Adapters
	K8sAdapter      *k8sAdaptersHttp.KubernetesAdapter
	ResponseAdapter *commonsHttp.ResponseAdapter
	RequestAdapter  *commonsHttp.RequestAdapter

	// Repositories (interfaces - implementations injected)
	ClusterRepo    k8sPorts.ClusterRepository
	NodeRepo       k8sPorts.NodeRepository
	NamespaceRepo  k8sPorts.NamespaceRepository
	DeploymentRepo k8sPorts.DeploymentRepository
	PodRepo        k8sPorts.PodRepository

	// Services (interfaces - implementations injected)
	ClientService          k8sPorts.KubernetesClientService
	MetricsService         k8sPorts.MetricsService
	EventService           k8sPorts.EventService
	HealthService          k8sPorts.HealthService
	ServiceContextResolver k8sPorts.ServiceContextResolver

	// Adapters for external integrations
	ServiceCatalogAdapter scPorts.KubernetesService // Adapter for service-catalog integration
}

// ModuleConfig represents configuration for the Kubernetes module
type ModuleConfig struct {
	// Repository implementations
	ClusterRepo    k8sPorts.ClusterRepository
	NodeRepo       k8sPorts.NodeRepository
	NamespaceRepo  k8sPorts.NamespaceRepository
	DeploymentRepo k8sPorts.DeploymentRepository
	PodRepo        k8sPorts.PodRepository

	// Service implementations
	ClientService  k8sPorts.KubernetesClientService
	MetricsService k8sPorts.MetricsService
	EventService   k8sPorts.EventService
	HealthService  k8sPorts.HealthService
}

// NewModule creates and initializes a new Kubernetes module
func NewModule(config *ModuleConfig) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Validate required dependencies
	if config.ClusterRepo == nil {
		return nil, fmt.Errorf("cluster repository is required")
	}
	if config.NodeRepo == nil {
		return nil, fmt.Errorf("node repository is required")
	}
	if config.NamespaceRepo == nil {
		return nil, fmt.Errorf("namespace repository is required")
	}
	if config.DeploymentRepo == nil {
		return nil, fmt.Errorf("deployment repository is required")
	}
	if config.PodRepo == nil {
		return nil, fmt.Errorf("pod repository is required")
	}

	// Initialize logic components
	healthCalculator := k8sLogic.NewHealthCalculator()

	// Initialize adapters
	k8sAdapter := k8sAdaptersHttp.NewKubernetesAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize service-catalog adapter
	serviceCatalogAdapter := k8sExternal.NewServiceCatalogAdapter(config.DeploymentRepo, config.ClusterRepo)

	// Initialize controller
	controller := kubernetes.NewKubernetesController(
		config.ClusterRepo,
		config.NodeRepo,
		config.NamespaceRepo,
		config.DeploymentRepo,
		config.PodRepo,
		healthCalculator,
	)

	// Initialize handler
	handler := handlers.NewHTTPHandler(
		controller,
		k8sAdapter,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		Controller:             controller,
		Handler:                handler,
		HealthCalculator:       healthCalculator,
		K8sAdapter:             k8sAdapter,
		ResponseAdapter:        responseAdapter,
		RequestAdapter:         requestAdapter,
		ClusterRepo:            config.ClusterRepo,
		NodeRepo:               config.NodeRepo,
		NamespaceRepo:          config.NamespaceRepo,
		DeploymentRepo:         config.DeploymentRepo,
		PodRepo:                config.PodRepo,
		ClientService:          config.ClientService,
		MetricsService:         config.MetricsService,
		EventService:           config.EventService,
		HealthService:          config.HealthService,
		ServiceContextResolver: nil, // TODO: Add ServiceContextResolver to ModuleConfig
		ServiceCatalogAdapter:  serviceCatalogAdapter,
	}, nil
}

// RegisterRoutes registers HTTP routes for the Kubernetes module
func (m *Module) RegisterRoutes(router *mux.Router) {
	// Create kubernetes prefix subrouter (consistent with other modules)
	k8sRouter := router.PathPrefix("/k8s").Subrouter()
	m.Handler.RegisterRoutes(k8sRouter)
}

// GetServiceCatalogAdapter returns the adapter for service-catalog integration
func (m *Module) GetServiceCatalogAdapter() scPorts.KubernetesService {
	return m.ServiceCatalogAdapter
}

// ParseKubernetesConfigFromFileConfig parses kubernetes config from file bytes
func ParseKubernetesConfigFromFileConfig(fileConfig []byte) ([]KubernetesConfig, error) {
	var dashYaml struct {
		Kubernetes []KubernetesConfig `yaml:"kubernetes"`
	}

	err := yaml.Unmarshal(fileConfig, &dashYaml)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubernetes config: %w", err)
	}

	if len(dashYaml.Kubernetes) == 0 {
		return nil, fmt.Errorf("no kubernetes configuration found")
	}

	// Return all kubernetes configs
	return dashYaml.Kubernetes, nil
}

// KubernetesConfig represents kubernetes configuration
type KubernetesConfig struct {
	Name       string     `yaml:"name"`
	Kubeconfig string     `yaml:"kubeconfig"`
	Context    string     `yaml:"context"`
	Permission Permission `yaml:"permission"`
	Listen     string     `yaml:"-"`
}

// Permission represents kubernetes permissions
type Permission struct {
	Deployments DeploymentsPermissions `yaml:"deployments" json:"deployments"`
}

// DeploymentsPermissions represents deployment permissions
type DeploymentsPermissions struct {
	Namespaces []string `yaml:"namespaces" json:"namespaces"`
	Restart    []string `yaml:"restart" json:"restart"`
	Scale      []string `yaml:"scale" json:"scale"`
}
