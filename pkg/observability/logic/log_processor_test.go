package logic

import (
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/stretchr/testify/assert"
)

func TestLogProcessor_ProcessLogs_WithValidLogs_ProcessesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewLogProcessor()
	logs := []models.LogEntry{
		{
			ID:        "log-1",
			Timestamp: time.Now(),
			Level:     "info",
			Message:   "Test log message",
			Service:   "test-service",
			Labels:    map[string]string{"env": "test"},
		},
		{
			ID:        "log-2",
			Timestamp: time.Now().Add(-1 * time.Hour),
			Level:     "error",
			Message:   "Error log message",
			Service:   "test-service",
			Labels:    map[string]string{"env": "test"},
		},
	}

	// Act
	result, err := processor.ProcessLogs(logs)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "log-1", result[0].ID)
	assert.Equal(t, "log-2", result[1].ID)
}

func TestLogProcessor_ProcessLogs_WithEmptyLogs_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewLogProcessor()
	logs := []models.LogEntry{}

	// Act
	result, err := processor.ProcessLogs(logs)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestLogProcessor_ProcessLogs_WithNilLogs_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewLogProcessor()

	// Act
	result, err := processor.ProcessLogs(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "logs cannot be nil")
}

func TestLogProcessor_EnrichLogs_WithServiceContext_EnrichesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewLogProcessor()
	logs := []models.LogEntry{
		{
			ID:        "log-1",
			Timestamp: time.Now(),
			Level:     "info",
			Message:   "Test log message",
			Service:   "test-service",
		},
	}
	serviceContext := &models.ServiceContext{
		ServiceName: "test-service",
		Namespace:   "default",
		Environment: "test",
		Team:        "platform",
	}

	// Act
	result, err := processor.EnrichLogs(logs, serviceContext)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "default", result[0].Namespace)
	assert.Equal(t, "test", result[0].Environment)
	assert.Equal(t, "platform", result[0].Team)
}

func TestLogProcessor_EnrichLogs_WithNilServiceContext_ReturnsOriginalLogs(t *testing.T) {
	// Arrange
	processor := NewLogProcessor()
	logs := []models.LogEntry{
		{
			ID:        "log-1",
			Timestamp: time.Now(),
			Level:     "info",
			Message:   "Test log message",
			Service:   "test-service",
		},
	}

	// Act
	result, err := processor.EnrichLogs(logs, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "log-1", result[0].ID)
	assert.Empty(t, result[0].Namespace)
	assert.Empty(t, result[0].Environment)
	assert.Empty(t, result[0].Team)
}
