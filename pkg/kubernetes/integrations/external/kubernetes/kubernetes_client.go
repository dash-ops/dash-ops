package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesClient handles communication with Kubernetes API
type KubernetesClient struct {
	clientSet *kubernetes.Clientset
	context   string
	config    *KubernetesConfig
}

// KubernetesConfig represents Kubernetes connection configuration
type KubernetesConfig struct {
	Kubeconfig string
	Context    string
}

// NewKubernetesClient creates a new Kubernetes client
func NewKubernetesClient(config *KubernetesConfig) (*KubernetesClient, error) {
	var restConfig *rest.Config
	var err error

	if config.Kubeconfig != "" {
		// Use kubeconfig file
		restConfig, err = clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
		}
	} else {
		// Use in-cluster config
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
		}
	}

	// Override context if specified
	if config.Context != "" {
		restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.Kubeconfig},
			&clientcmd.ConfigOverrides{CurrentContext: config.Context},
		).ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to override context: %w", err)
		}
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &KubernetesClient{
		clientSet: clientSet,
		context:   config.Context,
		config:    config,
	}, nil
}

// GetDeployment gets a deployment
func (kc *KubernetesClient) GetDeployment(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	return kc.clientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

// ListDeployments lists deployments
func (kc *KubernetesClient) ListDeployments(ctx context.Context, namespace string, options metav1.ListOptions) (*appsv1.DeploymentList, error) {
	return kc.clientSet.AppsV1().Deployments(namespace).List(ctx, options)
}

// GetPod gets a pod
func (kc *KubernetesClient) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return kc.clientSet.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

// ListPods lists pods
func (kc *KubernetesClient) ListPods(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.PodList, error) {
	return kc.clientSet.CoreV1().Pods(namespace).List(ctx, options)
}

// GetNode gets a node
func (kc *KubernetesClient) GetNode(ctx context.Context, name string) (*corev1.Node, error) {
	return kc.clientSet.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
}

// ListNodes lists nodes
func (kc *KubernetesClient) ListNodes(ctx context.Context, options metav1.ListOptions) (*corev1.NodeList, error) {
	return kc.clientSet.CoreV1().Nodes().List(ctx, options)
}

// GetNamespace gets a namespace
func (kc *KubernetesClient) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	return kc.clientSet.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
}

// ListNamespaces lists namespaces
func (kc *KubernetesClient) ListNamespaces(ctx context.Context, options metav1.ListOptions) (*corev1.NamespaceList, error) {
	return kc.clientSet.CoreV1().Namespaces().List(ctx, options)
}

// GetPodLogs gets pod logs
func (kc *KubernetesClient) GetPodLogs(ctx context.Context, namespace, podName string, options *corev1.PodLogOptions) (io.ReadCloser, error) {
	return kc.clientSet.CoreV1().Pods(namespace).GetLogs(podName, options).Stream(ctx)
}

// GetPodLogsString gets pod logs as string
func (kc *KubernetesClient) GetPodLogsString(ctx context.Context, namespace, podName string, options *corev1.PodLogOptions) (string, error) {
	stream, err := kc.GetPodLogs(ctx, namespace, podName, options)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, stream)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ScaleDeployment scales a deployment
func (kc *KubernetesClient) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	scale, err := kc.clientSet.AppsV1().Deployments(namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment scale: %w", err)
	}

	scale.Spec.Replicas = replicas
	_, err = kc.clientSet.AppsV1().Deployments(namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to scale deployment: %w", err)
	}

	return nil
}

// RestartDeployment restarts a deployment
func (kc *KubernetesClient) RestartDeployment(ctx context.Context, namespace, name string) error {
	deployment, err := kc.GetDeployment(ctx, namespace, name)
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// Add restart annotation
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = kc.clientSet.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to restart deployment: %w", err)
	}

	return nil
}

// DeletePod deletes a pod
func (kc *KubernetesClient) DeletePod(ctx context.Context, namespace, name string) error {
	return kc.clientSet.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// CreateNamespace creates a namespace
func (kc *KubernetesClient) CreateNamespace(ctx context.Context, name string, labels map[string]string) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}

	return kc.clientSet.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
}

// DeleteNamespace deletes a namespace
func (kc *KubernetesClient) DeleteNamespace(ctx context.Context, name string) error {
	return kc.clientSet.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}

// GetContext returns the current context
func (kc *KubernetesClient) GetContext() string {
	return kc.context
}

// GetConfig returns the client configuration
func (kc *KubernetesClient) GetConfig() *KubernetesConfig {
	return kc.config
}

// TestConnection tests the connection to Kubernetes API
func (kc *KubernetesClient) TestConnection(ctx context.Context) error {
	_, err := kc.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	return err
}

// GetServerVersion gets the Kubernetes server version
func (kc *KubernetesClient) GetServerVersion(ctx context.Context) (string, error) {
	version, err := kc.clientSet.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return version.String(), nil
}

// GetClusterInfo gets basic cluster information
func (kc *KubernetesClient) GetClusterInfo(ctx context.Context) (*ClusterInfo, error) {
	// Get server version
	version, err := kc.GetServerVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	// Get nodes count
	nodes, err := kc.ListNodes(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// Get namespaces count
	namespaces, err := kc.ListNamespaces(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	return &ClusterInfo{
		Version:        version,
		Context:        kc.context,
		NodeCount:      len(nodes.Items),
		NamespaceCount: len(namespaces.Items),
		LastChecked:    time.Now(),
	}, nil
}

// ClusterInfo represents basic cluster information
type ClusterInfo struct {
	Version        string    `json:"version"`
	Context        string    `json:"context"`
	NodeCount      int       `json:"node_count"`
	NamespaceCount int       `json:"namespace_count"`
	LastChecked    time.Time `json:"last_checked"`
}
