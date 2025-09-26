package controllers

import (
	"context"
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/repositories"
)

// NamespacesController handles namespaces business logic orchestration
type NamespacesController struct {
	repository *repositories.NamespacesRepository
}

// NewNamespacesController creates a new namespaces controller
func NewNamespacesController(repository *repositories.NamespacesRepository) *NamespacesController {
	return &NamespacesController{
		repository: repository,
	}
}

// GetNamespace gets a specific namespace with business logic validation
func (c *NamespacesController) GetNamespace(ctx context.Context, context, namespaceName string) (*k8sModels.Namespace, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if namespaceName == "" {
		return nil, fmt.Errorf("namespace name is required")
	}

	namespace, err := c.repository.GetNamespace(ctx, context, namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	return namespace, nil
}

// ListNamespaces lists all namespaces with business logic processing
func (c *NamespacesController) ListNamespaces(ctx context.Context, context string) ([]k8sModels.Namespace, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	namespaces, err := c.repository.ListNamespaces(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Apply business logic: sort namespaces by name for consistent ordering
	namespaces = c.sortNamespacesByName(namespaces)

	return namespaces, nil
}

// CreateNamespace creates a new namespace with business logic validation
func (c *NamespacesController) CreateNamespace(ctx context.Context, context, name string) (*k8sModels.Namespace, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if name == "" {
		return nil, fmt.Errorf("namespace name is required")
	}

	// Business logic: validate namespace name
	if err := c.validateNamespaceName(name); err != nil {
		return nil, fmt.Errorf("invalid namespace name: %w", err)
	}

	// Business logic: prevent creation of system namespaces
	if c.isSystemNamespace(name) {
		return nil, fmt.Errorf("cannot create namespace with system name: %s", name)
	}

	// Create namespace with default labels
	defaultLabels := map[string]string{
		"created-by": "dashops",
		"managed-by": "kubernetes-controller",
	}

	namespace, err := c.repository.CreateNamespace(ctx, context, name, defaultLabels)
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace: %w", err)
	}

	return namespace, nil
}

// DeleteNamespace deletes a namespace with business logic validation
func (c *NamespacesController) DeleteNamespace(ctx context.Context, context, name string) error {
	if context == "" {
		return fmt.Errorf("context is required")
	}
	if name == "" {
		return fmt.Errorf("namespace name is required")
	}

	// Business logic: prevent deletion of system namespaces
	if c.isSystemNamespace(name) {
		return fmt.Errorf("cannot delete system namespace: %s", name)
	}

	// Verify namespace exists before deletion
	_, err := c.repository.GetNamespace(ctx, context, name)
	if err != nil {
		return fmt.Errorf("namespace not found: %w", err)
	}

	err = c.repository.DeleteNamespace(ctx, context, name)
	if err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	return nil
}

// GetNamespacesSummary provides a summary of namespaces in a cluster
func (c *NamespacesController) GetNamespacesSummary(ctx context.Context, context string) (*NamespacesSummary, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	namespaces, err := c.repository.ListNamespaces(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces for summary: %w", err)
	}

	summary := &NamespacesSummary{
		Total:       len(namespaces),
		Active:      0,
		Terminating: 0,
		System:      0,
		User:        0,
	}

	for _, namespace := range namespaces {
		if namespace.Status == k8sModels.NamespaceStatusActive {
			summary.Active++
		} else if namespace.Status == k8sModels.NamespaceStatusTerminating {
			summary.Terminating++
		}

		if c.isSystemNamespace(namespace.Name) {
			summary.System++
		} else {
			summary.User++
		}
	}

	return summary, nil
}

// validateNamespaceName validates namespace name according to Kubernetes rules
func (c *NamespacesController) validateNamespaceName(name string) error {
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
func (c *NamespacesController) isSystemNamespace(name string) bool {
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

// sortNamespacesByName sorts namespaces by name for consistent ordering
func (c *NamespacesController) sortNamespacesByName(namespaces []k8sModels.Namespace) []k8sModels.Namespace {
	// Simple bubble sort for small lists (namespaces are typically < 50)
	for i := 0; i < len(namespaces)-1; i++ {
		for j := 0; j < len(namespaces)-i-1; j++ {
			if namespaces[j].Name > namespaces[j+1].Name {
				namespaces[j], namespaces[j+1] = namespaces[j+1], namespaces[j]
			}
		}
	}
	return namespaces
}

// NamespacesSummary represents a summary of namespaces in a cluster
type NamespacesSummary struct {
	Total       int `json:"total"`
	Active      int `json:"active"`
	Terminating int `json:"terminating"`
	System      int `json:"system"`
	User        int `json:"user"`
}
