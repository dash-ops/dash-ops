package servicecatalog

import (
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
	scConfigAdapter "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters/config"
	scAdaptersHttp "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters/http"
	scStorage "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters/storage"
	scControllers "github.com/dash-ops/dash-ops/pkg/service-catalog/controllers"
	"github.com/dash-ops/dash-ops/pkg/service-catalog/handlers"
	scInternal "github.com/dash-ops/dash-ops/pkg/service-catalog/integrations/kubernetes"
	scLogic "github.com/dash-ops/dash-ops/pkg/service-catalog/logic"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// Module represents the service catalog module - main entry point for the plugin
type Module struct {
	controller *scControllers.ServiceController
	handler    *handlers.HTTPHandler
	config     *scModels.ModuleConfig
}

// NewModule creates and initializes a new service catalog module (main factory)
func NewModule(fileConfig []byte) (*Module, error) {
	if fileConfig == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Parse configuration using adapter
	configAdapter := scConfigAdapter.NewConfigAdapter()
	moduleConfig, err := configAdapter.ParseModuleConfig(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Create filesystem repository (real implementation)
	serviceRepo, err := scStorage.NewFilesystemRepository(moduleConfig.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to create filesystem repository: %w", err)
	}

	// Initialize logic components
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	// Initialize adapters
	serviceAdapter := scAdaptersHttp.NewServiceAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize controller with injected dependencies
	controller := scControllers.NewServiceController(
		serviceRepo,
		nil, // VersioningRepo - can be added later
		nil, // KubernetesAdapter - can be added later
		nil, // GitHubService - can be added later
		validator,
		processor,
	)

	// Initialize handler
	handler := handlers.NewHTTPHandler(
		controller,
		serviceAdapter,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		controller: controller,
		handler:    handler,
		config:     moduleConfig,
	}, nil
}

// RegisterRoutes registers HTTP routes for the service catalog module
func (m *Module) RegisterRoutes(router *mux.Router) {
	// Create service-catalog prefix subrouter (consistent with other modules)
	serviceCatalogRouter := router.PathPrefix("/service-catalog").Subrouter()
	m.handler.RegisterRoutes(serviceCatalogRouter)
}

// LoadDependencies loads dependencies between modules after all modules are initialized
func (m *Module) LoadDependencies(modules map[string]interface{}) error {
	// Load kubernetes dependency if available
	if k8sModule, exists := modules["kubernetes"]; exists {
		if k8s, ok := k8sModule.(interface {
			GetServiceCatalogAdapter() scPorts.KubernetesService
		}); ok {
			if adapter := k8s.GetServiceCatalogAdapter(); adapter != nil {
				// Use the adapter directly from Kubernetes module
				m.controller.UpdateKubernetesService(adapter)
			}
		}
	}
	return nil
}

// GetServiceContextResolver returns the service context resolver for kubernetes integration
func (m *Module) GetServiceContextResolver() k8sPorts.ServiceContextResolver {
	// Use the new integration adapter for service context resolution
	return scInternal.NewServiceCatalogAdapter(m.controller.GetServiceRepository())
}
