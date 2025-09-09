package kubernetes

import "time"

// HealthStatus represents the health status of a resource
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// ResourceStatus represents the status of a resource
type ResourceStatus string

const (
	ResourceStatusHealthy  ResourceStatus = "healthy"
	ResourceStatusWarning  ResourceStatus = "warning"
	ResourceStatusCritical ResourceStatus = "critical"
	ResourceStatusUnknown  ResourceStatus = "unknown"
)

// ClusterHealth represents the overall health of a cluster
type ClusterHealth struct {
	Context     string         `json:"context"`
	Status      HealthStatus   `json:"status"`
	Nodes       []NodeHealth   `json:"nodes"`
	Summary     ClusterSummary `json:"summary"`
	LastUpdated time.Time      `json:"last_updated"`
}

// NodeHealth represents the health of a node
type NodeHealth struct {
	Name        string          `json:"name"`
	Status      NodeStatus      `json:"status"`
	Conditions  []NodeCondition `json:"conditions"`
	Resources   ResourceHealth  `json:"resources"`
	LastUpdated time.Time       `json:"last_updated"`
}

// ResourceHealth represents the health of resources
type ResourceHealth struct {
	CPU    ResourceHealthDetail `json:"cpu"`
	Memory ResourceHealthDetail `json:"memory"`
	Pods   ResourceHealthDetail `json:"pods,omitempty"`
}

// ResourceHealthDetail represents detailed resource health
type ResourceHealthDetail struct {
	Used               int64          `json:"used"`
	Available          int64          `json:"available"`
	Total              int64          `json:"total"`
	UtilizationPercent float64        `json:"utilization_percent"`
	Status             ResourceStatus `json:"status"`
}
