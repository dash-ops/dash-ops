package controllers

import (
	"context"
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/repositories"
)

// PodsController handles pods business logic orchestration
type PodsController struct {
	repository *repositories.PodsRepository
}

// NewPodsController creates a new pods controller
func NewPodsController(repository *repositories.PodsRepository) *PodsController {
	return &PodsController{
		repository: repository,
	}
}

// GetPod gets a specific pod with business logic validation
func (c *PodsController) GetPod(ctx context.Context, context, namespace, podName string) (*k8sModels.Pod, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if podName == "" {
		return nil, fmt.Errorf("pod name is required")
	}

	pod, err := c.repository.GetPod(ctx, context, namespace, podName)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	return pod, nil
}

// ListPods lists pods with optional filtering and business logic
func (c *PodsController) ListPods(ctx context.Context, context string, filter *k8sModels.PodFilter) (*k8sModels.PodList, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	podList, err := c.repository.ListPods(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	return podList, nil
}

// DeletePod deletes a pod with business logic validation
func (c *PodsController) DeletePod(ctx context.Context, context, namespace, podName string) error {
	if context == "" {
		return fmt.Errorf("context is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if podName == "" {
		return fmt.Errorf("pod name is required")
	}

	// Business logic: prevent deletion of pods in system namespaces
	if isSystemNamespace(namespace) {
		return fmt.Errorf("cannot delete pods in system namespace: %s", namespace)
	}

	// Verify pod exists before deletion
	_, err := c.repository.GetPod(ctx, context, namespace, podName)
	if err != nil {
		return fmt.Errorf("pod not found: %w", err)
	}

	err = c.repository.DeletePod(ctx, context, namespace, podName)
	if err != nil {
		return fmt.Errorf("failed to delete pod: %w", err)
	}

	return nil
}

// GetPodLogs gets logs for a pod/container with business logic processing
func (c *PodsController) GetPodLogs(ctx context.Context, context string, filter *k8sModels.LogFilter) ([]k8sModels.ContainerLog, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if filter == nil {
		return nil, fmt.Errorf("log filter is required")
	}
	if filter.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if filter.PodName == "" {
		return nil, fmt.Errorf("pod name is required")
	}

	// Business logic: validate log filter parameters
	if filter.TailLines < 0 {
		return nil, fmt.Errorf("tail lines must be non-negative")
	}

	// Set default tail lines if not specified
	if filter.TailLines == 0 {
		filter.TailLines = 100
	}

	// Business logic: limit maximum tail lines
	const maxTailLines = 10000
	if filter.TailLines > maxTailLines {
		filter.TailLines = maxTailLines
	}

	// Verify pod exists before getting logs
	_, err := c.repository.GetPod(ctx, context, filter.Namespace, filter.PodName)
	if err != nil {
		return nil, fmt.Errorf("pod not found: %w", err)
	}

	logs, err := c.repository.GetPodLogs(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod logs: %w", err)
	}

	return logs, nil
}

// GetPodsSummary provides a summary of pods in a namespace or cluster
func (c *PodsController) GetPodsSummary(ctx context.Context, context string, namespace string) (*PodsSummary, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	// Create filter for namespace if specified
	filter := &k8sModels.PodFilter{}
	if namespace != "" {
		filter.Namespace = namespace
	}

	podList, err := c.repository.ListPods(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods for summary: %w", err)
	}

	summary := &PodsSummary{
		Total:     len(podList.Pods),
		Running:   0,
		Pending:   0,
		Failed:    0,
		Succeeded: 0,
		Unknown:   0,
	}

	for _, pod := range podList.Pods {
		switch pod.Status {
		case k8sModels.PodStatusRunning:
			summary.Running++
		case k8sModels.PodStatusPending:
			summary.Pending++
		case k8sModels.PodStatusFailed:
			summary.Failed++
		case k8sModels.PodStatusSucceeded:
			summary.Succeeded++
		default:
			summary.Unknown++
		}
	}

	return summary, nil
}

// GetPodsByNode gets all pods running on a specific node
func (c *PodsController) GetPodsByNode(ctx context.Context, context, nodeName string) (*k8sModels.PodList, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	filter := &k8sModels.PodFilter{
		Node: nodeName,
	}

	podList, err := c.repository.ListPods(ctx, context, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods by node: %w", err)
	}

	return podList, nil
}

// isSystemNamespace checks if a namespace is a system namespace
func isSystemNamespace(name string) bool {
	systemNamespaces := []string{
		"kube-system",
		"kube-public",
		"kube-node-lease",
		"default",
		"kubernetes-dashboard",
		"kube-flannel",
		"kube-ingress",
		"kube-monitoring",
		"kube-storage",
		"kube-logging",
		"istio-system",
		"knative-serving",
		"knative-eventing",
	}

	for _, sysNs := range systemNamespaces {
		if name == sysNs {
			return true
		}
	}

	return false
}

// PodsSummary represents a summary of pods
type PodsSummary struct {
	Total     int `json:"total"`
	Running   int `json:"running"`
	Pending   int `json:"pending"`
	Failed    int `json:"failed"`
	Succeeded int `json:"succeeded"`
	Unknown   int `json:"unknown"`
}
