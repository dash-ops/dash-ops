package models

import (
	"time"
)

// LogEntry represents a log entry from any log source
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
	Source    string                 `json:"source"`
	Host      string                 `json:"host"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Labels    map[string]string      `json:"labels,omitempty"`
	Raw       string                 `json:"raw,omitempty"`
}

// ProcessedLogEntry represents a processed log entry with additional context
type ProcessedLogEntry struct {
	LogEntry
	ProcessedAt  time.Time              `json:"processed_at"`
	Enrichments  map[string]interface{} `json:"enrichments,omitempty"`
	Correlations []string               `json:"correlations,omitempty"`
}

// LogsConfig represents logs configuration
type LogsConfig struct {
	Enabled     bool     `json:"enabled"`
	Retention   string   `json:"retention"`
	Levels      []string `json:"levels"`
	Sources     []string `json:"sources"`
	QueryLimit  int      `json:"query_limit"`
	StreamLimit int      `json:"stream_limit"`
}
