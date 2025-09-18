package models

import (
	"time"
)

// Dashboard represents a dashboard configuration
type Dashboard struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Service     string    `json:"service,omitempty"`
	Charts      []Chart   `json:"charts"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Owner       string    `json:"owner,omitempty"`
	Public      bool      `json:"public"`
}

// Chart represents a chart configuration in a dashboard
type Chart struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Type         string                 `json:"type"` // line, area, bar, pie, table
	Metrics      []string               `json:"metrics"`
	ServiceScope string                 `json:"service_scope"` // all, specific
	TimeRange    string                 `json:"time_range"`
	Height       int                    `json:"height"`
	Config       map[string]interface{} `json:"config,omitempty"`
	Position     *ChartPosition         `json:"position,omitempty"`
}

// ChartPosition represents the position of a chart in a dashboard
type ChartPosition struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// DashboardTemplate represents a dashboard template
type DashboardTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Service     string   `json:"service,omitempty"`
	Charts      []Chart  `json:"charts"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags,omitempty"`
}

// DashboardData represents processed data for dashboard visualization
type DashboardData struct {
	DashboardID string                 `json:"dashboard_id"`
	Charts      []ChartData            `json:"charts"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	GeneratedAt time.Time              `json:"generated_at"`
}

// ChartData represents processed data for a specific chart
type ChartData struct {
	ChartID string                 `json:"chart_id"`
	Data    []interface{}          `json:"data"`
	Config  map[string]interface{} `json:"config,omitempty"`
}
