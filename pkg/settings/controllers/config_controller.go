package controllers

import (
	"context"
	"fmt"
	"strings"

	settingsLogic "github.com/dash-ops/dash-ops/pkg/settings/logic"
	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	settingsPorts "github.com/dash-ops/dash-ops/pkg/settings/ports"
)

// ConfigController handles configuration business logic orchestration.
type ConfigController struct {
	processor  *settingsLogic.ConfigProcessor
	configPath string
	cache      settingsPorts.ConfigCache
}

// NewConfigController creates a new config controller.
func NewConfigController(
	processor *settingsLogic.ConfigProcessor,
	cache settingsPorts.ConfigCache,
	configPath string,
) *ConfigController {
	return &ConfigController{
		processor:  processor,
		cache:      cache,
		configPath: configPath,
	}
}

// GetConfig returns the current configuration.
func (cc *ConfigController) GetConfig(ctx context.Context) (*settingsModels.DashConfig, error) {
	config := cc.cache.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}
	return config.Clone(), nil
}

// GetPlugins returns the list of enabled plugins.
func (cc *ConfigController) GetPlugins(ctx context.Context) (settingsModels.Plugins, error) {
	config := cc.cache.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}
	return append(settingsModels.Plugins(nil), config.Plugins...), nil
}

// IsPluginEnabled checks if a specific plugin is enabled.
func (cc *ConfigController) IsPluginEnabled(ctx context.Context, pluginName string) (bool, error) {
	config := cc.cache.GetConfig()
	if config == nil {
		return false, fmt.Errorf("configuration not loaded")
	}

	if pluginName == "" {
		return false, fmt.Errorf("plugin name cannot be empty")
	}

	return config.IsPluginEnabled(pluginName), nil
}

// ReloadConfig reloads configuration from file.
func (cc *ConfigController) ReloadConfig(ctx context.Context) (*settingsModels.DashConfig, error) {
	if strings.TrimSpace(cc.configPath) == "" {
		return nil, fmt.Errorf("configuration path is not set")
	}

	newConfig, err := cc.processor.LoadFromFile(cc.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to reload configuration: %w", err)
	}

	cc.cache.SetConfig(newConfig.Clone())
	return newConfig, nil
}

// ValidateConfig validates the current configuration.
func (cc *ConfigController) ValidateConfig(ctx context.Context) error {
	config := cc.cache.GetConfig()
	if config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	return config.Validate()
}

// GetSystemInfo returns system information.
func (cc *ConfigController) GetSystemInfo(ctx context.Context, version, environment, uptime string) (*SystemInfo, error) {
	config := cc.cache.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	return &SystemInfo{
		Version:     version,
		Environment: environment,
		Uptime:      uptime,
		Config:      config.Clone(),
	}, nil
}

// SystemInfo represents system information.
type SystemInfo struct {
	Version     string
	Environment string
	Uptime      string
	Config      *settingsModels.DashConfig
}
