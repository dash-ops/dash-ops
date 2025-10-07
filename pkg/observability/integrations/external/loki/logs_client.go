package loki

import (
	"context"
	"fmt"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// --- Implementation of ports.LogsClient interface ---

// QueryLogs implements ports.LogsClient - queries logs using standardized models
func (c *LokiClient) QueryLogs(ctx context.Context, query *models.LogQuery) ([]models.LogEntry, error) {
	// Transform models.LogQuery to wire.LokiQueryParams (provider-specific)
	lokiParams := c.transformToLokiParams(query)

	// Query Loki using existing queryRange method
	lokiResp, err := c.queryRange(ctx, lokiParams)
	if err != nil {
		return nil, fmt.Errorf("loki query failed: %w", err)
	}

	// Transform wire.LokiQueryResponse to []models.LogEntry (provider-specific)
	logs := c.transformLokiResponseToModels(lokiResp)

	return logs, nil
}

// GetLogLabels implements ports.LogsClient - retrieves all available labels
func (c *LokiClient) GetLogLabels(ctx context.Context) ([]string, error) {
	resp, err := c.listLabels(ctx, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetLogLevels implements ports.LogsClient - retrieves all available log levels
func (c *LokiClient) GetLogLevels(ctx context.Context) ([]string, error) {
	resp, err := c.getLabelValues(ctx, "level", time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// HealthCheck checks if Loki is healthy
func (c *LokiClient) HealthCheck(ctx context.Context) error {
	// Try to list labels as a health check
	_, err := c.listLabels(ctx, time.Time{}, time.Time{})
	return err
}
