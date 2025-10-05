package config

import (
	"fmt"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"gopkg.in/yaml.v2"
)

// ConfigAdapter handles observability configuration parsing
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ParseObservabilityConfigFromFileConfig parses observability config from file bytes
func (ca *ConfigAdapter) ParseObservabilityConfigFromFileConfig(fileConfig []byte) (*models.ObservabilityConfig, error) {
	var config struct {
		Observability models.ObservabilityConfig `yaml:"observability"`
	}

	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to parse observability configuration: %w", err)
	}

	obsConfig := &config.Observability

	// Set default values if not provided
	if err := ca.setDefaults(obsConfig); err != nil {
		return nil, fmt.Errorf("failed to set default values: %w", err)
	}

	// Validate configuration
	if err := ca.validate(obsConfig); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return obsConfig, nil
}

// setDefaults sets default values for the configuration
func (ca *ConfigAdapter) setDefaults(config *models.ObservabilityConfig) error {
	// Enable by default if not specified
	if !config.Enabled {
		config.Enabled = true
	}

	// Logs defaults
	if config.Logs.QueryLimit == 0 {
		config.Logs.QueryLimit = 1000
	}
	if config.Logs.StreamLimit == 0 {
		config.Logs.StreamLimit = 100
	}
	if len(config.Logs.Levels) == 0 {
		config.Logs.Levels = []string{"debug", "info", "warn", "error"}
	}

	// Traces defaults
	if config.Traces.QueryLimit == 0 {
		config.Traces.QueryLimit = 100
	}
	if config.Traces.SamplingRate == 0 {
		config.Traces.SamplingRate = 1.0
	}

	// Metrics defaults
	if config.Metrics.QueryLimit == 0 {
		config.Metrics.QueryLimit = 1000
	}

	// Cache defaults
	if config.Cache.TTL == 0 {
		config.Cache.TTL = 5 * time.Minute
	}
	if config.Cache.MaxSize == 0 {
		config.Cache.MaxSize = 104857600 // 100MB
	}
	if config.Cache.Cleanup == 0 {
		config.Cache.Cleanup = 10 * time.Minute
	}

	// UI defaults
	if config.UI.RefreshRate == 0 {
		config.UI.RefreshRate = 30 * time.Second
	}
	if config.UI.PageSize == 0 {
		config.UI.PageSize = 50
	}

	// Set default timeout for providers
	for i := range config.Logs.Providers {
		if config.Logs.Providers[i].Timeout == "" {
			config.Logs.Providers[i].Timeout = "30s"
		}
		if config.Logs.Providers[i].Auth.Type == "" {
			config.Logs.Providers[i].Auth.Type = "none"
		}
	}

	for i := range config.Traces.Providers {
		if config.Traces.Providers[i].Timeout == "" {
			config.Traces.Providers[i].Timeout = "30s"
		}
		if config.Traces.Providers[i].Auth.Type == "" {
			config.Traces.Providers[i].Auth.Type = "none"
		}
	}

	for i := range config.Metrics.Providers {
		if config.Metrics.Providers[i].Timeout == "" {
			config.Metrics.Providers[i].Timeout = "30s"
		}
		if config.Metrics.Providers[i].Auth.Type == "" {
			config.Metrics.Providers[i].Auth.Type = "none"
		}
	}

	return nil
}

// validate validates the configuration
func (ca *ConfigAdapter) validate(config *models.ObservabilityConfig) error {
	if !config.Enabled {
		return nil // Skip validation if observability is disabled
	}

	// Validate at least one provider is configured for each enabled component
	hasLogsProvider := len(config.Logs.Providers) > 0
	hasTracesProvider := len(config.Traces.Providers) > 0
	hasMetricsProvider := len(config.Metrics.Providers) > 0

	if !hasLogsProvider && !hasTracesProvider && !hasMetricsProvider {
		return fmt.Errorf("at least one provider must be configured")
	}

	// Validate provider configurations
	for _, provider := range config.Logs.Providers {
		if err := ca.validateProvider(provider); err != nil {
			return fmt.Errorf("invalid logs provider %s: %w", provider.Name, err)
		}
	}

	for _, provider := range config.Traces.Providers {
		if err := ca.validateProvider(provider); err != nil {
			return fmt.Errorf("invalid traces provider %s: %w", provider.Name, err)
		}
	}

	for _, provider := range config.Metrics.Providers {
		if err := ca.validateProvider(provider); err != nil {
			return fmt.Errorf("invalid metrics provider %s: %w", provider.Name, err)
		}
	}

	return nil
}

// validateProvider validates a single provider configuration
func (ca *ConfigAdapter) validateProvider(provider models.ProviderConfig) error {
	if provider.Name == "" {
		return fmt.Errorf("provider name is required")
	}

	if provider.Type == "" {
		return fmt.Errorf("provider type is required")
	}

	if provider.URL == "" {
		return fmt.Errorf("provider URL is required")
	}

	// Validate timeout format
	if provider.Timeout != "" {
		if _, err := time.ParseDuration(provider.Timeout); err != nil {
			return fmt.Errorf("invalid timeout format: %w", err)
		}
	}

	// Validate auth config
	if provider.Auth.Type != "none" && provider.Auth.Type != "basic" && provider.Auth.Type != "bearer" {
		return fmt.Errorf("invalid auth type: %s (must be 'none', 'basic', or 'bearer')", provider.Auth.Type)
	}

	if provider.Auth.Type == "basic" {
		if provider.Auth.Username == "" || provider.Auth.Password == "" {
			return fmt.Errorf("basic auth requires username and password")
		}
	}

	if provider.Auth.Type == "bearer" {
		if provider.Auth.Token == "" {
			return fmt.Errorf("bearer auth requires token")
		}
	}

	return nil
}
