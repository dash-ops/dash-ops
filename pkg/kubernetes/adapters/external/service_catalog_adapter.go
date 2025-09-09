package external

import (
	"context"
	"fmt"
	"time"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ServiceCatalogAdapter implements scPorts.KubernetesService interface
// This adapter bridges the kubernetes module with the service-catalog module
type ServiceCatalogAdapter struct {
	deploymentRepo k8sPorts.DeploymentRepository
	clusterRepo    k8sPorts.ClusterRepository
}

// NewServiceCatalogAdapter creates a new adapter for service-catalog integration
func NewServiceCatalogAdapter(deploymentRepo k8sPorts.DeploymentRepository, clusterRepo k8sPorts.ClusterRepository) scPorts.KubernetesService {
	return &ServiceCatalogAdapter{
		deploymentRepo: deploymentRepo,
		clusterRepo:    clusterRepo,
	}
}

// GetDeploymentHealth gets health information for a deployment
func (a *ServiceCatalogAdapter) GetDeploymentHealth(ctx context.Context, kubeContext, namespace, deploymentName string) (*scModels.DeploymentHealth, error) {
	// Get deployment from repository
	deployment, err := a.deploymentRepo.GetDeployment(ctx, kubeContext, namespace, deploymentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// Convert to service-catalog format
	// Determine service status based on health
	var status scModels.ServiceStatus
	if deployment.Replicas.Ready == deployment.Replicas.Desired && deployment.Replicas.Desired > 0 {
		status = scModels.StatusHealthy
	} else if deployment.Replicas.Ready > 0 {
		status = scModels.StatusDegraded
	} else {
		status = scModels.StatusCritical
	}

	health := &scModels.DeploymentHealth{
		Name:            deployment.Name,
		ReadyReplicas:   int(deployment.Replicas.Ready),
		DesiredReplicas: int(deployment.Replicas.Desired),
		Status:          status,
		LastUpdated:     time.Now(),
	}

	return health, nil
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
		deploymentHealth, err := a.GetDeploymentHealth(ctx, envSpec.Context, envSpec.Namespace, deployment.Name)
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
	// TODO: Implement namespace listing using namespace repository
	return []string{}, nil
}

// ListDeployments lists deployments in a namespace
func (a *ServiceCatalogAdapter) ListDeployments(ctx context.Context, kubeContext, namespace string) ([]string, error) {
	// Create filter for namespace
	filter := &k8sModels.DeploymentFilter{
		Namespace: namespace,
	}

	// Get deployments from repository
	deployments, err := a.deploymentRepo.ListDeployments(ctx, kubeContext, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	// Extract names
	var names []string
	for _, deployment := range deployments.Deployments {
		names = append(names, deployment.Name)
	}

	return names, nil
}

// ValidateContext validates if a Kubernetes context is accessible
func (a *ServiceCatalogAdapter) ValidateContext(ctx context.Context, kubeContext string) error {
	// Try to get clusters to validate context
	clusters, err := a.clusterRepo.ListClusters(ctx)
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}

	// Check if context exists
	for _, cluster := range clusters {
		if cluster.Context == kubeContext {
			return nil
		}
	}

	return fmt.Errorf("context %s not found", kubeContext)
}
