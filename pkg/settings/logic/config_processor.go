package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	"gopkg.in/yaml.v2"
)

// ConfigProcessor handles configuration processing logic.
type ConfigProcessor struct{}

// NewConfigProcessor creates a new config processor.
func NewConfigProcessor() *ConfigProcessor {
	return &ConfigProcessor{}
}

// LoadFromFile loads configuration from a YAML file.
func (cp *ConfigProcessor) LoadFromFile(filePath string) (*settingsModels.DashConfig, error) {
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

// ParseFromBytes parses configuration from byte data.
func (cp *ConfigProcessor) ParseFromBytes(data []byte) (*settingsModels.DashConfig, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("configuration data cannot be empty")
	}

	// Expand environment variables
	expandedData := os.ExpandEnv(string(data))

	// Parse YAML
	var config settingsModels.DashConfig
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

// GetConfigFilePath returns the configuration file path.
func (cp *ConfigProcessor) GetConfigFilePath() string {
	// Check environment variable first
	if path := os.Getenv("DASH_CONFIG"); path != "" {
		return path
	}

	// Default path
	return "./dash-ops.yaml"
}

// ResolveConfigFilePath returns a normalized configuration path, preferring the provided value.
func (cp *ConfigProcessor) ResolveConfigFilePath(providedPath string) string {
	if strings.TrimSpace(providedPath) != "" {
		return providedPath
	}
	return cp.GetConfigFilePath()
}

// DefaultConfig returns a default configuration used when no file exists.
func (cp *ConfigProcessor) DefaultConfig() *settingsModels.DashConfig {
	return &settingsModels.DashConfig{
		Port:    "8080",
		Origin:  "http://localhost:5173",
		Headers: []string{"Content-Type", "Authorization"},
		Front:   "front/dist",
		Plugins: settingsModels.Plugins{},
	}
}

// applyDefaults applies default values to configuration.
func (cp *ConfigProcessor) applyDefaults(config *settingsModels.DashConfig) {
	if config.Port == "" {
		config.Port = "8080"
	}

	if config.Origin == "" {
		config.Origin = "http://localhost:5173"
	}

	if len(config.Headers) == 0 {
		config.Headers = []string{"Content-Type", "Authorization"}
	}

	if config.Front == "" {
		config.Front = "front/dist"
	}

	// Normalize plugin names (remove duplicates and empty strings)
	config.Plugins = cp.normalizePlugins(config.Plugins)
}

// normalizePlugins removes duplicates and empty plugin names.
func (cp *ConfigProcessor) normalizePlugins(plugins settingsModels.Plugins) settingsModels.Plugins {
	seen := make(map[string]bool)
	var normalized settingsModels.Plugins

	for _, plugin := range plugins {
		plugin = strings.TrimSpace(plugin)
		if plugin != "" {
			key := strings.ToLower(plugin)
			if !seen[key] {
				seen[key] = true
				normalized = append(normalized, plugin)
			}
		}
	}

	return normalized
}

// MergeConfigs merges two configurations, with override taking precedence.
func (cp *ConfigProcessor) MergeConfigs(base, override *settingsModels.DashConfig) *settingsModels.DashConfig {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := base.Clone()

	// Override non-empty values
	if strings.TrimSpace(override.Port) != "" {
		merged.Port = override.Port
	}
	if strings.TrimSpace(override.Origin) != "" {
		merged.Origin = override.Origin
	}
	if strings.TrimSpace(override.Front) != "" {
		merged.Front = override.Front
	}
	if len(override.Headers) > 0 {
		merged.Headers = append([]string(nil), override.Headers...)
	}
	if override.Plugins.Count() > 0 {
		merged.Plugins = append(settingsModels.Plugins(nil), override.Plugins...)
	}

	return merged
}
