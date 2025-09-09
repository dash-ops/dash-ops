package aws

// InstanceOperationRequest represents instance operation request
type InstanceOperationRequest struct {
	InstanceIDs []string `json:"instance_ids" validate:"required,min=1,max=20"`
	Force       bool     `json:"force,omitempty"` // Force operation even if instance is in transitioning state
}

// InstanceFilterRequest represents instance filtering request
type InstanceFilterRequest struct {
	Account      string             `json:"account,omitempty"`
	Region       string             `json:"region,omitempty"`
	State        string             `json:"state,omitempty" validate:"omitempty,oneof=pending running shutting-down terminated stopping stopped"`
	InstanceType string             `json:"instance_type,omitempty"`
	Platform     string             `json:"platform,omitempty"`
	Tags         []TagFilterRequest `json:"tags,omitempty"`
	Search       string             `json:"search,omitempty"`
	Limit        int                `json:"limit,omitempty" validate:"omitempty,min=1,max=1000"`
	Offset       int                `json:"offset,omitempty" validate:"omitempty,min=0"`
}

// TagFilterRequest represents tag filtering request
type TagFilterRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value,omitempty"`
}

// BatchOperationRequest represents batch operation request
type BatchOperationRequest struct {
	Operation   string   `json:"operation" validate:"required,oneof=start stop restart"`
	InstanceIDs []string `json:"instance_ids" validate:"required,min=1,max=50"`
	Force       bool     `json:"force,omitempty"`
}

// MetricsRequest represents metrics request
type MetricsRequest struct {
	Period      string   `json:"period,omitempty" validate:"omitempty,oneof=5m 1h 6h 1d 7d 30d"`
	StartTime   string   `json:"start_time,omitempty"` // RFC3339 format
	EndTime     string   `json:"end_time,omitempty"`   // RFC3339 format
	MetricNames []string `json:"metric_names,omitempty"`
}

// CostAnalysisRequest represents cost analysis request
type CostAnalysisRequest struct {
	Period             string  `json:"period,omitempty" validate:"omitempty,oneof=1d 7d 30d 90d"`
	IncludeForecasting bool    `json:"include_forecasting,omitempty"`
	CostThreshold      float64 `json:"cost_threshold,omitempty"` // Alert threshold
}
