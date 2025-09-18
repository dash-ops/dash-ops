package wire

import (
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
	Traces []models.TraceInfo `json:"traces"`
	Total  int                `json:"total"`
	Query  string             `json:"query,omitempty"`
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
