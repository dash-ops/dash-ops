package http

import (
	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	settingsWire "github.com/dash-ops/dash-ops/pkg/settings/wire"
)

// ConfigAdapter handles transformation between models and wire formats.
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter.
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ModelToConfigResponse converts a DashConfig model to ConfigResponse.
func (ca *ConfigAdapter) ModelToConfigResponse(config *settingsModels.DashConfig) settingsWire.ConfigResponse {
	return settingsWire.ConfigResponse{
		Port:    config.GetPort(),
		Origin:  config.GetOrigin(),
		Headers: config.GetHeaders(),
		Front:   config.GetFront(),
		Plugins: config.Plugins.List(),
	}
}

// ModelToPluginsResponse converts plugins to PluginsResponse.
func (ca *ConfigAdapter) ModelToPluginsResponse(plugins settingsModels.Plugins) settingsWire.PluginsResponse {
	return settingsWire.PluginsResponse{
		Plugins: plugins.List(),
		Count:   plugins.Count(),
	}
}

// ModelToPluginsArray converts plugins to simple array for legacy compatibility.
func (ca *ConfigAdapter) ModelToPluginsArray(plugins settingsModels.Plugins) []string {
	return plugins.List()
}

// ModelToPluginStatusResponse converts plugin status to response.
func (ca *ConfigAdapter) ModelToPluginStatusResponse(pluginName string, config *settingsModels.DashConfig) settingsWire.PluginStatusResponse {
	return settingsWire.PluginStatusResponse{
		Name:    pluginName,
		Enabled: config.IsPluginEnabled(pluginName),
	}
}

// ModelToSystemInfoResponse converts config to system info response.
func (ca *ConfigAdapter) ModelToSystemInfoResponse(config *settingsModels.DashConfig, version, environment, uptime string) settingsWire.SystemInfoResponse {
	return settingsWire.SystemInfoResponse{
		Version:     version,
		Environment: environment,
		Uptime:      uptime,
		Plugins:     config.Plugins.List(),
	}
}
