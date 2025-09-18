package models

import (
	"time"
)

// ObservabilityConfig represents the global observability configuration
type ObservabilityConfig struct {
	Logs    LogsConfig    `json:"logs"`
	Metrics MetricsConfig `json:"metrics"`
	Traces  TracesConfig  `json:"traces"`
	Alerts  AlertsConfig  `json:"alerts"`
	Cache   CacheConfig   `json:"cache"`
	UI      UIConfig      `json:"ui"`
}

// ServiceObservabilityConfig represents observability configuration for a specific service
type ServiceObservabilityConfig struct {
	ServiceName string                 `json:"service_name"`
	Logs        *LogsConfig            `json:"logs,omitempty"`
	Metrics     *MetricsConfig         `json:"metrics,omitempty"`
	Traces      *TracesConfig          `json:"traces,omitempty"`
	Alerts      *AlertsConfig          `json:"alerts,omitempty"`
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
