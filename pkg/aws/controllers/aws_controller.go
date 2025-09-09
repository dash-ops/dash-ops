package aws

import (
	"context"
	"fmt"

	awsLogic "github.com/dash-ops/dash-ops/pkg/aws/logic"
	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
	awsPorts "github.com/dash-ops/dash-ops/pkg/aws/ports"
)

// AWSController handles AWS business logic orchestration
type AWSController struct {
	accountRepo    awsPorts.AccountRepository
	instanceRepo   awsPorts.InstanceRepository
	metricsRepo    awsPorts.MetricsRepository
	processor      *awsLogic.InstanceProcessor
	costCalculator *awsLogic.CostCalculator
}

// NewAWSController creates a new AWS controller
func NewAWSController(
	accountRepo awsPorts.AccountRepository,
	instanceRepo awsPorts.InstanceRepository,
	metricsRepo awsPorts.MetricsRepository,
	processor *awsLogic.InstanceProcessor,
	costCalculator *awsLogic.CostCalculator,
) *AWSController {
	return &AWSController{
		accountRepo:    accountRepo,
		instanceRepo:   instanceRepo,
		metricsRepo:    metricsRepo,
		processor:      processor,
		costCalculator: costCalculator,
	}
}

// ListAccounts lists all configured AWS accounts
func (ac *AWSController) ListAccounts(ctx context.Context) ([]awsModels.AWSAccount, error) {
	accounts, err := ac.accountRepo.ListAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

// GetAccount gets a specific AWS account
func (ac *AWSController) GetAccount(ctx context.Context, accountKey string) (*awsModels.AWSAccount, error) {
	if accountKey == "" {
		return nil, fmt.Errorf("account key is required")
	}

	account, err := ac.accountRepo.GetAccount(ctx, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

// ListInstances lists EC2 instances with filtering and permissions
func (ac *AWSController) ListInstances(ctx context.Context, accountKey, region string, filter *awsModels.InstanceFilter, userContext *awsPorts.UserContext) (*awsModels.InstanceList, error) {
	// Get account configuration
	account, err := ac.GetAccount(ctx, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Check view permissions
	if userContext != nil && !account.HasEC2ViewPermission(userContext.Groups) {
		return nil, fmt.Errorf("user does not have permission to view instances in account %s", accountKey)
	}

	// Get instances from repository
	instanceList, err := ac.instanceRepo.ListInstances(ctx, accountKey, region, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	// Process with skip list and permissions
	processedList := ac.processor.ProcessInstanceList(instanceList.Instances, filter, account.EC2Config.SkipList)
	processedList.Account = account.Name
	processedList.Region = region

	return processedList, nil
}

// GetInstance gets a specific EC2 instance
func (ac *AWSController) GetInstance(ctx context.Context, accountKey, region, instanceID string, userContext *awsPorts.UserContext) (*awsModels.EC2Instance, error) {
	if instanceID == "" {
		return nil, fmt.Errorf("instance ID is required")
	}

	// Validate instance ID format
	if err := ac.processor.ValidateInstanceID(instanceID); err != nil {
		return nil, fmt.Errorf("invalid instance ID: %w", err)
	}

	// Get account and check permissions
	account, err := ac.GetAccount(ctx, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if userContext != nil && !account.HasEC2ViewPermission(userContext.Groups) {
		return nil, fmt.Errorf("user does not have permission to view instances in account %s", accountKey)
	}

	// Get instance
	instance, err := ac.instanceRepo.GetInstance(ctx, accountKey, region, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return instance, nil
}

// StartInstance starts an EC2 instance
func (ac *AWSController) StartInstance(ctx context.Context, accountKey, region, instanceID string, userContext *awsPorts.UserContext) (*awsModels.InstanceOperation, error) {
	// Get account and validate permissions
	account, err := ac.GetAccount(ctx, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if userContext != nil && !account.HasEC2StartPermission(userContext.Groups) {
		return nil, fmt.Errorf("user does not have permission to start instances in account %s", accountKey)
	}

	// Get current instance state for validation
	instance, err := ac.GetInstance(ctx, accountKey, region, instanceID, userContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance for validation: %w", err)
	}

	// Validate operation
	if err := ac.processor.ValidateInstanceOperation(instance, "start"); err != nil {
		return nil, fmt.Errorf("operation validation failed: %w", err)
	}

	// Perform operation
	operation, err := ac.instanceRepo.StartInstance(ctx, accountKey, region, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to start instance: %w", err)
	}

	return operation, nil
}

// StopInstance stops an EC2 instance
func (ac *AWSController) StopInstance(ctx context.Context, accountKey, region, instanceID string, userContext *awsPorts.UserContext) (*awsModels.InstanceOperation, error) {
	// Get account and validate permissions
	account, err := ac.GetAccount(ctx, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if userContext != nil && !account.HasEC2StopPermission(userContext.Groups) {
		return nil, fmt.Errorf("user does not have permission to stop instances in account %s", accountKey)
	}

	// Get current instance state for validation
	instance, err := ac.GetInstance(ctx, accountKey, region, instanceID, userContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance for validation: %w", err)
	}

	// Validate operation
	if err := ac.processor.ValidateInstanceOperation(instance, "stop"); err != nil {
		return nil, fmt.Errorf("operation validation failed: %w", err)
	}

	// Perform operation
	operation, err := ac.instanceRepo.StopInstance(ctx, accountKey, region, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to stop instance: %w", err)
	}

	return operation, nil
}

// BatchOperation performs batch operations on multiple instances
func (ac *AWSController) BatchOperation(ctx context.Context, accountKey, region, operation string, instanceIDs []string, userContext *awsPorts.UserContext) (*awsModels.BatchOperation, error) {
	if len(instanceIDs) == 0 {
		return nil, fmt.Errorf("no instance IDs provided")
	}

	// Get account and validate permissions
	account, err := ac.GetAccount(ctx, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Check permissions based on operation
	if userContext != nil {
		switch operation {
		case "start":
			if !account.HasEC2StartPermission(userContext.Groups) {
				return nil, fmt.Errorf("user does not have permission to start instances")
			}
		case "stop":
			if !account.HasEC2StopPermission(userContext.Groups) {
				return nil, fmt.Errorf("user does not have permission to stop instances")
			}
		default:
			return nil, fmt.Errorf("unsupported batch operation: %s", operation)
		}
	}

	// Get instances for validation
	var instances []awsModels.EC2Instance
	for _, instanceID := range instanceIDs {
		instance, err := ac.GetInstance(ctx, accountKey, region, instanceID, userContext)
		if err != nil {
			// Skip instances that can't be found, but log the error
			continue
		}
		instances = append(instances, *instance)
	}

	// Execute batch operation
	result, err := ac.instanceRepo.BatchOperation(ctx, accountKey, region, operation, instanceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch operation: %w", err)
	}

	return result, nil
}

// GetAccountSummary gets account resource summary
func (ac *AWSController) GetAccountSummary(ctx context.Context, accountKey, region string, userContext *awsPorts.UserContext) (*awsModels.AccountSummary, error) {
	// Get account
	account, err := ac.GetAccount(ctx, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Check permissions
	if userContext != nil && !account.HasEC2ViewPermission(userContext.Groups) {
		return nil, fmt.Errorf("user does not have permission to view account summary")
	}

	// Get all instances
	instanceList, err := ac.instanceRepo.ListInstances(ctx, accountKey, region, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances for summary: %w", err)
	}

	// Calculate summary
	summary := ac.processor.CalculateAccountSummary(account, instanceList.Instances)
	return summary, nil
}

// GetCostSavings analyzes potential cost savings
func (ac *AWSController) GetCostSavings(ctx context.Context, accountKey, region string, userContext *awsPorts.UserContext) (*awsModels.CostSavings, error) {
	// Get instances
	instanceList, err := ac.ListInstances(ctx, accountKey, region, nil, userContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	// Calculate cost savings
	savings := ac.costCalculator.CalculateCostSavings(instanceList.Instances)
	return &savings, nil
}

// GetInstanceMetrics gets CloudWatch metrics for an instance
func (ac *AWSController) GetInstanceMetrics(ctx context.Context, accountKey, region, instanceID, period string, userContext *awsPorts.UserContext) (*awsModels.InstanceMetrics, error) {
	// Validate permissions
	_, err := ac.GetInstance(ctx, accountKey, region, instanceID, userContext)
	if err != nil {
		return nil, fmt.Errorf("permission or instance validation failed: %w", err)
	}

	// Get metrics
	if ac.metricsRepo == nil {
		return nil, fmt.Errorf("metrics repository not available")
	}

	metrics, err := ac.metricsRepo.GetInstanceMetrics(ctx, accountKey, region, instanceID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance metrics: %w", err)
	}

	return metrics, nil
}

// EstimateOperationCost estimates cost impact of an operation
func (ac *AWSController) EstimateOperationCost(ctx context.Context, accountKey, region, instanceID, operation string, userContext *awsPorts.UserContext) (*awsLogic.OperationCostEstimate, error) {
	// Get instance
	instance, err := ac.GetInstance(ctx, accountKey, region, instanceID, userContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	// Calculate cost estimate
	estimate := ac.costCalculator.EstimateOperationCost(instance, operation)
	return &estimate, nil
}
