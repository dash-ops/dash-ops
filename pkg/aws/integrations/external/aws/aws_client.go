package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// AWSClient handles communication with AWS APIs
type AWSClient struct {
	ec2Clients        map[string]ec2iface.EC2API
	cloudWatchClients map[string]*cloudwatch.CloudWatch
}

// NewAWSClient creates a new AWS client
func NewAWSClient() *AWSClient {
	return &AWSClient{
		ec2Clients:        make(map[string]ec2iface.EC2API),
		cloudWatchClients: make(map[string]*cloudwatch.CloudWatch),
	}
}

// GetEC2Client gets an EC2 client for a specific account
func (c *AWSClient) GetEC2Client(account *awsModels.AWSAccount) (ec2iface.EC2API, error) {
	if account == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}

	// Check if client already exists
	if client, exists := c.ec2Clients[account.Key]; exists {
		return client, nil
	}

	// Create new EC2 client
	client, err := c.createEC2Client(account)
	if err != nil {
		return nil, fmt.Errorf("failed to create EC2 client: %w", err)
	}

	// Cache the client
	c.ec2Clients[account.Key] = client
	return client, nil
}

// GetCloudWatchClient gets a CloudWatch client for a specific account
func (c *AWSClient) GetCloudWatchClient(account *awsModels.AWSAccount) (*cloudwatch.CloudWatch, error) {
	if account == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}

	// Check if client already exists
	if client, exists := c.cloudWatchClients[account.Key]; exists {
		return client, nil
	}

	// Create new CloudWatch client
	client, err := c.createCloudWatchClient(account)
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudWatch client: %w", err)
	}

	// Cache the client
	c.cloudWatchClients[account.Key] = client
	return client, nil
}

// ValidateCredentials validates AWS credentials
func (c *AWSClient) ValidateCredentials(account *awsModels.AWSAccount) error {
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

// createEC2Client creates an EC2 client for the given account
func (c *AWSClient) createEC2Client(account *awsModels.AWSAccount) (ec2iface.EC2API, error) {
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

	return ec2.New(awsSession), nil
}

// createCloudWatchClient creates a CloudWatch client for the given account
func (c *AWSClient) createCloudWatchClient(account *awsModels.AWSAccount) (*cloudwatch.CloudWatch, error) {
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

	return cloudwatch.New(awsSession), nil
}
