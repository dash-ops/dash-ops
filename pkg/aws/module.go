package aws

import (
	"fmt"

	"github.com/gorilla/mux"

	awsAdaptersExternal "github.com/dash-ops/dash-ops/pkg/aws/adapters/external"
	awsAdaptersHttp "github.com/dash-ops/dash-ops/pkg/aws/adapters/http"
	awsAdaptersStorage "github.com/dash-ops/dash-ops/pkg/aws/adapters/storage"
	aws "github.com/dash-ops/dash-ops/pkg/aws/controllers"
	"github.com/dash-ops/dash-ops/pkg/aws/handlers"
	awsLogic "github.com/dash-ops/dash-ops/pkg/aws/logic"
	awsPorts "github.com/dash-ops/dash-ops/pkg/aws/ports"
	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
)

// Module represents the AWS module with all its components
type Module struct {
	// Core components
	Controller *aws.AWSController
	Handler    *handlers.HTTPHandler

	// Logic components
	Processor      *awsLogic.InstanceProcessor
	CostCalculator *awsLogic.CostCalculator

	// Adapters
	AWSAdapter      *awsAdaptersHttp.AWSAdapter
	ResponseAdapter *commonsHttp.ResponseAdapter
	RequestAdapter  *commonsHttp.RequestAdapter

	// Repositories (interfaces - implementations injected)
	AccountRepo  awsPorts.AccountRepository
	InstanceRepo awsPorts.InstanceRepository
	MetricsRepo  awsPorts.MetricsRepository

	// Services (interfaces - implementations injected)
	ClientService       awsPorts.AWSClientService
	NotificationService awsPorts.NotificationService
	AuditService        awsPorts.AuditService
}

// ModuleConfig represents configuration for the AWS module
type ModuleConfig struct {
	// Repository implementations
	AccountRepo  awsPorts.AccountRepository
	InstanceRepo awsPorts.InstanceRepository
	MetricsRepo  awsPorts.MetricsRepository

	// Service implementations
	ClientService       awsPorts.AWSClientService
	NotificationService awsPorts.NotificationService
	AuditService        awsPorts.AuditService
}

// NewModule creates and initializes a new AWS module
func NewModule(config *ModuleConfig) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Validate required dependencies
	if config.AccountRepo == nil {
		return nil, fmt.Errorf("account repository is required")
	}
	if config.InstanceRepo == nil {
		return nil, fmt.Errorf("instance repository is required")
	}

	// Initialize logic components
	processor := awsLogic.NewInstanceProcessor()
	costCalculator := awsLogic.NewCostCalculator()

	// Initialize adapters
	awsAdapter := awsAdaptersHttp.NewAWSAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize controller
	controller := aws.NewAWSController(
		config.AccountRepo,
		config.InstanceRepo,
		config.MetricsRepo, // Can be nil
		processor,
		costCalculator,
	)

	// Initialize handler
	handler := handlers.NewHTTPHandler(
		controller,
		awsAdapter,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		Controller:          controller,
		Handler:             handler,
		Processor:           processor,
		CostCalculator:      costCalculator,
		AWSAdapter:          awsAdapter,
		ResponseAdapter:     responseAdapter,
		RequestAdapter:      requestAdapter,
		AccountRepo:         config.AccountRepo,
		InstanceRepo:        config.InstanceRepo,
		MetricsRepo:         config.MetricsRepo,
		ClientService:       config.ClientService,
		NotificationService: config.NotificationService,
		AuditService:        config.AuditService,
	}, nil
}

// NewMinimalModule creates a minimal module with only required dependencies
func NewMinimalModule(
	accountRepo awsPorts.AccountRepository,
	instanceRepo awsPorts.InstanceRepository,
) (*Module, error) {
	config := &ModuleConfig{
		AccountRepo:  accountRepo,
		InstanceRepo: instanceRepo,
		// Optional dependencies are nil
	}

	return NewModule(config)
}

// GetController returns the AWS controller
func (m *Module) GetController() *aws.AWSController {
	return m.Controller
}

// GetHandler returns the HTTP handler
func (m *Module) GetHandler() *handlers.HTTPHandler {
	return m.Handler
}

// GetProcessor returns the instance processor
func (m *Module) GetProcessor() *awsLogic.InstanceProcessor {
	return m.Processor
}

// GetCostCalculator returns the cost calculator
func (m *Module) GetCostCalculator() *awsLogic.CostCalculator {
	return m.CostCalculator
}

// WithMetrics adds metrics repository to the module
func (m *Module) WithMetrics(metricsRepo awsPorts.MetricsRepository) *Module {
	m.MetricsRepo = metricsRepo
	// TODO: Recreate controller with new dependencies
	return m
}

// WithNotifications adds notification service to the module
func (m *Module) WithNotifications(notificationService awsPorts.NotificationService) *Module {
	m.NotificationService = notificationService
	return m
}

// WithAudit adds audit service to the module
func (m *Module) WithAudit(auditService awsPorts.AuditService) *Module {
	m.AuditService = auditService
	return m
}

// RegisterRoutes registers HTTP routes for the AWS module
func (m *Module) RegisterRoutes(router *mux.Router) {
	if m.Handler == nil {
		return
	}
	m.Handler.RegisterRoutes(router)
}

// Validate validates the module configuration
func (m *Module) Validate() error {
	if m.AccountRepo == nil {
		return fmt.Errorf("account repository is required")
	}

	if m.InstanceRepo == nil {
		return fmt.Errorf("instance repository is required")
	}

	if m.Controller == nil {
		return fmt.Errorf("controller is not initialized")
	}

	if m.Handler == nil {
		return fmt.Errorf("handler is not initialized")
	}

	return nil
}

// NewAccountRepositoryAdapter creates a new account repository adapter
func NewAccountRepositoryAdapter(fileConfig []byte) (awsPorts.AccountRepository, error) {
	return awsAdaptersStorage.NewAccountRepositoryAdapter(fileConfig)
}

// NewAWSClientServiceAdapter creates a new AWS client service adapter
func NewAWSClientServiceAdapter() awsPorts.AWSClientService {
	return awsAdaptersExternal.NewAWSClientServiceAdapter()
}

// NewInstanceRepositoryAdapter creates a new instance repository adapter
func NewInstanceRepositoryAdapter(
	awsClientService awsPorts.AWSClientService,
	accountRepo awsPorts.AccountRepository,
) awsPorts.InstanceRepository {
	return awsAdaptersStorage.NewInstanceRepositoryAdapter(awsClientService, accountRepo)
}
