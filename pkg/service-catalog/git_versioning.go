package servicecatalog

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GitVersioning manages Git operations for service definitions
type GitVersioning struct {
	repoPath string
	enabled  bool
}

// NewGitVersioning creates a new Git versioning manager
func NewGitVersioning(repoPath string) *GitVersioning {
	return &GitVersioning{
		repoPath: repoPath,
		enabled:  true, // Always enabled according to roadmap
	}
}

// InitializeRepository initializes Git repository if it doesn't exist
func (gv *GitVersioning) InitializeRepository() error {
	if !gv.enabled {
		return nil
	}

	// Check if .git directory exists
	gitDir := filepath.Join(gv.repoPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		// Repository already exists
		return nil
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = gv.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Set default branch to main if git version supports it
	cmd = exec.Command("git", "config", "init.defaultBranch", "main")
	cmd.Dir = gv.repoPath
	cmd.Run() // Ignore error for older git versions

	// Create .gitignore if it doesn't exist
	gitignorePath := filepath.Join(gv.repoPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		gitignoreContent := `# Service Catalog
*.tmp
*.bak
.DS_Store
`
		if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore: %w", err)
		}
	}

	// Initial commit
	if err := gv.commitChanges("Initial service catalog repository", &UserContext{
		Username: "system",
		Name:     "Dash-Ops System",
		Email:    "system@dash-ops.local",
	}, "initialize"); err != nil {
		// Silently continue if initial commit fails
		_ = err
	}

	return nil
}

// CommitServiceChange creates a git commit for service changes
func (gv *GitVersioning) CommitServiceChange(service *Service, user *UserContext, action string) error {
	if !gv.enabled {
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

	// Stage the service file
	serviceFile := service.Metadata.Name + ".yaml"
	if err := gv.stageFile(serviceFile); err != nil {
		return fmt.Errorf("failed to stage service file: %w", err)
	}

	// Create commit message
	commitMsg := gv.buildCommitMessage(service, user, action)

	// Create commit with user info
	if err := gv.commitWithAuthor(commitMsg, user); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}

// CommitServiceDeletion creates a git commit for service deletion
func (gv *GitVersioning) CommitServiceDeletion(serviceName string, user *UserContext) error {
	if !gv.enabled {
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

	// Stage the deleted file
	serviceFile := serviceName + ".yaml"
	if err := gv.stageFile(serviceFile); err != nil {
		return fmt.Errorf("failed to stage deleted service file: %w", err)
	}

	// Create commit message
	commitMsg := fmt.Sprintf("Delete service '%s' by %s\n\n- Action: delete\n- User: %s (%s)\n- Timestamp: %s",
		serviceName, user.Name, user.Username, user.Email, time.Now().Format(time.RFC3339))

	// Create commit with user info
	if err := gv.commitWithAuthor(commitMsg, user); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}

// GetServiceHistory returns git history for a specific service
func (gv *GitVersioning) GetServiceHistory(serviceName string) ([]ServiceChange, error) {
	if !gv.enabled {
		return nil, fmt.Errorf("git versioning is disabled")
	}

	serviceFile := serviceName + ".yaml"

	// Get git log for specific service file
	cmd := exec.Command("git", "log", "--oneline", "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso", "--", serviceFile)
	cmd.Dir = gv.repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get service history: %w", err)
	}

	return gv.parseGitLog(string(output)), nil
}

// GetAllHistory returns complete git history for all services
func (gv *GitVersioning) GetAllHistory() ([]ServiceChange, error) {
	if !gv.enabled {
		return nil, fmt.Errorf("git versioning is disabled")
	}

	// Get complete git log
	cmd := exec.Command("git", "log", "--oneline", "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso")
	cmd.Dir = gv.repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git history: %w", err)
	}

	return gv.parseGitLog(string(output)), nil
}

// IsRepositoryClean checks if there are uncommitted changes
func (gv *GitVersioning) IsRepositoryClean() (bool, error) {
	if !gv.enabled {
		return true, nil
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = gv.repoPath

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(strings.TrimSpace(string(output))) == 0, nil
}

// stageFile stages a file for commit
func (gv *GitVersioning) stageFile(filename string) error {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = gv.repoPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage file %s: %w", filename, err)
	}

	return nil
}

// commitChanges creates a commit with all staged changes
func (gv *GitVersioning) commitChanges(message string, user *UserContext, action string) error {
	// Check if there are any changes to commit
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = gv.repoPath

	if err := cmd.Run(); err == nil {
		// No changes staged
		return nil
	}

	return gv.commitWithAuthor(message, user)
}

// commitWithAuthor creates a commit with specific author information
func (gv *GitVersioning) commitWithAuthor(message string, user *UserContext) error {
	cmd := exec.Command("git", "commit", "-m", message,
		"--author", fmt.Sprintf("%s <%s>", user.Name, user.Email))
	cmd.Dir = gv.repoPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}

// buildCommitMessage builds a standardized commit message
func (gv *GitVersioning) buildCommitMessage(service *Service, user *UserContext, action string) string {
	actionTitle := strings.Title(action)

	message := fmt.Sprintf("%s service '%s' by %s\n\n", actionTitle, service.Metadata.Name, user.Name)
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

// parseGitLog parses git log output into ServiceChange structures
func (gv *GitVersioning) parseGitLog(output string) []ServiceChange {
	var changes []ServiceChange

	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 5 {
			continue
		}

		// Parse timestamp
		timestamp, err := time.Parse("2006-01-02 15:04:05 -0700", parts[3])
		if err != nil {
			// Fallback to current time if parsing fails
			timestamp = time.Now()
		}

		changes = append(changes, ServiceChange{
			Commit:    parts[0],
			Author:    parts[1],
			Email:     parts[2],
			Timestamp: timestamp,
			Message:   parts[4],
		})
	}

	return changes
}

// ValidateGitInstallation checks if git is available
func ValidateGitInstallation() error {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git is not installed or not available in PATH")
	}
	return nil
}

// SetupGitConfig sets up basic git configuration if not present
func (gv *GitVersioning) SetupGitConfig() error {
	if !gv.enabled {
		return nil
	}

	// Check if user.name is configured
	cmd := exec.Command("git", "config", "user.name")
	cmd.Dir = gv.repoPath
	if err := cmd.Run(); err != nil {
		// Set default user.name
		cmd = exec.Command("git", "config", "user.name", "Dash-Ops Service Catalog")
		cmd.Dir = gv.repoPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set git user.name: %w", err)
		}
	}

	// Check if user.email is configured
	cmd = exec.Command("git", "config", "user.email")
	cmd.Dir = gv.repoPath
	if err := cmd.Run(); err != nil {
		// Set default user.email
		cmd = exec.Command("git", "config", "user.email", "service-catalog@dash-ops.local")
		cmd.Dir = gv.repoPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set git user.email: %w", err)
		}
	}

	return nil
}

// GetRepositoryStatus returns current repository status
func (gv *GitVersioning) GetRepositoryStatus() (string, error) {
	if !gv.enabled {
		return "Git versioning disabled", nil
	}

	cmd := exec.Command("git", "status", "--short")
	cmd.Dir = gv.repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get repository status: %w", err)
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return "Repository is clean", nil
	}

	return string(output), nil
}
