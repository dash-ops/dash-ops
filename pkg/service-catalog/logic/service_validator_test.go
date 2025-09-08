package servicecatalog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
)

func TestServiceValidator_ValidateForCreation(t *testing.T) {
	validator := NewServiceValidator()

	tests := []struct {
		name        string
		service     *scModels.Service
		expectError bool
	}{
		{
			name: "valid service",
			service: &scModels.Service{
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
			},
			expectError: false,
		},
		{
			name:        "nil service",
			service:     nil,
			expectError: true,
		},
		{
			name: "missing name",
			service: &scModels.Service{
				Metadata: scModels.ServiceMetadata{
					Tier: scModels.TierStandard,
				},
				Spec: scModels.ServiceSpec{
					Description: "Test service description",
					Team: scModels.ServiceTeam{
						GitHubTeam: "test-team",
					},
				},
			},
			expectError: true,
		},
		{
			name: "missing description",
			service: &scModels.Service{
				Metadata: scModels.ServiceMetadata{
					Name: "test-service",
					Tier: scModels.TierStandard,
				},
				Spec: scModels.ServiceSpec{
					Team: scModels.ServiceTeam{
						GitHubTeam: "test-team",
					},
				},
			},
			expectError: true,
		},
		{
			name: "invalid tier",
			service: &scModels.Service{
				Metadata: scModels.ServiceMetadata{
					Name: "test-service",
					Tier: "INVALID-TIER",
				},
				Spec: scModels.ServiceSpec{
					Description: "Test service description",
					Team: scModels.ServiceTeam{
						GitHubTeam: "test-team",
					},
				},
			},
			expectError: true,
		},
		{
			name: "missing github team",
			service: &scModels.Service{
				Metadata: scModels.ServiceMetadata{
					Name: "test-service",
					Tier: scModels.TierStandard,
				},
				Spec: scModels.ServiceSpec{
					Description: "Test service description",
					Team:        scModels.ServiceTeam{},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateForCreation(tt.service)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceValidator_validateServiceName(t *testing.T) {
	validator := NewServiceValidator()

	tests := []struct {
		name        string
		serviceName string
		expectError bool
	}{
		{
			name:        "valid name",
			serviceName: "test-service",
			expectError: false,
		},
		{
			name:        "valid name with underscores",
			serviceName: "test_service",
			expectError: false,
		},
		{
			name:        "valid name with numbers",
			serviceName: "test-service-123",
			expectError: false,
		},
		{
			name:        "empty name",
			serviceName: "",
			expectError: true,
		},
		{
			name:        "too short",
			serviceName: "ab",
			expectError: true,
		},
		{
			name:        "too long",
			serviceName: strings.Repeat("a", 101),
			expectError: true,
		},
		{
			name:        "invalid characters",
			serviceName: "test/service",
			expectError: true,
		},
		{
			name:        "starts with hyphen",
			serviceName: "-test-service",
			expectError: true,
		},
		{
			name:        "ends with hyphen",
			serviceName: "test-service-",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateServiceName(tt.serviceName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceValidator_validateTier(t *testing.T) {
	validator := NewServiceValidator()

	tests := []struct {
		name        string
		tier        scModels.ServiceTier
		expectError bool
	}{
		{
			name:        "valid TIER-1",
			tier:        scModels.TierCritical,
			expectError: false,
		},
		{
			name:        "valid TIER-2",
			tier:        scModels.TierImportant,
			expectError: false,
		},
		{
			name:        "valid TIER-3",
			tier:        scModels.TierStandard,
			expectError: false,
		},
		{
			name:        "invalid tier",
			tier:        "TIER-4",
			expectError: true,
		},
		{
			name:        "empty tier",
			tier:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateTier(tt.tier)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceValidator_validateCPUSpec(t *testing.T) {
	validator := NewServiceValidator()

	tests := []struct {
		name        string
		cpu         string
		expectError bool
	}{
		{
			name:        "valid millicore",
			cpu:         "100m",
			expectError: false,
		},
		{
			name:        "valid core",
			cpu:         "1",
			expectError: false,
		},
		{
			name:        "valid decimal core",
			cpu:         "0.5",
			expectError: false,
		},
		{
			name:        "invalid format",
			cpu:         "100cores",
			expectError: true,
		},
		{
			name:        "empty cpu",
			cpu:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateCPUSpec(tt.cpu)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceValidator_validateMemorySpec(t *testing.T) {
	validator := NewServiceValidator()

	tests := []struct {
		name        string
		memory      string
		expectError bool
	}{
		{
			name:        "valid Mi",
			memory:      "128Mi",
			expectError: false,
		},
		{
			name:        "valid Gi",
			memory:      "1Gi",
			expectError: false,
		},
		{
			name:        "valid M",
			memory:      "512M",
			expectError: false,
		},
		{
			name:        "invalid format",
			memory:      "128MB",
			expectError: true,
		},
		{
			name:        "empty memory",
			memory:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateMemorySpec(tt.memory)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
