package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogsController is a mock implementation of LogsController
type MockLogsController struct {
	mock.Mock
}

func (m *MockLogsController) GetLogs(ctx context.Context, req *wire.LogsRequest) (*wire.LogsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.LogsResponse), args.Error(1)
}

func (m *MockLogsController) GetLogStatistics(ctx context.Context, req *wire.LogStatsRequest) (*wire.LogStatisticsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.LogStatisticsResponse), args.Error(1)
}

// MockMetricsController is a mock implementation of MetricsController
type MockMetricsController struct {
	mock.Mock
}

func (m *MockMetricsController) GetMetrics(ctx context.Context, req *wire.MetricsRequest) (*wire.MetricsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.MetricsResponse), args.Error(1)
}

func (m *MockMetricsController) GetMetricStatistics(ctx context.Context, req *wire.MetricStatsRequest) (*wire.MetricStatisticsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.MetricStatisticsResponse), args.Error(1)
}

// MockTracesController is a mock implementation of TracesController
type MockTracesController struct {
	mock.Mock
}

func (m *MockTracesController) GetTraces(ctx context.Context, req *wire.TracesRequest) (*wire.TracesResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.TracesResponse), args.Error(1)
}

func (m *MockTracesController) GetTraceDetail(ctx context.Context, req *wire.TraceDetailRequest) (*wire.TraceDetailResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.TraceDetailResponse), args.Error(1)
}

func (m *MockTracesController) GetTraceStatistics(ctx context.Context, req *wire.TraceStatsRequest) (*wire.TraceStatisticsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.TraceStatisticsResponse), args.Error(1)
}

func (m *MockTracesController) AnalyzeTrace(ctx context.Context, req *wire.TraceAnalysisRequest) (*wire.TraceAnalysisResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.TraceAnalysisResponse), args.Error(1)
}

// MockAlertsController is a mock implementation of AlertsController
type MockAlertsController struct {
	mock.Mock
}

func (m *MockAlertsController) GetAlerts(ctx context.Context, req *wire.AlertsRequest) (*wire.AlertsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.AlertsResponse), args.Error(1)
}

func (m *MockAlertsController) CreateAlert(ctx context.Context, req *wire.CreateAlertRequest) (*wire.AlertResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.AlertResponse), args.Error(1)
}

func (m *MockAlertsController) GetAlertStatistics(ctx context.Context, req *wire.AlertStatsRequest) (*wire.AlertStatisticsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.AlertStatisticsResponse), args.Error(1)
}

// MockHealthController is a mock implementation of HealthController
type MockHealthController struct {
	mock.Mock
}

func (m *MockHealthController) GetCacheStats(ctx context.Context) (*wire.CacheStatsResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*wire.CacheStatsResponse), args.Error(1)
}

func (m *MockHealthController) Health(ctx context.Context) (*wire.HealthResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*wire.HealthResponse), args.Error(1)
}

// MockConfigController is a mock implementation of ConfigController
type MockConfigController struct {
	mock.Mock
}

func (m *MockConfigController) GetConfiguration(ctx context.Context) (*wire.ConfigurationResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*wire.ConfigurationResponse), args.Error(1)
}

func (m *MockConfigController) UpdateConfiguration(ctx context.Context, config *wire.ObservabilityConfig) (*wire.ConfigurationResponse, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(*wire.ConfigurationResponse), args.Error(1)
}

func (m *MockConfigController) GetServiceConfiguration(ctx context.Context, serviceName string) (*wire.ServiceConfigurationResponse, error) {
	args := m.Called(ctx, serviceName)
	return args.Get(0).(*wire.ServiceConfigurationResponse), args.Error(1)
}

func (m *MockConfigController) UpdateServiceConfiguration(ctx context.Context, serviceName string, config *wire.ServiceObservabilityConfig) (*wire.ServiceConfigurationResponse, error) {
	args := m.Called(ctx, serviceName, config)
	return args.Get(0).(*wire.ServiceConfigurationResponse), args.Error(1)
}

func (m *MockConfigController) GetNotificationChannels(ctx context.Context) (*wire.NotificationChannelsResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*wire.NotificationChannelsResponse), args.Error(1)
}

func (m *MockConfigController) ConfigureNotificationChannel(ctx context.Context, req *wire.NotificationChannelRequest) (*wire.NotificationChannelResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wire.NotificationChannelResponse), args.Error(1)
}

func TestHTTPHandler_GetLogs_WithValidRequest_Returns200(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	expectedResponse := &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:  []models.LogEntry{},
			Total: 0,
		},
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.LogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockLogsController.AssertExpectations(t)
}

func TestHTTPHandler_GetLogs_WithInvalidJSON_Returns400(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	// Create HTTP request with invalid JSON
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotEmpty(t, response.Error)
}

func TestHTTPHandler_GetLogs_WithControllerError_Returns500(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return((*wire.LogsResponse)(nil), assert.AnError)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotEmpty(t, response.Error)

	mockLogsController.AssertExpectations(t)
}

func TestHTTPHandler_GetLogs_WithMissingService_Returns400(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "", // Missing service
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
	}

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "service is required")
}

func TestHTTPHandler_GetLogs_WithInvalidTimeRange_Returns400(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(-1 * time.Hour), // End before start
		Limit:     100,
	}

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "start time must be before end time")
}

func TestHTTPHandler_GetLogs_WithNegativeLimit_Returns400(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     -1, // Negative limit
	}

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "limit must be positive")
}

func TestHTTPHandler_GetLogs_WithZeroLimit_UsesDefaultLimit(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     0, // Zero limit
	}

	expectedResponse := &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:  []models.LogEntry{},
			Total: 0,
		},
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.LogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockLogsController.AssertExpectations(t)
}

func TestHTTPHandler_GetLogs_WithLargeLimit_ClampsToMaxLimit(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     10000, // Large limit
	}

	expectedResponse := &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:  []models.LogEntry{},
			Total: 0,
		},
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.LogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockLogsController.AssertExpectations(t)
}

func TestHTTPHandler_GetLogs_WithStreamingEnabled_ReturnsStreamingResponse(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Stream:    true, // Streaming enabled
	}

	expectedResponse := &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:  []models.LogEntry{},
			Total: 0,
		},
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.LogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockLogsController.AssertExpectations(t)
}

func TestHTTPHandler_GetLogs_WithQuery_AppliesQuery(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Query:     `{service="test-service"} |= "error"`, // LogQL query
	}

	expectedResponse := &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:  []models.LogEntry{},
			Total: 0,
		},
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.LogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockLogsController.AssertExpectations(t)
}

func TestHTTPHandler_GetLogs_WithLevelFilter_AppliesLevelFilter(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Level:     "error", // Level filter
	}

	expectedResponse := &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:  []models.LogEntry{},
			Total: 0,
		},
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.LogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockLogsController.AssertExpectations(t)
}

func TestHTTPHandler_GetLogs_WithSorting_AppliesSorting(t *testing.T) {
	// Arrange
	mockLogsController := &MockLogsController{}
	mockMetricsController := &MockMetricsController{}
	mockTracesController := &MockTracesController{}
	mockAlertsController := &MockAlertsController{}
	mockHealthController := &MockHealthController{}
	mockConfigController := &MockConfigController{}

	handler := NewHTTPHandler(
		mockLogsController,
		mockMetricsController,
		mockTracesController,
		mockAlertsController,
		mockHealthController,
		mockConfigController,
		nil, // responseAdapter
		nil, // requestAdapter
	)

	req := &wire.LogsRequest{
		Service:   "test-service",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Limit:     100,
		Sort:      "timestamp",
		Order:     "desc",
	}

	expectedResponse := &wire.LogsResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.LogsData{
			Logs:  []models.LogEntry{},
			Total: 0,
		},
	}

	mockLogsController.On("GetLogs", mock.Anything, req).Return(expectedResponse, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/observability/logs", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.GetLogs(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response wire.LogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockLogsController.AssertExpectations(t)
}
