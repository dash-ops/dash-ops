package external

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

// KubernetesClient wraps the Kubernetes client-go
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
	kConfig, err := getConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &KubernetesClient{
		clientSet: clientSet,
		context:   config.Context,
		config:    config,
	}, nil
}

func getConfig(config *KubernetesConfig) (*rest.Config, error) {
	if config.Kubeconfig == "" {
		return rest.InClusterConfig()
	}

	if config.Context == "" {
		return clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.Kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: config.Context,
		}).ClientConfig()
}

// GetVersion gets the Kubernetes cluster version
func (kc *KubernetesClient) GetVersion(ctx context.Context) (string, error) {
	versionInfo, err := kc.clientSet.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get server version: %w", err)
	}
	return versionInfo.String(), nil
}

// ListNodes lists all nodes in the cluster
func (kc *KubernetesClient) ListNodes(ctx context.Context) ([]corev1.Node, error) {
	nodeList, err := kc.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}
	return nodeList.Items, nil
}

// GetNode gets a specific node
func (kc *KubernetesClient) GetNode(ctx context.Context, name string) (*corev1.Node, error) {
	node, err := kc.clientSet.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %w", name, err)
	}
	return node, nil
}

// ListNamespaces lists all namespaces
func (kc *KubernetesClient) ListNamespaces(ctx context.Context) ([]corev1.Namespace, error) {
	nsList, err := kc.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	return nsList.Items, nil
}

// GetNamespace gets a specific namespace
func (kc *KubernetesClient) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	ns, err := kc.clientSet.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace %s: %w", name, err)
	}
	return ns, nil
}

// CreateNamespace creates a new namespace
func (kc *KubernetesClient) CreateNamespace(ctx context.Context, name string, labels map[string]string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}

	createdNs, err := kc.clientSet.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace %s: %w", name, err)
	}
	return createdNs, nil
}

// DeleteNamespace deletes a namespace
func (kc *KubernetesClient) DeleteNamespace(ctx context.Context, name string) error {
	err := kc.clientSet.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete namespace %s: %w", name, err)
	}
	return nil
}

// ListDeployments lists deployments with optional namespace filter
func (kc *KubernetesClient) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	var deployments []appsv1.Deployment

	if namespace != "" {
		deploymentList, err := kc.clientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
		}
		deployments = deploymentList.Items
	} else {
		// List all deployments across all namespaces
		deploymentList, err := kc.clientSet.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list all deployments: %w", err)
		}
		deployments = deploymentList.Items
	}

	return deployments, nil
}

// GetDeployment gets a specific deployment
func (kc *KubernetesClient) GetDeployment(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	deployment, err := kc.clientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s/%s: %w", namespace, name, err)
	}
	return deployment, nil
}

// ScaleDeployment scales a deployment to specified replicas
func (kc *KubernetesClient) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	deployment, err := kc.GetDeployment(ctx, namespace, name)
	if err != nil {
		return err
	}

	deployment.Spec.Replicas = &replicas
	_, err = kc.clientSet.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to scale deployment %s/%s: %w", namespace, name, err)
	}
	return nil
}

// RestartDeployment restarts a deployment by updating annotations
func (kc *KubernetesClient) RestartDeployment(ctx context.Context, namespace, name string) error {
	deployment, err := kc.GetDeployment(ctx, namespace, name)
	if err != nil {
		return err
	}

	// Add or update restart annotation
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = kc.clientSet.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to restart deployment %s/%s: %w", namespace, name, err)
	}
	return nil
}

// ListPods lists pods with optional namespace filter
func (kc *KubernetesClient) ListPods(ctx context.Context, namespace string, labelSelector string) ([]corev1.Pod, error) {
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	var pods []corev1.Pod
	if namespace != "" {
		podList, err := kc.clientSet.CoreV1().Pods(namespace).List(ctx, listOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
		}
		pods = podList.Items
	} else {
		// List all pods across all namespaces
		podList, err := kc.clientSet.CoreV1().Pods("").List(ctx, listOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list all pods: %w", err)
		}
		pods = podList.Items
	}

	return pods, nil
}

// GetPod gets a specific pod
func (kc *KubernetesClient) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	pod, err := kc.clientSet.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s/%s: %w", namespace, name, err)
	}
	return pod, nil
}

// DeletePod deletes a pod
func (kc *KubernetesClient) DeletePod(ctx context.Context, namespace, name string) error {
	err := kc.clientSet.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod %s/%s: %w", namespace, name, err)
	}
	return nil
}

// GetPodLogs gets logs for a pod
func (kc *KubernetesClient) GetPodLogs(ctx context.Context, namespace, podName, containerName string, tailLines int64) (string, error) {
	podLogOptions := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &tailLines,
	}

	req := kc.clientSet.CoreV1().Pods(namespace).GetLogs(podName, podLogOptions)
	logs, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get logs for pod %s/%s: %w", namespace, podName, err)
	}
	defer logs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return buf.String(), nil
}
