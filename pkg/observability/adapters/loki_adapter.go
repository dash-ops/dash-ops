package adapters

import (
	"fmt"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LokiQueryAdapter transforms between domain models and Loki wire types
type LokiQueryAdapter struct{}

// NewLokiQueryAdapter creates a new Loki query adapter
func NewLokiQueryAdapter() *LokiQueryAdapter {
	return &LokiQueryAdapter{}
}

// ModelToLokiParams transforms models.LogQuery to wire.LokiQueryParams
func (a *LokiQueryAdapter) ModelToLokiParams(query *models.LogQuery) wire.LokiQueryParams {
	params := wire.LokiQueryParams{
		Start: query.StartTime,
		End:   query.EndTime,
		Limit: query.Limit,
	}

	// Build LogQL query from model
	params.Query = a.buildLogQLQuery(query)

	// Set direction based on sort order
	if query.Limit > 0 {
		params.Direction = "backward" // Most recent first
	}

	return params
}

// buildLogQLQuery builds a LogQL query string from domain model
func (a *LokiQueryAdapter) buildLogQLQuery(query *models.LogQuery) string {
	// If custom LogQL query is provided, use it
	if query.Query != "" {
		return query.Query
	}

	// Build query from filters
	filters := []string{}

	// Service filter
	if query.Service != "" {
		filters = append(filters, fmt.Sprintf(`service="%s"`, query.Service))
	}

	// Level filter
	if query.Level != "" {
		filters = append(filters, fmt.Sprintf(`level="%s"`, query.Level))
	}

	// Build LogQL selector
	if len(filters) == 0 {
		return `{job=~".+"}`
	}

	selector := "{"
	for i, filter := range filters {
		if i > 0 {
			selector += ","
		}
		selector += filter
	}
	selector += "}"

	return selector
}

// LokiStreamToModel transforms wire.LokiStream to models.LogEntry
func (a *LokiQueryAdapter) LokiStreamToModel(stream wire.LokiStream, timestamp int64, line string) models.LogEntry {
	entry := models.LogEntry{
		Timestamp: parseNanoTimestamp(timestamp),
		Message:   line,
		Labels:    stream.Stream,
		Level:     stream.Stream["level"],
		Service:   stream.Stream["service"],
		Host:      stream.Stream["host"],
		Source:    stream.Stream["source"],
	}

	// Extract trace/span IDs if present
	if traceID, ok := stream.Stream["trace_id"]; ok {
		entry.TraceID = traceID
	}
	if spanID, ok := stream.Stream["span_id"]; ok {
		entry.SpanID = spanID
	}

	return entry
}

// parseNanoTimestamp converts nanosecond timestamp to time.Time
func parseNanoTimestamp(nanos int64) time.Time {
	return time.Unix(0, nanos)
}
