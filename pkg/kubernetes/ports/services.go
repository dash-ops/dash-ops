package kubernetes

import (
	"context"
	"io"
	"time"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

// KubernetesClientService defines the interface for Kubernetes client operations
type KubernetesClientService interface {
	// GetClientset gets a Kubernetes clientset for a specific context
	GetClientset(context string) (KubernetesClientset, error)

	// ValidateContext validates if a context is accessible
	ValidateContext(context string) error

	// ListContexts lists all available contexts
	ListContexts() ([]string, error)

	// GetCurrentContext gets the current active context
	GetCurrentContext() (string, error)

	// SwitchContext switches to a different context
	SwitchContext(context string) error
}

// KubernetesClientset defines the interface for Kubernetes API operations
type KubernetesClientset interface {
	// Node operations
	GetNode(ctx context.Context, nodeName string) (*k8sModels.Node, error)
	ListNodes(ctx context.Context) ([]k8sModels.Node, error)

	// Namespace operations
	GetNamespace(ctx context.Context, namespaceName string) (*k8sModels.Namespace, error)
	ListNamespaces(ctx context.Context) ([]k8sModels.Namespace, error)
	CreateNamespace(ctx context.Context, namespaceName string, labels map[string]string) (*k8sModels.Namespace, error)
	DeleteNamespace(ctx context.Context, namespaceName string) error

	// Deployment operations
	GetDeployment(ctx context.Context, namespace, deploymentName string) (*k8sModels.Deployment, error)
	ListDeployments(ctx context.Context, namespace string) ([]k8sModels.Deployment, error)
	ScaleDeployment(ctx context.Context, namespace, deploymentName string, replicas int32) error
	RestartDeployment(ctx context.Context, namespace, deploymentName string) error

	// Pod operations
	GetPod(ctx context.Context, namespace, podName string) (*k8sModels.Pod, error)
	ListPods(ctx context.Context, namespace string) ([]k8sModels.Pod, error)
	DeletePod(ctx context.Context, namespace, podName string) error
	GetPodLogs(ctx context.Context, namespace, podName, containerName string, options *LogOptions) (io.ReadCloser, error)
}

// LogOptions represents options for getting pod logs
type LogOptions struct {
	Follow       bool       `json:"follow,omitempty"`
	TailLines    *int64     `json:"tail_lines,omitempty"`
	SinceSeconds *int64     `json:"since_seconds,omitempty"`
	SinceTime    *time.Time `json:"since_time,omitempty"`
	Previous     bool       `json:"previous,omitempty"`
	Timestamps   bool       `json:"timestamps,omitempty"`
}

// MetricsService defines the interface for Kubernetes metrics operations
type MetricsService interface {
	// GetNodeMetrics gets resource metrics for a node
	GetNodeMetrics(ctx context.Context, context, nodeName string) (*k8sModels.NodeResources, error)

	// GetPodMetrics gets resource metrics for a pod
	GetPodMetrics(ctx context.Context, context, namespace, podName string) (*PodMetrics, error)

	// GetNamespaceMetrics gets aggregated metrics for a namespace
	GetNamespaceMetrics(ctx context.Context, context, namespace string) (*NamespaceMetrics, error)

	// IsMetricsServerAvailable checks if metrics server is available
	IsMetricsServerAvailable(ctx context.Context, context string) (bool, error)
}

// PodMetrics represents pod resource metrics
type PodMetrics struct {
	PodName     string             `json:"pod_name"`
	Namespace   string             `json:"namespace"`
	Containers  []ContainerMetrics `json:"containers"`
	LastUpdated time.Time          `json:"last_updated"`
}

// ContainerMetrics represents container resource metrics
type ContainerMetrics struct {
	Name      string                       `json:"name"`
	Resources k8sModels.ContainerResources `json:"resources"`
	Usage     k8sModels.ResourceList       `json:"usage"`
}

// NamespaceMetrics represents namespace resource metrics
type NamespaceMetrics struct {
	Namespace   string                 `json:"namespace"`
	TotalPods   int                    `json:"total_pods"`
	RunningPods int                    `json:"running_pods"`
	Resources   k8sModels.ResourceList `json:"resources"`
	LastUpdated time.Time              `json:"last_updated"`
}

// EventService defines the interface for Kubernetes events
type EventService interface {
	// GetEvents gets events for a specific resource
	GetEvents(ctx context.Context, context, namespace, resourceType, resourceName string) ([]Event, error)

	// GetNamespaceEvents gets all events in a namespace
	GetNamespaceEvents(ctx context.Context, context, namespace string) ([]Event, error)

	// WatchEvents watches for new events
	WatchEvents(ctx context.Context, context, namespace string) (<-chan Event, error)
}

// Event represents a Kubernetes event
type Event struct {
	Type      string      `json:"type"`
	Reason    string      `json:"reason"`
	Message   string      `json:"message"`
	Object    EventObject `json:"object"`
	Source    EventSource `json:"source"`
	Count     int32       `json:"count"`
	FirstTime time.Time   `json:"first_time"`
	LastTime  time.Time   `json:"last_time"`
}

// EventObject represents the object that generated the event
type EventObject struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}

// EventSource represents the source of the event
type EventSource struct {
	Component string `json:"component"`
	Host      string `json:"host,omitempty"`
}

// HealthService defines the interface for health monitoring
type HealthService interface {
	// GetClusterHealth gets overall cluster health
	GetClusterHealth(ctx context.Context, context string) (*ClusterHealth, error)

	// GetNamespaceHealth gets health for a specific namespace
	GetNamespaceHealth(ctx context.Context, context, namespace string) (*NamespaceHealth, error)

	// GetDeploymentHealth gets health for a specific deployment
	GetDeploymentHealth(ctx context.Context, context, namespace, deploymentName string) (*DeploymentHealth, error)

	// MonitorHealth starts health monitoring for a cluster
	MonitorHealth(ctx context.Context, context string) (<-chan HealthUpdate, error)
}

// ClusterHealth represents overall cluster health
type ClusterHealth struct {
	Context     string                   `json:"context"`
	Status      k8sModels.ClusterStatus  `json:"status"`
	Nodes       []NodeHealth             `json:"nodes"`
	Namespaces  []NamespaceHealth        `json:"namespaces"`
	Summary     k8sModels.ClusterSummary `json:"summary"`
	LastUpdated time.Time                `json:"last_updated"`
}

// NodeHealth represents node health information
type NodeHealth struct {
	Name        string                    `json:"name"`
	Status      k8sModels.NodeStatus      `json:"status"`
	Conditions  []k8sModels.NodeCondition `json:"conditions"`
	Resources   ResourceHealth            `json:"resources"`
	LastUpdated time.Time                 `json:"last_updated"`
}

// NamespaceHealth represents namespace health information
type NamespaceHealth struct {
	Name        string                    `json:"name"`
	Status      k8sModels.NamespaceStatus `json:"status"`
	Deployments []DeploymentHealth        `json:"deployments"`
	Pods        []PodHealth               `json:"pods"`
	LastUpdated time.Time                 `json:"last_updated"`
}

// DeploymentHealth represents deployment health information
type DeploymentHealth struct {
	Name                string                          `json:"name"`
	Namespace           string                          `json:"namespace"`
	Status              string                          `json:"status"`
	Replicas            k8sModels.DeploymentReplicas    `json:"replicas"`
	Conditions          []k8sModels.DeploymentCondition `json:"conditions"`
	AvailabilityPercent float64                         `json:"availability_percent"`
	LastUpdated         time.Time                       `json:"last_updated"`
}

// PodHealth represents pod health information
type PodHealth struct {
	Name        string              `json:"name"`
	Namespace   string              `json:"namespace"`
	Status      k8sModels.PodStatus `json:"status"`
	Phase       string              `json:"phase"`
	Ready       bool                `json:"ready"`
	Restarts    int32               `json:"restarts"`
	Containers  []ContainerHealth   `json:"containers"`
	LastUpdated time.Time           `json:"last_updated"`
}

// ContainerHealth represents container health information
type ContainerHealth struct {
	Name         string    `json:"name"`
	Ready        bool      `json:"ready"`
	RestartCount int32     `json:"restart_count"`
	State        string    `json:"state"`
	LastUpdated  time.Time `json:"last_updated"`
}

// ResourceHealth represents resource health information
type ResourceHealth struct {
	CPU    ResourceHealthDetail `json:"cpu"`
	Memory ResourceHealthDetail `json:"memory"`
	Pods   ResourceHealthDetail `json:"pods,omitempty"`
}

// ResourceHealthDetail represents detailed resource health
type ResourceHealthDetail struct {
	Used               int64   `json:"used"`
	Available          int64   `json:"available"`
	Total              int64   `json:"total"`
	UtilizationPercent float64 `json:"utilization_percent"`
	Status             string  `json:"status"` // healthy, warning, critical
}

// HealthUpdate represents a health status update
type HealthUpdate struct {
	Type         string      `json:"type"` // cluster, namespace, deployment, pod
	Context      string      `json:"context"`
	Namespace    string      `json:"namespace,omitempty"`
	ResourceName string      `json:"resource_name,omitempty"`
	Status       string      `json:"status"`
	Details      interface{} `json:"details,omitempty"`
	Timestamp    time.Time   `json:"timestamp"`
}

// ServiceContextResolver defines the interface for resolving service context
type ServiceContextResolver interface {
	// ResolveDeploymentService resolves which service a deployment belongs to
	ResolveDeploymentService(deploymentName, namespace, context string) (*k8sModels.ServiceContext, error)
}
