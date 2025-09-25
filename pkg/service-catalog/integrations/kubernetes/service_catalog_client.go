package kubernetes

import (
	"context"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ServiceCatalogClient handles communication with Service Catalog module for Kubernetes
type ServiceCatalogClient struct {
	serviceRepo scPorts.ServiceRepository
}

// NewServiceCatalogClient creates a new Service Catalog client for Kubernetes
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

// GetServicesByNamespace gets services that have deployments in a specific namespace
func (c *ServiceCatalogClient) GetServicesByNamespace(ctx context.Context, namespace string) ([]scModels.Service, error) {
	// Get all services
	services, err := c.serviceRepo.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Filter services that have Kubernetes deployments in the specified namespace
	var filteredServices []scModels.Service
	for _, service := range services {
		if service.Spec.Kubernetes != nil {
			for _, env := range service.Spec.Kubernetes.Environments {
				if env.Namespace == namespace {
					filteredServices = append(filteredServices, service)
					break
				}
			}
		}
	}

	return filteredServices, nil
}

// GetServicesByContext gets services that have deployments in a specific Kubernetes context
func (c *ServiceCatalogClient) GetServicesByContext(ctx context.Context, kubeContext string) ([]scModels.Service, error) {
	// Get all services
	services, err := c.serviceRepo.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Filter services that have Kubernetes deployments in the specified context
	var filteredServices []scModels.Service
	for _, service := range services {
		if service.Spec.Kubernetes != nil {
			for _, env := range service.Spec.Kubernetes.Environments {
				if env.Context == kubeContext {
					filteredServices = append(filteredServices, service)
					break
				}
			}
		}
	}

	return filteredServices, nil
}
