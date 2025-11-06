package wire

// ExplorerQueryRequest represents a request to execute an explorer query
type ExplorerQueryRequest struct {
	Query         string `json:"query"`
	TimeRangeFrom string `json:"time_range_from,omitempty"` // ISO timestamp
	TimeRangeTo   string `json:"time_range_to,omitempty"`   // ISO timestamp
	Provider      string `json:"provider,omitempty"`
}

// ExplorerQueryResponse represents the response from an explorer query
type ExplorerQueryResponse struct {
	BaseResponse
	Data ExplorerQueryData `json:"data"`
}

// ExplorerQueryData represents the data portion of an explorer query response
type ExplorerQueryData struct {
	DataSource      string      `json:"data_source"` // "logs" | "traces" | "metrics"
	Results         interface{} `json:"results"`     // Can be []LogEntry, []TraceSpan, or []MetricData
	Total           int         `json:"total"`
	Query           string      `json:"query"`
	ExecutionTimeMs int64       `json:"execution_time_ms"`
}

// ParsedQuery represents the parsed query structure
type ParsedQuery struct {
	DataSource string                 // "logs" | "traces" | "metrics"
	Filters    map[string]interface{} // Parsed WHERE conditions
	RawQuery   string                 // Original query
}
