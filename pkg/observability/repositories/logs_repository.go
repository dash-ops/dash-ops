package repositories

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/integrations/external/loki"
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// LogsRepository implements log data access using Loki
type LogsRepository struct {
	lokiAdapter ports.LogRepository
}

// NewLogsRepository creates a new logs repository
func NewLogsRepository(lokiClient *loki.LokiClient) *LogsRepository {
	lokiAdapter := loki.NewLokiAdapter(lokiClient)

	return &LogsRepository{
		lokiAdapter: lokiAdapter,
	}
}

// QueryLogs queries logs from the repository
func (r *LogsRepository) QueryLogs(ctx context.Context, req *wire.LogsRequest) (*wire.LogsResponse, error) {
	return r.lokiAdapter.QueryLogs(ctx, req)
}

// StreamLogs streams logs from the repository
func (r *LogsRepository) StreamLogs(ctx context.Context, req *wire.LogsRequest) (<-chan models.LogEntry, error) {
	return r.lokiAdapter.StreamLogs(ctx, req)
}

// GetLogLabels retrieves available log labels
func (r *LogsRepository) GetLogLabels(ctx context.Context) ([]string, error) {
	return r.lokiAdapter.GetLogLabels(ctx)
}

// GetLogLevels retrieves available log levels
func (r *LogsRepository) GetLogLevels(ctx context.Context) ([]string, error) {
	return r.lokiAdapter.GetLogLevels(ctx)
}
