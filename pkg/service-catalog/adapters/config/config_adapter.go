package config

import (
	"fmt"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	"gopkg.in/yaml.v2"
)

// ConfigAdapter handles service catalog configuration parsing
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ParseServiceCatalogConfig parses service catalog configuration from YAML
func (ca *ConfigAdapter) ParseServiceCatalogConfig(fileConfig []byte) (*scModels.ParsedConfig, error) {
	var config scModels.ServiceCatalogConfig
	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Extract directory from filesystem config
	directory := config.ServiceCatalog.Storage.Filesystem.Directory
	if directory == "" {
		directory = "../services" // Default directory
	}

	return &scModels.ParsedConfig{
		Directory: directory,
	}, nil
}

// ParseModuleConfig parses the complete module configuration
func (ca *ConfigAdapter) ParseModuleConfig(fileConfig []byte) (*scModels.ModuleConfig, error) {
	parsedConfig, err := ca.ParseServiceCatalogConfig(fileConfig)
	if err != nil {
		return nil, err
	}

	return &scModels.ModuleConfig{
		Directory: parsedConfig.Directory,
	}, nil
}
