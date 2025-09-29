package repositories

import (
	"context"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/integrations/external/alertmanager"
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// AlertsRepository implements alerts data access using AlertManager
type AlertsRepository struct {
	alertManagerAdapter ports.AlertRepository
}

// NewAlertsRepository creates a new alerts repository
func NewAlertsRepository(alertManagerClient *alertmanager.AlertManagerClient) *AlertsRepository {
	alertManagerAdapter := alertmanager.NewAlertManagerAdapter(alertManagerClient)

	return &AlertsRepository{
		alertManagerAdapter: alertManagerAdapter,
	}
}

// GetAlerts retrieves alerts from the repository
func (r *AlertsRepository) GetAlerts(ctx context.Context, req *wire.AlertsRequest) (*wire.AlertsResponse, error) {
	return r.alertManagerAdapter.GetAlerts(ctx, req)
}

// CreateAlert creates a new alert
func (r *AlertsRepository) CreateAlert(ctx context.Context, req *wire.CreateAlertRequest) (*models.Alert, error) {
	return r.alertManagerAdapter.CreateAlert(ctx, req)
}

// UpdateAlert updates an existing alert
func (r *AlertsRepository) UpdateAlert(ctx context.Context, id string, req *wire.CreateAlertRequest) (*models.Alert, error) {
	return r.alertManagerAdapter.UpdateAlert(ctx, id, req)
}

// DeleteAlert deletes an alert
func (r *AlertsRepository) DeleteAlert(ctx context.Context, id string) error {
	return r.alertManagerAdapter.DeleteAlert(ctx, id)
}

// SilenceAlert silences an alert
func (r *AlertsRepository) SilenceAlert(ctx context.Context, id string, duration time.Duration) error {
	return r.alertManagerAdapter.SilenceAlert(ctx, id, duration)
}

// GetAlertRules retrieves alert rules
func (r *AlertsRepository) GetAlertRules(ctx context.Context) ([]models.AlertRule, error) {
	return r.alertManagerAdapter.GetAlertRules(ctx)
}
