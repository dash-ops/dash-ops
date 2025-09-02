package servicecatalog

import (
	"fmt"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons-new/adapters/http"
	scAdapters "github.com/dash-ops/dash-ops/pkg/service-catalog-new/adapters/http"
	servicecatalog "github.com/dash-ops/dash-ops/pkg/service-catalog-new/controllers"
	"github.com/dash-ops/dash-ops/pkg/service-catalog-new/handlers"
	scLogic "github.com/dash-ops/dash-ops/pkg/service-catalog-new/logic"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog-new/ports"
)

// Module represents the service catalog module with all its components
type Module struct {
	// Core components
	Controller *servicecatalog.ServiceController
	Handler    *handlers.HTTPHandler

	// Logic components
	Validator *scLogic.ServiceValidator
	Processor *scLogic.ServiceProcessor

	// Adapters
	ServiceAdapter  *scAdapters.ServiceAdapter
	ResponseAdapter *commonsHttp.ResponseAdapter
	RequestAdapter  *commonsHttp.RequestAdapter

	// Repositories (interfaces - implementations injected)
	ServiceRepo    scPorts.ServiceRepository
	VersioningRepo scPorts.VersioningRepository
	HealthRepo     scPorts.HealthRepository

	// Services (interfaces - implementations injected)
	KubernetesService scPorts.KubernetesService
	GitHubService     scPorts.GitHubService
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

// NewModule creates and initializes a new service catalog module
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

	// Initialize controller
	controller := servicecatalog.NewServiceController(
		config.ServiceRepo,
		config.VersioningRepo,    // Can be nil
		config.KubernetesService, // Can be nil
		config.GitHubService,     // Can be nil
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
		Controller:        controller,
		Handler:           handler,
		Validator:         validator,
		Processor:         processor,
		ServiceAdapter:    serviceAdapter,
		ResponseAdapter:   responseAdapter,
		RequestAdapter:    requestAdapter,
		ServiceRepo:       config.ServiceRepo,
		VersioningRepo:    config.VersioningRepo,
		HealthRepo:        config.HealthRepo,
		KubernetesService: config.KubernetesService,
		GitHubService:     config.GitHubService,
	}, nil
}

// NewMinimalModule creates a minimal module with only required dependencies
func NewMinimalModule(serviceRepo scPorts.ServiceRepository) (*Module, error) {
	config := &ModuleConfig{
		ServiceRepo: serviceRepo,
		// Optional dependencies are nil
	}

	return NewModule(config)
}

// GetController returns the service controller
func (m *Module) GetController() *servicecatalog.ServiceController {
	return m.Controller
}

// GetHandler returns the HTTP handler
func (m *Module) GetHandler() *handlers.HTTPHandler {
	return m.Handler
}

// GetValidator returns the service validator
func (m *Module) GetValidator() *scLogic.ServiceValidator {
	return m.Validator
}

// GetProcessor returns the service processor
func (m *Module) GetProcessor() *scLogic.ServiceProcessor {
	return m.Processor
}

// WithVersioning adds versioning repository to the module
func (m *Module) WithVersioning(versioningRepo scPorts.VersioningRepository) *Module {
	m.VersioningRepo = versioningRepo
	// TODO: Recreate controller with new dependencies
	return m
}

// WithKubernetes adds Kubernetes service to the module
func (m *Module) WithKubernetes(k8sService scPorts.KubernetesService) *Module {
	m.KubernetesService = k8sService
	// TODO: Recreate controller with new dependencies
	return m
}

// WithGitHub adds GitHub service to the module
func (m *Module) WithGitHub(githubService scPorts.GitHubService) *Module {
	m.GitHubService = githubService
	// TODO: Recreate controller with new dependencies
	return m
}

// Validate validates the module configuration
func (m *Module) Validate() error {
	if m.ServiceRepo == nil {
		return fmt.Errorf("service repository is required")
	}

	if m.Controller == nil {
		return fmt.Errorf("controller is not initialized")
	}

	if m.Handler == nil {
		return fmt.Errorf("handler is not initialized")
	}

	return nil
}
