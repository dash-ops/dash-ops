package kubernetes

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
)

// KubernetesAdapter implements Kubernetes service interfaces with data transformation
type KubernetesAdapter struct {
	client                 *KubernetesClient
	serviceContextResolver k8sPorts.ServiceContextResolver
}

// NewKubernetesAdapter creates a new Kubernetes adapter
func NewKubernetesAdapter(config *KubernetesConfig) (k8sPorts.KubernetesClientService, error) {
	client, err := NewKubernetesClient(config)
	if err != nil {
		return nil, err
	}

	return &KubernetesAdapter{
		client: client,
	}, nil
}

// GetClientset gets a Kubernetes clientset for a specific context
func (ka *KubernetesAdapter) GetClientset(context string) (k8sPorts.KubernetesClientset, error) {
	// For now, return the same client for all contexts
	// In a multi-context setup, this would manage different clients
	return ka, nil
}

// ValidateContext validates if a context is accessible
func (ka *KubernetesAdapter) ValidateContext(ctx string) error {
	return ka.client.TestConnection(context.Background())
}

// ListContexts lists all available contexts
func (ka *KubernetesAdapter) ListContexts() ([]string, error) {
	// For now, return the current context
	// In a multi-context setup, this would list all available contexts
	return []string{ka.client.GetContext()}, nil
}

// GetCurrentContext gets the current active context
func (ka *KubernetesAdapter) GetCurrentContext() (string, error) {
	return ka.client.GetContext(), nil
}

// SwitchContext switches to a different context
func (ka *KubernetesAdapter) SwitchContext(context string) error {
	// For now, this is a no-op
	// In a multi-context setup, this would switch the active context
	return nil
}

// Node operations
func (ka *KubernetesAdapter) GetNode(ctx context.Context, nodeName string) (*k8sModels.Node, error) {
	node, err := ka.client.GetNode(ctx, nodeName)
	if err != nil {
		return nil, err
	}
	return ka.convertNode(ctx, node), nil
}

func (ka *KubernetesAdapter) ListNodes(ctx context.Context) ([]k8sModels.Node, error) {
	nodeList, err := ka.client.ListNodes(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodes []k8sModels.Node
	for _, node := range nodeList.Items {
		nodes = append(nodes, *ka.convertNode(ctx, &node))
	}
	return nodes, nil
}

// Namespace operations
func (ka *KubernetesAdapter) GetNamespace(ctx context.Context, namespaceName string) (*k8sModels.Namespace, error) {
	namespace, err := ka.client.GetNamespace(ctx, namespaceName)
	if err != nil {
		return nil, err
	}
	return ka.convertNamespace(namespace), nil
}

func (ka *KubernetesAdapter) ListNamespaces(ctx context.Context) ([]k8sModels.Namespace, error) {
	namespaceList, err := ka.client.ListNamespaces(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []k8sModels.Namespace
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, *ka.convertNamespace(&namespace))
	}
	return namespaces, nil
}

func (ka *KubernetesAdapter) CreateNamespace(ctx context.Context, namespaceName string, labels map[string]string) (*k8sModels.Namespace, error) {
	namespace, err := ka.client.CreateNamespace(ctx, namespaceName, labels)
	if err != nil {
		return nil, err
	}
	return ka.convertNamespace(namespace), nil
}

func (ka *KubernetesAdapter) DeleteNamespace(ctx context.Context, namespaceName string) error {
	return ka.client.DeleteNamespace(ctx, namespaceName)
}

// Deployment operations
func (ka *KubernetesAdapter) GetDeployment(ctx context.Context, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	deployment, err := ka.client.GetDeployment(ctx, namespace, deploymentName)
	if err != nil {
		return nil, err
	}
	return ka.convertDeployment(deployment), nil
}

func (ka *KubernetesAdapter) ListDeployments(ctx context.Context, namespace string) ([]k8sModels.Deployment, error) {
	deploymentList, err := ka.client.ListDeployments(ctx, namespace, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deployments []k8sModels.Deployment
	for _, deployment := range deploymentList.Items {
		deployments = append(deployments, *ka.convertDeployment(&deployment))
	}
	return deployments, nil
}

func (ka *KubernetesAdapter) ScaleDeployment(ctx context.Context, namespace, deploymentName string, replicas int32) error {
	return ka.client.ScaleDeployment(ctx, namespace, deploymentName, replicas)
}

func (ka *KubernetesAdapter) RestartDeployment(ctx context.Context, namespace, deploymentName string) error {
	return ka.client.RestartDeployment(ctx, namespace, deploymentName)
}

// Pod operations
func (ka *KubernetesAdapter) GetPod(ctx context.Context, namespace, podName string) (*k8sModels.Pod, error) {
	pod, err := ka.client.GetPod(ctx, namespace, podName)
	if err != nil {
		return nil, err
	}
	return ka.convertPod(pod), nil
}

func (ka *KubernetesAdapter) ListPods(ctx context.Context, namespace string) ([]k8sModels.Pod, error) {
	podList, err := ka.client.ListPods(ctx, namespace, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var pods []k8sModels.Pod
	for _, pod := range podList.Items {
		pods = append(pods, *ka.convertPod(&pod))
	}
	return pods, nil
}

func (ka *KubernetesAdapter) DeletePod(ctx context.Context, namespace, podName string) error {
	return ka.client.DeletePod(ctx, namespace, podName)
}

func (ka *KubernetesAdapter) GetPodLogs(ctx context.Context, namespace, podName, containerName string, options *k8sPorts.LogOptions) (io.ReadCloser, error) {
	k8sOptions := &corev1.PodLogOptions{
		Container:    containerName,
		Follow:       options.Follow,
		Previous:     options.Previous,
		SinceSeconds: options.SinceSeconds,
		Timestamps:   options.Timestamps,
		TailLines:    options.TailLines,
	}

	// Only set SinceTime if it's not nil
	if options.SinceTime != nil {
		k8sOptions.SinceTime = &metav1.Time{Time: *options.SinceTime}
	}

	return ka.client.GetPodLogs(ctx, namespace, podName, k8sOptions)
}

// Data transformation methods
func (ka *KubernetesAdapter) convertDeployment(deployment *appsv1.Deployment) *k8sModels.Deployment {
	// Get replica status
	readyReplicas := deployment.Status.ReadyReplicas
	availableReplicas := deployment.Status.AvailableReplicas

	// Calculate pod info
	podInfo := k8sModels.PodInfo{
		Running: int(readyReplicas),
		Pending: int(*deployment.Spec.Replicas - availableReplicas),
		Failed:  0, // Would need to calculate from pod status
		Total:   int(*deployment.Spec.Replicas),
	}

	// Convert conditions
	var conditions []k8sModels.DeploymentCondition
	for _, condition := range deployment.Status.Conditions {
		conditions = append(conditions, k8sModels.DeploymentCondition{
			Type:           string(condition.Type),
			Status:         string(condition.Status),
			Reason:         condition.Reason,
			Message:        condition.Message,
			LastUpdateTime: condition.LastUpdateTime.Time,
		})
	}

	// Calculate age in the format expected by frontend
	age := time.Since(deployment.CreationTimestamp.Time)
	ageStr := age.String()

	return &k8sModels.Deployment{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
		PodInfo:   podInfo,
		Replicas: k8sModels.DeploymentReplicas{
			Desired:   *deployment.Spec.Replicas,
			Current:   availableReplicas,
			Ready:     readyReplicas,
			Available: availableReplicas,
		},
		Age:        ageStr,
		CreatedAt:  deployment.CreationTimestamp.Time,
		Conditions: conditions,
	}
}

func (ka *KubernetesAdapter) convertPod(pod *corev1.Pod) *k8sModels.Pod {
	// Get pod status
	status := k8sModels.PodStatus(pod.Status.Phase)

	// Get containers
	var containers []k8sModels.Container
	for _, cs := range pod.Status.ContainerStatuses {
		containers = append(containers, k8sModels.Container{
			Name:         cs.Name,
			Image:        "", // Would need to get from spec
			Ready:        cs.Ready,
			RestartCount: cs.RestartCount,
			State:        ka.convertContainerState(&cs.State),
		})
	}

	// Calculate ready status
	ready := "0/0"
	if len(pod.Status.ContainerStatuses) > 0 {
		readyCount := 0
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Ready {
				readyCount++
			}
		}
		ready = fmt.Sprintf("%d/%d", readyCount, len(pod.Status.ContainerStatuses))
	}

	// Calculate age in the format expected by frontend
	age := time.Since(pod.CreationTimestamp.Time)
	ageStr := age.String()

	return &k8sModels.Pod{
		Name:       pod.Name,
		Namespace:  pod.Namespace,
		Status:     status,
		Phase:      string(pod.Status.Phase),
		Node:       pod.Spec.NodeName,
		Age:        ageStr,
		Restarts:   int32(0), // Would need to calculate from container statuses
		Ready:      ready,
		IP:         pod.Status.PodIP,
		Containers: containers,
		CreatedAt:  pod.CreationTimestamp.Time,
	}
}

func (ka *KubernetesAdapter) convertNode(ctx context.Context, node *corev1.Node) *k8sModels.Node {
	// Get node status - determine from conditions (same logic as old repositories.go)
	status := k8sModels.NodeStatusNotReady
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
			status = k8sModels.NodeStatusReady
			break
		}
	}

	// Get node addresses
	var internalIP, externalIP string
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			internalIP = addr.Address
		} else if addr.Type == corev1.NodeExternalIP {
			externalIP = addr.Address
		}
	}

	// Get node roles from labels
	var roles []string
	for labelKey := range node.Labels {
		if labelKey == "node-role.kubernetes.io/control-plane" || labelKey == "node-role.kubernetes.io/master" {
			roles = append(roles, "control-plane")
		} else if labelKey == "node-role.kubernetes.io/worker" {
			roles = append(roles, "worker")
		}
	}
	if len(roles) == 0 {
		roles = []string{"worker"} // Default role
	}

	// Get capacity and allocatable from the node
	resources := k8sModels.NodeResources{
		Capacity: k8sModels.ResourceList{
			CPU:    node.Status.Capacity.Cpu().String(),
			Memory: node.Status.Capacity.Memory().String(),
			Pods:   node.Status.Capacity.Pods().String(),
		},
		Allocatable: k8sModels.ResourceList{
			CPU:    node.Status.Allocatable.Cpu().String(),
			Memory: node.Status.Allocatable.Memory().String(),
			Pods:   node.Status.Allocatable.Pods().String(),
		},
		// Used resources would require metrics-server to be installed
		// For now, we can at least count the pods
		Used: k8sModels.ResourceList{
			CPU:    "0",
			Memory: "0",
			Pods:   "0",
		},
	}

	// Try to get pod count and calculate resource usage on this node
	pods, err := ka.client.ListPods(ctx, "", metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name),
	})
	if err == nil {
		resources.Used.Pods = fmt.Sprintf("%d", len(pods.Items))

		// Calculate estimated CPU and Memory usage from pod requests
		var totalCPURequest, totalMemoryRequest int64
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				if container.Resources.Requests != nil {
					// CPU request in millicores
					if cpuRequest := container.Resources.Requests.Cpu(); cpuRequest != nil {
						totalCPURequest += cpuRequest.MilliValue()
					}
					// Memory request in bytes
					if memoryRequest := container.Resources.Requests.Memory(); memoryRequest != nil {
						totalMemoryRequest += memoryRequest.Value()
					}
				}
			}
		}

		// Convert to string format
		if totalCPURequest > 0 {
			resources.Used.CPU = fmt.Sprintf("%dm", totalCPURequest)
		}
		if totalMemoryRequest > 0 {
			// Convert bytes to Mi
			memoryMi := totalMemoryRequest / (1024 * 1024)
			resources.Used.Memory = fmt.Sprintf("%dMi", memoryMi)
		}
	}

	// Convert node conditions
	var conditions []k8sModels.NodeCondition
	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, k8sModels.NodeCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastTransitionTime: condition.LastTransitionTime.Time,
		})
	}

	// Calculate age in the format expected by frontend (e.g., "124h51m22.669937s")
	age := time.Since(node.CreationTimestamp.Time)
	ageStr := age.String() // This gives us the format like "124h51m22.669937s"

	return &k8sModels.Node{
		Name:       node.Name,
		Status:     status,
		Roles:      roles,
		Age:        ageStr,
		Version:    node.Status.NodeInfo.KubeletVersion,
		InternalIP: internalIP,
		ExternalIP: externalIP,
		Conditions: conditions,
		Resources:  resources,
		CreatedAt:  node.CreationTimestamp.Time,
	}
}

func (ka *KubernetesAdapter) convertNamespace(namespace *corev1.Namespace) *k8sModels.Namespace {
	// Extract labels
	labels := make(map[string]string)
	for key, value := range namespace.Labels {
		labels[key] = value
	}

	return &k8sModels.Namespace{
		Name:      namespace.Name,
		Status:    k8sModels.NamespaceStatus(namespace.Status.Phase),
		Labels:    labels,
		CreatedAt: namespace.CreationTimestamp.Time,
	}
}

func (ka *KubernetesAdapter) convertContainerState(state *corev1.ContainerState) k8sModels.ContainerState {
	if state.Running != nil {
		return k8sModels.ContainerState{
			Running: &k8sModels.ContainerStateRunning{
				StartedAt: state.Running.StartedAt.Time,
			},
		}
	}
	if state.Waiting != nil {
		return k8sModels.ContainerState{
			Waiting: &k8sModels.ContainerStateWaiting{
				Reason:  state.Waiting.Reason,
				Message: state.Waiting.Message,
			},
		}
	}
	if state.Terminated != nil {
		return k8sModels.ContainerState{
			Terminated: &k8sModels.ContainerStateTerminated{
				ExitCode:   state.Terminated.ExitCode,
				Reason:     state.Terminated.Reason,
				Message:    state.Terminated.Message,
				StartedAt:  state.Terminated.StartedAt.Time,
				FinishedAt: state.Terminated.FinishedAt.Time,
			},
		}
	}
	return k8sModels.ContainerState{}
}

// SetServiceContextResolver updates the service context resolver
func (ka *KubernetesAdapter) SetServiceContextResolver(resolver k8sPorts.ServiceContextResolver) {
	ka.serviceContextResolver = resolver
}

// Repository wrapper functions

// NewClusterRepository creates a cluster repository wrapper
func NewClusterRepository(adapter *KubernetesAdapter) k8sPorts.ClusterRepository {
	return &clusterRepositoryWrapper{adapter: adapter}
}

// NewNodeRepository creates a node repository wrapper
func NewNodeRepository(adapter *KubernetesAdapter) k8sPorts.NodeRepository {
	return &nodeRepositoryWrapper{adapter: adapter}
}

// NewNamespaceRepository creates a namespace repository wrapper
func NewNamespaceRepository(adapter *KubernetesAdapter) k8sPorts.NamespaceRepository {
	return &namespaceRepositoryWrapper{adapter: adapter}
}

// NewDeploymentRepository creates a deployment repository wrapper
func NewDeploymentRepository(adapter *KubernetesAdapter) k8sPorts.DeploymentRepository {
	return &deploymentRepositoryWrapper{adapter: adapter}
}

// NewPodRepository creates a pod repository wrapper
func NewPodRepository(adapter *KubernetesAdapter) k8sPorts.PodRepository {
	return &podRepositoryWrapper{adapter: adapter}
}

// Repository wrappers

type clusterRepositoryWrapper struct {
	adapter *KubernetesAdapter
}

func (w *clusterRepositoryWrapper) GetCluster(ctx context.Context, context string) (*k8sModels.Cluster, error) {
	// Validate context first
	if err := w.adapter.ValidateContext(context); err != nil {
		return nil, err
	}

	return &k8sModels.Cluster{
		Name:    context,
		Context: context,
		Version: "unknown", // Would need to get from server version
		Status:  k8sModels.ClusterStatusReady,
	}, nil
}

func (w *clusterRepositoryWrapper) ListClusters(ctx context.Context) ([]k8sModels.Cluster, error) {
	contexts, err := w.adapter.ListContexts()
	if err != nil {
		return nil, err
	}

	var clusters []k8sModels.Cluster
	for _, contextName := range contexts {
		clusters = append(clusters, k8sModels.Cluster{
			Name:    contextName,
			Context: contextName,
			Version: "unknown",
			Status:  k8sModels.ClusterStatusReady,
		})
	}

	return clusters, nil
}

func (w *clusterRepositoryWrapper) ValidateCluster(ctx context.Context, context string) error {
	return w.adapter.ValidateContext(context)
}

func (w *clusterRepositoryWrapper) GetClusterInfo(ctx context.Context, context string) (*k8sModels.ClusterInfo, error) {
	cluster, err := w.GetCluster(ctx, context)
	if err != nil {
		return nil, err
	}

	// Get nodes and namespaces
	nodes, err := w.adapter.ListNodes(ctx)
	if err != nil {
		nodes = []k8sModels.Node{} // Continue with empty if error
	}

	namespaces, err := w.adapter.ListNamespaces(ctx)
	if err != nil {
		namespaces = []k8sModels.Namespace{} // Continue with empty if error
	}

	return &k8sModels.ClusterInfo{
		Cluster:     *cluster,
		Nodes:       nodes,
		Namespaces:  namespaces,
		Summary:     k8sModels.ClusterSummary{},
		LastUpdated: time.Now(),
	}, nil
}

type nodeRepositoryWrapper struct {
	adapter *KubernetesAdapter
}

func (w *nodeRepositoryWrapper) GetNode(ctx context.Context, context, nodeName string) (*k8sModels.Node, error) {
	return w.adapter.GetNode(ctx, nodeName)
}

func (w *nodeRepositoryWrapper) ListNodes(ctx context.Context, context string) ([]k8sModels.Node, error) {
	return w.adapter.ListNodes(ctx)
}

func (w *nodeRepositoryWrapper) GetNodeMetrics(ctx context.Context, context, nodeName string) (*k8sModels.NodeResources, error) {
	node, err := w.GetNode(ctx, context, nodeName)
	if err != nil {
		return nil, err
	}

	return &node.Resources, nil
}

type namespaceRepositoryWrapper struct {
	adapter *KubernetesAdapter
}

func (w *namespaceRepositoryWrapper) GetNamespace(ctx context.Context, context, namespaceName string) (*k8sModels.Namespace, error) {
	return w.adapter.GetNamespace(ctx, namespaceName)
}

func (w *namespaceRepositoryWrapper) ListNamespaces(ctx context.Context, context string) ([]k8sModels.Namespace, error) {
	return w.adapter.ListNamespaces(ctx)
}

func (w *namespaceRepositoryWrapper) CreateNamespace(ctx context.Context, context, namespaceName string, labels map[string]string) (*k8sModels.Namespace, error) {
	return w.adapter.CreateNamespace(ctx, namespaceName, labels)
}

func (w *namespaceRepositoryWrapper) DeleteNamespace(ctx context.Context, context, namespaceName string) error {
	return w.adapter.DeleteNamespace(ctx, namespaceName)
}

type deploymentRepositoryWrapper struct {
	adapter *KubernetesAdapter
}

func (w *deploymentRepositoryWrapper) SetServiceContextResolver(resolver k8sPorts.ServiceContextResolver) {
	w.adapter.SetServiceContextResolver(resolver)
}

func (w *deploymentRepositoryWrapper) GetDeployment(ctx context.Context, context, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	deployment, err := w.adapter.GetDeployment(ctx, namespace, deploymentName)
	if err != nil {
		return nil, err
	}

	// Add service context if resolver is available
	if w.adapter.serviceContextResolver != nil {
		serviceContext, err := w.adapter.serviceContextResolver.ResolveDeploymentService(
			deploymentName,
			namespace,
			context,
		)
		if err == nil && serviceContext != nil {
			deployment.ServiceContext = serviceContext
		}
	}

	return deployment, nil
}

func (w *deploymentRepositoryWrapper) ListDeployments(ctx context.Context, context string, filter *k8sModels.DeploymentFilter) (*k8sModels.DeploymentList, error) {
	namespace := ""
	if filter != nil {
		namespace = filter.Namespace
	}

	deployments, err := w.adapter.ListDeployments(ctx, namespace)
	if err != nil {
		return nil, err
	}

	// Apply service context resolution
	for i := range deployments {
		if w.adapter.serviceContextResolver != nil {
			serviceContext, err := w.adapter.serviceContextResolver.ResolveDeploymentService(
				deployments[i].Name,
				deployments[i].Namespace,
				context,
			)
			if err == nil && serviceContext != nil {
				deployments[i].ServiceContext = serviceContext
			}
		}
	}

	return &k8sModels.DeploymentList{
		Deployments: deployments,
		Total:       len(deployments),
		Namespace:   namespace,
		Filter:      filter,
	}, nil
}

func (w *deploymentRepositoryWrapper) ScaleDeployment(ctx context.Context, context, namespace, deploymentName string, replicas int32) error {
	return w.adapter.ScaleDeployment(ctx, namespace, deploymentName, replicas)
}

func (w *deploymentRepositoryWrapper) RestartDeployment(ctx context.Context, context, namespace, deploymentName string) error {
	return w.adapter.RestartDeployment(ctx, namespace, deploymentName)
}

func (w *deploymentRepositoryWrapper) GetDeploymentStatus(ctx context.Context, context, namespace, deploymentName string) (*k8sPorts.DeploymentStatus, error) {
	deployment, err := w.GetDeployment(ctx, context, namespace, deploymentName)
	if err != nil {
		return nil, err
	}

	healthStatus := "healthy"
	if deployment.Replicas.Ready < deployment.Replicas.Desired {
		healthStatus = "degraded"
	}
	if deployment.Replicas.Ready == 0 {
		healthStatus = "unhealthy"
	}

	return &k8sPorts.DeploymentStatus{
		Name:         deployment.Name,
		Namespace:    deployment.Namespace,
		Replicas:     deployment.Replicas,
		Conditions:   deployment.Conditions,
		HealthStatus: healthStatus,
		LastUpdated:  time.Now(),
	}, nil
}

type podRepositoryWrapper struct {
	adapter *KubernetesAdapter
}

func (w *podRepositoryWrapper) GetPod(ctx context.Context, context, namespace, podName string) (*k8sModels.Pod, error) {
	return w.adapter.GetPod(ctx, namespace, podName)
}

func (w *podRepositoryWrapper) ListPods(ctx context.Context, context string, filter *k8sModels.PodFilter) (*k8sModels.PodList, error) {
	namespace := ""
	if filter != nil {
		namespace = filter.Namespace
	}

	pods, err := w.adapter.ListPods(ctx, namespace)
	if err != nil {
		return nil, err
	}

	// Apply additional filters
	var filteredPods []k8sModels.Pod
	for _, pod := range pods {
		if filter != nil {
			if filter.PodName != "" && pod.Name != filter.PodName {
				continue
			}
			if filter.Status != "" && string(pod.Status) != filter.Status {
				continue
			}
			if filter.Node != "" && pod.Node != filter.Node {
				continue
			}
		}
		filteredPods = append(filteredPods, pod)
	}

	return &k8sModels.PodList{
		Pods:      filteredPods,
		Total:     len(filteredPods),
		Namespace: namespace,
		Filter:    filter,
	}, nil
}

func (w *podRepositoryWrapper) DeletePod(ctx context.Context, context, namespace, podName string) error {
	return w.adapter.DeletePod(ctx, namespace, podName)
}

func (w *podRepositoryWrapper) GetPodLogs(ctx context.Context, context string, filter *k8sModels.LogFilter) ([]k8sModels.ContainerLog, error) {
	options := &k8sPorts.LogOptions{
		Follow:       false,
		Previous:     false,
		SinceSeconds: nil,
		SinceTime:    nil,
		Timestamps:   false,
		TailLines:    &filter.TailLines,
	}

	logs, err := w.adapter.GetPodLogs(ctx, filter.Namespace, filter.PodName, filter.ContainerName, options)
	if err != nil {
		return nil, err
	}
	defer logs.Close()

	// Parse logs into ContainerLog entries
	var containerLogs []k8sModels.ContainerLog

	// Read all logs from the stream
	buf := new(strings.Builder)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	logsContent := buf.String()
	lines := strings.Split(logsContent, "\n")
	for _, line := range lines {
		if line != "" {
			containerLogs = append(containerLogs, k8sModels.ContainerLog{
				Timestamp: time.Now(), // In real implementation, parse timestamp from log
				Message:   line,
			})
		}
	}

	return containerLogs, nil
}

func (w *podRepositoryWrapper) GetPodMetrics(ctx context.Context, context, namespace, podName string) (*k8sModels.PodMetrics, error) {
	// TODO: Implement metrics API client
	return &k8sModels.PodMetrics{
		PodName:     podName,
		Namespace:   namespace,
		Containers:  []k8sModels.ContainerMetrics{},
		LastUpdated: time.Now(),
	}, nil
}
