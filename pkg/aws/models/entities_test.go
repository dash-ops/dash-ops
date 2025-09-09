package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAWSAccount_Validate_WithValidAccount_ReturnsNoError(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Name:            "test-account",
		Region:          "us-east-1",
		AccessKeyID:     "AKIATEST123",
		SecretAccessKey: "secret-key",
	}

	// Act
	err := account.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestAWSAccount_Validate_WithMissingName_ReturnsError(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Region:          "us-east-1",
		AccessKeyID:     "AKIATEST123",
		SecretAccessKey: "secret-key",
	}

	// Act
	err := account.Validate()

	// Assert
	assert.Error(t, err)
}

func TestAWSAccount_Validate_WithMissingRegion_ReturnsError(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Name:            "test-account",
		AccessKeyID:     "AKIATEST123",
		SecretAccessKey: "secret-key",
	}

	// Act
	err := account.Validate()

	// Assert
	assert.Error(t, err)
}

func TestAWSAccount_Validate_WithMissingAccessKey_ReturnsError(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Name:            "test-account",
		Region:          "us-east-1",
		SecretAccessKey: "secret-key",
	}

	// Act
	err := account.Validate()

	// Assert
	assert.Error(t, err)
}

func TestAWSAccount_Validate_WithMissingSecretKey_ReturnsError(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Name:        "test-account",
		Region:      "us-east-1",
		AccessKeyID: "AKIATEST123",
	}

	// Act
	err := account.Validate()

	// Assert
	assert.Error(t, err)
}

func TestAWSAccount_GenerateKey_WithSimpleName_GeneratesLowercaseKey(t *testing.T) {
	// Arrange
	account := AWSAccount{Name: "Production"}

	// Act
	account.GenerateKey()

	// Assert
	assert.Equal(t, "production", account.Key)
}

func TestAWSAccount_GenerateKey_WithNameWithSpaces_GeneratesUnderscoreKey(t *testing.T) {
	// Arrange
	account := AWSAccount{Name: "Test Account"}

	// Act
	account.GenerateKey()

	// Assert
	assert.Equal(t, "test_account", account.Key)
}

func TestAWSAccount_GenerateKey_WithNameWithSpecialChars_GeneratesNormalizedKey(t *testing.T) {
	// Arrange
	account := AWSAccount{Name: "Dev-Environment"}

	// Act
	account.GenerateKey()

	// Assert
	assert.Equal(t, "dev-environment", account.Key)
}

func TestAWSAccount_HasEC2StartPermission_WithUserHavingAdminPermission_ReturnsTrue(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Permissions: AccountPermissions{
			EC2: EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"admin", "developers"}

	// Act
	result := account.HasEC2StartPermission(userGroups)

	// Assert
	assert.True(t, result)
}

func TestAWSAccount_HasEC2StartPermission_WithUserHavingEC2OperatorsPermission_ReturnsTrue(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Permissions: AccountPermissions{
			EC2: EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"ec2-operators"}

	// Act
	result := account.HasEC2StartPermission(userGroups)

	// Assert
	assert.True(t, result)
}

func TestAWSAccount_HasEC2StartPermission_WithUserHavingNoPermission_ReturnsFalse(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Permissions: AccountPermissions{
			EC2: EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
			},
		},
	}
	userGroups := []string{"developers", "viewers"}

	// Act
	result := account.HasEC2StartPermission(userGroups)

	// Assert
	assert.False(t, result)
}

func TestAWSAccount_HasEC2StartPermission_WithEmptyUserGroups_ReturnsFalse(t *testing.T) {
	// Arrange
	account := AWSAccount{
		Permissions: AccountPermissions{
			EC2: EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
			},
		},
	}
	userGroups := []string{}

	// Act
	result := account.HasEC2StartPermission(userGroups)

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_IsRunning_WithRunningInstance_ReturnsTrue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		State: InstanceState{Name: "running", Code: 16},
	}

	// Act
	result := instance.IsRunning()

	// Assert
	assert.True(t, result)
}

func TestEC2Instance_IsRunning_WithStoppedInstance_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		State: InstanceState{Name: "stopped", Code: 80},
	}

	// Act
	result := instance.IsRunning()

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_IsRunning_WithPendingInstance_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		State: InstanceState{Name: "pending", Code: 0},
	}

	// Act
	result := instance.IsRunning()

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_IsTransitioning_WithPendingInstance_ReturnsTrue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		State: InstanceState{Name: "pending", Code: 0},
	}

	// Act
	result := instance.IsTransitioning()

	// Assert
	assert.True(t, result)
}

func TestEC2Instance_IsTransitioning_WithStoppingInstance_ReturnsTrue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		State: InstanceState{Name: "stopping", Code: 64},
	}

	// Act
	result := instance.IsTransitioning()

	// Assert
	assert.True(t, result)
}

func TestEC2Instance_IsTransitioning_WithRunningInstance_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		State: InstanceState{Name: "running", Code: 16},
	}

	// Act
	result := instance.IsTransitioning()

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_IsTransitioning_WithStoppedInstance_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		State: InstanceState{Name: "stopped", Code: 80},
	}

	// Act
	result := instance.IsTransitioning()

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_GetTag_WithExistingTag_ReturnsTagValue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
			{Key: "Team", Value: "backend"},
		},
	}
	key := "Name"

	// Act
	result := instance.GetTag(key)

	// Assert
	assert.Equal(t, "test-instance", result)
}

func TestEC2Instance_GetTag_WithCaseInsensitiveKey_ReturnsTagValue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
			{Key: "Team", Value: "backend"},
		},
	}
	key := "environment"

	// Act
	result := instance.GetTag(key)

	// Assert
	assert.Equal(t, "production", result)
}

func TestEC2Instance_GetTag_WithNonExistingTag_ReturnsEmptyString(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
			{Key: "Team", Value: "backend"},
		},
	}
	key := "NonExistent"

	// Act
	result := instance.GetTag(key)

	// Assert
	assert.Equal(t, "", result)
}

func TestEC2Instance_GetTag_WithEmptyKey_ReturnsEmptyString(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
			{Key: "Team", Value: "backend"},
		},
	}
	key := ""

	// Act
	result := instance.GetTag(key)

	// Assert
	assert.Equal(t, "", result)
}

func TestEC2Instance_HasTag_WithExactMatch_ReturnsTrue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
		},
	}
	key := "Name"
	value := "test-instance"

	// Act
	result := instance.HasTag(key, value)

	// Assert
	assert.True(t, result)
}

func TestEC2Instance_HasTag_WithCaseInsensitiveMatch_ReturnsTrue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
		},
	}
	key := "environment"
	value := "PRODUCTION"

	// Act
	result := instance.HasTag(key, value)

	// Assert
	assert.True(t, result)
}

func TestEC2Instance_HasTag_WithKeyExistsButValueDoesNotMatch_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
		},
	}
	key := "Name"
	value := "other-instance"

	// Act
	result := instance.HasTag(key, value)

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_HasTag_WithKeyDoesNotExist_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
		},
	}
	key := "NonExistent"
	value := "value"

	// Act
	result := instance.HasTag(key, value)

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_ShouldSkip_WithMatchingSkipTag_ReturnsTrue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Skip", Value: "monitoring"},
			{Key: "Environment", Value: "test"},
		},
	}
	skipList := []string{"monitoring", "backup"}

	// Act
	result := instance.ShouldSkip(skipList)

	// Assert
	assert.True(t, result)
}

func TestEC2Instance_ShouldSkip_WithNoMatchingSkipTag_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Skip", Value: "monitoring"},
			{Key: "Environment", Value: "test"},
		},
	}
	skipList := []string{"backup", "archival"}

	// Act
	result := instance.ShouldSkip(skipList)

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_ShouldSkip_WithEmptySkipList_ReturnsFalse(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Skip", Value: "monitoring"},
			{Key: "Environment", Value: "test"},
		},
	}
	skipList := []string{}

	// Act
	result := instance.ShouldSkip(skipList)

	// Assert
	assert.False(t, result)
}

func TestEC2Instance_GetDisplayName_WithNameTag_ReturnsTagValue(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		InstanceID: "i-1234567890abcdef0",
		Tags: []Tag{
			{Key: "Name", Value: "web-server"},
		},
	}

	// Act
	result := instance.GetDisplayName()

	// Assert
	assert.Equal(t, "web-server", result)
}

func TestEC2Instance_GetDisplayName_WithoutNameTag_ReturnsInstanceID(t *testing.T) {
	// Arrange
	instance := EC2Instance{
		InstanceID: "i-1234567890abcdef0",
		Tags:       []Tag{},
	}

	// Act
	result := instance.GetDisplayName()

	// Assert
	assert.Equal(t, "i-1234567890abcdef0", result)
}

func TestEC2Instance_GetCostEstimate_WithT2Micro_ReturnsCorrectCost(t *testing.T) {
	// Arrange
	instance := EC2Instance{InstanceType: "t2.micro"}

	// Act
	result := instance.GetCostEstimate()

	// Assert
	assert.Equal(t, 8.76, result)
}

func TestEC2Instance_GetCostEstimate_WithT3Small_ReturnsCorrectCost(t *testing.T) {
	// Arrange
	instance := EC2Instance{InstanceType: "t3.small"}

	// Act
	result := instance.GetCostEstimate()

	// Assert
	assert.Equal(t, 16.80, result)
}

func TestEC2Instance_GetCostEstimate_WithUnknownType_ReturnsZero(t *testing.T) {
	// Arrange
	instance := EC2Instance{InstanceType: "unknown.type"}

	// Act
	result := instance.GetCostEstimate()

	// Assert
	assert.Equal(t, 0.0, result)
}
