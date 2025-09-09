package servicecatalog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
)

func TestServiceProcessor_PrepareForCreation_WithNilService_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	service := (*scModels.Service)(nil)
	user := &scModels.UserContext{Username: "testuser"}

	// Act
	result, err := processor.PrepareForCreation(service, user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestServiceProcessor_PrepareForCreation_WithValidServiceAndUser_ReturnsPreparedService(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	service := &scModels.Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: scModels.ServiceMetadata{
			Name: "Test Service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test service description",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}
	user := &scModels.UserContext{Username: "testuser"}

	// Act
	result, err := processor.PrepareForCreation(service, user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-service", result.Metadata.Name)     // Normalized name
	assert.Equal(t, "test-team", result.Spec.Team.GitHubTeam) // Normalized team
	assert.Equal(t, user.Username, result.Metadata.CreatedBy)
	assert.Equal(t, user.Username, result.Metadata.UpdatedBy)
}

func TestServiceProcessor_PrepareForCreation_WithValidServiceWithoutUser_ReturnsPreparedService(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	service := &scModels.Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: scModels.ServiceMetadata{
			Name: "Test Service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test service description",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}
	user := (*scModels.UserContext)(nil)

	// Act
	result, err := processor.PrepareForCreation(service, user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-service", result.Metadata.Name)     // Normalized name
	assert.Equal(t, "test-team", result.Spec.Team.GitHubTeam) // Normalized team
}

func TestServiceProcessor_PrepareForUpdate_WithNilService_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	service := (*scModels.Service)(nil)
	existingService := &scModels.Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: scModels.ServiceMetadata{
			Name:      "existing-service",
			CreatedAt: time.Now().Add(-time.Hour),
			CreatedBy: "original-user",
			Version:   1,
		},
		Spec: scModels.ServiceSpec{
			Description: "Original description",
		},
	}
	user := &scModels.UserContext{Username: "testuser"}

	// Act
	result, err := processor.PrepareForUpdate(service, existingService, user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestServiceProcessor_PrepareForUpdate_WithNilExistingService_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	service := &scModels.Service{}
	existingService := (*scModels.Service)(nil)
	user := &scModels.UserContext{Username: "testuser"}

	// Act
	result, err := processor.PrepareForUpdate(service, existingService, user)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestServiceProcessor_PrepareForUpdate_WithValidUpdate_ReturnsUpdatedService(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	service := &scModels.Service{
		Spec: scModels.ServiceSpec{
			Description: "Updated description",
			Team: scModels.ServiceTeam{
				GitHubTeam: "updated-team",
			},
		},
	}
	existingService := &scModels.Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: scModels.ServiceMetadata{
			Name:      "existing-service",
			CreatedAt: time.Now().Add(-time.Hour),
			CreatedBy: "original-user",
			Version:   1,
		},
		Spec: scModels.ServiceSpec{
			Description: "Original description",
		},
	}
	user := &scModels.UserContext{Username: "testuser"}

	// Act
	result, err := processor.PrepareForUpdate(service, existingService, user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Check that immutable fields are preserved
	assert.Equal(t, existingService.APIVersion, result.APIVersion)
	assert.Equal(t, existingService.Kind, result.Kind)
	assert.Equal(t, existingService.Metadata.Name, result.Metadata.Name)
	assert.Equal(t, existingService.Metadata.CreatedAt, result.Metadata.CreatedAt)
	assert.Equal(t, existingService.Metadata.CreatedBy, result.Metadata.CreatedBy)
	// Check that version is incremented
	assert.Equal(t, existingService.Metadata.Version+1, result.Metadata.Version)
	// Check that updated fields are set
	assert.Equal(t, user.Username, result.Metadata.UpdatedBy)
}

func TestServiceProcessor_CalculateServiceHealth_WithNoEnvironments_ReturnsUnknown(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	environments := []scModels.EnvironmentHealth{}
	tier := scModels.TierStandard

	// Act
	result := processor.CalculateServiceHealth(environments, tier)

	// Assert
	assert.Equal(t, scModels.StatusUnknown, result)
}

func TestServiceProcessor_CalculateServiceHealth_WithProductionHealthyAndCriticalTier_ReturnsHealthy(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	environments := []scModels.EnvironmentHealth{
		{Name: "production", Status: scModels.StatusHealthy},
		{Name: "staging", Status: scModels.StatusDegraded},
	}
	tier := scModels.TierCritical

	// Act
	result := processor.CalculateServiceHealth(environments, tier)

	// Assert
	assert.Equal(t, scModels.StatusHealthy, result)
}

func TestServiceProcessor_CalculateServiceHealth_WithProductionDegradedAndCriticalTier_ReturnsCritical(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	environments := []scModels.EnvironmentHealth{
		{Name: "production", Status: scModels.StatusDegraded},
		{Name: "staging", Status: scModels.StatusHealthy},
	}
	tier := scModels.TierCritical

	// Act
	result := processor.CalculateServiceHealth(environments, tier)

	// Assert
	assert.Equal(t, scModels.StatusCritical, result)
}

func TestServiceProcessor_CalculateServiceHealth_WithProductionDownAndCriticalTier_ReturnsCritical(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	environments := []scModels.EnvironmentHealth{
		{Name: "production", Status: scModels.StatusDown},
		{Name: "staging", Status: scModels.StatusHealthy},
	}
	tier := scModels.TierCritical

	// Act
	result := processor.CalculateServiceHealth(environments, tier)

	// Assert
	assert.Equal(t, scModels.StatusCritical, result)
}

func TestServiceProcessor_CalculateServiceHealth_WithProductionDownAndImportantTier_ReturnsDegraded(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	environments := []scModels.EnvironmentHealth{
		{Name: "production", Status: scModels.StatusDown},
		{Name: "staging", Status: scModels.StatusHealthy},
	}
	tier := scModels.TierImportant

	// Act
	result := processor.CalculateServiceHealth(environments, tier)

	// Assert
	assert.Equal(t, scModels.StatusDegraded, result)
}

func TestServiceProcessor_CalculateServiceHealth_WithProductionDownAndStandardTier_ReturnsDegraded(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	environments := []scModels.EnvironmentHealth{
		{Name: "production", Status: scModels.StatusDown},
		{Name: "staging", Status: scModels.StatusHealthy},
	}
	tier := scModels.TierStandard

	// Act
	result := processor.CalculateServiceHealth(environments, tier)

	// Assert
	assert.Equal(t, scModels.StatusDegraded, result)
}

func TestServiceProcessor_CalculateServiceHealth_WithNoProductionEnvironment_ReturnsFirstEnvironmentStatus(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	environments := []scModels.EnvironmentHealth{
		{Name: "staging", Status: scModels.StatusHealthy},
		{Name: "development", Status: scModels.StatusDegraded},
	}
	tier := scModels.TierStandard

	// Act
	result := processor.CalculateServiceHealth(environments, tier)

	// Assert
	assert.Equal(t, scModels.StatusHealthy, result)
}

func TestServiceProcessor_ProcessServiceList_WithNoFilter_ReturnsAllServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service1",
				Tier: scModels.TierCritical,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service2",
				Tier: scModels.TierStandard,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team2"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "another-service",
				Tier: scModels.TierImportant,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
	}
	filter := (*scModels.ServiceFilter)(nil)

	// Act
	result := processor.ProcessServiceList(services, filter)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 3)
	assert.Equal(t, 3, result.Total)
}

func TestServiceProcessor_ProcessServiceList_WithTeamFilter_ReturnsFilteredServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service1",
				Tier: scModels.TierCritical,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service2",
				Tier: scModels.TierStandard,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team2"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "another-service",
				Tier: scModels.TierImportant,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
	}
	filter := &scModels.ServiceFilter{
		Team: "team1",
	}

	// Act
	result := processor.ProcessServiceList(services, filter)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 2)
	assert.Equal(t, 2, result.Total)
}

func TestServiceProcessor_ProcessServiceList_WithTierFilter_ReturnsFilteredServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service1",
				Tier: scModels.TierCritical,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service2",
				Tier: scModels.TierStandard,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team2"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "another-service",
				Tier: scModels.TierImportant,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
	}
	filter := &scModels.ServiceFilter{
		Tier: "TIER-1",
	}

	// Act
	result := processor.ProcessServiceList(services, filter)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 1)
	assert.Equal(t, 1, result.Total)
}

func TestServiceProcessor_ProcessServiceList_WithSearchFilter_ReturnsFilteredServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service1",
				Tier: scModels.TierCritical,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service2",
				Tier: scModels.TierStandard,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team2"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "another-service",
				Tier: scModels.TierImportant,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
	}
	filter := &scModels.ServiceFilter{
		Search: "another",
	}

	// Act
	result := processor.ProcessServiceList(services, filter)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 1)
	assert.Equal(t, 1, result.Total)
}

func TestServiceProcessor_ProcessServiceList_WithPagination_ReturnsPaginatedServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service1",
				Tier: scModels.TierCritical,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "service2",
				Tier: scModels.TierStandard,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team2"},
			},
		},
		{
			Metadata: scModels.ServiceMetadata{
				Name: "another-service",
				Tier: scModels.TierImportant,
			},
			Spec: scModels.ServiceSpec{
				Team: scModels.ServiceTeam{GitHubTeam: "team1"},
			},
		},
	}
	filter := &scModels.ServiceFilter{
		Limit:  2,
		Offset: 1,
	}

	// Act
	result := processor.ProcessServiceList(services, filter)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 2)
	assert.Equal(t, 3, result.Total)
}

func TestServiceProcessor_normalizeServiceName_WithSimpleName_ReturnsSameName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "test-service"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "test-service", result)
}

func TestServiceProcessor_normalizeServiceName_WithSpaces_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "Test Service"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "test-service", result)
}

func TestServiceProcessor_normalizeServiceName_WithUnderscores_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "test_service"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "test-service", result)
}

func TestServiceProcessor_normalizeServiceName_WithMixedCase_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "Test_Service Name"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "test-service-name", result)
}

func TestServiceProcessor_normalizeServiceName_WithConsecutiveHyphens_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "test--service"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "test-service", result)
}

func TestServiceProcessor_normalizeServiceName_WithLeadingTrailingHyphens_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "-test-service-"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "test-service", result)
}

func TestServiceProcessor_normalizeServiceName_WithSpecialCharacters_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "test@service#name"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "test-service-name", result)
}

func TestServiceProcessor_normalizeServiceName_WithEmptyString_ReturnsEmptyString(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := ""

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "", result)
}

func TestServiceProcessor_normalizeServiceName_WithOnlySpecialChars_ReturnsEmptyString(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "---"

	// Act
	result := processor.normalizeServiceName(input)

	// Assert
	assert.Equal(t, "", result)
}

func TestServiceProcessor_GenerateServiceID_WithSimpleName_ReturnsSameName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "test-service"

	// Act
	result := processor.GenerateServiceID(input)

	// Assert
	assert.Equal(t, "test-service", result)
}

func TestServiceProcessor_GenerateServiceID_WithSpaces_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "Test Service"

	// Act
	result := processor.GenerateServiceID(input)

	// Assert
	assert.Equal(t, "test-service", result)
}

func TestServiceProcessor_GenerateServiceID_WithSpecialChars_ReturnsNormalizedName(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	input := "Test@Service#Name"

	// Act
	result := processor.GenerateServiceID(input)

	// Assert
	assert.Equal(t, "test-service-name", result)
}

func TestServiceProcessor_CompareServices_WithDifferentServices_ReturnsChanges(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	oldService := &scModels.Service{
		Spec: scModels.ServiceSpec{
			Description: "Old description",
			Team: scModels.ServiceTeam{
				GitHubTeam: "old-team",
			},
		},
		Metadata: scModels.ServiceMetadata{
			Tier: scModels.TierStandard,
		},
	}
	newService := &scModels.Service{
		Spec: scModels.ServiceSpec{
			Description: "New description",
			Team: scModels.ServiceTeam{
				GitHubTeam: "new-team",
			},
		},
		Metadata: scModels.ServiceMetadata{
			Tier: scModels.TierCritical,
		},
	}

	// Act
	changes := processor.CompareServices(oldService, newService)

	// Assert
	assert.Len(t, changes, 3)

	// Check description change
	descChange := findChangeByField(changes, "spec.description")
	assert.NotNil(t, descChange)
	assert.Equal(t, "Old description", descChange.OldValue)
	assert.Equal(t, "New description", descChange.NewValue)

	// Check tier change
	tierChange := findChangeByField(changes, "metadata.tier")
	assert.NotNil(t, tierChange)
	assert.Equal(t, scModels.TierStandard, tierChange.OldValue)
	assert.Equal(t, scModels.TierCritical, tierChange.NewValue)

	// Check team change
	teamChange := findChangeByField(changes, "spec.team.github_team")
	assert.NotNil(t, teamChange)
	assert.Equal(t, "old-team", teamChange.OldValue)
	assert.Equal(t, "new-team", teamChange.NewValue)
}

func TestServiceProcessor_applyPagination_WithFirstPage_ReturnsFirstTwoServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{Metadata: scModels.ServiceMetadata{Name: "service1"}},
		{Metadata: scModels.ServiceMetadata{Name: "service2"}},
		{Metadata: scModels.ServiceMetadata{Name: "service3"}},
		{Metadata: scModels.ServiceMetadata{Name: "service4"}},
		{Metadata: scModels.ServiceMetadata{Name: "service5"}},
	}
	serviceList := &scModels.ServiceList{
		Services: services,
		Total:    len(services),
		Filters:  &scModels.ServiceFilter{},
	}
	limit := 2
	offset := 0

	// Act
	result := processor.applyPagination(serviceList, limit, offset)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 2)
	assert.Equal(t, 5, result.Total)
}

func TestServiceProcessor_applyPagination_WithMiddlePage_ReturnsMiddleTwoServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{Metadata: scModels.ServiceMetadata{Name: "service1"}},
		{Metadata: scModels.ServiceMetadata{Name: "service2"}},
		{Metadata: scModels.ServiceMetadata{Name: "service3"}},
		{Metadata: scModels.ServiceMetadata{Name: "service4"}},
		{Metadata: scModels.ServiceMetadata{Name: "service5"}},
	}
	serviceList := &scModels.ServiceList{
		Services: services,
		Total:    len(services),
		Filters:  &scModels.ServiceFilter{},
	}
	limit := 2
	offset := 2

	// Act
	result := processor.applyPagination(serviceList, limit, offset)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 2)
	assert.Equal(t, 5, result.Total)
}

func TestServiceProcessor_applyPagination_WithLastPage_ReturnsLastService(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{Metadata: scModels.ServiceMetadata{Name: "service1"}},
		{Metadata: scModels.ServiceMetadata{Name: "service2"}},
		{Metadata: scModels.ServiceMetadata{Name: "service3"}},
		{Metadata: scModels.ServiceMetadata{Name: "service4"}},
		{Metadata: scModels.ServiceMetadata{Name: "service5"}},
	}
	serviceList := &scModels.ServiceList{
		Services: services,
		Total:    len(services),
		Filters:  &scModels.ServiceFilter{},
	}
	limit := 2
	offset := 4

	// Act
	result := processor.applyPagination(serviceList, limit, offset)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 1)
	assert.Equal(t, 5, result.Total)
}

func TestServiceProcessor_applyPagination_WithOffsetBeyondTotal_ReturnsEmptyList(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{Metadata: scModels.ServiceMetadata{Name: "service1"}},
		{Metadata: scModels.ServiceMetadata{Name: "service2"}},
		{Metadata: scModels.ServiceMetadata{Name: "service3"}},
		{Metadata: scModels.ServiceMetadata{Name: "service4"}},
		{Metadata: scModels.ServiceMetadata{Name: "service5"}},
	}
	serviceList := &scModels.ServiceList{
		Services: services,
		Total:    len(services),
		Filters:  &scModels.ServiceFilter{},
	}
	limit := 2
	offset := 10

	// Act
	result := processor.applyPagination(serviceList, limit, offset)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 0)
	assert.Equal(t, 0, result.Total)
}

func TestServiceProcessor_applyPagination_WithNoLimit_ReturnsAllServices(t *testing.T) {
	// Arrange
	processor := NewServiceProcessor()
	services := []scModels.Service{
		{Metadata: scModels.ServiceMetadata{Name: "service1"}},
		{Metadata: scModels.ServiceMetadata{Name: "service2"}},
		{Metadata: scModels.ServiceMetadata{Name: "service3"}},
		{Metadata: scModels.ServiceMetadata{Name: "service4"}},
		{Metadata: scModels.ServiceMetadata{Name: "service5"}},
	}
	serviceList := &scModels.ServiceList{
		Services: services,
		Total:    len(services),
		Filters:  &scModels.ServiceFilter{},
	}
	limit := 0
	offset := 0

	// Act
	result := processor.applyPagination(serviceList, limit, offset)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.Services, 5)
	assert.Equal(t, 5, result.Total)
}

// Helper function to find a change by field name
func findChangeByField(changes []scModels.ServiceFieldChange, field string) *scModels.ServiceFieldChange {
	for _, change := range changes {
		if change.Field == field {
			return &change
		}
	}
	return nil
}
