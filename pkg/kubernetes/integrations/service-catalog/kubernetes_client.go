package servicecatalog

import (
	"context"
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
)

// KubernetesClient handles communication with Kubernetes module for Service Catalog
type KubernetesClient struct {
	deploymentRepo k8sPorts.DeploymentRepository
	clusterRepo    k8sPorts.ClusterRepository
}

// NewKubernetesClient creates a new Kubernetes client for Service Catalog integration
func NewKubernetesClient(deploymentRepo k8sPorts.DeploymentRepository, clusterRepo k8sPorts.ClusterRepository) *KubernetesClient {
	return &KubernetesClient{
		deploymentRepo: deploymentRepo,
		clusterRepo:    clusterRepo,
	}
}

// GetDeploymentHealth gets health information for a deployment
func (c *KubernetesClient) GetDeploymentHealth(ctx context.Context, kubeContext, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	deployment, err := c.deploymentRepo.GetDeployment(ctx, kubeContext, namespace, deploymentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return deployment, nil
}

// ListDeployments lists deployments in a namespace
func (c *KubernetesClient) ListDeployments(ctx context.Context, kubeContext, namespace string) ([]string, error) {
	filter := &k8sModels.DeploymentFilter{
		Namespace: namespace,
	}
	deployments, err := c.deploymentRepo.ListDeployments(ctx, kubeContext, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	var names []string
	for _, deployment := range deployments.Deployments {
		names = append(names, deployment.Name)
	}
	return names, nil
}

// ValidateContext validates if a Kubernetes context is accessible
func (c *KubernetesClient) ValidateContext(ctx context.Context, kubeContext string) error {
	clusters, err := c.clusterRepo.ListClusters(ctx)
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}
	for _, cluster := range clusters {
		if cluster.Context == kubeContext {
			return nil
		}
	}
	return fmt.Errorf("context %s not found", kubeContext)
}
