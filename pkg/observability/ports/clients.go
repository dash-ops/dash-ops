package ports

import (
	"context"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LogsClient defines the interface for external log service clients (Loki, Splunk, etc.)
type LogsClient interface {
	// QueryLogs queries logs from the external service using standardized models
	QueryLogs(ctx context.Context, query *models.LogQuery) ([]models.LogEntry, error)

	// GetLogLabels retrieves all available labels
	GetLogLabels(ctx context.Context) ([]string, error)

	// GetLogLevels retrieves all available log levels
	GetLogLevels(ctx context.Context) ([]string, error)

	// HealthCheck checks if the external service is healthy
	HealthCheck(ctx context.Context) error
}

// MetricsClient defines the interface for external metrics service clients (Prometheus, etc.)
type MetricsClient interface {
	// QueryRange queries metrics from the external service within a time range
	QueryRange(ctx context.Context, query string, start, end time.Time, step string) (*wire.MetricsResponse, error)

	// Query performs an instant query at a single point in time
	Query(ctx context.Context, query string, ts time.Time) (*wire.MetricsResponse, error)

	// GetMetricNames returns available metric names
	GetMetricNames(ctx context.Context) ([]string, error)

	// GetLabelValues retrieves all values for a specific label
	GetLabelValues(ctx context.Context, label string) ([]string, error)

	// HealthCheck checks if the external service is healthy
	HealthCheck(ctx context.Context) error
}

// TracesClient defines the interface for external trace service clients (Tempo, Jaeger, etc.)
type TracesClient interface {
	// QueryTraces searches for traces based on criteria
	QueryTraces(ctx context.Context, req *wire.TracesRequest) (*wire.TracesResponse, error)

	// GetTraceDetail retrieves detailed information for a specific trace
	GetTraceDetail(ctx context.Context, traceID string) (*wire.TraceDetailResponse, error)

	// GetServices returns available services that have traces
	GetServices(ctx context.Context) ([]string, error)

	// HealthCheck checks if the external service is healthy
	HealthCheck(ctx context.Context) error
}

// AlertsClient defines the interface for external alert service clients (AlertManager, etc.)
type AlertsClient interface {
	// GetAlerts retrieves active alerts
	GetAlerts(ctx context.Context) (*wire.AlertsResponse, error)

	// GetSilences retrieves active silences
	GetSilences(ctx context.Context) ([]interface{}, error) // TODO: Define proper wire types

	// CreateSilence creates a new silence
	CreateSilence(ctx context.Context, silence interface{}) (interface{}, error) // TODO: Define proper wire types

	// DeleteSilence deletes a silence by ID
	DeleteSilence(ctx context.Context, silenceID string) error

	// HealthCheck checks if the external service is healthy
	HealthCheck(ctx context.Context) error
}
