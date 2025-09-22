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

// NewModule creates and initializes a new AWS module
func NewModule(fileConfig []byte) (*Module, error) {
	if fileConfig == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	accountRepo, err := awsAdaptersStorage.NewAccountRepositoryAdapter(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create account repository: %w", err)
	}

	// Create AWS client service
	awsClientService := awsAdaptersExternal.NewAWSClientServiceAdapter()

	// Create instance repository
	instanceRepo := awsAdaptersStorage.NewInstanceRepositoryAdapter(awsClientService, accountRepo)

	// Initialize logic components
	processor := awsLogic.NewInstanceProcessor()
	costCalculator := awsLogic.NewCostCalculator()

	// Initialize adapters
	awsAdapter := awsAdaptersHttp.NewAWSAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize controller
	controller := aws.NewAWSController(
		accountRepo,
		instanceRepo,
		nil, // TODO: Add metrics repository
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
		Controller:      controller,
		Handler:         handler,
		Processor:       processor,
		CostCalculator:  costCalculator,
		AWSAdapter:      awsAdapter,
		ResponseAdapter: responseAdapter,
		RequestAdapter:  requestAdapter,
	}, nil
}

// RegisterRoutes registers HTTP routes for the AWS module
func (m *Module) RegisterRoutes(router *mux.Router) {
	if m.Handler == nil {
		return
	}
	m.Handler.RegisterRoutes(router)
}
