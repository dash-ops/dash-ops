package http

import (
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes/wire"
)

// ClusterInfoToResponse converts ClusterInfo model to ClusterInfoResponse
func ClusterInfoToResponse(clusterInfo *k8sModels.ClusterInfo) k8sWire.ClusterInfoResponse {
	return k8sWire.ClusterInfoResponse{
		Cluster: k8sWire.ClusterResponse{
			Name:    clusterInfo.Cluster.Name,
			Context: clusterInfo.Cluster.Context,
			Server:  clusterInfo.Cluster.Server,
			Version: clusterInfo.Cluster.Version,
			Status:  string(clusterInfo.Cluster.Status),
		},
		Nodes:       NodesToResponse(clusterInfo.Nodes),
		Namespaces:  NamespacesToResponse(clusterInfo.Namespaces),
		Summary:     ClusterSummaryToResponse(&clusterInfo.Summary),
		LastUpdated: clusterInfo.LastUpdated,
	}
}

// ClusterSummaryToResponse converts ClusterSummary model to ClusterSummaryResponse
func ClusterSummaryToResponse(summary *k8sModels.ClusterSummary) k8sWire.ClusterSummaryResponse {
	return k8sWire.ClusterSummaryResponse{
		TotalNodes:       summary.TotalNodes,
		ReadyNodes:       summary.ReadyNodes,
		TotalNamespaces:  summary.TotalNamespaces,
		TotalDeployments: summary.TotalDeployments,
		TotalPods:        summary.TotalPods,
		RunningPods:      summary.RunningPods,
	}
}

// ClusterHealthToResponse converts ClusterHealth model to ClusterHealthResponse
func ClusterHealthToResponse(health *k8sModels.ClusterHealth) k8sWire.ClusterHealthResponse {
	var nodes []k8sWire.NodeHealthResponse
	for _, node := range health.Nodes {
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

		nodes = append(nodes, k8sWire.NodeHealthResponse{
			Name:       node.Name,
			Status:     string(node.Status),
			Conditions: conditions,
			Resources: k8sWire.ResourceHealthResponse{
				CPU: k8sWire.ResourceHealthDetailResponse{
					Used:               node.Resources.CPU.Used,
					Available:          node.Resources.CPU.Available,
					Total:              node.Resources.CPU.Total,
					UtilizationPercent: node.Resources.CPU.UtilizationPercent,
					Status:             string(node.Resources.CPU.Status),
				},
				Memory: k8sWire.ResourceHealthDetailResponse{
					Used:               node.Resources.Memory.Used,
					Available:          node.Resources.Memory.Available,
					Total:              node.Resources.Memory.Total,
					UtilizationPercent: node.Resources.Memory.UtilizationPercent,
					Status:             string(node.Resources.Memory.Status),
				},
			},
			LastUpdated: node.LastUpdated,
		})
	}

	return k8sWire.ClusterHealthResponse{
		Context:     health.Context,
		Status:      string(health.Status),
		Nodes:       nodes,
		Summary:     ClusterSummaryToResponse(&health.Summary),
		LastUpdated: health.LastUpdated,
	}
}

// ClustersToResponse converts Cluster slice to ClusterResponse slice
func ClustersToResponse(clusters []k8sModels.Cluster) []k8sWire.ClusterResponse {
	var response []k8sWire.ClusterResponse
	for _, cluster := range clusters {
		response = append(response, k8sWire.ClusterResponse{
			Name:    cluster.Name,
			Context: cluster.Context,
			Server:  cluster.Server,
			Version: cluster.Version,
			Status:  string(cluster.Status),
		})
	}
	return response
}

// ClusterListToResponse converts Cluster slice to ClusterListResponse
func ClusterListToResponse(clusters []k8sModels.Cluster) k8sWire.ClusterListResponse {
	return k8sWire.ClusterListResponse{
		Clusters: ClustersToResponse(clusters),
		Total:    len(clusters),
	}
}
