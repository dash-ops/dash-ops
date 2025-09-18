package models

import (
	"time"
)

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

// TracesConfig represents traces configuration
type TracesConfig struct {
	Enabled    bool     `json:"enabled"`
	Retention  string   `json:"retention"`
	Services   []string `json:"services"`
	Operations []string `json:"operations"`
	MaxDepth   int      `json:"max_depth"`
}
