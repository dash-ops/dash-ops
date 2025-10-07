package wire

import (
	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// DashboardsRequest represents a request for dashboard data
type DashboardsRequest struct {
	Service string `json:"service,omitempty"`
	Owner   string `json:"owner,omitempty"`
	Public  *bool  `json:"public,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
}

// DashboardsResponse represents the response for dashboard data
type DashboardsResponse struct {
	BaseResponse
	Data DashboardsData `json:"data"`
}

// DashboardsData represents the data portion of dashboards response
type DashboardsData struct {
	Dashboards []models.Dashboard `json:"dashboards"`
	Total      int                `json:"total"`
	Filters    DashboardFilters   `json:"filters,omitempty"`
}

// DashboardFilters represents filters for dashboard queries
type DashboardFilters struct {
	Owners   []string `json:"owners,omitempty"`
	Services []string `json:"services,omitempty"`
	Public   *bool    `json:"public,omitempty"`
}

// CreateDashboardRequest represents a request to create a dashboard
type CreateDashboardRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Service     string         `json:"service,omitempty"`
	Charts      []models.Chart `json:"charts"`
	Public      bool           `json:"public"`
}

// UpdateDashboardRequest represents a request to update a dashboard
type UpdateDashboardRequest struct {
	CreateDashboardRequest
}

// DeleteDashboardRequest represents a request to delete a dashboard
type DeleteDashboardRequest struct {
	ID string `json:"id"`
}

// GetDashboardRequest represents a request to get a specific dashboard
type GetDashboardRequest struct {
	ID string `json:"id"`
}

// DashboardResponse represents the response for a single dashboard operation
type DashboardResponse struct {
	BaseResponse
	Data models.Dashboard `json:"data"`
}

// DashboardTemplatesRequest represents a request for dashboard templates
type DashboardTemplatesRequest struct {
	Category string   `json:"category,omitempty"`
	Service  string   `json:"service,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// DashboardTemplatesResponse represents the response for dashboard templates
type DashboardTemplatesResponse struct {
	BaseResponse
	Data DashboardTemplatesData `json:"data"`
}

// DashboardTemplatesData represents the data portion of dashboard templates response
type DashboardTemplatesData struct {
	Templates  []models.DashboardTemplate `json:"templates"`
	Total      int                        `json:"total"`
	Categories []string                   `json:"categories,omitempty"`
	Tags       []string                   `json:"tags,omitempty"`
}

// DashboardStatsRequest represents a request for dashboard statistics
type DashboardStatsRequest struct {
	Owner string `json:"owner,omitempty"`
}

// DashboardStatisticsResponse represents the response for dashboard statistics
type DashboardStatisticsResponse struct {
	BaseResponse
	Data DashboardStatistics `json:"data"`
}

// DashboardStatistics represents dashboard statistics
type DashboardStatistics struct {
	TotalDashboards     int64                  `json:"total_dashboards"`
	DashboardsByOwner   map[string]int64       `json:"dashboards_by_owner"`
	DashboardsByService map[string]int64       `json:"dashboards_by_service"`
	PublicDashboards    int64                  `json:"public_dashboards"`
	PrivateDashboards   int64                  `json:"private_dashboards"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}
