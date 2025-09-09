package github

import (
	"strings"
	"time"
)

// UserProfile represents comprehensive user profile information
type UserProfile struct {
	User          GitHubUser           `json:"user"`
	Organizations []GitHubOrganization `json:"organizations"`
	Teams         []GitHubTeam         `json:"teams"`
	Repositories  []GitHubRepository   `json:"repositories,omitempty"`
	LastUpdated   time.Time            `json:"last_updated"`
}

// OrganizationProfile represents organization profile with members and teams
type OrganizationProfile struct {
	Organization GitHubOrganization `json:"organization"`
	Teams        []GitHubTeam       `json:"teams"`
	Members      []GitHubUser       `json:"members,omitempty"`
	Repositories []GitHubRepository `json:"repositories,omitempty"`
	LastUpdated  time.Time          `json:"last_updated"`
}

// TeamDetails represents detailed team information
type TeamDetails struct {
	Team         GitHubTeam         `json:"team"`
	Members      []GitHubUser       `json:"members"`
	Repositories []GitHubRepository `json:"repositories,omitempty"`
	LastUpdated  time.Time          `json:"last_updated"`
}

// RepositoryList represents a list of repositories with metadata
type RepositoryList struct {
	Repositories []GitHubRepository `json:"repositories"`
	Total        int                `json:"total"`
	Owner        string             `json:"owner,omitempty"`
	Filter       *RepositoryFilter  `json:"filter,omitempty"`
}

// RepositoryFilter represents filtering criteria for repositories
type RepositoryFilter struct {
	Owner     string   `json:"owner,omitempty"`
	Type      string   `json:"type,omitempty"` // all, owner, member, public, private
	Language  string   `json:"language,omitempty"`
	Sort      string   `json:"sort,omitempty"`      // created, updated, pushed, full_name
	Direction string   `json:"direction,omitempty"` // asc, desc
	Search    string   `json:"search,omitempty"`
	Topics    []string `json:"topics,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
}

// IssueList represents a list of issues with metadata
type IssueList struct {
	Issues     []GitHubIssue `json:"issues"`
	Total      int           `json:"total"`
	Repository string        `json:"repository,omitempty"`
	Filter     *IssueFilter  `json:"filter,omitempty"`
}

// IssueFilter represents filtering criteria for issues
type IssueFilter struct {
	Repository string     `json:"repository,omitempty"`
	State      string     `json:"state,omitempty"` // open, closed, all
	Labels     []string   `json:"labels,omitempty"`
	Assignee   string     `json:"assignee,omitempty"`
	Creator    string     `json:"creator,omitempty"`
	Sort       string     `json:"sort,omitempty"`      // created, updated, comments
	Direction  string     `json:"direction,omitempty"` // asc, desc
	Since      *time.Time `json:"since,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
}

// PullRequestList represents a list of pull requests with metadata
type PullRequestList struct {
	PullRequests []GitHubPullRequest `json:"pull_requests"`
	Total        int                 `json:"total"`
	Repository   string              `json:"repository,omitempty"`
	Filter       *PullRequestFilter  `json:"filter,omitempty"`
}

// PullRequestFilter represents filtering criteria for pull requests
type PullRequestFilter struct {
	Repository string     `json:"repository,omitempty"`
	State      string     `json:"state,omitempty"`     // open, closed, all
	Head       string     `json:"head,omitempty"`      // branch name
	Base       string     `json:"base,omitempty"`      // branch name
	Sort       string     `json:"sort,omitempty"`      // created, updated, popularity
	Direction  string     `json:"direction,omitempty"` // asc, desc
	Since      *time.Time `json:"since,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
}

// Methods for UserProfile

// GetOrganizationByLogin returns organization by login
func (up *UserProfile) GetOrganizationByLogin(login string) *GitHubOrganization {
	for _, org := range up.Organizations {
		if strings.EqualFold(org.Login, login) {
			return &org
		}
	}
	return nil
}

// GetTeamBySlug returns team by slug
func (up *UserProfile) GetTeamBySlug(orgLogin, teamSlug string) *GitHubTeam {
	for _, team := range up.Teams {
		if team.Organization != nil &&
			strings.EqualFold(team.Organization.Login, orgLogin) &&
			strings.EqualFold(team.Slug, teamSlug) {
			return &team
		}
	}
	return nil
}

// GetTeamsInOrganization returns all teams for a specific organization
func (up *UserProfile) GetTeamsInOrganization(orgLogin string) []GitHubTeam {
	teams := make([]GitHubTeam, 0)
	for _, team := range up.Teams {
		if team.Organization != nil && strings.EqualFold(team.Organization.Login, orgLogin) {
			teams = append(teams, team)
		}
	}
	return teams
}

// IsOrgAdmin checks if user is admin of any organization
func (up *UserProfile) IsOrgAdmin(orgLogin string) bool {
	org := up.GetOrganizationByLogin(orgLogin)
	return org != nil && org.IsUserAdmin()
}

// Methods for RepositoryList

// FilterByLanguage filters repositories by programming language
func (rl *RepositoryList) FilterByLanguage(language string) *RepositoryList {
	if language == "" {
		return rl
	}

	var filtered []GitHubRepository
	for _, repo := range rl.Repositories {
		if strings.EqualFold(repo.Language, language) {
			filtered = append(filtered, repo)
		}
	}

	return &RepositoryList{
		Repositories: filtered,
		Total:        len(filtered),
		Owner:        rl.Owner,
		Filter:       &RepositoryFilter{Language: language},
	}
}

// FilterByType filters repositories by type
func (rl *RepositoryList) FilterByType(repoType string) *RepositoryList {
	if repoType == "" {
		return rl
	}

	var filtered []GitHubRepository
	for _, repo := range rl.Repositories {
		switch repoType {
		case "public":
			if !repo.Private {
				filtered = append(filtered, repo)
			}
		case "private":
			if repo.Private {
				filtered = append(filtered, repo)
			}
		default:
			filtered = append(filtered, repo)
		}
	}

	return &RepositoryList{
		Repositories: filtered,
		Total:        len(filtered),
		Owner:        rl.Owner,
		Filter:       &RepositoryFilter{Type: repoType},
	}
}

// Search filters repositories by text search
func (rl *RepositoryList) Search(query string) *RepositoryList {
	if query == "" {
		return rl
	}

	query = strings.ToLower(query)
	var filtered []GitHubRepository

	for _, repo := range rl.Repositories {
		if rl.matchesSearch(repo, query) {
			filtered = append(filtered, repo)
		}
	}

	return &RepositoryList{
		Repositories: filtered,
		Total:        len(filtered),
		Owner:        rl.Owner,
		Filter:       &RepositoryFilter{Search: query},
	}
}

// matchesSearch checks if repository matches search query
func (rl *RepositoryList) matchesSearch(repo GitHubRepository, query string) bool {
	searchFields := []string{
		strings.ToLower(repo.Name),
		strings.ToLower(repo.FullName),
		strings.ToLower(repo.Description),
		strings.ToLower(repo.Language),
	}

	for _, field := range searchFields {
		if strings.Contains(field, query) {
			return true
		}
	}

	return false
}

// GetPublicRepositories returns only public repositories
func (rl *RepositoryList) GetPublicRepositories() []GitHubRepository {
	var public []GitHubRepository
	for _, repo := range rl.Repositories {
		if !repo.Private {
			public = append(public, repo)
		}
	}
	return public
}

// GetPrivateRepositories returns only private repositories
func (rl *RepositoryList) GetPrivateRepositories() []GitHubRepository {
	var private []GitHubRepository
	for _, repo := range rl.Repositories {
		if repo.Private {
			private = append(private, repo)
		}
	}
	return private
}

// Methods for IssueList

// FilterByState filters issues by state
func (il *IssueList) FilterByState(state string) *IssueList {
	if state == "" {
		return il
	}

	var filtered []GitHubIssue
	for _, issue := range il.Issues {
		if strings.EqualFold(issue.State, state) {
			filtered = append(filtered, issue)
		}
	}

	return &IssueList{
		Issues:     filtered,
		Total:      len(filtered),
		Repository: il.Repository,
		Filter:     &IssueFilter{State: state},
	}
}

// FilterByLabel filters issues by label
func (il *IssueList) FilterByLabel(labelName string) *IssueList {
	if labelName == "" {
		return il
	}

	var filtered []GitHubIssue
	for _, issue := range il.Issues {
		if issue.HasLabel(labelName) {
			filtered = append(filtered, issue)
		}
	}

	return &IssueList{
		Issues:     filtered,
		Total:      len(filtered),
		Repository: il.Repository,
		Filter:     &IssueFilter{Labels: []string{labelName}},
	}
}

// GetOpenIssues returns only open issues
func (il *IssueList) GetOpenIssues() []GitHubIssue {
	var open []GitHubIssue
	for _, issue := range il.Issues {
		if issue.IsOpen() {
			open = append(open, issue)
		}
	}
	return open
}

// GetUnassignedIssues returns only unassigned issues
func (il *IssueList) GetUnassignedIssues() []GitHubIssue {
	var unassigned []GitHubIssue
	for _, issue := range il.Issues {
		if !issue.IsAssigned() {
			unassigned = append(unassigned, issue)
		}
	}
	return unassigned
}

// Methods for PullRequestList

// FilterByState filters pull requests by state
func (prl *PullRequestList) FilterByState(state string) *PullRequestList {
	if state == "" {
		return prl
	}

	var filtered []GitHubPullRequest
	for _, pr := range prl.PullRequests {
		if strings.EqualFold(pr.State, state) {
			filtered = append(filtered, pr)
		}
	}

	return &PullRequestList{
		PullRequests: filtered,
		Total:        len(filtered),
		Repository:   prl.Repository,
		Filter:       &PullRequestFilter{State: state},
	}
}

// GetOpenPullRequests returns only open pull requests
func (prl *PullRequestList) GetOpenPullRequests() []GitHubPullRequest {
	var open []GitHubPullRequest
	for _, pr := range prl.PullRequests {
		if pr.IsOpen() {
			open = append(open, pr)
		}
	}
	return open
}

// GetMergedPullRequests returns only merged pull requests
func (prl *PullRequestList) GetMergedPullRequests() []GitHubPullRequest {
	var merged []GitHubPullRequest
	for _, pr := range prl.PullRequests {
		if pr.IsMerged() {
			merged = append(merged, pr)
		}
	}
	return merged
}

// GetLargePullRequests returns pull requests with significant changes
func (prl *PullRequestList) GetLargePullRequests() []GitHubPullRequest {
	var large []GitHubPullRequest
	for _, pr := range prl.PullRequests {
		if pr.IsLargePR() {
			large = append(large, pr)
		}
	}
	return large
}
