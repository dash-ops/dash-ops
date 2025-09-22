package logic

import (
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/stretchr/testify/assert"
)

func TestDashboardProcessor_ProcessDashboards_WithValidDashboards_ProcessesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewDashboardProcessor()
	dashboards := []models.Dashboard{
		{
			ID:          "dashboard-1",
			Name:        "Test Dashboard",
			Description: "A test dashboard",
			Service:     "test-service",
			Owner:       "test-user",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "dashboard-2",
			Name:        "Another Dashboard",
			Description: "Another test dashboard",
			Service:     "test-service",
			Owner:       "test-user",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
	}

	// Act
	result, err := processor.ProcessDashboards(dashboards)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "dashboard-1", result[0].ID)
	assert.Equal(t, "dashboard-2", result[1].ID)
}

func TestDashboardProcessor_ProcessDashboards_WithEmptyDashboards_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewDashboardProcessor()
	dashboards := []models.Dashboard{}

	// Act
	result, err := processor.ProcessDashboards(dashboards)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestDashboardProcessor_ProcessDashboards_WithNilDashboards_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewDashboardProcessor()

	// Act
	result, err := processor.ProcessDashboards(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "dashboards cannot be nil")
}

func TestDashboardProcessor_ValidateDashboard_WithValidDashboard_ValidatesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewDashboardProcessor()
	dashboard := &models.Dashboard{
		ID:          "dashboard-1",
		Name:        "Test Dashboard",
		Description: "A test dashboard",
		Service:     "test-service",
		Owner:       "test-user",
		Charts: []models.Chart{
			{
				ID:    "chart-1",
				Title: "CPU Usage",
				Type:  "line",
				Query: "cpu_usage",
			},
		},
	}

	// Act
	err := processor.ValidateDashboard(dashboard)

	// Assert
	assert.NoError(t, err)
}

func TestDashboardProcessor_ValidateDashboard_WithEmptyName_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewDashboardProcessor()
	dashboard := &models.Dashboard{
		ID:          "dashboard-1",
		Name:        "", // Empty name
		Description: "A test dashboard",
		Service:     "test-service",
		Owner:       "test-user",
	}

	// Act
	err := processor.ValidateDashboard(dashboard)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dashboard name is required")
}

func TestDashboardProcessor_ValidateDashboard_WithEmptyCharts_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewDashboardProcessor()
	dashboard := &models.Dashboard{
		ID:          "dashboard-1",
		Name:        "Test Dashboard",
		Description: "A test dashboard",
		Service:     "test-service",
		Owner:       "test-user",
		Charts:      []models.Chart{}, // Empty charts
	}

	// Act
	err := processor.ValidateDashboard(dashboard)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dashboard must have at least one chart")
}

func TestDashboardProcessor_EnrichDashboards_WithServiceContext_EnrichesSuccessfully(t *testing.T) {
	// Arrange
	processor := NewDashboardProcessor()
	dashboards := []models.Dashboard{
		{
			ID:      "dashboard-1",
			Name:    "Test Dashboard",
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
	result, err := processor.EnrichDashboards(dashboards, serviceContext)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "default", result[0].Namespace)
	assert.Equal(t, "test", result[0].Environment)
	assert.Equal(t, "platform", result[0].Team)
}
