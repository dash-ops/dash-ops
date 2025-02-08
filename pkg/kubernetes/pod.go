package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod Struct representing an k8s pod
type Pod struct {
	Name            string          `json:"name"`
	Namespace       string          `json:"namespace"`
	ConditionStatus PodStatus       `json:"condition_status"`
	RestartCount    int32           `json:"restart_count"`
	NodeName        string          `json:"node_name"`
	Requests        v1.ResourceList `json:"requests"`
	Limits          v1.ResourceList `json:"limits"`
}

// PodStatus representing an k8s pod status
type PodStatus struct {
	Status          string              `json:"status"`
	PodPhase        v1.PodPhase         `json:"pod_phase"`
	ContainerStates []v1.ContainerState `json:"container_states"`
}

type podFilter struct {
	Name      string
	Namespace string
}

// ContainerLog pod list container logs
type ContainerLog struct {
	Name string `json:"name"`
	Log  string `json:"log"`
}

func (kc client) GetPods(filter podFilter) ([]Pod, error) {
	var pods []Pod

	if filter.Namespace == "" {
		filter.Namespace = apiv1.NamespaceAll
	}

	podsList, err := kc.clientSet.
		CoreV1().
		Pods(filter.Namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get pods: %s", err)
	}

	for _, p := range podsList.Items {
		reqs, limits, _ := podRequestsAndLimits(&p)

		pods = append(pods, Pod{
			Name:            p.GetName(),
			Namespace:       p.GetNamespace(),
			ConditionStatus: getPodConditionStatus(p),
			RestartCount:    getPodRestartCount(p),
			NodeName:        p.Spec.NodeName,
			Requests:        reqs,
			Limits:          limits,
		})
	}

	return pods, nil
}

func (kc client) GetPodLogs(filter podFilter) ([]ContainerLog, error) {
	var logs []ContainerLog

	pod, err := kc.clientSet.CoreV1().Pods(filter.Namespace).Get(context.TODO(), filter.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get pod: %s", err)
	}

	for _, container := range pod.Spec.Containers {
		podLogOpts := v1.PodLogOptions{Container: container.Name}

		req := kc.clientSet.CoreV1().Pods(filter.Namespace).GetLogs(filter.Name, &podLogOpts)
		podLogs, err := req.Stream(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("Error in opening stream: %s", err)
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return nil, fmt.Errorf("Error in copy information from podLogs to buf: %s", err)
		}

		logs = append(logs, ContainerLog{
			Name: container.Name,
			Log:  buf.String(),
		})
	}

	return logs, nil
}

func getPodConditionStatus(pod v1.Pod) PodStatus {
	var states []v1.ContainerState
	for _, containerStatus := range pod.Status.ContainerStatuses {
		states = append(states, containerStatus.State)
	}

	return PodStatus{
		Status:          string(getPodStatusPhase(pod)),
		PodPhase:        pod.Status.Phase,
		ContainerStates: states,
	}
}

func getPodRestartCount(pod v1.Pod) int32 {
	var restartCount int32 = 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restartCount += containerStatus.RestartCount
	}
	return restartCount
}

func podRequestsAndLimits(pod *v1.Pod) (reqs, limits v1.ResourceList, err error) {
	reqs, limits = v1.ResourceList{}, v1.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}

	// Add overhead for running a pod to the sum of requests and to non-zero limits:
	if pod.Spec.Overhead != nil {
		addResourceList(reqs, pod.Spec.Overhead)

		for name, quantity := range pod.Spec.Overhead {
			if value, ok := limits[name]; ok && !value.IsZero() {
				value.Add(quantity)
				limits[name] = value
			}
		}
	}
	return
}

func addResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

func maxResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = quantity.DeepCopy()
			}
		}
	}
}

func getPodStatusPhase(pod v1.Pod) v1.PodPhase {
	if pod.Status.Phase == v1.PodFailed {
		return v1.PodFailed
	}

	if pod.Status.Phase == v1.PodSucceeded {
		return v1.PodSucceeded
	}

	ready := false
	initialized := false
	for _, c := range pod.Status.Conditions {
		if c.Type == v1.PodReady {
			ready = c.Status == v1.ConditionTrue
		}
		if c.Type == v1.PodInitialized {
			initialized = c.Status == v1.ConditionTrue
		}
	}

	if initialized && ready && pod.Status.Phase == v1.PodRunning {
		return v1.PodRunning
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		return v1.PodUnknown
	} else if pod.DeletionTimestamp != nil {
		return "Terminating"
	}

	return v1.PodPending
}
