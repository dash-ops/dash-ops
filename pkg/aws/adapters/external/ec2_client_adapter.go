package external

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
	awsPorts "github.com/dash-ops/dash-ops/pkg/aws/ports"
)

// EC2ClientAdapter implements EC2Client interface using AWS SDK
type EC2ClientAdapter struct {
	client ec2iface.EC2API
	region string
}

// NewEC2ClientAdapter creates a new EC2 client adapter
func NewEC2ClientAdapter(account *awsModels.AWSAccount) (*EC2ClientAdapter, error) {
	if account == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}

	// Create AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(account.Region),
		Credentials: credentials.NewStaticCredentials(
			account.AccessKeyID,
			account.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &EC2ClientAdapter{
		client: ec2.New(awsSession),
		region: account.Region,
	}, nil
}

// DescribeInstances describes EC2 instances
func (eca *EC2ClientAdapter) DescribeInstances(ctx context.Context, filter *awsPorts.EC2Filter) ([]awsModels.EC2Instance, error) {
	input := &ec2.DescribeInstancesInput{}

	// Apply filters
	if filter != nil {
		if len(filter.InstanceIDs) > 0 {
			input.InstanceIds = aws.StringSlice(filter.InstanceIDs)
		}

		if len(filter.States) > 0 {
			var filters []*ec2.Filter
			filters = append(filters, &ec2.Filter{
				Name:   aws.String("instance-state-name"),
				Values: aws.StringSlice(filter.States),
			})
			input.Filters = filters
		}

		if filter.MaxResults > 0 {
			input.MaxResults = aws.Int64(int64(filter.MaxResults))
		}
	}

	result, err := eca.client.DescribeInstancesWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	var instances []awsModels.EC2Instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, *eca.convertInstance(instance))
		}
	}

	return instances, nil
}

// DescribeInstance describes a specific EC2 instance
func (eca *EC2ClientAdapter) DescribeInstance(ctx context.Context, instanceID string) (*awsModels.EC2Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	result, err := eca.client.DescribeInstancesWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance %s: %w", instanceID, err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}

	return eca.convertInstance(result.Reservations[0].Instances[0]), nil
}

// StartInstances starts EC2 instances
func (eca *EC2ClientAdapter) StartInstances(ctx context.Context, instanceIDs []string) ([]awsModels.InstanceOperation, error) {
	if len(instanceIDs) == 0 {
		return nil, fmt.Errorf("no instance IDs provided")
	}

	input := &ec2.StartInstancesInput{
		InstanceIds: aws.StringSlice(instanceIDs),
	}

	result, err := eca.client.StartInstancesWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to start instances: %w", err)
	}

	var operations []awsModels.InstanceOperation
	for _, instance := range result.StartingInstances {
		operations = append(operations, awsModels.InstanceOperation{
			InstanceID:    aws.StringValue(instance.InstanceId),
			Operation:     "start",
			CurrentState:  eca.convertInstanceState(instance.CurrentState),
			PreviousState: eca.convertInstanceState(instance.PreviousState),
			Success:       true,
			Timestamp:     time.Now(),
		})
	}

	return operations, nil
}

// StopInstances stops EC2 instances
func (eca *EC2ClientAdapter) StopInstances(ctx context.Context, instanceIDs []string) ([]awsModels.InstanceOperation, error) {
	if len(instanceIDs) == 0 {
		return nil, fmt.Errorf("no instance IDs provided")
	}

	input := &ec2.StopInstancesInput{
		InstanceIds: aws.StringSlice(instanceIDs),
	}

	result, err := eca.client.StopInstancesWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to stop instances: %w", err)
	}

	var operations []awsModels.InstanceOperation
	for _, instance := range result.StoppingInstances {
		operations = append(operations, awsModels.InstanceOperation{
			InstanceID:    aws.StringValue(instance.InstanceId),
			Operation:     "stop",
			CurrentState:  eca.convertInstanceState(instance.CurrentState),
			PreviousState: eca.convertInstanceState(instance.PreviousState),
			Success:       true,
			Timestamp:     time.Now(),
		})
	}

	return operations, nil
}

// RebootInstances reboots EC2 instances
func (eca *EC2ClientAdapter) RebootInstances(ctx context.Context, instanceIDs []string) ([]awsModels.InstanceOperation, error) {
	if len(instanceIDs) == 0 {
		return nil, fmt.Errorf("no instance IDs provided")
	}

	input := &ec2.RebootInstancesInput{
		InstanceIds: aws.StringSlice(instanceIDs),
	}

	_, err := eca.client.RebootInstancesWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to reboot instances: %w", err)
	}

	// Reboot doesn't return state changes, so we create operations manually
	var operations []awsModels.InstanceOperation
	for _, instanceID := range instanceIDs {
		operations = append(operations, awsModels.InstanceOperation{
			InstanceID: instanceID,
			Operation:  "reboot",
			Success:    true,
			Message:    "Reboot initiated",
			Timestamp:  time.Now(),
		})
	}

	return operations, nil
}

// DescribeInstanceTypes describes available instance types
func (eca *EC2ClientAdapter) DescribeInstanceTypes(ctx context.Context) ([]awsPorts.InstanceTypeInfo, error) {
	// This is a simplified implementation
	// In production, you'd use the DescribeInstanceTypes API
	return []awsPorts.InstanceTypeInfo{
		{
			InstanceType: "t2.micro",
			VCPUs:        1,
			Memory:       1.0,
			Storage: awsPorts.InstanceTypeStorage{
				Type:         "EBS",
				EBSOptimized: false,
			},
			Network: awsPorts.InstanceTypeNetwork{
				Performance: "Low to Moderate",
				IPv6Support: true,
				ENASupport:  false,
			},
		},
		{
			InstanceType: "t3.small",
			VCPUs:        2,
			Memory:       2.0,
			Storage: awsPorts.InstanceTypeStorage{
				Type:         "EBS",
				EBSOptimized: true,
			},
			Network: awsPorts.InstanceTypeNetwork{
				Performance: "Up to 5 Gigabit",
				IPv6Support: true,
				ENASupport:  true,
			},
		},
	}, nil
}

// DescribeRegions describes available regions
func (eca *EC2ClientAdapter) DescribeRegions(ctx context.Context) ([]awsPorts.RegionInfo, error) {
	result, err := eca.client.DescribeRegionsWithContext(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %w", err)
	}

	var regions []awsPorts.RegionInfo
	for _, region := range result.Regions {
		regions = append(regions, awsPorts.RegionInfo{
			RegionName: aws.StringValue(region.RegionName),
			RegionCode: aws.StringValue(region.RegionName),
			Endpoint:   aws.StringValue(region.Endpoint),
			Available:  true,
		})
	}

	return regions, nil
}

// convertInstance converts AWS EC2 instance to domain model
func (eca *EC2ClientAdapter) convertInstance(instance *ec2.Instance) *awsModels.EC2Instance {
	// Extract tags
	var tags []awsModels.Tag
	var name string
	for _, tag := range instance.Tags {
		tagKey := aws.StringValue(tag.Key)
		tagValue := aws.StringValue(tag.Value)

		tags = append(tags, awsModels.Tag{
			Key:   tagKey,
			Value: tagValue,
		})

		if tagKey == "Name" {
			name = tagValue
		}
	}

	// Extract security groups
	var securityGroups []awsModels.SecurityGroup
	for _, sg := range instance.SecurityGroups {
		securityGroups = append(securityGroups, awsModels.SecurityGroup{
			GroupID:   aws.StringValue(sg.GroupId),
			GroupName: aws.StringValue(sg.GroupName),
		})
	}

	// Get instance type info for CPU/Memory
	cpu, memory := eca.getInstanceTypeResources(aws.StringValue(instance.InstanceType))

	return &awsModels.EC2Instance{
		InstanceID:     aws.StringValue(instance.InstanceId),
		Name:           name,
		State:          eca.convertInstanceState(instance.State),
		Platform:       eca.getPlatform(instance),
		InstanceType:   aws.StringValue(instance.InstanceType),
		PublicIP:       aws.StringValue(instance.PublicIpAddress),
		PrivateIP:      aws.StringValue(instance.PrivateIpAddress),
		SubnetID:       aws.StringValue(instance.SubnetId),
		VpcID:          aws.StringValue(instance.VpcId),
		CPU:            cpu,
		Memory:         memory,
		Tags:           tags,
		LaunchTime:     aws.TimeValue(instance.LaunchTime),
		Region:         eca.region,
		SecurityGroups: securityGroups,
	}
}

// convertInstanceState converts AWS instance state to domain model
func (eca *EC2ClientAdapter) convertInstanceState(state *ec2.InstanceState) awsModels.InstanceState {
	return awsModels.InstanceState{
		Name: aws.StringValue(state.Name),
		Code: int(aws.Int64Value(state.Code)),
	}
}

// getPlatform determines instance platform
func (eca *EC2ClientAdapter) getPlatform(instance *ec2.Instance) string {
	if instance.Platform != nil {
		return aws.StringValue(instance.Platform)
	}
	return "Linux/UNIX" // Default for non-Windows instances
}

// getInstanceTypeResources returns CPU and memory info for instance type
func (eca *EC2ClientAdapter) getInstanceTypeResources(instanceType string) (awsModels.InstanceCPU, awsModels.InstanceMemory) {
	// Simplified mapping - in production, use DescribeInstanceTypes API
	resourceMap := map[string]struct {
		vcpus  int
		memory float64
	}{
		"t2.nano":   {vcpus: 1, memory: 0.5},
		"t2.micro":  {vcpus: 1, memory: 1.0},
		"t2.small":  {vcpus: 1, memory: 2.0},
		"t2.medium": {vcpus: 2, memory: 4.0},
		"t3.micro":  {vcpus: 2, memory: 1.0},
		"t3.small":  {vcpus: 2, memory: 2.0},
		"t3.medium": {vcpus: 2, memory: 4.0},
		"m5.large":  {vcpus: 2, memory: 8.0},
		"m5.xlarge": {vcpus: 4, memory: 16.0},
	}

	if resources, exists := resourceMap[instanceType]; exists {
		return awsModels.InstanceCPU{
				VCPUs: resources.vcpus,
			}, awsModels.InstanceMemory{
				TotalGB: resources.memory,
			}
	}

	// Default for unknown types
	return awsModels.InstanceCPU{VCPUs: 1}, awsModels.InstanceMemory{TotalGB: 1.0}
}
