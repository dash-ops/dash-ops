package servicecatalog

import "time"

// Service represents a service in the catalog
type Service struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	DisplayName string            `json:"displayName" yaml:"displayName"`
	Description string            `json:"description" yaml:"description"`
	Tier        string            `json:"tier" yaml:"tier"` // tier-1, tier-2, tier-3
	Team        string            `json:"team" yaml:"team"`
	Squad       string            `json:"squad" yaml:"squad"`
	Owner       string            `json:"owner" yaml:"owner"`
	Tags        []string          `json:"tags" yaml:"tags"`
	Regions     []string          `json:"regions" yaml:"regions"`
	IngressType string            `json:"ingressType" yaml:"ingressType"` // internal, external
	Status      string            `json:"status" yaml:"status"`           // active, inactive, deprecated
	CreatedAt   time.Time         `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt" yaml:"updatedAt"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// ServiceSummary represents a simplified view for listing
type ServiceSummary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	Tier        string    `json:"tier"`
	Team        string    `json:"team"`
	Squad       string    `json:"squad"`
	Tags        []string  `json:"tags"`
	Regions     []string  `json:"regions"`
	IngressType string    `json:"ingressType"`
	Status      string    `json:"status"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateServiceRequest represents the request to create a new service
type CreateServiceRequest struct {
	Name        string   `json:"name" binding:"required"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description" binding:"required"`
	Tier        string   `json:"tier" binding:"required"`
	Team        string   `json:"team" binding:"required"`
	Squad       string   `json:"squad" binding:"required"`
	Tags        []string `json:"tags"`
}

// ServiceFilter represents filters for service listing
type ServiceFilter struct {
	Tier   string `json:"tier,omitempty"`
	Team   string `json:"team,omitempty"`
	Status string `json:"status,omitempty"`
	Search string `json:"search,omitempty"`
}
