package observability

import (
	"testing"
)

func TestNewModule(t *testing.T) {
	// Test with nil config
	module, err := NewModule(nil)
	if err == nil {
		t.Error("Expected error for nil config, got nil")
	}
	if module != nil {
		t.Error("Expected nil module for nil config")
	}

	// Test with valid config
	config := []byte(`
observability:
  enabled: true
  logs:
    providers:
      - name: "loki-local"
        type: "loki"
        url: "http://localhost:3100"
        timeout: "30s"
        enabled: true
  metrics:
    providers:
      - name: "prometheus-local"
        type: "prometheus"
        url: "http://localhost:9090"
        timeout: "30s"
        enabled: true
  traces:
    providers:
      - name: "tempo-local"
        type: "tempo"
        url: "http://localhost:3200"
        timeout: "30s"
        enabled: true
  alerts:
    providers:
      - name: "alertmanager-local"
        type: "alertmanager"
        url: "http://localhost:9093"
        timeout: "30s"
        enabled: true
`)

	module, err = NewModule(config)
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}
	if module == nil {
		t.Error("Expected non-nil module for valid config")
	}
	if module.Handler == nil {
		t.Error("Expected non-nil handler")
	}
}

func TestModule_RegisterRoutes(t *testing.T) {
	// Create a module with minimal config
	config := []byte(`
observability:
  enabled: true
  logs:
    providers:
      - name: "loki-local"
        type: "loki"
        url: "http://localhost:3100"
        timeout: "30s"
        enabled: true
`)

	module, err := NewModule(config)
	if err != nil {
		t.Fatalf("Failed to create module: %v", err)
	}

	// Test RegisterRoutes with nil router (should handle gracefully)
	module.RegisterRoutes(nil)
}
