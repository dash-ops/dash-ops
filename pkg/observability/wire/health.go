package wire

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// HealthRequest represents a request for health check
type HealthRequest struct {
	// Empty for basic health check
}

// HealthResponse represents the response for health check
type HealthResponse struct {
	BaseResponse
	Data HealthData `json:"data"`
}

// HealthData represents the data portion of health response
type HealthData struct {
	Status     string                     `json:"status"`
	Version    string                     `json:"version"`
	Uptime     time.Duration              `json:"uptime"`
	Components map[string]ComponentHealth `json:"components"`
	LastCheck  time.Time                  `json:"last_check"`
}

// ComponentHealth represents the health status of a component
type ComponentHealth struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	LastCheck time.Time              `json:"last_check"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// CacheStatsRequest represents a request for cache statistics
type CacheStatsRequest struct {
	// Empty for basic cache stats
}

// CacheStatsResponse represents the response for cache statistics
type CacheStatsResponse struct {
	BaseResponse
	Data models.CacheStats `json:"data"`
}
