package aws

import (
	"fmt"
	"strings"
	"time"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// InstanceProcessor handles EC2 instance processing logic
type InstanceProcessor struct{}

// NewInstanceProcessor creates a new instance processor
func NewInstanceProcessor() *InstanceProcessor {
	return &InstanceProcessor{}
}

// ProcessInstanceList processes a list of instances with filtering and enrichment
func (ip *InstanceProcessor) ProcessInstanceList(instances []awsModels.EC2Instance, filter *awsModels.InstanceFilter, skipList []string) *awsModels.InstanceList {
	// Filter out instances that should be skipped
	var filteredInstances []awsModels.EC2Instance
	for _, instance := range instances {
		if !instance.ShouldSkip(skipList) {
			filteredInstances = append(filteredInstances, instance)
		}
	}

	instanceList := &awsModels.InstanceList{
		Instances: filteredInstances,
		Total:     len(filteredInstances),
		Filter:    filter,
	}

	// Ensure Instances is never nil
	if instanceList.Instances == nil {
		instanceList.Instances = []awsModels.EC2Instance{}
	}

	if filter == nil {
		return instanceList
	}

	// Apply filters
	if filter.State != "" {
		instanceList = instanceList.FilterByState(filter.State)
	}

	if filter.InstanceType != "" {
		instanceList = instanceList.FilterByInstanceType(filter.InstanceType)
	}

	if filter.Search != "" {
		instanceList = instanceList.Search(filter.Search)
	}

	// Apply tag filters
	for _, tagFilter := range filter.Tags {
		instanceList = instanceList.FilterByTag(tagFilter.Key, tagFilter.Value)
	}

	// Apply pagination
	if filter.Limit > 0 {
		instanceList = ip.applyPagination(instanceList, filter.Limit, filter.Offset)
	}

	return instanceList
}

// ValidateInstanceOperation validates if an operation can be performed on an instance
func (ip *InstanceProcessor) ValidateInstanceOperation(instance *awsModels.EC2Instance, operation string) error {
	if instance == nil {
		return fmt.Errorf("instance cannot be nil")
	}

	switch strings.ToLower(operation) {
	case "start":
		if !instance.CanStart() {
			return fmt.Errorf("instance %s cannot be started (current state: %s)",
				instance.InstanceID, instance.State.Name)
		}
	case "stop":
		if !instance.CanStop() {
			return fmt.Errorf("instance %s cannot be stopped (current state: %s)",
				instance.InstanceID, instance.State.Name)
		}
	case "restart":
		if !instance.IsRunning() {
			return fmt.Errorf("instance %s cannot be restarted (current state: %s)",
				instance.InstanceID, instance.State.Name)
		}
	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}

	return nil
}

// ProcessBatchOperation processes a batch operation request
func (ip *InstanceProcessor) ProcessBatchOperation(instances []awsModels.EC2Instance, operation string) (*awsModels.BatchOperation, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances provided for batch operation")
	}

	batchOp := &awsModels.BatchOperation{
		Operation:  operation,
		TotalCount: len(instances),
		StartedAt:  time.Now(),
		Results:    []awsModels.InstanceOperation{},
	}

	// Extract instance IDs
	var instanceIDs []string
	for _, instance := range instances {
		instanceIDs = append(instanceIDs, instance.InstanceID)

		// Validate operation for each instance
		if err := ip.ValidateInstanceOperation(&instance, operation); err != nil {
			// Create failed operation result
			batchOp.Results = append(batchOp.Results, awsModels.InstanceOperation{
				InstanceID:    instance.InstanceID,
				Operation:     operation,
				CurrentState:  instance.State,
				PreviousState: instance.State,
				Success:       false,
				Message:       err.Error(),
				Timestamp:     time.Now(),
			})
			batchOp.FailureCount++
		}
	}

	batchOp.Instances = instanceIDs
	return batchOp, nil
}

// CalculateAccountSummary calculates account summary from instances
func (ip *InstanceProcessor) CalculateAccountSummary(account *awsModels.AWSAccount, instances []awsModels.EC2Instance) *awsModels.AccountSummary {
	summary := &awsModels.AccountSummary{
		Account: account.Name,
		Region:  account.Region,
	}

	summary.CalculateSummary(instances)
	return summary
}

// EnrichInstanceWithMetrics enriches instance with monitoring metrics
func (ip *InstanceProcessor) EnrichInstanceWithMetrics(instance *awsModels.EC2Instance, metrics *awsModels.InstanceMetrics) {
	if instance == nil || metrics == nil {
		return
	}

	// Extract CPU utilization
	for _, metric := range metrics.Metrics {
		if metric.MetricName == "CPUUtilization" && len(metric.DataPoints) > 0 {
			// Use latest data point
			latest := metric.DataPoints[len(metric.DataPoints)-1]
			instance.CPU.Utilization = latest.Value
		}
	}

	instance.Monitoring = awsModels.InstanceMonitoring{
		Enabled:     true,
		MetricsData: ip.convertMetricsToInstanceMetrics(metrics.Metrics),
		LastUpdated: metrics.LastUpdated,
	}
}

// FilterInstancesByPermissions filters instances based on user permissions
func (ip *InstanceProcessor) FilterInstancesByPermissions(instances []awsModels.EC2Instance, account *awsModels.AWSAccount, userGroups []string, operation string) []awsModels.EC2Instance {
	var permitted []awsModels.EC2Instance

	for _, instance := range instances {
		switch strings.ToLower(operation) {
		case "start":
			if account.HasEC2StartPermission(userGroups) {
				permitted = append(permitted, instance)
			}
		case "stop":
			if account.HasEC2StopPermission(userGroups) {
				permitted = append(permitted, instance)
			}
		case "view":
			if account.HasEC2ViewPermission(userGroups) {
				permitted = append(permitted, instance)
			}
		default:
			// For unknown operations, require view permission
			if account.HasEC2ViewPermission(userGroups) {
				permitted = append(permitted, instance)
			}
		}
	}

	return permitted
}

// applyPagination applies pagination to instance list
func (ip *InstanceProcessor) applyPagination(instanceList *awsModels.InstanceList, limit, offset int) *awsModels.InstanceList {
	total := len(instanceList.Instances)

	if offset >= total {
		return &awsModels.InstanceList{
			Instances: []awsModels.EC2Instance{},
			Total:     0,
			Account:   instanceList.Account,
			Region:    instanceList.Region,
			Filter:    instanceList.Filter,
		}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return &awsModels.InstanceList{
		Instances: instanceList.Instances[offset:end],
		Total:     total, // Keep original total for pagination info
		Account:   instanceList.Account,
		Region:    instanceList.Region,
		Filter:    instanceList.Filter,
	}
}

// convertMetricsToInstanceMetrics converts metrics data to instance metrics
func (ip *InstanceProcessor) convertMetricsToInstanceMetrics(metrics []awsModels.InstanceMetricData) []awsModels.InstanceMetric {
	var instanceMetrics []awsModels.InstanceMetric

	for _, metric := range metrics {
		for _, dataPoint := range metric.DataPoints {
			instanceMetrics = append(instanceMetrics, awsModels.InstanceMetric{
				MetricName: metric.MetricName,
				Value:      dataPoint.Value,
				Unit:       dataPoint.Unit,
				Timestamp:  dataPoint.Timestamp,
			})
		}
	}

	return instanceMetrics
}

// NormalizeInstanceType normalizes instance type format
func (ip *InstanceProcessor) NormalizeInstanceType(instanceType string) string {
	return strings.ToLower(strings.TrimSpace(instanceType))
}

// NormalizeRegion normalizes AWS region format
func (ip *InstanceProcessor) NormalizeRegion(region string) string {
	return strings.ToLower(strings.TrimSpace(region))
}

// ValidateInstanceID validates EC2 instance ID format
func (ip *InstanceProcessor) ValidateInstanceID(instanceID string) error {
	if instanceID == "" {
		return fmt.Errorf("instance ID cannot be empty")
	}

	if !strings.HasPrefix(instanceID, "i-") {
		return fmt.Errorf("invalid instance ID format: %s", instanceID)
	}

	if len(instanceID) != 19 && len(instanceID) != 10 { // New format: i-1234567890abcdef0, Old format: i-12345678
		return fmt.Errorf("invalid instance ID length: %s", instanceID)
	}

	return nil
}
