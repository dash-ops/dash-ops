package kubernetes

import (
	"testing"
)

// mockServiceContextResolver implements ServiceContextResolver for testing
type mockServiceContextResolver struct {
	contexts map[string]*ServiceContext
}

func newMockServiceContextResolver() *mockServiceContextResolver {
	return &mockServiceContextResolver{
		contexts: make(map[string]*ServiceContext),
	}
}

func (m *mockServiceContextResolver) AddContext(deploymentName, namespace, context string, ctx *ServiceContext) {
	key := context + "/" + namespace + "/" + deploymentName
	m.contexts[key] = ctx
}

func (m *mockServiceContextResolver) ResolveDeploymentService(deploymentName, namespace, context string) (*ServiceContext, error) {
	key := context + "/" + namespace + "/" + deploymentName
	if ctx, exists := m.contexts[key]; exists {
		return ctx, nil
	}
	return nil, nil
}

func TestServiceContextIntegration(t *testing.T) {
	// Create mock resolver with test data
	resolver := newMockServiceContextResolver()
	resolver.AddContext("auth-api", "auth", "docker-desktop", &ServiceContext{
		ServiceName: "user-authentication",
		ServiceTier: "TIER-1",
		Environment: "local",
		Context:     "docker-desktop",
		Team:        "auth-squad",
		Description: "Test authentication service",
	})

	// Create mock client (no expectations needed for this simple test)
	mockClient := &mockClient{}

	// Test the handler with context resolver
	handler := deploymentsHandlerWithContext(mockClient, resolver)

	// This test verifies that:
	// 1. The handler accepts a ServiceContextResolver
	// 2. The handler function can be created without errors
	// 3. The integration compiles and type-checks correctly

	if handler == nil {
		t.Fatal("Handler should not be nil")
	}

	t.Log("Service context integration test passed - handler creation successful")
}

func TestServiceContextResolver_Interface(t *testing.T) {
	// Test that mockServiceContextResolver implements the interface
	var resolver ServiceContextResolver = newMockServiceContextResolver()

	// Add test context
	mockResolver := resolver.(*mockServiceContextResolver)
	mockResolver.AddContext("test-deployment", "test-ns", "test-context", &ServiceContext{
		ServiceName: "test-service",
		ServiceTier: "TIER-2",
	})

	// Test resolution
	ctx, err := resolver.ResolveDeploymentService("test-deployment", "test-ns", "test-context")
	if err != nil {
		t.Fatal("Resolver should not return error:", err)
	}

	if ctx == nil {
		t.Fatal("Context should not be nil")
	}

	if ctx.ServiceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", ctx.ServiceName)
	}

	// Test non-existing deployment
	ctx2, err := resolver.ResolveDeploymentService("non-existing", "test-ns", "test-context")
	if err != nil {
		t.Fatal("Resolver should not return error for non-existing deployment:", err)
	}

	if ctx2 != nil {
		t.Fatal("Context should be nil for non-existing deployment")
	}

	t.Log("ServiceContextResolver interface test passed")
}
