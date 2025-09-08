package kubernetes

import (
	"context"
	"fmt"
	"time"

	k8sLogic "github.com/dash-ops/dash-ops/pkg/kubernetes/logic"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
)

// KubernetesController handles Kubernetes business logic orchestration
type KubernetesController struct {
	clusterRepo      k8sPorts.ClusterRepository
	nodeRepo         k8sPorts.NodeRepository
	namespaceRepo    k8sPorts.NamespaceRepository
	deploymentRepo   k8sPorts.DeploymentRepository
	podRepo          k8sPorts.PodRepository
	healthCalc       *k8sLogic.HealthCalculator
	clusterProcessor *k8sLogic.ClusterProcessor
}

// NewKubernetesController creates a new Kubernetes controller
func NewKubernetesController(
	clusterRepo k8sPorts.ClusterRepository,
	nodeRepo k8sPorts.NodeRepository,
	namespaceRepo k8sPorts.NamespaceRepository,
	deploymentRepo k8sPorts.DeploymentRepository,
	podRepo k8sPorts.PodRepository,
	healthCalc *k8sLogic.HealthCalculator,
) *KubernetesController {
	clusterProcessor := k8sLogic.NewClusterProcessor(clusterRepo)

	return &KubernetesController{
		clusterRepo:      clusterRepo,
		nodeRepo:         nodeRepo,
		namespaceRepo:    namespaceRepo,
		deploymentRepo:   deploymentRepo,
		podRepo:          podRepo,
		healthCalc:       healthCalc,
		clusterProcessor: clusterProcessor,
	}
}

// GetClusterInfo gets comprehensive cluster information
func (kc *KubernetesController) GetClusterInfo(ctx context.Context, context string) (*k8sModels.ClusterInfo, error) {
	// Get cluster basic info
	cluster, err := kc.clusterRepo.GetCluster(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}

	// Get nodes
	nodes, err := kc.nodeRepo.ListNodes(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	// Get namespaces
	namespaces, err := kc.namespaceRepo.ListNamespaces(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	clusterInfo := &k8sModels.ClusterInfo{
		Cluster:     *cluster,
		Nodes:       nodes,
		Namespaces:  namespaces,
		LastUpdated: time.Now(),
	}

	// Calculate summary
	clusterInfo.CalculateSummary()

	return clusterInfo, nil
}

// ListNodes lists all nodes in a cluster
func (kc *KubernetesController) ListNodes(ctx context.Context, context string) ([]k8sModels.Node, error) {
	nodes, err := kc.nodeRepo.ListNodes(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	return nodes, nil
}

// GetNode gets a specific node
func (kc *KubernetesController) GetNode(ctx context.Context, context, nodeName string) (*k8sModels.Node, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	node, err := kc.nodeRepo.GetNode(ctx, context, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	return node, nil
}

// ListNamespaces lists all namespaces in a cluster
func (kc *KubernetesController) ListNamespaces(ctx context.Context, context string) ([]k8sModels.Namespace, error) {
	namespaces, err := kc.namespaceRepo.ListNamespaces(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	return namespaces, nil
}

// GetNamespace gets a specific namespace
func (kc *KubernetesController) GetNamespace(ctx context.Context, context, namespaceName string) (*k8sModels.Namespace, error) {
	if namespaceName == "" {
		return nil, fmt.Errorf("namespace name is required")
	}

	namespace, err := kc.namespaceRepo.GetNamespace(ctx, context, namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	return namespace, nil
}

// ListDeployments lists deployments with optional filtering
func (kc *KubernetesController) ListDeployments(ctx context.Context, context string, filter *k8sModels.DeploymentFilter) (*k8sModels.DeploymentList, error) {
	deploymentList, err := kc.deploymentRepo.ListDeployments(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	// Enrich with health information
	for i := range deploymentList.Deployments {
		healthStatus := kc.healthCalc.CalculateDeploymentHealth(&deploymentList.Deployments[i])
		// Add health status to deployment (would need to extend model)
		_ = healthStatus
	}

	return deploymentList, nil
}

// GetDeployment gets a specific deployment
func (kc *KubernetesController) GetDeployment(ctx context.Context, context, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	if deploymentName == "" {
		return nil, fmt.Errorf("deployment name is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	deployment, err := kc.deploymentRepo.GetDeployment(ctx, context, namespace, deploymentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return deployment, nil
}

// ScaleDeployment scales a deployment to specified replicas
func (kc *KubernetesController) ScaleDeployment(ctx context.Context, context, namespace, deploymentName string, replicas int32) error {
	if deploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if replicas < 0 {
		return fmt.Errorf("replicas must be non-negative")
	}

	err := kc.deploymentRepo.ScaleDeployment(ctx, context, namespace, deploymentName, replicas)
	if err != nil {
		return fmt.Errorf("failed to scale deployment: %w", err)
	}

	return nil
}

// RestartDeployment restarts a deployment
func (kc *KubernetesController) RestartDeployment(ctx context.Context, context, namespace, deploymentName string) error {
	if deploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	err := kc.deploymentRepo.RestartDeployment(ctx, context, namespace, deploymentName)
	if err != nil {
		return fmt.Errorf("failed to restart deployment: %w", err)
	}

	return nil
}

// ListPods lists pods with optional filtering
func (kc *KubernetesController) ListPods(ctx context.Context, context string, filter *k8sModels.PodFilter) (*k8sModels.PodList, error) {
	podList, err := kc.podRepo.ListPods(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	return podList, nil
}

// GetPod gets a specific pod
func (kc *KubernetesController) GetPod(ctx context.Context, context, namespace, podName string) (*k8sModels.Pod, error) {
	if podName == "" {
		return nil, fmt.Errorf("pod name is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	pod, err := kc.podRepo.GetPod(ctx, context, namespace, podName)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	return pod, nil
}

// GetPodLogs gets logs for a pod/container
func (kc *KubernetesController) GetPodLogs(ctx context.Context, context string, filter *k8sModels.LogFilter) ([]k8sModels.ContainerLog, error) {
	if filter.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if filter.PodName == "" {
		return nil, fmt.Errorf("pod name is required")
	}

	logs, err := kc.podRepo.GetPodLogs(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod logs: %w", err)
	}

	return logs, nil
}

// DeletePod deletes a pod
func (kc *KubernetesController) DeletePod(ctx context.Context, context, namespace, podName string) error {
	if podName == "" {
		return fmt.Errorf("pod name is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	err := kc.podRepo.DeletePod(ctx, context, namespace, podName)
	if err != nil {
		return fmt.Errorf("failed to delete pod: %w", err)
	}

	return nil
}

// GetClusterHealth gets overall cluster health
func (kc *KubernetesController) GetClusterHealth(ctx context.Context, context string) (*k8sModels.ClusterHealth, error) {
	// Get cluster info
	clusterInfo, err := kc.GetClusterInfo(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}

	// Calculate health
	healthStatus := kc.healthCalc.CalculateClusterHealth(clusterInfo)

	// Build health response
	var nodeHealths []k8sModels.NodeHealth
	for _, node := range clusterInfo.Nodes {
		nodeHealth := k8sModels.NodeHealth{
			Name:       node.Name,
			Status:     node.Status,
			Conditions: node.Conditions,
			Resources: k8sModels.ResourceHealth{
				CPU: k8sModels.ResourceHealthDetail{
					Used:               0, // Would need actual calculation
					Available:          0,
					Total:              0,
					UtilizationPercent: 0,
					Status:             k8sModels.ResourceStatusHealthy,
				},
				Memory: k8sModels.ResourceHealthDetail{
					Used:               0, // Would need actual calculation
					Available:          0,
					Total:              0,
					UtilizationPercent: 0,
					Status:             k8sModels.ResourceStatusHealthy,
				},
			},
			LastUpdated: time.Now(),
		}
		nodeHealths = append(nodeHealths, nodeHealth)
	}

	return &k8sModels.ClusterHealth{
		Context:     context,
		Status:      k8sModels.HealthStatus(healthStatus),
		Nodes:       nodeHealths,
		Summary:     clusterInfo.Summary,
		LastUpdated: time.Now(),
	}, nil
}

// CreateNamespace creates a new namespace
func (kc *KubernetesController) CreateNamespace(ctx context.Context, context, name string) (*k8sModels.Namespace, error) {
	if name == "" {
		return nil, fmt.Errorf("namespace name is required")
	}

	// Create namespace with empty labels initially
	createdNamespace, err := kc.namespaceRepo.CreateNamespace(ctx, context, name, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace: %w", err)
	}

	return createdNamespace, nil
}

// DeleteNamespace deletes a namespace
func (kc *KubernetesController) DeleteNamespace(ctx context.Context, context, name string) error {
	if name == "" {
		return fmt.Errorf("namespace name is required")
	}

	err := kc.namespaceRepo.DeleteNamespace(ctx, context, name)
	if err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	return nil
}

// ListClusters lists all configured clusters
func (kc *KubernetesController) ListClusters(ctx context.Context) ([]k8sModels.Cluster, error) {
	clusters, err := kc.clusterProcessor.ListClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	return clusters, nil
}
