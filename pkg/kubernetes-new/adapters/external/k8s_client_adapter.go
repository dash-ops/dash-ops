package external

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes-new/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes-new/ports"
)

// KubernetesClientAdapter implements KubernetesClientset interface
type KubernetesClientAdapter struct {
	clientset *kubernetes.Clientset
	context   string
}

// NewKubernetesClientAdapter creates a new Kubernetes client adapter
func NewKubernetesClientAdapter(kubeconfig, context string) (*KubernetesClientAdapter, error) {
	config, err := buildConfig(kubeconfig, context)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &KubernetesClientAdapter{
		clientset: clientset,
		context:   context,
	}, nil
}

// GetNode gets a specific node
func (kca *KubernetesClientAdapter) GetNode(ctx context.Context, nodeName string) (*k8sModels.Node, error) {
	node, err := kca.clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}

	return kca.convertNode(node), nil
}

// ListNodes lists all nodes
func (kca *KubernetesClientAdapter) ListNodes(ctx context.Context) ([]k8sModels.Node, error) {
	nodeList, err := kca.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	var nodes []k8sModels.Node
	for _, node := range nodeList.Items {
		nodes = append(nodes, *kca.convertNode(&node))
	}

	return nodes, nil
}

// GetNamespace gets a specific namespace
func (kca *KubernetesClientAdapter) GetNamespace(ctx context.Context, namespaceName string) (*k8sModels.Namespace, error) {
	namespace, err := kca.clientset.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace %s: %w", namespaceName, err)
	}

	return kca.convertNamespace(namespace), nil
}

// ListNamespaces lists all namespaces
func (kca *KubernetesClientAdapter) ListNamespaces(ctx context.Context) ([]k8sModels.Namespace, error) {
	namespaceList, err := kca.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var namespaces []k8sModels.Namespace
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, *kca.convertNamespace(&namespace))
	}

	return namespaces, nil
}

// CreateNamespace creates a new namespace
func (kca *KubernetesClientAdapter) CreateNamespace(ctx context.Context, namespaceName string, labels map[string]string) (*k8sModels.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespaceName,
			Labels: labels,
		},
	}

	createdNamespace, err := kca.clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace %s: %w", namespaceName, err)
	}

	return kca.convertNamespace(createdNamespace), nil
}

// DeleteNamespace deletes a namespace
func (kca *KubernetesClientAdapter) DeleteNamespace(ctx context.Context, namespaceName string) error {
	err := kca.clientset.CoreV1().Namespaces().Delete(ctx, namespaceName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete namespace %s: %w", namespaceName, err)
	}
	return nil
}

// GetDeployment gets a specific deployment
func (kca *KubernetesClientAdapter) GetDeployment(ctx context.Context, namespace, deploymentName string) (*k8sModels.Deployment, error) {
	deployment, err := kca.clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s/%s: %w", namespace, deploymentName, err)
	}

	return kca.convertDeployment(deployment), nil
}

// ListDeployments lists deployments in a namespace
func (kca *KubernetesClientAdapter) ListDeployments(ctx context.Context, namespace string) ([]k8sModels.Deployment, error) {
	deploymentList, err := kca.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}

	var deployments []k8sModels.Deployment
	for _, deployment := range deploymentList.Items {
		deployments = append(deployments, *kca.convertDeployment(&deployment))
	}

	return deployments, nil
}

// ScaleDeployment scales a deployment to specified replicas
func (kca *KubernetesClientAdapter) ScaleDeployment(ctx context.Context, namespace, deploymentName string, replicas int32) error {
	scale, err := kca.clientset.AppsV1().Deployments(namespace).GetScale(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment scale: %w", err)
	}

	scale.Spec.Replicas = replicas
	_, err = kca.clientset.AppsV1().Deployments(namespace).UpdateScale(ctx, deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to scale deployment: %w", err)
	}

	return nil
}

// RestartDeployment restarts a deployment by updating its annotation
func (kca *KubernetesClientAdapter) RestartDeployment(ctx context.Context, namespace, deploymentName string) error {
	deployment, err := kca.clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// Add restart annotation to trigger rolling update
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = kca.clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to restart deployment: %w", err)
	}

	return nil
}

// GetPod gets a specific pod
func (kca *KubernetesClientAdapter) GetPod(ctx context.Context, namespace, podName string) (*k8sModels.Pod, error) {
	pod, err := kca.clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s/%s: %w", namespace, podName, err)
	}

	return kca.convertPod(pod), nil
}

// ListPods lists pods in a namespace
func (kca *KubernetesClientAdapter) ListPods(ctx context.Context, namespace string) ([]k8sModels.Pod, error) {
	podList, err := kca.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	var pods []k8sModels.Pod
	for _, pod := range podList.Items {
		pods = append(pods, *kca.convertPod(&pod))
	}

	return pods, nil
}

// DeletePod deletes a pod
func (kca *KubernetesClientAdapter) DeletePod(ctx context.Context, namespace, podName string) error {
	err := kca.clientset.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod %s/%s: %w", namespace, podName, err)
	}
	return nil
}

// GetPodLogs gets logs for a pod/container
func (kca *KubernetesClientAdapter) GetPodLogs(ctx context.Context, namespace, podName, containerName string, options *k8sPorts.LogOptions) (io.ReadCloser, error) {
	logOptions := &corev1.PodLogOptions{}

	if options != nil {
		logOptions.Follow = options.Follow
		logOptions.TailLines = options.TailLines
		logOptions.SinceSeconds = options.SinceSeconds
		logOptions.SinceTime = (*metav1.Time)(options.SinceTime)
		logOptions.Previous = options.Previous
		logOptions.Timestamps = options.Timestamps

		if containerName != "" {
			logOptions.Container = containerName
		}
	}

	req := kca.clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	return req.Stream(ctx)
}

// Conversion methods

func (kca *KubernetesClientAdapter) convertNode(node *corev1.Node) *k8sModels.Node {
	// Convert node status
	var status k8sModels.NodeStatus = k8sModels.NodeStatusUnknown
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				status = k8sModels.NodeStatusReady
			} else {
				status = k8sModels.NodeStatusNotReady
			}
			break
		}
	}

	// Extract roles
	var roles []string
	for label := range node.Labels {
		if label == "node-role.kubernetes.io/master" || label == "node-role.kubernetes.io/control-plane" {
			roles = append(roles, "master")
		}
		if label == "node-role.kubernetes.io/worker" {
			roles = append(roles, "worker")
		}
	}

	// Get IP addresses
	var internalIP, externalIP string
	for _, address := range node.Status.Addresses {
		switch address.Type {
		case corev1.NodeInternalIP:
			internalIP = address.Address
		case corev1.NodeExternalIP:
			externalIP = address.Address
		}
	}

	// Convert conditions
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

	return &k8sModels.Node{
		Name:       node.Name,
		Status:     status,
		Roles:      roles,
		Age:        time.Since(node.CreationTimestamp.Time).String(),
		Version:    node.Status.NodeInfo.KubeletVersion,
		InternalIP: internalIP,
		ExternalIP: externalIP,
		Conditions: conditions,
		Resources: k8sModels.NodeResources{
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
		},
		CreatedAt: node.CreationTimestamp.Time,
	}
}

func (kca *KubernetesClientAdapter) convertNamespace(namespace *corev1.Namespace) *k8sModels.Namespace {
	var status k8sModels.NamespaceStatus
	switch namespace.Status.Phase {
	case corev1.NamespaceActive:
		status = k8sModels.NamespaceStatusActive
	case corev1.NamespaceTerminating:
		status = k8sModels.NamespaceStatusTerminating
	default:
		status = k8sModels.NamespaceStatusActive
	}

	return &k8sModels.Namespace{
		Name:      namespace.Name,
		Status:    status,
		Labels:    namespace.Labels,
		Age:       time.Since(namespace.CreationTimestamp.Time).String(),
		CreatedAt: namespace.CreationTimestamp.Time,
	}
}

func (kca *KubernetesClientAdapter) convertDeployment(deployment *appsv1.Deployment) *k8sModels.Deployment {
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

	// Calculate pod info
	podInfo := k8sModels.PodInfo{
		Total:   int(deployment.Status.Replicas),
		Running: int(deployment.Status.ReadyReplicas),
		Pending: int(deployment.Status.Replicas - deployment.Status.ReadyReplicas),
		Failed:  0, // Would need to query pods for this
	}

	return &k8sModels.Deployment{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
		PodInfo:   podInfo,
		Replicas: k8sModels.DeploymentReplicas{
			Desired:   *deployment.Spec.Replicas,
			Current:   deployment.Status.Replicas,
			Ready:     deployment.Status.ReadyReplicas,
			Available: deployment.Status.AvailableReplicas,
		},
		Age:        time.Since(deployment.CreationTimestamp.Time).String(),
		CreatedAt:  deployment.CreationTimestamp.Time,
		Conditions: conditions,
	}
}

func (kca *KubernetesClientAdapter) convertPod(pod *corev1.Pod) *k8sModels.Pod {
	// Convert pod status
	var status k8sModels.PodStatus
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
	for _, container := range pod.Status.ContainerStatuses {
		containers = append(containers, k8sModels.Container{
			Name:         container.Name,
			Image:        container.Image,
			Ready:        container.Ready,
			RestartCount: container.RestartCount,
			State:        kca.convertContainerState(container.State),
		})
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

	// Calculate total restarts
	var totalRestarts int32
	for _, container := range containers {
		totalRestarts += container.RestartCount
	}

	return &k8sModels.Pod{
		Name:       pod.Name,
		Namespace:  pod.Namespace,
		Status:     status,
		Phase:      string(pod.Status.Phase),
		Node:       pod.Spec.NodeName,
		Age:        time.Since(pod.CreationTimestamp.Time).String(),
		Restarts:   totalRestarts,
		Ready:      kca.getPodReadyStatus(pod),
		IP:         pod.Status.PodIP,
		Containers: containers,
		Conditions: conditions,
		CreatedAt:  pod.CreationTimestamp.Time,
	}
}

func (kca *KubernetesClientAdapter) convertContainerState(state corev1.ContainerState) k8sModels.ContainerState {
	result := k8sModels.ContainerState{}

	if state.Running != nil {
		result.Running = &k8sModels.ContainerStateRunning{
			StartedAt: state.Running.StartedAt.Time,
		}
	}

	if state.Waiting != nil {
		result.Waiting = &k8sModels.ContainerStateWaiting{
			Reason:  state.Waiting.Reason,
			Message: state.Waiting.Message,
		}
	}

	if state.Terminated != nil {
		result.Terminated = &k8sModels.ContainerStateTerminated{
			ExitCode:   state.Terminated.ExitCode,
			Reason:     state.Terminated.Reason,
			Message:    state.Terminated.Message,
			StartedAt:  state.Terminated.StartedAt.Time,
			FinishedAt: state.Terminated.FinishedAt.Time,
		}
	}

	return result
}

func (kca *KubernetesClientAdapter) getPodReadyStatus(pod *corev1.Pod) string {
	readyContainers := 0
	totalContainers := len(pod.Status.ContainerStatuses)

	for _, container := range pod.Status.ContainerStatuses {
		if container.Ready {
			readyContainers++
		}
	}

	return fmt.Sprintf("%d/%d", readyContainers, totalContainers)
}

// buildConfig builds Kubernetes client configuration
func buildConfig(kubeconfig, context string) (*rest.Config, error) {
	if kubeconfig == "" {
		// Use in-cluster config
		return rest.InClusterConfig()
	}

	// Expand home directory
	if kubeconfig[0] == '~' {
		kubeconfig = filepath.Join(homeDir(), kubeconfig[1:])
	}

	// Load config from file
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Override context if specified
	if context != "" {
		config.CurrentContext = context
	}

	// Build REST config
	restConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build REST config: %w", err)
	}

	return restConfig, nil
}

// homeDir returns the home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // Windows
}
