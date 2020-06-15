package kubernetes

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployment Struct representing an k8s deployment
type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	PodCount  int    `json:"pod_count"`
}

type deploymentFilter struct {
	Namespace string
}

func (kc k8sClient) GetDeployments(filter deploymentFilter) ([]Deployment, error) {
	var deployments []Deployment

	if filter.Namespace == "" {
		filter.Namespace = apiv1.NamespaceAll
	}

	deploys, err := kc.clientSet.
		AppsV1().
		Deployments(filter.Namespace).
		List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get deployments: %s", err)
	}

	for _, deploy := range deploys.Items {
		deployments = append(deployments, Deployment{
			Name:      deploy.GetName(),
			Namespace: deploy.GetNamespace(),
			PodCount:  int(*deploy.Spec.Replicas),
		})
	}

	return deployments, nil
}

func (kc k8sClient) Scale(name string, ns string, replicas int32) error {
	deploy, err := kc.clientSet.AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deploy %s on ns %s: %s", name, ns, err)
	}
	deploy.Spec.Replicas = &replicas
	deployment, err := kc.clientSet.AppsV1().Deployments(deploy.GetNamespace()).Update(deploy)
	fmt.Println(deployment)
	return err
}
