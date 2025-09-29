package repositories

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// DashboardIaCRepository implementa ports.DashboardRepository lendo definições internas (IaC)
// Nesta versão, funções retornam estruturas vazias; futuras implementações irão carregar de arquivos/YAML

type DashboardIaCRepository struct{}

func NewDashboardIaCRepository() ports.DashboardRepository {
	return &DashboardIaCRepository{}
}

func (d *DashboardIaCRepository) GetDashboards(ctx context.Context) ([]models.Dashboard, error) {
	return []models.Dashboard{}, nil
}

func (d *DashboardIaCRepository) GetDashboard(ctx context.Context, id string) (*models.Dashboard, error) {
	return nil, nil
}

func (d *DashboardIaCRepository) CreateDashboard(ctx context.Context, req *wire.CreateDashboardRequest) (*models.Dashboard, error) {
	return &models.Dashboard{ID: "", Name: req.Name, Description: req.Description, Service: req.Service, Charts: req.Charts, Public: req.Public}, nil
}

func (d *DashboardIaCRepository) UpdateDashboard(ctx context.Context, id string, req *wire.CreateDashboardRequest) (*models.Dashboard, error) {
	return &models.Dashboard{ID: id, Name: req.Name, Description: req.Description, Service: req.Service, Charts: req.Charts, Public: req.Public}, nil
}

func (d *DashboardIaCRepository) DeleteDashboard(ctx context.Context, id string) error {
	return nil
}

func (d *DashboardIaCRepository) GetDashboardTemplates(ctx context.Context) ([]models.DashboardTemplate, error) {
	return []models.DashboardTemplate{}, nil
}
