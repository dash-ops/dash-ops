package ports

import (
	"context"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LogRepository defines the interface for log data operations
type LogRepository interface {
	// QueryLogs retrieves logs based on the provided criteria
	QueryLogs(ctx context.Context, req *wire.LogsRequest) (*wire.LogsResponse, error)

	// StreamLogs provides real-time log streaming
	StreamLogs(ctx context.Context, req *wire.LogsRequest) (<-chan models.LogEntry, error)

	// GetLogLabels returns available log labels for filtering
	GetLogLabels(ctx context.Context) ([]string, error)

	// GetLogLevels returns available log levels
	GetLogLevels(ctx context.Context) ([]string, error)
}

// MetricRepository defines the interface for metrics data operations
type MetricRepository interface {
	// QueryMetrics retrieves metrics based on the provided criteria
	QueryMetrics(ctx context.Context, req *wire.MetricsRequest) (*wire.MetricsResponse, error)

	// QueryPrometheus executes PromQL queries
	QueryPrometheus(ctx context.Context, req *wire.PrometheusQueryRequest) (*wire.MetricsResponse, error)

	// GetMetricNames returns available metric names
	GetMetricNames(ctx context.Context) ([]string, error)

	// GetMetricLabels returns available labels for a specific metric
	GetMetricLabels(ctx context.Context, metric string) ([]string, error)

	// GetMetricSeries returns time series data for a metric
	GetMetricSeries(ctx context.Context, metric string, start, end time.Time) (*wire.MetricSeries, error)
}

// TraceRepository defines the interface for trace data operations
type TraceRepository interface {
	// QueryTraces retrieves traces based on the provided criteria
	QueryTraces(ctx context.Context, req *wire.TracesRequest) (*wire.TracesResponse, error)

	// GetTraceDetail retrieves detailed information for a specific trace
	GetTraceDetail(ctx context.Context, traceID string) (*wire.TraceDetailResponse, error)

	// GetServices returns available services that have traces
	GetServices(ctx context.Context) ([]string, error)

	// GetOperations returns available operations for a specific service
	GetOperations(ctx context.Context, service string) ([]string, error)

	// GetTraceStatistics returns statistics about traces
	GetTraceStatistics(ctx context.Context, req *wire.TraceStatsRequest) (*wire.TraceStatistics, error)
}

// AlertRepository defines the interface for alert management operations
type AlertRepository interface {
	// GetAlerts retrieves alerts based on the provided criteria
	GetAlerts(ctx context.Context, req *wire.AlertsRequest) (*wire.AlertsResponse, error)

	// CreateAlert creates a new alert rule
	CreateAlert(ctx context.Context, req *wire.CreateAlertRequest) (*models.Alert, error)

	// UpdateAlert updates an existing alert rule
	UpdateAlert(ctx context.Context, id string, req *wire.CreateAlertRequest) (*models.Alert, error)

	// DeleteAlert deletes an alert rule
	DeleteAlert(ctx context.Context, id string) error

	// SilenceAlert silences an alert for a specified duration
	SilenceAlert(ctx context.Context, id string, duration time.Duration) error

	// GetAlertRules returns all configured alert rules
	GetAlertRules(ctx context.Context) ([]models.AlertRule, error)
}

// DashboardRepository defines the interface for dashboard operations
type DashboardRepository interface {
	// GetDashboards retrieves all dashboards
	GetDashboards(ctx context.Context) ([]models.Dashboard, error)

	// GetDashboard retrieves a specific dashboard by ID
	GetDashboard(ctx context.Context, id string) (*models.Dashboard, error)

	// CreateDashboard creates a new dashboard
	CreateDashboard(ctx context.Context, req *wire.CreateDashboardRequest) (*models.Dashboard, error)

	// UpdateDashboard updates an existing dashboard
	UpdateDashboard(ctx context.Context, id string, req *wire.CreateDashboardRequest) (*models.Dashboard, error)

	// DeleteDashboard deletes a dashboard
	DeleteDashboard(ctx context.Context, id string) error

	// GetDashboardTemplates returns available dashboard templates
	GetDashboardTemplates(ctx context.Context) ([]models.DashboardTemplate, error)
}

// ServiceContextRepository defines the interface for service context operations
type ServiceContextRepository interface {
	// GetServiceContext returns service context for observability queries
	GetServiceContext(ctx context.Context, serviceName string) (*models.ServiceContext, error)

	// GetServicesWithContext returns all services with their observability context
	GetServicesWithContext(ctx context.Context) ([]models.ServiceWithContext, error)

	// GetServiceHealth returns health status for a service
	GetServiceHealth(ctx context.Context, serviceName string) (*models.ServiceHealth, error)
}
