package http

import (
	awsModels "github.com/dash-ops/dash-ops/pkg/aws-new/models"
	awsWire "github.com/dash-ops/dash-ops/pkg/aws-new/wire"
)

// AWSAdapter handles transformation between models and wire formats
type AWSAdapter struct{}

// NewAWSAdapter creates a new AWS adapter
func NewAWSAdapter() *AWSAdapter {
	return &AWSAdapter{}
}

// AccountToResponse converts AWSAccount model to AccountResponse
func (aa *AWSAdapter) AccountToResponse(account *awsModels.AWSAccount) awsWire.AccountResponse {
	return awsWire.AccountResponse{
		Name:   account.Name,
		Key:    account.Key,
		Region: account.Region,
		Status: string(account.Status),
		Error:  account.Error,
	}
}

// AccountsToResponse converts AWSAccount slice to AccountListResponse
func (aa *AWSAdapter) AccountsToResponse(accounts []awsModels.AWSAccount) awsWire.AccountListResponse {
	var accountResponses []awsWire.AccountResponse
	for _, account := range accounts {
		accountResponses = append(accountResponses, aa.AccountToResponse(&account))
	}

	return awsWire.AccountListResponse{
		Accounts: accountResponses,
		Total:    len(accountResponses),
	}
}

// InstanceToResponse converts EC2Instance model to InstanceResponse
func (aa *AWSAdapter) InstanceToResponse(instance *awsModels.EC2Instance) awsWire.InstanceResponse {
	var tags []awsWire.TagResponse
	for _, tag := range instance.Tags {
		tags = append(tags, awsWire.TagResponse{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	var securityGroups []awsWire.SecurityGroupResponse
	for _, sg := range instance.SecurityGroups {
		securityGroups = append(securityGroups, awsWire.SecurityGroupResponse{
			GroupID:     sg.GroupID,
			GroupName:   sg.GroupName,
			Description: sg.Description,
		})
	}

	return awsWire.InstanceResponse{
		InstanceID: instance.InstanceID,
		Name:       instance.Name,
		State: awsWire.InstanceStateResponse{
			Name: instance.State.Name,
			Code: instance.State.Code,
		},
		Platform:     instance.Platform,
		InstanceType: instance.InstanceType,
		PublicIP:     instance.PublicIP,
		PrivateIP:    instance.PrivateIP,
		SubnetID:     instance.SubnetID,
		VpcID:        instance.VpcID,
		CPU: awsWire.InstanceCPUResponse{
			VCPUs:       instance.CPU.VCPUs,
			Utilization: instance.CPU.Utilization,
		},
		Memory: awsWire.InstanceMemoryResponse{
			TotalGB:     instance.Memory.TotalGB,
			Utilization: instance.Memory.Utilization,
		},
		Tags:           tags,
		LaunchTime:     instance.LaunchTime,
		Account:        instance.Account,
		Region:         instance.Region,
		SecurityGroups: securityGroups,
		CostEstimate:   instance.GetCostEstimate(),
	}
}

// InstanceListToResponse converts InstanceList model to InstanceListResponse
func (aa *AWSAdapter) InstanceListToResponse(instanceList *awsModels.InstanceList) awsWire.InstanceListResponse {
	var instances []awsWire.InstanceResponse
	for _, instance := range instanceList.Instances {
		instances = append(instances, aa.InstanceToResponse(&instance))
	}

	return awsWire.InstanceListResponse{
		Instances: instances,
		Total:     instanceList.Total,
		Account:   instanceList.Account,
		Region:    instanceList.Region,
		Filter:    instanceList.Filter,
	}
}

// InstanceOperationToResponse converts InstanceOperation model to InstanceOperationResponse
func (aa *AWSAdapter) InstanceOperationToResponse(operation *awsModels.InstanceOperation) awsWire.InstanceOperationResponse {
	return awsWire.InstanceOperationResponse{
		InstanceID: operation.InstanceID,
		Operation:  operation.Operation,
		CurrentState: awsWire.InstanceStateResponse{
			Name: operation.CurrentState.Name,
			Code: operation.CurrentState.Code,
		},
		PreviousState: awsWire.InstanceStateResponse{
			Name: operation.PreviousState.Name,
			Code: operation.PreviousState.Code,
		},
		Success:   operation.Success,
		Message:   operation.Message,
		Timestamp: operation.Timestamp,
	}
}

// BatchOperationToResponse converts BatchOperation model to BatchOperationResponse
func (aa *AWSAdapter) BatchOperationToResponse(batchOp *awsModels.BatchOperation) awsWire.BatchOperationResponse {
	var results []awsWire.InstanceOperationResponse
	for _, result := range batchOp.Results {
		results = append(results, aa.InstanceOperationToResponse(&result))
	}

	return awsWire.BatchOperationResponse{
		Operation:    batchOp.Operation,
		TotalCount:   batchOp.TotalCount,
		SuccessCount: batchOp.SuccessCount,
		FailureCount: batchOp.FailureCount,
		Results:      results,
		StartedAt:    batchOp.StartedAt,
		CompletedAt:  batchOp.CompletedAt,
		Duration:     batchOp.GetDuration().String(),
		SuccessRate:  batchOp.GetSuccessRate(),
	}
}

// AccountSummaryToResponse converts AccountSummary model to AccountSummaryResponse
func (aa *AWSAdapter) AccountSummaryToResponse(summary *awsModels.AccountSummary) awsWire.AccountSummaryResponse {
	return awsWire.AccountSummaryResponse{
		Account:              summary.Account,
		Region:               summary.Region,
		TotalInstances:       summary.TotalInstances,
		RunningInstances:     summary.RunningInstances,
		StoppedInstances:     summary.StoppedInstances,
		PendingInstances:     summary.PendingInstances,
		EstimatedMonthlyCost: summary.EstimatedMonthlyCost,
		LastUpdated:          summary.LastUpdated,
	}
}

// CostSavingsToResponse converts CostSavings model to CostSavingsResponse
func (aa *AWSAdapter) CostSavingsToResponse(savings *awsModels.CostSavings) awsWire.CostSavingsResponse {
	return awsWire.CostSavingsResponse{
		CurrentMonthlyCost: savings.CurrentMonthlyCost,
		PotentialSavings:   savings.PotentialSavings,
		StoppableInstances: savings.StoppableInstances,
		SavingsPercentage:  savings.SavingsPercentage,
		LastCalculated:     savings.LastCalculated,
	}
}

// RequestToInstanceFilter converts wire request to InstanceFilter model
func (aa *AWSAdapter) RequestToInstanceFilter(req awsWire.InstanceFilterRequest) *awsModels.InstanceFilter {
	var tagFilters []awsModels.TagFilter
	for _, tag := range req.Tags {
		tagFilters = append(tagFilters, awsModels.TagFilter{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	return &awsModels.InstanceFilter{
		Account:      req.Account,
		Region:       req.Region,
		State:        req.State,
		InstanceType: req.InstanceType,
		Platform:     req.Platform,
		Tags:         tagFilters,
		Search:       req.Search,
		Limit:        req.Limit,
		Offset:       req.Offset,
	}
}
