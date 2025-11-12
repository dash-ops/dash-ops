package logic

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/dash-ops/dash-ops/pkg/settings/models"
)

// YAMLProcessor handles YAML generation and validation
type YAMLProcessor struct{}

// NewYAMLProcessor creates a new YAML processor
func NewYAMLProcessor() *YAMLProcessor {
	return &YAMLProcessor{}
}

// GenerateYAML generates YAML from configuration data
func (yp *YAMLProcessor) GenerateYAML(config *models.DashConfig) ([]byte, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Apply defaults
	yp.applyDefaults(config)

	// Validate
	if err := yp.validateConfig(config); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return yamlData, nil
}

// ParseYAML parses YAML into configuration data
func (yp *YAMLProcessor) ParseYAML(data []byte) (*models.DashConfig, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("YAML data cannot be empty")
	}

	var config models.DashConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Apply defaults
	yp.applyDefaults(&config)

	// Validate
	if err := yp.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &config, nil
}

// applyDefaults applies default values to configuration
func (yp *YAMLProcessor) applyDefaults(config *models.DashConfig) {
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
}

// validateConfig validates configuration data
func (yp *YAMLProcessor) validateConfig(config *models.DashConfig) error {
	return config.Validate()
}
