package github

import (
	"fmt"
	"strings"
	"time"
)

// GitHubUser represents a GitHub user
type GitHubUser struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	HTMLURL   string    `json:"html_url"`
	Type      string    `json:"type"` // User, Organization
	SiteAdmin bool      `json:"site_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Additional user information
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Blog     string `json:"blog,omitempty"`

	// Statistics
	PublicRepos int `json:"public_repos"`
	Followers   int `json:"followers"`
	Following   int `json:"following"`
}

// GitHubOrganization represents a GitHub organization
type GitHubOrganization struct {
	ID          int64     `json:"id"`
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	AvatarURL   string    `json:"avatar_url"`
	HTMLURL     string    `json:"html_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Organization specific
	Company           string `json:"company,omitempty"`
	Location          string `json:"location,omitempty"`
	Email             string `json:"email,omitempty"`
	Blog              string `json:"blog,omitempty"`
	PublicRepos       int    `json:"public_repos"`
	TotalPrivateRepos int    `json:"total_private_repos,omitempty"`

	// User's role in organization
	Role string `json:"role,omitempty"` // admin, member
}

// GitHubTeam represents a GitHub team
type GitHubTeam struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	Slug         string              `json:"slug"`
	Description  string              `json:"description,omitempty"`
	Privacy      string              `json:"privacy"`    // closed, secret
	Permission   string              `json:"permission"` // pull, push, admin
	HTMLURL      string              `json:"html_url"`
	Organization *GitHubOrganization `json:"organization"`
	Members      []GitHubUser        `json:"members,omitempty"`
	MemberCount  int                 `json:"member_count"`

	// User's role in team
	Role string `json:"role,omitempty"` // maintainer, member
}

// GitHubRepository represents a GitHub repository
type GitHubRepository struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	Description string      `json:"description,omitempty"`
	Private     bool        `json:"private"`
	HTMLURL     string      `json:"html_url"`
	CloneURL    string      `json:"clone_url"`
	SSHURL      string      `json:"ssh_url"`
	Owner       *GitHubUser `json:"owner"`

	// Repository metadata
	Language      string    `json:"language,omitempty"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PushedAt      time.Time `json:"pushed_at,omitempty"`

	// Statistics
	Size            int `json:"size"`
	StargazersCount int `json:"stargazers_count"`
	WatchersCount   int `json:"watchers_count"`
	ForksCount      int `json:"forks_count"`
	OpenIssuesCount int `json:"open_issues_count"`

	// Features
	HasIssues    bool `json:"has_issues"`
	HasProjects  bool `json:"has_projects"`
	HasWiki      bool `json:"has_wiki"`
	HasPages     bool `json:"has_pages"`
	HasDownloads bool `json:"has_downloads"`

	// User's permission level
	Permissions *RepositoryPermissions `json:"permissions,omitempty"`
}

// RepositoryPermissions represents user's permissions on a repository
type RepositoryPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

// GitHubIssue represents a GitHub issue
type GitHubIssue struct {
	ID       int64         `json:"id"`
	Number   int           `json:"number"`
	Title    string        `json:"title"`
	Body     string        `json:"body,omitempty"`
	State    string        `json:"state"` // open, closed
	HTMLURL  string        `json:"html_url"`
	User     *GitHubUser   `json:"user"`
	Assignee *GitHubUser   `json:"assignee,omitempty"`
	Labels   []GitHubLabel `json:"labels,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`

	// Statistics
	Comments int `json:"comments"`
}

// GitHubLabel represents a GitHub label
type GitHubLabel struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

// GitHubPullRequest represents a GitHub pull request
type GitHubPullRequest struct {
	ID      int64       `json:"id"`
	Number  int         `json:"number"`
	Title   string      `json:"title"`
	Body    string      `json:"body,omitempty"`
	State   string      `json:"state"` // open, closed
	HTMLURL string      `json:"html_url"`
	User    *GitHubUser `json:"user"`

	// PR specific
	Head      *GitHubBranch `json:"head"`
	Base      *GitHubBranch `json:"base"`
	Merged    bool          `json:"merged"`
	Mergeable *bool         `json:"mergeable,omitempty"`
	MergedAt  *time.Time    `json:"merged_at,omitempty"`
	MergedBy  *GitHubUser   `json:"merged_by,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`

	// Statistics
	Comments       int `json:"comments"`
	ReviewComments int `json:"review_comments"`
	Commits        int `json:"commits"`
	Additions      int `json:"additions"`
	Deletions      int `json:"deletions"`
	ChangedFiles   int `json:"changed_files"`
}

// GitHubBranch represents a GitHub branch
type GitHubBranch struct {
	Name string            `json:"name"`
	SHA  string            `json:"sha"`
	Repo *GitHubRepository `json:"repo,omitempty"`
}

// Domain methods for GitHubUser

// IsOrganization checks if user is an organization
func (u *GitHubUser) IsOrganization() bool {
	return u.Type == "Organization"
}

// GetDisplayName returns display name (Name or Login)
func (u *GitHubUser) GetDisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Login
}

// Validate validates user data
func (u *GitHubUser) Validate() error {
	if u.Login == "" {
		return fmt.Errorf("user login is required")
	}
	if u.ID == 0 {
		return fmt.Errorf("user ID is required")
	}
	return nil
}

// Domain methods for GitHubOrganization

// Validate validates organization data
func (o *GitHubOrganization) Validate() error {
	if o.Login == "" {
		return fmt.Errorf("organization login is required")
	}
	if o.ID == 0 {
		return fmt.Errorf("organization ID is required")
	}
	return nil
}

// GetDisplayName returns display name (Name or Login)
func (o *GitHubOrganization) GetDisplayName() string {
	if o.Name != "" {
		return o.Name
	}
	return o.Login
}

// IsUserAdmin checks if user has admin role
func (o *GitHubOrganization) IsUserAdmin() bool {
	return o.Role == "admin"
}

// Domain methods for GitHubTeam

// Validate validates team data
func (t *GitHubTeam) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("team name is required")
	}
	if t.Slug == "" {
		return fmt.Errorf("team slug is required")
	}
	if t.Organization == nil {
		return fmt.Errorf("team organization is required")
	}
	return nil
}

// IsSecret checks if team is secret
func (t *GitHubTeam) IsSecret() bool {
	return t.Privacy == "secret"
}

// IsMaintainer checks if user is a maintainer of the team
func (t *GitHubTeam) IsMaintainer() bool {
	return t.Role == "maintainer"
}

// HasMember checks if user is a member of the team
func (t *GitHubTeam) HasMember(userLogin string) bool {
	for _, member := range t.Members {
		if strings.EqualFold(member.Login, userLogin) {
			return true
		}
	}
	return false
}

// GetFullName returns team full name (org/team)
func (t *GitHubTeam) GetFullName() string {
	if t.Organization != nil {
		return fmt.Sprintf("%s/%s", t.Organization.Login, t.Slug)
	}
	return t.Slug
}

// Domain methods for GitHubRepository

// Validate validates repository data
func (r *GitHubRepository) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("repository name is required")
	}
	if r.FullName == "" {
		return fmt.Errorf("repository full name is required")
	}
	if r.Owner == nil {
		return fmt.Errorf("repository owner is required")
	}
	return nil
}

// IsArchived checks if repository is archived
func (r *GitHubRepository) IsArchived() bool {
	// Would need to add Archived field to the struct
	return false
}

// IsForked checks if repository is a fork
func (r *GitHubRepository) IsForked() bool {
	// Would need to add Fork field to the struct
	return false
}

// CanUserPush checks if user can push to repository
func (r *GitHubRepository) CanUserPush() bool {
	return r.Permissions != nil && r.Permissions.Push
}

// CanUserAdmin checks if user can admin repository
func (r *GitHubRepository) CanUserAdmin() bool {
	return r.Permissions != nil && r.Permissions.Admin
}

// GetLanguageOrDefault returns language or default
func (r *GitHubRepository) GetLanguageOrDefault() string {
	if r.Language != "" {
		return r.Language
	}
	return "Unknown"
}

// Domain methods for GitHubIssue

// IsOpen checks if issue is open
func (i *GitHubIssue) IsOpen() bool {
	return i.State == "open"
}

// IsClosed checks if issue is closed
func (i *GitHubIssue) IsClosed() bool {
	return i.State == "closed"
}

// IsAssigned checks if issue is assigned
func (i *GitHubIssue) IsAssigned() bool {
	return i.Assignee != nil
}

// HasLabels checks if issue has labels
func (i *GitHubIssue) HasLabels() bool {
	return len(i.Labels) > 0
}

// HasLabel checks if issue has a specific label
func (i *GitHubIssue) HasLabel(labelName string) bool {
	for _, label := range i.Labels {
		if strings.EqualFold(label.Name, labelName) {
			return true
		}
	}
	return false
}

// GetAge returns issue age
func (i *GitHubIssue) GetAge() time.Duration {
	return time.Since(i.CreatedAt)
}

// Domain methods for GitHubPullRequest

// IsOpen checks if PR is open
func (pr *GitHubPullRequest) IsOpen() bool {
	return pr.State == "open"
}

// IsClosed checks if PR is closed
func (pr *GitHubPullRequest) IsClosed() bool {
	return pr.State == "closed"
}

// IsMerged checks if PR is merged
func (pr *GitHubPullRequest) IsMerged() bool {
	return pr.Merged
}

// IsDraft checks if PR is a draft
func (pr *GitHubPullRequest) IsDraft() bool {
	// Would need to add Draft field to the struct
	return false
}

// GetAge returns PR age
func (pr *GitHubPullRequest) GetAge() time.Duration {
	return time.Since(pr.CreatedAt)
}

// GetChangeSize returns total lines changed
func (pr *GitHubPullRequest) GetChangeSize() int {
	return pr.Additions + pr.Deletions
}

// IsLargePR checks if PR is considered large (>500 lines changed)
func (pr *GitHubPullRequest) IsLargePR() bool {
	return pr.GetChangeSize() > 500
}
