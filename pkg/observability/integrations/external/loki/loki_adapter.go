package loki

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LokiAdapter implements log repository using Loki
type LokiAdapter struct {
	client *LokiClient
}

// NewLokiAdapter creates a new Loki adapter
func NewLokiAdapter(client *LokiClient) ports.LogRepository {
	return &LokiAdapter{
		client: client,
	}
}

// QueryLogs implements LogRepository
func (l *LokiAdapter) QueryLogs(ctx context.Context, req *wire.LogsRequest) (*wire.LogsResponse, error) {
	// Build Loki query params from request
	params := wire.LokiQueryParams{
		Query:     buildLokiQuery(req),
		Start:     req.StartTime,
		End:       req.EndTime,
		Limit:     req.Limit,
		Direction: req.Order, // asc -> forward, desc -> backward
	}

	// Set defaults
	if params.Start.IsZero() {
		params.Start = time.Now().Add(-1 * time.Hour)
	}
	if params.End.IsZero() {
		params.End = time.Now()
	}
	if params.Limit == 0 {
		params.Limit = 100
	}
	if params.Direction == "" || params.Direction == "desc" {
		params.Direction = "backward"
	} else {
		params.Direction = "forward"
	}

	// Query Loki (returns wire.LokiQueryResponse)
	lokiResp, err := l.client.QueryRange(ctx, params)
	if err != nil {
		return &wire.LogsResponse{
			BaseResponse: wire.BaseResponse{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	// Transform wire.LokiQueryResponse -> []models.LogEntry
	logs := l.transformLokiStreamsToLogEntries(lokiResp.Data.Result)

	return &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{Success: true},
		Data: wire.LogsData{
			Logs:    logs,
			Total:   len(logs),
			HasMore: len(logs) >= params.Limit,
		},
	}, nil
}

// StreamLogs implements LogRepository
func (l *LokiAdapter) StreamLogs(ctx context.Context, req *wire.LogsRequest) (<-chan models.LogEntry, error) {
	ch := make(chan models.LogEntry, 100)

	go func() {
		defer close(ch)

		// Build Loki query params
		params := wire.LokiQueryParams{
			Query:     buildLokiQuery(req),
			Start:     req.StartTime,
			End:       req.EndTime,
			Limit:     req.Limit,
			Direction: "backward",
		}

		// Set defaults
		if params.Start.IsZero() {
			params.Start = time.Now().Add(-1 * time.Hour)
		}
		if params.End.IsZero() {
			params.End = time.Now()
		}

		// Query Loki
		lokiResp, err := l.client.QueryRange(ctx, params)
		if err != nil {
			return
		}

		// Transform and stream logs
		logs := l.transformLokiStreamsToLogEntries(lokiResp.Data.Result)
		for _, log := range logs {
			select {
			case ch <- log:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

// GetLogLabels implements LogRepository
func (l *LokiAdapter) GetLogLabels(ctx context.Context) ([]string, error) {
	resp, err := l.client.ListLabels(ctx, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetLogLevels implements LogRepository
func (l *LokiAdapter) GetLogLevels(ctx context.Context) ([]string, error) {
	resp, err := l.client.GetLabelValues(ctx, "level", time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// --- Transformation Functions (wire -> models) ---

// transformLokiStreamsToLogEntries converts Loki streams to domain log entries
func (l *LokiAdapter) transformLokiStreamsToLogEntries(streams []wire.LokiStream) []models.LogEntry {
	var entries []models.LogEntry

	for _, stream := range streams {
		for _, value := range stream.Values {
			if len(value) < 2 {
				continue
			}

			// Parse timestamp (nanoseconds)
			tsNano, err := parseTimestamp(value[0])
			if err != nil {
				continue
			}

			// Extract log fields from labels
			entry := models.LogEntry{
				Timestamp: time.Unix(0, tsNano),
				Message:   value[1],
				Labels:    stream.Stream,
				Level:     stream.Stream["level"],
				Service:   stream.Stream["service"],
				Host:      stream.Stream["host"],
				Source:    stream.Stream["source"],
			}

			// Try to parse trace/span IDs if present
			if traceID, ok := stream.Stream["trace_id"]; ok {
				entry.TraceID = traceID
			}
			if spanID, ok := stream.Stream["span_id"]; ok {
				entry.SpanID = spanID
			}

			entries = append(entries, entry)
		}
	}

	return entries
}

// buildLokiQuery converts a logs request to Loki LogQL query format
func buildLokiQuery(req *wire.LogsRequest) string {
	// If custom query is provided, use it
	if req.Query != "" {
		return req.Query
	}

	// Build query from filters
	filters := []string{}

	// Service filter
	if req.Service != "" {
		filters = append(filters, fmt.Sprintf(`service="%s"`, req.Service))
	}

	// Level filter
	if req.Level != "" {
		filters = append(filters, fmt.Sprintf(`level="%s"`, req.Level))
	}

	// Build LogQL selector
	if len(filters) == 0 {
		return `{job=~".+"}`
	}

	return fmt.Sprintf("{%s}", joinFilters(filters))
}

// joinFilters joins filter strings with commas
func joinFilters(filters []string) string {
	result := ""
	for i, filter := range filters {
		if i > 0 {
			result += ","
		}
		result += filter
	}
	return result
}

// parseTimestamp parses a timestamp string (nanoseconds) to int64
func parseTimestamp(ts string) (int64, error) {
	return strconv.ParseInt(ts, 10, 64)
}
