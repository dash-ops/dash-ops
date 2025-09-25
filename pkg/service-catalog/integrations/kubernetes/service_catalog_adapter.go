package kubernetes

import (
	"context"
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ServiceCatalogAdapter adapts Service Catalog service to Kubernetes needs
type ServiceCatalogAdapter struct {
	client *ServiceCatalogClient
}

// NewServiceCatalogAdapter creates a new adapter for Kubernetes integration
func NewServiceCatalogAdapter(serviceRepo scPorts.ServiceRepository) *ServiceCatalogAdapter {
	return &ServiceCatalogAdapter{
		client: NewServiceCatalogClient(serviceRepo),
	}
}

// ValidateContext validates if a Kubernetes context is accessible
func (a *ServiceCatalogAdapter) ValidateContext(ctx context.Context, kubeContext string) error {
	services, err := a.client.ListServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	for _, service := range services {
		if service.Spec.Kubernetes != nil {
			for _, env := range service.Spec.Kubernetes.Environments {
				if env.Context == kubeContext {
					return nil
				}
			}
		}
	}

	if len(services) == 0 {
		return fmt.Errorf("no services found for context %s", kubeContext)
	}

	return nil
}

// ResolveDeploymentService resolves which service a deployment belongs to
func (a *ServiceCatalogAdapter) ResolveDeploymentService(deploymentName, namespace, kubeContext string) (*k8sModels.ServiceContext, error) {
	// Get services for this namespace
	services, err := a.client.GetServicesByNamespace(context.Background(), namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get services for namespace %s: %w", namespace, err)
	}

	// Look for a service that has this deployment
	for _, service := range services {
		if service.Spec.Kubernetes != nil {
			for _, env := range service.Spec.Kubernetes.Environments {
				if env.Context == kubeContext && env.Namespace == namespace {
					for _, deployment := range env.Resources.Deployments {
						if deployment.Name == deploymentName {
							// Found the service, return service context
							return &k8sModels.ServiceContext{
								ServiceName: service.Metadata.Name,
								ServiceTier: string(service.Metadata.Tier),
								Team:        "unknown", // ServiceMetadata doesn't have Team field
								Environment: env.Name,
								Context:     kubeContext,
								Description: service.Spec.Description,
								Found:       true,
							}, nil
						}
					}
				}
			}
		}
	}

	// If no service found, return nil (not an error)
	return nil, nil
}
