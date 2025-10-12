package tempo

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// QueryTraces searches for traces based on criteria using standardized models
func (c *TempoClient) QueryTraces(ctx context.Context, query *models.TraceQuery) ([]models.TraceSummary, error) {
	// Transform models.TraceQuery to wire.TempoQueryParams
	params := transformToTempoParams(query)

	// Execute search
	response, err := c.search(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search traces: %w", err)
	}

	// Transform response to models
	summaries := transformSearchResponseToModels(response)

	return summaries, nil
}

// GetTraceDetail retrieves detailed information for a specific trace
func (c *TempoClient) GetTraceDetail(ctx context.Context, traceID string) (*models.Trace, error) {
	// Execute get trace by ID
	response, err := c.getTraceByID(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trace detail: %w", err)
	}

	// Transform response to models
	trace := transformTraceByIDResponseToModel(response)
	if trace == nil {
		return nil, fmt.Errorf("trace not found or empty")
	}

	return trace, nil
}

// GetServices returns available services that have traces
func (c *TempoClient) GetServices(ctx context.Context) ([]string, error) {
	// Get the "service.name" tag values
	response, err := c.getTagValues(ctx, "service.name")
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	return response.TagValues, nil
}

// HealthCheck checks if the Tempo service is healthy
func (c *TempoClient) HealthCheck(ctx context.Context) error {
	// Try to get tags (lightweight operation)
	_, err := c.getTags(ctx)
	if err != nil {
		return fmt.Errorf("tempo health check failed: %w", err)
	}
	return nil
}
