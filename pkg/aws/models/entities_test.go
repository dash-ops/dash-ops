package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAWSAccount_Validate(t *testing.T) {
	tests := []struct {
		name        string
		account     AWSAccount
		expectError bool
	}{
		{
			name: "valid account",
			account: AWSAccount{
				Name:            "test-account",
				Region:          "us-east-1",
				AccessKeyID:     "AKIATEST123",
				SecretAccessKey: "secret-key",
			},
			expectError: false,
		},
		{
			name: "missing name",
			account: AWSAccount{
				Region:          "us-east-1",
				AccessKeyID:     "AKIATEST123",
				SecretAccessKey: "secret-key",
			},
			expectError: true,
		},
		{
			name: "missing region",
			account: AWSAccount{
				Name:            "test-account",
				AccessKeyID:     "AKIATEST123",
				SecretAccessKey: "secret-key",
			},
			expectError: true,
		},
		{
			name: "missing access key",
			account: AWSAccount{
				Name:            "test-account",
				Region:          "us-east-1",
				SecretAccessKey: "secret-key",
			},
			expectError: true,
		},
		{
			name: "missing secret key",
			account: AWSAccount{
				Name:        "test-account",
				Region:      "us-east-1",
				AccessKeyID: "AKIATEST123",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAWSAccount_GenerateKey(t *testing.T) {
	tests := []struct {
		name        string
		accountName string
		expected    string
	}{
		{
			name:        "simple name",
			accountName: "Production",
			expected:    "production",
		},
		{
			name:        "name with spaces",
			accountName: "Test Account",
			expected:    "test_account",
		},
		{
			name:        "name with special chars",
			accountName: "Dev-Environment",
			expected:    "dev-environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := AWSAccount{Name: tt.accountName}
			account.GenerateKey()
			assert.Equal(t, tt.expected, account.Key)
		})
	}
}

func TestAWSAccount_HasEC2StartPermission(t *testing.T) {
	account := AWSAccount{
		Permissions: AccountPermissions{
			EC2: EC2Permissions{
				Start: []string{"admin", "ec2-operators"},
			},
		},
	}

	tests := []struct {
		name       string
		userGroups []string
		expected   bool
	}{
		{
			name:       "user has admin permission",
			userGroups: []string{"admin", "developers"},
			expected:   true,
		},
		{
			name:       "user has ec2-operators permission",
			userGroups: []string{"ec2-operators"},
			expected:   true,
		},
		{
			name:       "user has no permission",
			userGroups: []string{"developers", "viewers"},
			expected:   false,
		},
		{
			name:       "empty user groups",
			userGroups: []string{},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := account.HasEC2StartPermission(tt.userGroups)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEC2Instance_IsRunning(t *testing.T) {
	tests := []struct {
		name     string
		instance EC2Instance
		expected bool
	}{
		{
			name: "running instance",
			instance: EC2Instance{
				State: InstanceState{Name: "running", Code: 16},
			},
			expected: true,
		},
		{
			name: "stopped instance",
			instance: EC2Instance{
				State: InstanceState{Name: "stopped", Code: 80},
			},
			expected: false,
		},
		{
			name: "pending instance",
			instance: EC2Instance{
				State: InstanceState{Name: "pending", Code: 0},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.instance.IsRunning()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEC2Instance_IsTransitioning(t *testing.T) {
	tests := []struct {
		name     string
		instance EC2Instance
		expected bool
	}{
		{
			name: "pending instance",
			instance: EC2Instance{
				State: InstanceState{Name: "pending", Code: 0},
			},
			expected: true,
		},
		{
			name: "stopping instance",
			instance: EC2Instance{
				State: InstanceState{Name: "stopping", Code: 64},
			},
			expected: true,
		},
		{
			name: "running instance",
			instance: EC2Instance{
				State: InstanceState{Name: "running", Code: 16},
			},
			expected: false,
		},
		{
			name: "stopped instance",
			instance: EC2Instance{
				State: InstanceState{Name: "stopped", Code: 80},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.instance.IsTransitioning()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEC2Instance_GetTag(t *testing.T) {
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
			{Key: "Team", Value: "backend"},
		},
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "existing tag",
			key:      "Name",
			expected: "test-instance",
		},
		{
			name:     "case insensitive",
			key:      "environment",
			expected: "production",
		},
		{
			name:     "non-existing tag",
			key:      "NonExistent",
			expected: "",
		},
		{
			name:     "empty key",
			key:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := instance.GetTag(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEC2Instance_HasTag(t *testing.T) {
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Environment", Value: "production"},
		},
	}

	tests := []struct {
		name     string
		key      string
		value    string
		expected bool
	}{
		{
			name:     "exact match",
			key:      "Name",
			value:    "test-instance",
			expected: true,
		},
		{
			name:     "case insensitive match",
			key:      "environment",
			value:    "PRODUCTION",
			expected: true,
		},
		{
			name:     "key exists but value doesn't match",
			key:      "Name",
			value:    "other-instance",
			expected: false,
		},
		{
			name:     "key doesn't exist",
			key:      "NonExistent",
			value:    "value",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := instance.HasTag(tt.key, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEC2Instance_ShouldSkip(t *testing.T) {
	instance := EC2Instance{
		Tags: []Tag{
			{Key: "Name", Value: "test-instance"},
			{Key: "Skip", Value: "monitoring"},
			{Key: "Environment", Value: "test"},
		},
	}

	tests := []struct {
		name     string
		skipList []string
		expected bool
	}{
		{
			name:     "should skip - monitoring",
			skipList: []string{"monitoring", "backup"},
			expected: true,
		},
		{
			name:     "should not skip",
			skipList: []string{"backup", "archival"},
			expected: false,
		},
		{
			name:     "empty skip list",
			skipList: []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := instance.ShouldSkip(tt.skipList)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEC2Instance_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		instance EC2Instance
		expected string
	}{
		{
			name: "with name tag",
			instance: EC2Instance{
				InstanceID: "i-1234567890abcdef0",
				Tags: []Tag{
					{Key: "Name", Value: "web-server"},
				},
			},
			expected: "web-server",
		},
		{
			name: "without name tag",
			instance: EC2Instance{
				InstanceID: "i-1234567890abcdef0",
				Tags:       []Tag{},
			},
			expected: "i-1234567890abcdef0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.instance.GetDisplayName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEC2Instance_GetCostEstimate(t *testing.T) {
	tests := []struct {
		name         string
		instanceType string
		expected     float64
	}{
		{
			name:         "t2.micro",
			instanceType: "t2.micro",
			expected:     8.76,
		},
		{
			name:         "t3.small",
			instanceType: "t3.small",
			expected:     16.80,
		},
		{
			name:         "unknown type",
			instanceType: "unknown.type",
			expected:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := EC2Instance{InstanceType: tt.instanceType}
			result := instance.GetCostEstimate()
			assert.Equal(t, tt.expected, result)
		})
	}
}
