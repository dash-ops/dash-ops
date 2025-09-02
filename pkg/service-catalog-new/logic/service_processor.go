package servicecatalog

import (
	"fmt"
	"strings"
	"time"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog-new/models"
)

// ServiceProcessor handles service processing logic
type ServiceProcessor struct{}

// NewServiceProcessor creates a new service processor
func NewServiceProcessor() *ServiceProcessor {
	return &ServiceProcessor{}
}

// PrepareForCreation prepares a service for creation
func (sp *ServiceProcessor) PrepareForCreation(service *scModels.Service, user *scModels.UserContext) (*scModels.Service, error) {
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	// Create a copy to avoid modifying the original
	prepared := *service

	// Set defaults
	prepared.SetDefaults()

	// Set audit metadata
	if user != nil {
		prepared.Metadata.CreatedBy = user.Username
		prepared.Metadata.UpdatedBy = user.Username
	}

	// Normalize service name
	prepared.Metadata.Name = sp.normalizeServiceName(prepared.Metadata.Name)

	// Process team information
	sp.processTeamInfo(&prepared.Spec.Team)

	// Process dependencies
	prepared.Spec.Business.Dependencies = sp.normalizeDependencies(prepared.Spec.Business.Dependencies)

	return &prepared, nil
}

// PrepareForUpdate prepares a service for update
func (sp *ServiceProcessor) PrepareForUpdate(service *scModels.Service, existingService *scModels.Service, user *scModels.UserContext) (*scModels.Service, error) {
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	if existingService == nil {
		return nil, fmt.Errorf("existing service cannot be nil")
	}

	// Create a copy to avoid modifying the original
	prepared := *service

	// Preserve immutable fields
	prepared.APIVersion = existingService.APIVersion
	prepared.Kind = existingService.Kind
	prepared.Metadata.Name = existingService.Metadata.Name
	prepared.Metadata.CreatedAt = existingService.Metadata.CreatedAt
	prepared.Metadata.CreatedBy = existingService.Metadata.CreatedBy

	// Update audit metadata
	prepared.Metadata.UpdatedAt = time.Now()
	prepared.Metadata.Version = existingService.Metadata.Version + 1

	if user != nil {
		prepared.Metadata.UpdatedBy = user.Username
	}

	// Process team information
	sp.processTeamInfo(&prepared.Spec.Team)

	// Process dependencies
	prepared.Spec.Business.Dependencies = sp.normalizeDependencies(prepared.Spec.Business.Dependencies)

	return &prepared, nil
}

// CalculateServiceHealth calculates overall service health
func (sp *ServiceProcessor) CalculateServiceHealth(environments []scModels.EnvironmentHealth, tier scModels.ServiceTier) scModels.ServiceStatus {
	if len(environments) == 0 {
		return scModels.StatusUnknown
	}

	// Find production environment
	var prodHealth scModels.ServiceStatus
	var hasProduction bool

	for _, env := range environments {
		if strings.EqualFold(env.Name, "production") {
			prodHealth = env.Status
			hasProduction = true
			break
		}
	}

	// If no production environment, use first environment
	if !hasProduction && len(environments) > 0 {
		prodHealth = environments[0].Status
	}

	// Apply tier-based logic
	switch tier {
	case scModels.TierCritical:
		// Critical services: any production issues are critical
		if prodHealth == scModels.StatusDown || prodHealth == scModels.StatusDegraded {
			return scModels.StatusCritical
		}
		return prodHealth

	case scModels.TierImportant:
		// Important services: production issues are degraded
		if prodHealth == scModels.StatusDown {
			return scModels.StatusDegraded
		}
		return prodHealth

	case scModels.TierStandard:
		// Standard services: only complete failure is concerning
		if prodHealth == scModels.StatusDown {
			return scModels.StatusDegraded
		}
		return scModels.StatusHealthy

	default:
		return prodHealth
	}
}

// ProcessServiceList processes a list of services with filtering and pagination
func (sp *ServiceProcessor) ProcessServiceList(services []scModels.Service, filter *scModels.ServiceFilter) *scModels.ServiceList {
	serviceList := &scModels.ServiceList{
		Services: services,
		Total:    len(services),
		Filters:  filter,
	}

	if filter == nil {
		return serviceList
	}

	// Apply filters
	if filter.Team != "" {
		serviceList = serviceList.FilterByTeam(filter.Team)
	}

	if filter.Tier != "" {
		serviceList = serviceList.FilterByTier(filter.Tier)
	}

	if filter.Search != "" {
		serviceList = serviceList.Search(filter.Search)
	}

	// Apply pagination
	if filter.Limit > 0 {
		serviceList = sp.applyPagination(serviceList, filter.Limit, filter.Offset)
	}

	return serviceList
}

// normalizeServiceName normalizes service name format
func (sp *ServiceProcessor) normalizeServiceName(name string) string {
	// Convert to lowercase and replace spaces/special chars with hyphens
	normalized := strings.ToLower(name)
	normalized = strings.ReplaceAll(normalized, " ", "-")
	normalized = strings.ReplaceAll(normalized, "_", "-")

	// Remove consecutive hyphens
	for strings.Contains(normalized, "--") {
		normalized = strings.ReplaceAll(normalized, "--", "-")
	}

	// Trim hyphens from start and end
	normalized = strings.Trim(normalized, "-")

	return normalized
}

// processTeamInfo processes team information
func (sp *ServiceProcessor) processTeamInfo(team *scModels.ServiceTeam) {
	// Normalize team name
	team.GitHubTeam = strings.TrimSpace(team.GitHubTeam)
	team.GitHubTeam = strings.ToLower(team.GitHubTeam)
}

// normalizeDependencies normalizes service dependencies
func (sp *ServiceProcessor) normalizeDependencies(dependencies []string) []string {
	seen := make(map[string]bool)
	var normalized []string

	for _, dep := range dependencies {
		dep = strings.TrimSpace(dep)
		dep = sp.normalizeServiceName(dep)

		if dep != "" && !seen[dep] {
			seen[dep] = true
			normalized = append(normalized, dep)
		}
	}

	return normalized
}

// applyPagination applies pagination to service list
func (sp *ServiceProcessor) applyPagination(serviceList *scModels.ServiceList, limit, offset int) *scModels.ServiceList {
	total := len(serviceList.Services)

	if offset >= total {
		return &scModels.ServiceList{
			Services: []scModels.Service{},
			Total:    0,
			Filters:  serviceList.Filters,
		}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return &scModels.ServiceList{
		Services: serviceList.Services[offset:end],
		Total:    total, // Keep original total for pagination info
		Filters:  serviceList.Filters,
	}
}

// GenerateServiceID generates a unique service ID
func (sp *ServiceProcessor) GenerateServiceID(serviceName string) string {
	// For now, use normalized service name as ID
	// In production, this could be a UUID or hash
	return sp.normalizeServiceName(serviceName)
}

// CompareServices compares two services and returns field changes
func (sp *ServiceProcessor) CompareServices(oldService, newService *scModels.Service) []scModels.ServiceFieldChange {
	var changes []scModels.ServiceFieldChange

	// Compare basic fields
	if oldService.Spec.Description != newService.Spec.Description {
		changes = append(changes, scModels.ServiceFieldChange{
			Field:    "spec.description",
			OldValue: oldService.Spec.Description,
			NewValue: newService.Spec.Description,
		})
	}

	if oldService.Metadata.Tier != newService.Metadata.Tier {
		changes = append(changes, scModels.ServiceFieldChange{
			Field:    "metadata.tier",
			OldValue: oldService.Metadata.Tier,
			NewValue: newService.Metadata.Tier,
		})
	}

	if oldService.Spec.Team.GitHubTeam != newService.Spec.Team.GitHubTeam {
		changes = append(changes, scModels.ServiceFieldChange{
			Field:    "spec.team.github_team",
			OldValue: oldService.Spec.Team.GitHubTeam,
			NewValue: newService.Spec.Team.GitHubTeam,
		})
	}

	// Add more field comparisons as needed...

	return changes
}
