package servicecatalog

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
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
	// Configuration data
	Directory string `yaml:"directory" json:"directory"`
}

// NewModule creates and initializes a new service catalog module (main factory)
func NewModule(fileConfig []byte) (*Module, error) {
	if fileConfig == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Parse configuration
	config, err := ParseServiceCatalogConfig(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Create filesystem repository (real implementation)
	serviceRepo, err := scStorage.NewFilesystemRepository(config.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to create filesystem repository: %w", err)
	}

	// Initialize logic components
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	// Initialize adapters
	serviceAdapter := scAdapters.NewServiceAdapter()
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
		controller:        controller,
		handler:           handler,
		kubernetesAdapter: nil, // Can be added later
		config: &ModuleConfig{
			Directory: config.Directory,
		},
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
func (k *KubernetesServiceContextAdapter) ResolveDeploymentService(deploymentName, namespace, contextName string) (*k8sModels.ServiceContext, error) {
	// Use the controller to resolve the service context
	serviceContext, err := k.controller.ResolveDeploymentService(context.TODO(), deploymentName, namespace, contextName)
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

// LoadDependencies loads dependencies between modules after all modules are initialized
func (m *Module) LoadDependencies(modules map[string]interface{}) error {
	// Load kubernetes dependency if available
	if k8sModule, exists := modules["kubernetes"]; exists {
		if k8s, ok := k8sModule.(interface {
			GetServiceCatalogAdapter() scPorts.KubernetesService
		}); ok {
			if adapter := k8s.GetServiceCatalogAdapter(); adapter != nil {
				m.kubernetesAdapter = scK8sAdapter.NewKubernetesAdapter(adapter)
				m.controller.UpdateKubernetesService(m.kubernetesAdapter)
			}
		}
	}
	return nil
}

// GetServiceContextResolver returns the service context resolver for kubernetes integration
func (m *Module) GetServiceContextResolver() k8sPorts.ServiceContextResolver {
	// Always return the KubernetesServiceContextAdapter which handles service resolution
	return NewKubernetesServiceContextAdapter(m.controller)
}

// NewFilesystemRepository creates a new filesystem repository (exported for main.go)
func NewFilesystemRepository(directory string) (scPorts.ServiceRepository, error) {
	return scStorage.NewFilesystemRepository(directory)
}
