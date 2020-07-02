package kubernetes

import (
	"bytes"
	"fmt"
	"io"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod Struct representing an k8s pod
type Pod struct {
	Name         string             `json:"name"`
	Namespace    string             `json:"namespace"`
	Status       v1.ConditionStatus `json:"status"`
	RestartCount int32              `json:"restart_count"`
	NodeName     string             `json:"node_name"`
	Requests     v1.ResourceList    `json:"requests"`
	Limits       v1.ResourceList    `json:"limits"`
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

func (kc k8sClient) GetPods(filter podFilter) ([]Pod, error) {
	var pods []Pod

	if filter.Namespace == "" {
		filter.Namespace = apiv1.NamespaceAll
	}

	podsList, err := kc.clientSet.
		CoreV1().
		Pods(filter.Namespace).
		List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get pods: %s", err)
	}

	for _, p := range podsList.Items {
		reqs, limits, _ := podRequestsAndLimits(&p)

		pods = append(pods, Pod{
			Name:         p.GetName(),
			Namespace:    p.GetNamespace(),
			Status:       getPodConditionStatus(p, v1.PodReady),
			RestartCount: getPodRestartCount(p),
			NodeName:     p.Spec.NodeName,
			Requests:     reqs,
			Limits:       limits,
		})
	}

	return pods, nil
}

func (kc k8sClient) GetPodLogs(filter podFilter) ([]ContainerLog, error) {
	var logs []ContainerLog

	pod, err := kc.clientSet.CoreV1().Pods(filter.Namespace).Get(filter.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get pod: %s", err)
	}

	for _, container := range pod.Spec.Containers {
		podLogOpts := v1.PodLogOptions{Container: container.Name}

		req := kc.clientSet.CoreV1().Pods(filter.Namespace).GetLogs(filter.Name, &podLogOpts)
		podLogs, err := req.Stream()
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

func getPodConditionStatus(pod v1.Pod, conditionType v1.PodConditionType) v1.ConditionStatus {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return v1.ConditionUnknown
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
