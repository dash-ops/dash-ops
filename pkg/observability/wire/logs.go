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

// --- Loki-specific DTOs (External API responses) ---

// LokiQueryResponse represents the response from Loki query API
type LokiQueryResponse struct {
	Status string         `json:"status"`
	Data   LokiResultData `json:"data"`
}

// LokiResultData represents the data field in Loki response
type LokiResultData struct {
	ResultType string       `json:"resultType"` // streams or matrix
	Result     []LokiStream `json:"result"`
	Stats      *LokiStats   `json:"stats,omitempty"`
}

// LokiStream represents a log stream with its values
type LokiStream struct {
	Stream map[string]string `json:"stream"` // labels
	Values [][]string        `json:"values"` // [timestamp_ns, log_line]
}

// LokiStats represents query statistics from Loki
type LokiStats struct {
	Summary  LokiStatsSummary  `json:"summary"`
	Querier  LokiStatsQuerier  `json:"querier"`
	Ingester LokiStatsIngester `json:"ingester"`
}

// LokiStatsSummary represents summary statistics
type LokiStatsSummary struct {
	BytesProcessedPerSecond int     `json:"bytesProcessedPerSecond"`
	LinesProcessedPerSecond int     `json:"linesProcessedPerSecond"`
	TotalBytesProcessed     int     `json:"totalBytesProcessed"`
	TotalLinesProcessed     int     `json:"totalLinesProcessed"`
	ExecTime                float64 `json:"execTime"`
}

// LokiStatsQuerier represents querier statistics
type LokiStatsQuerier struct {
	Store LokiStatsStore `json:"store"`
}

// LokiStatsStore represents store statistics
type LokiStatsStore struct {
	TotalChunksRef        int `json:"totalChunksRef"`
	TotalChunksDownloaded int `json:"totalChunksDownloaded"`
}

// LokiStatsIngester represents ingester statistics
type LokiStatsIngester struct {
	TotalReached       int `json:"totalReached"`
	TotalChunksMatched int `json:"totalChunksMatched"`
	TotalLinesSent     int `json:"totalLinesSent"`
}

// LokiLabelsResponse represents the response from labels API
type LokiLabelsResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

// LokiLabelValuesResponse represents the response from label values API
type LokiLabelValuesResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

// LokiSeriesResponse represents the response from series API
type LokiSeriesResponse struct {
	Status string              `json:"status"`
	Data   []map[string]string `json:"data"`
}

// LokiError represents an error response from Loki
type LokiError struct {
	Status    string `json:"status"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}

// LokiQueryParams represents parameters for a Loki query
type LokiQueryParams struct {
	Query     string    `json:"query"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Limit     int       `json:"limit,omitempty"`
	Direction string    `json:"direction,omitempty"` // forward or backward
	Step      string    `json:"step,omitempty"`      // for range queries
}
