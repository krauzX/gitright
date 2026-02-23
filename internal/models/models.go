package models

import (
	"time"
)

// User represents a GitHub user
type User struct {
	ID             int64     `json:"id" db:"id"`
	GitHubID       int64     `json:"github_id" db:"github_id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	AvatarURL      string    `json:"avatar_url" db:"avatar_url"`
	Bio            string    `json:"bio" db:"bio"`
	Location       string    `json:"location" db:"location"`
	Company        string    `json:"company" db:"company"`
	Blog           string    `json:"blog" db:"blog"`
	AccessToken    string    `json:"-" db:"access_token"`
	RefreshToken   string    `json:"-" db:"refresh_token"`
	TokenExpiresAt time.Time `json:"-" db:"token_expires_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt    time.Time `json:"last_login_at" db:"last_login_at"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID              int64     `json:"id"`
	GitHubID        int64     `json:"github_id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	Private         bool      `json:"private"`
	Fork            bool      `json:"fork"`
	Language        string    `json:"language"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	DefaultBranch   string    `json:"default_branch"`
	Topics          []string  `json:"topics"`
	HTMLURL         string    `json:"html_url"`
	CloneURL        string    `json:"clone_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
}

// RepositoryAnalysis contains analyzed repository information
type RepositoryAnalysis struct {
	Repository       *Repository         `json:"repository"`
	Languages        map[string]int      `json:"languages"`
	Files            []string            `json:"files"`
	Dependencies     map[string][]string `json:"dependencies"`
	KeyFiles         map[string]string   `json:"key_files"`
	CommitCount      int                 `json:"commit_count"`
	ContributorCount int                 `json:"contributor_count"`
}

// Project represents a user's selected project for their profile
type Project struct {
	ID               int64     `json:"id" db:"id"`
	UserID           int64     `json:"user_id" db:"user_id"`
	GitHubID         int64     `json:"github_id" db:"github_id"` // Direct GitHub repo ID
	FullName         string    `json:"full_name" db:"full_name"` // owner/repo format
	Priority         int       `json:"priority" db:"priority"`
	FocusTag         string    `json:"focus_tag" db:"focus_tag"` // "best_performance", "team_project", "personal_favorite"
	CustomSummary    string    `json:"custom_summary" db:"custom_summary"`
	GeneratedSummary string    `json:"generated_summary" db:"generated_summary"`
	IncludeInProfile bool      `json:"include_in_profile" db:"include_in_profile"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ProfileConfig represents user's profile customization settings
type ProfileConfig struct {
	ID               int64              `json:"id" db:"id"`
	UserID           int64              `json:"user_id" db:"user_id"`
	TargetRole       string             `json:"target_role" db:"target_role"`
	SkillsEmphasis   []string           `json:"skills_emphasis" db:"skills_emphasis"`
	ToneOfVoice      string             `json:"tone_of_voice" db:"tone_of_voice"` // "professional", "friendly", "technical"
	TemplateID       string             `json:"template_id" db:"template_id"`     // "technical_deep_dive", "hiring_manager_scan", "community_contributor"
	ContactPrefs     ContactPreferences `json:"contact_prefs" db:"contact_prefs"`
	ShowPrivateRepos bool               `json:"show_private_repos" db:"show_private_repos"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`
}

// ContactPreferences represents contact preferences
type ContactPreferences struct {
	LinkedIn        string   `json:"linkedin"`
	PersonalWebsite string   `json:"personal_website"`
	Email           string   `json:"email"`
	Twitter         string   `json:"twitter"`
	PreferredOrder  []string `json:"preferred_order"`
}

// GeneratedProfile represents the final generated README content
type GeneratedProfile struct {
	ID              int64      `json:"id" db:"id"`
	UserID          int64      `json:"user_id" db:"user_id"`
	ConfigID        int64      `json:"config_id" db:"config_id"`
	Content         string     `json:"content" db:"content"`
	MarkdownPreview string     `json:"markdown_preview" db:"markdown_preview"`
	Deployed        bool       `json:"deployed" db:"deployed"`
	DeployedAt      *time.Time `json:"deployed_at,omitempty" db:"deployed_at"`
	Version         int        `json:"version" db:"version"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

type ContentGenerationRequest struct {
	TargetRole       string               `json:"target_role" validate:"required"`
	EmphasizedSkills []string             `json:"emphasized_skills"`
	ToneOfVoice      string               `json:"tone_of_voice" validate:"required"`
	ContactPrefs     ContactPreferences   `json:"contact_prefs"`
	Projects         []RepositoryAnalysis `json:"projects" validate:"required,min=1"`
	UserAPIKey       string               `json:"user_api_key" validate:"required"`
}

type ContentGenerationResponse struct {
	Markdown        string   `json:"markdown"`
	ExtractedSkills []string `json:"extracted_skills"`
	SuggestedBadges []Badge  `json:"suggested_badges"`
	Confidence      float64  `json:"confidence"`
}

// Badge represents a technology badge
type Badge struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Color string `json:"color"`
}

// ProfilePitch represents the generated introduction paragraph
type ProfilePitch struct {
	Content    string  `json:"content"`
	Confidence float64 `json:"confidence"`
}

// TemplateData represents data for template rendering
type TemplateData struct {
	User         *User            `json:"user"`
	Config       *ProfileConfig   `json:"config"`
	Projects     []ProjectSummary `json:"projects"`
	Skills       []string         `json:"skills"`
	Badges       []Badge          `json:"badges"`
	ProfilePitch *ProfilePitch    `json:"profile_pitch"`
}

// ProjectSummary combines project details with generated content
type ProjectSummary struct {
	Project    *Project    `json:"project"`
	Repository *Repository `json:"repository"`
	Summary    string      `json:"summary"`
	TechStack  []string    `json:"tech_stack"`
	Highlights []string    `json:"highlights"`
}
