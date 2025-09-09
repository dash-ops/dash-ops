package aws

import (
	"context"
	"time"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// AWSClientService defines the interface for AWS SDK operations
type AWSClientService interface {
	// GetEC2Client gets an EC2 client for a specific account
	GetEC2Client(account *awsModels.AWSAccount) (EC2Client, error)

	// GetCloudWatchClient gets a CloudWatch client for a specific account
	GetCloudWatchClient(account *awsModels.AWSAccount) (CloudWatchClient, error)

	// ValidateCredentials validates AWS credentials
	ValidateCredentials(account *awsModels.AWSAccount) error

	// GetAccountInfo gets basic account information
	GetAccountInfo(account *awsModels.AWSAccount) (*AccountInfo, error)
}

// EC2Client defines the interface for EC2 operations
type EC2Client interface {
	// DescribeInstances describes EC2 instances
	DescribeInstances(ctx context.Context, filter *EC2Filter) ([]awsModels.EC2Instance, error)

	// DescribeInstance describes a specific EC2 instance
	DescribeInstance(ctx context.Context, instanceID string) (*awsModels.EC2Instance, error)

	// StartInstances starts EC2 instances
	StartInstances(ctx context.Context, instanceIDs []string) ([]awsModels.InstanceOperation, error)

	// StopInstances stops EC2 instances
	StopInstances(ctx context.Context, instanceIDs []string) ([]awsModels.InstanceOperation, error)

	// RebootInstances reboots EC2 instances
	RebootInstances(ctx context.Context, instanceIDs []string) ([]awsModels.InstanceOperation, error)

	// DescribeInstanceTypes describes available instance types
	DescribeInstanceTypes(ctx context.Context) ([]InstanceTypeInfo, error)

	// DescribeRegions describes available regions
	DescribeRegions(ctx context.Context) ([]RegionInfo, error)
}

// CloudWatchClient defines the interface for CloudWatch operations
type CloudWatchClient interface {
	// GetInstanceMetrics gets CloudWatch metrics for an instance
	GetInstanceMetrics(ctx context.Context, instanceID string, period time.Duration, startTime, endTime time.Time) (*awsModels.InstanceMetrics, error)

	// GetMetricStatistics gets metric statistics
	GetMetricStatistics(ctx context.Context, metricName, namespace string, dimensions map[string]string, period time.Duration, startTime, endTime time.Time) ([]awsModels.MetricDataPoint, error)

	// ListMetrics lists available metrics
	ListMetrics(ctx context.Context, namespace string) ([]MetricInfo, error)
}

// EC2Filter represents EC2 API filtering options
type EC2Filter struct {
	InstanceIDs   []string          `json:"instance_ids,omitempty"`
	States        []string          `json:"states,omitempty"`
	InstanceTypes []string          `json:"instance_types,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	MaxResults    int               `json:"max_results,omitempty"`
}

// AccountInfo represents basic AWS account information
type AccountInfo struct {
	AccountID   string    `json:"account_id"`
	Alias       string    `json:"alias,omitempty"`
	Region      string    `json:"region"`
	Status      string    `json:"status"`
	LastChecked time.Time `json:"last_checked"`
}

// InstanceTypeInfo represents EC2 instance type information
type InstanceTypeInfo struct {
	InstanceType string              `json:"instance_type"`
	VCPUs        int                 `json:"vcpus"`
	Memory       float64             `json:"memory_gb"`
	Storage      InstanceTypeStorage `json:"storage"`
	Network      InstanceTypeNetwork `json:"network"`
	Pricing      InstanceTypePricing `json:"pricing,omitempty"`
}

// InstanceTypeStorage represents instance type storage information
type InstanceTypeStorage struct {
	Type         string `json:"type"` // EBS, Instance Store
	SizeGB       int    `json:"size_gb,omitempty"`
	IOPS         int    `json:"iops,omitempty"`
	EBSOptimized bool   `json:"ebs_optimized"`
}

// InstanceTypeNetwork represents instance type network information
type InstanceTypeNetwork struct {
	Performance string `json:"performance"` // Low, Moderate, High, etc.
	IPv6Support bool   `json:"ipv6_support"`
	ENASupport  bool   `json:"ena_support"`
}

// InstanceTypePricing represents instance type pricing information
type InstanceTypePricing struct {
	OnDemand PricingInfo `json:"on_demand"`
	Reserved PricingInfo `json:"reserved,omitempty"`
	Spot     PricingInfo `json:"spot,omitempty"`
}

// PricingInfo represents pricing information
type PricingInfo struct {
	HourlyRate  float64   `json:"hourly_rate"`
	MonthlyRate float64   `json:"monthly_rate"`
	Currency    string    `json:"currency"`
	LastUpdated time.Time `json:"last_updated"`
}

// RegionInfo represents AWS region information
type RegionInfo struct {
	RegionName string `json:"region_name"`
	RegionCode string `json:"region_code"`
	Endpoint   string `json:"endpoint"`
	Available  bool   `json:"available"`
}

// MetricInfo represents CloudWatch metric information
type MetricInfo struct {
	MetricName string            `json:"metric_name"`
	Namespace  string            `json:"namespace"`
	Dimensions []MetricDimension `json:"dimensions"`
	Unit       string            `json:"unit"`
}

// MetricDimension represents a CloudWatch metric dimension
type MetricDimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NotificationService defines the interface for AWS notifications
type NotificationService interface {
	// NotifyInstanceStateChange notifies about instance state changes
	NotifyInstanceStateChange(ctx context.Context, operation *awsModels.InstanceOperation) error

	// NotifyBatchOperationComplete notifies about batch operation completion
	NotifyBatchOperationComplete(ctx context.Context, batchOp *awsModels.BatchOperation) error

	// NotifyHighCostAlert notifies about high cost usage
	NotifyHighCostAlert(ctx context.Context, account string, cost float64, threshold float64) error

	// NotifyAccountError notifies about account connectivity issues
	NotifyAccountError(ctx context.Context, account string, error string) error
}

// AuditService defines the interface for AWS audit logging
type AuditService interface {
	// LogInstanceOperation logs instance operations
	LogInstanceOperation(ctx context.Context, operation *awsModels.InstanceOperation, userContext *UserContext) error

	// LogBatchOperation logs batch operations
	LogBatchOperation(ctx context.Context, batchOp *awsModels.BatchOperation, userContext *UserContext) error

	// LogAccountAccess logs account access events
	LogAccountAccess(ctx context.Context, account string, userContext *UserContext, action string) error

	// LogCostAlert logs cost-related alerts
	LogCostAlert(ctx context.Context, account string, cost float64, threshold float64) error
}

// UserContext represents user information for audit and permissions
type UserContext struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
	IP       string   `json:"ip,omitempty"`
}

// CostOptimizationService defines the interface for cost optimization
type CostOptimizationService interface {
	// AnalyzeCostOptimization analyzes cost optimization opportunities
	AnalyzeCostOptimization(ctx context.Context, accountKey, region string) (*CostOptimizationReport, error)

	// GetUnderutilizedInstances identifies underutilized instances
	GetUnderutilizedInstances(ctx context.Context, accountKey, region string, threshold float64) ([]awsModels.EC2Instance, error)

	// GetRightSizingRecommendations provides right-sizing recommendations
	GetRightSizingRecommendations(ctx context.Context, accountKey, region string) ([]RightSizingRecommendation, error)
}

// CostOptimizationReport represents cost optimization analysis
type CostOptimizationReport struct {
	Account                string                  `json:"account"`
	Region                 string                  `json:"region"`
	CurrentMonthlyCost     float64                 `json:"current_monthly_cost"`
	OptimizedMonthlyCost   float64                 `json:"optimized_monthly_cost"`
	PotentialSavings       float64                 `json:"potential_savings"`
	SavingsPercentage      float64                 `json:"savings_percentage"`
	Recommendations        []CostRecommendation    `json:"recommendations"`
	UnderutilizedInstances []awsModels.EC2Instance `json:"underutilized_instances"`
	LastAnalyzed           time.Time               `json:"last_analyzed"`
}

// CostRecommendation represents a cost optimization recommendation
type CostRecommendation struct {
	Type          string  `json:"type"` // stop, resize, schedule, etc.
	InstanceID    string  `json:"instance_id"`
	CurrentCost   float64 `json:"current_cost"`
	OptimizedCost float64 `json:"optimized_cost"`
	Savings       float64 `json:"savings"`
	Description   string  `json:"description"`
	Priority      string  `json:"priority"` // high, medium, low
}

// RightSizingRecommendation represents instance right-sizing recommendation
type RightSizingRecommendation struct {
	InstanceID          string  `json:"instance_id"`
	CurrentInstanceType string  `json:"current_instance_type"`
	RecommendedType     string  `json:"recommended_type"`
	CurrentCost         float64 `json:"current_cost"`
	RecommendedCost     float64 `json:"recommended_cost"`
	Savings             float64 `json:"savings"`
	CPUUtilization      float64 `json:"cpu_utilization"`
	MemoryUtilization   float64 `json:"memory_utilization"`
	Confidence          string  `json:"confidence"` // high, medium, low
	Reason              string  `json:"reason"`
}
