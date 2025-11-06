package servicecatalog

import (
	"context"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ServiceCatalogClient handles communication with Service Catalog module
type ServiceCatalogClient struct {
	serviceRepo scPorts.ServiceRepository
}

// NewServiceCatalogClient creates a new Service Catalog client
func NewServiceCatalogClient(serviceRepo scPorts.ServiceRepository) *ServiceCatalogClient {
	return &ServiceCatalogClient{
		serviceRepo: serviceRepo,
	}
}

// GetService gets a service by name
func (c *ServiceCatalogClient) GetService(ctx context.Context, serviceName string) (*scModels.Service, error) {
	return c.serviceRepo.GetByName(ctx, serviceName)
}

// ListServices lists all services
func (c *ServiceCatalogClient) ListServices(ctx context.Context) ([]scModels.Service, error) {
	return c.serviceRepo.List(ctx, nil)
}
