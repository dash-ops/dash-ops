package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AwsClient Aws module interface
type AwsClient interface {
	GetInstances() ([]Instance, error)
	StartInstance(instanceID string) (InstanceOutput, error)
	StopInstance(instanceID string) (InstanceOutput, error)
}

type awsClient struct {
	session   *session.Session
	blacklist []string
}

// NewAwsClient Create a new aws region access session
func NewAwsClient(config awsConfig) (AwsClient, error) {
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return awsClient{
		awsSession,
		config.EC2Config.Blacklist,
	}, nil
}
