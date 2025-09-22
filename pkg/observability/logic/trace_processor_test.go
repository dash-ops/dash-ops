package logic

import (
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/stretchr/testify/assert"
)

func TestTraceProcessor_ProcessTraces_WithValidTraces_ProcessesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewTraceProcessor()
	traces := []models.TraceInfo{
		{
			TraceID:   "trace-1",
			StartTime: time.Now(),
			Duration:  100 * time.Millisecond,
			Service:   "test-service",
			Operation: "test-operation",
			Status:    "ok",
		},
		{
			TraceID:   "trace-2",
			StartTime: time.Now().Add(-1 * time.Hour),
			Duration:  200 * time.Millisecond,
			Service:   "test-service",
			Operation: "test-operation",
			Status:    "error",
		},
	}

	// Act
	result, err := processor.ProcessTraces(traces)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "trace-1", result[0].TraceID)
	assert.Equal(t, "trace-2", result[1].TraceID)
}

func TestTraceProcessor_ProcessTraces_WithEmptyTraces_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewTraceProcessor()
	traces := []models.TraceInfo{}

	// Act
	result, err := processor.ProcessTraces(traces)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestTraceProcessor_ProcessTraces_WithNilTraces_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewTraceProcessor()

	// Act
	result, err := processor.ProcessTraces(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "traces cannot be nil")
}

func TestTraceProcessor_AnalyzeTrace_WithValidTrace_AnalyzesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewTraceProcessor()
	trace := &models.TraceInfo{
		TraceID:   "trace-1",
		StartTime: time.Now(),
		Duration:  100 * time.Millisecond,
		Service:   "test-service",
		Operation: "test-operation",
		Status:    "ok",
	}

	// Act
	result, err := processor.AnalyzeTrace(trace)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "trace-1", result.TraceID)
	assert.Equal(t, int64(100), result.TotalDuration)
	assert.Equal(t, 0.0, result.ErrorRate)
}

func TestTraceProcessor_AnalyzeTrace_WithErrorTrace_CalculatesErrorRate(t *testing.T) {
	// Arrange
	processor := NewTraceProcessor()
	trace := &models.TraceInfo{
		TraceID:   "trace-1",
		StartTime: time.Now(),
		Duration:  100 * time.Millisecond,
		Service:   "test-service",
		Operation: "test-operation",
		Status:    "error",
	}

	// Act
	result, err := processor.AnalyzeTrace(trace)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "trace-1", result.TraceID)
	assert.Equal(t, int64(100), result.TotalDuration)
	assert.Equal(t, 1.0, result.ErrorRate)
}

func TestTraceProcessor_EnrichTraces_WithServiceContext_EnrichesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewTraceProcessor()
	traces := []models.TraceInfo{
		{
			TraceID:   "trace-1",
			StartTime: time.Now(),
			Duration:  100 * time.Millisecond,
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
	result, err := processor.EnrichTraces(traces, serviceContext)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "default", result[0].Namespace)
	assert.Equal(t, "test", result[0].Environment)
	assert.Equal(t, "platform", result[0].Team)
}
