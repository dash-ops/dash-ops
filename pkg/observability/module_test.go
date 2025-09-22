package observability

import (
	"context"
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogRepository is a mock implementation of LogRepository
type MockLogRepository struct {
	mock.Mock
}

func (m *MockLogRepository) GetLogs(ctx context.Context, req *wire.LogsRequest) ([]models.LogEntry, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.LogEntry), args.Error(1)
}

func (m *MockLogRepository) GetLogStatistics(ctx context.Context, req *wire.LogStatsRequest) (*wire.LogStatistics, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.LogStatistics), args.Error(1)
}

// MockMetricRepository is a mock implementation of MetricRepository
type MockMetricRepository struct {
	mock.Mock
}

func (m *MockMetricRepository) GetMetrics(ctx context.Context, req *wire.MetricsRequest) ([]models.MetricData, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.MetricData), args.Error(1)
}

func (m *MockMetricRepository) GetMetricStatistics(ctx context.Context, req *wire.MetricStatsRequest) (*wire.MetricStatistics, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.MetricStatistics), args.Error(1)
}

// MockTraceRepository is a mock implementation of TraceRepository
type MockTraceRepository struct {
	mock.Mock
}

func (m *MockTraceRepository) GetTraces(ctx context.Context, req *wire.TracesRequest) ([]models.TraceInfo, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.TraceInfo), args.Error(1)
}

func (m *MockTraceRepository) GetTraceDetail(ctx context.Context, req *wire.TraceDetailRequest) ([]models.TraceSpan, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.TraceSpan), args.Error(1)
}

func (m *MockTraceRepository) GetTraceStatistics(ctx context.Context, req *wire.TraceStatsRequest) (*wire.TraceStatistics, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.TraceStatistics), args.Error(1)
}

// MockAlertRepository is a mock implementation of AlertRepository
type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) GetAlerts(ctx context.Context, req *wire.AlertsRequest) ([]models.Alert, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.Alert), args.Error(1)
}

func (m *MockAlertRepository) CreateAlert(ctx context.Context, req *wire.CreateAlertRequest) (*models.Alert, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Alert), args.Error(1)
}

func (m *MockAlertRepository) UpdateAlert(ctx context.Context, req *wire.UpdateAlertRequest) (*models.Alert, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Alert), args.Error(1)
}

func (m *MockAlertRepository) DeleteAlert(ctx context.Context, req *wire.DeleteAlertRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAlertRepository) GetAlertStatistics(ctx context.Context, req *wire.AlertStatsRequest) (*wire.AlertStatistics, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.AlertStatistics), args.Error(1)
}

// MockDashboardRepository is a mock implementation of DashboardRepository
type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetDashboards(ctx context.Context, req *wire.DashboardsRequest) ([]models.Dashboard, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.Dashboard), args.Error(1)
}

func (m *MockDashboardRepository) CreateDashboard(ctx context.Context, req *wire.CreateDashboardRequest) (*models.Dashboard, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Dashboard), args.Error(1)
}

func (m *MockDashboardRepository) UpdateDashboard(ctx context.Context, req *wire.UpdateDashboardRequest) (*models.Dashboard, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Dashboard), args.Error(1)
}

func (m *MockDashboardRepository) DeleteDashboard(ctx context.Context, req *wire.DeleteDashboardRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockDashboardRepository) GetDashboard(ctx context.Context, req *wire.GetDashboardRequest) (*models.Dashboard, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Dashboard), args.Error(1)
}

func (m *MockDashboardRepository) GetDashboardTemplates(ctx context.Context, req *wire.DashboardTemplatesRequest) ([]models.DashboardTemplate, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.DashboardTemplate), args.Error(1)
}

func (m *MockDashboardRepository) GetDashboardStatistics(ctx context.Context, req *wire.DashboardStatsRequest) (*wire.DashboardStatistics, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.DashboardStatistics), args.Error(1)
}

// MockServiceContextRepository is a mock implementation of ServiceContextRepository
type MockServiceContextRepository struct {
	mock.Mock
}

func (m *MockServiceContextRepository) GetServiceContext(ctx context.Context, serviceName string) (*models.ServiceContext, error) {
	args := m.Called(ctx, serviceName)
	return args.Get(0).(*models.ServiceContext), args.Error(1)
}

func (m *MockServiceContextRepository) GetServicesWithContext(ctx context.Context, req *wire.ServicesWithContextRequest) ([]models.ServiceWithContext, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.ServiceWithContext), args.Error(1)
}

func (m *MockServiceContextRepository) GetServiceHealth(ctx context.Context, req *wire.ServiceHealthRequest) (*models.ServiceHealth, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.ServiceHealth), args.Error(1)
}

// MockLogService is a mock implementation of LogService
type MockLogService struct {
	mock.Mock
}

func (m *MockLogService) ProcessLogs(ctx context.Context, logs []models.LogEntry) ([]models.ProcessedLogEntry, error) {
	args := m.Called(ctx, logs)
	return args.Get(0).([]models.ProcessedLogEntry), args.Error(1)
}

func (m *MockLogService) StreamLogs(ctx context.Context, req *wire.LogsRequest) (<-chan models.LogEntry, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(<-chan models.LogEntry), args.Error(1)
}

// MockMetricService is a mock implementation of MetricService
type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) ProcessMetrics(ctx context.Context, metrics []models.MetricData) ([]models.ProcessedMetric, error) {
	args := m.Called(ctx, metrics)
	return args.Get(0).([]models.ProcessedMetric), args.Error(1)
}

func (m *MockMetricService) CalculateDerivedMetrics(ctx context.Context, metrics []models.MetricData) ([]models.DerivedMetric, error) {
	args := m.Called(ctx, metrics)
	return args.Get(0).([]models.DerivedMetric), args.Error(1)
}

func (m *MockMetricService) StreamMetrics(ctx context.Context, req *wire.MetricsRequest) (<-chan models.MetricData, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(<-chan models.MetricData), args.Error(1)
}

// MockTraceService is a mock implementation of TraceService
type MockTraceService struct {
	mock.Mock
}

func (m *MockTraceService) ProcessTraces(ctx context.Context, traces []models.TraceInfo) ([]models.ProcessedTrace, error) {
	args := m.Called(ctx, traces)
	return args.Get(0).([]models.ProcessedTrace), args.Error(1)
}

func (m *MockTraceService) AnalyzeTrace(ctx context.Context, trace *models.TraceInfo) (*models.TracePerformance, error) {
	args := m.Called(ctx, trace)
	return args.Get(0).(*models.TracePerformance), args.Error(1)
}

func (m *MockTraceService) StreamTraces(ctx context.Context, req *wire.TracesRequest) (<-chan models.TraceInfo, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(<-chan models.TraceInfo), args.Error(1)
}

// MockAlertService is a mock implementation of AlertService
type MockAlertService struct {
	mock.Mock
}

func (m *MockAlertService) ProcessAlerts(ctx context.Context, alerts []models.Alert) ([]models.ProcessedAlert, error) {
	args := m.Called(ctx, alerts)
	return args.Get(0).([]models.ProcessedAlert), args.Error(1)
}

func (m *MockAlertService) EvaluateAlert(ctx context.Context, alert *models.Alert) (*models.AlertEvaluation, error) {
	args := m.Called(ctx, alert)
	return args.Get(0).(*models.AlertEvaluation), args.Error(1)
}

func (m *MockAlertService) StreamAlerts(ctx context.Context, req *wire.AlertsRequest) (<-chan models.Alert, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(<-chan models.Alert), args.Error(1)
}

// MockDashboardService is a mock implementation of DashboardService
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) ProcessDashboards(ctx context.Context, dashboards []models.Dashboard) ([]models.Dashboard, error) {
	args := m.Called(ctx, dashboards)
	return args.Get(0).([]models.Dashboard), args.Error(1)
}

func (m *MockDashboardService) ValidateDashboard(ctx context.Context, dashboard *models.Dashboard) error {
	args := m.Called(ctx, dashboard)
	return args.Error(0)
}

func (m *MockDashboardService) GetDashboardData(ctx context.Context, dashboard *models.Dashboard) (*models.DashboardData, error) {
	args := m.Called(ctx, dashboard)
	return args.Get(0).(*models.DashboardData), args.Error(1)
}

// MockNotificationService is a mock implementation of NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(ctx context.Context, channel string, message string) error {
	args := m.Called(ctx, channel, message)
	return args.Error(0)
}

func (m *MockNotificationService) GetNotificationChannels(ctx context.Context) ([]models.NotificationChannel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.NotificationChannel), args.Error(1)
}

func (m *MockNotificationService) ConfigureNotificationChannel(ctx context.Context, channel *models.NotificationChannel) error {
	args := m.Called(ctx, channel)
	return args.Error(0)
}

// MockCacheService is a mock implementation of CacheService
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Error(1)
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) GetStats(ctx context.Context) (*models.CacheStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.CacheStats), args.Error(1)
}

// MockConfigurationService is a mock implementation of ConfigurationService
type MockConfigurationService struct {
	mock.Mock
}

func (m *MockConfigurationService) GetConfiguration(ctx context.Context) (*models.ObservabilityConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.ObservabilityConfig), args.Error(1)
}

func (m *MockConfigurationService) UpdateConfiguration(ctx context.Context, config *models.ObservabilityConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockConfigurationService) GetServiceConfiguration(ctx context.Context, serviceName string) (*models.ServiceObservabilityConfig, error) {
	args := m.Called(ctx, serviceName)
	return args.Get(0).(*models.ServiceObservabilityConfig), args.Error(1)
}

func (m *MockConfigurationService) UpdateServiceConfiguration(ctx context.Context, serviceName string, config *models.ServiceObservabilityConfig) error {
	args := m.Called(ctx, serviceName, config)
	return args.Error(0)
}

func TestNewModule_WithValidConfig_CreatesModuleSuccessfully(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	// Act
	module, err := NewModule(config)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.NotNil(t, module.LogsController)
	assert.NotNil(t, module.MetricsController)
	assert.NotNil(t, module.TracesController)
	assert.NotNil(t, module.AlertsController)
	assert.NotNil(t, module.HealthController)
	assert.NotNil(t, module.ConfigController)
	assert.NotNil(t, module.Handler)
	assert.NotNil(t, module.LogProcessor)
	assert.NotNil(t, module.MetricProcessor)
	assert.NotNil(t, module.TraceProcessor)
	assert.NotNil(t, module.AlertProcessor)
	assert.NotNil(t, module.DashboardProcessor)
	assert.NotNil(t, module.ResponseAdapter)
	assert.NotNil(t, module.RequestAdapter)
}

func TestNewModule_WithNilConfig_ReturnsError(t *testing.T) {
	// Arrange
	var config *ModuleConfig = nil

	// Act
	module, err := NewModule(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "module config cannot be nil")
}

func TestNewModule_WithNilLogRepo_ReturnsError(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              nil, // Missing required dependency
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	// Act
	module, err := NewModule(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "log repository is required")
}

func TestNewModule_WithNilMetricRepo_ReturnsError(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           nil, // Missing required dependency
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	// Act
	module, err := NewModule(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "metric repository is required")
}

func TestNewModule_WithNilTraceRepo_ReturnsError(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            nil, // Missing required dependency
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	// Act
	module, err := NewModule(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "trace repository is required")
}

func TestNewModule_WithNilAlertRepo_ReturnsError(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            nil, // Missing required dependency
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	// Act
	module, err := NewModule(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "alert repository is required")
}

func TestNewModule_WithNilDashboardRepo_ReturnsError(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        nil, // Missing required dependency
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	// Act
	module, err := NewModule(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "dashboard repository is required")
}

func TestNewModule_WithNilServiceRepo_ReturnsError(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          nil, // Missing required dependency
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	// Act
	module, err := NewModule(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "service context repository is required")
}

func TestModule_WithValidConfig_HasHandler(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	module, _ := NewModule(config)

	// Assert
	assert.NotNil(t, module.Handler)
}

func TestModule_WithValidConfig_HasAllControllers(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	module, _ := NewModule(config)

	// Assert
	assert.NotNil(t, module.LogsController)
	assert.NotNil(t, module.MetricsController)
	assert.NotNil(t, module.TracesController)
	assert.NotNil(t, module.AlertsController)
	assert.NotNil(t, module.HealthController)
	assert.NotNil(t, module.ConfigController)
}

func TestModule_RegisterRoutes_WithValidHandler_RegistersRoutes(t *testing.T) {
	// Arrange
	config := &ModuleConfig{
		LogRepo:              &MockLogRepository{},
		MetricRepo:           &MockMetricRepository{},
		TraceRepo:            &MockTraceRepository{},
		AlertRepo:            &MockAlertRepository{},
		DashboardRepo:        &MockDashboardRepository{},
		ServiceRepo:          &MockServiceContextRepository{},
		LogService:           &MockLogService{},
		MetricService:        &MockMetricService{},
		TraceService:         &MockTraceService{},
		AlertService:         &MockAlertService{},
		DashboardService:     &MockDashboardService{},
		NotificationService:  &MockNotificationService{},
		CacheService:         &MockCacheService{},
		ConfigurationService: &MockConfigurationService{},
	}

	module, _ := NewModule(config)
	router := mux.NewRouter()

	// Act
	module.RegisterRoutes(router)

	// Assert
	assert.NotNil(t, router)
	// Note: In a real test, you would verify that specific routes are registered
	// This would require more complex setup and route verification
}

func TestModule_RegisterRoutes_WithNilHandler_DoesNotPanic(t *testing.T) {
	// Arrange
	module := &Module{
		Handler: nil, // Nil handler
	}
	router := mux.NewRouter()

	// Act & Assert
	assert.NotPanics(t, func() {
		module.RegisterRoutes(router)
	})
}
