package wire

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// BaseResponse represents a base response structure
type BaseResponse struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// LogsResponse represents the response for log data
type LogsResponse struct {
	BaseResponse
	Data LogsData `json:"data"`
}

// LogsData represents the data portion of logs response
type LogsData struct {
	Logs       []models.LogEntry `json:"logs"`
	Total      int               `json:"total"`
	HasMore    bool              `json:"has_more"`
	NextOffset int               `json:"next_offset,omitempty"`
	Filters    LogFilters        `json:"filters,omitempty"`
}

// MetricsResponse represents the response for metric data
type MetricsResponse struct {
	BaseResponse
	Data MetricsData `json:"data"`
}

// MetricsData represents the data portion of metrics response
type MetricsData struct {
	Metrics []models.MetricData `json:"metrics"`
	Labels  []string            `json:"labels,omitempty"`
	Query   string              `json:"query,omitempty"`
	Step    string              `json:"step,omitempty"`
}

// TracesResponse represents the response for trace data
type TracesResponse struct {
	BaseResponse
	Data TracesData `json:"data"`
}

// TracesData represents the data portion of traces response
type TracesData struct {
	Traces []models.TraceInfo `json:"traces"`
	Total  int                `json:"total"`
	Query  string             `json:"query,omitempty"`
}

// TraceDetailResponse represents the response for detailed trace information
type TraceDetailResponse struct {
	BaseResponse
	Data TraceDetailData `json:"data"`
}

// TraceDetailData represents the data portion of trace detail response
type TraceDetailData struct {
	TraceID  string             `json:"trace_id"`
	Spans    []models.TraceSpan `json:"spans"`
	Total    int                `json:"total"`
	Timeline TraceTimeline      `json:"timeline,omitempty"`
}

// TraceTimeline represents timeline information for a trace
type TraceTimeline struct {
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	Duration  int64    `json:"duration"`
	Services  []string `json:"services"`
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

// AlertResponse represents the response for a single alert operation
type AlertResponse struct {
	BaseResponse
	Data models.Alert `json:"data"`
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

// DashboardsResponse represents the response for dashboard data
type DashboardsResponse struct {
	BaseResponse
	Data DashboardsData `json:"data"`
}

// DashboardsData represents the data portion of dashboards response
type DashboardsData struct {
	Dashboards []models.Dashboard `json:"dashboards"`
	Total      int                `json:"total"`
	Filters    DashboardFilters   `json:"filters,omitempty"`
}

// DashboardFilters represents filters for dashboard queries
type DashboardFilters struct {
	Owners   []string `json:"owners,omitempty"`
	Services []string `json:"services,omitempty"`
	Public   *bool    `json:"public,omitempty"`
}

// DashboardResponse represents the response for a single dashboard operation
type DashboardResponse struct {
	BaseResponse
	Data models.Dashboard `json:"data"`
}

// DashboardTemplatesResponse represents the response for dashboard templates
type DashboardTemplatesResponse struct {
	BaseResponse
	Data DashboardTemplatesData `json:"data"`
}

// DashboardTemplatesData represents the data portion of dashboard templates response
type DashboardTemplatesData struct {
	Templates  []models.DashboardTemplate `json:"templates"`
	Total      int                        `json:"total"`
	Categories []string                   `json:"categories,omitempty"`
	Tags       []string                   `json:"tags,omitempty"`
}

// ServiceContextResponse represents the response for service context
type ServiceContextResponse struct {
	BaseResponse
	Data models.ServiceContext `json:"data"`
}

// ServicesWithContextResponse represents the response for services with context
type ServicesWithContextResponse struct {
	BaseResponse
	Data ServicesWithContextData `json:"data"`
}

// ServicesWithContextData represents the data portion of services with context response
type ServicesWithContextData struct {
	Services []models.ServiceWithContext `json:"services"`
	Total    int                         `json:"total"`
	Summary  ServiceSummary              `json:"summary,omitempty"`
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

// ServiceHealthResponse represents the response for service health
type ServiceHealthResponse struct {
	BaseResponse
	Data models.ServiceHealth `json:"data"`
}

// LogStatisticsResponse represents the response for log statistics
type LogStatisticsResponse struct {
	BaseResponse
	Data LogStatistics `json:"data"`
}

// MetricStatisticsResponse represents the response for metric statistics
type MetricStatisticsResponse struct {
	BaseResponse
	Data MetricStatistics `json:"data"`
}

// TraceStatisticsResponse represents the response for trace statistics
type TraceStatisticsResponse struct {
	BaseResponse
	Data TraceStatistics `json:"data"`
}

// AlertStatisticsResponse represents the response for alert statistics
type AlertStatisticsResponse struct {
	BaseResponse
	Data AlertStatistics `json:"data"`
}

// DashboardStatisticsResponse represents the response for dashboard statistics
type DashboardStatisticsResponse struct {
	BaseResponse
	Data DashboardStatistics `json:"data"`
}

// TraceAnalysisResponse represents the response for trace analysis
type TraceAnalysisResponse struct {
	BaseResponse
	Data TraceAnalysis `json:"data"`
}

// ConfigurationResponse represents the response for configuration
type ConfigurationResponse struct {
	BaseResponse
	Data models.ObservabilityConfig `json:"data"`
}

// ServiceConfigurationResponse represents the response for service configuration
type ServiceConfigurationResponse struct {
	BaseResponse
	Data models.ServiceObservabilityConfig `json:"data"`
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

// NotificationChannelResponse represents the response for a single notification channel operation
type NotificationChannelResponse struct {
	BaseResponse
	Data models.NotificationChannel `json:"data"`
}

// CacheStatsResponse represents the response for cache statistics
type CacheStatsResponse struct {
	BaseResponse
	Data models.CacheStats `json:"data"`
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	BaseResponse
	ErrorCode    string                 `json:"error_code,omitempty"`
	ErrorDetails map[string]interface{} `json:"error_details,omitempty"`
	RequestID    string                 `json:"request_id,omitempty"`
	Timestamp    time.Time              `json:"timestamp"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// QueryInfo represents query information
type QueryInfo struct {
	Query     string                 `json:"query"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	TimeRange TimeRange              `json:"time_range"`
	Limit     int                    `json:"limit,omitempty"`
	Offset    int                    `json:"offset,omitempty"`
}

// PerformanceInfo represents performance information
type PerformanceInfo struct {
	QueryTime   time.Duration `json:"query_time"`
	ProcessTime time.Duration `json:"process_time"`
	TotalTime   time.Duration `json:"total_time"`
	CacheHit    bool          `json:"cache_hit"`
	CacheTime   time.Duration `json:"cache_time,omitempty"`
	ResultSize  int           `json:"result_size"`
}
