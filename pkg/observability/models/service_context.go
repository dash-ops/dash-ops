package models

import (
	"time"
)

// ServiceContext represents service context for observability
type ServiceContext struct {
	ServiceName string                 `json:"service_name"`
	Namespace   string                 `json:"namespace"`
	Cluster     string                 `json:"cluster"`
	Labels      map[string]string      `json:"labels"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Health      *ServiceHealth         `json:"health,omitempty"`
}

// ServiceWithContext represents a service with its observability context
type ServiceWithContext struct {
	ServiceContext
	LogCount    int64 `json:"log_count"`
	MetricCount int64 `json:"metric_count"`
	TraceCount  int64 `json:"trace_count"`
	AlertCount  int64 `json:"alert_count"`
}

// ServiceHealth represents the health status of a service
type ServiceHealth struct {
	Status    string                 `json:"status"`
	LastCheck time.Time              `json:"last_check"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Metrics   map[string]float64     `json:"metrics,omitempty"`
	Alerts    []string               `json:"alerts,omitempty"`
}
