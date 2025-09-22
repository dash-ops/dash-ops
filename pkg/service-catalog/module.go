package servicecatalog

import (
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	scK8sAdapter "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters"
	scAdapters "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters/http"
	scStorage "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters/storage"
	scControllers "github.com/dash-ops/dash-ops/pkg/service-catalog/controllers"
	"github.com/dash-ops/dash-ops/pkg/service-catalog/handlers"
	scLogic "github.com/dash-ops/dash-ops/pkg/service-catalog/logic"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
	"gopkg.in/yaml.v2"
)

// Module represents the service catalog module - main entry point for the plugin
type Module struct {
	controller        *scControllers.ServiceController
	handler           *handlers.HTTPHandler
	kubernetesAdapter *scK8sAdapter.KubernetesAdapter
	config            *ModuleConfig
}

// ModuleConfig represents configuration for the service catalog module
type ModuleConfig struct {
	// Repository implementations
	ServiceRepo    scPorts.ServiceRepository
	VersioningRepo scPorts.VersioningRepository
	HealthRepo     scPorts.HealthRepository

	// Service implementations
	KubernetesService scPorts.KubernetesService
	GitHubService     scPorts.GitHubService
}

// NewModule creates and initializes a new service catalog module (main factory)
func NewModule(config *ModuleConfig) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Validate required dependencies
	if config.ServiceRepo == nil {
		return nil, fmt.Errorf("service repository is required")
	}

	// Initialize logic components
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	// Initialize adapters
	serviceAdapter := scAdapters.NewServiceAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize kubernetes adapter if kubernetes service is available
	var kubernetesAdapter *scK8sAdapter.KubernetesAdapter
	if config.KubernetesService != nil {
		kubernetesAdapter = scK8sAdapter.NewKubernetesAdapter(config.KubernetesService)
	}

	// Initialize controller with injected dependencies
	controller := scControllers.NewServiceController(
		config.ServiceRepo,
		config.VersioningRepo, // Can be nil
		kubernetesAdapter,     // Use adapter instead of direct service
		config.GitHubService,  // Can be nil
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
		controller:        controller,
		handler:           handler,
		kubernetesAdapter: kubernetesAdapter,
		config:            config,
	}, nil
}

// RegisterRoutes registers HTTP routes for the service catalog module
func (m *Module) RegisterRoutes(router *mux.Router) {
	// Create service-catalog prefix subrouter (consistent with other modules)
	serviceCatalogRouter := router.PathPrefix("/service-catalog").Subrouter()
	m.handler.RegisterRoutes(serviceCatalogRouter)
}

// GetKubernetesAdapter returns adapter for kubernetes integration (for main.go compatibility)
func (m *Module) GetKubernetesAdapter() *KubernetesServiceContextAdapter {
	return NewKubernetesServiceContextAdapter(m.controller)
}

// UpdateKubernetesService updates the kubernetes service dependency
func (m *Module) UpdateKubernetesService(k8sService scPorts.KubernetesService) {
	// Create new adapter with the updated service
	m.kubernetesAdapter = scK8sAdapter.NewKubernetesAdapter(k8sService)
	m.controller.UpdateKubernetesService(m.kubernetesAdapter)
}

// KubernetesServiceContextAdapter adapts service catalog controller to kubernetes ServiceContextResolver
type KubernetesServiceContextAdapter struct {
	controller *scControllers.ServiceController
}

// NewKubernetesServiceContextAdapter creates a new adapter
func NewKubernetesServiceContextAdapter(controller *scControllers.ServiceController) *KubernetesServiceContextAdapter {
	return &KubernetesServiceContextAdapter{
		controller: controller,
	}
}

// ResolveDeploymentService implements kubernetes.ServiceContextResolver interface
func (k *KubernetesServiceContextAdapter) ResolveDeploymentService(deploymentName, namespace, context string) (*k8sModels.ServiceContext, error) {
	// Use the controller to resolve the service context
	serviceContext, err := k.controller.ResolveDeploymentService(nil, deploymentName, namespace, context)
	if err != nil {
		// Return nil if no service found (not an error)
		return nil, nil
	}

	if serviceContext == nil {
		return nil, nil
	}

	// Convert from service-catalog ServiceContext to kubernetes ServiceContext
	return &k8sModels.ServiceContext{
		ServiceName: serviceContext.Service.Metadata.Name,
		ServiceTier: string(serviceContext.Service.Metadata.Tier),
		Environment: serviceContext.Environment,
		Context:     serviceContext.Context,
		Team:        serviceContext.Service.Spec.Team.GitHubTeam,
		Description: serviceContext.Service.Spec.Description,
		Found:       serviceContext.Found,
	}, nil
}

// ServiceCatalog type alias for compatibility with main.go
type ServiceCatalog = Module

// Configuration parsing functions

type ServiceCatalogConfig struct {
	ServiceCatalog struct {
		Storage struct {
			Provider   string `yaml:"provider"`
			Filesystem struct {
				Directory string `yaml:"directory"`
			} `yaml:"filesystem"`
		} `yaml:"storage"`
	} `yaml:"service_catalog"`
}

type ParsedConfig struct {
	Directory string
}

// ParseServiceCatalogConfig parses service catalog configuration from YAML (exported for main.go)
func ParseServiceCatalogConfig(fileConfig []byte) (*ParsedConfig, error) {
	var config ServiceCatalogConfig
	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Extract directory from filesystem config
	directory := config.ServiceCatalog.Storage.Filesystem.Directory
	if directory == "" {
		directory = "../services" // Default directory
	}

	return &ParsedConfig{
		Directory: directory,
	}, nil
}

// NewFilesystemRepository creates a new filesystem repository (exported for main.go)
func NewFilesystemRepository(directory string) (scPorts.ServiceRepository, error) {
	return scStorage.NewFilesystemRepository(directory)
}
