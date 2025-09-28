package aws

import (
	"fmt"

	"github.com/gorilla/mux"

	awsAdaptersConfig "github.com/dash-ops/dash-ops/pkg/aws/adapters/config"
	"github.com/dash-ops/dash-ops/pkg/aws/handlers"
	awsIntegrations "github.com/dash-ops/dash-ops/pkg/aws/integrations/external/aws"
	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
)

// Module represents the AWS module with all its components
type Module struct {
	Handler *handlers.HTTPHandler
}

// NewModule creates and initializes a new AWS module
func NewModule(fileConfig []byte) (*Module, error) {
	if fileConfig == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Parse AWS configuration
	configAdapter := awsAdaptersConfig.NewConfigAdapter()
	accounts, err := configAdapter.ParseAWSConfigFromFileConfig(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AWS configuration: %w", err)
	}

	// Create AWS client service
	awsClientService := awsIntegrations.NewAWSAdapter()

	// Initialize adapters
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize handler with DI
	handler := handlers.NewHTTPHandler(
		awsClientService,
		accounts,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		Handler: handler,
	}, nil
}

// LoadDependencies loads dependencies between modules after all modules are initialized
func (m *Module) LoadDependencies(modules map[string]interface{}) error {
	// AWS module doesn't have cross-module dependencies
	return nil
}

// RegisterRoutes registers HTTP routes for the AWS module
func (m *Module) RegisterRoutes(router *mux.Router) {
	if m.Handler == nil {
		return
	}
	m.Handler.RegisterRoutes(router)
}
