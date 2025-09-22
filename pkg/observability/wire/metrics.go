package wire

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

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

// PrometheusQueryRequest represents a Prometheus query request
type PrometheusQueryRequest struct {
	Query     string        `json:"query"` // PromQL query
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Step      string        `json:"step,omitempty"`
	Timeout   time.Duration `json:"timeout,omitempty"`
}

// MetricStatsRequest represents a request for metric statistics
type MetricStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	Metric    string    `json:"metric,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
}

// MetricStatisticsResponse represents the response for metric statistics
type MetricStatisticsResponse struct {
	BaseResponse
	Data MetricStatistics `json:"data"`
}

// MetricStatistics represents metric statistics
type MetricStatistics struct {
	TotalMetrics     int64                  `json:"total_metrics"`
	MetricsByType    map[string]int64       `json:"metrics_by_type"`
	MetricsByService map[string]int64       `json:"metrics_by_service"`
	TimeSeries       []TimeSeriesData       `json:"time_series,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// MetricSeries represents a metric time series
type MetricSeries struct {
	Metric   string                 `json:"metric"`
	Labels   map[string]string      `json:"labels"`
	Values   []models.MetricData    `json:"values"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
