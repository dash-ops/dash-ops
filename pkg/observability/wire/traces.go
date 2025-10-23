package wire

import (
	"encoding/json"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

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

// TracesResponse represents the response for trace data
type TracesResponse struct {
	BaseResponse
	Data TracesData `json:"data"`
}

// TracesData represents the data portion of traces response
type TracesData struct {
	Traces []models.TraceSummary `json:"traces"`
	Total  int                   `json:"total"`
	Query  string                `json:"query,omitempty"`
}

// TraceDetailRequest represents a request for detailed trace information
type TraceDetailRequest struct {
	TraceID string `json:"trace_id"`
}

// TraceDetailResponse represents the response for detailed trace information
type TraceDetailResponse struct {
	BaseResponse
	Data TraceDetailData `json:"data"`
}

// TraceDetailData represents the data portion of trace detail response
type TraceDetailData struct {
	TraceID  string             `json:"trace_id"`
	Spans    []models.TraceSpan `json:"spans"`
	Total    int                `json:"total"`
	Timeline TraceTimeline      `json:"timeline,omitempty"`
}

// TraceTimeline represents timeline information for a trace
type TraceTimeline struct {
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	Duration  int64    `json:"duration"`
	Services  []string `json:"services"`
}

// TraceStatsRequest represents a request for trace statistics
type TraceStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
}

// TraceStatisticsResponse represents the response for trace statistics
type TraceStatisticsResponse struct {
	BaseResponse
	Data TraceStatistics `json:"data"`
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

// TraceAnalysisRequest represents a request for trace analysis
type TraceAnalysisRequest struct {
	TraceID string `json:"trace_id"`
}

// TraceAnalysisResponse represents the response for trace analysis
type TraceAnalysisResponse struct {
	BaseResponse
	Data TraceAnalysis `json:"data"`
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

// ProviderInfo is defined in logs.go (shared across logs, traces, metrics, alerts)

// =============================================================================
// Tempo-specific DTOs (external API structures)
// =============================================================================

// TempoSearchResponse represents the response from Tempo's /api/search endpoint
type TempoSearchResponse struct {
	Traces  []TempoTraceSearchResult `json:"traces"`
	Metrics TempoSearchMetrics       `json:"metrics,omitempty"`
}

// TempoTraceSearchResult represents a single trace result from Tempo search
type TempoTraceSearchResult struct {
	TraceID           string        `json:"traceID"`
	RootServiceName   string        `json:"rootServiceName"`
	RootTraceName     string        `json:"rootTraceName"`
	StartTimeUnixNano string        `json:"startTimeUnixNano"`
	DurationMs        int64         `json:"durationMs"`
	SpanSet           *TempoSpanSet `json:"spanSet,omitempty"`
}

// TempoSpanSet represents a set of spans in Tempo
type TempoSpanSet struct {
	Spans   []TempoSpan `json:"spans"`
	Matched int         `json:"matched"`
}

// TempoSpan represents a span in Tempo's search response
type TempoSpan struct {
	SpanID            string                 `json:"spanID"`
	StartTimeUnixNano string                 `json:"startTimeUnixNano"`
	DurationNanos     string                 `json:"durationNanos"`
	Attributes        map[string]interface{} `json:"attributes,omitempty"`
}

// TempoSearchMetrics represents metrics from Tempo search
type TempoSearchMetrics struct {
	InspectedTraces int64  `json:"inspectedTraces,omitempty"`
	InspectedBytes  string `json:"inspectedBytes,omitempty"` // Tempo returns this as string
	InspectedBlocks int64  `json:"inspectedBlocks,omitempty"`
	TotalBlockBytes string `json:"totalBlockBytes,omitempty"` // Tempo returns this as string
	CompletedJobs   int64  `json:"completedJobs,omitempty"`
	TotalJobs       int64  `json:"totalJobs,omitempty"`
}

// TempoTraceByIDResponse represents the response from Tempo's /api/traces/{traceID} endpoint
type TempoTraceByIDResponse struct {
	Batches []TempoBatch `json:"batches"`
}

// TempoBatch represents a batch of spans from a single resource
type TempoBatch struct {
	Resource                    TempoResource                      `json:"resource"`
	ScopeSpans                  []TempoScopeSpans                  `json:"scopeSpans"`
	InstrumentationLibrarySpans []TempoInstrumentationLibrarySpans `json:"instrumentationLibrarySpans,omitempty"` // deprecated
}

// TempoResource represents a resource in Tempo
type TempoResource struct {
	Attributes []TempoAttribute `json:"attributes"`
}

// TempoScopeSpans represents scope spans in Tempo
type TempoScopeSpans struct {
	Scope TempoScope       `json:"scope,omitempty"`
	Spans []TempoTraceSpan `json:"spans"`
}

// TempoInstrumentationLibrarySpans represents instrumentation library spans (deprecated)
type TempoInstrumentationLibrarySpans struct {
	InstrumentationLibrary TempoInstrumentationLibrary `json:"instrumentationLibrary,omitempty"`
	Spans                  []TempoTraceSpan            `json:"spans"`
}

// TempoScope represents a scope in Tempo
type TempoScope struct {
	Name       string           `json:"name"`
	Version    string           `json:"version,omitempty"`
	Attributes []TempoAttribute `json:"attributes,omitempty"`
}

// TempoInstrumentationLibrary represents an instrumentation library (deprecated)
type TempoInstrumentationLibrary struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// TempoSpanKind represents a span kind that can be either string or int
type TempoSpanKind struct {
	Value int
}

// UnmarshalJSON implements custom JSON unmarshaling for TempoSpanKind
func (k *TempoSpanKind) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as int first
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		k.Value = intValue
		return nil
	}

	// Try to unmarshal as string
	var strValue string
	if err := json.Unmarshal(data, &strValue); err != nil {
		return err
	}

	// Map string values to int values
	switch strValue {
	case "SPAN_KIND_UNSPECIFIED":
		k.Value = 0
	case "SPAN_KIND_INTERNAL":
		k.Value = 1
	case "SPAN_KIND_SERVER":
		k.Value = 2
	case "SPAN_KIND_CLIENT":
		k.Value = 3
	case "SPAN_KIND_PRODUCER":
		k.Value = 4
	case "SPAN_KIND_CONSUMER":
		k.Value = 5
	default:
		k.Value = 0
	}

	return nil
}

// TempoTraceSpan represents a single span in Tempo's trace response
type TempoTraceSpan struct {
	TraceID           string           `json:"traceId"`
	SpanID            string           `json:"spanId"`
	ParentSpanID      string           `json:"parentSpanId,omitempty"`
	Name              string           `json:"name"`
	Kind              TempoSpanKind    `json:"kind,omitempty"` // Can be string or int
	StartTimeUnixNano string           `json:"startTimeUnixNano"`
	EndTimeUnixNano   string           `json:"endTimeUnixNano"`
	Attributes        []TempoAttribute `json:"attributes,omitempty"`
	Events            []TempoEvent     `json:"events,omitempty"`
	Links             []TempoLink      `json:"links,omitempty"`
	Status            *TempoStatus     `json:"status,omitempty"`
}

// TempoAttribute represents a key-value attribute in Tempo
type TempoAttribute struct {
	Key   string     `json:"key"`
	Value TempoValue `json:"value"`
}

// TempoValue represents a value in Tempo (can be different types)
type TempoValue struct {
	StringValue string            `json:"stringValue,omitempty"`
	IntValue    string            `json:"intValue,omitempty"`
	DoubleValue float64           `json:"doubleValue,omitempty"`
	BoolValue   bool              `json:"boolValue,omitempty"`
	ArrayValue  *TempoArrayValue  `json:"arrayValue,omitempty"`
	KvlistValue *TempoKvListValue `json:"kvlistValue,omitempty"`
	BytesValue  string            `json:"bytesValue,omitempty"`
}

// TempoArrayValue represents an array value in Tempo
type TempoArrayValue struct {
	Values []TempoValue `json:"values"`
}

// TempoKvListValue represents a key-value list in Tempo
type TempoKvListValue struct {
	Values []TempoAttribute `json:"values"`
}

// TempoEvent represents an event in a span
type TempoEvent struct {
	TimeUnixNano           string           `json:"timeUnixNano"`
	Name                   string           `json:"name"`
	Attributes             []TempoAttribute `json:"attributes,omitempty"`
	DroppedAttributesCount int              `json:"droppedAttributesCount,omitempty"`
}

// TempoLink represents a link to another span
type TempoLink struct {
	TraceID                string           `json:"traceId"`
	SpanID                 string           `json:"spanId"`
	TraceState             string           `json:"traceState,omitempty"`
	Attributes             []TempoAttribute `json:"attributes,omitempty"`
	DroppedAttributesCount int              `json:"droppedAttributesCount,omitempty"`
}

// TempoStatus represents the status of a span
type TempoStatus struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"` // 0=UNSET, 1=OK, 2=ERROR
}

// TempoSearchTagsResponse represents the response from Tempo's /api/search/tags endpoint
type TempoSearchTagsResponse struct {
	TagNames []string `json:"tagNames"`
}

// TempoSearchTagValuesResponse represents the response from Tempo's /api/search/tag/{tagName}/values endpoint
type TempoSearchTagValuesResponse struct {
	TagValues []string `json:"tagValues"`
}

// TempoQueryParams represents query parameters for Tempo API calls
type TempoQueryParams struct {
	Query     string // TraceQL query
	Start     int64  // Start time in Unix nanoseconds
	End       int64  // End time in Unix nanoseconds
	Limit     int    // Limit number of results
	SpanLimit int    // Limit number of spans per trace
}
