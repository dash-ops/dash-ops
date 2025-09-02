package kubernetes

import (
	"context"
	"time"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes-new/models"
)

// ClusterRepository defines the interface for cluster data access
type ClusterRepository interface {
	// GetCluster gets cluster information
	GetCluster(ctx context.Context, context string) (*k8sModels.Cluster, error)

	// ListClusters lists all configured clusters
	ListClusters(ctx context.Context) ([]k8sModels.Cluster, error)

	// ValidateCluster validates cluster connectivity
	ValidateCluster(ctx context.Context, context string) error

	// GetClusterInfo gets comprehensive cluster information
	GetClusterInfo(ctx context.Context, context string) (*k8sModels.ClusterInfo, error)
}

// NodeRepository defines the interface for node data access
type NodeRepository interface {
	// GetNode gets a specific node
	GetNode(ctx context.Context, context, nodeName string) (*k8sModels.Node, error)

	// ListNodes lists all nodes in a cluster
	ListNodes(ctx context.Context, context string) ([]k8sModels.Node, error)

	// GetNodeMetrics gets node resource metrics
	GetNodeMetrics(ctx context.Context, context, nodeName string) (*k8sModels.NodeResources, error)
}

// NamespaceRepository defines the interface for namespace data access
type NamespaceRepository interface {
	// GetNamespace gets a specific namespace
	GetNamespace(ctx context.Context, context, namespaceName string) (*k8sModels.Namespace, error)

	// ListNamespaces lists all namespaces in a cluster
	ListNamespaces(ctx context.Context, context string) ([]k8sModels.Namespace, error)

	// CreateNamespace creates a new namespace
	CreateNamespace(ctx context.Context, context, namespaceName string, labels map[string]string) (*k8sModels.Namespace, error)

	// DeleteNamespace deletes a namespace
	DeleteNamespace(ctx context.Context, context, namespaceName string) error
}

// DeploymentRepository defines the interface for deployment data access
type DeploymentRepository interface {
	// GetDeployment gets a specific deployment
	GetDeployment(ctx context.Context, context, namespace, deploymentName string) (*k8sModels.Deployment, error)

	// ListDeployments lists deployments with optional filtering
	ListDeployments(ctx context.Context, context string, filter *k8sModels.DeploymentFilter) (*k8sModels.DeploymentList, error)

	// ScaleDeployment scales a deployment to specified replicas
	ScaleDeployment(ctx context.Context, context, namespace, deploymentName string, replicas int32) error

	// RestartDeployment restarts a deployment
	RestartDeployment(ctx context.Context, context, namespace, deploymentName string) error

	// GetDeploymentStatus gets deployment status and health
	GetDeploymentStatus(ctx context.Context, context, namespace, deploymentName string) (*DeploymentStatus, error)
}

// PodRepository defines the interface for pod data access
type PodRepository interface {
	// GetPod gets a specific pod
	GetPod(ctx context.Context, context, namespace, podName string) (*k8sModels.Pod, error)

	// ListPods lists pods with optional filtering
	ListPods(ctx context.Context, context string, filter *k8sModels.PodFilter) (*k8sModels.PodList, error)

	// DeletePod deletes a pod
	DeletePod(ctx context.Context, context, namespace, podName string) error

	// GetPodLogs gets logs for a pod/container
	GetPodLogs(ctx context.Context, context string, filter *k8sModels.LogFilter) ([]k8sModels.ContainerLog, error)

	// GetPodMetrics gets pod resource metrics
	GetPodMetrics(ctx context.Context, context, namespace, podName string) (*k8sModels.PodMetrics, error)
}

// DeploymentStatus represents deployment status information
type DeploymentStatus struct {
	Name         string                          `json:"name"`
	Namespace    string                          `json:"namespace"`
	Replicas     k8sModels.DeploymentReplicas    `json:"replicas"`
	Conditions   []k8sModels.DeploymentCondition `json:"conditions"`
	HealthStatus string                          `json:"health_status"`
	LastUpdated  time.Time                       `json:"last_updated"`
}
