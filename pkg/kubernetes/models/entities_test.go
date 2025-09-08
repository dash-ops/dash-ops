package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCluster_IsConnected(t *testing.T) {
	tests := []struct {
		name     string
		cluster  Cluster
		expected bool
	}{
		{
			name: "connected cluster",
			cluster: Cluster{
				Name:    "test-cluster",
				Context: "test-context",
				Status:  ClusterStatusConnected,
			},
			expected: true,
		},
		{
			name: "disconnected cluster",
			cluster: Cluster{
				Name:    "test-cluster",
				Context: "test-context",
				Status:  ClusterStatusDisconnected,
			},
			expected: false,
		},
		{
			name: "error cluster",
			cluster: Cluster{
				Name:    "test-cluster",
				Context: "test-context",
				Status:  ClusterStatusError,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cluster.IsConnected()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCluster_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cluster     Cluster
		expectError bool
	}{
		{
			name: "valid cluster",
			cluster: Cluster{
				Name:    "test-cluster",
				Context: "test-context",
			},
			expectError: false,
		},
		{
			name: "missing name",
			cluster: Cluster{
				Context: "test-context",
			},
			expectError: true,
		},
		{
			name: "missing context",
			cluster: Cluster{
				Name: "test-cluster",
			},
			expectError: true,
		},
		{
			name:        "empty cluster",
			cluster:     Cluster{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cluster.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_IsReady(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected bool
	}{
		{
			name: "ready node",
			node: Node{
				Name:   "test-node",
				Status: NodeStatusReady,
			},
			expected: true,
		},
		{
			name: "not ready node",
			node: Node{
				Name:   "test-node",
				Status: NodeStatusNotReady,
			},
			expected: false,
		},
		{
			name: "unknown status node",
			node: Node{
				Name:   "test-node",
				Status: NodeStatusUnknown,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.IsReady()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNode_HasRole(t *testing.T) {
	node := Node{
		Name:  "test-node",
		Roles: []string{"master", "control-plane"},
	}

	assert.True(t, node.HasRole("master"))
	assert.True(t, node.HasRole("control-plane"))
	assert.False(t, node.HasRole("worker"))
	assert.False(t, node.HasRole(""))
}

func TestNode_IsMaster(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected bool
	}{
		{
			name: "master node",
			node: Node{
				Name:  "master-node",
				Roles: []string{"master"},
			},
			expected: true,
		},
		{
			name: "control-plane node",
			node: Node{
				Name:  "control-plane-node",
				Roles: []string{"control-plane"},
			},
			expected: true,
		},
		{
			name: "worker node",
			node: Node{
				Name:  "worker-node",
				Roles: []string{"worker"},
			},
			expected: false,
		},
		{
			name: "node with no roles",
			node: Node{
				Name:  "no-role-node",
				Roles: []string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.IsMaster()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeployment_IsHealthy(t *testing.T) {
	tests := []struct {
		name       string
		deployment Deployment
		expected   bool
	}{
		{
			name: "healthy deployment",
			deployment: Deployment{
				Name: "test-deployment",
				Replicas: DeploymentReplicas{
					Desired: 3,
					Ready:   3,
				},
			},
			expected: true,
		},
		{
			name: "degraded deployment",
			deployment: Deployment{
				Name: "test-deployment",
				Replicas: DeploymentReplicas{
					Desired: 3,
					Ready:   2,
				},
			},
			expected: false,
		},
		{
			name: "stopped deployment",
			deployment: Deployment{
				Name: "test-deployment",
				Replicas: DeploymentReplicas{
					Desired: 0,
					Ready:   0,
				},
			},
			expected: false,
		},
		{
			name: "failed deployment",
			deployment: Deployment{
				Name: "test-deployment",
				Replicas: DeploymentReplicas{
					Desired: 3,
					Ready:   0,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.deployment.IsHealthy()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeployment_GetAvailabilityPercentage(t *testing.T) {
	tests := []struct {
		name       string
		deployment Deployment
		expected   float64
	}{
		{
			name: "100% available",
			deployment: Deployment{
				Replicas: DeploymentReplicas{
					Desired: 3,
					Ready:   3,
				},
			},
			expected: 100.0,
		},
		{
			name: "66.67% available",
			deployment: Deployment{
				Replicas: DeploymentReplicas{
					Desired: 3,
					Ready:   2,
				},
			},
			expected: 66.66666666666666,
		},
		{
			name: "0% available",
			deployment: Deployment{
				Replicas: DeploymentReplicas{
					Desired: 3,
					Ready:   0,
				},
			},
			expected: 0.0,
		},
		{
			name: "stopped deployment",
			deployment: Deployment{
				Replicas: DeploymentReplicas{
					Desired: 0,
					Ready:   0,
				},
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.deployment.GetAvailabilityPercentage()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPod_IsRunning(t *testing.T) {
	tests := []struct {
		name     string
		pod      Pod
		expected bool
	}{
		{
			name: "running pod",
			pod: Pod{
				Name:   "test-pod",
				Status: PodStatusRunning,
			},
			expected: true,
		},
		{
			name: "pending pod",
			pod: Pod{
				Name:   "test-pod",
				Status: PodStatusPending,
			},
			expected: false,
		},
		{
			name: "failed pod",
			pod: Pod{
				Name:   "test-pod",
				Status: PodStatusFailed,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pod.IsRunning()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPod_IsReady(t *testing.T) {
	tests := []struct {
		name     string
		pod      Pod
		expected bool
	}{
		{
			name: "all containers ready",
			pod: Pod{
				Containers: []Container{
					{Name: "container1", Ready: true},
					{Name: "container2", Ready: true},
				},
			},
			expected: true,
		},
		{
			name: "some containers not ready",
			pod: Pod{
				Containers: []Container{
					{Name: "container1", Ready: true},
					{Name: "container2", Ready: false},
				},
			},
			expected: false,
		},
		{
			name: "no containers",
			pod: Pod{
				Containers: []Container{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pod.IsReady()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPod_GetTotalRestarts(t *testing.T) {
	pod := Pod{
		Containers: []Container{
			{Name: "container1", RestartCount: 2},
			{Name: "container2", RestartCount: 3},
			{Name: "container3", RestartCount: 0},
		},
	}

	result := pod.GetTotalRestarts()
	assert.Equal(t, int32(5), result)
}

func TestContainer_IsRunning(t *testing.T) {
	tests := []struct {
		name      string
		container Container
		expected  bool
	}{
		{
			name: "running container",
			container: Container{
				State: ContainerState{
					Running: &ContainerStateRunning{
						StartedAt: time.Now(),
					},
				},
			},
			expected: true,
		},
		{
			name: "waiting container",
			container: Container{
				State: ContainerState{
					Waiting: &ContainerStateWaiting{
						Reason: "ImagePullBackOff",
					},
				},
			},
			expected: false,
		},
		{
			name: "terminated container",
			container: Container{
				State: ContainerState{
					Terminated: &ContainerStateTerminated{
						ExitCode: 0,
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.container.IsRunning()
			assert.Equal(t, tt.expected, result)
		})
	}
}
