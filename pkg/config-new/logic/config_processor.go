package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	configModels "github.com/dash-ops/dash-ops/pkg/config-new/models"
)

// ConfigProcessor handles configuration processing logic
type ConfigProcessor struct{}

// NewConfigProcessor creates a new config processor
func NewConfigProcessor() *ConfigProcessor {
	return &ConfigProcessor{}
}

// LoadFromFile loads configuration from a YAML file
func (cp *ConfigProcessor) LoadFromFile(filePath string) (*configModels.DashConfig, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Read file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", absPath, err)
	}

	return cp.ParseFromBytes(data)
}

// ParseFromBytes parses configuration from byte data
func (cp *ConfigProcessor) ParseFromBytes(data []byte) (*configModels.DashConfig, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("configuration data cannot be empty")
	}

	// Expand environment variables
	expandedData := os.ExpandEnv(string(data))

	// Parse YAML
	var config configModels.DashConfig
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML configuration: %w", err)
	}

	// Apply defaults
	cp.applyDefaults(&config)

	// Validate
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// GetConfigFilePath returns the configuration file path
func (cp *ConfigProcessor) GetConfigFilePath() string {
	// Check environment variable first
	if path := os.Getenv("DASH_CONFIG"); path != "" {
		return path
	}

	// Default path
	return "./dash-ops.yaml"
}

// applyDefaults applies default values to configuration
func (cp *ConfigProcessor) applyDefaults(config *configModels.DashConfig) {
	if config.Port == "" {
		config.Port = "8080"
	}

	if config.Origin == "" {
		config.Origin = "http://localhost:3000"
	}

	if len(config.Headers) == 0 {
		config.Headers = []string{"Content-Type", "Authorization"}
	}

	// Normalize plugin names (remove duplicates and empty strings)
	config.Plugins = cp.normalizePlugins(config.Plugins)
}

// normalizePlugins removes duplicates and empty plugin names
func (cp *ConfigProcessor) normalizePlugins(plugins configModels.Plugins) configModels.Plugins {
	seen := make(map[string]bool)
	var normalized configModels.Plugins

	for _, plugin := range plugins {
		plugin = strings.TrimSpace(plugin)
		if plugin != "" && !seen[strings.ToLower(plugin)] {
			seen[strings.ToLower(plugin)] = true
			normalized = append(normalized, plugin)
		}
	}

	return normalized
}

// MergeConfigs merges two configurations, with override taking precedence
func (cp *ConfigProcessor) MergeConfigs(base, override *configModels.DashConfig) *configModels.DashConfig {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := *base // Copy base

	// Override non-empty values
	if override.Port != "" {
		merged.Port = override.Port
	}
	if override.Origin != "" {
		merged.Origin = override.Origin
	}
	if override.Front != "" {
		merged.Front = override.Front
	}
	if len(override.Headers) > 0 {
		merged.Headers = override.Headers
	}
	if len(override.Plugins) > 0 {
		merged.Plugins = override.Plugins
	}

	return &merged
}
