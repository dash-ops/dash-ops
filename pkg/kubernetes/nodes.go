package kubernetes

import (
	"fmt"

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

	nodes, err := kc.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get nodes: %s", err)
	}

	for _, node := range nodes.Items {
		pods, err := getNodePods(kc.clientSet, node)
		if err != nil {
			return nil, fmt.Errorf("Failed to get pods in node: %s", err)
		}

		allocatedResources, err := getNodeAllocatedResources(node, pods)
		list = append(list, Node{
			Name:               node.GetName(),
			Ready:              getNodeConditionStatus(node, v1.NodeReady),
			AllocatedResources: allocatedResources,
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

func getNodePods(client *kubernetes.Clientset, node v1.Node) (*v1.PodList, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name +
		",status.phase!=" + string(v1.PodSucceeded) +
		",status.phase!=" + string(v1.PodFailed))

	if err != nil {
		return nil, err
	}

	return client.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{
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
