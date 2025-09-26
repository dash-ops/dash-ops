package controllers

import (
	"context"
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	"github.com/dash-ops/dash-ops/pkg/kubernetes/repositories"
)

// NodesController handles nodes business logic orchestration
type NodesController struct {
	repository *repositories.NodesRepository
}

// NewNodesController creates a new nodes controller
func NewNodesController(repository *repositories.NodesRepository) *NodesController {
	return &NodesController{
		repository: repository,
	}
}

// GetNode gets a specific node with business logic validation
func (c *NodesController) GetNode(ctx context.Context, context, nodeName string) (*k8sModels.Node, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	node, err := c.repository.GetNode(ctx, context, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	return node, nil
}

// ListNodes lists all nodes in a cluster with optional aggregation
func (c *NodesController) ListNodes(ctx context.Context, context string) ([]k8sModels.Node, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	nodes, err := c.repository.ListNodes(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// Apply business logic: sort nodes by name for consistent ordering
	nodes = c.sortNodesByName(nodes)

	return nodes, nil
}

// GetNodeMetrics gets node resource metrics with business logic processing
func (c *NodesController) GetNodeMetrics(ctx context.Context, context, nodeName string) (*k8sModels.NodeResources, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}

	// First verify the node exists
	_, err := c.repository.GetNode(ctx, context, nodeName)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	metrics, err := c.repository.GetNodeMetrics(ctx, context, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics: %w", err)
	}

	return metrics, nil
}

// GetNodesSummary provides a summary of nodes in the cluster
func (c *NodesController) GetNodesSummary(ctx context.Context, context string) (*NodesSummary, error) {
	if context == "" {
		return nil, fmt.Errorf("context is required")
	}

	nodes, err := c.repository.ListNodes(ctx, context)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes for summary: %w", err)
	}

	summary := &NodesSummary{
		Total:    len(nodes),
		Ready:    0,
		NotReady: 0,
		Master:   0,
		Worker:   0,
	}

	for _, node := range nodes {
		if node.IsReady() {
			summary.Ready++
		} else {
			summary.NotReady++
		}

		if node.IsMaster() {
			summary.Master++
		} else {
			summary.Worker++
		}
	}

	return summary, nil
}

// sortNodesByName sorts nodes by name for consistent ordering
func (c *NodesController) sortNodesByName(nodes []k8sModels.Node) []k8sModels.Node {
	// Simple bubble sort for small lists (nodes are typically < 100)
	for i := 0; i < len(nodes)-1; i++ {
		for j := 0; j < len(nodes)-i-1; j++ {
			if nodes[j].Name > nodes[j+1].Name {
				nodes[j], nodes[j+1] = nodes[j+1], nodes[j]
			}
		}
	}
	return nodes
}

// NodesSummary represents a summary of nodes in a cluster
type NodesSummary struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
	Master   int `json:"master"`
	Worker   int `json:"worker"`
}
