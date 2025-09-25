package kubernetes

import (
	"context"
	"fmt"
	"time"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ServiceCatalogAdapter adapts Service Catalog service to Kubernetes needs
type ServiceCatalogAdapter struct {
	client *ServiceCatalogClient
}

// NewServiceCatalogAdapter creates a new Service Catalog adapter
func NewServiceCatalogAdapter(serviceRepo scPorts.ServiceRepository) *ServiceCatalogAdapter {
	return &ServiceCatalogAdapter{
		client: NewServiceCatalogClient(serviceRepo),
	}
}

// GetDeploymentHealth gets health information for a deployment
func (a *ServiceCatalogAdapter) GetDeploymentHealth(ctx context.Context, namespace, deploymentName, kubeContext string) (*scModels.DeploymentHealth, error) {
	// This method would typically be called by Kubernetes module
	// to get health information from Service Catalog perspective
	// For now, we'll return a basic health status

	// Try to find a service that has this deployment
	services, err := a.client.GetServicesByNamespace(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get services for namespace: %w", err)
	}

	// Look for a service that has this deployment
	for _, service := range services {
		if service.Spec.Kubernetes != nil {
			for _, env := range service.Spec.Kubernetes.Environments {
				if env.Context == kubeContext && env.Namespace == namespace {
					for _, deployment := range env.Resources.Deployments {
						if deployment.Name == deploymentName {
							// Found the service, return basic health info
							return &scModels.DeploymentHealth{
								Name:            deploymentName,
								ReadyReplicas:   0, // This would be filled by Kubernetes module
								DesiredReplicas: 0, // This would be filled by Kubernetes module
								Status:          scModels.StatusUnknown,
								LastUpdated:     time.Now(),
							}, nil
						}
					}
				}
			}
		}
	}

	// If no service found, return unknown status
	return &scModels.DeploymentHealth{
		Name:            deploymentName,
		ReadyReplicas:   0,
		DesiredReplicas: 0,
		Status:          scModels.StatusUnknown,
		LastUpdated:     time.Now(),
	}, nil
}

// GetEnvironmentHealth gets health information for all deployments in an environment
func (a *ServiceCatalogAdapter) GetEnvironmentHealth(ctx context.Context, service *scModels.Service, environment string) (*scModels.EnvironmentHealth, error) {
	if service.Spec.Kubernetes == nil {
		return &scModels.EnvironmentHealth{
			Name:        environment,
			Context:     "",
			Status:      scModels.StatusUnknown,
			Deployments: []scModels.DeploymentHealth{},
		}, nil
	}

	// Find the environment in service spec
	var envSpec *scModels.KubernetesEnvironment
	for _, env := range service.Spec.Kubernetes.Environments {
		if env.Name == environment {
			envSpec = &env
			break
		}
	}

	if envSpec == nil {
		return &scModels.EnvironmentHealth{
			Name:        environment,
			Context:     "",
			Status:      scModels.StatusUnknown,
			Deployments: []scModels.DeploymentHealth{},
		}, nil
	}

	// Get health for each deployment in this environment
	var deployments []scModels.DeploymentHealth
	overallStatus := scModels.StatusHealthy

	for _, deployment := range envSpec.Resources.Deployments {
		deploymentHealth, err := a.GetDeploymentHealth(ctx, envSpec.Namespace, deployment.Name, envSpec.Context)
		if err != nil {
			// If we can't get health for a deployment, mark as unknown
			deploymentHealth = &scModels.DeploymentHealth{
				Name:            deployment.Name,
				ReadyReplicas:   0,
				DesiredReplicas: 0,
				Status:          scModels.StatusUnknown,
				LastUpdated:     time.Now(),
			}
			overallStatus = scModels.StatusUnknown
		}

		deployments = append(deployments, *deploymentHealth)

		// Update overall status based on deployment status
		if deploymentHealth.Status == scModels.StatusCritical {
			overallStatus = scModels.StatusCritical
		} else if deploymentHealth.Status == scModels.StatusUnknown && overallStatus != scModels.StatusCritical {
			overallStatus = scModels.StatusUnknown
		}
	}

	return &scModels.EnvironmentHealth{
		Name:        environment,
		Context:     envSpec.Context,
		Status:      overallStatus,
		Deployments: deployments,
	}, nil
}

// GetServiceHealth gets aggregated health information for a service
func (a *ServiceCatalogAdapter) GetServiceHealth(ctx context.Context, service *scModels.Service) (*scModels.ServiceHealth, error) {
	if service.Spec.Kubernetes == nil {
		return &scModels.ServiceHealth{
			ServiceName:   service.Metadata.Name,
			OverallStatus: scModels.StatusUnknown,
			Environments:  []scModels.EnvironmentHealth{},
			LastUpdated:   time.Now(),
		}, nil
	}

	// Get health for each environment
	var environments []scModels.EnvironmentHealth
	overallStatus := scModels.StatusHealthy

	for _, env := range service.Spec.Kubernetes.Environments {
		envHealth, err := a.GetEnvironmentHealth(ctx, service, env.Name)
		if err != nil {
			// If we can't get health for an environment, mark as unknown
			envHealth = &scModels.EnvironmentHealth{
				Name:        env.Name,
				Context:     env.Context,
				Status:      scModels.StatusUnknown,
				Deployments: []scModels.DeploymentHealth{},
			}
			overallStatus = scModels.StatusUnknown
		}

		environments = append(environments, *envHealth)

		// Update overall status based on environment status
		if envHealth.Status == scModels.StatusCritical {
			overallStatus = scModels.StatusCritical
		} else if envHealth.Status == scModels.StatusUnknown && overallStatus != scModels.StatusCritical {
			overallStatus = scModels.StatusUnknown
		}
	}

	return &scModels.ServiceHealth{
		ServiceName:   service.Metadata.Name,
		OverallStatus: overallStatus,
		Environments:  environments,
		LastUpdated:   time.Now(),
	}, nil
}

// ListNamespaces lists available namespaces in a context
func (a *ServiceCatalogAdapter) ListNamespaces(ctx context.Context, kubeContext string) ([]string, error) {
	// Get services for this context
	services, err := a.client.GetServicesByContext(ctx, kubeContext)
	if err != nil {
		return nil, err
	}

	// Extract unique namespaces
	namespaceMap := make(map[string]bool)
	for _, service := range services {
		if service.Spec.Kubernetes != nil {
			for _, env := range service.Spec.Kubernetes.Environments {
				if env.Context == kubeContext {
					namespaceMap[env.Namespace] = true
				}
			}
		}
	}

	// Convert map to slice
	var namespaces []string
	for namespace := range namespaceMap {
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

// ListDeployments lists deployments in a namespace
func (a *ServiceCatalogAdapter) ListDeployments(ctx context.Context, namespace, kubeContext string) ([]string, error) {
	// Get services for this namespace
	services, err := a.client.GetServicesByNamespace(ctx, namespace)
	if err != nil {
		return nil, err
	}

	// Extract deployment names
	var deployments []string
	for _, service := range services {
		if service.Spec.Kubernetes != nil {
			for _, env := range service.Spec.Kubernetes.Environments {
				if env.Context == kubeContext && env.Namespace == namespace {
					for _, deployment := range env.Resources.Deployments {
						deployments = append(deployments, deployment.Name)
					}
				}
			}
		}
	}

	return deployments, nil
}

// ValidateContext validates if a Kubernetes context is accessible
func (a *ServiceCatalogAdapter) ValidateContext(ctx context.Context, kubeContext string) error {
	// Check if we have any services using this context
	services, err := a.client.GetServicesByContext(ctx, kubeContext)
	if err != nil {
		return err
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
