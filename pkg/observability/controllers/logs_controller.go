package controllers

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/observability/logic"
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/repositories"
)

// LogsController handles log-related use cases
type LogsController struct {
	// Dependencies
	LogRepo     *repositories.LogsRepository
	ServiceRepo ports.ServiceContextRepository
	LogService  ports.LogService
	Cache       ports.CacheService

	// Processors
	LogProcessor *logic.LogProcessor
}

func NewLogsController(
	logRepo *repositories.LogsRepository,
	serviceRepo ports.ServiceContextRepository,
	logService ports.LogService,
	cache ports.CacheService,
	logProcessor *logic.LogProcessor,
) *LogsController {
	return &LogsController{
		LogRepo:      logRepo,
		ServiceRepo:  serviceRepo,
		LogService:   logService,
		Cache:        cache,
		LogProcessor: logProcessor,
	}
}

// QueryLogs retrieves logs based on the provided query (works with models, not wire)
func (c *LogsController) QueryLogs(ctx context.Context, provider string, query *models.LogQuery) ([]models.LogEntry, error) {
	// Validate query
	if err := c.validateQuery(query); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	// Query repository - this already returns models.LogEntry
	logs, err := c.LogRepo.QueryLogsWithModel(ctx, provider, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs from provider '%s': %w", provider, err)
	}

	// Process/enrich logs with business logic
	if c.LogProcessor != nil {
		logs = c.LogProcessor.EnrichLogs(logs)
	}

	return logs, nil
}

// GetLogLabels retrieves available log labels from a specific provider
func (c *LogsController) GetLogLabels(ctx context.Context, provider string) ([]string, error) {
	return c.LogRepo.GetLogLabels(ctx, provider)
}

// GetLogLevels retrieves available log levels from a specific provider
func (c *LogsController) GetLogLevels(ctx context.Context, provider string) ([]string, error) {
	return c.LogRepo.GetLogLevels(ctx, provider)
}

// validateQuery validates the log query
func (c *LogsController) validateQuery(query *models.LogQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	if query.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}

	if query.Limit > 10000 {
		return fmt.Errorf("limit cannot exceed 10000")
	}

	if query.StartTime.After(query.EndTime) {
		return fmt.Errorf("start time cannot be after end time")
	}

	return nil
}
