package repositories

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dash-ops/dash-ops/pkg/kubernetes/integrations/external/kubernetes"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

// DeploymentsRepository handles deployment-related data access
type DeploymentsRepository struct {
	client *kubernetes.KubernetesClient
}

// NewDeploymentsRepository creates a new deployments repository
func NewDeploymentsRepository(client *kubernetes.KubernetesClient) *DeploymentsRepository {
	return &DeploymentsRepository{
		client: client,
	}
}

// GetDeployment gets a specific deployment
func (r *DeploymentsRepository) GetDeployment(ctx context.Context, context, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	if deploymentName == "" {
		return nil, fmt.Errorf("deployment name is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	deployment, err := r.client.GetDeployment(ctx, namespace, deploymentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s/%s: %w", namespace, deploymentName, err)
	}

	return r.convertDeployment(deployment), nil
}

// ListDeployments lists deployments with optional filtering
func (r *DeploymentsRepository) ListDeployments(ctx context.Context, context string, filter *k8sModels.DeploymentFilter) (*k8sModels.DeploymentList, error) {
	// Build list options based on filter
	listOptions := metav1.ListOptions{}
	if filter != nil && filter.LabelSelector != "" {
		listOptions.LabelSelector = filter.LabelSelector
	}

	var deployments []k8sModels.Deployment

	if filter != nil && filter.Namespace != "" {
		// List deployments in specific namespace
		deploymentList, err := r.client.ListDeployments(ctx, filter.Namespace, listOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", filter.Namespace, err)
		}

		for _, deployment := range deploymentList.Items {
			deployments = append(deployments, *r.convertDeployment(&deployment))
		}
	} else {
		// For cross-namespace listing, we need to list all namespaces first
		// and then get deployments from each namespace
		namespaces, err := r.client.ListNamespaces(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}

		for _, namespace := range namespaces.Items {
			deploymentList, err := r.client.ListDeployments(ctx, namespace.Name, listOptions)
			if err != nil {
				// Log error but continue with other namespaces
				continue
			}

			for _, deployment := range deploymentList.Items {
				deployments = append(deployments, *r.convertDeployment(&deployment))
			}
		}
	}

	// Apply additional filters
	if filter != nil {
		deployments = r.filterDeployments(deployments, filter)
	}

	return &k8sModels.DeploymentList{
		Deployments: deployments,
		Total:       len(deployments),
	}, nil
}

// ScaleDeployment scales a deployment to specified replicas
func (r *DeploymentsRepository) ScaleDeployment(ctx context.Context, context, namespace, deploymentName string, replicas int32) error {
	if deploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if replicas < 0 {
		return fmt.Errorf("replicas must be non-negative")
	}

	err := r.client.ScaleDeployment(ctx, namespace, deploymentName, replicas)
	if err != nil {
		return fmt.Errorf("failed to scale deployment %s/%s to %d replicas: %w", namespace, deploymentName, replicas, err)
	}

	return nil
}

// RestartDeployment restarts a deployment
func (r *DeploymentsRepository) RestartDeployment(ctx context.Context, context, namespace, deploymentName string) error {
	if deploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	err := r.client.RestartDeployment(ctx, namespace, deploymentName)
	if err != nil {
		return fmt.Errorf("failed to restart deployment %s/%s: %w", namespace, deploymentName, err)
	}

	return nil
}

// filterDeployments applies additional filters to the deployment list
func (r *DeploymentsRepository) filterDeployments(deployments []k8sModels.Deployment, filter *k8sModels.DeploymentFilter) []k8sModels.Deployment {
	var filtered []k8sModels.Deployment

	for _, deployment := range deployments {
		// Apply service name filter
		if filter.ServiceName != "" && deployment.ServiceContext != nil {
			if deployment.ServiceContext.ServiceName != filter.ServiceName {
				continue
			}
		}

		// Apply status filter
		if filter.Status != "" {
			status := r.getDeploymentStatus(&deployment)
			if status != filter.Status {
				continue
			}
		}

		// Apply limit if specified
		if filter.Limit > 0 && len(filtered) >= filter.Limit {
			break
		}

		filtered = append(filtered, deployment)
	}

	return filtered
}

// getDeploymentStatus determines the status of a deployment
func (r *DeploymentsRepository) getDeploymentStatus(deployment *k8sModels.Deployment) string {
	if deployment.Replicas.Ready == deployment.Replicas.Desired && deployment.Replicas.Desired > 0 {
		return "healthy"
	}
	if deployment.Replicas.Ready == 0 && deployment.Replicas.Desired > 0 {
		return "unhealthy"
	}
	return "scaling"
}

// convertDeployment converts a Kubernetes deployment to our domain model
func (r *DeploymentsRepository) convertDeployment(deployment *appsv1.Deployment) *k8sModels.Deployment {
	// Convert pod info
	podInfo := k8sModels.PodInfo{
		Running: int(deployment.Status.ReadyReplicas),
		Pending: int(deployment.Status.Replicas - deployment.Status.ReadyReplicas),
		Failed:  int(deployment.Status.Replicas - deployment.Status.ReadyReplicas - deployment.Status.ReadyReplicas),
		Total:   int(deployment.Status.Replicas),
	}

	// Convert replicas
	replicas := k8sModels.DeploymentReplicas{
		Desired:   *deployment.Spec.Replicas,
		Current:   deployment.Status.Replicas,
		Ready:     deployment.Status.ReadyReplicas,
		Available: deployment.Status.AvailableReplicas,
	}

	// Convert conditions
	var conditions []k8sModels.DeploymentCondition
	for _, condition := range deployment.Status.Conditions {
		conditions = append(conditions, k8sModels.DeploymentCondition{
			Type:           string(condition.Type),
			Status:         string(condition.Status),
			Reason:         condition.Reason,
			Message:        condition.Message,
			LastUpdateTime: condition.LastUpdateTime.Time,
		})
	}

	return &k8sModels.Deployment{
		Name:       deployment.Name,
		Namespace:  deployment.Namespace,
		PodInfo:    podInfo,
		Replicas:   replicas,
		Age:        time.Since(deployment.CreationTimestamp.Time).Round(time.Second).String(),
		CreatedAt:  deployment.CreationTimestamp.Time,
		Conditions: conditions,
		// ServiceContext will be populated by the controller if service-catalog integration is available
	}
}
