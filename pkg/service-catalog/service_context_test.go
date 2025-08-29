package servicecatalog

import (
	"testing"
)

func TestMatchesDeploymentName(t *testing.T) {
	tests := []struct {
		name           string
		configuredName string
		actualName     string
		expected       bool
	}{
		{
			name:           "exact match",
			configuredName: "auth-api",
			actualName:     "auth-api",
			expected:       true,
		},
		{
			name:           "case insensitive match",
			configuredName: "auth-api",
			actualName:     "Auth-API",
			expected:       true,
		},
		{
			name:           "no match",
			configuredName: "auth-api",
			actualName:     "payment-api",
			expected:       false,
		},
		{
			name:           "empty strings",
			configuredName: "",
			actualName:     "",
			expected:       true,
		},
		{
			name:           "one empty string",
			configuredName: "auth-api",
			actualName:     "",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesDeploymentName(tt.configuredName, tt.actualName)
			if result != tt.expected {
				t.Errorf("matchesDeploymentName(%s, %s) = %v, want %v",
					tt.configuredName, tt.actualName, result, tt.expected)
			}
		})
	}
}

func TestCreateDeploymentKey(t *testing.T) {
	tests := []struct {
		name           string
		deploymentName string
		namespace      string
		context        string
		expected       string
	}{
		{
			name:           "full key",
			deploymentName: "auth-api",
			namespace:      "auth",
			context:        "docker-desktop",
			expected:       "docker-desktop/auth/auth-api",
		},
		{
			name:           "with empty context",
			deploymentName: "auth-api",
			namespace:      "auth",
			context:        "",
			expected:       "/auth/auth-api",
		},
		{
			name:           "all empty",
			deploymentName: "",
			namespace:      "",
			context:        "",
			expected:       "//",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createDeploymentKey(tt.deploymentName, tt.namespace, tt.context)
			if result != tt.expected {
				t.Errorf("createDeploymentKey(%s, %s, %s) = %s, want %s",
					tt.deploymentName, tt.namespace, tt.context, result, tt.expected)
			}
		})
	}
}

func TestValidateServiceContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *ServiceContext
		wantErr bool
	}{
		{
			name: "valid context",
			ctx: &ServiceContext{
				ServiceName: "user-authentication",
				ServiceTier: "TIER-1",
				Environment: "production",
				Context:     "prod-cluster",
				Team:        "auth-squad",
			},
			wantErr: false,
		},
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: true,
		},
		{
			name: "missing service name",
			ctx: &ServiceContext{
				ServiceTier: "TIER-1",
				Environment: "production",
			},
			wantErr: true,
		},
		{
			name: "missing service tier",
			ctx: &ServiceContext{
				ServiceName: "user-authentication",
				Environment: "production",
			},
			wantErr: true,
		},
		{
			name: "missing environment",
			ctx: &ServiceContext{
				ServiceName: "user-authentication",
				ServiceTier: "TIER-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServiceContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Mock functions for integration testing would go here
// For now, we'll focus on unit tests for the helper functions
