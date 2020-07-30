package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Instance Struct representing an ec2 instance
type Instance struct {
	InstanceID string `json:"instance_id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Platform   string `json:"platform"`
	PublicIP   string `json:"public_ip"`
	PrivateIP  string `json:"private_ip"`
}

// InstanceOutput Struct that represents the return of the state change of an instance
type InstanceOutput struct {
	CurrentState  string `json:"current_state"`
	PreviousState string `json:"previous_state"`
}

type instanceTags struct {
	Name string
	Skip bool
}

func (ac client) GetInstances() ([]Instance, error) {
	var instances []Instance

	params := &ec2.DescribeInstancesInput{}

	resp, err := ac.ec2.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	for _, reservation := range resp.Reservations {
		for _, inst := range reservation.Instances {
			var instance = Instance{
				InstanceID: *inst.InstanceId,
			}
			if inst.Platform != nil {
				instance.Platform = *inst.Platform
			}
			if inst.PublicIpAddress != nil {
				instance.PublicIP = *inst.PublicIpAddress
			}
			if inst.PrivateIpAddress != nil {
				instance.PrivateIP = *inst.PrivateIpAddress
			}
			if inst.State != nil {
				instance.State = *inst.State.Name
			}
			if inst.Tags != nil {
				it := getTagsInstance(inst.Tags, ac.skipList)
				instance.Name = it.Name
				if it.Skip {
					break
				}
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func getTagsInstance(tags []*ec2.Tag, skipList []string) instanceTags {
	var it instanceTags

	for _, tag := range tags {
		if *tag.Key == "Name" {
			it.Name = *tag.Value
		}
	}

	for i := range skipList {
		if it.Name == skipList[i] {
			it.Skip = true
		}
	}

	return it
}

func (ac client) StartInstance(instanceID string) (InstanceOutput, error) {
	dryRun := false
	params := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: &dryRun,
	}

	output := InstanceOutput{}
	result, err := ac.ec2.StartInstances(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return output, aerr
			}
		}
		return output, err
	}

	for _, instance := range result.StartingInstances {
		if instance.CurrentState != nil {
			output.CurrentState = *instance.CurrentState.Name
		}
		if instance.PreviousState != nil {
			output.PreviousState = *instance.PreviousState.Name
		}
	}
	return output, nil
}

func (ac client) StopInstance(instanceID string) (InstanceOutput, error) {
	dryRun := false
	params := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: &dryRun,
	}

	output := InstanceOutput{}
	result, err := ac.ec2.StopInstances(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return output, aerr
			}
		}
		return output, err
	}

	for _, instance := range result.StoppingInstances {
		if instance.CurrentState != nil {
			output.CurrentState = *instance.CurrentState.Name
		}
		if instance.PreviousState != nil {
			output.PreviousState = *instance.PreviousState.Name
		}
	}
	return output, nil
}
