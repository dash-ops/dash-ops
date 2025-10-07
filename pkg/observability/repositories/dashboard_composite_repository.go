package repositories

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// DashboardCompositeRepository agrega múltiplas fontes (IaC + Grafana)

type DashboardCompositeRepository struct {
	providers []ports.DashboardRepository
}

func NewDashboardCompositeRepository(providers ...ports.DashboardRepository) ports.DashboardRepository {
	return &DashboardCompositeRepository{providers: providers}
}

func (r *DashboardCompositeRepository) GetDashboards(ctx context.Context) ([]models.Dashboard, error) {
	var out []models.Dashboard
	for _, p := range r.providers {
		dash, err := p.GetDashboards(ctx)
		if err != nil {
			continue
		}
		out = append(out, dash...)
	}
	return out, nil
}

func (r *DashboardCompositeRepository) GetDashboard(ctx context.Context, id string) (*models.Dashboard, error) {
	for _, p := range r.providers {
		d, _ := p.GetDashboard(ctx, id)
		if d != nil {
			return d, nil
		}
	}
	return nil, nil
}

func (r *DashboardCompositeRepository) CreateDashboard(ctx context.Context, req *wire.CreateDashboardRequest) (*models.Dashboard, error) {
	// Por padrão, cria apenas no primeiro provider (ex.: IaC)
	if len(r.providers) == 0 {
		return nil, nil
	}
	return r.providers[0].CreateDashboard(ctx, req)
}

func (r *DashboardCompositeRepository) UpdateDashboard(ctx context.Context, id string, req *wire.CreateDashboardRequest) (*models.Dashboard, error) {
	if len(r.providers) == 0 {
		return nil, nil
	}
	return r.providers[0].UpdateDashboard(ctx, id, req)
}

func (r *DashboardCompositeRepository) DeleteDashboard(ctx context.Context, id string) error {
	if len(r.providers) == 0 {
		return nil
	}
	return r.providers[0].DeleteDashboard(ctx, id)
}

func (r *DashboardCompositeRepository) GetDashboardTemplates(ctx context.Context) ([]models.DashboardTemplate, error) {
	// Agrega templates de todos
	var out []models.DashboardTemplate
	for _, p := range r.providers {
		t, err := p.GetDashboardTemplates(ctx)
		if err != nil {
			continue
		}
		out = append(out, t...)
	}
	return out, nil
}
