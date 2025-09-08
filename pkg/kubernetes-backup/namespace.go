package kubernetes

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Namespace Struct representing an k8s namespace
type Namespace struct {
	Name string `json:"name"`
}

func (kc client) GetNamespaces() ([]Namespace, error) {
	var list []Namespace

	namespaces, err := kc.clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to get namespace: %s", err)
	}

	for _, ns := range namespaces.Items {
		list = append(list, Namespace{
			Name: ns.GetName(),
		})
	}

	return list, nil
}
