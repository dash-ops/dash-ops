package controllers

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/settings/logic"
	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	settingsPorts "github.com/dash-ops/dash-ops/pkg/settings/ports"
)

// SetupController handles setup-related use cases
type SetupController struct {
	configRepo     settingsPorts.ConfigRepository
	statusRepo     settingsPorts.ConfigStatusRepository
	setupProcessor *logic.SetupProcessor
	configCache    settingsPorts.ConfigCache
}

// NewSetupController creates a new setup controller
func NewSetupController(
	configRepo settingsPorts.ConfigRepository,
	statusRepo settingsPorts.ConfigStatusRepository,
	setupProcessor *logic.SetupProcessor,
	configCache settingsPorts.ConfigCache,
) *SetupController {
	return &SetupController{
		configRepo:     configRepo,
		statusRepo:     statusRepo,
		setupProcessor: setupProcessor,
		configCache:    configCache,
	}
}

// GetSetupStatus returns the current setup status
func (sc *SetupController) GetSetupStatus(ctx context.Context) (*settingsModels.SetupStatus, error) {
	return sc.statusRepo.GetSetupStatus(ctx)
}

// ConfigureSetup configures the initial setup
func (sc *SetupController) ConfigureSetup(ctx context.Context, setupConfig *settingsModels.SetupConfig) (string, error) {
	// Check if plugins are already configured
	hasPlugins, err := sc.statusRepo.HasPluginsConfigured(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check plugins status: %w", err)
	}

	if hasPlugins {
		return "", fmt.Errorf("plugins are already configured. Use settings endpoint to modify configuration")
	}

	// Process setup config
	dashConfig, err := sc.setupProcessor.ProcessSetupConfig(setupConfig)
	if err != nil {
		return "", fmt.Errorf("failed to process setup config: %w", err)
	}

	// Get config file path
	configPath := sc.configRepo.GetDefaultConfigPath()

	// Save configuration
	if err := sc.configRepo.SaveConfig(ctx, dashConfig, configPath); err != nil {
		return "", fmt.Errorf("failed to save configuration: %w", err)
	}

	// Update in-memory configuration cache
	sc.configCache.SetConfig(dashConfig.Clone())

	return configPath, nil
}
