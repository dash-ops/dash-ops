package aws

import (
	"fmt"
	"strings"
	"time"
)

// AWSAccount represents an AWS account configuration
type AWSAccount struct {
	Name            string `yaml:"name" json:"name"`
	Key             string `json:"key"` // Normalized name for API usage
	Region          string `yaml:"region" json:"region"`
	AccessKeyID     string `yaml:"accessKeyId" json:"access_key_id"`
	SecretAccessKey string `yaml:"secretAccessKey" json:"-"` // Don't serialize secrets

	// Permissions and configuration
	Permissions AccountPermissions `yaml:"permission" json:"permissions"`
	EC2Config   EC2Config          `yaml:"ec2Config" json:"ec2_config"`

	// Status and metadata
	Status      AccountStatus `json:"status"`
	LastChecked time.Time     `json:"last_checked,omitempty"`
	Error       string        `json:"error,omitempty"`
}

// AccountStatus represents AWS account status
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusError    AccountStatus = "error"
	AccountStatusUnknown  AccountStatus = "unknown"
)

// AccountPermissions represents account-level permissions
type AccountPermissions struct {
	EC2 EC2Permissions `yaml:"ec2" json:"ec2"`
}

// EC2Permissions represents EC2 operation permissions
type EC2Permissions struct {
	Start []string `yaml:"start" json:"start"`
	Stop  []string `yaml:"stop" json:"stop"`
	View  []string `yaml:"view,omitempty" json:"view,omitempty"`
}

// EC2Config represents EC2-specific configuration
type EC2Config struct {
	SkipList    []string `yaml:"skipList" json:"skip_list"`
	DefaultTags []Tag    `yaml:"defaultTags,omitempty" json:"default_tags,omitempty"`
}

// EC2Instance represents an EC2 instance
type EC2Instance struct {
	InstanceID   string        `json:"instance_id"`
	Name         string        `json:"name"`
	State        InstanceState `json:"state"`
	Platform     string        `json:"platform"`
	InstanceType string        `json:"instance_type"`
	PublicIP     string        `json:"public_ip"`
	PrivateIP    string        `json:"private_ip"`
	SubnetID     string        `json:"subnet_id,omitempty"`
	VpcID        string        `json:"vpc_id,omitempty"`

	// Resource information
	CPU     InstanceCPU       `json:"cpu"`
	Memory  InstanceMemory    `json:"memory"`
	Storage []InstanceStorage `json:"storage,omitempty"`

	// Metadata
	Tags       []Tag     `json:"tags"`
	LaunchTime time.Time `json:"launch_time"`
	Account    string    `json:"account"`
	Region     string    `json:"region"`

	// Monitoring
	Monitoring     InstanceMonitoring `json:"monitoring,omitempty"`
	SecurityGroups []SecurityGroup    `json:"security_groups,omitempty"`
}

// InstanceState represents EC2 instance state
type InstanceState struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

// Common instance states
var (
	InstanceStatePending      = InstanceState{Name: "pending", Code: 0}
	InstanceStateRunning      = InstanceState{Name: "running", Code: 16}
	InstanceStateShuttingDown = InstanceState{Name: "shutting-down", Code: 32}
	InstanceStateTerminated   = InstanceState{Name: "terminated", Code: 48}
	InstanceStateStopping     = InstanceState{Name: "stopping", Code: 64}
	InstanceStateStopped      = InstanceState{Name: "stopped", Code: 80}
)

// InstanceCPU represents instance CPU information
type InstanceCPU struct {
	VCPUs          int     `json:"vcpus"`
	CoreCount      int     `json:"core_count,omitempty"`
	ThreadsPerCore int     `json:"threads_per_core,omitempty"`
	Utilization    float64 `json:"utilization,omitempty"` // Percentage
}

// InstanceMemory represents instance memory information
type InstanceMemory struct {
	TotalGB     float64 `json:"total_gb"`
	AvailableGB float64 `json:"available_gb,omitempty"`
	Utilization float64 `json:"utilization,omitempty"` // Percentage
}

// InstanceStorage represents instance storage information
type InstanceStorage struct {
	DeviceName string `json:"device_name"`
	VolumeID   string `json:"volume_id"`
	VolumeType string `json:"volume_type"`
	SizeGB     int    `json:"size_gb"`
	IOPS       int    `json:"iops,omitempty"`
	Encrypted  bool   `json:"encrypted"`
}

// InstanceMonitoring represents instance monitoring information
type InstanceMonitoring struct {
	Enabled     bool             `json:"enabled"`
	MetricsData []InstanceMetric `json:"metrics_data,omitempty"`
	LastUpdated time.Time        `json:"last_updated,omitempty"`
}

// InstanceMetric represents a single metric data point
type InstanceMetric struct {
	MetricName string    `json:"metric_name"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
}

// SecurityGroup represents an EC2 security group
type SecurityGroup struct {
	GroupID     string `json:"group_id"`
	GroupName   string `json:"group_name"`
	Description string `json:"description,omitempty"`
}

// Tag represents an AWS resource tag
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// InstanceOperation represents an instance operation result
type InstanceOperation struct {
	InstanceID    string        `json:"instance_id"`
	Operation     string        `json:"operation"`
	CurrentState  InstanceState `json:"current_state"`
	PreviousState InstanceState `json:"previous_state"`
	Success       bool          `json:"success"`
	Message       string        `json:"message,omitempty"`
	Timestamp     time.Time     `json:"timestamp"`
}

// Domain methods for AWSAccount

// Validate validates the AWS account configuration
func (acc *AWSAccount) Validate() error {
	if acc.Name == "" {
		return fmt.Errorf("account name is required")
	}

	if acc.Region == "" {
		return fmt.Errorf("region is required")
	}

	if acc.AccessKeyID == "" {
		return fmt.Errorf("access key ID is required")
	}

	if acc.SecretAccessKey == "" {
		return fmt.Errorf("secret access key is required")
	}

	return nil
}

// GenerateKey generates a normalized key from account name
func (acc *AWSAccount) GenerateKey() {
	acc.Key = strings.ToLower(strings.ReplaceAll(acc.Name, " ", "_"))
}

// IsActive checks if account is active
func (acc *AWSAccount) IsActive() bool {
	return acc.Status == AccountStatusActive
}

// HasEC2StartPermission checks if user has permission to start instances
func (acc *AWSAccount) HasEC2StartPermission(userGroups []string) bool {
	return acc.hasPermission(acc.Permissions.EC2.Start, userGroups)
}

// HasEC2StopPermission checks if user has permission to stop instances
func (acc *AWSAccount) HasEC2StopPermission(userGroups []string) bool {
	return acc.hasPermission(acc.Permissions.EC2.Stop, userGroups)
}

// HasEC2ViewPermission checks if user has permission to view instances
func (acc *AWSAccount) HasEC2ViewPermission(userGroups []string) bool {
	return acc.hasPermission(acc.Permissions.EC2.View, userGroups)
}

// hasPermission checks if user has any of the required permissions
func (acc *AWSAccount) hasPermission(requiredPerms []string, userGroups []string) bool {
	if len(requiredPerms) == 0 {
		return true // No permissions required
	}

	for _, required := range requiredPerms {
		for _, userGroup := range userGroups {
			if strings.EqualFold(required, userGroup) {
				return true
			}
		}
	}

	return false
}

// Domain methods for EC2Instance

// IsRunning checks if instance is in running state
func (inst *EC2Instance) IsRunning() bool {
	return inst.State.Name == "running"
}

// IsStopped checks if instance is in stopped state
func (inst *EC2Instance) IsStopped() bool {
	return inst.State.Name == "stopped"
}

// IsTransitioning checks if instance is in a transitioning state
func (inst *EC2Instance) IsTransitioning() bool {
	transitioning := []string{"pending", "shutting-down", "stopping", "starting"}
	for _, state := range transitioning {
		if inst.State.Name == state {
			return true
		}
	}
	return false
}

// CanStart checks if instance can be started
func (inst *EC2Instance) CanStart() bool {
	return inst.State.Name == "stopped"
}

// CanStop checks if instance can be stopped
func (inst *EC2Instance) CanStop() bool {
	return inst.State.Name == "running"
}

// GetTag returns the value of a specific tag
func (inst *EC2Instance) GetTag(key string) string {
	for _, tag := range inst.Tags {
		if strings.EqualFold(tag.Key, key) {
			return tag.Value
		}
	}
	return ""
}

// HasTag checks if instance has a specific tag
func (inst *EC2Instance) HasTag(key, value string) bool {
	for _, tag := range inst.Tags {
		if strings.EqualFold(tag.Key, key) && strings.EqualFold(tag.Value, value) {
			return true
		}
	}
	return false
}

// ShouldSkip checks if instance should be skipped based on tags
func (inst *EC2Instance) ShouldSkip(skipList []string) bool {
	for _, skipTag := range skipList {
		if inst.HasTag("Skip", skipTag) || inst.HasTag("skip", skipTag) {
			return true
		}
	}
	return false
}

// GetDisplayName returns the display name (Name tag or instance ID)
func (inst *EC2Instance) GetDisplayName() string {
	if name := inst.GetTag("Name"); name != "" {
		return name
	}
	return inst.InstanceID
}

// GetCostEstimate estimates monthly cost (simplified calculation)
func (inst *EC2Instance) GetCostEstimate() float64 {
	// This is a simplified cost estimation
	// In production, this would use AWS Pricing API
	baseCost := map[string]float64{
		"t2.micro":  8.76, // USD per month
		"t2.small":  17.52,
		"t2.medium": 35.04,
		"t3.micro":  8.40,
		"t3.small":  16.80,
		"t3.medium": 33.60,
		"m5.large":  87.60,
		"m5.xlarge": 175.20,
	}

	if cost, exists := baseCost[inst.InstanceType]; exists {
		return cost
	}

	return 0.0 // Unknown instance type
}

// Domain methods for InstanceOperation

// IsSuccessful checks if operation was successful
func (op *InstanceOperation) IsSuccessful() bool {
	return op.Success
}

// GetDuration returns operation duration (if timestamps are available)
func (op *InstanceOperation) GetDuration() time.Duration {
	// In a real implementation, we'd track start/end times
	return 0
}
