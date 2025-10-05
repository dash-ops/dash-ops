package adapters

import (
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LogsAdapter handles transformation between wire DTOs and domain models for logs
type LogsAdapter struct{}

// NewLogsAdapter creates a new logs adapter
func NewLogsAdapter() *LogsAdapter {
	return &LogsAdapter{}
}

// WireRequestToModel transforms wire.LogsRequest to a format usable by business logic
func (a *LogsAdapter) WireRequestToModel(req *wire.LogsRequest) *models.LogQuery {
	return &models.LogQuery{
		Service:   req.Service,
		Level:     req.Level,
		Query:     req.Query,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}
}

// ModelToWireResponse transforms domain models to wire.LogsResponse
func (a *LogsAdapter) ModelToWireResponse(logs []models.LogEntry, total int, hasMore bool) *wire.LogsResponse {
	return &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:    logs,
			Total:   total,
			HasMore: hasMore,
		},
	}
}

// ErrorToWireResponse transforms an error to wire.LogsResponse
func (a *LogsAdapter) ErrorToWireResponse(err error) *wire.LogsResponse {
	return &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: false,
			Error:   err.Error(),
		},
	}
}
