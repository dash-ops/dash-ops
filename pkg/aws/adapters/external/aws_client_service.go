package external

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
	awsPorts "github.com/dash-ops/dash-ops/pkg/aws/ports"
)

// AWSClientServiceAdapter implements AWSClientService interface using AWS SDK
type AWSClientServiceAdapter struct {
	ec2Clients        map[string]awsPorts.EC2Client
	cloudWatchClients map[string]awsPorts.CloudWatchClient
}

// NewAWSClientServiceAdapter creates a new AWS client service adapter
func NewAWSClientServiceAdapter() *AWSClientServiceAdapter {
	return &AWSClientServiceAdapter{
		ec2Clients:        make(map[string]awsPorts.EC2Client),
		cloudWatchClients: make(map[string]awsPorts.CloudWatchClient),
	}
}

// GetEC2Client gets an EC2 client for a specific account
func (acsa *AWSClientServiceAdapter) GetEC2Client(account *awsModels.AWSAccount) (awsPorts.EC2Client, error) {
	if account == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}

	// Check if client already exists
	if client, exists := acsa.ec2Clients[account.Key]; exists {
		return client, nil
	}

	// Create new EC2 client
	client, err := NewEC2ClientAdapter(account)
	if err != nil {
		return nil, fmt.Errorf("failed to create EC2 client: %w", err)
	}

	// Cache the client
	acsa.ec2Clients[account.Key] = client
	return client, nil
}

// GetCloudWatchClient gets a CloudWatch client for a specific account
func (acsa *AWSClientServiceAdapter) GetCloudWatchClient(account *awsModels.AWSAccount) (awsPorts.CloudWatchClient, error) {
	if account == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}

	// Check if client already exists
	if client, exists := acsa.cloudWatchClients[account.Key]; exists {
		return client, nil
	}

	// Create new CloudWatch client
	client, err := NewCloudWatchClientAdapter(account)
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudWatch client: %w", err)
	}

	// Cache the client
	acsa.cloudWatchClients[account.Key] = client
	return client, nil
}

// ValidateCredentials validates AWS credentials
func (acsa *AWSClientServiceAdapter) ValidateCredentials(account *awsModels.AWSAccount) error {
	if account == nil {
		return fmt.Errorf("account cannot be nil")
	}

	// Create a session to validate credentials
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(account.Region),
		Credentials: credentials.NewStaticCredentials(
			account.AccessKeyID,
			account.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Try to make a simple API call to validate credentials
	ec2Client := ec2.New(awsSession)
	_, err = ec2Client.DescribeRegionsWithContext(context.Background(), &ec2.DescribeRegionsInput{})
	if err != nil {
		return fmt.Errorf("credential validation failed: %w", err)
	}

	return nil
}

// GetAccountInfo gets basic account information
func (acsa *AWSClientServiceAdapter) GetAccountInfo(account *awsModels.AWSAccount) (*awsPorts.AccountInfo, error) {
	if account == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}

	// For now, return basic info from the account model
	// In a real implementation, this would call STS GetCallerIdentity
	return &awsPorts.AccountInfo{
		AccountID:   "unknown", // Would be retrieved from STS
		Alias:       account.Name,
		Region:      account.Region,
		Status:      "active",
		LastChecked: time.Now(),
	}, nil
}

// CloudWatchClientAdapter implements CloudWatchClient interface using AWS SDK
type CloudWatchClientAdapter struct {
	client cloudwatch.CloudWatch
	region string
}

// NewCloudWatchClientAdapter creates a new CloudWatch client adapter
func NewCloudWatchClientAdapter(account *awsModels.AWSAccount) (*CloudWatchClientAdapter, error) {
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

	return &CloudWatchClientAdapter{
		client: *cloudwatch.New(awsSession),
		region: account.Region,
	}, nil
}

// GetInstanceMetrics gets CloudWatch metrics for an instance
func (cwca *CloudWatchClientAdapter) GetInstanceMetrics(ctx context.Context, instanceID string, period time.Duration, startTime, endTime time.Time) (*awsModels.InstanceMetrics, error) {
	// This is a simplified implementation
	// In production, you would call CloudWatch GetMetricStatistics API

	// For now, return mock data
	metrics := &awsModels.InstanceMetrics{
		InstanceID: instanceID,
		Region:     cwca.region,
		Period:     period.String(),
		Metrics: []awsModels.InstanceMetricData{
			{
				MetricName: "CPUUtilization",
				Unit:       "Percent",
				DataPoints: []awsModels.MetricDataPoint{
					{
						Timestamp: time.Now().Add(-1 * time.Hour),
						Value:     45.5,
						Unit:      "Percent",
					},
					{
						Timestamp: time.Now(),
						Value:     52.3,
						Unit:      "Percent",
					},
				},
			},
			{
				MetricName: "NetworkIn",
				Unit:       "Bytes",
				DataPoints: []awsModels.MetricDataPoint{
					{
						Timestamp: time.Now().Add(-1 * time.Hour),
						Value:     1024000,
						Unit:      "Bytes",
					},
					{
						Timestamp: time.Now(),
						Value:     2048000,
						Unit:      "Bytes",
					},
				},
			},
		},
		LastUpdated: time.Now(),
	}

	return metrics, nil
}

// GetMetricStatistics gets metric statistics
func (cwca *CloudWatchClientAdapter) GetMetricStatistics(ctx context.Context, metricName, namespace string, dimensions map[string]string, period time.Duration, startTime, endTime time.Time) ([]awsModels.MetricDataPoint, error) {
	// This is a simplified implementation
	// In production, you would call CloudWatch GetMetricStatistics API

	// For now, return mock data
	return []awsModels.MetricDataPoint{
		{
			Timestamp: startTime,
			Value:     50.0,
			Unit:      "Percent",
		},
		{
			Timestamp: endTime,
			Value:     55.0,
			Unit:      "Percent",
		},
	}, nil
}

// ListMetrics lists available metrics
func (cwca *CloudWatchClientAdapter) ListMetrics(ctx context.Context, namespace string) ([]awsPorts.MetricInfo, error) {
	// This is a simplified implementation
	// In production, you would call CloudWatch ListMetrics API

	// For now, return mock data
	return []awsPorts.MetricInfo{
		{
			MetricName: "CPUUtilization",
			Namespace:  namespace,
			Unit:       "Percent",
			Dimensions: []awsPorts.MetricDimension{
				{Name: "InstanceId", Value: "i-1234567890abcdef0"},
			},
		},
		{
			MetricName: "NetworkIn",
			Namespace:  namespace,
			Unit:       "Bytes",
			Dimensions: []awsPorts.MetricDimension{
				{Name: "InstanceId", Value: "i-1234567890abcdef0"},
			},
		},
	}, nil
}
