package repositories

import (
	"context"
	"fmt"
	"time"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
	awsPorts "github.com/dash-ops/dash-ops/pkg/aws/ports"
)

// InstanceRepository implements instance data access using AWS SDK
type InstanceRepository struct {
	awsClientService awsPorts.AWSClientService
}

// NewInstanceRepository creates a new instance repository
func NewInstanceRepository(awsClientService awsPorts.AWSClientService) *InstanceRepository {
	return &InstanceRepository{
		awsClientService: awsClientService,
	}
}

// GetInstance gets a specific EC2 instance
func (ir *InstanceRepository) GetInstance(ctx context.Context, account *awsModels.AWSAccount, region, instanceID string) (*awsModels.EC2Instance, error) {
	// Get AWS client
	awsClient, err := ir.awsClientService.GetEC2Client(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS client: %w", err)
	}

	// Get instance from AWS
	instance, err := awsClient.DescribeInstance(ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %w", err)
	}

	// Set account and region
	instance.Account = account.Name
	instance.Region = region

	return instance, nil
}

// ListInstances lists EC2 instances with optional filtering
func (ir *InstanceRepository) ListInstances(ctx context.Context, account *awsModels.AWSAccount, region string, filter *awsModels.InstanceFilter) (*awsModels.InstanceList, error) {
	// Get AWS client
	awsClient, err := ir.awsClientService.GetEC2Client(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS client: %w", err)
	}

	// Convert filter to AWS filter
	awsFilter := ir.convertToAWSFilter(filter)

	// Get instances from AWS
	instances, err := awsClient.DescribeInstances(ctx, awsFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	// Set account and region for all instances
	for i := range instances {
		instances[i].Account = account.Name
		instances[i].Region = region
	}

	// Create instance list
	instanceList := &awsModels.InstanceList{
		Instances: instances,
		Total:     len(instances),
		Account:   account.Name,
		Region:    region,
		Filter:    filter,
	}

	// Ensure Instances is never nil
	if instanceList.Instances == nil {
		instanceList.Instances = []awsModels.EC2Instance{}
	}

	return instanceList, nil
}

// StartInstance starts an EC2 instance
func (ir *InstanceRepository) StartInstance(ctx context.Context, account *awsModels.AWSAccount, region, instanceID string) (*awsModels.InstanceOperation, error) {
	// Get AWS client
	awsClient, err := ir.awsClientService.GetEC2Client(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS client: %w", err)
	}

	// Start instance
	operations, err := awsClient.StartInstances(ctx, []string{instanceID})
	if err != nil {
		return nil, fmt.Errorf("failed to start instance: %w", err)
	}

	if len(operations) == 0 {
		return nil, fmt.Errorf("no operation result returned")
	}

	return &operations[0], nil
}

// StopInstance stops an EC2 instance
func (ir *InstanceRepository) StopInstance(ctx context.Context, account *awsModels.AWSAccount, region, instanceID string) (*awsModels.InstanceOperation, error) {
	// Get AWS client
	awsClient, err := ir.awsClientService.GetEC2Client(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS client: %w", err)
	}

	// Stop instance
	operations, err := awsClient.StopInstances(ctx, []string{instanceID})
	if err != nil {
		return nil, fmt.Errorf("failed to stop instance: %w", err)
	}

	if len(operations) == 0 {
		return nil, fmt.Errorf("no operation result returned")
	}

	return &operations[0], nil
}

// RestartInstance restarts an EC2 instance
func (ir *InstanceRepository) RestartInstance(ctx context.Context, account *awsModels.AWSAccount, region, instanceID string) (*awsModels.InstanceOperation, error) {
	// Get AWS client
	awsClient, err := ir.awsClientService.GetEC2Client(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS client: %w", err)
	}

	// Restart instance
	operations, err := awsClient.RebootInstances(ctx, []string{instanceID})
	if err != nil {
		return nil, fmt.Errorf("failed to restart instance: %w", err)
	}

	if len(operations) == 0 {
		return nil, fmt.Errorf("no operation result returned")
	}

	return &operations[0], nil
}

// GetInstanceStatus gets current instance status
func (ir *InstanceRepository) GetInstanceStatus(ctx context.Context, account *awsModels.AWSAccount, region, instanceID string) (*awsModels.InstanceState, error) {
	// Get instance to get current state
	instance, err := ir.GetInstance(ctx, account, region, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return &instance.State, nil
}

// BatchOperation performs batch operations on multiple instances
func (ir *InstanceRepository) BatchOperation(ctx context.Context, account *awsModels.AWSAccount, region string, operation string, instanceIDs []string) (*awsModels.BatchOperation, error) {
	if len(instanceIDs) == 0 {
		return nil, fmt.Errorf("no instance IDs provided")
	}

	// Get AWS client
	awsClient, err := ir.awsClientService.GetEC2Client(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS client: %w", err)
	}

	// Create batch operation result
	batchOp := &awsModels.BatchOperation{
		Operation:    operation,
		Instances:    instanceIDs,
		Account:      account.Name,
		Region:       region,
		TotalCount:   len(instanceIDs),
		SuccessCount: 0,
		FailureCount: 0,
		StartedAt:    time.Now(),
		Results:      []awsModels.InstanceOperation{},
	}

	// Execute operation based on type
	var operations []awsModels.InstanceOperation
	switch operation {
	case "start":
		operations, err = awsClient.StartInstances(ctx, instanceIDs)
	case "stop":
		operations, err = awsClient.StopInstances(ctx, instanceIDs)
	case "restart":
		operations, err = awsClient.RebootInstances(ctx, instanceIDs)
	default:
		return nil, fmt.Errorf("unsupported batch operation: %s", operation)
	}

	if err != nil {
		// Mark all as failed
		for _, instanceID := range instanceIDs {
			batchOp.Results = append(batchOp.Results, awsModels.InstanceOperation{
				InstanceID: instanceID,
				Operation:  operation,
				Success:    false,
				Message:    err.Error(),
				Timestamp:  time.Now(),
			})
		}
		batchOp.FailureCount = len(instanceIDs)
	} else {
		// Process results
		for _, op := range operations {
			batchOp.Results = append(batchOp.Results, op)
			if op.Success {
				batchOp.SuccessCount++
			} else {
				batchOp.FailureCount++
			}
		}
	}

	batchOp.CompletedAt = time.Now()
	return batchOp, nil
}

// convertToAWSFilter converts domain filter to AWS filter
func (ir *InstanceRepository) convertToAWSFilter(filter *awsModels.InstanceFilter) *awsPorts.EC2Filter {
	if filter == nil {
		return nil
	}

	awsFilter := &awsPorts.EC2Filter{
		MaxResults: filter.Limit,
	}

	// Convert state filter
	if filter.State != "" {
		awsFilter.States = []string{filter.State}
	}

	// Convert instance type filter
	if filter.InstanceType != "" {
		awsFilter.InstanceTypes = []string{filter.InstanceType}
	}

	// Convert tag filters
	if len(filter.Tags) > 0 {
		awsFilter.Tags = make(map[string]string)
		for _, tag := range filter.Tags {
			awsFilter.Tags[tag.Key] = tag.Value
		}
	}

	return awsFilter
}
