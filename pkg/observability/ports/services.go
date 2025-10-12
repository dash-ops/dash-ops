package ports

import (
	"context"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LogService defines the interface for log processing business logic
type LogService interface {
	// ProcessLogs processes and enriches log entries
	ProcessLogs(ctx context.Context, logs []models.LogEntry) ([]models.ProcessedLogEntry, error)

	// SearchLogs performs advanced log search with query parsing
	SearchLogs(ctx context.Context, query string, filters wire.LogFilters) (*wire.LogsResponse, error)

	// GetLogStatistics returns statistics about logs
	GetLogStatistics(ctx context.Context, req *wire.LogStatsRequest) (*wire.LogStatistics, error)

	// CorrelateLogsWithTraces correlates logs with trace data
	CorrelateLogsWithTraces(ctx context.Context, logs []models.LogEntry) ([]models.LogEntry, error)
}

// MetricService defines the interface for metrics processing business logic
type MetricService interface {
	// ProcessMetrics processes and aggregates metric data
	ProcessMetrics(ctx context.Context, metrics []models.MetricData) ([]models.ProcessedMetric, error)

	// CalculateDerivedMetrics calculates derived metrics from base metrics
	CalculateDerivedMetrics(ctx context.Context, baseMetrics []models.MetricData) ([]models.DerivedMetric, error)

	// GetMetricStatistics returns statistics about metrics
	GetMetricStatistics(ctx context.Context, req *wire.MetricStatsRequest) (*wire.MetricStatistics, error)

	// CorrelateMetricsWithLogs correlates metrics with log data
	CorrelateMetricsWithLogs(ctx context.Context, metrics []models.MetricData) ([]models.MetricData, error)
}

// TraceService defines the interface for trace processing business logic
type TraceService interface {
	// ProcessTraces processes and enriches trace data
	ProcessTraces(ctx context.Context, traces []models.TraceSpan) ([]models.TraceSpan, error)

	// AnalyzeTracePerformance analyzes trace performance characteristics
	AnalyzeTracePerformance(ctx context.Context, traceID string) (*wire.TraceAnalysis, error)

	// GetTraceStatistics returns statistics about traces
	GetTraceStatistics(ctx context.Context, req *wire.TraceStatsRequest) (*wire.TraceStatistics, error)

	// CorrelateTracesWithLogs correlates traces with log data
	CorrelateTracesWithLogs(ctx context.Context, traces []models.TraceSpan) ([]models.TraceSpan, error)
}

// AlertService defines the interface for alert processing business logic
type AlertService interface {
	// ProcessAlerts processes and evaluates alert conditions
	ProcessAlerts(ctx context.Context, alerts []models.Alert) ([]models.ProcessedAlert, error)

	// EvaluateAlertRules evaluates alert rules against current data
	EvaluateAlertRules(ctx context.Context, rules []models.AlertRule) ([]models.AlertEvaluation, error)

	// GetAlertStatistics returns statistics about alerts
	GetAlertStatistics(ctx context.Context, req *wire.AlertStatsRequest) (*wire.AlertStatistics, error)

	// CorrelateAlertsWithMetrics correlates alerts with metric data
	CorrelateAlertsWithMetrics(ctx context.Context, alerts []models.Alert) ([]models.Alert, error)
}

// DashboardService defines the interface for dashboard business logic
type DashboardService interface {
	// GenerateDashboard generates a dashboard based on service context
	GenerateDashboard(ctx context.Context, serviceName string, template string) (*models.Dashboard, error)

	// ProcessDashboardData processes data for dashboard visualization
	ProcessDashboardData(ctx context.Context, dashboard *models.Dashboard) (*models.DashboardData, error)

	// GetDashboardStatistics returns statistics about dashboards
	GetDashboardStatistics(ctx context.Context, req *wire.DashboardStatsRequest) (*wire.DashboardStatistics, error)

	// ValidateDashboard validates dashboard configuration
	ValidateDashboard(ctx context.Context, dashboard *models.Dashboard) error
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	// SendAlertNotification sends alert notifications
	SendAlertNotification(ctx context.Context, alert *models.Alert, channels []string) error

	// SendDashboardNotification sends dashboard notifications
	SendDashboardNotification(ctx context.Context, dashboard *models.Dashboard, channels []string) error

	// GetNotificationChannels returns available notification channels
	GetNotificationChannels(ctx context.Context) ([]models.NotificationChannel, error)

	// ConfigureNotificationChannel configures a notification channel
	ConfigureNotificationChannel(ctx context.Context, channel *models.NotificationChannel) error
}

// CacheService defines the interface for caching operations
type CacheService interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (interface{}, error)

	// Set stores a value in cache
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Clear clears all cache entries
	Clear(ctx context.Context) error

	// GetStats returns cache statistics
	GetStats(ctx context.Context) (*models.CacheStats, error)
}

// ConfigurationService defines the interface for configuration operations
type ConfigurationService interface {
	// GetObservabilityConfig returns observability configuration
	GetObservabilityConfig(ctx context.Context) (*models.ObservabilityConfig, error)

	// UpdateObservabilityConfig updates observability configuration
	UpdateObservabilityConfig(ctx context.Context, config *models.ObservabilityConfig) error

	// GetServiceConfig returns configuration for a specific service
	GetServiceConfig(ctx context.Context, serviceName string) (*models.ServiceObservabilityConfig, error)

	// UpdateServiceConfig updates configuration for a specific service
	UpdateServiceConfig(ctx context.Context, serviceName string, config *models.ServiceObservabilityConfig) error
}
