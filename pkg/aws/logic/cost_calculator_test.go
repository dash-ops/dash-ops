package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

func TestCostCalculator_CalculateInstanceMonthlyCost_WithNilInstance_ReturnsZero(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()

	// Act
	result := calculator.CalculateInstanceMonthlyCost(nil)

	// Assert
	assert.Equal(t, 0.0, result)
}

func TestCostCalculator_CalculateInstanceMonthlyCost_WithRunningT2Micro_ReturnsCorrectCost(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instance := &awsModels.EC2Instance{
		InstanceID:   "i-123",
		InstanceType: "t2.micro",
		State:        awsModels.InstanceStateRunning,
	}

	// Act
	result := calculator.CalculateInstanceMonthlyCost(instance)

	// Assert
	assert.Equal(t, 8.468, result)
}

func TestCostCalculator_CalculateInstanceMonthlyCost_WithStoppedT2Micro_ReturnsZero(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instance := &awsModels.EC2Instance{
		InstanceID:   "i-123",
		InstanceType: "t2.micro",
		State:        awsModels.InstanceStateStopped,
	}

	// Act
	result := calculator.CalculateInstanceMonthlyCost(instance)

	// Assert
	assert.Equal(t, 0.0, result) // Stopped instances don't incur compute costs
}

func TestCostCalculator_CalculateInstanceMonthlyCost_WithUnknownInstanceType_ReturnsZero(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instance := &awsModels.EC2Instance{
		InstanceID:   "i-123",
		InstanceType: "unknown.type",
		State:        awsModels.InstanceStateRunning,
	}

	// Act
	result := calculator.CalculateInstanceMonthlyCost(instance)

	// Assert
	assert.Equal(t, 0.0, result)
}

func TestCostCalculator_CalculateAccountMonthlyCost_WithMixedInstances_ReturnsCorrectTotal(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instances := []awsModels.EC2Instance{
		{
			InstanceID:   "i-1",
			InstanceType: "t2.micro",
			State:        awsModels.InstanceStateRunning,
		},
		{
			InstanceID:   "i-2",
			InstanceType: "t2.small",
			State:        awsModels.InstanceStateRunning,
		},
		{
			InstanceID:   "i-3",
			InstanceType: "t2.micro",
			State:        awsModels.InstanceStateStopped, // Should not count
		},
	}
	expected := 8.468 + 16.79 // Only running instances

	// Act
	result := calculator.CalculateAccountMonthlyCost(instances)

	// Assert
	assert.Equal(t, expected, result)
}

func TestCostCalculator_CalculateCostSavings_WithMixedInstances_ReturnsCorrectSavings(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instances := []awsModels.EC2Instance{
		{
			InstanceID:   "i-1",
			InstanceType: "t2.micro",
			State:        awsModels.InstanceStateRunning, // Can be stopped
		},
		{
			InstanceID:   "i-2",
			InstanceType: "t2.small",
			State:        awsModels.InstanceStateRunning, // Can be stopped
		},
		{
			InstanceID:   "i-3",
			InstanceType: "t2.micro",
			State:        awsModels.InstanceStateStopped, // Already stopped
		},
	}
	expectedCurrent := 8.468 + 16.79 // Only running instances
	expectedSavings := 8.468 + 16.79 // Both running instances can be stopped

	// Act
	savings := calculator.CalculateCostSavings(instances)

	// Assert
	assert.Equal(t, expectedCurrent, savings.CurrentMonthlyCost)
	assert.Equal(t, expectedSavings, savings.PotentialSavings)
	assert.Equal(t, 2, savings.StoppableInstances)
	assert.Equal(t, 100.0, savings.SavingsPercentage) // 100% savings possible
}

func TestCostCalculator_EstimateOperationCost_WithStartStoppedInstance_ReturnsIncrease(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instance := &awsModels.EC2Instance{
		InstanceID:   "i-123",
		InstanceType: "t2.micro",
		State:        awsModels.InstanceStateStopped,
	}
	operation := "start"

	// Act
	estimate := calculator.EstimateOperationCost(instance, operation)

	// Assert
	assert.Equal(t, instance.InstanceID, estimate.InstanceID)
	assert.Equal(t, operation, estimate.Operation)
	assert.Equal(t, 0.0116, estimate.CostImpact) // Hourly rate
	assert.Equal(t, "increase", estimate.ImpactType)
}

func TestCostCalculator_EstimateOperationCost_WithStopRunningInstance_ReturnsDecrease(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instance := &awsModels.EC2Instance{
		InstanceID:   "i-123",
		InstanceType: "t2.micro",
		State:        awsModels.InstanceStateRunning,
	}
	operation := "stop"

	// Act
	estimate := calculator.EstimateOperationCost(instance, operation)

	// Assert
	assert.Equal(t, instance.InstanceID, estimate.InstanceID)
	assert.Equal(t, operation, estimate.Operation)
	assert.Equal(t, -0.0116, estimate.CostImpact) // Negative hourly rate
	assert.Equal(t, "decrease", estimate.ImpactType)
}

func TestCostCalculator_EstimateOperationCost_WithStartAlreadyRunningInstance_ReturnsNoImpact(t *testing.T) {
	// Arrange
	calculator := NewCostCalculator()
	instance := &awsModels.EC2Instance{
		InstanceID:   "i-123",
		InstanceType: "t2.micro",
		State:        awsModels.InstanceStateRunning,
	}
	operation := "start"

	// Act
	estimate := calculator.EstimateOperationCost(instance, operation)

	// Assert
	assert.Equal(t, instance.InstanceID, estimate.InstanceID)
	assert.Equal(t, operation, estimate.Operation)
	assert.Equal(t, 0.0, estimate.CostImpact) // No impact
	assert.Equal(t, "", estimate.ImpactType)
}

func TestInstanceProcessor_ProcessInstanceList_WithNoFilter_ReturnsFilteredInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{
			InstanceID: "i-1",
			Name:       "web-server",
			State:      awsModels.InstanceStateRunning,
			Tags: []awsModels.Tag{
				{Key: "Skip", Value: "monitoring"},
			},
		},
		{
			InstanceID: "i-2",
			Name:       "api-server",
			State:      awsModels.InstanceStateRunning,
			Tags: []awsModels.Tag{
				{Key: "Environment", Value: "production"},
			},
		},
		{
			InstanceID: "i-3",
			Name:       "test-server",
			State:      awsModels.InstanceStateStopped,
			Tags: []awsModels.Tag{
				{Key: "Skip", Value: "backup"},
			},
		},
	}
	skipList := []string{"monitoring"}

	// Act
	result := processor.ProcessInstanceList(instances, nil, skipList)

	// Assert
	assert.Equal(t, 2, len(result.Instances)) // i-1 skipped due to monitoring tag
}

func TestInstanceProcessor_ProcessInstanceList_WithStateFilter_ReturnsFilteredInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{
			InstanceID: "i-1",
			Name:       "web-server",
			State:      awsModels.InstanceStateRunning,
			Tags: []awsModels.Tag{
				{Key: "Skip", Value: "monitoring"},
			},
		},
		{
			InstanceID: "i-2",
			Name:       "api-server",
			State:      awsModels.InstanceStateRunning,
			Tags: []awsModels.Tag{
				{Key: "Environment", Value: "production"},
			},
		},
		{
			InstanceID: "i-3",
			Name:       "test-server",
			State:      awsModels.InstanceStateStopped,
			Tags: []awsModels.Tag{
				{Key: "Skip", Value: "backup"},
			},
		},
	}
	skipList := []string{"monitoring"}
	filter := &awsModels.InstanceFilter{
		State: "running",
	}

	// Act
	result := processor.ProcessInstanceList(instances, filter, skipList)

	// Assert
	assert.Equal(t, 1, len(result.Instances)) // Only i-2 (i-1 skipped)
}

func TestInstanceProcessor_ProcessInstanceList_WithSearchFilter_ReturnsFilteredInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{
			InstanceID: "i-1",
			Name:       "web-server",
			State:      awsModels.InstanceStateRunning,
			Tags: []awsModels.Tag{
				{Key: "Skip", Value: "monitoring"},
			},
		},
		{
			InstanceID: "i-2",
			Name:       "api-server",
			State:      awsModels.InstanceStateRunning,
			Tags: []awsModels.Tag{
				{Key: "Environment", Value: "production"},
			},
		},
		{
			InstanceID: "i-3",
			Name:       "test-server",
			State:      awsModels.InstanceStateStopped,
			Tags: []awsModels.Tag{
				{Key: "Skip", Value: "backup"},
			},
		},
	}
	skipList := []string{"monitoring"}
	filter := &awsModels.InstanceFilter{
		Search: "api",
	}

	// Act
	result := processor.ProcessInstanceList(instances, filter, skipList)

	// Assert
	assert.Equal(t, 1, len(result.Instances)) // Only i-2
}
