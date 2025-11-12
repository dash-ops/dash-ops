package logic

import (
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/settings/models"
)

// SetupProcessor handles setup-related business logic
type SetupProcessor struct {
	yamlProcessor *YAMLProcessor
}

// NewSetupProcessor creates a new setup processor
func NewSetupProcessor(yamlProcessor *YAMLProcessor) *SetupProcessor {
	return &SetupProcessor{
		yamlProcessor: yamlProcessor,
	}
}

// ProcessSetupConfig processes setup configuration and converts it to DashConfig
func (sp *SetupProcessor) ProcessSetupConfig(setupConfig *models.SetupConfig) (*models.DashConfig, error) {
	if setupConfig == nil {
		return nil, fmt.Errorf("setup config cannot be nil")
	}

	// Build enabled plugins list
	enabledPlugins := setupConfig.Plugins.EnabledPlugins

	// Create DashConfig
	dashConfig := &models.DashConfig{
		Port:    setupConfig.Port,
		Origin:  setupConfig.Origin,
		Headers: setupConfig.Headers,
		Front:   setupConfig.Front,
		Plugins: models.Plugins(enabledPlugins),
	}

	// Add plugin configurations
	if setupConfig.Plugins.Auth != nil && len(setupConfig.Plugins.Auth.Provider) > 0 {
		dashConfig.Auth = []models.AuthConfig{*setupConfig.Plugins.Auth}
	}

	if len(setupConfig.Plugins.Kubernetes) > 0 {
		dashConfig.Kubernetes = setupConfig.Plugins.Kubernetes
	}

	if len(setupConfig.Plugins.AWS) > 0 {
		dashConfig.AWS = setupConfig.Plugins.AWS
	}

	if setupConfig.Plugins.ServiceCatalog != nil {
		dashConfig.ServiceCatalog = setupConfig.Plugins.ServiceCatalog
	}

	if setupConfig.Plugins.Observability != nil {
		dashConfig.Observability = setupConfig.Plugins.Observability
	}

	// Generate YAML to validate
	_, err := sp.yamlProcessor.GenerateYAML(dashConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate YAML: %w", err)
	}

	return dashConfig, nil
}
