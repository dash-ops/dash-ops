package config

import (
	"fmt"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons-new/adapters/http"
	httpAdapter "github.com/dash-ops/dash-ops/pkg/config-new/adapters/http"
	config "github.com/dash-ops/dash-ops/pkg/config-new/controllers"
	"github.com/dash-ops/dash-ops/pkg/config-new/handlers"
	config "github.com/dash-ops/dash-ops/pkg/config-new/logic"
	config "github.com/dash-ops/dash-ops/pkg/config-new/models"
)

// Module represents the config module with all its components
type Module struct {
	// Core components
	Config     *models.DashConfig
	Controller *config.ConfigController
	Handler    *handlers.HTTPHandler

	// Logic components
	Processor *logic.ConfigProcessor

	// Adapters
	ConfigAdapter   *httpAdapter.ConfigAdapter
	ResponseAdapter *commonsHttp.ResponseAdapter
}

// NewModule creates and initializes a new config module
func NewModule(configFilePath string) (*Module, error) {
	// Initialize logic components
	processor := logic.NewConfigProcessor()

	// Load configuration
	var config *models.DashConfig
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
	controller := config.NewConfigController(processor, config)

	// Initialize handler
	handler := handlers.NewHTTPHandler(controller, configAdapter, responseAdapter)

	return &Module{
		Config:          config,
		Controller:      controller,
		Handler:         handler,
		Processor:       processor,
		ConfigAdapter:   configAdapter,
		ResponseAdapter: responseAdapter,
	}, nil
}

// NewModuleFromBytes creates a module from configuration bytes
func NewModuleFromBytes(configData []byte) (*Module, error) {
	// Initialize logic components
	processor := logic.NewConfigProcessor()

	// Parse configuration
	config, err := processor.ParseFromBytes(configData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Initialize adapters
	configAdapter := httpAdapter.NewConfigAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()

	// Initialize controller
	controller := config.NewConfigController(processor, config)

	// Initialize handler
	handler := handlers.NewHTTPHandler(controller, configAdapter, responseAdapter)

	return &Module{
		Config:          config,
		Controller:      controller,
		Handler:         handler,
		Processor:       processor,
		ConfigAdapter:   configAdapter,
		ResponseAdapter: responseAdapter,
	}, nil
}

// GetConfig returns the current configuration
func (m *Module) GetConfig() *models.DashConfig {
	return m.Config
}

// GetController returns the config controller
func (m *Module) GetController() *config.ConfigController {
	return m.Controller
}

// GetHandler returns the HTTP handler
func (m *Module) GetHandler() *handlers.HTTPHandler {
	return m.Handler
}

// Reload reloads the configuration
func (m *Module) Reload() error {
	newConfig, err := m.Controller.ReloadConfig(nil)
	if err != nil {
		return fmt.Errorf("failed to reload configuration: %w", err)
	}

	m.Config = newConfig
	return nil
}
