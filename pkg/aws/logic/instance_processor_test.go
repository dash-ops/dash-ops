package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

func TestInstanceProcessor_ValidateInstanceOperation_WithStartStoppedInstance_ReturnsNoError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instance := &awsModels.EC2Instance{
		InstanceID: "i-123456789",
		State:      awsModels.InstanceStateStopped,
	}
	operation := "start"

	// Act
	err := processor.ValidateInstanceOperation(instance, operation)

	// Assert
	assert.NoError(t, err)
}

func TestInstanceProcessor_ValidateInstanceOperation_WithStartRunningInstance_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instance := &awsModels.EC2Instance{
		InstanceID: "i-123456789",
		State:      awsModels.InstanceStateRunning,
	}
	operation := "start"

	// Act
	err := processor.ValidateInstanceOperation(instance, operation)

	// Assert
	assert.Error(t, err)
}

func TestInstanceProcessor_ValidateInstanceOperation_WithStopRunningInstance_ReturnsNoError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instance := &awsModels.EC2Instance{
		InstanceID: "i-123456789",
		State:      awsModels.InstanceStateRunning,
	}
	operation := "stop"

	// Act
	err := processor.ValidateInstanceOperation(instance, operation)

	// Assert
	assert.NoError(t, err)
}

func TestInstanceProcessor_ValidateInstanceOperation_WithStopStoppedInstance_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instance := &awsModels.EC2Instance{
		InstanceID: "i-123456789",
		State:      awsModels.InstanceStateStopped,
	}
	operation := "stop"

	// Act
	err := processor.ValidateInstanceOperation(instance, operation)

	// Assert
	assert.Error(t, err)
}

func TestInstanceProcessor_ValidateInstanceOperation_WithNilInstance_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	operation := "start"

	// Act
	err := processor.ValidateInstanceOperation(nil, operation)

	// Assert
	assert.Error(t, err)
}

func TestInstanceProcessor_ValidateInstanceOperation_WithInvalidOperation_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instance := &awsModels.EC2Instance{
		InstanceID: "i-123456789",
		State:      awsModels.InstanceStateRunning,
	}
	operation := "invalid"

	// Act
	err := processor.ValidateInstanceOperation(instance, operation)

	// Assert
	assert.Error(t, err)
}

func TestInstanceProcessor_FilterInstancesByPermissions_WithAdminStart_ReturnsAllInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{InstanceID: "i-1", State: awsModels.InstanceStateRunning},
		{InstanceID: "i-2", State: awsModels.InstanceStateStopped},
	}
	account := &awsModels.AWSAccount{
		Name: "test-account",
		Permissions: awsModels.AccountPermissions{
			EC2: awsModels.EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
				Stop:  []string{"admin"},
				View:  []string{"admin", "viewers", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"admin"}
	operation := "start"

	// Act
	result := processor.FilterInstancesByPermissions(instances, account, userGroups, operation)

	// Assert
	assert.Equal(t, 2, len(result))
}

func TestInstanceProcessor_FilterInstancesByPermissions_WithEC2OperatorsStart_ReturnsAllInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{InstanceID: "i-1", State: awsModels.InstanceStateRunning},
		{InstanceID: "i-2", State: awsModels.InstanceStateStopped},
	}
	account := &awsModels.AWSAccount{
		Name: "test-account",
		Permissions: awsModels.AccountPermissions{
			EC2: awsModels.EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
				Stop:  []string{"admin"},
				View:  []string{"admin", "viewers", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"ec2-operators"}
	operation := "start"

	// Act
	result := processor.FilterInstancesByPermissions(instances, account, userGroups, operation)

	// Assert
	assert.Equal(t, 2, len(result))
}

func TestInstanceProcessor_FilterInstancesByPermissions_WithViewersStart_ReturnsNoInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{InstanceID: "i-1", State: awsModels.InstanceStateRunning},
		{InstanceID: "i-2", State: awsModels.InstanceStateStopped},
	}
	account := &awsModels.AWSAccount{
		Name: "test-account",
		Permissions: awsModels.AccountPermissions{
			EC2: awsModels.EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
				Stop:  []string{"admin"},
				View:  []string{"admin", "viewers", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"viewers"}
	operation := "start"

	// Act
	result := processor.FilterInstancesByPermissions(instances, account, userGroups, operation)

	// Assert
	assert.Equal(t, 0, len(result))
}

func TestInstanceProcessor_FilterInstancesByPermissions_WithAdminStop_ReturnsAllInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{InstanceID: "i-1", State: awsModels.InstanceStateRunning},
		{InstanceID: "i-2", State: awsModels.InstanceStateStopped},
	}
	account := &awsModels.AWSAccount{
		Name: "test-account",
		Permissions: awsModels.AccountPermissions{
			EC2: awsModels.EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
				Stop:  []string{"admin"},
				View:  []string{"admin", "viewers", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"admin"}
	operation := "stop"

	// Act
	result := processor.FilterInstancesByPermissions(instances, account, userGroups, operation)

	// Assert
	assert.Equal(t, 2, len(result))
}

func TestInstanceProcessor_FilterInstancesByPermissions_WithEC2OperatorsStop_ReturnsNoInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{InstanceID: "i-1", State: awsModels.InstanceStateRunning},
		{InstanceID: "i-2", State: awsModels.InstanceStateStopped},
	}
	account := &awsModels.AWSAccount{
		Name: "test-account",
		Permissions: awsModels.AccountPermissions{
			EC2: awsModels.EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
				Stop:  []string{"admin"},
				View:  []string{"admin", "viewers", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"ec2-operators"}
	operation := "stop"

	// Act
	result := processor.FilterInstancesByPermissions(instances, account, userGroups, operation)

	// Assert
	assert.Equal(t, 0, len(result))
}

func TestInstanceProcessor_FilterInstancesByPermissions_WithViewersView_ReturnsAllInstances(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instances := []awsModels.EC2Instance{
		{InstanceID: "i-1", State: awsModels.InstanceStateRunning},
		{InstanceID: "i-2", State: awsModels.InstanceStateStopped},
	}
	account := &awsModels.AWSAccount{
		Name: "test-account",
		Permissions: awsModels.AccountPermissions{
			EC2: awsModels.EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
				Stop:  []string{"admin"},
				View:  []string{"admin", "viewers", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"viewers"}
	operation := "view"

	// Act
	result := processor.FilterInstancesByPermissions(instances, account, userGroups, operation)

	// Assert
	assert.Equal(t, 2, len(result))
}

func TestInstanceProcessor_ValidateInstanceID_WithValidNewFormat_ReturnsNoError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceID := "i-1234567890abcdef0"

	// Act
	err := processor.ValidateInstanceID(instanceID)

	// Assert
	assert.NoError(t, err)
}

func TestInstanceProcessor_ValidateInstanceID_WithValidOldFormat_ReturnsNoError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceID := "i-12345678"

	// Act
	err := processor.ValidateInstanceID(instanceID)

	// Assert
	assert.NoError(t, err)
}

func TestInstanceProcessor_ValidateInstanceID_WithEmptyInstanceID_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceID := ""

	// Act
	err := processor.ValidateInstanceID(instanceID)

	// Assert
	assert.Error(t, err)
}

func TestInstanceProcessor_ValidateInstanceID_WithInvalidPrefix_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceID := "inst-1234567890abcdef0"

	// Act
	err := processor.ValidateInstanceID(instanceID)

	// Assert
	assert.Error(t, err)
}

func TestInstanceProcessor_ValidateInstanceID_WithInvalidLength_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceID := "i-123"

	// Act
	err := processor.ValidateInstanceID(instanceID)

	// Assert
	assert.Error(t, err)
}

func TestInstanceProcessor_NormalizeInstanceType_WithUppercase_ReturnsLowercase(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceType := "T2.MICRO"

	// Act
	result := processor.NormalizeInstanceType(instanceType)

	// Assert
	assert.Equal(t, "t2.micro", result)
}

func TestInstanceProcessor_NormalizeInstanceType_WithMixedCase_ReturnsLowercase(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceType := "T3.Small"

	// Act
	result := processor.NormalizeInstanceType(instanceType)

	// Assert
	assert.Equal(t, "t3.small", result)
}

func TestInstanceProcessor_NormalizeInstanceType_WithWhitespace_ReturnsTrimmed(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceType := " m5.large "

	// Act
	result := processor.NormalizeInstanceType(instanceType)

	// Assert
	assert.Equal(t, "m5.large", result)
}

func TestInstanceProcessor_NormalizeInstanceType_WithAlreadyNormalized_ReturnsSame(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	instanceType := "t2.micro"

	// Act
	result := processor.NormalizeInstanceType(instanceType)

	// Assert
	assert.Equal(t, "t2.micro", result)
}

func TestInstanceProcessor_NormalizeRegion_WithUppercase_ReturnsLowercase(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	region := "US-EAST-1"

	// Act
	result := processor.NormalizeRegion(region)

	// Assert
	assert.Equal(t, "us-east-1", result)
}

func TestInstanceProcessor_NormalizeRegion_WithWhitespace_ReturnsTrimmed(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	region := " eu-west-1 "

	// Act
	result := processor.NormalizeRegion(region)

	// Assert
	assert.Equal(t, "eu-west-1", result)
}

func TestInstanceProcessor_NormalizeRegion_WithAlreadyNormalized_ReturnsSame(t *testing.T) {
	// Arrange
	processor := NewInstanceProcessor()
	region := "us-west-2"

	// Act
	result := processor.NormalizeRegion(region)

	// Assert
	assert.Equal(t, "us-west-2", result)
}
