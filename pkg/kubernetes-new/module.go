package kubernetes

import (
	"fmt"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons-new/adapters/http"
	k8sAdaptersHttp "github.com/dash-ops/dash-ops/pkg/kubernetes-new/adapters/http"
	kubernetes "github.com/dash-ops/dash-ops/pkg/kubernetes-new/controllers"
	"github.com/dash-ops/dash-ops/pkg/kubernetes-new/handlers"
	k8sLogic "github.com/dash-ops/dash-ops/pkg/kubernetes-new/logic"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes-new/ports"
)

// Module represents the Kubernetes module with all its components
type Module struct {
	// Core components
	Controller *kubernetes.KubernetesController
	Handler    *handlers.HTTPHandler

	// Logic components
	HealthCalculator   *k8sLogic.HealthCalculator
	ResourceCalculator *k8sLogic.ResourceCalculator

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
	ClientService  k8sPorts.KubernetesClientService
	MetricsService k8sPorts.MetricsService
	EventService   k8sPorts.EventService
	HealthService  k8sPorts.HealthService
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
	resourceCalculator := k8sLogic.NewResourceCalculator()

	// Initialize adapters
	k8sAdapter := k8sAdaptersHttp.NewKubernetesAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize controller
	controller := kubernetes.NewKubernetesController(
		config.ClusterRepo,
		config.NodeRepo,
		config.NamespaceRepo,
		config.DeploymentRepo,
		config.PodRepo,
		healthCalculator,
		resourceCalculator,
	)

	// Initialize handler
	handler := handlers.NewHTTPHandler(
		controller,
		k8sAdapter,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		Controller:         controller,
		Handler:            handler,
		HealthCalculator:   healthCalculator,
		ResourceCalculator: resourceCalculator,
		K8sAdapter:         k8sAdapter,
		ResponseAdapter:    responseAdapter,
		RequestAdapter:     requestAdapter,
		ClusterRepo:        config.ClusterRepo,
		NodeRepo:           config.NodeRepo,
		NamespaceRepo:      config.NamespaceRepo,
		DeploymentRepo:     config.DeploymentRepo,
		PodRepo:            config.PodRepo,
		ClientService:      config.ClientService,
		MetricsService:     config.MetricsService,
		EventService:       config.EventService,
		HealthService:      config.HealthService,
	}, nil
}

// NewMinimalModule creates a minimal module with only required dependencies
func NewMinimalModule(
	clusterRepo k8sPorts.ClusterRepository,
	nodeRepo k8sPorts.NodeRepository,
	namespaceRepo k8sPorts.NamespaceRepository,
	deploymentRepo k8sPorts.DeploymentRepository,
	podRepo k8sPorts.PodRepository,
) (*Module, error) {
	config := &ModuleConfig{
		ClusterRepo:    clusterRepo,
		NodeRepo:       nodeRepo,
		NamespaceRepo:  namespaceRepo,
		DeploymentRepo: deploymentRepo,
		PodRepo:        podRepo,
		// Optional dependencies are nil
	}

	return NewModule(config)
}

// GetController returns the Kubernetes controller
func (m *Module) GetController() *kubernetes.KubernetesController {
	return m.Controller
}

// GetHandler returns the HTTP handler
func (m *Module) GetHandler() *handlers.HTTPHandler {
	return m.Handler
}

// GetHealthCalculator returns the health calculator
func (m *Module) GetHealthCalculator() *k8sLogic.HealthCalculator {
	return m.HealthCalculator
}

// GetResourceCalculator returns the resource calculator
func (m *Module) GetResourceCalculator() *k8sLogic.ResourceCalculator {
	return m.ResourceCalculator
}

// WithMetrics adds metrics service to the module
func (m *Module) WithMetrics(metricsService k8sPorts.MetricsService) *Module {
	m.MetricsService = metricsService
	return m
}

// WithEvents adds event service to the module
func (m *Module) WithEvents(eventService k8sPorts.EventService) *Module {
	m.EventService = eventService
	return m
}

// WithHealth adds health service to the module
func (m *Module) WithHealth(healthService k8sPorts.HealthService) *Module {
	m.HealthService = healthService
	return m
}

// Validate validates the module configuration
func (m *Module) Validate() error {
	if m.ClusterRepo == nil {
		return fmt.Errorf("cluster repository is required")
	}

	if m.NodeRepo == nil {
		return fmt.Errorf("node repository is required")
	}

	if m.NamespaceRepo == nil {
		return fmt.Errorf("namespace repository is required")
	}

	if m.DeploymentRepo == nil {
		return fmt.Errorf("deployment repository is required")
	}

	if m.PodRepo == nil {
		return fmt.Errorf("pod repository is required")
	}

	if m.Controller == nil {
		return fmt.Errorf("controller is not initialized")
	}

	if m.Handler == nil {
		return fmt.Errorf("handler is not initialized")
	}

	return nil
}
