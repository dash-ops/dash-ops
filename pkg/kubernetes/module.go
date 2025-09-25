package kubernetes

import (
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/adapters/config"
	k8sAdaptersHttp "github.com/dash-ops/dash-ops/pkg/kubernetes/adapters/http"
	kubernetes "github.com/dash-ops/dash-ops/pkg/kubernetes/controllers"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/handlers"
	k8sExternalIntegration "github.com/dash-ops/dash-ops/pkg/kubernetes/integrations/external/kubernetes"
	k8sInternal "github.com/dash-ops/dash-ops/pkg/kubernetes/integrations/service-catalog"
	k8sLogic "github.com/dash-ops/dash-ops/pkg/kubernetes/logic"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
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

	// External integrations
	KubernetesAdapter k8sPorts.KubernetesClientService
}

// NewModule creates and initializes a new Kubernetes module
func NewModule(fileConfig []byte) (*Module, error) {
	if fileConfig == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Parse configuration using adapter
	configAdapter := config.NewConfigAdapter()
	moduleConfig, err := configAdapter.ParseModuleConfig(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Initialize logic components
	healthCalculator := k8sLogic.NewHealthCalculator()

	// Initialize external integrations
	var kubernetesAdapter k8sPorts.KubernetesClientService
	if len(moduleConfig.Configs) > 0 {
		// Use the first config for now (in a multi-cluster setup, this would be different)
		config := &k8sExternalIntegration.KubernetesConfig{
			Kubeconfig: moduleConfig.Configs[0].Kubeconfig,
			Context:    moduleConfig.Configs[0].Context,
		}
		adapter, err := k8sExternalIntegration.NewKubernetesAdapter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes adapter: %w", err)
		}
		kubernetesAdapter = adapter
	}

	// Initialize adapters
	k8sAdapter := k8sAdaptersHttp.NewKubernetesAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Use wrapper functions to create repositories from the adapter
	var clusterRepo k8sPorts.ClusterRepository
	var nodeRepo k8sPorts.NodeRepository
	var namespaceRepo k8sPorts.NamespaceRepository
	var deploymentRepo k8sPorts.DeploymentRepository
	var podRepo k8sPorts.PodRepository

	if kubernetesAdapter != nil {
		// Cast to concrete type to access wrapper functions
		if adapter, ok := kubernetesAdapter.(*k8sExternalIntegration.KubernetesAdapter); ok {
			clusterRepo = k8sExternalIntegration.NewClusterRepository(adapter)
			nodeRepo = k8sExternalIntegration.NewNodeRepository(adapter)
			namespaceRepo = k8sExternalIntegration.NewNamespaceRepository(adapter)
			deploymentRepo = k8sExternalIntegration.NewDeploymentRepository(adapter)
			podRepo = k8sExternalIntegration.NewPodRepository(adapter)
		}
	}

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
		KubernetesAdapter:      kubernetesAdapter,
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
				if dr, ok := m.DeploymentRepo.(interface {
					SetServiceContextResolver(k8sPorts.ServiceContextResolver)
				}); ok {
					dr.SetServiceContextResolver(resolver)
				}
			}
		}
	}
	return nil
}

// GetServiceCatalogAdapter returns the adapter for service-catalog integration
func (m *Module) GetServiceCatalogAdapter() scPorts.KubernetesService {
	return k8sInternal.NewKubernetesAdapter(m.DeploymentRepo, m.ClusterRepo)
}
