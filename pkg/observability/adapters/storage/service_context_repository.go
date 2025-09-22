package storage

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
)

// ServiceContextRepositoryAdapter implements service context repository using storage
type ServiceContextRepositoryAdapter struct {
	// Add storage implementation here
}

// NewServiceContextRepositoryAdapter creates a new service context repository adapter
func NewServiceContextRepositoryAdapter() ports.ServiceContextRepository {
	return &ServiceContextRepositoryAdapter{}
}

// GetServiceContext implements ServiceContextRepository
func (s *ServiceContextRepositoryAdapter) GetServiceContext(ctx context.Context, serviceName string) (*models.ServiceContext, error) {
	return &models.ServiceContext{ServiceName: serviceName, Labels: map[string]string{}}, nil
}

// GetServicesWithContext implements ServiceContextRepository
func (s *ServiceContextRepositoryAdapter) GetServicesWithContext(ctx context.Context) ([]models.ServiceWithContext, error) {
	return []models.ServiceWithContext{}, nil
}

// GetServiceHealth implements ServiceContextRepository
func (s *ServiceContextRepositoryAdapter) GetServiceHealth(ctx context.Context, serviceName string) (*models.ServiceHealth, error) {
	return &models.ServiceHealth{Status: "healthy"}, nil
}
