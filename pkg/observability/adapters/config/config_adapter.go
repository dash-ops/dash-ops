package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// ObservabilityConfig represents the observability configuration structure
type ObservabilityConfig struct {
	Loki         ServiceConfig `yaml:"loki"`
	Prometheus   ServiceConfig `yaml:"prometheus"`
	Tempo        ServiceConfig `yaml:"tempo"`
	AlertManager ServiceConfig `yaml:"alertmanager"`
}

// ServiceConfig represents configuration for an observability service
type ServiceConfig struct {
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
	Enabled bool   `yaml:"enabled"`
}

// ConfigAdapter handles observability configuration parsing
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ParseObservabilityConfigFromFileConfig parses observability config from file bytes
func (ca *ConfigAdapter) ParseObservabilityConfigFromFileConfig(fileConfig []byte) (*ObservabilityConfig, error) {
	var config struct {
		Observability ObservabilityConfig `yaml:"observability"`
	}

	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to parse observability configuration: %w", err)
	}

	// Set default values if not provided
	obsConfig := &config.Observability

	if obsConfig.Loki.Timeout == 0 {
		obsConfig.Loki.Timeout = 30
	}
	if obsConfig.Prometheus.Timeout == 0 {
		obsConfig.Prometheus.Timeout = 30
	}
	if obsConfig.Tempo.Timeout == 0 {
		obsConfig.Tempo.Timeout = 30
	}
	if obsConfig.AlertManager.Timeout == 0 {
		obsConfig.AlertManager.Timeout = 30
	}

	return obsConfig, nil
}
