package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Client Aws module interface
type Client interface {
	GetInstances() ([]Instance, error)
	StartInstance(instanceID string) (InstanceOutput, error)
	StopInstance(instanceID string) (InstanceOutput, error)
}

type client struct {
	session  *session.Session
	skipList []string
}

// NewClient Create a new aws region access session
func NewClient(config config) (Client, error) {
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return client{
		awsSession,
		config.EC2Config.SkipList,
	}, nil
}
