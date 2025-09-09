package servicecatalog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
)

func TestServiceValidator_ValidateForCreation_WithValidService_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	service := &scModels.Service{
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
	}

	// Act
	err := validator.ValidateForCreation(service)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_ValidateForCreation_WithNilService_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	service := (*scModels.Service)(nil)

	// Act
	err := validator.ValidateForCreation(service)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_ValidateForCreation_WithMissingName_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test service description",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}

	// Act
	err := validator.ValidateForCreation(service)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_ValidateForCreation_WithMissingDescription_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "test-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}

	// Act
	err := validator.ValidateForCreation(service)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_ValidateForCreation_WithInvalidTier_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	service := &scModels.Service{
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
	}

	// Act
	err := validator.ValidateForCreation(service)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_ValidateForCreation_WithMissingGitHubTeam_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "test-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test service description",
			Team:        scModels.ServiceTeam{},
		},
	}

	// Act
	err := validator.ValidateForCreation(service)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateServiceName_WithValidName_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := "test-service"

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateServiceName_WithValidNameWithUnderscores_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := "test_service"

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateServiceName_WithValidNameWithNumbers_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := "test-service-123"

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateServiceName_WithEmptyName_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := ""

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateServiceName_WithTooShortName_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := "ab"

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateServiceName_WithTooLongName_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := strings.Repeat("a", 101)

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateServiceName_WithInvalidCharacters_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := "test/service"

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateServiceName_WithNameStartingWithHyphen_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := "-test-service"

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateServiceName_WithNameEndingWithHyphen_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	serviceName := "test-service-"

	// Act
	err := validator.validateServiceName(serviceName)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateTier_WithValidTIER1_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	tier := scModels.TierCritical

	// Act
	err := validator.validateTier(tier)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateTier_WithValidTIER2_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	tier := scModels.TierImportant

	// Act
	err := validator.validateTier(tier)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateTier_WithValidTIER3_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	tier := scModels.TierStandard

	// Act
	err := validator.validateTier(tier)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateTier_WithInvalidTier_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	tier := scModels.ServiceTier("TIER-4")

	// Act
	err := validator.validateTier(tier)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateTier_WithEmptyTier_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	tier := scModels.ServiceTier("")

	// Act
	err := validator.validateTier(tier)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateCPUSpec_WithValidMillicore_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	cpu := "100m"

	// Act
	err := validator.validateCPUSpec(cpu)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateCPUSpec_WithValidCore_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	cpu := "1"

	// Act
	err := validator.validateCPUSpec(cpu)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateCPUSpec_WithValidDecimalCore_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	cpu := "0.5"

	// Act
	err := validator.validateCPUSpec(cpu)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateCPUSpec_WithInvalidFormat_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	cpu := "100cores"

	// Act
	err := validator.validateCPUSpec(cpu)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateCPUSpec_WithEmptyCPU_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	cpu := ""

	// Act
	err := validator.validateCPUSpec(cpu)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateMemorySpec_WithValidMi_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	memory := "128Mi"

	// Act
	err := validator.validateMemorySpec(memory)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateMemorySpec_WithValidGi_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	memory := "1Gi"

	// Act
	err := validator.validateMemorySpec(memory)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateMemorySpec_WithValidM_ReturnsNoError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	memory := "512M"

	// Act
	err := validator.validateMemorySpec(memory)

	// Assert
	assert.NoError(t, err)
}

func TestServiceValidator_validateMemorySpec_WithInvalidFormat_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	memory := "128MB"

	// Act
	err := validator.validateMemorySpec(memory)

	// Assert
	assert.Error(t, err)
}

func TestServiceValidator_validateMemorySpec_WithEmptyMemory_ReturnsError(t *testing.T) {
	// Arrange
	validator := NewServiceValidator()
	memory := ""

	// Act
	err := validator.validateMemorySpec(memory)

	// Assert
	assert.Error(t, err)
}
