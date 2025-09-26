package http

import (
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes/wire"
)

// NodesToResponse converts Node models to NodeResponse slice
func NodesToResponse(nodes []k8sModels.Node) []k8sWire.NodeResponse {
	var response []k8sWire.NodeResponse
	for _, node := range nodes {
		response = append(response, *NodeToResponse(&node))
	}
	return response
}

// NodeToResponse converts Node model to NodeResponse
func NodeToResponse(node *k8sModels.Node) *k8sWire.NodeResponse {
	var conditions []k8sWire.NodeConditionResponse
	for _, condition := range node.Conditions {
		conditions = append(conditions, k8sWire.NodeConditionResponse{
			Type:               condition.Type,
			Status:             condition.Status,
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime,
		})
	}

	return &k8sWire.NodeResponse{
		Name:       node.Name,
		Status:     string(node.Status),
		Roles:      node.Roles,
		Age:        node.Age,
		Version:    node.Version,
		InternalIP: node.InternalIP,
		ExternalIP: node.ExternalIP,
		Conditions: conditions,
		Resources: k8sWire.NodeResourcesResponse{
			Capacity: k8sWire.ResourceListResponse{
				CPU:    node.Resources.Capacity.CPU,
				Memory: node.Resources.Capacity.Memory,
				Pods:   node.Resources.Capacity.Pods,
			},
			Allocatable: k8sWire.ResourceListResponse{
				CPU:    node.Resources.Allocatable.CPU,
				Memory: node.Resources.Allocatable.Memory,
				Pods:   node.Resources.Allocatable.Pods,
			},
			Used: k8sWire.ResourceListResponse{
				CPU:    node.Resources.Used.CPU,
				Memory: node.Resources.Used.Memory,
				Pods:   node.Resources.Used.Pods,
			},
		},
		CreatedAt: node.CreatedAt,
	}
}
