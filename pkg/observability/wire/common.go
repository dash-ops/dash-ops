package wire

import (
	"time"
)

// BaseResponse represents a base response structure
type BaseResponse struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TimeSeriesData represents time series data for statistics
type TimeSeriesData struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Label     string    `json:"label,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	BaseResponse
	ErrorCode    string                 `json:"error_code,omitempty"`
	ErrorDetails map[string]interface{} `json:"error_details,omitempty"`
	RequestID    string                 `json:"request_id,omitempty"`
	Timestamp    time.Time              `json:"timestamp"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// QueryInfo represents query information
type QueryInfo struct {
	Query     string                 `json:"query"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	TimeRange TimeRange              `json:"time_range"`
	Limit     int                    `json:"limit,omitempty"`
	Offset    int                    `json:"offset,omitempty"`
}

// PerformanceInfo represents performance information
type PerformanceInfo struct {
	QueryTime   time.Duration `json:"query_time"`
	ProcessTime time.Duration `json:"process_time"`
	TotalTime   time.Duration `json:"total_time"`
	CacheHit    bool          `json:"cache_hit"`
	CacheTime   time.Duration `json:"cache_time,omitempty"`
	ResultSize  int           `json:"result_size"`
}
