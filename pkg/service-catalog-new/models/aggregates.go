package servicecatalog

import (
	"strings"
	"time"
)

// ServiceList represents a list of services for API responses
type ServiceList struct {
	Services []Service      `json:"services"`
	Total    int            `json:"total"`
	Filters  *ServiceFilter `json:"filters,omitempty"`
}

// ServiceFilter represents filtering criteria for services
type ServiceFilter struct {
	Team   string        `json:"team,omitempty"`
	Tier   ServiceTier   `json:"tier,omitempty"`
	Status ServiceStatus `json:"status,omitempty"`
	Search string        `json:"search,omitempty"`
	Limit  int           `json:"limit,omitempty"`
	Offset int           `json:"offset,omitempty"`
}

// ServiceHealth represents aggregated service health status
type ServiceHealth struct {
	ServiceName   string              `json:"service_name"`
	OverallStatus ServiceStatus       `json:"overall_status"`
	Environments  []EnvironmentHealth `json:"environments"`
	LastUpdated   time.Time           `json:"last_updated"`
}

// EnvironmentHealth represents health status for a specific environment
type EnvironmentHealth struct {
	Name        string             `json:"name"`
	Context     string             `json:"context"`
	Status      ServiceStatus      `json:"status"`
	Deployments []DeploymentHealth `json:"deployments"`
}

// DeploymentHealth represents health status for a specific deployment
type DeploymentHealth struct {
	Name            string        `json:"name"`
	ReadyReplicas   int           `json:"ready_replicas"`
	DesiredReplicas int           `json:"desired_replicas"`
	Status          ServiceStatus `json:"status"`
	LastUpdated     time.Time     `json:"last_updated"`
}

// ServiceHistory represents service change history
type ServiceHistory struct {
	ServiceName string          `json:"service_name"`
	History     []ServiceChange `json:"history"`
}

// ServiceChange represents a single change in service history
type ServiceChange struct {
	Commit    string               `json:"commit"`
	Author    string               `json:"author"`
	Email     string               `json:"email"`
	Timestamp time.Time            `json:"timestamp"`
	Message   string               `json:"message"`
	Changes   []ServiceFieldChange `json:"changes,omitempty"`
}

// ServiceFieldChange represents a specific field change
type ServiceFieldChange struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

// ServiceContext represents service context resolution
type ServiceContext struct {
	Service     *Service `json:"service"`
	Environment string   `json:"environment"`
	Namespace   string   `json:"namespace"`
	Context     string   `json:"context"`
	Found       bool     `json:"found"`
}

// UserContext represents user information from OAuth2
type UserContext struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Teams    []string `json:"teams,omitempty"`
}

// Methods for ServiceList

// FilterByTeam filters services by team
func (sl *ServiceList) FilterByTeam(team string) *ServiceList {
	if team == "" {
		return sl
	}

	var filtered []Service
	for _, service := range sl.Services {
		if strings.EqualFold(service.Spec.Team.GitHubTeam, team) {
			filtered = append(filtered, service)
		}
	}

	return &ServiceList{
		Services: filtered,
		Total:    len(filtered),
		Filters:  &ServiceFilter{Team: team},
	}
}

// FilterByTier filters services by tier
func (sl *ServiceList) FilterByTier(tier ServiceTier) *ServiceList {
	if tier == "" {
		return sl
	}

	var filtered []Service
	for _, service := range sl.Services {
		if service.Metadata.Tier == tier {
			filtered = append(filtered, service)
		}
	}

	return &ServiceList{
		Services: filtered,
		Total:    len(filtered),
		Filters:  &ServiceFilter{Tier: tier},
	}
}

// FilterByStatus filters services by status
func (sl *ServiceList) FilterByStatus(status ServiceStatus) *ServiceList {
	// Note: This would require health information to be available
	// For now, return all services
	return sl
}

// Search filters services by text search
func (sl *ServiceList) Search(query string) *ServiceList {
	if query == "" {
		return sl
	}

	query = strings.ToLower(query)
	var filtered []Service

	for _, service := range sl.Services {
		if sl.matchesSearch(service, query) {
			filtered = append(filtered, service)
		}
	}

	return &ServiceList{
		Services: filtered,
		Total:    len(filtered),
		Filters:  &ServiceFilter{Search: query},
	}
}

// matchesSearch checks if service matches search query
func (sl *ServiceList) matchesSearch(service Service, query string) bool {
	searchFields := []string{
		strings.ToLower(service.Metadata.Name),
		strings.ToLower(service.Spec.Description),
		strings.ToLower(service.Spec.Team.GitHubTeam),
		strings.ToLower(service.Spec.Technology.Language),
		strings.ToLower(service.Spec.Technology.Framework),
	}

	for _, field := range searchFields {
		if strings.Contains(field, query) {
			return true
		}
	}

	return false
}

// Methods for ServiceHealth

// CalculateOverallStatus calculates overall status based on environments and tier
func (sh *ServiceHealth) CalculateOverallStatus() ServiceStatus {
	if len(sh.Environments) == 0 {
		return StatusUnknown
	}

	// Find production environment
	var prodStatus ServiceStatus
	var hasProduction bool

	for _, env := range sh.Environments {
		if env.Name == "production" {
			prodStatus = env.Status
			hasProduction = true
			break
		}
	}

	// If no production environment, use first environment
	if !hasProduction && len(sh.Environments) > 0 {
		prodStatus = sh.Environments[0].Status
	}

	return prodStatus
}

// Methods for UserContext

// HasTeam checks if user belongs to a specific team
func (uc *UserContext) HasTeam(team string) bool {
	for _, userTeam := range uc.Teams {
		if strings.EqualFold(userTeam, team) {
			return true
		}
	}
	return false
}

// HasAnyTeam checks if user belongs to any of the specified teams
func (uc *UserContext) HasAnyTeam(teams []string) bool {
	for _, team := range teams {
		if uc.HasTeam(team) {
			return true
		}
	}
	return false
}
