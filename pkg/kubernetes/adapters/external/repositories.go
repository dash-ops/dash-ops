package external

import (
	"context"
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
)

// ClusterRepositoryImpl implements ClusterRepository interface
type ClusterRepositoryImpl struct {
	clients map[string]*KubernetesClient
}

// NewClusterRepository creates a new cluster repository
func NewClusterRepository(configs []KubernetesConfig) (k8sPorts.ClusterRepository, error) {
	clients := make(map[string]*KubernetesClient)

	for _, config := range configs {
		client, err := NewKubernetesClient(&config)
		if err != nil {
			return nil, fmt.Errorf("failed to create client for context %s: %w", config.Context, err)
		}

		contextName := config.Context
		if contextName == "" {
			contextName = "default"
		}
		clients[contextName] = client
	}

	return &ClusterRepositoryImpl{
		clients: clients,
	}, nil
}

// GetCluster gets cluster information
func (cr *ClusterRepositoryImpl) GetCluster(ctx context.Context, context string) (*k8sModels.Cluster, error) {
	client, ok := cr.clients[context]
	if !ok {
		return nil, fmt.Errorf("no client found for context: %s", context)
	}

	version, err := client.GetVersion(ctx)
	if err != nil {
		return nil, err
	}

	return &k8sModels.Cluster{
		Name:    context,
		Context: context,
		Version: version,
		Status:  k8sModels.ClusterStatusReady,
	}, nil
}

// ListClusters lists all configured clusters
func (cr *ClusterRepositoryImpl) ListClusters(ctx context.Context) ([]k8sModels.Cluster, error) {
	var clusters []k8sModels.Cluster

	for context, client := range cr.clients {
		version, _ := client.GetVersion(ctx)
		clusters = append(clusters, k8sModels.Cluster{
			Name:    context,
			Context: context,
			Version: version,
			Status:  k8sModels.ClusterStatusReady,
		})
	}

	return clusters, nil
}

// ValidateCluster validates cluster connectivity
func (cr *ClusterRepositoryImpl) ValidateCluster(ctx context.Context, context string) error {
	client, ok := cr.clients[context]
	if !ok {
		return fmt.Errorf("no client found for context: %s", context)
	}

	_, err := client.GetVersion(ctx)
	return err
}

// GetClusterInfo gets comprehensive cluster information
func (cr *ClusterRepositoryImpl) GetClusterInfo(ctx context.Context, context string) (*k8sModels.ClusterInfo, error) {
	cluster, err := cr.GetCluster(ctx, context)
	if err != nil {
		return nil, err
	}

	// For now, return basic info - could be expanded to include nodes, namespaces, etc.
	return &k8sModels.ClusterInfo{
		Cluster:     *cluster,
		Nodes:       []k8sModels.Node{},
		Namespaces:  []k8sModels.Namespace{},
		Summary:     k8sModels.ClusterSummary{},
		LastUpdated: time.Now(),
	}, nil
}

// GetClient gets the kubernetes client for a context
func (cr *ClusterRepositoryImpl) GetClient(context string) (*KubernetesClient, error) {
	client, ok := cr.clients[context]
	if !ok {
		return nil, fmt.Errorf("no client found for context: %s", context)
	}
	return client, nil
}

// NodeRepositoryImpl implements NodeRepository interface
type NodeRepositoryImpl struct {
	clusterRepo *ClusterRepositoryImpl
}

// NewNodeRepository creates a new node repository
func NewNodeRepository(clusterRepo k8sPorts.ClusterRepository) k8sPorts.NodeRepository {
	return &NodeRepositoryImpl{
		clusterRepo: clusterRepo.(*ClusterRepositoryImpl),
	}
}

// ListNodes lists all nodes in a cluster
func (nr *NodeRepositoryImpl) ListNodes(ctx context.Context, context string) ([]k8sModels.Node, error) {
	client, err := nr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sNodes, err := client.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	var nodes []k8sModels.Node
	for _, k8sNode := range k8sNodes {
		nodes = append(nodes, nr.convertNode(ctx, client, k8sNode))
	}

	return nodes, nil
}

// GetNode gets a specific node
func (nr *NodeRepositoryImpl) GetNode(ctx context.Context, context, name string) (*k8sModels.Node, error) {
	client, err := nr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sNode, err := client.GetNode(ctx, name)
	if err != nil {
		return nil, err
	}

	node := nr.convertNode(ctx, client, *k8sNode)
	return &node, nil
}

// GetNodeMetrics gets node resource metrics
func (nr *NodeRepositoryImpl) GetNodeMetrics(ctx context.Context, context, nodeName string) (*k8sModels.NodeResources, error) {
	client, err := nr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sNode, err := client.GetNode(ctx, nodeName)
	if err != nil {
		return nil, err
	}

	// Get capacity and allocatable from the node
	resources := &k8sModels.NodeResources{
		Capacity: k8sModels.ResourceList{
			CPU:    k8sNode.Status.Capacity.Cpu().String(),
			Memory: k8sNode.Status.Capacity.Memory().String(),
			Pods:   k8sNode.Status.Capacity.Pods().String(),
		},
		Allocatable: k8sModels.ResourceList{
			CPU:    k8sNode.Status.Allocatable.Cpu().String(),
			Memory: k8sNode.Status.Allocatable.Memory().String(),
			Pods:   k8sNode.Status.Allocatable.Pods().String(),
		},
		// Used resources would require metrics-server to be installed
		// For now, we can at least count the pods
		Used: k8sModels.ResourceList{
			CPU:    "0",
			Memory: "0",
			Pods:   "0",
		},
	}

	// Try to get pod count on this node
	pods, err := client.ListPods(ctx, "", fmt.Sprintf("spec.nodeName=%s", nodeName))
	if err == nil {
		resources.Used.Pods = fmt.Sprintf("%d", len(pods))
	}

	return resources, nil
}

// convertNode converts k8s node to model
func (nr *NodeRepositoryImpl) convertNode(ctx context.Context, client *KubernetesClient, k8sNode corev1.Node) k8sModels.Node {
	var roles []string
	for label := range k8sNode.Labels {
		if strings.HasPrefix(label, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(label, "node-role.kubernetes.io/")
			roles = append(roles, role)
		}
	}

	var conditions []k8sModels.NodeCondition
	for _, cond := range k8sNode.Status.Conditions {
		conditions = append(conditions, k8sModels.NodeCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastTransitionTime: cond.LastTransitionTime.Time,
		})
	}

	status := k8sModels.NodeStatusNotReady
	for _, cond := range k8sNode.Status.Conditions {
		if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
			status = k8sModels.NodeStatusReady
			break
		}
	}

	var internalIP, externalIP string
	for _, addr := range k8sNode.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			internalIP = addr.Address
		}
		if addr.Type == corev1.NodeExternalIP {
			externalIP = addr.Address
		}
	}

	// Get node resources from capacity and allocatable
	resources := k8sModels.NodeResources{
		Capacity: k8sModels.ResourceList{
			CPU:    k8sNode.Status.Capacity.Cpu().String(),
			Memory: k8sNode.Status.Capacity.Memory().String(),
			Pods:   k8sNode.Status.Capacity.Pods().String(),
		},
		Allocatable: k8sModels.ResourceList{
			CPU:    k8sNode.Status.Allocatable.Cpu().String(),
			Memory: k8sNode.Status.Allocatable.Memory().String(),
			Pods:   k8sNode.Status.Allocatable.Pods().String(),
		},
		// Calculate used resources from pods on this node
		Used: nr.calculateUsedResources(ctx, client, k8sNode),
	}

	return k8sModels.Node{
		Name:       k8sNode.Name,
		Status:     status,
		Roles:      roles,
		Age:        time.Since(k8sNode.CreationTimestamp.Time).String(),
		Version:    k8sNode.Status.NodeInfo.KubeletVersion,
		InternalIP: internalIP,
		ExternalIP: externalIP,
		Conditions: conditions,
		Resources:  resources,
		CreatedAt:  k8sNode.CreationTimestamp.Time,
	}
}

// NamespaceRepositoryImpl implements NamespaceRepository interface
type NamespaceRepositoryImpl struct {
	clusterRepo *ClusterRepositoryImpl
}

// NewNamespaceRepository creates a new namespace repository
func NewNamespaceRepository(clusterRepo k8sPorts.ClusterRepository) k8sPorts.NamespaceRepository {
	return &NamespaceRepositoryImpl{
		clusterRepo: clusterRepo.(*ClusterRepositoryImpl),
	}
}

// ListNamespaces lists all namespaces
func (nr *NamespaceRepositoryImpl) ListNamespaces(ctx context.Context, context string) ([]k8sModels.Namespace, error) {
	client, err := nr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sNamespaces, err := client.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	var namespaces []k8sModels.Namespace
	for _, k8sNs := range k8sNamespaces {
		namespaces = append(namespaces, nr.convertNamespace(k8sNs))
	}

	return namespaces, nil
}

// GetNamespace gets a specific namespace
func (nr *NamespaceRepositoryImpl) GetNamespace(ctx context.Context, context, name string) (*k8sModels.Namespace, error) {
	client, err := nr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sNs, err := client.GetNamespace(ctx, name)
	if err != nil {
		return nil, err
	}

	ns := nr.convertNamespace(*k8sNs)
	return &ns, nil
}

// CreateNamespace creates a new namespace
func (nr *NamespaceRepositoryImpl) CreateNamespace(ctx context.Context, context, name string, labels map[string]string) (*k8sModels.Namespace, error) {
	client, err := nr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sNs, err := client.CreateNamespace(ctx, name, labels)
	if err != nil {
		return nil, err
	}

	ns := nr.convertNamespace(*k8sNs)
	return &ns, nil
}

// DeleteNamespace deletes a namespace
func (nr *NamespaceRepositoryImpl) DeleteNamespace(ctx context.Context, context, name string) error {
	client, err := nr.clusterRepo.GetClient(context)
	if err != nil {
		return err
	}

	return client.DeleteNamespace(ctx, name)
}

// convertNamespace converts k8s namespace to model
func (nr *NamespaceRepositoryImpl) convertNamespace(k8sNs corev1.Namespace) k8sModels.Namespace {
	status := k8sModels.NamespaceStatusActive
	if k8sNs.Status.Phase == corev1.NamespaceTerminating {
		status = k8sModels.NamespaceStatusTerminating
	}

	return k8sModels.Namespace{
		Name:      k8sNs.Name,
		Status:    status,
		Labels:    k8sNs.Labels,
		Age:       time.Since(k8sNs.CreationTimestamp.Time).String(),
		CreatedAt: k8sNs.CreationTimestamp.Time,
	}
}

// DeploymentRepositoryImpl implements DeploymentRepository interface
type DeploymentRepositoryImpl struct {
	clusterRepo            *ClusterRepositoryImpl
	serviceContextResolver k8sPorts.ServiceContextResolver
}

// NewDeploymentRepository creates a new deployment repository
func NewDeploymentRepository(clusterRepo k8sPorts.ClusterRepository, serviceContextResolver k8sPorts.ServiceContextResolver) k8sPorts.DeploymentRepository {
	return &DeploymentRepositoryImpl{
		clusterRepo:            clusterRepo.(*ClusterRepositoryImpl),
		serviceContextResolver: serviceContextResolver,
	}
}

// ListDeployments lists deployments with optional filtering
func (dr *DeploymentRepositoryImpl) ListDeployments(ctx context.Context, context string, filter *k8sModels.DeploymentFilter) (*k8sModels.DeploymentList, error) {
	client, err := dr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	namespace := ""
	if filter != nil {
		namespace = filter.Namespace
	}

	k8sDeployments, err := client.ListDeployments(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var deployments []k8sModels.Deployment
	for _, k8sDep := range k8sDeployments {
		deployment := dr.convertDeployment(k8sDep, dr.serviceContextResolver, context)

		// Apply additional filters if needed
		if filter != nil {
			if filter.ServiceName != "" {
				// Skip if doesn't match service name filter
				// This would need service context resolution
			}
			if filter.Status != "" {
				// Skip if doesn't match status filter
			}
		}

		deployments = append(deployments, deployment)
	}

	return &k8sModels.DeploymentList{
		Deployments: deployments,
		Total:       len(deployments),
		Namespace:   namespace,
		Filter:      filter,
	}, nil
}

// GetDeployment gets a specific deployment
func (dr *DeploymentRepositoryImpl) GetDeployment(ctx context.Context, context, namespace, name string) (*k8sModels.Deployment, error) {
	client, err := dr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sDep, err := client.GetDeployment(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	deployment := dr.convertDeployment(*k8sDep, dr.serviceContextResolver, context)
	return &deployment, nil
}

// ScaleDeployment scales a deployment
func (dr *DeploymentRepositoryImpl) ScaleDeployment(ctx context.Context, context, namespace, name string, replicas int32) error {
	client, err := dr.clusterRepo.GetClient(context)
	if err != nil {
		return err
	}

	return client.ScaleDeployment(ctx, namespace, name, replicas)
}

// RestartDeployment restarts a deployment
func (dr *DeploymentRepositoryImpl) RestartDeployment(ctx context.Context, context, namespace, name string) error {
	client, err := dr.clusterRepo.GetClient(context)
	if err != nil {
		return err
	}

	return client.RestartDeployment(ctx, namespace, name)
}

// GetDeploymentStatus gets deployment status and health
func (dr *DeploymentRepositoryImpl) GetDeploymentStatus(ctx context.Context, context, namespace, name string) (*k8sPorts.DeploymentStatus, error) {
	deployment, err := dr.GetDeployment(ctx, context, namespace, name)
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

// convertDeployment converts k8s deployment to model
func (dr *DeploymentRepositoryImpl) convertDeployment(k8sDep appsv1.Deployment, serviceContextResolver k8sPorts.ServiceContextResolver, context string) k8sModels.Deployment {
	var conditions []k8sModels.DeploymentCondition
	for _, cond := range k8sDep.Status.Conditions {
		conditions = append(conditions, k8sModels.DeploymentCondition{
			Type:           string(cond.Type),
			Status:         string(cond.Status),
			Reason:         cond.Reason,
			Message:        cond.Message,
			LastUpdateTime: cond.LastUpdateTime.Time,
		})
	}

	desired := int32(1)
	if k8sDep.Spec.Replicas != nil {
		desired = *k8sDep.Spec.Replicas
	}

	// Calculate PodInfo based on deployment status
	// Running = ReadyReplicas (pods that are ready and running)
	// Pending = Replicas - ReadyReplicas (pods that are created but not ready yet)
	// Failed = Replicas - AvailableReplicas (pods that failed to become available)
	// Total = Replicas (total number of pods)
	podInfo := k8sModels.PodInfo{
		Running: int(k8sDep.Status.ReadyReplicas),
		Pending: int(k8sDep.Status.Replicas - k8sDep.Status.ReadyReplicas),
		Failed:  int(k8sDep.Status.Replicas - k8sDep.Status.AvailableReplicas),
		Total:   int(k8sDep.Status.Replicas),
	}

	// Resolve service context if resolver is available
	var serviceContext *k8sModels.ServiceContext
	if serviceContextResolver != nil {
		resolvedContext, err := serviceContextResolver.ResolveDeploymentService(
			k8sDep.Name,
			k8sDep.Namespace,
			context,
		)
		if err == nil && resolvedContext != nil {
			serviceContext = resolvedContext
		}
		// Ignore errors - service context is optional
	}

	return k8sModels.Deployment{
		Name:      k8sDep.Name,
		Namespace: k8sDep.Namespace,
		PodInfo:   podInfo,
		Replicas: k8sModels.DeploymentReplicas{
			Desired:   desired,
			Current:   k8sDep.Status.Replicas,
			Ready:     k8sDep.Status.ReadyReplicas,
			Available: k8sDep.Status.AvailableReplicas,
		},
		Age:            time.Since(k8sDep.CreationTimestamp.Time).String(),
		CreatedAt:      k8sDep.CreationTimestamp.Time,
		Conditions:     conditions,
		ServiceContext: serviceContext,
	}
}

// PodRepositoryImpl implements PodRepository interface
type PodRepositoryImpl struct {
	clusterRepo *ClusterRepositoryImpl
}

// NewPodRepository creates a new pod repository
func NewPodRepository(clusterRepo k8sPorts.ClusterRepository) k8sPorts.PodRepository {
	return &PodRepositoryImpl{
		clusterRepo: clusterRepo.(*ClusterRepositoryImpl),
	}
}

// ListPods lists pods with optional filtering
func (pr *PodRepositoryImpl) ListPods(ctx context.Context, context string, filter *k8sModels.PodFilter) (*k8sModels.PodList, error) {
	client, err := pr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	namespace := ""
	labelSelector := ""
	if filter != nil {
		namespace = filter.Namespace
		labelSelector = filter.LabelSelector
	}

	k8sPods, err := client.ListPods(ctx, namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	var pods []k8sModels.Pod
	for _, k8sPod := range k8sPods {
		pod := pr.convertPod(k8sPod)

		// Apply additional filters if needed
		if filter != nil {
			if filter.PodName != "" && !strings.Contains(pod.Name, filter.PodName) {
				continue
			}
			if filter.Status != "" && string(pod.Status) != filter.Status {
				continue
			}
			if filter.Node != "" && pod.Node != filter.Node {
				continue
			}
		}

		pods = append(pods, pod)
	}

	return &k8sModels.PodList{
		Pods:      pods,
		Total:     len(pods),
		Namespace: namespace,
		Filter:    filter,
	}, nil
}

// GetPod gets a specific pod
func (pr *PodRepositoryImpl) GetPod(ctx context.Context, context, namespace, name string) (*k8sModels.Pod, error) {
	client, err := pr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	k8sPod, err := client.GetPod(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	pod := pr.convertPod(*k8sPod)
	return &pod, nil
}

// DeletePod deletes a pod
func (pr *PodRepositoryImpl) DeletePod(ctx context.Context, context, namespace, name string) error {
	client, err := pr.clusterRepo.GetClient(context)
	if err != nil {
		return err
	}

	return client.DeletePod(ctx, namespace, name)
}

// GetPodLogs gets pod logs
func (pr *PodRepositoryImpl) GetPodLogs(ctx context.Context, context string, filter *k8sModels.LogFilter) ([]k8sModels.ContainerLog, error) {
	client, err := pr.clusterRepo.GetClient(context)
	if err != nil {
		return nil, err
	}

	tailLines := int64(100)
	if filter.TailLines > 0 {
		tailLines = filter.TailLines
	}

	logs, err := client.GetPodLogs(ctx, filter.Namespace, filter.PodName, filter.ContainerName, tailLines)
	if err != nil {
		return nil, err
	}

	// Parse logs into ContainerLog entries
	var containerLogs []k8sModels.ContainerLog
	lines := strings.Split(logs, "\n")
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

// GetPodMetrics gets pod resource metrics
func (pr *PodRepositoryImpl) GetPodMetrics(ctx context.Context, context, namespace, podName string) (*k8sModels.PodMetrics, error) {
	// TODO: Implement metrics API client
	// For now, return empty metrics
	return &k8sModels.PodMetrics{
		PodName:     podName,
		Namespace:   namespace,
		Containers:  []k8sModels.ContainerMetrics{},
		LastUpdated: time.Now(),
	}, nil
}

// convertPod converts k8s pod to model
func (pr *PodRepositoryImpl) convertPod(k8sPod corev1.Pod) k8sModels.Pod {
	var containers []k8sModels.Container
	var totalRestarts int32

	for _, container := range k8sPod.Status.ContainerStatuses {
		totalRestarts += container.RestartCount

		var state k8sModels.ContainerState
		if container.State.Running != nil {
			state.Running = &k8sModels.ContainerStateRunning{
				StartedAt: container.State.Running.StartedAt.Time,
			}
		} else if container.State.Waiting != nil {
			state.Waiting = &k8sModels.ContainerStateWaiting{
				Reason:  container.State.Waiting.Reason,
				Message: container.State.Waiting.Message,
			}
		} else if container.State.Terminated != nil {
			state.Terminated = &k8sModels.ContainerStateTerminated{
				ExitCode:   container.State.Terminated.ExitCode,
				Reason:     container.State.Terminated.Reason,
				Message:    container.State.Terminated.Message,
				StartedAt:  container.State.Terminated.StartedAt.Time,
				FinishedAt: container.State.Terminated.FinishedAt.Time,
			}
		}

		containers = append(containers, k8sModels.Container{
			Name:         container.Name,
			Image:        container.Image,
			Ready:        container.Ready,
			RestartCount: container.RestartCount,
			State:        state,
		})
	}

	var conditions []k8sModels.PodCondition
	for _, cond := range k8sPod.Status.Conditions {
		conditions = append(conditions, k8sModels.PodCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastTransitionTime: cond.LastTransitionTime.Time,
		})
	}

	// Determine pod status
	var status k8sModels.PodStatus
	switch k8sPod.Status.Phase {
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

	readyCount := 0
	for _, container := range k8sPod.Status.ContainerStatuses {
		if container.Ready {
			readyCount++
		}
	}

	return k8sModels.Pod{
		Name:       k8sPod.Name,
		Namespace:  k8sPod.Namespace,
		Status:     status,
		Phase:      string(k8sPod.Status.Phase),
		Node:       k8sPod.Spec.NodeName,
		Age:        time.Since(k8sPod.CreationTimestamp.Time).String(),
		Restarts:   totalRestarts,
		Ready:      fmt.Sprintf("%d/%d", readyCount, len(k8sPod.Spec.Containers)),
		IP:         k8sPod.Status.PodIP,
		Containers: containers,
		Conditions: conditions,
		CreatedAt:  k8sPod.CreationTimestamp.Time,
	}
}

// calculateUsedResources calculates used resources for a node based on pods
func (nr *NodeRepositoryImpl) calculateUsedResources(ctx context.Context, client *KubernetesClient, k8sNode corev1.Node) k8sModels.ResourceList {
	// Get pods on this node
	pods, err := nr.getPodsOnNode(ctx, client, k8sNode.Name)
	if err != nil {
		// If we can't get pods, return zero usage
		return k8sModels.ResourceList{
			CPU:    "0",
			Memory: "0",
			Pods:   "0",
		}
	}

	// Calculate total resource requests from all pods
	var totalCPU, totalMemory resource.Quantity
	podCount := len(pods)

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if cpu := container.Resources.Requests[corev1.ResourceCPU]; !cpu.IsZero() {
				totalCPU.Add(cpu)
			}
			if memory := container.Resources.Requests[corev1.ResourceMemory]; !memory.IsZero() {
				totalMemory.Add(memory)
			}
		}
		// Also check init containers
		for _, container := range pod.Spec.InitContainers {
			if cpu := container.Resources.Requests[corev1.ResourceCPU]; !cpu.IsZero() {
				totalCPU.Add(cpu)
			}
			if memory := container.Resources.Requests[corev1.ResourceMemory]; !memory.IsZero() {
				totalMemory.Add(memory)
			}
		}
	}

	return k8sModels.ResourceList{
		CPU:    totalCPU.String(),
		Memory: totalMemory.String(),
		Pods:   fmt.Sprintf("%d", podCount),
	}
}

// getPodsOnNode gets all pods running on a specific node
func (nr *NodeRepositoryImpl) getPodsOnNode(ctx context.Context, client *KubernetesClient, nodeName string) ([]corev1.Pod, error) {
	// List all pods and filter by node name
	pods, err := client.ListPods(ctx, "", "")
	if err != nil {
		return nil, err
	}

	var nodePods []corev1.Pod
	for _, pod := range pods {
		if pod.Spec.NodeName == nodeName {
			nodePods = append(nodePods, pod)
		}
	}

	return nodePods, nil
}
