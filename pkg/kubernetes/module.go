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
	// Configuration data
	Configs []KubernetesConfig `yaml:"configs" json:"configs"`
}

// NewModule creates and initializes a new Kubernetes module
func NewModule(fileConfig []byte) (*Module, error) {
	if fileConfig == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Parse configuration
	configs, err := ParseKubernetesConfigFromFileConfig(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Convert to external config format
	var externalConfigs []k8sExternal.KubernetesConfig
	for _, config := range configs {
		externalConfigs = append(externalConfigs, k8sExternal.KubernetesConfig{
			Kubeconfig: config.Kubeconfig,
			Context:    config.Context,
		})
	}

	// Initialize logic components
	healthCalculator := k8sLogic.NewHealthCalculator()

	// Initialize adapters
	k8sAdapter := k8sAdaptersHttp.NewKubernetesAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize repositories with real configs
	clusterRepo, err := k8sExternal.NewClusterRepository(externalConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster repository: %w", err)
	}
	nodeRepo := k8sExternal.NewNodeRepository(clusterRepo)
	namespaceRepo := k8sExternal.NewNamespaceRepository(clusterRepo)
	deploymentRepo := k8sExternal.NewDeploymentRepository(clusterRepo, nil) // ServiceContextResolver will be set later
	podRepo := k8sExternal.NewPodRepository(clusterRepo)

	// Initialize service-catalog adapter
	serviceCatalogAdapter := k8sExternal.NewServiceCatalogAdapter(deploymentRepo, clusterRepo)

	// Initialize controller
	controller := kubernetes.NewKubernetesController(
		clusterRepo,
		nodeRepo,
		namespaceRepo,
		deploymentRepo,
		podRepo,
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
		ClusterRepo:            clusterRepo,
		NodeRepo:               nodeRepo,
		NamespaceRepo:          namespaceRepo,
		DeploymentRepo:         deploymentRepo,
		PodRepo:                podRepo,
		ClientService:          nil, // Can be added later
		MetricsService:         nil, // Can be added later
		EventService:           nil, // Can be added later
		HealthService:          nil, // Can be added later
		ServiceContextResolver: nil, // Will be set via LoadDependencies
		ServiceCatalogAdapter:  serviceCatalogAdapter,
	}, nil
}

// RegisterRoutes registers HTTP routes for the Kubernetes module
func (m *Module) RegisterRoutes(router *mux.Router) {
	// Create kubernetes prefix subrouter (consistent with other modules)
	k8sRouter := router.PathPrefix("/k8s").Subrouter()
	m.Handler.RegisterRoutes(k8sRouter)
}

// LoadDependencies loads dependencies between modules after all modules are initialized
func (m *Module) LoadDependencies(modules map[string]interface{}) error {
	// Load service-catalog dependency if available
	if scModule, exists := modules["service-catalog"]; exists {
		if sc, ok := scModule.(interface {
			GetServiceContextResolver() k8sPorts.ServiceContextResolver
		}); ok {
			if resolver := sc.GetServiceContextResolver(); resolver != nil {
				m.ServiceContextResolver = resolver
				// Update deployment repository with the resolver
				if dr, ok := m.DeploymentRepo.(*k8sExternal.DeploymentRepositoryImpl); ok {
					dr.SetServiceContextResolver(resolver)
				}
			}
		}
	}
	return nil
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
