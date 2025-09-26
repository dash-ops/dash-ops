package controllers

import (
	"context"
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/repositories"
)

// DeploymentsController handles deployments business logic orchestration
type DeploymentsController struct {
	repository             *repositories.DeploymentsRepository
	serviceContextResolver k8sPorts.ServiceContextResolver
}

// NewDeploymentsController creates a new deployments controller
func NewDeploymentsController(repository *repositories.DeploymentsRepository) *DeploymentsController {
	return &DeploymentsController{
		repository: repository,
	}
}

// SetServiceContextResolver sets the service context resolver for enrichment
func (c *DeploymentsController) SetServiceContextResolver(resolver k8sPorts.ServiceContextResolver) {
	c.serviceContextResolver = resolver
}

// GetDeployment gets a specific deployment with business logic validation
func (c *DeploymentsController) GetDeployment(ctx context.Context, context, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if deploymentName == "" {
		return nil, fmt.Errorf("deployment name is required")
	}

	deployment, err := c.repository.GetDeployment(ctx, context, namespace, deploymentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// Enrich with service context if resolver is available
	if c.serviceContextResolver != nil {
		serviceContext, err := c.serviceContextResolver.ResolveDeploymentService(deploymentName, namespace, context)
		if err == nil && serviceContext != nil {
			deployment.ServiceContext = serviceContext
		}
	}

	return deployment, nil
}

// ListDeployments lists deployments with optional filtering and business logic
func (c *DeploymentsController) ListDeployments(ctx context.Context, context string, filter *k8sModels.DeploymentFilter) (*k8sModels.DeploymentList, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	deploymentList, err := c.repository.ListDeployments(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	// Enrich with service context if resolver is available
	if c.serviceContextResolver != nil {
		for i := range deploymentList.Deployments {
			deployment := &deploymentList.Deployments[i]
			serviceContext, err := c.serviceContextResolver.ResolveDeploymentService(deployment.Name, deployment.Namespace, context)
			if err == nil && serviceContext != nil {
				deployment.ServiceContext = serviceContext
			}
		}
	}

	return deploymentList, nil
}

// ScaleDeployment scales a deployment with business logic validation
func (c *DeploymentsController) ScaleDeployment(ctx context.Context, context, namespace, deploymentName string, replicas int32) error {
	if context == "" {
		return fmt.Errorf("context is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if deploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}
	if replicas < 0 {
		return fmt.Errorf("replicas must be non-negative")
	}

	// Business logic: validate scaling limits
	const maxReplicas = 1000 // Reasonable upper limit
	if replicas > maxReplicas {
		return fmt.Errorf("replicas cannot exceed %d", maxReplicas)
	}

	// Verify deployment exists before scaling
	_, err := c.repository.GetDeployment(ctx, context, namespace, deploymentName)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	err = c.repository.ScaleDeployment(ctx, context, namespace, deploymentName, replicas)
	if err != nil {
		return fmt.Errorf("failed to scale deployment: %w", err)
	}

	return nil
}

// RestartDeployment restarts a deployment with business logic validation
func (c *DeploymentsController) RestartDeployment(ctx context.Context, context, namespace, deploymentName string) error {
	if context == "" {
		return fmt.Errorf("context is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if deploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}

	// Verify deployment exists before restarting
	_, err := c.repository.GetDeployment(ctx, context, namespace, deploymentName)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	err = c.repository.RestartDeployment(ctx, context, namespace, deploymentName)
	if err != nil {
		return fmt.Errorf("failed to restart deployment: %w", err)
	}

	return nil
}

// GetDeploymentsSummary provides a summary of deployments in a namespace or cluster
func (c *DeploymentsController) GetDeploymentsSummary(ctx context.Context, context string, namespace string) (*DeploymentsSummary, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	// Create filter for namespace if specified
	filter := &k8sModels.DeploymentFilter{}
	if namespace != "" {
		filter.Namespace = namespace
	}

	deploymentList, err := c.repository.ListDeployments(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments for summary: %w", err)
	}

	summary := &DeploymentsSummary{
		Total:        len(deploymentList.Deployments),
		Healthy:      0,
		Unhealthy:    0,
		Scaling:      0,
		ZeroReplicas: 0,
	}

	for _, deployment := range deploymentList.Deployments {
		if deployment.IsHealthy() {
			summary.Healthy++
		} else if deployment.Replicas.Desired == 0 {
			summary.ZeroReplicas++
		} else if deployment.Replicas.Ready == 0 {
			summary.Unhealthy++
		} else {
			summary.Scaling++
		}
	}

	return summary, nil
}

// DeploymentsSummary represents a summary of deployments
type DeploymentsSummary struct {
	Total        int `json:"total"`
	Healthy      int `json:"healthy"`
	Unhealthy    int `json:"unhealthy"`
	Scaling      int `json:"scaling"`
	ZeroReplicas int `json:"zero_replicas"`
}
