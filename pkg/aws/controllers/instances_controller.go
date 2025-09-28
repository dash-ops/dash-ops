package aws

import (
	"context"
	"fmt"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
	awsRepositories "github.com/dash-ops/dash-ops/pkg/aws/repositories"
)

// InstancesController orchestrates EC2 instance operations
type InstancesController struct {
	instanceRepo *awsRepositories.InstanceRepository
	accounts     []awsModels.AWSAccount
}

func NewInstancesController(instanceRepo *awsRepositories.InstanceRepository, accounts []awsModels.AWSAccount) *InstancesController {
	return &InstancesController{
		instanceRepo: instanceRepo,
		accounts:     accounts,
	}
}

func (c *InstancesController) ListInstances(ctx context.Context, accountKey, region string, filter *awsModels.InstanceFilter) (*awsModels.InstanceList, error) {
	account, err := c.getAccount(accountKey)
	if err != nil {
		return nil, err
	}
	return c.instanceRepo.ListInstances(ctx, account, region, filter)
}

func (c *InstancesController) GetInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.EC2Instance, error) {
	account, err := c.getAccount(accountKey)
	if err != nil {
		return nil, err
	}
	return c.instanceRepo.GetInstance(ctx, account, region, instanceID)
}

func (c *InstancesController) StartInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceOperation, error) {
	account, err := c.getAccount(accountKey)
	if err != nil {
		return nil, err
	}
	return c.instanceRepo.StartInstance(ctx, account, region, instanceID)
}

func (c *InstancesController) StopInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceOperation, error) {
	account, err := c.getAccount(accountKey)
	if err != nil {
		return nil, err
	}
	return c.instanceRepo.StopInstance(ctx, account, region, instanceID)
}

func (c *InstancesController) RestartInstance(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceOperation, error) {
	account, err := c.getAccount(accountKey)
	if err != nil {
		return nil, err
	}
	return c.instanceRepo.RestartInstance(ctx, account, region, instanceID)
}

func (c *InstancesController) GetInstanceStatus(ctx context.Context, accountKey, region, instanceID string) (*awsModels.InstanceState, error) {
	account, err := c.getAccount(accountKey)
	if err != nil {
		return nil, err
	}
	return c.instanceRepo.GetInstanceStatus(ctx, account, region, instanceID)
}

func (c *InstancesController) BatchOperation(ctx context.Context, accountKey, region string, operation string, instanceIDs []string) (*awsModels.BatchOperation, error) {
	account, err := c.getAccount(accountKey)
	if err != nil {
		return nil, err
	}
	return c.instanceRepo.BatchOperation(ctx, account, region, operation, instanceIDs)
}

// getAccount finds an account by key
func (c *InstancesController) getAccount(accountKey string) (*awsModels.AWSAccount, error) {
	if accountKey == "" {
		return nil, fmt.Errorf("account key is required")
	}

	for _, account := range c.accounts {
		if account.Key == accountKey {
			return &account, nil
		}
	}

	return nil, fmt.Errorf("account %s not found", accountKey)
}
