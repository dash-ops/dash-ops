package kubernetes

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

// Node Struct representing an k8s nodes
type Node struct {
	Name               string                 `json:"name"`
	Ready              v1.ConditionStatus     `json:"ready"`
	AllocatedResources NodeAllocatedResources `json:"allocated_resources"`
	Conditions         []NodeCondition        `json:"conditions"`
	Capacity           NodeCapacity           `json:"capacity"`
	Age                string                 `json:"age"`
	CreatedAt          time.Time              `json:"created_at"`
	Version            string                 `json:"version"`
	Roles              []string               `json:"roles"`
	Taints             int                    `json:"taints"`
}

// NodeCondition represents a node condition
type NodeCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// NodeCapacity represents node capacity information
type NodeCapacity struct {
	Storage          string  `json:"storage,omitempty"`
	EphemeralStorage string  `json:"ephemeral_storage,omitempty"`
	DiskPressure     bool    `json:"disk_pressure"`
	DiskUsagePercent float64 `json:"disk_usage_percent"`
}

// NodeAllocatedResources describes node allocated resources.
type NodeAllocatedResources struct {
	CPURequests            int64   `json:"cpu_requests"`
	CPURequestsFraction    float64 `json:"cpu_requests_fraction"`
	CPULimits              int64   `json:"cpu_limits"`
	CPULimitsFraction      float64 `json:"cpu_limits_fraction"`
	CPUCapacity            int64   `json:"cpu_capacity"`
	MemoryRequests         int64   `json:"memory_requests"`
	MemoryRequestsFraction float64 `json:"memory_requests_fraction"`
	MemoryLimits           int64   `json:"memory_limits"`
	MemoryLimitsFraction   float64 `json:"memory_limits_fraction"`
	MemoryCapacity         int64   `json:"memory_capacity"`
	AllocatedPods          int     `json:"allocated_pods"`
	PodCapacity            int64   `json:"pod_capacity"`
	PodFraction            float64 `json:"pod_fraction"`
}

func (kc client) GetNodes() ([]Node, error) {
	var list []Node
	nodes, err := kc.clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %s", err)
	}

	for _, node := range nodes.Items {
		pods, err := getNodePods(kc.clientSet, node)
		if err != nil {
			return nil, fmt.Errorf("failed to get pods in node: %s", err)
		}

		allocatedResources, err := getNodeAllocatedResources(node, pods)
		if err != nil {
			return nil, fmt.Errorf("failed to get node allocated resources: %s", err)
		}
		conditions := getNodeConditions(node)
		capacity := getNodeCapacity(node, pods)
		age := calculateNodeAge(node.CreationTimestamp.Time)
		version := getNodeVersion(node)
		roles := getNodeRoles(node)
		taints := len(node.Spec.Taints)

		list = append(list, Node{
			Name:               node.GetName(),
			Ready:              getNodeConditionStatus(node, v1.NodeReady),
			AllocatedResources: allocatedResources,
			Conditions:         conditions,
			Capacity:           capacity,
			Age:                age,
			CreatedAt:          node.CreationTimestamp.Time,
			Version:            version,
			Roles:              roles,
			Taints:             taints,
		})
	}

	return list, nil
}

func getNodeConditionStatus(node v1.Node, conditionType v1.NodeConditionType) v1.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return v1.ConditionUnknown
}

func getNodeConditions(node v1.Node) []NodeCondition {
	var conditions []NodeCondition
	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, NodeCondition{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  condition.Reason,
			Message: condition.Message,
		})
	}
	return conditions
}

func getNodeCapacity(node v1.Node, pods *v1.PodList) NodeCapacity {
	capacity := NodeCapacity{}

	if storage := node.Status.Capacity.Storage(); storage != nil {
		capacity.Storage = storage.String()
	}

	if ephemeralStorage := node.Status.Capacity.StorageEphemeral(); ephemeralStorage != nil {
		capacity.EphemeralStorage = ephemeralStorage.String()
	}

	// Check for disk pressure condition
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeDiskPressure && condition.Status == v1.ConditionTrue {
			capacity.DiskPressure = true
			break
		}
	}

	// Estimate disk usage based on pod count and disk pressure
	// This is a rough estimation since we don't have direct disk metrics
	podCount := len(pods.Items)
	diskUsage := float64(podCount) * 2.0 // Estimate ~2% per pod as baseline

	if capacity.DiskPressure {
		diskUsage = 85.0 // If disk pressure, assume high usage
	} else if diskUsage > 75 {
		diskUsage = 75.0 // Cap at 75% if no pressure reported
	}

	capacity.DiskUsagePercent = diskUsage

	return capacity
}

func calculateNodeAge(createdAt time.Time) string {
	now := time.Now()
	duration := now.Sub(createdAt)

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	} else if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

func getNodeVersion(node v1.Node) string {
	return node.Status.NodeInfo.KubeletVersion
}

func getNodeRoles(node v1.Node) []string {
	var roles []string
	for label := range node.Labels {
		if label == "node-role.kubernetes.io/control-plane" || label == "node-role.kubernetes.io/master" {
			roles = append(roles, "control-plane")
		} else if label == "node-role.kubernetes.io/worker" {
			roles = append(roles, "worker")
		}
	}

	// If no specific role found, check for generic role labels
	if len(roles) == 0 {
		for label := range node.Labels {
			if label == "kubernetes.io/role" {
				if role, exists := node.Labels[label]; exists {
					roles = append(roles, role)
				}
			}
		}
	}

	// Default to worker if no role found
	if len(roles) == 0 {
		roles = append(roles, "worker")
	}

	return roles
}

func getNodePods(client *kubernetes.Clientset, node v1.Node) (*v1.PodList, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name +
		",status.phase!=" + string(v1.PodSucceeded) +
		",status.phase!=" + string(v1.PodFailed))

	if err != nil {
		return nil, err
	}

	return client.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})
}

func getNodeAllocatedResources(node v1.Node, podList *v1.PodList) (NodeAllocatedResources, error) {
	reqs, limits := map[v1.ResourceName]resource.Quantity{}, map[v1.ResourceName]resource.Quantity{}

	for _, pod := range podList.Items {
		podReqs, podLimits, err := podRequestsAndLimits(&pod)
		if err != nil {
			return NodeAllocatedResources{}, err
		}
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}

	cpuRequests, cpuLimits, memoryRequests, memoryLimits := reqs[v1.ResourceCPU],
		limits[v1.ResourceCPU], reqs[v1.ResourceMemory], limits[v1.ResourceMemory]

	var cpuRequestsFraction, cpuLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Capacity.Cpu().MilliValue()); capacity > 0 {
		cpuRequestsFraction = float64(cpuRequests.MilliValue()) / capacity * 100
		cpuLimitsFraction = float64(cpuLimits.MilliValue()) / capacity * 100
	}

	var memoryRequestsFraction, memoryLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Capacity.Memory().MilliValue()); capacity > 0 {
		memoryRequestsFraction = float64(memoryRequests.MilliValue()) / capacity * 100
		memoryLimitsFraction = float64(memoryLimits.MilliValue()) / capacity * 100
	}

	var podFraction float64 = 0
	var podCapacity int64 = node.Status.Capacity.Pods().Value()
	if podCapacity > 0 {
		podFraction = float64(len(podList.Items)) / float64(podCapacity) * 100
	}

	return NodeAllocatedResources{
		CPURequests:            cpuRequests.MilliValue(),
		CPURequestsFraction:    cpuRequestsFraction,
		CPULimits:              cpuLimits.MilliValue(),
		CPULimitsFraction:      cpuLimitsFraction,
		CPUCapacity:            node.Status.Capacity.Cpu().MilliValue(),
		MemoryRequests:         memoryRequests.Value(),
		MemoryRequestsFraction: memoryRequestsFraction,
		MemoryLimits:           memoryLimits.Value(),
		MemoryLimitsFraction:   memoryLimitsFraction,
		MemoryCapacity:         node.Status.Capacity.Memory().Value(),
		AllocatedPods:          len(podList.Items),
		PodCapacity:            podCapacity,
		PodFraction:            podFraction,
	}, nil
}
