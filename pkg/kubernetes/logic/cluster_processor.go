package logic

import (
	"context"
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
)

// ClusterProcessor handles cluster business logic
type ClusterProcessor struct {
	clusterRepo k8sPorts.ClusterRepository
}

// NewClusterProcessor creates a new cluster processor
func NewClusterProcessor(clusterRepo k8sPorts.ClusterRepository) *ClusterProcessor {
	return &ClusterProcessor{
		clusterRepo: clusterRepo,
	}
}

// ListClusters lists all configured clusters with their status
func (cp *ClusterProcessor) ListClusters(ctx context.Context) ([]k8sModels.Cluster, error) {
	clusters, err := cp.clusterRepo.ListClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	// Update cluster status by validating connectivity
	for i := range clusters {
		cluster := &clusters[i]

		// Validate cluster connectivity
		if err := cp.clusterRepo.ValidateCluster(ctx, cluster.Context); err != nil {
			cluster.Status = k8sModels.ClusterStatusError
		} else {
			cluster.Status = k8sModels.ClusterStatusConnected
		}
	}

	return clusters, nil
}

// GetClusterStatus gets the status of a specific cluster
func (cp *ClusterProcessor) GetClusterStatus(ctx context.Context, context string) (k8sModels.ClusterStatus, error) {
	// Validate cluster connectivity
	if err := cp.clusterRepo.ValidateCluster(ctx, context); err != nil {
		return k8sModels.ClusterStatusError, nil
	}

	return k8sModels.ClusterStatusConnected, nil
}

// ValidateClusterConfig validates cluster configuration
func (cp *ClusterProcessor) ValidateClusterConfig(cluster *k8sModels.Cluster) error {
	if err := cluster.Validate(); err != nil {
		return fmt.Errorf("cluster validation failed: %w", err)
	}

	return nil
}
