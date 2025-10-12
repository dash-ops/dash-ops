package models

import "time"

// TraceQuery represents a standardized query for traces (provider-agnostic)
type TraceQuery struct {
	// Service name to filter traces
	Service string

	// Operation name to filter traces
	Operation string

	// Tags to filter traces (key-value pairs)
	Tags map[string]string

	// Minimum duration for traces (e.g., "100ms", "1s")
	MinDuration string

	// Maximum duration for traces (e.g., "1s", "10s")
	MaxDuration string

	// Time range
	StartTime time.Time
	EndTime   time.Time

	// Limit the number of results
	Limit int
}

// TraceSpan represents a single span in a trace
type TraceSpan struct {
	TraceID       string                 `json:"trace_id"`
	SpanID        string                 `json:"span_id"`
	ParentSpanID  string                 `json:"parent_span_id,omitempty"`
	OperationName string                 `json:"operation_name"`
	ServiceName   string                 `json:"service_name"`
	StartTime     time.Time              `json:"start_time"`
	Duration      time.Duration          `json:"duration"`
	Tags          map[string]interface{} `json:"tags,omitempty"`
	Logs          []SpanLog              `json:"logs,omitempty"`
	References    []SpanReference        `json:"references,omitempty"`
	Status        SpanStatus             `json:"status"`
}

// SpanLog represents a log entry within a span
type SpanLog struct {
	Timestamp time.Time              `json:"timestamp"`
	Fields    map[string]interface{} `json:"fields"`
}

// SpanReference represents a reference to another span
type SpanReference struct {
	RefType string `json:"ref_type"` // CHILD_OF or FOLLOWS_FROM
	TraceID string `json:"trace_id"`
	SpanID  string `json:"span_id"`
}

// SpanStatus represents the status of a span
type SpanStatus struct {
	Code    int    `json:"code"` // 0=OK, 1=ERROR, 2=UNSET
	Message string `json:"message,omitempty"`
}

// Trace represents a complete trace with all its spans
type Trace struct {
	TraceID   string        `json:"trace_id"`
	Spans     []TraceSpan   `json:"spans"`
	Duration  time.Duration `json:"duration"`
	StartTime time.Time     `json:"start_time"`
	Services  []string      `json:"services"` // List of services involved in the trace
}

// TraceSummary represents a summary of a trace (for list views)
type TraceSummary struct {
	TraceID       string        `json:"trace_id"`
	RootService   string        `json:"root_service"`
	RootOperation string        `json:"root_operation"`
	StartTime     time.Time     `json:"start_time"`
	Duration      time.Duration `json:"duration"`
	SpanCount     int           `json:"span_count"`
	ErrorCount    int           `json:"error_count"`
	Services      []string      `json:"services"`
}
