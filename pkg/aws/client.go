package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Client Aws module interface
type Client interface {
	GetInstances() ([]Instance, error)
	StartInstance(instanceID string) (InstanceOutput, error)
	StopInstance(instanceID string) (InstanceOutput, error)
}

type client struct {
	ec2      *ec2.EC2
	skipList []string
}

// NewClient Create a new aws region access session
func NewClient(config config) (Client, error) {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	}))

	return client{
		ec2.New(awsSession),
		config.EC2Config.SkipList,
	}, nil
}
