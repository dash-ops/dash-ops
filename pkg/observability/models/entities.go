package models

import (
	"time"
)

// LogEntry represents a log entry from any log source
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
	Source    string                 `json:"source"`
	Host      string                 `json:"host"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Labels    map[string]string      `json:"labels,omitempty"`
	Raw       string                 `json:"raw,omitempty"`
}

// ProcessedLogEntry represents a processed log entry with additional context
type ProcessedLogEntry struct {
	LogEntry
	ProcessedAt  time.Time              `json:"processed_at"`
	Enrichments  map[string]interface{} `json:"enrichments,omitempty"`
	Correlations []string               `json:"correlations,omitempty"`
}

// MetricData represents a metric data point
type MetricData struct {
	Timestamp time.Time              `json:"timestamp"`
	Value     float64                `json:"value"`
	Labels    map[string]string      `json:"labels,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Metric    string                 `json:"metric"`
	Type      string                 `json:"type"` // counter, gauge, histogram, summary
}

// ProcessedMetric represents a processed metric with additional context
type ProcessedMetric struct {
	MetricData
	ProcessedAt time.Time              `json:"processed_at"`
	Enrichments map[string]interface{} `json:"enrichments,omitempty"`
	Derived     bool                   `json:"derived"`
}

// DerivedMetric represents a derived metric calculated from base metrics
type DerivedMetric struct {
	MetricData
	Formula      string    `json:"formula"`
	BaseMetrics  []string  `json:"base_metrics"`
	CalculatedAt time.Time `json:"calculated_at"`
}

// TraceSpan represents a span in a distributed trace
type TraceSpan struct {
	ID            string                 `json:"id"`
	TraceID       string                 `json:"trace_id"`
	OperationName string                 `json:"operation_name"`
	Service       string                 `json:"service"`
	StartTime     int64                  `json:"start_time"` // microseconds
	Duration      int64                  `json:"duration"`   // microseconds
	Status        string                 `json:"status"`
	Tags          map[string]interface{} `json:"tags"`
	ParentID      string                 `json:"parent_id,omitempty"`
	Depth         int                    `json:"depth"`
	Logs          []LogEntry             `json:"logs,omitempty"`
}

// ProcessedTrace represents a processed trace with additional context
type ProcessedTrace struct {
	TraceSpan
	ProcessedAt  time.Time              `json:"processed_at"`
	Enrichments  map[string]interface{} `json:"enrichments,omitempty"`
	Correlations []string               `json:"correlations,omitempty"`
	Performance  *TracePerformance      `json:"performance,omitempty"`
}

// TraceInfo represents a summarized view of a trace used in listings
type TraceInfo struct {
	TraceID       string    `json:"trace_id"`
	RootOperation string    `json:"root_operation"`
	TotalDuration int64     `json:"total_duration"`
	SpanCount     int       `json:"span_count"`
	ServiceCount  int       `json:"service_count"`
	Status        string    `json:"status"` // ok, error
	Timestamp     time.Time `json:"timestamp"`
	Errors        int       `json:"errors"`
	Services      []string  `json:"services"`
}

// TracePerformance represents performance characteristics of a trace
type TracePerformance struct {
	TotalDuration int64    `json:"total_duration"`
	CriticalPath  []string `json:"critical_path"`
	Bottlenecks   []string `json:"bottlenecks"`
	SlowestSpans  []string `json:"slowest_spans"`
	ErrorRate     float64  `json:"error_rate"`
	Throughput    float64  `json:"throughput"`
}

// Alert represents an alert instance
type Alert struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	Severity     string                 `json:"severity"`
	Service      string                 `json:"service"`
	Labels       map[string]string      `json:"labels"`
	Annotations  map[string]string      `json:"annotations"`
	StartsAt     time.Time              `json:"starts_at"`
	EndsAt       *time.Time             `json:"ends_at,omitempty"`
	GeneratorURL string                 `json:"generator_url,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	RuleID       string                 `json:"rule_id,omitempty"`
}

// ProcessedAlert represents a processed alert with additional context
type ProcessedAlert struct {
	Alert
	ProcessedAt  time.Time              `json:"processed_at"`
	Enrichments  map[string]interface{} `json:"enrichments,omitempty"`
	Correlations []string               `json:"correlations,omitempty"`
	Actions      []AlertAction          `json:"actions,omitempty"`
}

// AlertAction represents an action taken on an alert
type AlertAction struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	User        string                 `json:"user,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AlertRule represents an alert rule configuration
type AlertRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Query       string            `json:"query"`
	Threshold   float64           `json:"threshold"`
	Severity    string            `json:"severity"`
	Service     string            `json:"service,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Enabled     bool              `json:"enabled"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AlertEvaluation represents the evaluation of an alert rule
type AlertEvaluation struct {
	RuleID      string    `json:"rule_id"`
	EvaluatedAt time.Time `json:"evaluated_at"`
	Result      bool      `json:"result"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Message     string    `json:"message,omitempty"`
}

// Dashboard represents a dashboard configuration
type Dashboard struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Service     string    `json:"service,omitempty"`
	Charts      []Chart   `json:"charts"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Owner       string    `json:"owner,omitempty"`
	Public      bool      `json:"public"`
}

// Chart represents a chart configuration in a dashboard
type Chart struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Type         string                 `json:"type"` // line, area, bar, pie, table
	Metrics      []string               `json:"metrics"`
	ServiceScope string                 `json:"service_scope"` // all, specific
	TimeRange    string                 `json:"time_range"`
	Height       int                    `json:"height"`
	Config       map[string]interface{} `json:"config,omitempty"`
	Position     *ChartPosition         `json:"position,omitempty"`
}

// ChartPosition represents the position of a chart in a dashboard
type ChartPosition struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// DashboardTemplate represents a dashboard template
type DashboardTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Service     string   `json:"service,omitempty"`
	Charts      []Chart  `json:"charts"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags,omitempty"`
}

// DashboardData represents processed data for dashboard visualization
type DashboardData struct {
	DashboardID string                 `json:"dashboard_id"`
	Charts      []ChartData            `json:"charts"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	GeneratedAt time.Time              `json:"generated_at"`
}

// ChartData represents processed data for a specific chart
type ChartData struct {
	ChartID string                 `json:"chart_id"`
	Data    []interface{}          `json:"data"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

// ServiceContext represents service context for observability
type ServiceContext struct {
	ServiceName string                 `json:"service_name"`
	Namespace   string                 `json:"namespace"`
	Cluster     string                 `json:"cluster"`
	Labels      map[string]string      `json:"labels"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Health      *ServiceHealth         `json:"health,omitempty"`
}

// ServiceWithContext represents a service with its observability context
type ServiceWithContext struct {
	ServiceContext
	LogCount    int64 `json:"log_count"`
	MetricCount int64 `json:"metric_count"`
	TraceCount  int64 `json:"trace_count"`
	AlertCount  int64 `json:"alert_count"`
}

// ServiceHealth represents the health status of a service
type ServiceHealth struct {
	Status    string                 `json:"status"`
	LastCheck time.Time              `json:"last_check"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Metrics   map[string]float64     `json:"metrics,omitempty"`
	Alerts    []string               `json:"alerts,omitempty"`
}

// NotificationChannel represents a notification channel configuration
type NotificationChannel struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"` // email, slack, webhook, etc.
	Config    map[string]interface{} `json:"config"`
	Enabled   bool                   `json:"enabled"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

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

// LogsConfig represents logs configuration
type LogsConfig struct {
	Enabled     bool     `json:"enabled"`
	Retention   string   `json:"retention"`
	Levels      []string `json:"levels"`
	Sources     []string `json:"sources"`
	QueryLimit  int      `json:"query_limit"`
	StreamLimit int      `json:"stream_limit"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled     bool     `json:"enabled"`
	Retention   string   `json:"retention"`
	Metrics     []string `json:"metrics"`
	Aggregation string   `json:"aggregation"`
	Step        string   `json:"step"`
}

// TracesConfig represents traces configuration
type TracesConfig struct {
	Enabled    bool     `json:"enabled"`
	Retention  string   `json:"retention"`
	Services   []string `json:"services"`
	Operations []string `json:"operations"`
	MaxDepth   int      `json:"max_depth"`
}

// AlertsConfig represents alerts configuration
type AlertsConfig struct {
	Enabled    bool     `json:"enabled"`
	Channels   []string `json:"channels"`
	Severities []string `json:"severities"`
	Cooldown   string   `json:"cooldown"`
	MaxAlerts  int      `json:"max_alerts"`
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
