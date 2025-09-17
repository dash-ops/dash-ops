package wire

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// LogsRequest represents a request for log data
type LogsRequest struct {
	Service   string    `json:"service,omitempty"`
	Level     string    `json:"level,omitempty"` // error, warn, info, debug
	Query     string    `json:"query,omitempty"` // Loki query syntax
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Limit     int       `json:"limit,omitempty"` // default: 100
	Offset    int       `json:"offset,omitempty"`
	Stream    bool      `json:"stream,omitempty"` // real-time streaming
	Sort      string    `json:"sort,omitempty"`   // timestamp, level, service
	Order     string    `json:"order,omitempty"`  // asc, desc
}

// LogsResponse moved to responses.go

// MetricsRequest represents a request for metric data
type MetricsRequest struct {
	Service     string            `json:"service,omitempty"`
	Metric      string            `json:"metric,omitempty"` // cpu, memory, requests, etc.
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Step        string            `json:"step,omitempty"`        // 1m, 5m, 1h, etc.
	Aggregation string            `json:"aggregation,omitempty"` // sum, avg, max, min
	Labels      map[string]string `json:"labels,omitempty"`
	GroupBy     []string          `json:"group_by,omitempty"`
}

// MetricsResponse moved to responses.go

// PrometheusQueryRequest represents a Prometheus query request
type PrometheusQueryRequest struct {
	Query     string        `json:"query"` // PromQL query
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Step      string        `json:"step,omitempty"`
	Timeout   time.Duration `json:"timeout,omitempty"`
}

// TracesRequest represents a request for trace data
type TracesRequest struct {
	Service     string    `json:"service,omitempty"`
	Operation   string    `json:"operation,omitempty"`
	TraceID     string    `json:"trace_id,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status,omitempty"`       // ok, error
	MinDuration string    `json:"min_duration,omitempty"` // 100ms, 1s, etc.
	MaxDuration string    `json:"max_duration,omitempty"`
	Limit       int       `json:"limit,omitempty"`
	Sort        string    `json:"sort,omitempty"`  // timestamp, duration
	Order       string    `json:"order,omitempty"` // asc, desc
}

// TracesResponse moved to responses.go

// TraceDetailRequest represents a request for detailed trace information
type TraceDetailRequest struct {
	TraceID string `json:"trace_id"`
}

// TraceDetailResponse moved to responses.go

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

// AlertsResponse moved to responses.go

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

// DashboardsRequest represents a request for dashboard data
type DashboardsRequest struct {
	Service string `json:"service,omitempty"`
	Owner   string `json:"owner,omitempty"`
	Public  *bool  `json:"public,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
}

// DashboardsResponse moved to responses.go

// CreateDashboardRequest represents a request to create a dashboard
type CreateDashboardRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Service     string         `json:"service,omitempty"`
	Charts      []models.Chart `json:"charts"`
	Public      bool           `json:"public"`
}

// UpdateDashboardRequest represents a request to update a dashboard
type UpdateDashboardRequest struct {
	CreateDashboardRequest
}

// DeleteDashboardRequest represents a request to delete a dashboard
type DeleteDashboardRequest struct {
	ID string `json:"id"`
}

// GetDashboardRequest represents a request to get a specific dashboard
type GetDashboardRequest struct {
	ID string `json:"id"`
}

// DashboardTemplatesRequest represents a request for dashboard templates
type DashboardTemplatesRequest struct {
	Category string   `json:"category,omitempty"`
	Service  string   `json:"service,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// DashboardTemplatesResponse moved to responses.go

// ServiceContextRequest represents a request for service context
type ServiceContextRequest struct {
	ServiceName string `json:"service_name"`
}

// ServiceContextResponse moved to responses.go

// ServicesWithContextRequest represents a request for services with context
type ServicesWithContextRequest struct {
	IncludeHealth bool `json:"include_health,omitempty"`
	IncludeStats  bool `json:"include_stats,omitempty"`
}

// ServicesWithContextResponse moved to responses.go

// ServiceHealthRequest represents a request for service health
type ServiceHealthRequest struct {
	ServiceName string `json:"service_name"`
}

// ServiceHealthResponse moved to responses.go

// LogFilters represents filters for log queries
type LogFilters struct {
	Levels    []string          `json:"levels,omitempty"`
	Services  []string          `json:"services,omitempty"`
	Hosts     []string          `json:"hosts,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	TimeRange string            `json:"time_range,omitempty"`
	Query     string            `json:"query,omitempty"`
}

// LogStatsRequest represents a request for log statistics
type LogStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
}

// LogStatistics represents log statistics
type LogStatistics struct {
	TotalLogs     int64                  `json:"total_logs"`
	LogsByLevel   map[string]int64       `json:"logs_by_level"`
	LogsByService map[string]int64       `json:"logs_by_service"`
	LogsByHost    map[string]int64       `json:"logs_by_host"`
	TimeSeries    []TimeSeriesData       `json:"time_series,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// MetricStatsRequest represents a request for metric statistics
type MetricStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	Metric    string    `json:"metric,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
}

// MetricStatistics represents metric statistics
type MetricStatistics struct {
	TotalMetrics     int64                  `json:"total_metrics"`
	MetricsByType    map[string]int64       `json:"metrics_by_type"`
	MetricsByService map[string]int64       `json:"metrics_by_service"`
	TimeSeries       []TimeSeriesData       `json:"time_series,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// TraceStatsRequest represents a request for trace statistics
type TraceStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
}

// TraceStatistics represents trace statistics
type TraceStatistics struct {
	TotalTraces     int64                  `json:"total_traces"`
	TracesByStatus  map[string]int64       `json:"traces_by_status"`
	TracesByService map[string]int64       `json:"traces_by_service"`
	AvgDuration     float64                `json:"avg_duration"`
	MaxDuration     float64                `json:"max_duration"`
	MinDuration     float64                `json:"min_duration"`
	TimeSeries      []TimeSeriesData       `json:"time_series,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// AlertStatsRequest represents a request for alert statistics
type AlertStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
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

// DashboardStatsRequest represents a request for dashboard statistics
type DashboardStatsRequest struct {
	Owner string `json:"owner,omitempty"`
}

// DashboardStatistics represents dashboard statistics
type DashboardStatistics struct {
	TotalDashboards     int64                  `json:"total_dashboards"`
	DashboardsByOwner   map[string]int64       `json:"dashboards_by_owner"`
	DashboardsByService map[string]int64       `json:"dashboards_by_service"`
	PublicDashboards    int64                  `json:"public_dashboards"`
	PrivateDashboards   int64                  `json:"private_dashboards"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// TimeSeriesData represents time series data for statistics
type TimeSeriesData struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Label     string    `json:"label,omitempty"`
}

// MetricSeries represents a metric time series
type MetricSeries struct {
	Metric   string                 `json:"metric"`
	Labels   map[string]string      `json:"labels"`
	Values   []models.MetricData    `json:"values"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TraceAnalysis represents trace analysis results
type TraceAnalysis struct {
	TraceID         string                 `json:"trace_id"`
	TotalDuration   int64                  `json:"total_duration"`
	CriticalPath    []string               `json:"critical_path"`
	Bottlenecks     []string               `json:"bottlenecks"`
	SlowestSpans    []string               `json:"slowest_spans"`
	ErrorRate       float64                `json:"error_rate"`
	Throughput      float64                `json:"throughput"`
	Recommendations []string               `json:"recommendations,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}
