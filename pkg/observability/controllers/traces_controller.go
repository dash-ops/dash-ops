package controllers

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/observability/logic"
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/repositories"
)

// TracesController handles trace-related use cases
type TracesController struct {
	// Dependencies
	TracesRepo *repositories.TracesRepository

	// Processors
	TraceProcessor *logic.TraceProcessor
}

// NewTracesController creates a new traces controller
func NewTracesController(
	tracesRepo *repositories.TracesRepository,
	traceProcessor *logic.TraceProcessor,
) *TracesController {
	return &TracesController{
		TracesRepo:     tracesRepo,
		TraceProcessor: traceProcessor,
	}
}

// QueryTraces retrieves traces based on the provided query
func (c *TracesController) QueryTraces(ctx context.Context, provider string, query *models.TraceQuery) ([]models.TraceSummary, error) {
	// Validate query
	if err := c.validateQuery(query); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	// Query repository
	traces, err := c.TracesRepo.QueryTraces(ctx, provider, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query traces from provider '%s': %w", provider, err)
	}

	// Process/enrich traces with business logic
	if c.TraceProcessor != nil {
		traces = c.TraceProcessor.EnrichTraceSummaries(traces)
	}

	return traces, nil
}

// GetTraceDetail retrieves detailed information for a specific trace
func (c *TracesController) GetTraceDetail(ctx context.Context, provider string, traceID string) (*models.Trace, error) {
	// Validate trace ID
	if traceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	// Get trace detail from repository
	trace, err := c.TracesRepo.GetTraceDetail(ctx, provider, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trace detail from provider '%s': %w", provider, err)
	}

	// Process/enrich trace with business logic
	if c.TraceProcessor != nil {
		trace = c.TraceProcessor.EnrichTrace(trace)
	}

	return trace, nil
}

// GetServices retrieves available services from a specific provider
func (c *TracesController) GetServices(ctx context.Context, provider string) ([]string, error) {
	return c.TracesRepo.GetServices(ctx, provider)
}

// validateQuery validates the trace query
func (c *TracesController) validateQuery(query *models.TraceQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	// Validate time range
	if query.EndTime.Before(query.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	// Validate limit
	if query.Limit < 0 {
		return fmt.Errorf("limit must be non-negative")
	}
	if query.Limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}

	return nil
}
