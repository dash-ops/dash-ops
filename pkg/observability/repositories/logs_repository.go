package repositories

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
)

// LogsRepository implements ports.LogRepository with multiple providers support
type LogsRepository struct {
	// Map of provider name to LogsClient implementation
	clients map[string]ports.LogsClient
}

// NewLogsRepository creates a new LogsRepository with multiple clients
func NewLogsRepository(clients map[string]ports.LogsClient) *LogsRepository {
	return &LogsRepository{
		clients: clients,
	}
}

// QueryLogsWithModel retrieves logs from the specified provider
func (r *LogsRepository) QueryLogsWithModel(ctx context.Context, provider string, query *models.LogQuery) ([]models.LogEntry, error) {
	client, exists := r.clients[provider]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", provider)
	}

	// Client handles all provider-specific transformations internally
	return client.QueryLogs(ctx, query)
}

// GetLogLabels returns available log labels from the specified provider
func (r *LogsRepository) GetLogLabels(ctx context.Context, provider string) ([]string, error) {
	client, exists := r.clients[provider]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", provider)
	}

	return client.GetLogLabels(ctx)
}

// GetLogLevels returns available log levels from the specified provider
func (r *LogsRepository) GetLogLevels(ctx context.Context, provider string) ([]string, error) {
	client, exists := r.clients[provider]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", provider)
	}

	return client.GetLogLevels(ctx)
}

// GetAvailableProviders returns list of available provider names
func (r *LogsRepository) GetAvailableProviders() []string {
	providers := make([]string, 0, len(r.clients))
	for provider := range r.clients {
		providers = append(providers, provider)
	}
	return providers
}

// HealthCheck checks health of all providers
func (r *LogsRepository) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)
	for provider, client := range r.clients {
		results[provider] = client.HealthCheck(ctx)
	}
	return results
}
