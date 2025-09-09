package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

func TestCostCalculator_CalculateInstanceMonthlyCost(t *testing.T) {
	calculator := NewCostCalculator()

	tests := []struct {
		name     string
		instance *awsModels.EC2Instance
		expected float64
	}{
		{
			name:     "nil instance",
			instance: nil,
			expected: 0,
		},
		{
			name: "running t2.micro",
			instance: &awsModels.EC2Instance{
				InstanceID:   "i-123",
				InstanceType: "t2.micro",
				State:        awsModels.InstanceStateRunning,
			},
			expected: 8.468,
		},
		{
			name: "stopped t2.micro",
			instance: &awsModels.EC2Instance{
				InstanceID:   "i-123",
				InstanceType: "t2.micro",
				State:        awsModels.InstanceStateStopped,
			},
			expected: 0, // Stopped instances don't incur compute costs
		},
		{
			name: "unknown instance type",
			instance: &awsModels.EC2Instance{
				InstanceID:   "i-123",
				InstanceType: "unknown.type",
				State:        awsModels.InstanceStateRunning,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateInstanceMonthlyCost(tt.instance)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCostCalculator_CalculateAccountMonthlyCost(t *testing.T) {
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
	result := calculator.CalculateAccountMonthlyCost(instances)
	assert.Equal(t, expected, result)
}

func TestCostCalculator_CalculateCostSavings(t *testing.T) {
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

	savings := calculator.CalculateCostSavings(instances)

	expectedCurrent := 8.468 + 16.79 // Only running instances
	expectedSavings := 8.468 + 16.79 // Both running instances can be stopped

	assert.Equal(t, expectedCurrent, savings.CurrentMonthlyCost)
	assert.Equal(t, expectedSavings, savings.PotentialSavings)
	assert.Equal(t, 2, savings.StoppableInstances)
	assert.Equal(t, 100.0, savings.SavingsPercentage) // 100% savings possible
}

func TestCostCalculator_EstimateOperationCost(t *testing.T) {
	calculator := NewCostCalculator()

	tests := []struct {
		name           string
		instance       *awsModels.EC2Instance
		operation      string
		expectedImpact float64
		expectedType   string
	}{
		{
			name: "start stopped instance",
			instance: &awsModels.EC2Instance{
				InstanceID:   "i-123",
				InstanceType: "t2.micro",
				State:        awsModels.InstanceStateStopped,
			},
			operation:      "start",
			expectedImpact: 0.0116, // Hourly rate
			expectedType:   "increase",
		},
		{
			name: "stop running instance",
			instance: &awsModels.EC2Instance{
				InstanceID:   "i-123",
				InstanceType: "t2.micro",
				State:        awsModels.InstanceStateRunning,
			},
			operation:      "stop",
			expectedImpact: -0.0116, // Negative hourly rate
			expectedType:   "decrease",
		},
		{
			name: "start already running instance",
			instance: &awsModels.EC2Instance{
				InstanceID:   "i-123",
				InstanceType: "t2.micro",
				State:        awsModels.InstanceStateRunning,
			},
			operation:      "start",
			expectedImpact: 0, // No impact
			expectedType:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate := calculator.EstimateOperationCost(tt.instance, tt.operation)

			assert.Equal(t, tt.instance.InstanceID, estimate.InstanceID)
			assert.Equal(t, tt.operation, estimate.Operation)
			assert.Equal(t, tt.expectedImpact, estimate.CostImpact)
			assert.Equal(t, tt.expectedType, estimate.ImpactType)
		})
	}
}

func TestInstanceProcessor_ProcessInstanceList(t *testing.T) {
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

	tests := []struct {
		name          string
		filter        *awsModels.InstanceFilter
		expectedCount int
	}{
		{
			name:          "no filter",
			filter:        nil,
			expectedCount: 2, // i-1 skipped due to monitoring tag
		},
		{
			name: "filter by state",
			filter: &awsModels.InstanceFilter{
				State: "running",
			},
			expectedCount: 1, // Only i-2 (i-1 skipped)
		},
		{
			name: "search by name",
			filter: &awsModels.InstanceFilter{
				Search: "api",
			},
			expectedCount: 1, // Only i-2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ProcessInstanceList(instances, tt.filter, skipList)
			assert.Equal(t, tt.expectedCount, len(result.Instances))
		})
	}
}
