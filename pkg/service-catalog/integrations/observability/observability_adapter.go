package observability

import (
	"context"

	obsModels "github.com/dash-ops/dash-ops/pkg/observability/models"
	obsPorts "github.com/dash-ops/dash-ops/pkg/observability/ports"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ObservabilityAdapter adapts Service Catalog to Observability needs
type ObservabilityAdapter struct {
	serviceRepo scPorts.ServiceRepository
}

// NewObservabilityAdapter creates a new adapter for Observability integration
func NewObservabilityAdapter(serviceRepo scPorts.ServiceRepository) obsPorts.ServiceContextRepository {
	return &ObservabilityAdapter{
		serviceRepo: serviceRepo,
	}
}

// GetServiceContext returns service context for observability queries
func (a *ObservabilityAdapter) GetServiceContext(ctx context.Context, serviceName string) (*obsModels.ServiceContext, error) {
	service, err := a.serviceRepo.GetByName(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	return a.transformToServiceContext(service), nil
}

// GetServicesWithContext returns all services with their observability context
func (a *ObservabilityAdapter) GetServicesWithContext(ctx context.Context) ([]obsModels.ServiceWithContext, error) {
	services, err := a.serviceRepo.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []obsModels.ServiceWithContext
	for _, service := range services {
		serviceContext := a.transformToServiceContext(&service)
		result = append(result, obsModels.ServiceWithContext{
			ServiceContext: *serviceContext,
		})
	}
	return result, nil
}

// GetServiceHealth returns health status for a service
func (a *ObservabilityAdapter) GetServiceHealth(ctx context.Context, serviceName string) (*obsModels.ServiceHealth, error) {
	service, err := a.serviceRepo.GetByName(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return &obsModels.ServiceHealth{
		Status:    "unknown",
		LastCheck: service.Metadata.UpdatedAt,
		Details: map[string]interface{}{
			"tier":        string(service.Metadata.Tier),
			"description": service.Spec.Description,
		},
	}, nil
}

// transformToServiceContext transforms ServiceCatalog service to Observability context
func (a *ObservabilityAdapter) transformToServiceContext(service *scModels.Service) *obsModels.ServiceContext {
	ctx := &obsModels.ServiceContext{
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
