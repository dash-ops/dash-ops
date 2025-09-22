package wire

import (
	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// ConfigurationRequest represents a request for configuration
type ConfigurationRequest struct {
	// Empty for GET requests, contains config for PUT requests
}

// ConfigurationResponse represents the response for configuration
type ConfigurationResponse struct {
	BaseResponse
	Data models.ObservabilityConfig `json:"data"`
}

// ServiceConfigurationRequest represents a request for service configuration
type ServiceConfigurationRequest struct {
	ServiceName string `json:"service_name"`
}

// ServiceConfigurationResponse represents the response for service configuration
type ServiceConfigurationResponse struct {
	BaseResponse
	Data models.ServiceObservabilityConfig `json:"data"`
}

// NotificationChannelsRequest represents a request for notification channels
type NotificationChannelsRequest struct {
	Type    string `json:"type,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
}

// NotificationChannelsResponse represents the response for notification channels
type NotificationChannelsResponse struct {
	BaseResponse
	Data NotificationChannelsData `json:"data"`
}

// NotificationChannelsData represents the data portion of notification channels response
type NotificationChannelsData struct {
	Channels []models.NotificationChannel `json:"channels"`
	Total    int                          `json:"total"`
}

// NotificationChannelRequest represents a request for a single notification channel operation
type NotificationChannelRequest struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"name"`
	Type    string                 `json:"type"` // email, slack, webhook, etc.
	Config  map[string]interface{} `json:"config"`
	Enabled bool                   `json:"enabled"`
}

// NotificationChannelResponse represents the response for a single notification channel operation
type NotificationChannelResponse struct {
	BaseResponse
	Data models.NotificationChannel `json:"data"`
}
