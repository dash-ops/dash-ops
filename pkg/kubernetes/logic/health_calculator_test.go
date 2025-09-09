package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

func TestHealthCalculator_CalculateDeploymentHealth_WithNilDeployment_ReturnsUnknown(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := (*k8sModels.Deployment)(nil)

	// Act
	result := calculator.CalculateDeploymentHealth(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusUnknown, result)
}

func TestHealthCalculator_CalculateDeploymentHealth_WithHealthyDeployment_ReturnsHealthy(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := &k8sModels.Deployment{
		Replicas: k8sModels.DeploymentReplicas{
			Desired: 3,
			Ready:   3,
		},
	}

	// Act
	result := calculator.CalculateDeploymentHealth(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusHealthy, result)
}

func TestHealthCalculator_CalculateDeploymentHealth_WithDegradedDeployment_ReturnsDegraded(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := &k8sModels.Deployment{
		Replicas: k8sModels.DeploymentReplicas{
			Desired: 3,
			Ready:   2,
		},
	}

	// Act
	result := calculator.CalculateDeploymentHealth(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusDegraded, result)
}

func TestHealthCalculator_CalculateDeploymentHealth_WithUnhealthyDeployment_ReturnsUnhealthy(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := &k8sModels.Deployment{
		Replicas: k8sModels.DeploymentReplicas{
			Desired: 3,
			Ready:   0,
		},
	}

	// Act
	result := calculator.CalculateDeploymentHealth(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusUnhealthy, result)
}

func TestHealthCalculator_CalculateDeploymentHealth_WithStoppedDeployment_ReturnsStopped(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := &k8sModels.Deployment{
		Replicas: k8sModels.DeploymentReplicas{
			Desired: 0,
			Ready:   0,
		},
	}

	// Act
	result := calculator.CalculateDeploymentHealth(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusStopped, result)
}

func TestHealthCalculator_CalculatePodHealth_WithNilPod_ReturnsUnknown(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	pod := (*k8sModels.Pod)(nil)

	// Act
	result := calculator.CalculatePodHealth(pod)

	// Assert
	assert.Equal(t, PodHealthStatusUnknown, result)
}

func TestHealthCalculator_CalculatePodHealth_WithHealthyRunningPod_ReturnsHealthy(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	pod := &k8sModels.Pod{
		Status: k8sModels.PodStatusRunning,
		Containers: []k8sModels.Container{
			{Ready: true},
			{Ready: true},
		},
	}

	// Act
	result := calculator.CalculatePodHealth(pod)

	// Assert
	assert.Equal(t, PodHealthStatusHealthy, result)
}

func TestHealthCalculator_CalculatePodHealth_WithDegradedRunningPod_ReturnsDegraded(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	pod := &k8sModels.Pod{
		Status: k8sModels.PodStatusRunning,
		Containers: []k8sModels.Container{
			{Ready: true},
			{Ready: false},
		},
	}

	// Act
	result := calculator.CalculatePodHealth(pod)

	// Assert
	assert.Equal(t, PodHealthStatusDegraded, result)
}

func TestHealthCalculator_CalculatePodHealth_WithPendingPod_ReturnsPending(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	pod := &k8sModels.Pod{
		Status: k8sModels.PodStatusPending,
	}

	// Act
	result := calculator.CalculatePodHealth(pod)

	// Assert
	assert.Equal(t, PodHealthStatusPending, result)
}

func TestHealthCalculator_CalculatePodHealth_WithFailedPod_ReturnsFailed(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	pod := &k8sModels.Pod{
		Status: k8sModels.PodStatusFailed,
	}

	// Act
	result := calculator.CalculatePodHealth(pod)

	// Assert
	assert.Equal(t, PodHealthStatusFailed, result)
}

func TestHealthCalculator_CalculatePodHealth_WithSucceededPod_ReturnsCompleted(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	pod := &k8sModels.Pod{
		Status: k8sModels.PodStatusSucceeded,
	}

	// Act
	result := calculator.CalculatePodHealth(pod)

	// Assert
	assert.Equal(t, PodHealthStatusCompleted, result)
}

func TestHealthCalculator_CalculateNodeHealth_WithNilNode_ReturnsUnknown(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	node := (*k8sModels.Node)(nil)

	// Act
	result := calculator.CalculateNodeHealth(node)

	// Assert
	assert.Equal(t, NodeHealthStatusUnknown, result)
}

func TestHealthCalculator_CalculateNodeHealth_WithReadyNode_ReturnsHealthy(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	node := &k8sModels.Node{
		Status: k8sModels.NodeStatusReady,
	}

	// Act
	result := calculator.CalculateNodeHealth(node)

	// Assert
	assert.Equal(t, NodeHealthStatusHealthy, result)
}

func TestHealthCalculator_CalculateNodeHealth_WithNotReadyNode_ReturnsUnhealthy(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	node := &k8sModels.Node{
		Status: k8sModels.NodeStatusNotReady,
	}

	// Act
	result := calculator.CalculateNodeHealth(node)

	// Assert
	assert.Equal(t, NodeHealthStatusUnhealthy, result)
}

func TestHealthCalculator_CalculateNodeHealth_WithUnknownStatusNode_ReturnsUnknown(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	node := &k8sModels.Node{
		Status: k8sModels.NodeStatusUnknown,
	}

	// Act
	result := calculator.CalculateNodeHealth(node)

	// Assert
	assert.Equal(t, NodeHealthStatusUnknown, result)
}

func TestHealthCalculator_CalculateClusterHealth_WithNilClusterInfo_ReturnsDisconnected(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	clusterInfo := (*k8sModels.ClusterInfo)(nil)

	// Act
	result := calculator.CalculateClusterHealth(clusterInfo)

	// Assert
	assert.Equal(t, ClusterHealthStatusDisconnected, result)
}

func TestHealthCalculator_CalculateClusterHealth_WithDisconnectedCluster_ReturnsDisconnected(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	clusterInfo := &k8sModels.ClusterInfo{
		Cluster: k8sModels.Cluster{
			Status: k8sModels.ClusterStatusDisconnected,
		},
	}

	// Act
	result := calculator.CalculateClusterHealth(clusterInfo)

	// Assert
	assert.Equal(t, ClusterHealthStatusDisconnected, result)
}

func TestHealthCalculator_CalculateClusterHealth_WithHealthyCluster_ReturnsHealthy(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	clusterInfo := &k8sModels.ClusterInfo{
		Cluster: k8sModels.Cluster{
			Status: k8sModels.ClusterStatusConnected,
		},
		Nodes: []k8sModels.Node{
			{Status: k8sModels.NodeStatusReady},
			{Status: k8sModels.NodeStatusReady},
			{Status: k8sModels.NodeStatusReady},
		},
	}

	// Act
	result := calculator.CalculateClusterHealth(clusterInfo)

	// Assert
	assert.Equal(t, ClusterHealthStatusHealthy, result)
}

func TestHealthCalculator_CalculateClusterHealth_WithDegradedCluster_ReturnsDegraded(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	clusterInfo := &k8sModels.ClusterInfo{
		Cluster: k8sModels.Cluster{
			Status: k8sModels.ClusterStatusConnected,
		},
		Nodes: []k8sModels.Node{
			{Status: k8sModels.NodeStatusReady},
			{Status: k8sModels.NodeStatusReady},
			{Status: k8sModels.NodeStatusReady},
			{Status: k8sModels.NodeStatusReady},
			{Status: k8sModels.NodeStatusNotReady}, // 80% ready
		},
	}

	// Act
	result := calculator.CalculateClusterHealth(clusterInfo)

	// Assert
	assert.Equal(t, ClusterHealthStatusDegraded, result)
}

func TestHealthCalculator_CalculateClusterHealth_WithCriticalCluster_ReturnsCritical(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	clusterInfo := &k8sModels.ClusterInfo{
		Cluster: k8sModels.Cluster{
			Status: k8sModels.ClusterStatusConnected,
		},
		Nodes: []k8sModels.Node{
			{Status: k8sModels.NodeStatusReady},
			{Status: k8sModels.NodeStatusNotReady},
			{Status: k8sModels.NodeStatusNotReady},
			{Status: k8sModels.NodeStatusNotReady},
			{Status: k8sModels.NodeStatusNotReady}, // 20% ready
		},
	}

	// Act
	result := calculator.CalculateClusterHealth(clusterInfo)

	// Assert
	assert.Equal(t, ClusterHealthStatusCritical, result)
}

func TestHealthCalculator_GetDeploymentHealthSummary_WithNilDeployment_ReturnsUnknownWithIssues(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := (*k8sModels.Deployment)(nil)

	// Act
	summary := calculator.GetDeploymentHealthSummary(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusUnknown, summary.Status)
	assert.Len(t, summary.Issues, 1)
}

func TestHealthCalculator_GetDeploymentHealthSummary_WithHealthyDeployment_ReturnsHealthyWithNoIssues(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := &k8sModels.Deployment{
		Replicas: k8sModels.DeploymentReplicas{
			Desired: 3,
			Ready:   3,
		},
	}

	// Act
	summary := calculator.GetDeploymentHealthSummary(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusHealthy, summary.Status)
	assert.Len(t, summary.Issues, 0)
}

func TestHealthCalculator_GetDeploymentHealthSummary_WithScaledToZeroDeployment_ReturnsStoppedWithIssues(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := &k8sModels.Deployment{
		Replicas: k8sModels.DeploymentReplicas{
			Desired: 0,
			Ready:   0,
		},
	}

	// Act
	summary := calculator.GetDeploymentHealthSummary(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusStopped, summary.Status)
	assert.Len(t, summary.Issues, 1)
}

func TestHealthCalculator_GetDeploymentHealthSummary_WithNoReplicasReadyDeployment_ReturnsUnhealthyWithIssues(t *testing.T) {
	// Arrange
	calculator := NewHealthCalculator()
	deployment := &k8sModels.Deployment{
		Replicas: k8sModels.DeploymentReplicas{
			Desired: 3,
			Ready:   0,
		},
	}

	// Act
	summary := calculator.GetDeploymentHealthSummary(deployment)

	// Assert
	assert.Equal(t, DeploymentStatusUnhealthy, summary.Status)
	assert.Len(t, summary.Issues, 1)
}
