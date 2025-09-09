package logic

import (
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

// HealthCalculator provides health calculation logic for Kubernetes resources
type HealthCalculator struct{}

// NewHealthCalculator creates a new health calculator
func NewHealthCalculator() *HealthCalculator {
	return &HealthCalculator{}
}

// CalculateDeploymentHealth calculates health status for a deployment
func (hc *HealthCalculator) CalculateDeploymentHealth(deployment *k8sModels.Deployment) DeploymentHealthStatus {
	if deployment == nil {
		return DeploymentStatusUnknown
	}

	// Check if deployment has desired replicas
	if deployment.Replicas.Desired == 0 {
		return DeploymentStatusStopped
	}

	// Check if all replicas are ready
	if deployment.Replicas.Ready == deployment.Replicas.Desired {
		return DeploymentStatusHealthy
	}

	// Check if some replicas are ready
	if deployment.Replicas.Ready > 0 {
		return DeploymentStatusDegraded
	}

	// No replicas ready
	return DeploymentStatusUnhealthy
}

// CalculatePodHealth calculates health status for a pod
func (hc *HealthCalculator) CalculatePodHealth(pod *k8sModels.Pod) PodHealthStatus {
	if pod == nil {
		return PodHealthStatusUnknown
	}

	// Check pod phase
	switch pod.Status {
	case k8sModels.PodStatusRunning:
		if pod.IsReady() {
			return PodHealthStatusHealthy
		}
		return PodHealthStatusDegraded

	case k8sModels.PodStatusPending:
		return PodHealthStatusPending

	case k8sModels.PodStatusSucceeded:
		return PodHealthStatusCompleted

	case k8sModels.PodStatusFailed:
		return PodHealthStatusFailed

	default:
		return PodHealthStatusUnknown
	}
}

// CalculateNodeHealth calculates health status for a node
func (hc *HealthCalculator) CalculateNodeHealth(node *k8sModels.Node) NodeHealthStatus {
	if node == nil {
		return NodeHealthStatusUnknown
	}

	switch node.Status {
	case k8sModels.NodeStatusReady:
		return NodeHealthStatusHealthy
	case k8sModels.NodeStatusNotReady:
		return NodeHealthStatusUnhealthy
	default:
		return NodeHealthStatusUnknown
	}
}

// CalculateClusterHealth calculates overall cluster health
func (hc *HealthCalculator) CalculateClusterHealth(clusterInfo *k8sModels.ClusterInfo) ClusterHealthStatus {
	if clusterInfo == nil || !clusterInfo.Cluster.IsConnected() {
		return ClusterHealthStatusDisconnected
	}

	// Calculate node health ratio
	if len(clusterInfo.Nodes) == 0 {
		return ClusterHealthStatusUnknown
	}

	readyNodes := 0
	for _, node := range clusterInfo.Nodes {
		if node.IsReady() {
			readyNodes++
		}
	}

	healthRatio := float64(readyNodes) / float64(len(clusterInfo.Nodes))

	// Determine cluster health based on node health ratio
	switch {
	case healthRatio >= 0.9:
		return ClusterHealthStatusHealthy
	case healthRatio >= 0.7:
		return ClusterHealthStatusDegraded
	case healthRatio >= 0.5:
		return ClusterHealthStatusUnhealthy
	default:
		return ClusterHealthStatusCritical
	}
}

// GetDeploymentHealthSummary provides a summary of deployment health issues
func (hc *HealthCalculator) GetDeploymentHealthSummary(deployment *k8sModels.Deployment) *DeploymentHealthSummary {
	if deployment == nil {
		return &DeploymentHealthSummary{
			Status: DeploymentStatusUnknown,
			Issues: []string{"Deployment is nil"},
		}
	}

	summary := &DeploymentHealthSummary{
		Status: hc.CalculateDeploymentHealth(deployment),
		Issues: []string{},
	}

	// Identify specific issues
	if deployment.Replicas.Desired == 0 {
		summary.Issues = append(summary.Issues, "Deployment is scaled to 0 replicas")
	} else if deployment.Replicas.Ready == 0 {
		summary.Issues = append(summary.Issues, "No replicas are ready")
	} else if deployment.Replicas.Ready < deployment.Replicas.Desired {
		summary.Issues = append(summary.Issues,
			fmt.Sprintf("Only %d of %d replicas are ready",
				deployment.Replicas.Ready, deployment.Replicas.Desired))
	}

	// Check deployment conditions for additional issues
	for _, condition := range deployment.Conditions {
		if condition.Status == "False" && condition.Type == "Available" {
			summary.Issues = append(summary.Issues,
				fmt.Sprintf("Availability issue: %s", condition.Message))
		}
		if condition.Status == "False" && condition.Type == "Progressing" {
			summary.Issues = append(summary.Issues,
				fmt.Sprintf("Progress issue: %s", condition.Message))
		}
	}

	return summary
}

// DeploymentHealthSummary represents deployment health summary
type DeploymentHealthSummary struct {
	Status DeploymentHealthStatus `json:"status"`
	Issues []string               `json:"issues"`
}

// Health status enums (adding to models)

// DeploymentHealthStatus represents deployment health status
type DeploymentHealthStatus string

const (
	DeploymentStatusHealthy   DeploymentHealthStatus = "healthy"
	DeploymentStatusDegraded  DeploymentHealthStatus = "degraded"
	DeploymentStatusUnhealthy DeploymentHealthStatus = "unhealthy"
	DeploymentStatusStopped   DeploymentHealthStatus = "stopped"
	DeploymentStatusUnknown   DeploymentHealthStatus = "unknown"
)

// PodHealthStatus represents pod health status
type PodHealthStatus string

const (
	PodHealthStatusHealthy   PodHealthStatus = "healthy"
	PodHealthStatusDegraded  PodHealthStatus = "degraded"
	PodHealthStatusPending   PodHealthStatus = "pending"
	PodHealthStatusCompleted PodHealthStatus = "completed"
	PodHealthStatusFailed    PodHealthStatus = "failed"
	PodHealthStatusUnknown   PodHealthStatus = "unknown"
)

// NodeHealthStatus represents node health status
type NodeHealthStatus string

const (
	NodeHealthStatusHealthy   NodeHealthStatus = "healthy"
	NodeHealthStatusUnhealthy NodeHealthStatus = "unhealthy"
	NodeHealthStatusUnknown   NodeHealthStatus = "unknown"
)

// ClusterHealthStatus represents cluster health status
type ClusterHealthStatus string

const (
	ClusterHealthStatusHealthy      ClusterHealthStatus = "healthy"
	ClusterHealthStatusDegraded     ClusterHealthStatus = "degraded"
	ClusterHealthStatusUnhealthy    ClusterHealthStatus = "unhealthy"
	ClusterHealthStatusCritical     ClusterHealthStatus = "critical"
	ClusterHealthStatusDisconnected ClusterHealthStatus = "disconnected"
	ClusterHealthStatusUnknown      ClusterHealthStatus = "unknown"
)
