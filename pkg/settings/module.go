package settings

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	settingsAdaptersHttp "github.com/dash-ops/dash-ops/pkg/settings/adapters/http"
	settingsControllers "github.com/dash-ops/dash-ops/pkg/settings/controllers"
	"github.com/dash-ops/dash-ops/pkg/settings/handlers"
	settingsLogic "github.com/dash-ops/dash-ops/pkg/settings/logic"
	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	settingsPorts "github.com/dash-ops/dash-ops/pkg/settings/ports"
	settingsRepositories "github.com/dash-ops/dash-ops/pkg/settings/repositories"
)

// Module represents the settings module - main entry point.
type Module struct {
	config          *settingsModels.DashConfig
	configPath      string
	configRepo      settingsPorts.ConfigRepository
	configProcessor *settingsLogic.ConfigProcessor
	settingsHandler *handlers.HTTPHandler
}

var _ settingsPorts.ConfigCache = (*Module)(nil)

// NewModule creates and initializes a new settings module.
func NewModule(configFilePath string) (*Module, error) {
	// Initialize logic components
	yamlProcessor := settingsLogic.NewYAMLProcessor()
	setupProcessor := settingsLogic.NewSetupProcessor(yamlProcessor)
	settingsProcessor := settingsLogic.NewSettingsProcessor(yamlProcessor)
	configProcessor := settingsLogic.NewConfigProcessor()

	// Resolve config path and initialize repository
	configPath := configProcessor.ResolveConfigFilePath(configFilePath)
	configRepo := settingsRepositories.NewFileRepository(configPath)

	module := &Module{
		configRepo:      configRepo,
		configProcessor: configProcessor,
		configPath:      configPath,
	}

	// Load configuration (or default) into cache
	if err := module.initializeConfig(context.Background()); err != nil {
		return nil, err
	}

	// Initialize repositories / controllers
	statusRepo := settingsRepositories.NewConfigStatusRepository(module)

	setupController := settingsControllers.NewSetupController(
		configRepo,
		statusRepo,
		setupProcessor,
		module,
	)

	settingsController := settingsControllers.NewSettingsController(
		configRepo,
		statusRepo,
		settingsProcessor,
		module,
	)

	configController := settingsControllers.NewConfigController(
		configProcessor,
		module,
		configPath,
	)

	// Initialize adapters
	settingsAdapter := settingsAdaptersHttp.NewSettingsAdapter()
	configAdapter := settingsAdaptersHttp.NewConfigAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize handlers
	module.settingsHandler = handlers.NewHTTPHandler(
		setupController,
		settingsController,
		configController,
		statusRepo,
		settingsAdapter,
		configAdapter,
		responseAdapter,
		requestAdapter,
	)

	return module, nil
}

// initializeConfig loads the configuration from disk or falls back to defaults.
func (m *Module) initializeConfig(ctx context.Context) error {
	exists, err := m.configRepo.ConfigExists(ctx, m.configPath)
	if err != nil {
		return fmt.Errorf("failed to check configuration existence: %w", err)
	}

	if exists {
		config, err := m.configRepo.LoadConfig(ctx, m.configPath)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
		m.SetConfig(config)
		return nil
	}

	// Use default configuration if file does not exist yet
	m.SetConfig(m.configProcessor.DefaultConfig())
	return nil
}

// RegisterRoutes registers HTTP routes for the settings module.
func (m *Module) RegisterRoutes(router *mux.Router) {
	if m.settingsHandler != nil {
		m.settingsHandler.RegisterRoutes(router)
	}
}

// GetConfig returns a copy of the current configuration.
func (m *Module) GetConfig() *settingsModels.DashConfig {
	if m.config == nil {
		return nil
	}
	return m.config.Clone()
}

// GetPlugins returns the current enabled plugins.
func (m *Module) GetPlugins() settingsModels.Plugins {
	if m.config == nil {
		return settingsModels.Plugins{}
	}
	return append(settingsModels.Plugins(nil), m.config.Plugins...)
}

// SetConfig updates the in-memory configuration cache.
func (m *Module) SetConfig(config *settingsModels.DashConfig) {
	if config == nil {
		m.config = nil
		return
	}
	m.config = config.Clone()
}

// GetConfigPath returns the current configuration file path.
func (m *Module) GetConfigPath() string {
	return m.configPath
}

// ReloadFromDisk reloads configuration from disk into the cache.
func (m *Module) ReloadFromDisk() error {
	if strings.TrimSpace(m.configPath) == "" {
		return fmt.Errorf("configuration path is not set")
	}

	config, err := m.configProcessor.LoadFromFile(m.configPath)
	if err != nil {
		return err
	}

	m.SetConfig(config)
	return nil
}

// GetFileConfigBytes returns raw config file bytes with environment expansion.
func (m *Module) GetFileConfigBytes() []byte {
	if strings.TrimSpace(m.configPath) == "" {
		return nil
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		panic(err)
	}

	return []byte(os.ExpandEnv(string(data)))
}
