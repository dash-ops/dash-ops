package wire

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// LogsRequest represents a request for log data
type LogsRequest struct {
	Service   string    `json:"service,omitempty"`
	Level     string    `json:"level,omitempty"` // error, warn, info, debug
	Query     string    `json:"query,omitempty"` // Loki query syntax
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Limit     int       `json:"limit,omitempty"` // default: 100
	Offset    int       `json:"offset,omitempty"`
	Stream    bool      `json:"stream,omitempty"` // real-time streaming
	Sort      string    `json:"sort,omitempty"`   // timestamp, level, service
	Order     string    `json:"order,omitempty"`  // asc, desc
}

// LogsResponse represents the response for log data
type LogsResponse struct {
	BaseResponse
	Data LogsData `json:"data"`
}

// LogsData represents the data portion of logs response
type LogsData struct {
	Logs       []models.LogEntry `json:"logs"`
	Total      int               `json:"total"`
	HasMore    bool              `json:"has_more"`
	NextOffset int               `json:"next_offset,omitempty"`
	Filters    LogFilters        `json:"filters,omitempty"`
}

// LogFilters represents filters for log queries
type LogFilters struct {
	Levels    []string          `json:"levels,omitempty"`
	Services  []string          `json:"services,omitempty"`
	Hosts     []string          `json:"hosts,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	TimeRange string            `json:"time_range,omitempty"`
	Query     string            `json:"query,omitempty"`
}

// LogStatsRequest represents a request for log statistics
type LogStatsRequest struct {
	Service   string    `json:"service,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   []string  `json:"group_by,omitempty"`
}

// LogStatisticsResponse represents the response for log statistics
type LogStatisticsResponse struct {
	BaseResponse
	Data LogStatistics `json:"data"`
}

// LogStatistics represents log statistics
type LogStatistics struct {
	TotalLogs     int64                  `json:"total_logs"`
	LogsByLevel   map[string]int64       `json:"logs_by_level"`
	LogsByService map[string]int64       `json:"logs_by_service"`
	LogsByHost    map[string]int64       `json:"logs_by_host"`
	TimeSeries    []TimeSeriesData       `json:"time_series,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}
