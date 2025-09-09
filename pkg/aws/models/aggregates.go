package aws

import (
	"strings"
	"time"
)

// InstanceList represents a list of EC2 instances with metadata
type InstanceList struct {
	Instances []EC2Instance   `json:"instances"`
	Total     int             `json:"total"`
	Account   string          `json:"account"`
	Region    string          `json:"region"`
	Filter    *InstanceFilter `json:"filter,omitempty"`
}

// InstanceFilter represents filtering criteria for instances
type InstanceFilter struct {
	Account      string      `json:"account,omitempty"`
	Region       string      `json:"region,omitempty"`
	State        string      `json:"state,omitempty"`
	InstanceType string      `json:"instance_type,omitempty"`
	Platform     string      `json:"platform,omitempty"`
	Tags         []TagFilter `json:"tags,omitempty"`
	Search       string      `json:"search,omitempty"`
	Limit        int         `json:"limit,omitempty"`
	Offset       int         `json:"offset,omitempty"`
}

// TagFilter represents tag-based filtering
type TagFilter struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// AccountSummary represents account resource summary
type AccountSummary struct {
	Account              string    `json:"account"`
	Region               string    `json:"region"`
	TotalInstances       int       `json:"total_instances"`
	RunningInstances     int       `json:"running_instances"`
	StoppedInstances     int       `json:"stopped_instances"`
	PendingInstances     int       `json:"pending_instances"`
	EstimatedMonthlyCost float64   `json:"estimated_monthly_cost"`
	LastUpdated          time.Time `json:"last_updated"`
}

// InstanceMetrics represents instance monitoring metrics
type InstanceMetrics struct {
	InstanceID  string               `json:"instance_id"`
	Account     string               `json:"account"`
	Region      string               `json:"region"`
	Metrics     []InstanceMetricData `json:"metrics"`
	Period      string               `json:"period"` // 5m, 1h, 1d, etc.
	LastUpdated time.Time            `json:"last_updated"`
}

// InstanceMetricData represents metric data points
type InstanceMetricData struct {
	MetricName string            `json:"metric_name"`
	Unit       string            `json:"unit"`
	DataPoints []MetricDataPoint `json:"data_points"`
}

// MetricDataPoint represents a single metric data point
type MetricDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
}

// BatchOperation represents a batch operation on multiple instances
type BatchOperation struct {
	Operation    string              `json:"operation"`
	Instances    []string            `json:"instances"`
	Account      string              `json:"account"`
	Region       string              `json:"region"`
	Results      []InstanceOperation `json:"results"`
	TotalCount   int                 `json:"total_count"`
	SuccessCount int                 `json:"success_count"`
	FailureCount int                 `json:"failure_count"`
	StartedAt    time.Time           `json:"started_at"`
	CompletedAt  time.Time           `json:"completed_at,omitempty"`
}

// Methods for AWSAccount

// GetDisplayName returns account display name
func (acc *AWSAccount) GetDisplayName() string {
	if acc.Name != "" {
		return acc.Name
	}
	return acc.Key
}

// Methods for InstanceList

// FilterByState filters instances by state
func (il *InstanceList) FilterByState(state string) *InstanceList {
	if state == "" {
		return il
	}

	var filtered []EC2Instance
	for _, instance := range il.Instances {
		if strings.EqualFold(instance.State.Name, state) {
			filtered = append(filtered, instance)
		}
	}

	return &InstanceList{
		Instances: filtered,
		Total:     len(filtered),
		Account:   il.Account,
		Region:    il.Region,
		Filter:    &InstanceFilter{State: state},
	}
}

// FilterByInstanceType filters instances by instance type
func (il *InstanceList) FilterByInstanceType(instanceType string) *InstanceList {
	if instanceType == "" {
		return il
	}

	var filtered []EC2Instance
	for _, instance := range il.Instances {
		if strings.EqualFold(instance.InstanceType, instanceType) {
			filtered = append(filtered, instance)
		}
	}

	return &InstanceList{
		Instances: filtered,
		Total:     len(filtered),
		Account:   il.Account,
		Region:    il.Region,
		Filter:    &InstanceFilter{InstanceType: instanceType},
	}
}

// FilterByTag filters instances by tag
func (il *InstanceList) FilterByTag(key, value string) *InstanceList {
	if key == "" {
		return il
	}

	var filtered []EC2Instance
	for _, instance := range il.Instances {
		if value == "" {
			// Just check if tag key exists
			if instance.GetTag(key) != "" {
				filtered = append(filtered, instance)
			}
		} else {
			// Check for exact key-value match
			if instance.HasTag(key, value) {
				filtered = append(filtered, instance)
			}
		}
	}

	return &InstanceList{
		Instances: filtered,
		Total:     len(filtered),
		Account:   il.Account,
		Region:    il.Region,
		Filter:    &InstanceFilter{Tags: []TagFilter{{Key: key, Value: value}}},
	}
}

// Search filters instances by text search
func (il *InstanceList) Search(query string) *InstanceList {
	if query == "" {
		return il
	}

	query = strings.ToLower(query)
	var filtered []EC2Instance

	for _, instance := range il.Instances {
		if il.matchesSearch(instance, query) {
			filtered = append(filtered, instance)
		}
	}

	return &InstanceList{
		Instances: filtered,
		Total:     len(filtered),
		Account:   il.Account,
		Region:    il.Region,
		Filter:    &InstanceFilter{Search: query},
	}
}

// matchesSearch checks if instance matches search query
func (il *InstanceList) matchesSearch(instance EC2Instance, query string) bool {
	searchFields := []string{
		strings.ToLower(instance.InstanceID),
		strings.ToLower(instance.Name),
		strings.ToLower(instance.InstanceType),
		strings.ToLower(instance.Platform),
		strings.ToLower(instance.PublicIP),
		strings.ToLower(instance.PrivateIP),
	}

	// Search in tags
	for _, tag := range instance.Tags {
		searchFields = append(searchFields,
			strings.ToLower(tag.Key),
			strings.ToLower(tag.Value))
	}

	for _, field := range searchFields {
		if strings.Contains(field, query) {
			return true
		}
	}

	return false
}

// GetRunningInstances returns only running instances
func (il *InstanceList) GetRunningInstances() []EC2Instance {
	var running []EC2Instance
	for _, instance := range il.Instances {
		if instance.IsRunning() {
			running = append(running, instance)
		}
	}
	return running
}

// GetStoppedInstances returns only stopped instances
func (il *InstanceList) GetStoppedInstances() []EC2Instance {
	var stopped []EC2Instance
	for _, instance := range il.Instances {
		if instance.IsStopped() {
			stopped = append(stopped, instance)
		}
	}
	return stopped
}

// GetTransitioningInstances returns instances in transitioning states
func (il *InstanceList) GetTransitioningInstances() []EC2Instance {
	var transitioning []EC2Instance
	for _, instance := range il.Instances {
		if instance.IsTransitioning() {
			transitioning = append(transitioning, instance)
		}
	}
	return transitioning
}

// CalculateEstimatedCost calculates total estimated monthly cost
func (il *InstanceList) CalculateEstimatedCost() float64 {
	var total float64
	for _, instance := range il.Instances {
		if instance.IsRunning() { // Only count running instances
			total += instance.GetCostEstimate()
		}
	}
	return total
}

// Methods for AccountSummary

// CalculateSummary calculates summary from instance list
func (as *AccountSummary) CalculateSummary(instances []EC2Instance) {
	as.TotalInstances = len(instances)
	as.RunningInstances = 0
	as.StoppedInstances = 0
	as.PendingInstances = 0
	as.EstimatedMonthlyCost = 0

	for _, instance := range instances {
		switch {
		case instance.IsRunning():
			as.RunningInstances++
			as.EstimatedMonthlyCost += instance.GetCostEstimate()
		case instance.IsStopped():
			as.StoppedInstances++
		case instance.IsTransitioning():
			as.PendingInstances++
		}
	}

	as.LastUpdated = time.Now()
}

// Methods for BatchOperation

// IsCompleted checks if batch operation is completed
func (bo *BatchOperation) IsCompleted() bool {
	return !bo.CompletedAt.IsZero()
}

// GetSuccessRate returns success rate percentage
func (bo *BatchOperation) GetSuccessRate() float64 {
	if bo.TotalCount == 0 {
		return 0
	}
	return float64(bo.SuccessCount) / float64(bo.TotalCount) * 100
}

// GetDuration returns operation duration
func (bo *BatchOperation) GetDuration() time.Duration {
	if bo.CompletedAt.IsZero() {
		return time.Since(bo.StartedAt)
	}
	return bo.CompletedAt.Sub(bo.StartedAt)
}

// CostSavings represents potential cost savings analysis
type CostSavings struct {
	CurrentMonthlyCost float64   `json:"current_monthly_cost"`
	PotentialSavings   float64   `json:"potential_savings"`
	StoppableInstances int       `json:"stoppable_instances"`
	SavingsPercentage  float64   `json:"savings_percentage"`
	LastCalculated     time.Time `json:"last_calculated"`
}
