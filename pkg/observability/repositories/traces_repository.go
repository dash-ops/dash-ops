package repositories

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/integrations/external/tempo"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// TracesRepository implements traces data access using Tempo
type TracesRepository struct {
	tempoAdapter ports.TraceRepository
}

// NewTracesRepository creates a new traces repository
func NewTracesRepository(tempoClient *tempo.TempoClient) *TracesRepository {
	tempoAdapter := tempo.NewTempoAdapter(tempoClient)

	return &TracesRepository{
		tempoAdapter: tempoAdapter,
	}
}

// QueryTraces queries traces from the repository
func (r *TracesRepository) QueryTraces(ctx context.Context, req *wire.TracesRequest) (*wire.TracesResponse, error) {
	return r.tempoAdapter.QueryTraces(ctx, req)
}

// GetTraceDetail retrieves a specific trace by ID
func (r *TracesRepository) GetTraceDetail(ctx context.Context, traceID string) (*wire.TraceDetailResponse, error) {
	return r.tempoAdapter.GetTraceDetail(ctx, traceID)
}

// GetServices retrieves available trace services
func (r *TracesRepository) GetServices(ctx context.Context) ([]string, error) {
	return r.tempoAdapter.GetServices(ctx)
}

// GetOperations retrieves operations for a service
func (r *TracesRepository) GetOperations(ctx context.Context, service string) ([]string, error) {
	return r.tempoAdapter.GetOperations(ctx, service)
}

// GetTraceStatistics retrieves trace statistics
func (r *TracesRepository) GetTraceStatistics(ctx context.Context, req *wire.TraceStatsRequest) (*wire.TraceStatistics, error) {
	return r.tempoAdapter.GetTraceStatistics(ctx, req)
}
