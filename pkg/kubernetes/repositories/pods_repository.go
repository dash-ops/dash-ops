package repositories

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dash-ops/dash-ops/pkg/kubernetes/integrations/external/kubernetes"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

// PodsRepository handles pod-related data access
type PodsRepository struct {
	client *kubernetes.KubernetesClient
}

// NewPodsRepository creates a new pods repository
func NewPodsRepository(client *kubernetes.KubernetesClient) *PodsRepository {
	return &PodsRepository{
		client: client,
	}
}

// GetPod gets a specific pod
func (r *PodsRepository) GetPod(ctx context.Context, context, namespace, podName string) (*k8sModels.Pod, error) {
	if podName == "" {
		return nil, fmt.Errorf("pod name is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	pod, err := r.client.GetPod(ctx, namespace, podName)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s/%s: %w", namespace, podName, err)
	}

	return r.convertPod(pod), nil
}

// ListPods lists pods with optional filtering
func (r *PodsRepository) ListPods(ctx context.Context, context string, filter *k8sModels.PodFilter) (*k8sModels.PodList, error) {
	// Build list options based on filter
	listOptions := metav1.ListOptions{}
	if filter != nil {
		if filter.LabelSelector != "" {
			listOptions.LabelSelector = filter.LabelSelector
		}
		if filter.FieldSelector != "" {
			listOptions.FieldSelector = filter.FieldSelector
		}
	}

	var pods []k8sModels.Pod

	if filter != nil && filter.Namespace != "" {
		// List pods in specific namespace
		podList, err := r.client.ListPods(ctx, filter.Namespace, listOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list pods in namespace %s: %w", filter.Namespace, err)
		}

		for _, pod := range podList.Items {
			pods = append(pods, *r.convertPod(&pod))
		}
	} else {
		// For cross-namespace listing, we need to list all namespaces first
		// and then get pods from each namespace
		namespaces, err := r.client.ListNamespaces(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}

		for _, namespace := range namespaces.Items {
			podList, err := r.client.ListPods(ctx, namespace.Name, listOptions)
			if err != nil {
				// Log error but continue with other namespaces
				continue
			}

			for _, pod := range podList.Items {
				pods = append(pods, *r.convertPod(&pod))
			}
		}
	}

	// Apply additional filters
	if filter != nil {
		pods = r.filterPods(pods, filter)
	}

	return &k8sModels.PodList{
		Pods:  pods,
		Total: len(pods),
	}, nil
}

// DeletePod deletes a pod
func (r *PodsRepository) DeletePod(ctx context.Context, context, namespace, podName string) error {
	if podName == "" {
		return fmt.Errorf("pod name is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	err := r.client.DeletePod(ctx, namespace, podName)
	if err != nil {
		return fmt.Errorf("failed to delete pod %s/%s: %w", namespace, podName, err)
	}

	return nil
}

// GetPodLogs gets logs for a pod/container
func (r *PodsRepository) GetPodLogs(ctx context.Context, context string, filter *k8sModels.LogFilter) ([]k8sModels.ContainerLog, error) {
	if filter.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if filter.PodName == "" {
		return nil, fmt.Errorf("pod name is required")
	}

	// Get pod to determine available containers
	pod, err := r.client.GetPod(ctx, filter.Namespace, filter.PodName)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod for logs: %w", err)
	}

	var logs []k8sModels.ContainerLog

	// If no specific container is specified, get logs from all containers
	if filter.ContainerName == "" {
		for _, container := range pod.Spec.Containers {
			containerLogs, err := r.getContainerLogs(ctx, filter.Namespace, filter.PodName, container.Name, filter.TailLines)
			if err != nil {
				// Log error but continue with other containers
				fmt.Printf("Warning: failed to get logs for container %s: %v\n", container.Name, err)
				continue
			}
			logs = append(logs, containerLogs...)
		}
	} else {
		// Get logs from specific container
		containerLogs, err := r.getContainerLogs(ctx, filter.Namespace, filter.PodName, filter.ContainerName, filter.TailLines)
		if err != nil {
			return nil, fmt.Errorf("failed to get logs for container %s: %w", filter.ContainerName, err)
		}
		logs = append(logs, containerLogs...)
	}

	return logs, nil
}

// getContainerLogs gets logs for a specific container
func (r *PodsRepository) getContainerLogs(ctx context.Context, namespace, podName, containerName string, tailLines int64) ([]k8sModels.ContainerLog, error) {
	// Get logs stream
	logOptions := &corev1.PodLogOptions{
		Container:  containerName,
		TailLines:  &tailLines,
		Timestamps: true,
	}

	logStream, err := r.client.GetPodLogs(ctx, namespace, podName, logOptions)
	if err != nil {
		return nil, err
	}
	defer logStream.Close()

	// Parse logs
	return r.parseLogs(logStream, namespace, podName, containerName), nil
}

// parseLogs parses log stream into structured log entries
func (r *PodsRepository) parseLogs(stream io.ReadCloser, namespace, podName, containerName string) []k8sModels.ContainerLog {
	var logs []k8sModels.ContainerLog
	scanner := bufio.NewScanner(stream)

	// Regex to match Kubernetes log format with timestamps
	// Format: 2023-12-19T10:30:45.123456789Z stdout P log message
	timestampRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z)\s+(\w+)\s+([FP])\s+(.*)$`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		log := k8sModels.ContainerLog{
			ContainerName: containerName,
			PodName:       podName,
			Namespace:     namespace,
			Message:       line,
			Level:         "INFO",     // Default level
			Timestamp:     time.Now(), // Default to current time
		}

		// Try to parse structured log format
		if matches := timestampRegex.FindStringSubmatch(line); len(matches) == 5 {
			if timestamp, err := time.Parse(time.RFC3339Nano, matches[1]); err == nil {
				log.Timestamp = timestamp
			}
			log.Level = r.determineLogLevel(matches[4])
			log.Message = matches[4]
		} else {
			// Fallback: try to extract log level from message
			log.Level = r.determineLogLevel(line)
		}

		logs = append(logs, log)
	}

	return logs
}

// determineLogLevel determines log level from message content
func (r *PodsRepository) determineLogLevel(message string) string {
	message = strings.ToUpper(message)

	switch {
	case strings.Contains(message, "ERROR") || strings.Contains(message, "FATAL"):
		return "ERROR"
	case strings.Contains(message, "WARN"):
		return "WARN"
	case strings.Contains(message, "DEBUG"):
		return "DEBUG"
	case strings.Contains(message, "TRACE"):
		return "TRACE"
	default:
		return "INFO"
	}
}

// filterPods applies additional filters to the pod list
func (r *PodsRepository) filterPods(pods []k8sModels.Pod, filter *k8sModels.PodFilter) []k8sModels.Pod {
	var filtered []k8sModels.Pod

	for _, pod := range pods {
		// Apply pod name filter
		if filter.PodName != "" && pod.Name != filter.PodName {
			continue
		}

		// Apply container name filter (check if pod has the specified container)
		if filter.ContainerName != "" {
			found := false
			for _, container := range pod.Containers {
				if container.Name == filter.ContainerName {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply status filter
		if filter.Status != "" {
			if string(pod.Status) != filter.Status {
				continue
			}
		}

		// Apply node filter
		if filter.Node != "" && pod.Node != filter.Node {
			continue
		}

		// Apply limit if specified
		if filter.Limit > 0 && len(filtered) >= filter.Limit {
			break
		}

		filtered = append(filtered, pod)
	}

	return filtered
}

// convertPod converts a Kubernetes pod to our domain model
func (r *PodsRepository) convertPod(pod *corev1.Pod) *k8sModels.Pod {
	// Determine pod status and phase
	var status k8sModels.PodStatus
	phase := string(pod.Status.Phase)

	switch pod.Status.Phase {
	case corev1.PodRunning:
		status = k8sModels.PodStatusRunning
	case corev1.PodPending:
		status = k8sModels.PodStatusPending
	case corev1.PodSucceeded:
		status = k8sModels.PodStatusSucceeded
	case corev1.PodFailed:
		status = k8sModels.PodStatusFailed
	default:
		status = k8sModels.PodStatusUnknown
	}

	// Convert containers
	var containers []k8sModels.Container
	var totalRestarts int32
	var readyCount int

	for _, container := range pod.Spec.Containers {
		// Find container status
		var containerStatus *corev1.ContainerStatus
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Name == container.Name {
				containerStatus = &cs
				break
			}
		}

		containerModel := k8sModels.Container{
			Name:         container.Name,
			Image:        container.Image,
			Ready:        false,
			RestartCount: 0,
		}

		if containerStatus != nil {
			containerModel.Ready = containerStatus.Ready
			containerModel.RestartCount = containerStatus.RestartCount
			totalRestarts += containerStatus.RestartCount

			// Convert container state
			if containerStatus.State.Running != nil {
				containerModel.State = k8sModels.ContainerState{
					Running: &k8sModels.ContainerStateRunning{
						StartedAt: containerStatus.State.Running.StartedAt.Time,
					},
				}
			} else if containerStatus.State.Waiting != nil {
				containerModel.State = k8sModels.ContainerState{
					Waiting: &k8sModels.ContainerStateWaiting{
						Reason:  containerStatus.State.Waiting.Reason,
						Message: containerStatus.State.Waiting.Message,
					},
				}
			} else if containerStatus.State.Terminated != nil {
				containerModel.State = k8sModels.ContainerState{
					Terminated: &k8sModels.ContainerStateTerminated{
						ExitCode:   containerStatus.State.Terminated.ExitCode,
						Reason:     containerStatus.State.Terminated.Reason,
						Message:    containerStatus.State.Terminated.Message,
						StartedAt:  containerStatus.State.Terminated.StartedAt.Time,
						FinishedAt: containerStatus.State.Terminated.FinishedAt.Time,
					},
				}
			}

			// Convert resources
			if container.Resources.Requests != nil || container.Resources.Limits != nil {
				containerModel.Resources = k8sModels.ContainerResources{
					Requests: k8sModels.ResourceList{
						CPU:    container.Resources.Requests.Cpu().String(),
						Memory: container.Resources.Requests.Memory().String(),
					},
					Limits: k8sModels.ResourceList{
						CPU:    container.Resources.Limits.Cpu().String(),
						Memory: container.Resources.Limits.Memory().String(),
					},
				}
			}

			if containerStatus.Ready {
				readyCount++
			}
		}

		containers = append(containers, containerModel)
	}

	// Convert conditions
	var conditions []k8sModels.PodCondition
	for _, condition := range pod.Status.Conditions {
		conditions = append(conditions, k8sModels.PodCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime.Time,
		})
	}

	// Determine QoS class
	qosClass := string(pod.Status.QOSClass)

	// Prepare ready string
	readyStr := fmt.Sprintf("%d/%d", readyCount, len(containers))

	return &k8sModels.Pod{
		Name:       pod.Name,
		Namespace:  pod.Namespace,
		Status:     status,
		Phase:      phase,
		Node:       pod.Spec.NodeName,
		Age:        time.Since(pod.CreationTimestamp.Time).Round(time.Second).String(),
		Restarts:   totalRestarts,
		Ready:      readyStr,
		IP:         pod.Status.PodIP,
		Containers: containers,
		Conditions: conditions,
		CreatedAt:  pod.CreationTimestamp.Time,
		QoSClass:   qosClass,
	}
}
