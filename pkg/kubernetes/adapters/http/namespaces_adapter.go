package http

import (
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes/wire"
)

// NamespacesToResponse converts Namespace models to NamespaceResponse slice
func NamespacesToResponse(namespaces []k8sModels.Namespace) []k8sWire.NamespaceResponse {
	var response []k8sWire.NamespaceResponse
	for _, namespace := range namespaces {
		response = append(response, k8sWire.NamespaceResponse{
			Name:      namespace.Name,
			Status:    string(namespace.Status),
			Labels:    namespace.Labels,
			Age:       namespace.Age,
			CreatedAt: namespace.CreatedAt,
		})
	}
	return response
}

// NamespaceToResponse converts Namespace model to NamespaceResponse
func NamespaceToResponse(namespace *k8sModels.Namespace) *k8sWire.NamespaceResponse {
	return &k8sWire.NamespaceResponse{
		Name:      namespace.Name,
		Status:    string(namespace.Status),
		Labels:    namespace.Labels,
		Age:       namespace.Age,
		CreatedAt: namespace.CreatedAt,
	}
}
