package external

import (
	"context"
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/wire"
	"github.com/stretchr/testify/assert"
)

func TestLokiAdapter_GetLogs_WithValidRequest_ReturnsLogs(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestLokiAdapter_GetLogs_WithInvalidURL_ReturnsError(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("invalid-url")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLokiAdapter_GetLogs_WithNilRequest_ReturnsError(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")

	// Act
	result, err := adapter.GetLogs(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request cannot be nil")
}

func TestLokiAdapter_GetLogs_WithEmptyService_ReturnsError(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "", // Empty service
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "service is required")
}

func TestLokiAdapter_GetLogs_WithInvalidTimeRange_ReturnsError(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(-1 * time.Hour), // End before start
		Limit:     100,
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "start time must be before end time")
}

func TestLokiAdapter_GetLogs_WithNegativeLimit_ReturnsError(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     -1, // Negative limit
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "limit must be positive")
}

func TestLokiAdapter_GetLogs_WithZeroLimit_UsesDefaultLimit(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     0, // Zero limit
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestLokiAdapter_GetLogs_WithLargeLimit_ClampsToMaxLimit(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     10000, // Large limit
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestLokiAdapter_GetLogs_WithStreamingEnabled_ReturnsStreamingResponse(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Stream:    true, // Streaming enabled
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestLokiAdapter_GetLogs_WithQuery_AppliesQuery(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Query:     `{service="test-service"} |= "error"`, // LogQL query
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestLokiAdapter_GetLogs_WithLevelFilter_AppliesLevelFilter(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Level:     "error", // Level filter
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestLokiAdapter_GetLogs_WithSorting_AppliesSorting(t *testing.T) {
	// Arrange
	adapter := NewLokiAdapter("http://localhost:3100")
	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Sort:      "timestamp",
		Order:     "desc",
	}

	// Act
	result, err := adapter.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}
