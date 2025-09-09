package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	scAdapters "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters/http"
	servicecatalog "github.com/dash-ops/dash-ops/pkg/service-catalog/controllers"
	scLogic "github.com/dash-ops/dash-ops/pkg/service-catalog/logic"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scWire "github.com/dash-ops/dash-ops/pkg/service-catalog/wire"
)

// MockServiceRepository is a mock implementation of ServiceRepository
type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) Create(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
	args := m.Called(ctx, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*scModels.Service), args.Error(1)
}

func (m *MockServiceRepository) GetByName(ctx context.Context, name string) (*scModels.Service, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*scModels.Service), args.Error(1)
}

func (m *MockServiceRepository) Update(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
	args := m.Called(ctx, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*scModels.Service), args.Error(1)
}

func (m *MockServiceRepository) Delete(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockServiceRepository) List(ctx context.Context, filter *scModels.ServiceFilter) ([]scModels.Service, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]scModels.Service), args.Error(1)
}

func (m *MockServiceRepository) Exists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockServiceRepository) ListByTeam(ctx context.Context, team string) ([]scModels.Service, error) {
	args := m.Called(ctx, team)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]scModels.Service), args.Error(1)
}

func (m *MockServiceRepository) ListByTier(ctx context.Context, tier scModels.ServiceTier) ([]scModels.Service, error) {
	args := m.Called(ctx, tier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]scModels.Service), args.Error(1)
}

func (m *MockServiceRepository) Search(ctx context.Context, query string, limit int) ([]scModels.Service, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]scModels.Service), args.Error(1)
}

func TestHTTPHandler_CreateService_SuccessfulCreation_ReturnsCreated(t *testing.T) {
	// Arrange
	mockRepo := &MockServiceRepository{}
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	controller := servicecatalog.NewServiceController(
		mockRepo, nil, nil, nil, validator, processor,
	)

	serviceAdapter := scAdapters.NewServiceAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	handler := NewHTTPHandler(controller, serviceAdapter, responseAdapter, requestAdapter)

	requestBody := scWire.CreateServiceRequest{
		Name:        "test-service",
		Description: "Test service description",
		Tier:        "TIER-3",
		Team: scWire.TeamRequest{
			GitHubTeam: "test-team",
		},
	}

	mockRepo.On("Exists", mock.Anything, "test-service").Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(&scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "test-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test service description",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}, nil)

	requestBodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/services", bytes.NewReader(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.createServiceHandler(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response scWire.ServiceResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, requestBody.Name, response.Metadata.Name)
}

func TestHTTPHandler_CreateService_MissingRequiredFields_ReturnsBadRequest(t *testing.T) {
	// Arrange
	mockRepo := &MockServiceRepository{}
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	controller := servicecatalog.NewServiceController(
		mockRepo, nil, nil, nil, validator, processor,
	)

	serviceAdapter := scAdapters.NewServiceAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	handler := NewHTTPHandler(controller, serviceAdapter, responseAdapter, requestAdapter)

	requestBody := scWire.CreateServiceRequest{
		Name: "test-service",
		// Missing description and team
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/services", bytes.NewReader(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.createServiceHandler(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPHandler_GetService_ExistingService_ReturnsOK(t *testing.T) {
	// Arrange
	mockRepo := &MockServiceRepository{}
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	controller := servicecatalog.NewServiceController(
		mockRepo, nil, nil, nil, validator, processor,
	)

	serviceAdapter := scAdapters.NewServiceAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	handler := NewHTTPHandler(controller, serviceAdapter, responseAdapter, requestAdapter)

	serviceName := "test-service"
	mockRepo.On("GetByName", mock.Anything, serviceName).Return(&scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: serviceName,
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test service",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}, nil)

	req := httptest.NewRequest("GET", "/services/"+serviceName, nil)
	req = mux.SetURLVars(req, map[string]string{"name": serviceName})
	w := httptest.NewRecorder()

	// Act
	handler.getServiceHandler(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response scWire.ServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, serviceName, response.Metadata.Name)
}

func TestHTTPHandler_GetService_NonExistingService_ReturnsNotFound(t *testing.T) {
	// Arrange
	mockRepo := &MockServiceRepository{}
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	controller := servicecatalog.NewServiceController(
		mockRepo, nil, nil, nil, validator, processor,
	)

	serviceAdapter := scAdapters.NewServiceAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	handler := NewHTTPHandler(controller, serviceAdapter, responseAdapter, requestAdapter)

	serviceName := "non-existent"
	mockRepo.On("GetByName", mock.Anything, serviceName).Return(nil, fmt.Errorf("service not found"))

	req := httptest.NewRequest("GET", "/services/"+serviceName, nil)
	req = mux.SetURLVars(req, map[string]string{"name": serviceName})
	w := httptest.NewRecorder()

	// Act
	handler.getServiceHandler(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHTTPHandler_GetService_EmptyServiceName_ReturnsBadRequest(t *testing.T) {
	// Arrange
	mockRepo := &MockServiceRepository{}
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	controller := servicecatalog.NewServiceController(
		mockRepo, nil, nil, nil, validator, processor,
	)

	serviceAdapter := scAdapters.NewServiceAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	handler := NewHTTPHandler(controller, serviceAdapter, responseAdapter, requestAdapter)

	serviceName := ""
	req := httptest.NewRequest("GET", "/services/"+serviceName, nil)
	req = mux.SetURLVars(req, map[string]string{"name": serviceName})
	w := httptest.NewRecorder()

	// Act
	handler.getServiceHandler(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
