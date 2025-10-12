package models

import (
	"time"
)

// ObservabilityConfig represents the global observability configuration
type ObservabilityConfig struct {
	Enabled  bool                    `yaml:"enabled" json:"enabled"`
	Logs     LogsProviderConfig      `yaml:"logs" json:"logs"`
	Metrics  MetricsProviderConfig   `yaml:"metrics" json:"metrics"`
	Traces   TracesProviderConfig    `yaml:"traces" json:"traces"`
	Alerts   AlertsProviderConfig    `yaml:"alerts" json:"alerts"`
	Cache    CacheConfig             `yaml:"cache" json:"cache"`
	UI       UIConfig                `yaml:"ui" json:"ui"`
	Services []ServiceOverrideConfig `yaml:"services,omitempty" json:"services,omitempty"`
}

// LogsProviderConfig represents logs provider configuration
type LogsProviderConfig struct {
	Providers   []ProviderConfig `yaml:"providers" json:"providers"`
	Retention   string           `yaml:"retention" json:"retention"`
	QueryLimit  int              `yaml:"query_limit" json:"query_limit"`
	StreamLimit int              `yaml:"stream_limit" json:"stream_limit"`
	Levels      []string         `yaml:"levels" json:"levels"`
}

// MetricsProviderConfig represents metrics provider configuration
type MetricsProviderConfig struct {
	Providers      []ProviderConfig `yaml:"providers" json:"providers"`
	Retention      string           `yaml:"retention" json:"retention"`
	QueryLimit     int              `yaml:"query_limit" json:"query_limit"`
	ScrapeInterval string           `yaml:"scrape_interval" json:"scrape_interval"`
}

// TracesProviderConfig represents traces provider configuration
type TracesProviderConfig struct {
	Providers    []ProviderConfig `yaml:"providers" json:"providers"`
	Retention    string           `yaml:"retention" json:"retention"`
	QueryLimit   int              `yaml:"query_limit" json:"query_limit"`
	SamplingRate float64          `yaml:"sampling_rate" json:"sampling_rate"`
}

// AlertsProviderConfig represents alerts provider configuration
type AlertsProviderConfig struct {
	Enabled   bool             `yaml:"enabled" json:"enabled"`
	Providers []ProviderConfig `yaml:"providers" json:"providers"`
}

// ProviderConfig represents a generic provider configuration
type ProviderConfig struct {
	Name     string                 `yaml:"name" json:"name"`
	Type     string                 `yaml:"type" json:"type"` // loki, tempo, prometheus, datadog, etc.
	URL      string                 `yaml:"url" json:"url"`
	Timeout  string                 `yaml:"timeout" json:"timeout"`
	Auth     AuthConfig             `yaml:"auth" json:"auth"`
	Enabled  bool                   `yaml:"enabled" json:"enabled"`
	Labels   map[string]string      `yaml:"labels,omitempty" json:"labels,omitempty"`
	Metadata map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type     string `yaml:"type" json:"type"` // basic, bearer, none
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
	Token    string `yaml:"token,omitempty" json:"token,omitempty"`
}

// ServiceOverrideConfig represents service-specific configuration overrides
type ServiceOverrideConfig struct {
	Name    string                 `yaml:"name" json:"name"`
	Logs    *LogsOverrideConfig    `yaml:"logs,omitempty" json:"logs,omitempty"`
	Traces  *TracesOverrideConfig  `yaml:"traces,omitempty" json:"traces,omitempty"`
	Metrics *MetricsOverrideConfig `yaml:"metrics,omitempty" json:"metrics,omitempty"`
}

// LogsOverrideConfig represents logs override configuration
type LogsOverrideConfig struct {
	Levels []string `yaml:"levels,omitempty" json:"levels,omitempty"`
}

// TracesOverrideConfig represents traces override configuration
type TracesOverrideConfig struct {
	SamplingRate *float64 `yaml:"sampling_rate,omitempty" json:"sampling_rate,omitempty"`
}

// MetricsOverrideConfig represents metrics override configuration
type MetricsOverrideConfig struct {
	ScrapeInterval *string `yaml:"scrape_interval,omitempty" json:"scrape_interval,omitempty"`
}

// ServiceObservabilityConfig represents observability configuration for a specific service
type ServiceObservabilityConfig struct {
	ServiceName string                 `json:"service_name"`
	Logs        *LogsProviderConfig    `json:"logs,omitempty"`
	Metrics     *MetricsProviderConfig `json:"metrics,omitempty"`
	Traces      *TracesProviderConfig  `json:"traces,omitempty"`
	Alerts      *AlertsProviderConfig  `json:"alerts,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	Enabled bool          `json:"enabled"`
	TTL     time.Duration `json:"ttl"`
	MaxSize int64         `json:"max_size"`
	Cleanup time.Duration `json:"cleanup"`
}

// UIConfig represents UI configuration
type UIConfig struct {
	Theme       string                 `json:"theme"`
	RefreshRate time.Duration          `json:"refresh_rate"`
	PageSize    int                    `json:"page_size"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits      int64   `json:"hits"`
	Misses    int64   `json:"misses"`
	Size      int64   `json:"size"`
	MaxSize   int64   `json:"max_size"`
	HitRate   float64 `json:"hit_rate"`
	Evictions int64   `json:"evictions"`
}
