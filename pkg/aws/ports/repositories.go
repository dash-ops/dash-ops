package aws

import (
	"context"
	"time"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// AccountRepository defines the interface for AWS account data access
type AccountRepository interface {
	// GetAccount gets a specific AWS account configuration
	GetAccount(ctx context.Context, accountKey string) (*awsModels.AWSAccount, error)

	// ListAccounts lists all configured AWS accounts
	ListAccounts(ctx context.Context) ([]awsModels.AWSAccount, error)

	// ValidateAccount validates account credentials and connectivity
	ValidateAccount(ctx context.Context, account *awsModels.AWSAccount) error

	// UpdateAccountStatus updates account status and error information
	UpdateAccountStatus(ctx context.Context, accountKey string, status awsModels.AccountStatus, errorMsg string) error
}

// InstanceRepository defines the interface for EC2 instance data access
type InstanceRepository interface {
	// GetInstance gets a specific EC2 instance
	GetInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.EC2Instance, error)

	// ListInstances lists EC2 instances with optional filtering
	ListInstances(ctx context.Context, accountKey, region string, filter *awsModels.InstanceFilter) (*awsModels.InstanceList, error)

	// StartInstance starts an EC2 instance
	StartInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceOperation, error)

	// StopInstance stops an EC2 instance
	StopInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceOperation, error)

	// RestartInstance restarts an EC2 instance
	RestartInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceOperation, error)

	// GetInstanceStatus gets current instance status
	GetInstanceStatus(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceState, error)

	// BatchOperation performs batch operations on multiple instances
	BatchOperation(ctx context.Context, accountKey, region string, operation string, instanceIDs []string) (*awsModels.BatchOperation, error)
}

// MetricsRepository defines the interface for CloudWatch metrics data access
type MetricsRepository interface {
	// GetInstanceMetrics gets CloudWatch metrics for an instance
	GetInstanceMetrics(ctx context.Context, accountKey, region, instanceID string, period string) (*awsModels.InstanceMetrics, error)

	// GetAccountMetrics gets aggregated metrics for an account
	GetAccountMetrics(ctx context.Context, accountKey, region string) (*AccountMetrics, error)

	// GetCostMetrics gets cost and billing metrics
	GetCostMetrics(ctx context.Context, accountKey string, period string) (*CostMetrics, error)
}

// AccountMetrics represents aggregated account metrics
type AccountMetrics struct {
	Account     string                      `json:"account"`
	Region      string                      `json:"region"`
	Summary     awsModels.AccountSummary    `json:"summary"`
	Instances   []awsModels.InstanceMetrics `json:"instances"`
	LastUpdated time.Time                   `json:"last_updated"`
}

// CostMetrics represents cost and billing metrics
type CostMetrics struct {
	Account      string        `json:"account"`
	Period       string        `json:"period"`
	TotalCost    float64       `json:"total_cost"`
	ServiceCosts []ServiceCost `json:"service_costs"`
	DailyCosts   []DailyCost   `json:"daily_costs"`
	LastUpdated  time.Time     `json:"last_updated"`
}

// ServiceCost represents cost breakdown by AWS service
type ServiceCost struct {
	ServiceName string  `json:"service_name"`
	Cost        float64 `json:"cost"`
	Percentage  float64 `json:"percentage"`
}

// DailyCost represents daily cost data
type DailyCost struct {
	Date string  `json:"date"`
	Cost float64 `json:"cost"`
}
