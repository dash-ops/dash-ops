package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client K8S module interface
type Client interface {
	GetNodes() ([]Node, error)
	GetNamespaces() ([]Namespace, error)
	GetDeployments(filters deploymentFilter) ([]Deployment, error)
	GetDeploymentsWithContext(filters deploymentFilter, resolver ServiceContextResolver) ([]Deployment, error)
	Scale(name string, ns string, replicas int32) error
	RestartDeployment(name string, ns string) error
	GetPods(filters podFilter) ([]Pod, error)
	GetPodLogs(filters podFilter) ([]ContainerLog, error)
}

type client struct {
	clientSet *kubernetes.Clientset
	context   string
}

// NewClient Create a new k8s client
func NewClient(config config) (Client, error) {
	kConfig, err := getConfig(config)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		return nil, err
	}

	return client{
		clientSet: clientSet,
		context:   config.Context,
	}, nil
}

func getConfig(config config) (*rest.Config, error) {
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
