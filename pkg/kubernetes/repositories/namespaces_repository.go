package repositories

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dash-ops/dash-ops/pkg/kubernetes/integrations/external/kubernetes"
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

// NamespacesRepository handles namespace-related data access
type NamespacesRepository struct {
	client *kubernetes.KubernetesClient
}

// NewNamespacesRepository creates a new namespaces repository
func NewNamespacesRepository(client *kubernetes.KubernetesClient) *NamespacesRepository {
	return &NamespacesRepository{
		client: client,
	}
}

// GetNamespace gets a specific namespace
func (r *NamespacesRepository) GetNamespace(ctx context.Context, context, namespaceName string) (*k8sModels.Namespace, error) {
	if namespaceName == "" {
		return nil, fmt.Errorf("namespace name is required")
	}

	namespace, err := r.client.GetNamespace(ctx, namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace %s: %w", namespaceName, err)
	}

	return r.convertNamespace(namespace), nil
}

// ListNamespaces lists all namespaces in a cluster
func (r *NamespacesRepository) ListNamespaces(ctx context.Context, context string) ([]k8sModels.Namespace, error) {
	namespaceList, err := r.client.ListNamespaces(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var namespaces []k8sModels.Namespace
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, *r.convertNamespace(&namespace))
	}

	return namespaces, nil
}

// CreateNamespace creates a new namespace
func (r *NamespacesRepository) CreateNamespace(ctx context.Context, context, namespaceName string, labels map[string]string) (*k8sModels.Namespace, error) {
	if namespaceName == "" {
		return nil, fmt.Errorf("namespace name is required")
	}

	// Validate namespace name
	if err := validateNamespaceName(namespaceName); err != nil {
		return nil, fmt.Errorf("invalid namespace name: %w", err)
	}

	createdNamespace, err := r.client.CreateNamespace(ctx, namespaceName, labels)
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace %s: %w", namespaceName, err)
	}

	return r.convertNamespace(createdNamespace), nil
}

// DeleteNamespace deletes a namespace
func (r *NamespacesRepository) DeleteNamespace(ctx context.Context, context, namespaceName string) error {
	if namespaceName == "" {
		return fmt.Errorf("namespace name is required")
	}

	// Prevent deletion of system namespaces
	if isSystemNamespace(namespaceName) {
		return fmt.Errorf("cannot delete system namespace: %s", namespaceName)
	}

	err := r.client.DeleteNamespace(ctx, namespaceName)
	if err != nil {
		return fmt.Errorf("failed to delete namespace %s: %w", namespaceName, err)
	}

	return nil
}

// validateNamespaceName validates namespace name according to Kubernetes rules
func validateNamespaceName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("namespace name cannot be empty")
	}
	if len(name) > 63 {
		return fmt.Errorf("namespace name cannot be longer than 63 characters")
	}

	// Check for valid characters (alphanumeric and hyphens)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("namespace name contains invalid character: %c", char)
		}
	}

	// Cannot start or end with hyphen
	if name[0] == '-' || name[len(name)-1] == '-' {
		return fmt.Errorf("namespace name cannot start or end with hyphen")
	}

	return nil
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

// convertNamespace converts a Kubernetes namespace to our domain model
func (r *NamespacesRepository) convertNamespace(namespace *corev1.Namespace) *k8sModels.Namespace {
	// Determine namespace status
	var status k8sModels.NamespaceStatus
	if namespace.Status.Phase == corev1.NamespaceActive {
		status = k8sModels.NamespaceStatusActive
	} else if namespace.Status.Phase == corev1.NamespaceTerminating {
		status = k8sModels.NamespaceStatusTerminating
	} else {
		status = k8sModels.NamespaceStatusActive // Default to active for other phases
	}

	// Convert labels
	labels := make(map[string]string)
	if namespace.Labels != nil {
		for key, value := range namespace.Labels {
			labels[key] = value
		}
	}

	return &k8sModels.Namespace{
		Name:      namespace.Name,
		Status:    status,
		Labels:    labels,
		Age:       time.Since(namespace.CreationTimestamp.Time).Round(time.Second).String(),
		CreatedAt: namespace.CreationTimestamp.Time,
	}
}
