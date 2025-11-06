package wire

import (
	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// ServiceContextRequest represents a request for service context
type ServiceContextRequest struct {
	ServiceName string `json:"service_name"`
}

// ServiceContextResponse represents the response for service context
type ServiceContextResponse struct {
	BaseResponse
	Data models.ServiceContext `json:"data"`
}

// ServicesWithContextRequest represents a request for services with context
type ServicesWithContextRequest struct {
	Search        string `json:"search,omitempty"`
	Limit         int    `json:"limit,omitempty"` // default: 50
	Offset        int    `json:"offset,omitempty"`
	IncludeHealth bool   `json:"include_health,omitempty"`
	IncludeStats  bool   `json:"include_stats,omitempty"`
}

// ServicesWithContextResponse represents the response for services with context
type ServicesWithContextResponse struct {
	BaseResponse
	Data ServicesWithContextData `json:"data"`
}

// ServicesWithContextData represents the data portion of services with context response
type ServicesWithContextData struct {
	Services   []models.ServiceWithContext `json:"services"`
	Total      int                         `json:"total"`
	HasMore    bool                        `json:"has_more"`
	NextOffset int                         `json:"next_offset,omitempty"`
	Summary    ServiceSummary              `json:"summary,omitempty"`
}

// ServiceSummary represents a summary of services
type ServiceSummary struct {
	HealthyServices  int   `json:"healthy_services"`
	WarningServices  int   `json:"warning_services"`
	CriticalServices int   `json:"critical_services"`
	TotalLogs        int64 `json:"total_logs"`
	TotalMetrics     int64 `json:"total_metrics"`
	TotalTraces      int64 `json:"total_traces"`
	TotalAlerts      int64 `json:"total_alerts"`
}

// ServiceHealthRequest represents a request for service health
type ServiceHealthRequest struct {
	ServiceName string `json:"service_name"`
}

// ServiceHealthResponse represents the response for service health
type ServiceHealthResponse struct {
	BaseResponse
	Data models.ServiceHealth `json:"data"`
}
