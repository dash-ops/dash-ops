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

func TestHTTPHandler_CreateService(t *testing.T) {
	// Setup test dependencies
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

	tests := []struct {
		name           string
		requestBody    scWire.CreateServiceRequest
		expectedStatus int
		setupMocks     func(*MockServiceRepository)
	}{
		{
			name: "successful creation",
			requestBody: scWire.CreateServiceRequest{
				Name:        "test-service",
				Description: "Test service description",
				Tier:        "TIER-3",
				Team: scWire.TeamRequest{
					GitHubTeam: "test-team",
				},
			},
			expectedStatus: http.StatusCreated,
			setupMocks: func(repo *MockServiceRepository) {
				repo.On("Exists", mock.Anything, "test-service").Return(false, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(&scModels.Service{
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
			},
		},
		{
			name: "missing required fields",
			requestBody: scWire.CreateServiceRequest{
				Name: "test-service",
				// Missing description and team
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(repo *MockServiceRepository) {
				// No mocks needed for validation errors
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks(mockRepo)

			// Create request
			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/services", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			handler.createServiceHandler(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response scWire.ServiceResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.requestBody.Name, response.Metadata.Name)
			}

			// Clear mock expectations
			mockRepo.ExpectedCalls = nil
		})
	}
}

func TestHTTPHandler_GetService(t *testing.T) {
	// Setup test dependencies
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

	tests := []struct {
		name           string
		serviceName    string
		expectedStatus int
		setupMocks     func(*MockServiceRepository)
	}{
		{
			name:           "existing service",
			serviceName:    "test-service",
			expectedStatus: http.StatusOK,
			setupMocks: func(repo *MockServiceRepository) {
				repo.On("GetByName", mock.Anything, "test-service").Return(&scModels.Service{
					Metadata: scModels.ServiceMetadata{
						Name: "test-service",
						Tier: scModels.TierStandard,
					},
					Spec: scModels.ServiceSpec{
						Description: "Test service",
						Team: scModels.ServiceTeam{
							GitHubTeam: "test-team",
						},
					},
				}, nil)
			},
		},
		{
			name:           "non-existing service",
			serviceName:    "non-existent",
			expectedStatus: http.StatusNotFound,
			setupMocks: func(repo *MockServiceRepository) {
				repo.On("GetByName", mock.Anything, "non-existent").Return(nil, fmt.Errorf("service not found"))
			},
		},
		{
			name:           "empty service name",
			serviceName:    "",
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(repo *MockServiceRepository) {
				// No mocks needed for validation errors
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks(mockRepo)

			// Create request
			req := httptest.NewRequest("GET", "/services/"+tt.serviceName, nil)
			req = mux.SetURLVars(req, map[string]string{"name": tt.serviceName})
			w := httptest.NewRecorder()

			// Execute
			handler.getServiceHandler(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response scWire.ServiceResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.serviceName, response.Metadata.Name)
			}

			// Clear mock expectations
			mockRepo.ExpectedCalls = nil
		})
	}
}
