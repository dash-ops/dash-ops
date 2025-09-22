package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/logic"
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
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

func TestLogsController_GetLogs_WithValidRequest_ReturnsLogs(t *testing.T) {
	// Arrange
	mockLogRepo := &MockLogRepository{}
	mockServiceRepo := &MockServiceContextRepository{}
	mockLogService := &MockLogService{}
	mockCacheService := &MockCacheService{}
	processor := logic.NewLogProcessor()

	controller := NewLogsController(
		mockLogRepo,
		mockServiceRepo,
		mockLogService,
		mockCacheService,
		processor,
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	expectedLogs := []models.LogEntry{
		{
			ID:        "log-1",
			Timestamp: time.Now(),
			Level:     "info",
			Message:   "Test log message",
			Service:   "test-service",
		},
	}

	mockLogRepo.On("GetLogs", mock.Anything, req).Return(expectedLogs, nil)
	mockCacheService.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, assert.AnError)
	mockCacheService.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// Act
	result, err := controller.GetLogs(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Len(t, result.Data.Logs, 1)
	assert.Equal(t, "log-1", result.Data.Logs[0].ID)

	mockLogRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

func TestLogsController_GetLogs_WithRepositoryError_ReturnsError(t *testing.T) {
	// Arrange
	mockLogRepo := &MockLogRepository{}
	mockServiceRepo := &MockServiceContextRepository{}
	mockLogService := &MockLogService{}
	mockCacheService := &MockCacheService{}
	processor := logic.NewLogProcessor()

	controller := NewLogsController(
		mockLogRepo,
		mockServiceRepo,
		mockLogService,
		mockCacheService,
		processor,
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	mockLogRepo.On("GetLogs", mock.Anything, req).Return([]models.LogEntry{}, assert.AnError)
	mockCacheService.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	// Act
	result, err := controller.GetLogs(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockLogRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

func TestLogsController_GetLogStatistics_WithValidRequest_ReturnsStatistics(t *testing.T) {
	// Arrange
	mockLogRepo := &MockLogRepository{}
	mockServiceRepo := &MockServiceContextRepository{}
	mockLogService := &MockLogService{}
	mockCacheService := &MockCacheService{}
	processor := logic.NewLogProcessor()

	controller := NewLogsController(
		mockLogRepo,
		mockServiceRepo,
		mockLogService,
		mockCacheService,
		processor,
	)

	req := &wire.LogStatsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
	}

	expectedStats := &wire.LogStatistics{
		TotalLogs:     100,
		LogsByLevel:   map[string]int64{"info": 80, "error": 20},
		LogsByService: map[string]int64{"test-service": 100},
		LogsByHost:    map[string]int64{"host-1": 60, "host-2": 40},
	}

	mockLogRepo.On("GetLogStatistics", mock.Anything, req).Return(expectedStats, nil)
	mockCacheService.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, assert.AnError)
	mockCacheService.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// Act
	result, err := controller.GetLogStatistics(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, int64(100), result.Data.TotalLogs)
	assert.Equal(t, int64(80), result.Data.LogsByLevel["info"])
	assert.Equal(t, int64(20), result.Data.LogsByLevel["error"])

	mockLogRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

func TestLogsController_GetLogStatistics_WithRepositoryError_ReturnsError(t *testing.T) {
	// Arrange
	mockLogRepo := &MockLogRepository{}
	mockServiceRepo := &MockServiceContextRepository{}
	mockLogService := &MockLogService{}
	mockCacheService := &MockCacheService{}
	processor := logic.NewLogProcessor()

	controller := NewLogsController(
		mockLogRepo,
		mockServiceRepo,
		mockLogService,
		mockCacheService,
		processor,
	)

	req := &wire.LogStatsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
	}

	mockLogRepo.On("GetLogStatistics", mock.Anything, req).Return((*wire.LogStatistics)(nil), assert.AnError)
	mockCacheService.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	// Act
	result, err := controller.GetLogStatistics(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockLogRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}
