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

// NewKubernetesClient creates a new Kubernetes client for Service Catalog
func NewKubernetesClient(deploymentRepo k8sPorts.DeploymentRepository, clusterRepo k8sPorts.ClusterRepository) *KubernetesClient {
	return &KubernetesClient{
		deploymentRepo: deploymentRepo,
		clusterRepo:    clusterRepo,
	}
}

// GetDeployment gets a deployment from Kubernetes
func (c *KubernetesClient) GetDeployment(ctx context.Context, kubeContext, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	return c.deploymentRepo.GetDeployment(ctx, kubeContext, namespace, deploymentName)
}

// ListDeployments lists deployments in a namespace
func (c *KubernetesClient) ListDeployments(ctx context.Context, kubeContext, namespace string) ([]string, error) {
	// Create filter for namespace
	filter := &k8sModels.DeploymentFilter{
		Namespace: namespace,
	}

	// Get deployments from repository
	deployments, err := c.deploymentRepo.ListDeployments(ctx, kubeContext, filter)
	if err != nil {
		return nil, err
	}

	// Extract names
	var names []string
	for _, deployment := range deployments.Deployments {
		names = append(names, deployment.Name)
	}

	return names, nil
}

// ListClusters lists available clusters
func (c *KubernetesClient) ListClusters(ctx context.Context) ([]k8sModels.Cluster, error) {
	return c.clusterRepo.ListClusters(ctx)
}

// ValidateContext validates if a Kubernetes context is accessible
func (c *KubernetesClient) ValidateContext(ctx context.Context, kubeContext string) error {
	// Try to get clusters to validate context
	clusters, err := c.clusterRepo.ListClusters(ctx)
	if err != nil {
		return err
	}

	// Check if context exists
	for _, cluster := range clusters {
		if cluster.Context == kubeContext {
			return nil
		}
	}

	return fmt.Errorf("context %s not found", kubeContext)
}
