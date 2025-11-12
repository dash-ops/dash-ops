package controllers

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/settings/logic"
	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	settingsPorts "github.com/dash-ops/dash-ops/pkg/settings/ports"
)

// SettingsController handles settings-related use cases
type SettingsController struct {
	configRepo        settingsPorts.ConfigRepository
	statusRepo        settingsPorts.ConfigStatusRepository
	settingsProcessor *logic.SettingsProcessor
	configCache       settingsPorts.ConfigCache
}

// NewSettingsController creates a new settings controller
func NewSettingsController(
	configRepo settingsPorts.ConfigRepository,
	statusRepo settingsPorts.ConfigStatusRepository,
	settingsProcessor *logic.SettingsProcessor,
	configCache settingsPorts.ConfigCache,
) *SettingsController {
	return &SettingsController{
		configRepo:        configRepo,
		statusRepo:        statusRepo,
		settingsProcessor: settingsProcessor,
		configCache:       configCache,
	}
}

// GetSettings returns the current settings configuration
func (sc *SettingsController) GetSettings(ctx context.Context) (*settingsModels.SettingsConfig, error) {
	// Get config file path
	configPath := sc.configRepo.GetDefaultConfigPath()

	currentConfig := sc.configCache.GetConfig()
	if currentConfig == nil {
		loadedConfig, err := sc.configRepo.LoadConfig(ctx, configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %w", err)
		}
		currentConfig = loadedConfig
		sc.configCache.SetConfig(loadedConfig.Clone())
	}

	return &settingsModels.SettingsConfig{
		Config:  currentConfig,
		Plugins: currentConfig.Plugins.List(),
		CanEdit: true,
	}, nil
}

// UpdateSettings updates the settings configuration
func (sc *SettingsController) UpdateSettings(
	ctx context.Context,
	request *settingsModels.UpdateSettingsRequest,
) (*settingsModels.UpdateSettingsResponse, error) {
	// Get config file path
	configPath := sc.configRepo.GetDefaultConfigPath()

	// Load current configuration
	currentConfig := sc.configCache.GetConfig()
	if currentConfig == nil {
		var err error
		currentConfig, err = sc.configRepo.LoadConfig(ctx, configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load current configuration: %w", err)
		}
	}

	// Process update request
	updatedConfig, err := sc.settingsProcessor.ProcessUpdateRequest(currentConfig, request)
	if err != nil {
		return nil, fmt.Errorf("failed to process update request: %w", err)
	}

	// Check if restart is required
	requiresRestart := sc.settingsProcessor.RequiresRestart(currentConfig, updatedConfig)

	// Save updated configuration
	if err := sc.configRepo.SaveConfig(ctx, updatedConfig, configPath); err != nil {
		return nil, fmt.Errorf("failed to save updated configuration: %w", err)
	}

	// Update in-memory cache with cloned config to avoid unintended mutations
	sc.configCache.SetConfig(updatedConfig.Clone())

	return &settingsModels.UpdateSettingsResponse{
		Success:         true,
		Message:         "Settings updated successfully",
		RequiresRestart: requiresRestart,
	}, nil
}
