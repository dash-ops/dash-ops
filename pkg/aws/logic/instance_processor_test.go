package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

func TestInstanceProcessor_ValidateInstanceOperation(t *testing.T) {
	processor := NewInstanceProcessor()

	tests := []struct {
		name        string
		instance    *awsModels.EC2Instance
		operation   string
		expectError bool
	}{
		{
			name: "start stopped instance",
			instance: &awsModels.EC2Instance{
				InstanceID: "i-123456789",
				State:      awsModels.InstanceStateStopped,
			},
			operation:   "start",
			expectError: false,
		},
		{
			name: "start running instance",
			instance: &awsModels.EC2Instance{
				InstanceID: "i-123456789",
				State:      awsModels.InstanceStateRunning,
			},
			operation:   "start",
			expectError: true,
		},
		{
			name: "stop running instance",
			instance: &awsModels.EC2Instance{
				InstanceID: "i-123456789",
				State:      awsModels.InstanceStateRunning,
			},
			operation:   "stop",
			expectError: false,
		},
		{
			name: "stop stopped instance",
			instance: &awsModels.EC2Instance{
				InstanceID: "i-123456789",
				State:      awsModels.InstanceStateStopped,
			},
			operation:   "stop",
			expectError: true,
		},
		{
			name:        "nil instance",
			instance:    nil,
			operation:   "start",
			expectError: true,
		},
		{
			name: "invalid operation",
			instance: &awsModels.EC2Instance{
				InstanceID: "i-123456789",
				State:      awsModels.InstanceStateRunning,
			},
			operation:   "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.ValidateInstanceOperation(tt.instance, tt.operation)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInstanceProcessor_FilterInstancesByPermissions(t *testing.T) {
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

	tests := []struct {
		name          string
		userGroups    []string
		operation     string
		expectedCount int
	}{
		{
			name:          "admin can start",
			userGroups:    []string{"admin"},
			operation:     "start",
			expectedCount: 2,
		},
		{
			name:          "ec2-operators can start",
			userGroups:    []string{"ec2-operators"},
			operation:     "start",
			expectedCount: 2,
		},
		{
			name:          "viewers cannot start",
			userGroups:    []string{"viewers"},
			operation:     "start",
			expectedCount: 0,
		},
		{
			name:          "admin can stop",
			userGroups:    []string{"admin"},
			operation:     "stop",
			expectedCount: 2,
		},
		{
			name:          "ec2-operators cannot stop",
			userGroups:    []string{"ec2-operators"},
			operation:     "stop",
			expectedCount: 0,
		},
		{
			name:          "viewers can view",
			userGroups:    []string{"viewers"},
			operation:     "view",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.FilterInstancesByPermissions(instances, account, tt.userGroups, tt.operation)
			assert.Equal(t, tt.expectedCount, len(result))
		})
	}
}

func TestInstanceProcessor_ValidateInstanceID(t *testing.T) {
	processor := NewInstanceProcessor()

	tests := []struct {
		name        string
		instanceID  string
		expectError bool
	}{
		{
			name:        "valid new format",
			instanceID:  "i-1234567890abcdef0",
			expectError: false,
		},
		{
			name:        "valid old format",
			instanceID:  "i-12345678",
			expectError: false,
		},
		{
			name:        "empty instance ID",
			instanceID:  "",
			expectError: true,
		},
		{
			name:        "invalid prefix",
			instanceID:  "inst-1234567890abcdef0",
			expectError: true,
		},
		{
			name:        "invalid length",
			instanceID:  "i-123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.ValidateInstanceID(tt.instanceID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInstanceProcessor_NormalizeInstanceType(t *testing.T) {
	processor := NewInstanceProcessor()

	tests := []struct {
		name         string
		instanceType string
		expected     string
	}{
		{
			name:         "uppercase",
			instanceType: "T2.MICRO",
			expected:     "t2.micro",
		},
		{
			name:         "mixed case",
			instanceType: "T3.Small",
			expected:     "t3.small",
		},
		{
			name:         "with whitespace",
			instanceType: " m5.large ",
			expected:     "m5.large",
		},
		{
			name:         "already normalized",
			instanceType: "t2.micro",
			expected:     "t2.micro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.NormalizeInstanceType(tt.instanceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInstanceProcessor_NormalizeRegion(t *testing.T) {
	processor := NewInstanceProcessor()

	tests := []struct {
		name     string
		region   string
		expected string
	}{
		{
			name:     "uppercase",
			region:   "US-EAST-1",
			expected: "us-east-1",
		},
		{
			name:     "with whitespace",
			region:   " eu-west-1 ",
			expected: "eu-west-1",
		},
		{
			name:     "already normalized",
			region:   "us-west-2",
			expected: "us-west-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.NormalizeRegion(tt.region)
			assert.Equal(t, tt.expected, result)
		})
	}
}
