package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
)

// ServiceContextController handles service context operations
type ServiceContextController struct {
	serviceRepo ports.ServiceContextRepository
}

// NewServiceContextController creates a new service context controller
func NewServiceContextController(serviceRepo ports.ServiceContextRepository) *ServiceContextController {
	return &ServiceContextController{
		serviceRepo: serviceRepo,
	}
}

// GetServicesWithContext retrieves all services with context (paginated and searchable)
func (c *ServiceContextController) GetServicesWithContext(ctx context.Context,
	search string, limit, offset int) ([]models.ServiceWithContext, int, error) {

	if c.serviceRepo == nil {
		return []models.ServiceWithContext{}, 0, fmt.Errorf("service catalog integration not available")
	}

	// Get all services from Service Catalog
	services, err := c.serviceRepo.GetServicesWithContext(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get services: %w", err)
	}

	// Apply search filter (case-insensitive)
	if search != "" {
		filtered := make([]models.ServiceWithContext, 0)
		searchLower := strings.ToLower(search)
		for _, service := range services {
			if strings.Contains(strings.ToLower(service.ServiceName), searchLower) ||
				strings.Contains(strings.ToLower(service.Namespace), searchLower) ||
				strings.Contains(strings.ToLower(service.Cluster), searchLower) {
				filtered = append(filtered, service)
			}
		}
		services = filtered
	}

	// Apply pagination
	totalFiltered := len(services)
	if offset >= totalFiltered {
		return []models.ServiceWithContext{}, totalFiltered, nil
	}

	end := offset + limit
	if end > totalFiltered {
		end = totalFiltered
	}

	return services[offset:end], totalFiltered, nil
}
