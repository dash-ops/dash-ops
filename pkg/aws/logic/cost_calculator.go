package aws

import (
	"strings"
	"time"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// CostCalculator provides AWS cost calculation logic
type CostCalculator struct {
	pricingData map[string]InstancePricing
}

// InstancePricing represents pricing information for an instance type
type InstancePricing struct {
	InstanceType    string    `json:"instance_type"`
	HourlyRate      float64   `json:"hourly_rate"`  // USD per hour
	MonthlyRate     float64   `json:"monthly_rate"` // USD per month (730 hours)
	Region          string    `json:"region"`
	OperatingSystem string    `json:"operating_system"`
	LastUpdated     time.Time `json:"last_updated"`
}

// NewCostCalculator creates a new cost calculator with default pricing
func NewCostCalculator() *CostCalculator {
	return &CostCalculator{
		pricingData: getDefaultPricing(),
	}
}

// CalculateInstanceMonthlyCost calculates monthly cost for a single instance
func (cc *CostCalculator) CalculateInstanceMonthlyCost(instance *awsModels.EC2Instance) float64 {
	if instance == nil {
		return 0
	}

	pricing, exists := cc.pricingData[instance.InstanceType]
	if !exists {
		return 0 // Unknown instance type
	}

	// Only calculate cost for running instances
	if !instance.IsRunning() {
		return 0
	}

	return pricing.MonthlyRate
}

// CalculateInstanceHourlyCost calculates hourly cost for a single instance
func (cc *CostCalculator) CalculateInstanceHourlyCost(instance *awsModels.EC2Instance) float64 {
	if instance == nil {
		return 0
	}

	pricing, exists := cc.pricingData[instance.InstanceType]
	if !exists {
		return 0
	}

	if !instance.IsRunning() {
		return 0
	}

	return pricing.HourlyRate
}

// CalculateAccountMonthlyCost calculates total monthly cost for an account
func (cc *CostCalculator) CalculateAccountMonthlyCost(instances []awsModels.EC2Instance) float64 {
	var totalCost float64
	for _, instance := range instances {
		totalCost += cc.CalculateInstanceMonthlyCost(&instance)
	}
	return totalCost
}

// CalculateCostSavings calculates potential cost savings if instances were stopped
func (cc *CostCalculator) CalculateCostSavings(instances []awsModels.EC2Instance) awsModels.CostSavings {
	var currentCost, potentialSavings float64
	var stoppableInstances int

	for _, instance := range instances {
		cost := cc.CalculateInstanceMonthlyCost(&instance)
		currentCost += cost

		if instance.CanStop() {
			potentialSavings += cost
			stoppableInstances++
		}
	}

	return awsModels.CostSavings{
		CurrentMonthlyCost: currentCost,
		PotentialSavings:   potentialSavings,
		StoppableInstances: stoppableInstances,
		SavingsPercentage:  cc.calculatePercentage(potentialSavings, currentCost),
		LastCalculated:     time.Now(),
	}
}

// EstimateOperationCost estimates cost impact of an operation
func (cc *CostCalculator) EstimateOperationCost(instance *awsModels.EC2Instance, operation string) OperationCostEstimate {
	if instance == nil {
		return OperationCostEstimate{}
	}

	// For cost estimation, we need the potential hourly cost regardless of current state
	pricing, exists := cc.pricingData[instance.InstanceType]
	if !exists {
		return OperationCostEstimate{
			InstanceID: instance.InstanceID,
			Operation:  operation,
		}
	}

	hourlyCost := pricing.HourlyRate
	monthlyCost := pricing.MonthlyRate

	estimate := OperationCostEstimate{
		InstanceID:     instance.InstanceID,
		Operation:      operation,
		HourlyCost:     hourlyCost,
		MonthlyCost:    monthlyCost,
		LastCalculated: time.Now(),
	}

	switch strings.ToLower(operation) {
	case "start":
		if instance.IsStopped() {
			estimate.CostImpact = hourlyCost
			estimate.ImpactType = "increase"
			estimate.Description = "Starting instance will incur hourly charges"
		}
	case "stop":
		if instance.IsRunning() {
			estimate.CostImpact = -hourlyCost
			estimate.ImpactType = "decrease"
			estimate.Description = "Stopping instance will save hourly charges"
		}
	}

	return estimate
}

// GetInstanceTypePricing returns pricing information for an instance type
func (cc *CostCalculator) GetInstanceTypePricing(instanceType string) (InstancePricing, bool) {
	pricing, exists := cc.pricingData[instanceType]
	return pricing, exists
}

// UpdatePricing updates pricing data for an instance type
func (cc *CostCalculator) UpdatePricing(instanceType string, pricing InstancePricing) {
	cc.pricingData[instanceType] = pricing
}

// calculatePercentage calculates percentage safely
func (cc *CostCalculator) calculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (part / total) * 100
}

// OperationCostEstimate represents cost estimate for an operation
type OperationCostEstimate struct {
	InstanceID     string    `json:"instance_id"`
	Operation      string    `json:"operation"`
	HourlyCost     float64   `json:"hourly_cost"`
	MonthlyCost    float64   `json:"monthly_cost"`
	CostImpact     float64   `json:"cost_impact"` // Positive = increase, Negative = decrease
	ImpactType     string    `json:"impact_type"` // increase, decrease, none
	Description    string    `json:"description"`
	LastCalculated time.Time `json:"last_calculated"`
}

// getDefaultPricing returns default pricing data (US East 1, Linux)
func getDefaultPricing() map[string]InstancePricing {
	return map[string]InstancePricing{
		"t2.nano": {
			InstanceType:    "t2.nano",
			HourlyRate:      0.0058,
			MonthlyRate:     4.234,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"t2.micro": {
			InstanceType:    "t2.micro",
			HourlyRate:      0.0116,
			MonthlyRate:     8.468,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"t2.small": {
			InstanceType:    "t2.small",
			HourlyRate:      0.023,
			MonthlyRate:     16.79,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"t2.medium": {
			InstanceType:    "t2.medium",
			HourlyRate:      0.0464,
			MonthlyRate:     33.872,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"t3.micro": {
			InstanceType:    "t3.micro",
			HourlyRate:      0.0104,
			MonthlyRate:     7.592,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"t3.small": {
			InstanceType:    "t3.small",
			HourlyRate:      0.0208,
			MonthlyRate:     15.184,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"t3.medium": {
			InstanceType:    "t3.medium",
			HourlyRate:      0.0416,
			MonthlyRate:     30.368,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"m5.large": {
			InstanceType:    "m5.large",
			HourlyRate:      0.096,
			MonthlyRate:     70.08,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
		"m5.xlarge": {
			InstanceType:    "m5.xlarge",
			HourlyRate:      0.192,
			MonthlyRate:     140.16,
			Region:          "us-east-1",
			OperatingSystem: "Linux",
			LastUpdated:     time.Now(),
		},
	}
}
