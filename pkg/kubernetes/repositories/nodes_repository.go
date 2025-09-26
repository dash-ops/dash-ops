package repositories

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dash-ops/dash-ops/pkg/kubernetes/integrations/external/kubernetes"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

// NodesRepository handles node-related data access
type NodesRepository struct {
	client *kubernetes.KubernetesClient
}

// NewNodesRepository creates a new nodes repository
func NewNodesRepository(client *kubernetes.KubernetesClient) *NodesRepository {
	return &NodesRepository{
		client: client,
	}
}

// GetNode gets a specific node
func (r *NodesRepository) GetNode(ctx context.Context, context, nodeName string) (*k8sModels.Node, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	node, err := r.client.GetNode(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}

	return r.convertNode(ctx, node), nil
}

// ListNodes lists all nodes in a cluster
func (r *NodesRepository) ListNodes(ctx context.Context, context string) ([]k8sModels.Node, error) {
	nodeList, err := r.client.ListNodes(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	var nodes []k8sModels.Node
	for _, node := range nodeList.Items {
		nodes = append(nodes, *r.convertNode(ctx, &node))
	}

	return nodes, nil
}

// GetNodeMetrics gets node resource metrics
func (r *NodesRepository) GetNodeMetrics(ctx context.Context, context, nodeName string) (*k8sModels.NodeResources, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	node, err := r.client.GetNode(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics for %s: %w", nodeName, err)
	}

	// Extract resource information from node status
	resources := &k8sModels.NodeResources{
		Capacity: k8sModels.ResourceList{
			CPU:    node.Status.Capacity.Cpu().String(),
			Memory: node.Status.Capacity.Memory().String(),
			Pods:   node.Status.Capacity.Pods().String(),
		},
		Allocatable: k8sModels.ResourceList{
			CPU:    node.Status.Allocatable.Cpu().String(),
			Memory: node.Status.Allocatable.Memory().String(),
			Pods:   node.Status.Allocatable.Pods().String(),
		},
	}

	// TODO: Add actual usage metrics when metrics server is available
	// For now, we'll return empty used resources
	resources.Used = k8sModels.ResourceList{}

	return resources, nil
}

// convertNode converts a Kubernetes node to our domain model
func (r *NodesRepository) convertNode(ctx context.Context, node *corev1.Node) *k8sModels.Node {
	// Determine node status
	var status k8sModels.NodeStatus
	for _, condition := range node.Status.Conditions {
		if condition.Type == "Ready" {
			if condition.Status == "True" {
				status = k8sModels.NodeStatusReady
			} else {
				status = k8sModels.NodeStatusNotReady
			}
			break
		}
	}

	// Extract roles from labels
	var roles []string
	if _, exists := node.Labels["node-role.kubernetes.io/control-plane"]; exists {
		roles = append(roles, "control-plane")
	}
	if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
		roles = append(roles, "master")
	}
	if _, exists := node.Labels["node-role.kubernetes.io/worker"]; exists {
		roles = append(roles, "worker")
	}

	// Extract IP addresses
	var internalIP, externalIP string
	for _, address := range node.Status.Addresses {
		switch address.Type {
		case "InternalIP":
			internalIP = address.Address
		case "ExternalIP":
			externalIP = address.Address
		}
	}

	// Convert conditions
	var conditions []k8sModels.NodeCondition
	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, k8sModels.NodeCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime.Time,
		})
	}

	// Build resources
	resources := k8sModels.NodeResources{
		Capacity: k8sModels.ResourceList{
			CPU:    node.Status.Capacity.Cpu().String(),
			Memory: node.Status.Capacity.Memory().String(),
			Pods:   node.Status.Capacity.Pods().String(),
		},
		Allocatable: k8sModels.ResourceList{
			CPU:    node.Status.Allocatable.Cpu().String(),
			Memory: node.Status.Allocatable.Memory().String(),
			Pods:   node.Status.Allocatable.Pods().String(),
		},
		// Used resources would require metrics-server to be installed
		// For now, we can at least count the pods
		Used: k8sModels.ResourceList{
			CPU:    "0",
			Memory: "0",
			Pods:   "0",
		},
	}

	// Try to get pod count and calculate resource usage on this node
	pods, err := r.client.ListPods(ctx, "", metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name),
	})
	if err == nil {
		resources.Used.Pods = fmt.Sprintf("%d", len(pods.Items))

		// Calculate estimated CPU and Memory usage from pod requests
		var totalCPURequest, totalMemoryRequest int64
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				if container.Resources.Requests != nil {
					// CPU request in millicores
					if cpuRequest := container.Resources.Requests.Cpu(); cpuRequest != nil {
						totalCPURequest += cpuRequest.MilliValue()
					}
					// Memory request in bytes
					if memoryRequest := container.Resources.Requests.Memory(); memoryRequest != nil {
						totalMemoryRequest += memoryRequest.Value()
					}
				}
			}
		}

		// Convert to string format
		if totalCPURequest > 0 {
			resources.Used.CPU = fmt.Sprintf("%dm", totalCPURequest)
		}
		if totalMemoryRequest > 0 {
			// Convert bytes to Mi
			memoryMi := totalMemoryRequest / (1024 * 1024)
			resources.Used.Memory = fmt.Sprintf("%dMi", memoryMi)
		}
	}

	return &k8sModels.Node{
		Name:       node.Name,
		Status:     status,
		Roles:      roles,
		Age:        time.Since(node.CreationTimestamp.Time).Round(time.Second).String(),
		Version:    node.Status.NodeInfo.KubeletVersion,
		InternalIP: internalIP,
		ExternalIP: externalIP,
		Conditions: conditions,
		Resources:  resources,
		CreatedAt:  node.CreationTimestamp.Time,
	}
}
