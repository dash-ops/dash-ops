package adapters

import (
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// TracesAdapter handles transformations between HTTP wire DTOs and domain models for traces
type TracesAdapter struct{}

// NewTracesAdapter creates a new traces adapter
func NewTracesAdapter() *TracesAdapter {
	return &TracesAdapter{}
}

// WireRequestToModel transforms wire.TracesRequest to models.TraceQuery
func (a *TracesAdapter) WireRequestToModel(req *wire.TracesRequest) *models.TraceQuery {
	query := &models.TraceQuery{
		Service:     req.Service,
		Operation:   req.Operation,
		MinDuration: req.MinDuration,
		MaxDuration: req.MaxDuration,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Limit:       req.Limit,
		Tags:        make(map[string]string),
	}

	// Set default time range if not provided
	if query.StartTime.IsZero() {
		query.StartTime = time.Now().Add(-1 * time.Hour)
	}
	if query.EndTime.IsZero() {
		query.EndTime = time.Now()
	}

	// Set default limit if not provided
	if query.Limit == 0 {
		query.Limit = 20
	}

	// TODO: Extract tags from req if needed in the future

	return query
}

// ModelToWireResponse transforms []models.TraceSummary to wire.TracesResponse
func (a *TracesAdapter) ModelToWireResponse(traces []models.TraceSummary, provider string, providerType string) *wire.TracesResponse {
	return &wire.TracesResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.TracesData{
			Traces: traces,
			Total:  len(traces),
		},
	}
}

// ModelToWireDetailResponse transforms models.Trace to wire.TraceDetailResponse
func (a *TracesAdapter) ModelToWireDetailResponse(trace *models.Trace) *wire.TraceDetailResponse {
	// Convert models.TraceSpan to wire (reusing models.TraceSpan directly)
	spans := make([]models.TraceSpan, len(trace.Spans))
	copy(spans, trace.Spans)

	return &wire.TraceDetailResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.TraceDetailData{
			TraceID: trace.TraceID,
			Spans:   spans,
			Total:   len(spans),
			Timeline: wire.TraceTimeline{
				StartTime: trace.StartTime.UnixNano(),
				EndTime:   trace.StartTime.Add(trace.Duration).UnixNano(),
				Duration:  trace.Duration.Nanoseconds(),
				Services:  trace.Services,
			},
		},
	}
}

// ErrorToWireResponse transforms an error to wire.TracesResponse
func (a *TracesAdapter) ErrorToWireResponse(err error) *wire.TracesResponse {
	return &wire.TracesResponse{
		BaseResponse: wire.BaseResponse{
			Success: false,
			Error:   err.Error(),
		},
	}
}

// ErrorToWireDetailResponse transforms an error to wire.TraceDetailResponse
func (a *TracesAdapter) ErrorToWireDetailResponse(err error) *wire.TraceDetailResponse {
	return &wire.TraceDetailResponse{
		BaseResponse: wire.BaseResponse{
			Success: false,
			Error:   err.Error(),
		},
	}
}
