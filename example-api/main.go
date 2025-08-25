package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

// APIInfo represents basic API information
type APIInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
	GoVersion string    `json:"go_version"`
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

var startTime time.Time

func main() {
	startTime = time.Now()

	router := mux.NewRouter()

	// Health endpoints (for Kubernetes probes)
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/ready", readinessHandler).Methods("GET")

	// API endpoints
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/info", infoHandler).Methods("GET")
	router.HandleFunc("/api/status", statusHandler).Methods("GET")
	router.HandleFunc("/api/version", versionHandler).Methods("GET")

	// Get port from environment or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 User Authentication API starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	log.Printf("📋 API Info: http://localhost:%s/info", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

// rootHandler handles the root endpoint
func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "👋 Welcome to User Authentication API",
		"version": "v1.0.0",
		"docs":    "/info",
		"health":  "/health",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// healthHandler handles liveness probe
func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks: map[string]string{
			"database": "connected",
			"cache":    "operational",
			"api":      "running",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

// readinessHandler handles readiness probe
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate some startup time
	if time.Since(startTime) < 10*time.Second {
		http.Error(w, "Service not ready yet", http.StatusServiceUnavailable)
		return
	}

	ready := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Checks: map[string]string{
			"startup":    "completed",
			"migrations": "applied",
			"config":     "loaded",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ready)
}

// infoHandler provides detailed API information
func infoHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()

	info := APIInfo{
		Name:      "User Authentication API",
		Version:   "v1.0.0",
		Uptime:    time.Since(startTime).String(),
		GoVersion: runtime.Version(),
		Timestamp: time.Now(),
		Host:      hostname,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// statusHandler provides service status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"service": "user-authentication",
		"status":  "running",
		"uptime":  time.Since(startTime).String(),
		"tier":    "TIER-1",
		"team":    "auth-squad",
		"environment": map[string]interface{}{
			"NODE_ENV": getEnvWithDefault("NODE_ENV", "production"),
			"PORT":     getEnvWithDefault("PORT", "8080"),
		},
		"resources": map[string]interface{}{
			"cpu_usage":    "12%",
			"memory_usage": "45MB",
			"connections":  23,
		},
		"dependencies": map[string]string{
			"database":      "healthy",
			"email-service": "healthy",
			"cache":         "healthy",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// versionHandler returns API version
func versionHandler(w http.ResponseWriter, r *http.Request) {
	version := map[string]string{
		"version":    "v1.0.0",
		"build_time": "2024-01-20T15:30:00Z",
		"git_commit": "abc123d",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

// getEnvWithDefault returns environment variable or default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
