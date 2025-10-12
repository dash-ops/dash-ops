package repositories

import (
	"context"
	"fmt"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
)

// TracesRepository handles trace operations across multiple providers
type TracesRepository struct {
	clients map[string]ports.TracesClient
}

// NewTracesRepository creates a new traces repository with multiple clients
func NewTracesRepository(clients map[string]ports.TracesClient) *TracesRepository {
	return &TracesRepository{
		clients: clients,
	}
}

// QueryTraces retrieves traces from a specific provider using standardized models
func (r *TracesRepository) QueryTraces(ctx context.Context, provider string, query *models.TraceQuery) ([]models.TraceSummary, error) {
	client, ok := r.clients[provider]
	if !ok {
		return nil, fmt.Errorf("traces provider '%s' not found", provider)
	}

	return client.QueryTraces(ctx, query)
}

// GetTraceDetail retrieves detailed information for a specific trace from a provider
func (r *TracesRepository) GetTraceDetail(ctx context.Context, provider string, traceID string) (*models.Trace, error) {
	client, ok := r.clients[provider]
	if !ok {
		return nil, fmt.Errorf("traces provider '%s' not found", provider)
	}

	return client.GetTraceDetail(ctx, traceID)
}

// GetServices retrieves available services from a specific provider
func (r *TracesRepository) GetServices(ctx context.Context, provider string) ([]string, error) {
	client, ok := r.clients[provider]
	if !ok {
		return nil, fmt.Errorf("traces provider '%s' not found", provider)
	}

	return client.GetServices(ctx)
}

// HealthCheck checks if a specific provider is healthy
func (r *TracesRepository) HealthCheck(ctx context.Context, provider string) error {
	client, ok := r.clients[provider]
	if !ok {
		return fmt.Errorf("traces provider '%s' not found", provider)
	}

	return client.HealthCheck(ctx)
}
