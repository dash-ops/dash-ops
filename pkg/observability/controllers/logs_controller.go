package controllers

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/logic"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LogsController handles log-related use cases
type LogsController struct {
	// Dependencies
	LogRepo     ports.LogRepository
	ServiceRepo ports.ServiceContextRepository
	LogService  ports.LogService
	Cache       ports.CacheService

	// Processors
	LogProcessor *logic.LogProcessor
}

func NewLogsController(
	logRepo ports.LogRepository,
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

// GetLogs retrieves logs based on the provided criteria
func (c *LogsController) GetLogs(ctx context.Context, req *wire.LogsRequest) (*wire.LogsResponse, error) {
	// Call repository (which will call Loki adapter)
	response, err := c.LogRepo.QueryLogs(ctx, req)
	if err != nil {
		return &wire.LogsResponse{
			BaseResponse: wire.BaseResponse{
				Success: false,
				Error:   err.Error(),
			},
		}, err
	}

	// Process logs if processor is available
	if c.LogProcessor != nil && response.Data.Logs != nil {
		// Could add enrichment, filtering, or other processing here
		processedLogs := c.LogProcessor.EnrichLogs(response.Data.Logs)
		response.Data.Logs = processedLogs
	}

	return response, nil
}

// GetLogStatistics retrieves log statistics
func (c *LogsController) GetLogStatistics(ctx context.Context, req *wire.LogStatsRequest) (*wire.LogStatisticsResponse, error) {
	// TODO: Implement log statistics logic
	return nil, nil
}

// GetLogLabels retrieves available log labels
func (c *LogsController) GetLogLabels(ctx context.Context) ([]string, error) {
	return c.LogRepo.GetLogLabels(ctx)
}

// GetLogLevels retrieves available log levels
func (c *LogsController) GetLogLevels(ctx context.Context) ([]string, error) {
	return c.LogRepo.GetLogLevels(ctx)
}
