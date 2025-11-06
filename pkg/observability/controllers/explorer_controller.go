package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/logic"
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// ExplorerController handles explorer query operations
type ExplorerController struct {
	logsClients   map[string]ports.LogsClient
	tracesClients map[string]ports.TracesClient
	queryParser   *logic.QueryParser
}

// NewExplorerController creates a new explorer controller
func NewExplorerController(
	logsClients map[string]ports.LogsClient,
	tracesClients map[string]ports.TracesClient,
) *ExplorerController {
	return &ExplorerController{
		logsClients:   logsClients,
		tracesClients: tracesClients,
		queryParser:   logic.NewQueryParser(),
	}
}

// ExecuteQuery executes an explorer query and returns results
func (c *ExplorerController) ExecuteQuery(
	ctx context.Context,
	query string,
	timeRangeFrom, timeRangeTo *time.Time,
	provider string,
) (dataSource string, results interface{}, total int, executionTimeMs int64, err error) {
	startTime := time.Now()

	// Parse query to determine data source
	parsed, err := c.queryParser.Parse(query)
	if err != nil {
		return "", nil, 0, 0, fmt.Errorf("failed to parse query: %w", err)
	}

	dataSource = parsed.DataSource

	// Validate that provider is specified
	if provider == "" {
		return "", nil, 0, 0, fmt.Errorf("provider is required")
	}

	// Execute query based on data source
	switch dataSource {
	case "logs":
		results, total, err = c.executeLogsQuery(ctx, parsed, timeRangeFrom, timeRangeTo, provider)
	case "traces":
		results, total, err = c.executeTracesQuery(ctx, parsed, timeRangeFrom, timeRangeTo, provider)
	case "metrics":
		// TODO: Implement metrics query
		return dataSource, []interface{}{}, 0, 0, fmt.Errorf("metrics queries not yet implemented")
	default:
		return "", nil, 0, 0, fmt.Errorf("unsupported data source: %s", dataSource)
	}

	if err != nil {
		return dataSource, nil, 0, 0, err
	}

	executionTimeMs = time.Since(startTime).Milliseconds()
	return dataSource, results, total, executionTimeMs, nil
}

// executeLogsQuery executes a logs query
func (c *ExplorerController) executeLogsQuery(
	ctx context.Context,
	parsed *wire.ParsedQuery,
	timeRangeFrom, timeRangeTo *time.Time,
	provider string,
) ([]models.LogEntry, int, error) {
	client, exists := c.logsClients[provider]
	if !exists {
		return nil, 0, fmt.Errorf("logs provider not found: %s", provider)
	}

	// Build query parameters from parsed filters
	logQuery := &models.LogQuery{
		Limit: 100, // Default limit
	}

	if timeRangeFrom != nil {
		logQuery.StartTime = *timeRangeFrom
	} else {
		logQuery.StartTime = time.Now().Add(-1 * time.Hour)
	}

	if timeRangeTo != nil {
		logQuery.EndTime = *timeRangeTo
	} else {
		logQuery.EndTime = time.Now()
	}

	// Apply filters from parsed query
	if level, ok := parsed.Filters["level"].(string); ok {
		logQuery.Level = level
	}
	if service, ok := parsed.Filters["service"].(string); ok {
		logQuery.Service = service
	}

	// Use raw query if it's a LogQL-style query
	if parsed.RawQuery != "" {
		logQuery.Query = parsed.RawQuery
	}

	logs, err := client.QueryLogs(ctx, logQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query logs: %w", err)
	}

	return logs, len(logs), nil
}

// executeTracesQuery executes a traces query
func (c *ExplorerController) executeTracesQuery(
	ctx context.Context,
	parsed *wire.ParsedQuery,
	timeRangeFrom, timeRangeTo *time.Time,
	provider string,
) ([]models.TraceSpan, int, error) {
	client, exists := c.tracesClients[provider]
	if !exists {
		return nil, 0, fmt.Errorf("traces provider not found: %s", provider)
	}

	// Build query parameters from parsed filters
	traceQuery := &models.TraceQuery{
		Limit: 100, // Default limit - we'll fetch details for up to 100 traces
	}

	if timeRangeFrom != nil {
		traceQuery.StartTime = *timeRangeFrom
	} else {
		traceQuery.StartTime = time.Now().Add(-1 * time.Hour)
	}

	if timeRangeTo != nil {
		traceQuery.EndTime = *timeRangeTo
	} else {
		traceQuery.EndTime = time.Now()
	}

	// Apply filters from parsed query
	if service, ok := parsed.Filters["service"].(string); ok {
		traceQuery.Service = service
	}
	if operation, ok := parsed.Filters["operation"].(string); ok {
		traceQuery.Operation = operation
	}

	// Step 1: Query traces to get summaries
	summaries, err := client.QueryTraces(ctx, traceQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query traces: %w", err)
	}

	// Step 2: For each summary, get the trace detail (which contains spans)
	var allSpans []models.TraceSpan
	for _, summary := range summaries {
		trace, err := client.GetTraceDetail(ctx, summary.TraceID)
		if err != nil {
			// Log error but continue with other traces
			continue
		}

		// Extract spans from the trace
		if trace != nil && trace.Spans != nil {
			allSpans = append(allSpans, trace.Spans...)
		}
	}

	return allSpans, len(allSpans), nil
}
