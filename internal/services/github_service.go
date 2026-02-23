package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/krauzx/gitright/internal/github"
	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/repository"
)

type GitHubService struct {
	githubClient  *github.Client
	analyzer      *github.Analyzer
	repoCacheRepo *repository.RepositoryCacheRepository
}

func NewGitHubService(
	githubClient *github.Client,
	analyzer *github.Analyzer,
	repoCacheRepo *repository.RepositoryCacheRepository,
) *GitHubService {
	return &GitHubService{
		githubClient:  githubClient,
		analyzer:      analyzer,
		repoCacheRepo: repoCacheRepo,
	}
}

// ListUserRepositories fetches all repositories for a user with caching
func (s *GitHubService) ListUserRepositories(ctx context.Context, userID int64, accessToken string, includePrivate bool) ([]*models.Repository, error) {
	cachedRepos, err := s.repoCacheRepo.GetRepositoryList(ctx, userID, includePrivate)
	if err == nil && cachedRepos != nil {
		return cachedRepos, nil
	}

	githubRepos, err := s.githubClient.ListRepositories(ctx, accessToken, includePrivate)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	repos := make([]*models.Repository, 0, len(githubRepos))
	for _, gr := range githubRepos {
		if !includePrivate && gr.GetPrivate() {
			continue
		}

		repo := &models.Repository{
			ID:              gr.GetID(),
			GitHubID:        gr.GetID(),
			Name:            gr.GetName(),
			FullName:        gr.GetFullName(),
			Description:     gr.GetDescription(),
			Private:         gr.GetPrivate(),
			Fork:            gr.GetFork(),
			Language:        gr.GetLanguage(),
			StargazersCount: gr.GetStargazersCount(),
			ForksCount:      gr.GetForksCount(),
			OpenIssuesCount: gr.GetOpenIssuesCount(),
			DefaultBranch:   gr.GetDefaultBranch(),
			Topics:          gr.Topics,
			HTMLURL:         gr.GetHTMLURL(),
			CloneURL:        gr.GetCloneURL(),
			CreatedAt:       gr.GetCreatedAt().Time,
			UpdatedAt:       gr.GetUpdatedAt().Time,
			PushedAt:        gr.GetPushedAt().Time,
		}
		repos = append(repos, repo)
	}

	if err := s.repoCacheRepo.SetRepositoryList(ctx, userID, includePrivate, repos); err != nil {
		slog.Warn("Failed to cache repository list", "userID", userID, "error", err)
	}

	return repos, nil
}

// AnalyzeRepository performs deep analysis on a repository
func (s *GitHubService) AnalyzeRepository(ctx context.Context, accessToken, owner, repo string) (*models.RepositoryAnalysis, error) {
	repoInfo, err := s.githubClient.GetRepository(ctx, accessToken, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	githubID := repoInfo.GetID()
	fullName := fmt.Sprintf("%s/%s", owner, repo)

	cachedAnalysis, err := s.repoCacheRepo.GetRepositoryAnalysis(ctx, githubID)
	if err == nil && cachedAnalysis != nil {
		return cachedAnalysis, nil
	}

	analysis, err := s.analyzer.AnalyzeRepository(ctx, accessToken, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze repository: %w", err)
	}

	if err := s.repoCacheRepo.SetRepositoryAnalysis(ctx, githubID, fullName, analysis); err != nil {
		slog.Warn("Failed to cache repository analysis", "repo", fullName, "error", err)
	}

	return analysis, nil
}

// GetRepository fetches a single repository details
func (s *GitHubService) GetRepository(ctx context.Context, accessToken, owner, repo string) (*models.Repository, error) {
	gr, err := s.githubClient.GetRepository(ctx, accessToken, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	repository := &models.Repository{
		ID:              gr.GetID(),
		GitHubID:        gr.GetID(),
		Name:            gr.GetName(),
		FullName:        gr.GetFullName(),
		Description:     gr.GetDescription(),
		Private:         gr.GetPrivate(),
		Fork:            gr.GetFork(),
		Language:        gr.GetLanguage(),
		StargazersCount: gr.GetStargazersCount(),
		ForksCount:      gr.GetForksCount(),
		OpenIssuesCount: gr.GetOpenIssuesCount(),
		DefaultBranch:   gr.GetDefaultBranch(),
		Topics:          gr.Topics,
		HTMLURL:         gr.GetHTMLURL(),
		CloneURL:        gr.GetCloneURL(),
		CreatedAt:       gr.GetCreatedAt().Time,
		UpdatedAt:       gr.GetUpdatedAt().Time,
		PushedAt:        gr.GetPushedAt().Time,
	}

	return repository, nil
}

// DeployProfileREADME creates or updates the profile README on GitHub
func (s *GitHubService) DeployProfileREADME(ctx context.Context, accessToken, username, content string) error {
	currentSHA := ""
	sha, err := s.githubClient.GetProfileReadmeSHA(ctx, accessToken, username)
	if err == nil {
		currentSHA = sha
	}

	message := "Update profile README via GitRight"
	if currentSHA == "" {
		message = "Create profile README via GitRight"
	}

	if err := s.githubClient.CreateOrUpdateFile(ctx, accessToken, username, username, "README.md", message, content, currentSHA); err != nil {
		return fmt.Errorf("failed to deploy README: %w", err)
	}

	return nil
}

// ClearUserCache removes all cached data for a user
func (s *GitHubService) ClearUserCache(ctx context.Context, userID int64) error {
	if err := s.repoCacheRepo.InvalidateAllRepositoryLists(ctx, userID); err != nil {
		slog.Warn("Failed to invalidate repository list cache", "userID", userID, "error", err)
	}
	return nil
}

// ValidateRepositoryAccess checks if user has access to a repository
func (s *GitHubService) ValidateRepositoryAccess(ctx context.Context, accessToken, owner, repo string) error {
	_, err := s.githubClient.GetRepository(ctx, accessToken, owner, repo)
	if err != nil {
		return fmt.Errorf("access denied or repository not found: %w", err)
	}
	return nil
}

// BatchAnalyzeRepositories analyzes multiple repositories. It returns partial
// results when individual repositories fail; the caller receives both the
// successful analyses and a joined error listing every failure.
func (s *GitHubService) BatchAnalyzeRepositories(ctx context.Context, accessToken string, repos []string) (map[string]*models.RepositoryAnalysis, error) {
	results := make(map[string]*models.RepositoryAnalysis, len(repos))
	var errs []error

	for _, fullName := range repos {
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			errs = append(errs, fmt.Errorf("invalid repository name %q: must be owner/repo", fullName))
			continue
		}

		owner, repo := parts[0], parts[1]

		analysis, err := s.AnalyzeRepository(ctx, accessToken, owner, repo)
		if err != nil {
			slog.Warn("Failed to analyze repository in batch", "repo", fullName, "error", err)
			errs = append(errs, fmt.Errorf("repository %q: %w", fullName, err))
			continue
		}

		results[fullName] = analysis
	}

	return results, errors.Join(errs...)
}
