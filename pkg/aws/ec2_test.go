package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEc2API struct {
	ec2iface.EC2API
	mock.Mock
}

func (m mockEc2API) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called()
	return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
}

func (m mockEc2API) StartInstances(input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	args := m.Called()
	return args.Get(0).(*ec2.StartInstancesOutput), args.Error(1)
}

func (m mockEc2API) StopInstances(input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	args := m.Called()
	return args.Get(0).(*ec2.StopInstancesOutput), args.Error(1)
}

func TestGetInstances(t *testing.T) {
	mockInstances := []Instance{
		{
			InstanceID: "xpto",
			Name:       "MyInstanceEC2",
			State:      "running",
			Platform:   "",
			PublicIP:   "99.99.99.99",
			PrivateIP:  "10.10.10.10",
		},
	}

	ec2 := new(mockEc2API)
	ec2.On("DescribeInstances").Return(mockEc2DescribeInstance(mockInstances), nil)

	awsClient := client{ec2: ec2}

	instances, err := awsClient.GetInstances()

	assert.Nil(t, err)
	assert.Equal(t, mockInstances[0].Name, instances[0].Name, "return instance")
}

func TestGetInstancesWithError(t *testing.T) {
	mockErr := errors.New("message error")

	ec2 := new(mockEc2API)
	ec2.On("DescribeInstances").Return(mockEc2DescribeInstance([]Instance{}), mockErr)

	awsClient := client{ec2: ec2}

	_, err := awsClient.GetInstances()

	assert.Equal(t, mockErr, err, "return error message")
}

func TestGetInstanceWithSkippingItem(t *testing.T) {
	skipList := []string{"ScaleInstanceXpto"}

	mockInstances := []Instance{
		{
			InstanceID: "xpto1",
			Name:       "ScaleInstanceXpto",
			State:      "running",
			Platform:   "",
			PublicIP:   "99.99.99.99",
			PrivateIP:  "10.10.10.10",
		},
		{
			InstanceID: "xpto2",
			Name:       "InstanceOK",
			State:      "running",
			Platform:   "",
			PublicIP:   "99.99.99.99",
			PrivateIP:  "10.10.10.10",
		},
	}

	ec2 := new(mockEc2API)
	ec2.On("DescribeInstances").Return(mockEc2DescribeInstance(mockInstances), nil)

	awsClient := client{
		ec2:      ec2,
		skipList: skipList,
	}

	instances, err := awsClient.GetInstances()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(instances), "returns only the instance that is not in the skip list")
	assert.Equal(t, mockInstances[1].Name, instances[0].Name, "return name instance")
}

func TestStartInstance(t *testing.T) {
	mockOutput := InstanceOutput{
		CurrentState:  "pending",
		PreviousState: "stopped",
	}

	var mockCode int64 = 55
	mockStartOutput := &ec2.StartInstancesOutput{
		StartingInstances: []*ec2.InstanceStateChange{
			{
				CurrentState: &ec2.InstanceState{
					Name: &mockOutput.CurrentState,
					Code: &mockCode,
				},
			},
			{
				PreviousState: &ec2.InstanceState{
					Name: &mockOutput.PreviousState,
					Code: &mockCode,
				},
			},
		},
	}

	ec2 := new(mockEc2API)
	ec2.On("StartInstances").Return(mockStartOutput, nil)

	awsClient := client{
		ec2: ec2,
	}

	output, err := awsClient.StartInstance("5555555")

	assert.Nil(t, err)
	assert.Equal(t, mockOutput.CurrentState, output.CurrentState)
	assert.Equal(t, mockOutput.CurrentState, output.CurrentState)
}

func TestStopInstance(t *testing.T) {
	mockOutput := InstanceOutput{
		CurrentState:  "stopping",
		PreviousState: "running",
	}

	var mockCode int64 = 55
	mockStopOutput := &ec2.StopInstancesOutput{
		StoppingInstances: []*ec2.InstanceStateChange{
			{
				CurrentState: &ec2.InstanceState{
					Name: &mockOutput.CurrentState,
					Code: &mockCode,
				},
			},
			{
				PreviousState: &ec2.InstanceState{
					Name: &mockOutput.PreviousState,
					Code: &mockCode,
				},
			},
		},
	}

	ec2 := new(mockEc2API)
	ec2.On("StopInstances").Return(mockStopOutput, nil)

	awsClient := client{
		ec2: ec2,
	}

	output, err := awsClient.StopInstance("5555555")

	assert.Nil(t, err)
	assert.Equal(t, mockOutput.CurrentState, output.CurrentState)
	assert.Equal(t, mockOutput.CurrentState, output.CurrentState)
}

func mockEc2DescribeInstance(instances []Instance) *ec2.DescribeInstancesOutput {
	keyTag := "Name"

	var ec2Instances []*ec2.Instance

	for i := 0; i < len(instances); i++ {
		var state ec2.InstanceState
		state = ec2.InstanceState{
			Name: &instances[i].State,
		}

		var tag ec2.Tag
		tag = ec2.Tag{
			Key:   &keyTag,
			Value: &instances[i].Name,
		}

		var ec2Instance ec2.Instance
		ec2Instance = ec2.Instance{
			InstanceId:       &instances[i].InstanceID,
			Platform:         &instances[i].Platform,
			PublicIpAddress:  &instances[i].PublicIP,
			PrivateIpAddress: &instances[i].PrivateIP,
			State:            &state,
			Tags:             []*ec2.Tag{&tag},
		}

		ec2Instances = append(ec2Instances, &ec2Instance)
	}

	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{Instances: ec2Instances},
		},
	}
}
