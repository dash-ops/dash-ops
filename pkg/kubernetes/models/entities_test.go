package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCluster_IsConnected_WithConnectedCluster_ReturnsTrue(t *testing.T) {
	// Arrange
	cluster := Cluster{
		Name:    "test-cluster",
		Context: "test-context",
		Status:  ClusterStatusConnected,
	}

	// Act
	result := cluster.IsConnected()

	// Assert
	assert.True(t, result)
}

func TestCluster_IsConnected_WithDisconnectedCluster_ReturnsFalse(t *testing.T) {
	// Arrange
	cluster := Cluster{
		Name:    "test-cluster",
		Context: "test-context",
		Status:  ClusterStatusDisconnected,
	}

	// Act
	result := cluster.IsConnected()

	// Assert
	assert.False(t, result)
}

func TestCluster_IsConnected_WithErrorCluster_ReturnsFalse(t *testing.T) {
	// Arrange
	cluster := Cluster{
		Name:    "test-cluster",
		Context: "test-context",
		Status:  ClusterStatusError,
	}

	// Act
	result := cluster.IsConnected()

	// Assert
	assert.False(t, result)
}

func TestCluster_Validate_WithValidCluster_ReturnsNoError(t *testing.T) {
	// Arrange
	cluster := Cluster{
		Name:    "test-cluster",
		Context: "test-context",
	}

	// Act
	err := cluster.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestCluster_Validate_WithMissingName_ReturnsError(t *testing.T) {
	// Arrange
	cluster := Cluster{
		Context: "test-context",
	}

	// Act
	err := cluster.Validate()

	// Assert
	assert.Error(t, err)
}

func TestCluster_Validate_WithMissingContext_ReturnsError(t *testing.T) {
	// Arrange
	cluster := Cluster{
		Name: "test-cluster",
	}

	// Act
	err := cluster.Validate()

	// Assert
	assert.Error(t, err)
}

func TestCluster_Validate_WithEmptyCluster_ReturnsError(t *testing.T) {
	// Arrange
	cluster := Cluster{}

	// Act
	err := cluster.Validate()

	// Assert
	assert.Error(t, err)
}

func TestNode_IsReady_WithReadyNode_ReturnsTrue(t *testing.T) {
	// Arrange
	node := Node{
		Name:   "test-node",
		Status: NodeStatusReady,
	}

	// Act
	result := node.IsReady()

	// Assert
	assert.True(t, result)
}

func TestNode_IsReady_WithNotReadyNode_ReturnsFalse(t *testing.T) {
	// Arrange
	node := Node{
		Name:   "test-node",
		Status: NodeStatusNotReady,
	}

	// Act
	result := node.IsReady()

	// Assert
	assert.False(t, result)
}

func TestNode_IsReady_WithUnknownStatusNode_ReturnsFalse(t *testing.T) {
	// Arrange
	node := Node{
		Name:   "test-node",
		Status: NodeStatusUnknown,
	}

	// Act
	result := node.IsReady()

	// Assert
	assert.False(t, result)
}

func TestNode_HasRole_WithMasterRole_ReturnsTrue(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "test-node",
		Roles: []string{"master", "control-plane"},
	}
	role := "master"

	// Act
	result := node.HasRole(role)

	// Assert
	assert.True(t, result)
}

func TestNode_HasRole_WithControlPlaneRole_ReturnsTrue(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "test-node",
		Roles: []string{"master", "control-plane"},
	}
	role := "control-plane"

	// Act
	result := node.HasRole(role)

	// Assert
	assert.True(t, result)
}

func TestNode_HasRole_WithWorkerRole_ReturnsFalse(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "test-node",
		Roles: []string{"master", "control-plane"},
	}
	role := "worker"

	// Act
	result := node.HasRole(role)

	// Assert
	assert.False(t, result)
}

func TestNode_HasRole_WithEmptyRole_ReturnsFalse(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "test-node",
		Roles: []string{"master", "control-plane"},
	}
	role := ""

	// Act
	result := node.HasRole(role)

	// Assert
	assert.False(t, result)
}

func TestNode_IsMaster_WithMasterNode_ReturnsTrue(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "master-node",
		Roles: []string{"master"},
	}

	// Act
	result := node.IsMaster()

	// Assert
	assert.True(t, result)
}

func TestNode_IsMaster_WithControlPlaneNode_ReturnsTrue(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "control-plane-node",
		Roles: []string{"control-plane"},
	}

	// Act
	result := node.IsMaster()

	// Assert
	assert.True(t, result)
}

func TestNode_IsMaster_WithWorkerNode_ReturnsFalse(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "worker-node",
		Roles: []string{"worker"},
	}

	// Act
	result := node.IsMaster()

	// Assert
	assert.False(t, result)
}

func TestNode_IsMaster_WithNodeWithNoRoles_ReturnsFalse(t *testing.T) {
	// Arrange
	node := Node{
		Name:  "no-role-node",
		Roles: []string{},
	}

	// Act
	result := node.IsMaster()

	// Assert
	assert.False(t, result)
}

func TestDeployment_IsHealthy_WithHealthyDeployment_ReturnsTrue(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Name: "test-deployment",
		Replicas: DeploymentReplicas{
			Desired: 3,
			Ready:   3,
		},
	}

	// Act
	result := deployment.IsHealthy()

	// Assert
	assert.True(t, result)
}

func TestDeployment_IsHealthy_WithDegradedDeployment_ReturnsFalse(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Name: "test-deployment",
		Replicas: DeploymentReplicas{
			Desired: 3,
			Ready:   2,
		},
	}

	// Act
	result := deployment.IsHealthy()

	// Assert
	assert.False(t, result)
}

func TestDeployment_IsHealthy_WithStoppedDeployment_ReturnsFalse(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Name: "test-deployment",
		Replicas: DeploymentReplicas{
			Desired: 0,
			Ready:   0,
		},
	}

	// Act
	result := deployment.IsHealthy()

	// Assert
	assert.False(t, result)
}

func TestDeployment_IsHealthy_WithFailedDeployment_ReturnsFalse(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Name: "test-deployment",
		Replicas: DeploymentReplicas{
			Desired: 3,
			Ready:   0,
		},
	}

	// Act
	result := deployment.IsHealthy()

	// Assert
	assert.False(t, result)
}

func TestDeployment_GetAvailabilityPercentage_With100PercentAvailable_Returns100(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Replicas: DeploymentReplicas{
			Desired: 3,
			Ready:   3,
		},
	}

	// Act
	result := deployment.GetAvailabilityPercentage()

	// Assert
	assert.Equal(t, 100.0, result)
}

func TestDeployment_GetAvailabilityPercentage_With66PercentAvailable_Returns66(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Replicas: DeploymentReplicas{
			Desired: 3,
			Ready:   2,
		},
	}

	// Act
	result := deployment.GetAvailabilityPercentage()

	// Assert
	assert.Equal(t, 66.66666666666666, result)
}

func TestDeployment_GetAvailabilityPercentage_With0PercentAvailable_Returns0(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Replicas: DeploymentReplicas{
			Desired: 3,
			Ready:   0,
		},
	}

	// Act
	result := deployment.GetAvailabilityPercentage()

	// Assert
	assert.Equal(t, 0.0, result)
}

func TestDeployment_GetAvailabilityPercentage_WithStoppedDeployment_Returns0(t *testing.T) {
	// Arrange
	deployment := Deployment{
		Replicas: DeploymentReplicas{
			Desired: 0,
			Ready:   0,
		},
	}

	// Act
	result := deployment.GetAvailabilityPercentage()

	// Assert
	assert.Equal(t, 0.0, result)
}

func TestPod_IsRunning_WithRunningPod_ReturnsTrue(t *testing.T) {
	// Arrange
	pod := Pod{
		Name:   "test-pod",
		Status: PodStatusRunning,
	}

	// Act
	result := pod.IsRunning()

	// Assert
	assert.True(t, result)
}

func TestPod_IsRunning_WithPendingPod_ReturnsFalse(t *testing.T) {
	// Arrange
	pod := Pod{
		Name:   "test-pod",
		Status: PodStatusPending,
	}

	// Act
	result := pod.IsRunning()

	// Assert
	assert.False(t, result)
}

func TestPod_IsRunning_WithFailedPod_ReturnsFalse(t *testing.T) {
	// Arrange
	pod := Pod{
		Name:   "test-pod",
		Status: PodStatusFailed,
	}

	// Act
	result := pod.IsRunning()

	// Assert
	assert.False(t, result)
}

func TestPod_IsReady_WithAllContainersReady_ReturnsTrue(t *testing.T) {
	// Arrange
	pod := Pod{
		Containers: []Container{
			{Name: "container1", Ready: true},
			{Name: "container2", Ready: true},
		},
	}

	// Act
	result := pod.IsReady()

	// Assert
	assert.True(t, result)
}

func TestPod_IsReady_WithSomeContainersNotReady_ReturnsFalse(t *testing.T) {
	// Arrange
	pod := Pod{
		Containers: []Container{
			{Name: "container1", Ready: true},
			{Name: "container2", Ready: false},
		},
	}

	// Act
	result := pod.IsReady()

	// Assert
	assert.False(t, result)
}

func TestPod_IsReady_WithNoContainers_ReturnsFalse(t *testing.T) {
	// Arrange
	pod := Pod{
		Containers: []Container{},
	}

	// Act
	result := pod.IsReady()

	// Assert
	assert.False(t, result)
}

func TestPod_GetTotalRestarts_WithMultipleContainers_ReturnsSum(t *testing.T) {
	// Arrange
	pod := Pod{
		Containers: []Container{
			{Name: "container1", RestartCount: 2},
			{Name: "container2", RestartCount: 3},
			{Name: "container3", RestartCount: 0},
		},
	}

	// Act
	result := pod.GetTotalRestarts()

	// Assert
	assert.Equal(t, int32(5), result)
}

func TestContainer_IsRunning_WithRunningContainer_ReturnsTrue(t *testing.T) {
	// Arrange
	container := Container{
		State: ContainerState{
			Running: &ContainerStateRunning{
				StartedAt: time.Now(),
			},
		},
	}

	// Act
	result := container.IsRunning()

	// Assert
	assert.True(t, result)
}

func TestContainer_IsRunning_WithWaitingContainer_ReturnsFalse(t *testing.T) {
	// Arrange
	container := Container{
		State: ContainerState{
			Waiting: &ContainerStateWaiting{
				Reason: "ImagePullBackOff",
			},
		},
	}

	// Act
	result := container.IsRunning()

	// Assert
	assert.False(t, result)
}

func TestContainer_IsRunning_WithTerminatedContainer_ReturnsFalse(t *testing.T) {
	// Arrange
	container := Container{
		State: ContainerState{
			Terminated: &ContainerStateTerminated{
				ExitCode: 0,
			},
		},
	}

	// Act
	result := container.IsRunning()

	// Assert
	assert.False(t, result)
}
