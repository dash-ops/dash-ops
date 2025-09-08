package adapters

import (
	"context"
	"fmt"
	"time"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// KubernetesAdapter adapts Kubernetes service to service-catalog needs
type KubernetesAdapter struct {
	kubernetesService scPorts.KubernetesService
}

// NewKubernetesAdapter creates a new Kubernetes adapter
func NewKubernetesAdapter(kubernetesService scPorts.KubernetesService) *KubernetesAdapter {
	return &KubernetesAdapter{
		kubernetesService: kubernetesService,
	}
}

// GetDeploymentHealth gets health information for a deployment
func (a *KubernetesAdapter) GetDeploymentHealth(ctx context.Context, namespace, deploymentName, kubeContext string) (*scModels.DeploymentHealth, error) {
	// Get deployment health from kubernetes service
	k8sHealth, err := a.kubernetesService.GetDeploymentHealth(ctx, kubeContext, namespace, deploymentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment health: %w", err)
	}

	// Convert to service-catalog format
	// Determine service status based on health
	var status scModels.ServiceStatus
	if k8sHealth.ReadyReplicas == k8sHealth.DesiredReplicas && k8sHealth.DesiredReplicas > 0 {
		status = scModels.StatusHealthy
	} else if k8sHealth.ReadyReplicas > 0 {
		status = scModels.StatusDegraded
	} else {
		status = scModels.StatusCritical
	}

	return &scModels.DeploymentHealth{
		Name:            k8sHealth.Name,
		ReadyReplicas:   int(k8sHealth.ReadyReplicas),
		DesiredReplicas: int(k8sHealth.DesiredReplicas),
		Status:          status,
		LastUpdated:     k8sHealth.LastUpdated,
	}, nil
}

// GetEnvironmentHealth gets health information for all deployments in an environment
func (a *KubernetesAdapter) GetEnvironmentHealth(ctx context.Context, service *scModels.Service, environment string) (*scModels.EnvironmentHealth, error) {
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
func (a *KubernetesAdapter) GetServiceHealth(ctx context.Context, service *scModels.Service) (*scModels.ServiceHealth, error) {
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
func (a *KubernetesAdapter) ListNamespaces(ctx context.Context, kubeContext string) ([]string, error) {
	// This would need to be implemented by the kubernetes service
	// For now, return empty list
	return []string{}, nil
}

// ListDeployments lists deployments in a namespace
func (a *KubernetesAdapter) ListDeployments(ctx context.Context, namespace, kubeContext string) ([]string, error) {
	return a.kubernetesService.ListDeployments(ctx, kubeContext, namespace)
}

// ValidateContext validates if a Kubernetes context is accessible
func (a *KubernetesAdapter) ValidateContext(ctx context.Context, kubeContext string) error {
	return a.kubernetesService.ValidateContext(ctx, kubeContext)
}
