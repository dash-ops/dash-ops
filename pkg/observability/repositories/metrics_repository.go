package repositories

import (
	"context"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/integrations/external/prometheus"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// MetricsRepository implements metrics data access using Prometheus
type MetricsRepository struct {
	prometheusAdapter ports.MetricRepository
}

// NewMetricsRepository creates a new metrics repository
func NewMetricsRepository(prometheusClient *prometheus.PrometheusClient) *MetricsRepository {
	prometheusAdapter := prometheus.NewPrometheusAdapter(prometheusClient)

	return &MetricsRepository{
		prometheusAdapter: prometheusAdapter,
	}
}

// QueryMetrics queries metrics from the repository
func (r *MetricsRepository) QueryMetrics(ctx context.Context, req *wire.MetricsRequest) (*wire.MetricsResponse, error) {
	return r.prometheusAdapter.QueryMetrics(ctx, req)
}

// QueryPrometheus executes direct PromQL queries
func (r *MetricsRepository) QueryPrometheus(ctx context.Context, req *wire.PrometheusQueryRequest) (*wire.MetricsResponse, error) {
	return r.prometheusAdapter.QueryPrometheus(ctx, req)
}

// GetMetricNames retrieves available metric names
func (r *MetricsRepository) GetMetricNames(ctx context.Context) ([]string, error) {
	return r.prometheusAdapter.GetMetricNames(ctx)
}

// GetMetricLabels retrieves labels for a specific metric
func (r *MetricsRepository) GetMetricLabels(ctx context.Context, metric string) ([]string, error) {
	return r.prometheusAdapter.GetMetricLabels(ctx, metric)
}

// GetMetricSeries retrieves time series for a metric
func (r *MetricsRepository) GetMetricSeries(ctx context.Context, metric string, start, end time.Time) (*wire.MetricSeries, error) {
	return r.prometheusAdapter.GetMetricSeries(ctx, metric, start, end)
}
