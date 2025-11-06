package servicecatalog

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	obsPorts "github.com/dash-ops/dash-ops/pkg/observability/ports"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ServiceCatalogAdapter adapts Service Catalog to Observability needs
type ServiceCatalogAdapter struct {
	serviceRepo scPorts.ServiceRepository
}

// NewServiceCatalogAdapter creates a new adapter for Observability integration
func NewServiceCatalogAdapter(serviceRepo scPorts.ServiceRepository) obsPorts.ServiceContextRepository {
	return &ServiceCatalogAdapter{
		serviceRepo: serviceRepo,
	}
}

// GetServiceContext returns service context for observability queries
func (a *ServiceCatalogAdapter) GetServiceContext(ctx context.Context, serviceName string) (*models.ServiceContext, error) {
	service, err := a.serviceRepo.GetByName(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	return a.transformToServiceContext(service), nil
}

// GetServicesWithContext returns all services with their observability context
func (a *ServiceCatalogAdapter) GetServicesWithContext(ctx context.Context) ([]models.ServiceWithContext, error) {
	services, err := a.serviceRepo.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []models.ServiceWithContext
	for _, service := range services {
		serviceContext := a.transformToServiceContext(&service)
		result = append(result, models.ServiceWithContext{
			ServiceContext: *serviceContext,
		})
	}
	return result, nil
}

// GetServiceHealth returns health status for a service
func (a *ServiceCatalogAdapter) GetServiceHealth(ctx context.Context, serviceName string) (*models.ServiceHealth, error) {
	service, err := a.serviceRepo.GetByName(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return &models.ServiceHealth{
		Status:    "unknown",
		LastCheck: service.Metadata.UpdatedAt,
		Details: map[string]interface{}{
			"tier":        string(service.Metadata.Tier),
			"description": service.Spec.Description,
		},
	}, nil
}

// transformToServiceContext transforms ServiceCatalog service to Observability context
func (a *ServiceCatalogAdapter) transformToServiceContext(service *scModels.Service) *models.ServiceContext {
	ctx := &models.ServiceContext{
		ServiceName: service.Metadata.Name,
		Labels: map[string]string{
			"service_name": service.Metadata.Name,
			"tier":         string(service.Metadata.Tier),
		},
		Metadata: map[string]interface{}{
			"description": service.Spec.Description,
			"tier":        string(service.Metadata.Tier),
		},
	}

	if service.Spec.Kubernetes != nil && len(service.Spec.Kubernetes.Environments) > 0 {
		env := service.Spec.Kubernetes.Environments[0]
		ctx.Namespace = env.Namespace
		ctx.Cluster = env.Context
	}

	return ctx
}
