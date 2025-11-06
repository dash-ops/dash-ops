package adapters

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// ExplorerAdapter handles transformation between wire DTOs and domain models for explorer queries
type ExplorerAdapter struct {
	logsAdapter   *LogsAdapter
	tracesAdapter *TracesAdapter
}

// NewExplorerAdapter creates a new explorer adapter
func NewExplorerAdapter(logsAdapter *LogsAdapter, tracesAdapter *TracesAdapter) *ExplorerAdapter {
	return &ExplorerAdapter{
		logsAdapter:   logsAdapter,
		tracesAdapter: tracesAdapter,
	}
}

// WireRequestToModel transforms wire.ExplorerQueryRequest to parameters for controller
func (a *ExplorerAdapter) WireRequestToModel(req *wire.ExplorerQueryRequest) (
	query string,
	timeRangeFrom, timeRangeTo *time.Time,
	provider string,
) {
	query = req.Query

	// Parse time range using RFC3339 format (respects timezone in the string)
	if req.TimeRangeFrom != "" {
		if t, err := time.Parse(time.RFC3339, req.TimeRangeFrom); err == nil {
			utc := t.UTC()
			timeRangeFrom = &utc
		}
	}
	if req.TimeRangeTo != "" {
		if t, err := time.Parse(time.RFC3339, req.TimeRangeTo); err == nil {
			utc := t.UTC()
			timeRangeTo = &utc
		}
	}

	provider = req.Provider
	return
}

// ModelToWireResponse transforms domain models to wire.ExplorerQueryResponse
func (a *ExplorerAdapter) ModelToWireResponse(
	dataSource string,
	results interface{},
	total int,
	query string,
	executionTimeMs int64,
) *wire.ExplorerQueryResponse {
	// Transform results based on data source
	var transformedResults interface{}

	switch dataSource {
	case "logs":
		if logs, ok := results.([]models.LogEntry); ok {
			// Wire format uses models.LogEntry directly
			transformedResults = logs
		} else {
			transformedResults = []models.LogEntry{}
		}
	case "traces":
		if spans, ok := results.([]models.TraceSpan); ok {
			// Wire format uses models.TraceSpan directly
			transformedResults = spans
		} else {
			transformedResults = []models.TraceSpan{}
		}
	case "metrics":
		// TODO: Transform metrics
		transformedResults = []interface{}{}
	default:
		transformedResults = []interface{}{}
	}

	return &wire.ExplorerQueryResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.ExplorerQueryData{
			DataSource:      dataSource,
			Results:         transformedResults,
			Total:           total,
			Query:           query,
			ExecutionTimeMs: executionTimeMs,
		},
	}
}

// ErrorToWireResponse transforms an error to wire.ExplorerQueryResponse
func (a *ExplorerAdapter) ErrorToWireResponse(err error) *wire.ExplorerQueryResponse {
	return &wire.ExplorerQueryResponse{
		BaseResponse: wire.BaseResponse{
			Success: false,
			Error:   err.Error(),
		},
		Data: wire.ExplorerQueryData{
			DataSource:      "",
			Results:         []interface{}{},
			Total:           0,
			Query:           "",
			ExecutionTimeMs: 0,
		},
	}
}
