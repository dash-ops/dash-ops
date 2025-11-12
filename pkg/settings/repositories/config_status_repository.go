package repositories

import (
	"context"

	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	settingsPorts "github.com/dash-ops/dash-ops/pkg/settings/ports"
)

// ConfigStatusRepository implements ConfigStatusRepository using config module
type ConfigStatusRepository struct {
	configProvider settingsPorts.ConfigProvider
}

// NewConfigStatusRepository creates a new config status repository
func NewConfigStatusRepository(configProvider settingsPorts.ConfigProvider) settingsPorts.ConfigStatusRepository {
	return &ConfigStatusRepository{
		configProvider: configProvider,
	}
}

// GetSetupStatus returns the current setup status
func (csr *ConfigStatusRepository) GetSetupStatus(ctx context.Context) (*settingsModels.SetupStatus, error) {
	config := csr.configProvider.GetConfig()
	plugins := csr.configProvider.GetPlugins()

	// If config is nil or no plugins configured, setup is required
	if config == nil || plugins.Count() == 0 {
		return &settingsModels.SetupStatus{
			SetupRequired: true,
			PluginsCount:  0,
			HasAuth:       false,
		}, nil
	}

	pluginsCount := plugins.Count()
	hasAuth := config.IsPluginEnabled("Auth") || config.IsPluginEnabled("auth")

	return &settingsModels.SetupStatus{
		SetupRequired: false,
		PluginsCount:  pluginsCount,
		HasAuth:       hasAuth,
	}, nil
}

// HasPluginsConfigured checks if any plugins are configured
func (csr *ConfigStatusRepository) HasPluginsConfigured(ctx context.Context) (bool, error) {
	plugins := csr.configProvider.GetPlugins()
	return plugins.Count() > 0, nil
}

// HasAuthConfigured checks if auth is configured
func (csr *ConfigStatusRepository) HasAuthConfigured(ctx context.Context) (bool, error) {
	config := csr.configProvider.GetConfig()
	if config == nil {
		return false, nil
	}

	return config.IsPluginEnabled("Auth") || config.IsPluginEnabled("auth"), nil
}
