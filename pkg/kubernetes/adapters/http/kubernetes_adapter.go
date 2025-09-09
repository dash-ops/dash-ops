package http

import (
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes/wire"
)

// KubernetesAdapter handles transformation between models and wire formats
type KubernetesAdapter struct{}

// NewKubernetesAdapter creates a new Kubernetes adapter
func NewKubernetesAdapter() *KubernetesAdapter {
	return &KubernetesAdapter{}
}

// ClusterInfoToResponse converts ClusterInfo model to ClusterInfoResponse
func (ka *KubernetesAdapter) ClusterInfoToResponse(clusterInfo *k8sModels.ClusterInfo) k8sWire.ClusterInfoResponse {
	return k8sWire.ClusterInfoResponse{
		Cluster: k8sWire.ClusterResponse{
			Name:    clusterInfo.Cluster.Name,
			Context: clusterInfo.Cluster.Context,
			Server:  clusterInfo.Cluster.Server,
			Version: clusterInfo.Cluster.Version,
			Status:  string(clusterInfo.Cluster.Status),
		},
		Nodes:       ka.NodesToResponse(clusterInfo.Nodes),
		Namespaces:  ka.NamespacesToResponse(clusterInfo.Namespaces),
		Summary:     ka.ClusterSummaryToResponse(&clusterInfo.Summary),
		LastUpdated: clusterInfo.LastUpdated,
	}
}

// NodesToResponse converts Node models to NodeResponse slice
func (ka *KubernetesAdapter) NodesToResponse(nodes []k8sModels.Node) []k8sWire.NodeResponse {
	var response []k8sWire.NodeResponse
	for _, node := range nodes {
		response = append(response, *ka.NodeToResponse(&node))
	}
	return response
}

// NodeToResponse converts Node model to NodeResponse
func (ka *KubernetesAdapter) NodeToResponse(node *k8sModels.Node) *k8sWire.NodeResponse {
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

// NamespacesToResponse converts Namespace models to NamespaceResponse slice
func (ka *KubernetesAdapter) NamespacesToResponse(namespaces []k8sModels.Namespace) []k8sWire.NamespaceResponse {
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

// ClusterSummaryToResponse converts ClusterSummary model to ClusterSummaryResponse
func (ka *KubernetesAdapter) ClusterSummaryToResponse(summary *k8sModels.ClusterSummary) k8sWire.ClusterSummaryResponse {
	return k8sWire.ClusterSummaryResponse{
		TotalNodes:       summary.TotalNodes,
		ReadyNodes:       summary.ReadyNodes,
		TotalNamespaces:  summary.TotalNamespaces,
		TotalDeployments: summary.TotalDeployments,
		TotalPods:        summary.TotalPods,
		RunningPods:      summary.RunningPods,
	}
}

// DeploymentListToResponse converts DeploymentList model to DeploymentListResponse
func (ka *KubernetesAdapter) DeploymentListToResponse(deploymentList *k8sModels.DeploymentList) k8sWire.DeploymentListResponse {
	var deployments []k8sWire.DeploymentResponse
	for _, deployment := range deploymentList.Deployments {
		deployments = append(deployments, *ka.DeploymentToResponse(&deployment))
	}

	return k8sWire.DeploymentListResponse{
		Deployments: deployments,
		Total:       deploymentList.Total,
		Namespace:   deploymentList.Namespace,
		Filter:      deploymentList.Filter,
	}
}

// DeploymentToResponse converts Deployment model to DeploymentResponse
func (ka *KubernetesAdapter) DeploymentToResponse(deployment *k8sModels.Deployment) *k8sWire.DeploymentResponse {
	var conditions []k8sWire.DeploymentConditionResponse
	for _, condition := range deployment.Conditions {
		conditions = append(conditions, k8sWire.DeploymentConditionResponse{
			Type:           condition.Type,
			Status:         condition.Status,
			Reason:         condition.Reason,
			Message:        condition.Message,
			LastUpdateTime: condition.LastUpdateTime,
		})
	}

	var serviceContext *k8sWire.ServiceContextResponse
	if deployment.ServiceContext != nil {
		serviceContext = &k8sWire.ServiceContextResponse{
			ServiceName: deployment.ServiceContext.ServiceName,
			ServiceTier: deployment.ServiceContext.ServiceTier,
			Environment: deployment.ServiceContext.Environment,
			Context:     deployment.ServiceContext.Context,
			Team:        deployment.ServiceContext.Team,
			Description: deployment.ServiceContext.Description,
			Found:       deployment.ServiceContext.Found,
		}
	}

	return &k8sWire.DeploymentResponse{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
		PodInfo: k8sWire.PodInfoResponse{
			Running: deployment.PodInfo.Running,
			Pending: deployment.PodInfo.Pending,
			Failed:  deployment.PodInfo.Failed,
			Total:   deployment.PodInfo.Total,
		},
		Replicas: k8sWire.DeploymentReplicasResponse{
			Desired:   deployment.Replicas.Desired,
			Current:   deployment.Replicas.Current,
			Ready:     deployment.Replicas.Ready,
			Available: deployment.Replicas.Available,
		},
		Age:                 deployment.Age,
		CreatedAt:           deployment.CreatedAt,
		Conditions:          conditions,
		ServiceContext:      serviceContext,
		AvailabilityPercent: deployment.GetAvailabilityPercentage(),
	}
}

// PodListToResponse converts PodList model to PodListResponse
func (ka *KubernetesAdapter) PodListToResponse(podList *k8sModels.PodList) k8sWire.PodListResponse {
	var pods []k8sWire.PodResponse
	for _, pod := range podList.Pods {
		pods = append(pods, *ka.PodToResponse(&pod))
	}

	return k8sWire.PodListResponse{
		Pods:      pods,
		Total:     podList.Total,
		Namespace: podList.Namespace,
		Filter:    podList.Filter,
	}
}

// PodToResponse converts Pod model to PodResponse
func (ka *KubernetesAdapter) PodToResponse(pod *k8sModels.Pod) *k8sWire.PodResponse {
	var containers []k8sWire.ContainerResponse
	for _, container := range pod.Containers {
		containers = append(containers, k8sWire.ContainerResponse{
			Name:         container.Name,
			Image:        container.Image,
			Ready:        container.Ready,
			RestartCount: container.RestartCount,
			State:        ka.containerStateToResponse(container.State),
		})
	}

	var conditions []k8sWire.PodConditionResponse
	for _, condition := range pod.Conditions {
		conditions = append(conditions, k8sWire.PodConditionResponse{
			Type:               condition.Type,
			Status:             condition.Status,
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime,
		})
	}

	return &k8sWire.PodResponse{
		Name:       pod.Name,
		Namespace:  pod.Namespace,
		Status:     string(pod.Status),
		Phase:      pod.Phase,
		Node:       pod.Node,
		Age:        pod.Age,
		Restarts:   pod.Restarts,
		Ready:      pod.Ready,
		IP:         pod.IP,
		Containers: containers,
		Conditions: conditions,
		CreatedAt:  pod.CreatedAt,
	}
}

// containerStateToResponse converts ContainerState model to ContainerStateResponse
func (ka *KubernetesAdapter) containerStateToResponse(state k8sModels.ContainerState) k8sWire.ContainerStateResponse {
	response := k8sWire.ContainerStateResponse{}

	if state.Running != nil {
		response.Running = &k8sWire.ContainerStateRunningResponse{
			StartedAt: state.Running.StartedAt,
		}
	}

	if state.Waiting != nil {
		response.Waiting = &k8sWire.ContainerStateWaitingResponse{
			Reason:  state.Waiting.Reason,
			Message: state.Waiting.Message,
		}
	}

	if state.Terminated != nil {
		response.Terminated = &k8sWire.ContainerStateTerminatedResponse{
			ExitCode:   state.Terminated.ExitCode,
			Reason:     state.Terminated.Reason,
			Message:    state.Terminated.Message,
			StartedAt:  state.Terminated.StartedAt,
			FinishedAt: state.Terminated.FinishedAt,
		}
	}

	return response
}

// ClusterHealthToResponse converts ClusterHealth model to ClusterHealthResponse
func (ka *KubernetesAdapter) ClusterHealthToResponse(health *k8sModels.ClusterHealth) k8sWire.ClusterHealthResponse {
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
		Summary:     ka.ClusterSummaryToResponse(&health.Summary),
		LastUpdated: health.LastUpdated,
	}
}

// NamespaceToResponse converts Namespace model to NamespaceResponse
func (ka *KubernetesAdapter) NamespaceToResponse(namespace *k8sModels.Namespace) *k8sWire.NamespaceResponse {
	return &k8sWire.NamespaceResponse{
		Name:      namespace.Name,
		Status:    string(namespace.Status),
		Labels:    namespace.Labels,
		Age:       namespace.Age,
		CreatedAt: namespace.CreatedAt,
	}
}

// DeploymentsToResponse converts Deployment slice to DeploymentResponse slice
func (ka *KubernetesAdapter) DeploymentsToResponse(deployments []k8sModels.Deployment) []k8sWire.DeploymentResponse {
	var response []k8sWire.DeploymentResponse
	for _, deployment := range deployments {
		response = append(response, *ka.DeploymentToResponse(&deployment))
	}
	return response
}

// PodsToResponse converts Pod slice to PodResponse slice
func (ka *KubernetesAdapter) PodsToResponse(pods []k8sModels.Pod) []k8sWire.PodResponse {
	var response []k8sWire.PodResponse
	for _, pod := range pods {
		response = append(response, *ka.PodToResponse(&pod))
	}
	return response
}

// ClustersToResponse converts Cluster slice to ClusterResponse slice
func (ka *KubernetesAdapter) ClustersToResponse(clusters []k8sModels.Cluster) []k8sWire.ClusterResponse {
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
func (ka *KubernetesAdapter) ClusterListToResponse(clusters []k8sModels.Cluster) k8sWire.ClusterListResponse {
	return k8sWire.ClusterListResponse{
		Clusters: ka.ClustersToResponse(clusters),
		Total:    len(clusters),
	}
}
