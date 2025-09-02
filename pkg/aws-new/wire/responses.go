package aws

import "time"

// AccountResponse represents AWS account response
type AccountResponse struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Region string `json:"region"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// AccountListResponse represents account list response
type AccountListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
	Total    int               `json:"total"`
}

// InstanceResponse represents EC2 instance response
type InstanceResponse struct {
	InstanceID     string                  `json:"instance_id"`
	Name           string                  `json:"name"`
	State          InstanceStateResponse   `json:"state"`
	Platform       string                  `json:"platform"`
	InstanceType   string                  `json:"instance_type"`
	PublicIP       string                  `json:"public_ip"`
	PrivateIP      string                  `json:"private_ip"`
	SubnetID       string                  `json:"subnet_id,omitempty"`
	VpcID          string                  `json:"vpc_id,omitempty"`
	CPU            InstanceCPUResponse     `json:"cpu"`
	Memory         InstanceMemoryResponse  `json:"memory"`
	Tags           []TagResponse           `json:"tags"`
	LaunchTime     time.Time               `json:"launch_time"`
	Account        string                  `json:"account"`
	Region         string                  `json:"region"`
	SecurityGroups []SecurityGroupResponse `json:"security_groups,omitempty"`
	CostEstimate   float64                 `json:"cost_estimate"`
}

// InstanceStateResponse represents instance state response
type InstanceStateResponse struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

// InstanceCPUResponse represents instance CPU response
type InstanceCPUResponse struct {
	VCPUs       int     `json:"vcpus"`
	Utilization float64 `json:"utilization,omitempty"`
}

// InstanceMemoryResponse represents instance memory response
type InstanceMemoryResponse struct {
	TotalGB     float64 `json:"total_gb"`
	Utilization float64 `json:"utilization,omitempty"`
}

// TagResponse represents tag response
type TagResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SecurityGroupResponse represents security group response
type SecurityGroupResponse struct {
	GroupID     string `json:"group_id"`
	GroupName   string `json:"group_name"`
	Description string `json:"description,omitempty"`
}

// InstanceListResponse represents instance list response
type InstanceListResponse struct {
	Instances []InstanceResponse `json:"instances"`
	Total     int                `json:"total"`
	Account   string             `json:"account"`
	Region    string             `json:"region"`
	Filter    interface{}        `json:"filter,omitempty"`
}

// InstanceOperationResponse represents instance operation response
type InstanceOperationResponse struct {
	InstanceID    string                `json:"instance_id"`
	Operation     string                `json:"operation"`
	CurrentState  InstanceStateResponse `json:"current_state"`
	PreviousState InstanceStateResponse `json:"previous_state"`
	Success       bool                  `json:"success"`
	Message       string                `json:"message,omitempty"`
	Timestamp     time.Time             `json:"timestamp"`
}

// BatchOperationResponse represents batch operation response
type BatchOperationResponse struct {
	Operation    string                      `json:"operation"`
	TotalCount   int                         `json:"total_count"`
	SuccessCount int                         `json:"success_count"`
	FailureCount int                         `json:"failure_count"`
	Results      []InstanceOperationResponse `json:"results"`
	StartedAt    time.Time                   `json:"started_at"`
	CompletedAt  time.Time                   `json:"completed_at,omitempty"`
	Duration     string                      `json:"duration"`
	SuccessRate  float64                     `json:"success_rate"`
}

// AccountSummaryResponse represents account summary response
type AccountSummaryResponse struct {
	Account              string    `json:"account"`
	Region               string    `json:"region"`
	TotalInstances       int       `json:"total_instances"`
	RunningInstances     int       `json:"running_instances"`
	StoppedInstances     int       `json:"stopped_instances"`
	PendingInstances     int       `json:"pending_instances"`
	EstimatedMonthlyCost float64   `json:"estimated_monthly_cost"`
	LastUpdated          time.Time `json:"last_updated"`
}

// InstanceMetricsResponse represents instance metrics response
type InstanceMetricsResponse struct {
	InstanceID  string                       `json:"instance_id"`
	Account     string                       `json:"account"`
	Region      string                       `json:"region"`
	Metrics     []InstanceMetricDataResponse `json:"metrics"`
	Period      string                       `json:"period"`
	LastUpdated time.Time                    `json:"last_updated"`
}

// InstanceMetricDataResponse represents metric data response
type InstanceMetricDataResponse struct {
	MetricName string                    `json:"metric_name"`
	Unit       string                    `json:"unit"`
	DataPoints []MetricDataPointResponse `json:"data_points"`
}

// MetricDataPointResponse represents metric data point response
type MetricDataPointResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
}

// CostSavingsResponse represents cost savings analysis response
type CostSavingsResponse struct {
	CurrentMonthlyCost float64   `json:"current_monthly_cost"`
	PotentialSavings   float64   `json:"potential_savings"`
	StoppableInstances int       `json:"stoppable_instances"`
	SavingsPercentage  float64   `json:"savings_percentage"`
	LastCalculated     time.Time `json:"last_calculated"`
}

// OperationCostEstimateResponse represents operation cost estimate response
type OperationCostEstimateResponse struct {
	InstanceID     string    `json:"instance_id"`
	Operation      string    `json:"operation"`
	HourlyCost     float64   `json:"hourly_cost"`
	MonthlyCost    float64   `json:"monthly_cost"`
	CostImpact     float64   `json:"cost_impact"`
	ImpactType     string    `json:"impact_type"`
	Description    string    `json:"description"`
	LastCalculated time.Time `json:"last_calculated"`
}

// RegionInfoResponse represents region information response
type RegionInfoResponse struct {
	RegionName string `json:"region_name"`
	RegionCode string `json:"region_code"`
	Endpoint   string `json:"endpoint"`
	Available  bool   `json:"available"`
}

// InstanceTypeInfoResponse represents instance type information response
type InstanceTypeInfoResponse struct {
	InstanceType string                      `json:"instance_type"`
	VCPUs        int                         `json:"vcpus"`
	Memory       float64                     `json:"memory_gb"`
	Storage      InstanceTypeStorageResponse `json:"storage"`
	Network      InstanceTypeNetworkResponse `json:"network"`
	Pricing      InstanceTypePricingResponse `json:"pricing,omitempty"`
}

// InstanceTypeStorageResponse represents instance type storage response
type InstanceTypeStorageResponse struct {
	Type         string `json:"type"`
	SizeGB       int    `json:"size_gb,omitempty"`
	IOPS         int    `json:"iops,omitempty"`
	EBSOptimized bool   `json:"ebs_optimized"`
}

// InstanceTypeNetworkResponse represents instance type network response
type InstanceTypeNetworkResponse struct {
	Performance string `json:"performance"`
	IPv6Support bool   `json:"ipv6_support"`
	ENASupport  bool   `json:"ena_support"`
}

// InstanceTypePricingResponse represents instance type pricing response
type InstanceTypePricingResponse struct {
	OnDemand PricingInfoResponse `json:"on_demand"`
	Reserved PricingInfoResponse `json:"reserved,omitempty"`
	Spot     PricingInfoResponse `json:"spot,omitempty"`
}

// PricingInfoResponse represents pricing information response
type PricingInfoResponse struct {
	HourlyRate  float64   `json:"hourly_rate"`
	MonthlyRate float64   `json:"monthly_rate"`
	Currency    string    `json:"currency"`
	LastUpdated time.Time `json:"last_updated"`
}
