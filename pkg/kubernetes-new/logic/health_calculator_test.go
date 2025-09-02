package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes-new/models"
)

func TestHealthCalculator_CalculateDeploymentHealth(t *testing.T) {
	calculator := NewHealthCalculator()

	tests := []struct {
		name       string
		deployment *k8sModels.Deployment
		expected   DeploymentHealthStatus
	}{
		{
			name:       "nil deployment",
			deployment: nil,
			expected:   DeploymentStatusUnknown,
		},
		{
			name: "healthy deployment",
			deployment: &k8sModels.Deployment{
				Replicas: k8sModels.DeploymentReplicas{
					Desired: 3,
					Ready:   3,
				},
			},
			expected: DeploymentStatusHealthy,
		},
		{
			name: "degraded deployment",
			deployment: &k8sModels.Deployment{
				Replicas: k8sModels.DeploymentReplicas{
					Desired: 3,
					Ready:   2,
				},
			},
			expected: DeploymentStatusDegraded,
		},
		{
			name: "unhealthy deployment",
			deployment: &k8sModels.Deployment{
				Replicas: k8sModels.DeploymentReplicas{
					Desired: 3,
					Ready:   0,
				},
			},
			expected: DeploymentStatusUnhealthy,
		},
		{
			name: "stopped deployment",
			deployment: &k8sModels.Deployment{
				Replicas: k8sModels.DeploymentReplicas{
					Desired: 0,
					Ready:   0,
				},
			},
			expected: DeploymentStatusStopped,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateDeploymentHealth(tt.deployment)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHealthCalculator_CalculatePodHealth(t *testing.T) {
	calculator := NewHealthCalculator()

	tests := []struct {
		name     string
		pod      *k8sModels.Pod
		expected PodHealthStatus
	}{
		{
			name:     "nil pod",
			pod:      nil,
			expected: PodHealthStatusUnknown,
		},
		{
			name: "healthy running pod",
			pod: &k8sModels.Pod{
				Status: k8sModels.PodStatusRunning,
				Containers: []k8sModels.Container{
					{Ready: true},
					{Ready: true},
				},
			},
			expected: PodHealthStatusHealthy,
		},
		{
			name: "degraded running pod",
			pod: &k8sModels.Pod{
				Status: k8sModels.PodStatusRunning,
				Containers: []k8sModels.Container{
					{Ready: true},
					{Ready: false},
				},
			},
			expected: PodHealthStatusDegraded,
		},
		{
			name: "pending pod",
			pod: &k8sModels.Pod{
				Status: k8sModels.PodStatusPending,
			},
			expected: PodHealthStatusPending,
		},
		{
			name: "failed pod",
			pod: &k8sModels.Pod{
				Status: k8sModels.PodStatusFailed,
			},
			expected: PodHealthStatusFailed,
		},
		{
			name: "succeeded pod",
			pod: &k8sModels.Pod{
				Status: k8sModels.PodStatusSucceeded,
			},
			expected: PodHealthStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculatePodHealth(tt.pod)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHealthCalculator_CalculateNodeHealth(t *testing.T) {
	calculator := NewHealthCalculator()

	tests := []struct {
		name     string
		node     *k8sModels.Node
		expected NodeHealthStatus
	}{
		{
			name:     "nil node",
			node:     nil,
			expected: NodeHealthStatusUnknown,
		},
		{
			name: "ready node",
			node: &k8sModels.Node{
				Status: k8sModels.NodeStatusReady,
			},
			expected: NodeHealthStatusHealthy,
		},
		{
			name: "not ready node",
			node: &k8sModels.Node{
				Status: k8sModels.NodeStatusNotReady,
			},
			expected: NodeHealthStatusUnhealthy,
		},
		{
			name: "unknown status node",
			node: &k8sModels.Node{
				Status: k8sModels.NodeStatusUnknown,
			},
			expected: NodeHealthStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateNodeHealth(tt.node)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHealthCalculator_CalculateClusterHealth(t *testing.T) {
	calculator := NewHealthCalculator()

	tests := []struct {
		name        string
		clusterInfo *k8sModels.ClusterInfo
		expected    ClusterHealthStatus
	}{
		{
			name:        "nil cluster info",
			clusterInfo: nil,
			expected:    ClusterHealthStatusDisconnected,
		},
		{
			name: "disconnected cluster",
			clusterInfo: &k8sModels.ClusterInfo{
				Cluster: k8sModels.Cluster{
					Status: k8sModels.ClusterStatusDisconnected,
				},
			},
			expected: ClusterHealthStatusDisconnected,
		},
		{
			name: "healthy cluster (all nodes ready)",
			clusterInfo: &k8sModels.ClusterInfo{
				Cluster: k8sModels.Cluster{
					Status: k8sModels.ClusterStatusConnected,
				},
				Nodes: []k8sModels.Node{
					{Status: k8sModels.NodeStatusReady},
					{Status: k8sModels.NodeStatusReady},
					{Status: k8sModels.NodeStatusReady},
				},
			},
			expected: ClusterHealthStatusHealthy,
		},
		{
			name: "degraded cluster (80% nodes ready)",
			clusterInfo: &k8sModels.ClusterInfo{
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
			},
			expected: ClusterHealthStatusDegraded,
		},
		{
			name: "critical cluster (20% nodes ready)",
			clusterInfo: &k8sModels.ClusterInfo{
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
			},
			expected: ClusterHealthStatusCritical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateClusterHealth(tt.clusterInfo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHealthCalculator_GetDeploymentHealthSummary(t *testing.T) {
	calculator := NewHealthCalculator()

	tests := []struct {
		name           string
		deployment     *k8sModels.Deployment
		expectedStatus DeploymentHealthStatus
		expectedIssues int
	}{
		{
			name:           "nil deployment",
			deployment:     nil,
			expectedStatus: DeploymentStatusUnknown,
			expectedIssues: 1,
		},
		{
			name: "healthy deployment",
			deployment: &k8sModels.Deployment{
				Replicas: k8sModels.DeploymentReplicas{
					Desired: 3,
					Ready:   3,
				},
			},
			expectedStatus: DeploymentStatusHealthy,
			expectedIssues: 0,
		},
		{
			name: "scaled to zero",
			deployment: &k8sModels.Deployment{
				Replicas: k8sModels.DeploymentReplicas{
					Desired: 0,
					Ready:   0,
				},
			},
			expectedStatus: DeploymentStatusStopped,
			expectedIssues: 1,
		},
		{
			name: "no replicas ready",
			deployment: &k8sModels.Deployment{
				Replicas: k8sModels.DeploymentReplicas{
					Desired: 3,
					Ready:   0,
				},
			},
			expectedStatus: DeploymentStatusUnhealthy,
			expectedIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := calculator.GetDeploymentHealthSummary(tt.deployment)
			assert.Equal(t, tt.expectedStatus, summary.Status)
			assert.Equal(t, tt.expectedIssues, len(summary.Issues))
		})
	}
}
