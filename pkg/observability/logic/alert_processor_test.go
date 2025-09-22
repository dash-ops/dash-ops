package logic

import (
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/stretchr/testify/assert"
)

func TestAlertProcessor_ProcessAlerts_WithValidAlerts_ProcessesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewAlertProcessor()
	alerts := []models.Alert{
		{
			ID:          "alert-1",
			Name:        "High CPU Usage",
			Description: "CPU usage is above 80%",
			Severity:    "warning",
			Status:      "active",
			Service:     "test-service",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "alert-2",
			Name:        "Memory Usage Critical",
			Description: "Memory usage is above 90%",
			Severity:    "critical",
			Status:      "active",
			Service:     "test-service",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
	}

	// Act
	result, err := processor.ProcessAlerts(alerts)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "alert-1", result[0].ID)
	assert.Equal(t, "alert-2", result[1].ID)
}

func TestAlertProcessor_ProcessAlerts_WithEmptyAlerts_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewAlertProcessor()
	alerts := []models.Alert{}

	// Act
	result, err := processor.ProcessAlerts(alerts)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestAlertProcessor_ProcessAlerts_WithNilAlerts_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewAlertProcessor()

	// Act
	result, err := processor.ProcessAlerts(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "alerts cannot be nil")
}

func TestAlertProcessor_EvaluateAlert_WithValidAlert_EvaluatesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewAlertProcessor()
	alert := &models.Alert{
		ID:           "alert-1",
		Name:         "High CPU Usage",
		Description:  "CPU usage is above 80%",
		Severity:     "warning",
		Status:       "active",
		Service:      "test-service",
		Threshold:    80.0,
		CurrentValue: 85.0,
		CreatedAt:    time.Now(),
	}

	// Act
	result, err := processor.EvaluateAlert(alert)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "alert-1", result.AlertID)
	assert.True(t, result.Triggered)
	assert.Equal(t, 85.0, result.CurrentValue)
	assert.Equal(t, 80.0, result.Threshold)
}

func TestAlertProcessor_EvaluateAlert_WithBelowThreshold_DoesNotTrigger(t *testing.T) {
	// Arrange
	processor := NewAlertProcessor()
	alert := &models.Alert{
		ID:           "alert-1",
		Name:         "High CPU Usage",
		Description:  "CPU usage is above 80%",
		Severity:     "warning",
		Status:       "active",
		Service:      "test-service",
		Threshold:    80.0,
		CurrentValue: 75.0,
		CreatedAt:    time.Now(),
	}

	// Act
	result, err := processor.EvaluateAlert(alert)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "alert-1", result.AlertID)
	assert.False(t, result.Triggered)
	assert.Equal(t, 75.0, result.CurrentValue)
	assert.Equal(t, 80.0, result.Threshold)
}

func TestAlertProcessor_EnrichAlerts_WithServiceContext_EnrichesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewAlertProcessor()
	alerts := []models.Alert{
		{
			ID:      "alert-1",
			Name:    "High CPU Usage",
			Service: "test-service",
		},
	}
	serviceContext := &models.ServiceContext{
		ServiceName: "test-service",
		Namespace:   "default",
		Environment: "test",
		Team:        "platform",
	}

	// Act
	result, err := processor.EnrichAlerts(alerts, serviceContext)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "default", result[0].Namespace)
	assert.Equal(t, "test", result[0].Environment)
	assert.Equal(t, "platform", result[0].Team)
}
