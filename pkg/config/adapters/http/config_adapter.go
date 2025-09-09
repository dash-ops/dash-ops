package http

import (
	configModels "github.com/dash-ops/dash-ops/pkg/config/models"
	configWire "github.com/dash-ops/dash-ops/pkg/config/wire"
)

// ConfigAdapter handles transformation between models and wire formats
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ModelToConfigResponse converts a DashConfig model to ConfigResponse
func (ca *ConfigAdapter) ModelToConfigResponse(config *configModels.DashConfig) configWire.ConfigResponse {
	return configWire.ConfigResponse{
		Port:    config.GetPort(),
		Origin:  config.GetOrigin(),
		Headers: config.GetHeaders(),
		Front:   config.Front,
		Plugins: config.Plugins.List(),
	}
}

// ModelToPluginsResponse converts plugins to PluginsResponse
func (ca *ConfigAdapter) ModelToPluginsResponse(plugins configModels.Plugins) configWire.PluginsResponse {
	return configWire.PluginsResponse{
		Plugins: plugins.List(),
		Count:   plugins.Count(),
	}
}

// ModelToPluginsArray converts plugins to simple array for legacy compatibility
func (ca *ConfigAdapter) ModelToPluginsArray(plugins configModels.Plugins) []string {
	return plugins.List()
}

// ModelToPluginStatusResponse converts plugin status to response
func (ca *ConfigAdapter) ModelToPluginStatusResponse(pluginName string, config *configModels.DashConfig) configWire.PluginStatusResponse {
	return configWire.PluginStatusResponse{
		Name:    pluginName,
		Enabled: config.IsPluginEnabled(pluginName),
	}
}

// ModelToSystemInfoResponse converts config to system info response
func (ca *ConfigAdapter) ModelToSystemInfoResponse(config *configModels.DashConfig, version, environment, uptime string) configWire.SystemInfoResponse {
	return configWire.SystemInfoResponse{
		Version:     version,
		Environment: environment,
		Uptime:      uptime,
		Plugins:     config.Plugins.List(),
	}
}
