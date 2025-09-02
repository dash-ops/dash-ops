package servicecatalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// SimpleVersioning provides basic versioning without Git
// Stores change history in JSON files alongside service definitions
type SimpleVersioning struct {
	directory string
	enabled   bool
}

// NewSimpleVersioning creates a new simple versioning provider
func NewSimpleVersioning(directory string) *SimpleVersioning {
	return &SimpleVersioning{
		directory: directory,
		enabled:   true,
	}
}

// Initialize sets up the simple versioning system (implements VersioningProvider)
func (sv *SimpleVersioning) Initialize() error {
	if !sv.enabled {
		return nil
	}

	// Ensure directory exists
	if err := os.MkdirAll(sv.directory, 0755); err != nil {
		return fmt.Errorf("failed to create versioning directory: %w", err)
	}

	// Create .history directory if it doesn't exist
	historyDir := filepath.Join(sv.directory, ".history")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	return nil
}

// CommitServiceChange records a service change (implements VersioningProvider)
func (sv *SimpleVersioning) CommitServiceChange(service *Service, user *UserContext, action string) error {
	if !sv.enabled {
		return nil
	}

	// Ensure we have user context
	if user == nil {
		user = &UserContext{
			Username: "anonymous",
			Name:     "Anonymous User",
			Email:    "anonymous@dash-ops.local",
		}
	}

	// Create service change record
	change := ServiceChange{
		Commit:    sv.generateChangeID(service.Metadata.Name, action),
		Author:    user.Name,
		Email:     user.Email,
		Timestamp: time.Now(),
		Message:   sv.buildChangeMessage(service, user, action),
	}

	// Save to history file
	return sv.saveServiceChange(service.Metadata.Name, change)
}

// CommitServiceDeletion records a service deletion (implements VersioningProvider)
func (sv *SimpleVersioning) CommitServiceDeletion(serviceName string, user *UserContext) error {
	if !sv.enabled {
		return nil
	}

	// Ensure we have user context
	if user == nil {
		user = &UserContext{
			Username: "anonymous",
			Name:     "Anonymous User",
			Email:    "anonymous@dash-ops.local",
		}
	}

	// Create deletion change record
	change := ServiceChange{
		Commit:    sv.generateChangeID(serviceName, "delete"),
		Author:    user.Name,
		Email:     user.Email,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Delete service '%s' by %s", serviceName, user.Name),
	}

	// Save to history file
	return sv.saveServiceChange(serviceName, change)
}

// GetServiceHistory returns change history for a specific service (implements VersioningProvider)
func (sv *SimpleVersioning) GetServiceHistory(serviceName string) ([]ServiceChange, error) {
	if !sv.enabled {
		return []ServiceChange{}, nil
	}

	historyFile := sv.getHistoryFilePath(serviceName)

	// Check if history file exists
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return []ServiceChange{}, nil
	}

	// Read history file
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var history []ServiceChange
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse history file: %w", err)
	}

	// Sort by timestamp (newest first)
	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp.After(history[j].Timestamp)
	})

	return history, nil
}

// GetAllHistory returns complete change history (implements VersioningProvider)
func (sv *SimpleVersioning) GetAllHistory() ([]ServiceChange, error) {
	if !sv.enabled {
		return []ServiceChange{}, nil
	}

	historyDir := filepath.Join(sv.directory, ".history")

	// Read all history files
	files, err := os.ReadDir(historyDir)
	if err != nil {
		return []ServiceChange{}, nil // Return empty if directory doesn't exist
	}

	var allChanges []ServiceChange

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Extract service name from filename
		serviceName := file.Name()[:len(file.Name())-5] // Remove .json extension

		serviceHistory, err := sv.GetServiceHistory(serviceName)
		if err != nil {
			continue // Skip files with errors
		}

		allChanges = append(allChanges, serviceHistory...)
	}

	// Sort all changes by timestamp (newest first)
	sort.Slice(allChanges, func(i, j int) bool {
		return allChanges[i].Timestamp.After(allChanges[j].Timestamp)
	})

	return allChanges, nil
}

// IsEnabled returns whether simple versioning is enabled (implements VersioningProvider)
func (sv *SimpleVersioning) IsEnabled() bool {
	return sv.enabled
}

// GetStatus returns current versioning system status (implements VersioningProvider)
func (sv *SimpleVersioning) GetStatus() (string, error) {
	if !sv.enabled {
		return "Simple versioning disabled", nil
	}

	historyDir := filepath.Join(sv.directory, ".history")
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		return "History directory not found", nil
	}

	// Count history files
	files, err := os.ReadDir(historyDir)
	if err != nil {
		return "Failed to read history directory", err
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			count++
		}
	}

	return fmt.Sprintf("Simple versioning active, tracking %d services", count), nil
}

// getHistoryFilePath returns the file path for a service's history
func (sv *SimpleVersioning) getHistoryFilePath(serviceName string) string {
	return filepath.Join(sv.directory, ".history", serviceName+".json")
}

// saveServiceChange saves a service change to the history file
func (sv *SimpleVersioning) saveServiceChange(serviceName string, change ServiceChange) error {
	historyFile := sv.getHistoryFilePath(serviceName)

	// Read existing history
	var history []ServiceChange
	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &history)
	}

	// Add new change
	history = append(history, change)

	// Sort by timestamp (newest first)
	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp.After(history[j].Timestamp)
	})

	// Limit history to last 100 entries to prevent unlimited growth
	if len(history) > 100 {
		history = history[:100]
	}

	// Write back to file
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	if err := os.WriteFile(historyFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

// generateChangeID generates a unique change ID
func (sv *SimpleVersioning) generateChangeID(serviceName, action string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%s-%d", serviceName, action, timestamp)
}

// buildChangeMessage builds a standardized change message
func (sv *SimpleVersioning) buildChangeMessage(service *Service, user *UserContext, action string) string {
	message := fmt.Sprintf("%s service '%s' by %s\n\n",
		capitalizeFirst(action), service.Metadata.Name, user.Name)
	message += fmt.Sprintf("- Action: %s\n", action)
	message += fmt.Sprintf("- User: %s (%s)\n", user.Username, user.Email)
	message += fmt.Sprintf("- Timestamp: %s\n", time.Now().Format(time.RFC3339))
	message += fmt.Sprintf("- Tier: %s\n", service.Metadata.Tier)
	message += fmt.Sprintf("- Team: %s\n", service.Spec.Team.GitHubTeam)

	if service.Metadata.Version > 0 {
		message += fmt.Sprintf("- Version: %d\n", service.Metadata.Version)
	}

	return message
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
