package logic

import (
	"context"
	"errors"
	"testing"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
)

// MockClusterRepository is a mock implementation of ClusterRepository
type MockClusterRepository struct {
	clusters []k8sModels.Cluster
	errors   map[string]error
}

func (m *MockClusterRepository) GetCluster(ctx context.Context, context string) (*k8sModels.Cluster, error) {
	if err, exists := m.errors["GetCluster"]; exists {
		return nil, err
	}

	for _, cluster := range m.clusters {
		if cluster.Context == context {
			return &cluster, nil
		}
	}
	return nil, errors.New("cluster not found")
}

func (m *MockClusterRepository) ListClusters(ctx context.Context) ([]k8sModels.Cluster, error) {
	if err, exists := m.errors["ListClusters"]; exists {
		return nil, err
	}
	return m.clusters, nil
}

func (m *MockClusterRepository) ValidateCluster(ctx context.Context, context string) error {
	if err, exists := m.errors["ValidateCluster"]; exists {
		return err
	}
	return nil
}

func (m *MockClusterRepository) GetClusterInfo(ctx context.Context, context string) (*k8sModels.ClusterInfo, error) {
	return nil, errors.New("not implemented")
}

func TestClusterProcessor_ListClusters(t *testing.T) {
	tests := []struct {
		name           string
		clusters       []k8sModels.Cluster
		validateErrors map[string]error
		expectedStatus k8sModels.ClusterStatus
		expectError    bool
	}{
		{
			name: "successful list with connected clusters",
			clusters: []k8sModels.Cluster{
				{
					Name:    "test-cluster-1",
					Context: "test-context-1",
					Status:  k8sModels.ClusterStatusUnknown,
				},
				{
					Name:    "test-cluster-2",
					Context: "test-context-2",
					Status:  k8sModels.ClusterStatusUnknown,
				},
			},
			expectedStatus: k8sModels.ClusterStatusConnected,
			expectError:    false,
		},
		{
			name: "list with validation errors",
			clusters: []k8sModels.Cluster{
				{
					Name:    "test-cluster-1",
					Context: "test-context-1",
					Status:  k8sModels.ClusterStatusUnknown,
				},
			},
			validateErrors: map[string]error{
				"ValidateCluster": errors.New("connection failed"),
			},
			expectedStatus: k8sModels.ClusterStatusError,
			expectError:    false,
		},
		{
			name:     "repository error",
			clusters: []k8sModels.Cluster{},
			validateErrors: map[string]error{
				"ListClusters": errors.New("repository error"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockClusterRepository{
				clusters: tt.clusters,
				errors:   tt.validateErrors,
			}

			processor := NewClusterProcessor(mockRepo)

			result, err := processor.ListClusters(context.Background())

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.clusters) {
				t.Errorf("expected %d clusters, got %d", len(tt.clusters), len(result))
			}

			for _, cluster := range result {
				if cluster.Status != tt.expectedStatus {
					t.Errorf("expected status %s, got %s", tt.expectedStatus, cluster.Status)
				}
			}
		})
	}
}

func TestClusterProcessor_GetClusterStatus(t *testing.T) {
	tests := []struct {
		name           string
		context        string
		validateErrors map[string]error
		expectedStatus k8sModels.ClusterStatus
		expectError    bool
	}{
		{
			name:           "connected cluster",
			context:        "test-context",
			expectedStatus: k8sModels.ClusterStatusConnected,
			expectError:    false,
		},
		{
			name:    "error cluster",
			context: "test-context",
			validateErrors: map[string]error{
				"ValidateCluster": errors.New("connection failed"),
			},
			expectedStatus: k8sModels.ClusterStatusError,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockClusterRepository{
				errors: tt.validateErrors,
			}

			processor := NewClusterProcessor(mockRepo)

			status, err := processor.GetClusterStatus(context.Background(), tt.context)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, status)
			}
		})
	}
}

func TestClusterProcessor_ValidateClusterConfig(t *testing.T) {
	tests := []struct {
		name        string
		cluster     *k8sModels.Cluster
		expectError bool
	}{
		{
			name: "valid cluster",
			cluster: &k8sModels.Cluster{
				Name:    "test-cluster",
				Context: "test-context",
			},
			expectError: false,
		},
		{
			name: "missing name",
			cluster: &k8sModels.Cluster{
				Context: "test-context",
			},
			expectError: true,
		},
		{
			name: "missing context",
			cluster: &k8sModels.Cluster{
				Name: "test-cluster",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockClusterRepository{}
			processor := NewClusterProcessor(mockRepo)

			err := processor.ValidateClusterConfig(tt.cluster)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
