package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClient K8S module interface
type K8sClient interface {
	GetDeployments() ([]Deployment, error)
	Scale(name string, ns string, replicas int32) error
}

type k8sClient struct {
	clientSet *kubernetes.Clientset
}

// NewK8sClient Create a new k8s client
func NewK8sClient(config kubernetesConfig) (K8sClient, error) {
	kConfig, err := getConfig(config.Kubeconfig)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		return nil, err
	}

	return k8sClient{clientSet}, nil
}

func getConfig(config string) (*rest.Config, error) {
	if config == "" {
		return rest.InClusterConfig()
	}

	return clientcmd.BuildConfigFromFlags("", config)
}
