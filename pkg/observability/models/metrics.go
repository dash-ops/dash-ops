package models

import (
	"time"
)

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

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled     bool     `json:"enabled"`
	Retention   string   `json:"retention"`
	Metrics     []string `json:"metrics"`
	Aggregation string   `json:"aggregation"`
	Step        string   `json:"step"`
}
