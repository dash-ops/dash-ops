package config

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	httpAdapter "github.com/dash-ops/dash-ops/pkg/config/adapters/http"
	configControllers "github.com/dash-ops/dash-ops/pkg/config/controllers"
	"github.com/dash-ops/dash-ops/pkg/config/handlers"
	configLogic "github.com/dash-ops/dash-ops/pkg/config/logic"
	configModels "github.com/dash-ops/dash-ops/pkg/config/models"
)

// Module represents the config module - main entry point for the plugin
type Module struct {
	config     *configModels.DashConfig
	handler    *handlers.HTTPHandler
	controller *configControllers.ConfigController
}

// NewModule creates and initializes a new config module (main factory)
func NewModule(configFilePath string) (*Module, error) {
	// Initialize logic components
	processor := configLogic.NewConfigProcessor()

	// Load configuration
	var config *configModels.DashConfig
	var err error

	if configFilePath != "" {
		config, err = processor.LoadFromFile(configFilePath)
	} else {
		defaultPath := processor.GetConfigFilePath()
		config, err = processor.LoadFromFile(defaultPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize adapters
	configAdapter := httpAdapter.NewConfigAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()

	// Initialize controller
	controller := configControllers.NewConfigController(processor, config)

	// Initialize handler
	handler := handlers.NewHTTPHandler(controller, configAdapter, responseAdapter)

	return &Module{
		config:     config,
		controller: controller,
		handler:    handler,
	}, nil
}

// NewModuleFromBytes creates a module from configuration bytes
func NewModuleFromBytes(configData []byte) (*Module, error) {
	// Initialize logic components
	processor := configLogic.NewConfigProcessor()

	// Parse configuration
	config, err := processor.ParseFromBytes(configData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Initialize adapters
	configAdapter := httpAdapter.NewConfigAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()

	// Initialize controller
	controller := configControllers.NewConfigController(processor, config)

	// Initialize handler
	handler := handlers.NewHTTPHandler(controller, configAdapter, responseAdapter)

	return &Module{
		config:     config,
		controller: controller,
		handler:    handler,
	}, nil
}

// RegisterRoutes registers HTTP routes for the config module
func (m *Module) RegisterRoutes(router *mux.Router) {
	m.handler.RegisterRoutes(router)
}

// GetConfig returns the current configuration (for main.go compatibility)
func (m *Module) GetConfig() *configModels.DashConfig {
	return m.config
}

// GetPlugins returns enabled plugins (for main.go compatibility)
func (m *Module) GetPlugins() configModels.Plugins {
	return m.config.Plugins
}

// GetFileConfigBytes returns raw config file bytes (for dependency injection)
func (m *Module) GetFileConfigBytes() []byte {
	return GetFileGlobalConfig() // Use existing function
}

// Legacy compatibility functions for existing main.go

// GetFileGlobalConfig reads configuration file - legacy compatibility
func GetFileGlobalConfig() []byte {
	processor := configLogic.NewConfigProcessor()
	configPath := processor.GetConfigFilePath()

	// Read raw file bytes for compatibility with other modules
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err) // Maintain same behavior as original
	}

	// Expand environment variables like the original did
	return []byte(os.ExpandEnv(string(data)))
}
