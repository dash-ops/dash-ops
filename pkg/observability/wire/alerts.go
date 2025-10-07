package wire

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// AlertsRequest represents a request for alert data
type AlertsRequest struct {
	Service   string    `json:"service,omitempty"`
	Status    string    `json:"status,omitempty"`   // active, resolved, silenced
	Severity  string    `json:"severity,omitempty"` // critical, warning, info
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

// AlertsResponse represents the response for alert data
type AlertsResponse struct {
	BaseResponse
	Data AlertsData `json:"data"`
}

// AlertsData represents the data portion of alerts response
type AlertsData struct {
	Alerts  []models.Alert `json:"alerts"`
	Total   int            `json:"total"`
	Filters AlertFilters   `json:"filters,omitempty"`
}

// AlertFilters represents filters for alert queries
type AlertFilters struct {
	Statuses   []string `json:"statuses,omitempty"`
	Severities []string `json:"severities,omitempty"`
	Services   []string `json:"services,omitempty"`
}

// CreateAlertRequest represents a request to create an alert
type CreateAlertRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Query       string            `json:"query"` // PromQL or LogQL
	Threshold   float64           `json:"threshold"`
	Severity    string            `json:"severity"`
	Service     string            `json:"service,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Enabled     bool              `json:"enabled"`
}

// UpdateAlertRequest represents a request to update an alert
type UpdateAlertRequest struct {
	CreateAlertRequest
}

// DeleteAlertRequest represents a request to delete an alert
type DeleteAlertRequest struct {
	ID string `json:"id"`
}

// SilenceAlertRequest represents a request to silence an alert
type SilenceAlertRequest struct {
	ID       string        `json:"id"`
	Duration time.Duration `json:"duration"`
	Reason   string        `json:"reason,omitempty"`
}

// AlertResponse represents the response for a single alert operation
type AlertResponse struct {
	BaseResponse
	Data models.Alert `json:"data"`
}

// AlertRulesRequest represents a request for alert rules
type AlertRulesRequest struct {
	Service string `json:"service,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

// AlertRulesResponse represents the response for alert rules
type AlertRulesResponse struct {
	BaseResponse
	Data AlertRulesData `json:"data"`
}

// AlertRulesData represents the data portion of alert rules response
type AlertRulesData struct {
	Rules []models.AlertRule `json:"rules"`
	Total int                `json:"total"`
}

// AlertStatsRequest represents a request for alert statistics
type AlertStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
}

// AlertStatisticsResponse represents the response for alert statistics
type AlertStatisticsResponse struct {
	BaseResponse
	Data AlertStatistics `json:"data"`
}

// AlertStatistics represents alert statistics
type AlertStatistics struct {
	TotalAlerts      int64                  `json:"total_alerts"`
	AlertsByStatus   map[string]int64       `json:"alerts_by_status"`
	AlertsBySeverity map[string]int64       `json:"alerts_by_severity"`
	AlertsByService  map[string]int64       `json:"alerts_by_service"`
	TimeSeries       []TimeSeriesData       `json:"time_series,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}
