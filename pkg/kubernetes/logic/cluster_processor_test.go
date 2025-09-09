package logic

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestClusterProcessor_ListClusters_WithSuccessfulList_ReturnsConnectedClusters(t *testing.T) {
	// Arrange
	clusters := []k8sModels.Cluster{
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
	}
	mockRepo := &MockClusterRepository{
		clusters: clusters,
		errors:   map[string]error{},
	}
	processor := NewClusterProcessor(mockRepo)

	// Act
	result, err := processor.ListClusters(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, cluster := range result {
		assert.Equal(t, k8sModels.ClusterStatusConnected, cluster.Status)
	}
}

func TestClusterProcessor_ListClusters_WithValidationErrors_ReturnsErrorClusters(t *testing.T) {
	// Arrange
	clusters := []k8sModels.Cluster{
		{
			Name:    "test-cluster-1",
			Context: "test-context-1",
			Status:  k8sModels.ClusterStatusUnknown,
		},
	}
	validateErrors := map[string]error{
		"ValidateCluster": errors.New("connection failed"),
	}
	mockRepo := &MockClusterRepository{
		clusters: clusters,
		errors:   validateErrors,
	}
	processor := NewClusterProcessor(mockRepo)

	// Act
	result, err := processor.ListClusters(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, k8sModels.ClusterStatusError, result[0].Status)
}

func TestClusterProcessor_ListClusters_WithRepositoryError_ReturnsError(t *testing.T) {
	// Arrange
	validateErrors := map[string]error{
		"ListClusters": errors.New("repository error"),
	}
	mockRepo := &MockClusterRepository{
		clusters: []k8sModels.Cluster{},
		errors:   validateErrors,
	}
	processor := NewClusterProcessor(mockRepo)

	// Act
	result, err := processor.ListClusters(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestClusterProcessor_GetClusterStatus_WithConnectedCluster_ReturnsConnected(t *testing.T) {
	// Arrange
	contextName := "test-context"
	mockRepo := &MockClusterRepository{
		errors: map[string]error{},
	}
	processor := NewClusterProcessor(mockRepo)

	// Act
	status, err := processor.GetClusterStatus(context.Background(), contextName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, k8sModels.ClusterStatusConnected, status)
}

func TestClusterProcessor_GetClusterStatus_WithErrorCluster_ReturnsError(t *testing.T) {
	// Arrange
	contextName := "test-context"
	validateErrors := map[string]error{
		"ValidateCluster": errors.New("connection failed"),
	}
	mockRepo := &MockClusterRepository{
		errors: validateErrors,
	}
	processor := NewClusterProcessor(mockRepo)

	// Act
	status, err := processor.GetClusterStatus(context.Background(), contextName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, k8sModels.ClusterStatusError, status)
}

func TestClusterProcessor_ValidateClusterConfig_WithValidCluster_ReturnsNoError(t *testing.T) {
	// Arrange
	cluster := &k8sModels.Cluster{
		Name:    "test-cluster",
		Context: "test-context",
	}
	mockRepo := &MockClusterRepository{}
	processor := NewClusterProcessor(mockRepo)

	// Act
	err := processor.ValidateClusterConfig(cluster)

	// Assert
	assert.NoError(t, err)
}

func TestClusterProcessor_ValidateClusterConfig_WithMissingName_ReturnsError(t *testing.T) {
	// Arrange
	cluster := &k8sModels.Cluster{
		Context: "test-context",
	}
	mockRepo := &MockClusterRepository{}
	processor := NewClusterProcessor(mockRepo)

	// Act
	err := processor.ValidateClusterConfig(cluster)

	// Assert
	assert.Error(t, err)
}

func TestClusterProcessor_ValidateClusterConfig_WithMissingContext_ReturnsError(t *testing.T) {
	// Arrange
	cluster := &k8sModels.Cluster{
		Name: "test-cluster",
	}
	mockRepo := &MockClusterRepository{}
	processor := NewClusterProcessor(mockRepo)

	// Act
	err := processor.ValidateClusterConfig(cluster)

	// Assert
	assert.Error(t, err)
}
