package config

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/config-new/logic"
	"github.com/dash-ops/dash-ops/pkg/config-new/models"
)

// ConfigController handles configuration business logic orchestration
type ConfigController struct {
	processor *logic.ConfigProcessor
	config    *models.DashConfig
}

// NewConfigController creates a new config controller
func NewConfigController(processor *logic.ConfigProcessor, config *models.DashConfig) *ConfigController {
	return &ConfigController{
		processor: processor,
		config:    config,
	}
}

// GetConfig returns the current configuration
func (cc *ConfigController) GetConfig(ctx context.Context) (*models.DashConfig, error) {
	if cc.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}
	return cc.config, nil
}

// GetPlugins returns the list of enabled plugins
func (cc *ConfigController) GetPlugins(ctx context.Context) (models.Plugins, error) {
	if cc.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}
	return cc.config.Plugins, nil
}

// IsPluginEnabled checks if a specific plugin is enabled
func (cc *ConfigController) IsPluginEnabled(ctx context.Context, pluginName string) (bool, error) {
	if cc.config == nil {
		return false, fmt.Errorf("configuration not loaded")
	}

	if pluginName == "" {
		return false, fmt.Errorf("plugin name cannot be empty")
	}

	return cc.config.IsPluginEnabled(pluginName), nil
}

// ReloadConfig reloads configuration from file
func (cc *ConfigController) ReloadConfig(ctx context.Context) (*models.DashConfig, error) {
	filePath := cc.processor.GetConfigFilePath()

	newConfig, err := cc.processor.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to reload configuration: %w", err)
	}

	cc.config = newConfig
	return cc.config, nil
}

// ValidateConfig validates the current configuration
func (cc *ConfigController) ValidateConfig(ctx context.Context) error {
	if cc.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	return cc.config.Validate()
}

// GetSystemInfo returns system information including configuration
func (cc *ConfigController) GetSystemInfo(ctx context.Context, version, environment, uptime string) (*SystemInfo, error) {
	if cc.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	return &SystemInfo{
		Version:     version,
		Environment: environment,
		Uptime:      uptime,
		Config:      cc.config,
	}, nil
}

// SystemInfo represents system information
type SystemInfo struct {
	Version     string
	Environment string
	Uptime      string
	Config      *models.DashConfig
}
