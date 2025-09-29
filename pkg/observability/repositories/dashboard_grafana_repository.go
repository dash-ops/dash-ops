package repositories

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/integrations/external/grafana"
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// DashboardGrafanaRepository implementa ports.DashboardRepository usando Grafana

type DashboardGrafanaRepository struct {
	client *grafana.Client
}

func NewDashboardGrafanaRepository(client *grafana.Client) ports.DashboardRepository {
	return &DashboardGrafanaRepository{client: client}
}

func (r *DashboardGrafanaRepository) GetDashboards(ctx context.Context) ([]models.Dashboard, error) {
	_, _ = r.client.ListDashboards(ctx)
	return []models.Dashboard{}, nil
}

func (r *DashboardGrafanaRepository) GetDashboard(ctx context.Context, id string) (*models.Dashboard, error) {
	_, _ = r.client.GetDashboard(ctx, id)
	return nil, nil
}

func (r *DashboardGrafanaRepository) CreateDashboard(ctx context.Context, req *wire.CreateDashboardRequest) (*models.Dashboard, error) {
	// Normalmente criação seria no IaC; mantemos noop aqui
	return &models.Dashboard{ID: "", Name: req.Name, Description: req.Description, Service: req.Service, Charts: req.Charts, Public: req.Public}, nil
}

func (r *DashboardGrafanaRepository) UpdateDashboard(ctx context.Context, id string, req *wire.CreateDashboardRequest) (*models.Dashboard, error) {
	return &models.Dashboard{ID: id, Name: req.Name, Description: req.Description, Service: req.Service, Charts: req.Charts, Public: req.Public}, nil
}

func (r *DashboardGrafanaRepository) DeleteDashboard(ctx context.Context, id string) error {
	return nil
}

func (r *DashboardGrafanaRepository) GetDashboardTemplates(ctx context.Context) ([]models.DashboardTemplate, error) {
	return []models.DashboardTemplate{}, nil
}
