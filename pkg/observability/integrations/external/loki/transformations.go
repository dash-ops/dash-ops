package loki

import (
	"fmt"
	"strings"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// --- Provider-specific transformation methods ---

// transformToLokiParams transforms models.LogQuery to wire.LokiQueryParams
func (c *LokiClient) transformToLokiParams(query *models.LogQuery) wire.LokiQueryParams {
	// Build Loki LogQL query
	logqlQuery := c.buildLogQLQuery(query)

	// Set defaults for time range
	startTime := query.StartTime
	endTime := query.EndTime
	if startTime.IsZero() {
		startTime = time.Now().Add(-1 * time.Hour)
	}
	if endTime.IsZero() {
		endTime = time.Now()
	}

	// Set limit
	limit := query.Limit
	if limit == 0 {
		limit = 100
	}

	return wire.LokiQueryParams{
		Query:     logqlQuery,
		Start:     startTime,
		End:       endTime,
		Limit:     limit,
		Direction: "backward", // Default direction
	}
}

// buildLogQLQuery builds a LogQL query from models.LogQuery
func (c *LokiClient) buildLogQLQuery(query *models.LogQuery) string {
	// If custom query is provided, use it
	if query.Query != "" {
		return query.Query
	}

	// Build query from filters
	var filters []string

	// Service filter - look for app label (Promtail uses this)
	if query.Service != "" {
		filters = append(filters, fmt.Sprintf(`app="%s"`, query.Service))
	}

	// Level filter - look for level label
	if query.Level != "" {
		filters = append(filters, fmt.Sprintf(`level="%s"`, query.Level))
	}

	// Build LogQL selector
	if len(filters) == 0 {
		return `{job=~".+"}` // Default query to get all logs
	}

	return fmt.Sprintf("{%s}", strings.Join(filters, ","))
}

// transformLokiResponseToModels transforms wire.LokiQueryResponse to []models.LogEntry
func (c *LokiClient) transformLokiResponseToModels(resp *wire.LokiQueryResponse) []models.LogEntry {
	var entries []models.LogEntry

	for _, stream := range resp.Data.Result {
		for _, value := range stream.Values {
			if len(value) < 2 {
				continue
			}

			// Parse timestamp (nanoseconds)
			tsNano, err := c.parseTimestamp(value[0])
			if err != nil {
				continue
			}

			// Create log entry
			entry := models.LogEntry{
				Timestamp: time.Unix(0, tsNano),
				Message:   value[1],
				Labels:    stream.Stream,
				Level:     stream.Stream["level"],
				Service:   stream.Stream["app"], // Promtail uses 'app' label
				Host:      stream.Stream["hostname"],
				Source:    stream.Stream["filename"],
			}

			// Try to parse trace/span IDs if present
			if traceID, ok := stream.Stream["trace_id"]; ok {
				entry.TraceID = traceID
			}
			if spanID, ok := stream.Stream["span_id"]; ok {
				entry.SpanID = spanID
			}

			// Generate ID from timestamp and message hash
			entry.ID = fmt.Sprintf("%d_%d", tsNano, len(value[1]))

			entries = append(entries, entry)
		}
	}

	return entries
}

// parseTimestamp parses a timestamp string (nanoseconds) to int64
func (c *LokiClient) parseTimestamp(ts string) (int64, error) {
	// Try parsing as nanoseconds
	var result int64
	_, err := fmt.Sscanf(ts, "%d", &result)
	return result, err
}
