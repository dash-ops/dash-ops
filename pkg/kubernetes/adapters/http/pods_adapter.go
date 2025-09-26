package http

import (
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes/wire"
)

// PodListToResponse converts PodList model to PodListResponse
func PodListToResponse(podList *k8sModels.PodList) k8sWire.PodListResponse {
	var pods []k8sWire.PodResponse
	for _, pod := range podList.Pods {
		pods = append(pods, *PodToResponse(&pod))
	}

	return k8sWire.PodListResponse{
		Pods:      pods,
		Total:     podList.Total,
		Namespace: podList.Namespace,
		Filter:    podList.Filter,
	}
}

// PodToResponse converts Pod model to PodResponse
func PodToResponse(pod *k8sModels.Pod) *k8sWire.PodResponse {
	var containers []k8sWire.ContainerResponse
	for _, container := range pod.Containers {
		containers = append(containers, k8sWire.ContainerResponse{
			Name:         container.Name,
			Image:        container.Image,
			Ready:        container.Ready,
			RestartCount: container.RestartCount,
			State:        containerStateToResponse(container.State),
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
func containerStateToResponse(state k8sModels.ContainerState) k8sWire.ContainerStateResponse {
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

// PodsToResponse converts Pod slice to PodResponse slice
func PodsToResponse(pods []k8sModels.Pod) []k8sWire.PodResponse {
	var response []k8sWire.PodResponse
	for _, pod := range pods {
		response = append(response, *PodToResponse(&pod))
	}
	return response
}
