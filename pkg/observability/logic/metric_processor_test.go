package logic

import (
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/stretchr/testify/assert"
)

func TestMetricProcessor_ProcessMetrics_WithValidMetrics_ProcessesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewMetricProcessor()
	metrics := []models.MetricData{
		{
			ID:        "metric-1",
			Timestamp: time.Now(),
			Name:      "cpu_usage",
			Value:     75.5,
			Labels:    map[string]string{"instance": "server-1"},
		},
		{
			ID:        "metric-2",
			Timestamp: time.Now().Add(-1 * time.Hour),
			Name:      "memory_usage",
			Value:     60.2,
			Labels:    map[string]string{"instance": "server-1"},
		},
	}

	// Act
	result, err := processor.ProcessMetrics(metrics)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "metric-1", result[0].ID)
	assert.Equal(t, "metric-2", result[1].ID)
}

func TestMetricProcessor_ProcessMetrics_WithEmptyMetrics_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewMetricProcessor()
	metrics := []models.MetricData{}

	// Act
	result, err := processor.ProcessMetrics(metrics)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestMetricProcessor_ProcessMetrics_WithNilMetrics_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewMetricProcessor()

	// Act
	result, err := processor.ProcessMetrics(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "metrics cannot be nil")
}

func TestMetricProcessor_CalculateDerivedMetrics_WithValidMetrics_CalculatesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewMetricProcessor()
	metrics := []models.MetricData{
		{Name: "cpu_usage", Value: 75.5, Timestamp: time.Now()},
		{Name: "cpu_usage", Value: 80.2, Timestamp: time.Now().Add(-1 * time.Minute)},
		{Name: "cpu_usage", Value: 70.1, Timestamp: time.Now().Add(-2 * time.Minute)},
	}

	// Act
	result, err := processor.CalculateDerivedMetrics(metrics)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "cpu_usage_avg", result[0].Name)
	assert.Equal(t, 75.27, result[0].Value, 0.01) // Average of 75.5, 80.2, 70.1
}

func TestMetricProcessor_CalculateDerivedMetrics_WithEmptyMetrics_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewMetricProcessor()
	metrics := []models.MetricData{}

	// Act
	result, err := processor.CalculateDerivedMetrics(metrics)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestMetricProcessor_EnrichMetrics_WithServiceContext_EnrichesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewMetricProcessor()
	metrics := []models.MetricData{
		{
			ID:        "metric-1",
			Timestamp: time.Now(),
			Name:      "cpu_usage",
			Value:     75.5,
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
	result, err := processor.EnrichMetrics(metrics, serviceContext)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "default", result[0].Namespace)
	assert.Equal(t, "test", result[0].Environment)
	assert.Equal(t, "platform", result[0].Team)
}
