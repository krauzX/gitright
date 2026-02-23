package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/krauzx/gitright/internal/models"
)

type RepositoryCacheRepository struct {
	db *sql.DB
}

func NewRepositoryCacheRepository(db *sql.DB) *RepositoryCacheRepository {
	return &RepositoryCacheRepository{db: db}
}

func (r *RepositoryCacheRepository) GetRepositoryList(ctx context.Context, userID int64, includePrivate bool) ([]*models.Repository, error) {
	query := `
		SELECT repositories
		FROM repository_list_cache
		WHERE user_id = $1
		  AND include_private = $2
		  AND expires_at > NOW()
		LIMIT 1
	`

	var reposJSON []byte
	err := r.db.QueryRowContext(ctx, query, userID, includePrivate).Scan(&reposJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cached repository list: %w", err)
	}

	var repos []*models.Repository
	if err := json.Unmarshal(reposJSON, &repos); err != nil {
		return nil, nil
	}

	return repos, nil
}

func (r *RepositoryCacheRepository) SetRepositoryList(ctx context.Context, userID int64, includePrivate bool, repos []*models.Repository) error {
	reposJSON, err := json.Marshal(repos)
	if err != nil {
		return fmt.Errorf("failed to marshal repositories: %w", err)
	}

	query := `
		INSERT INTO repository_list_cache
			(user_id, include_private, repositories, expires_at)
		VALUES
			($1, $2, $3, NOW() + INTERVAL '5 minutes')
		ON CONFLICT (user_id, include_private) DO UPDATE
		SET
			repositories = EXCLUDED.repositories,
			cached_at = NOW(),
			expires_at = NOW() + INTERVAL '5 minutes'
	`

	_, err = r.db.ExecContext(ctx, query, userID, includePrivate, reposJSON)
	if err != nil {
		return fmt.Errorf("failed to set repository list cache: %w", err)
	}

	return nil
}

func (r *RepositoryCacheRepository) InvalidateRepositoryList(ctx context.Context, userID int64, includePrivate bool) error {
	query := `DELETE FROM repository_list_cache WHERE user_id = $1 AND include_private = $2`
	_, err := r.db.ExecContext(ctx, query, userID, includePrivate)
	if err != nil {
		return fmt.Errorf("failed to invalidate repository list cache: %w", err)
	}
	return nil
}

func (r *RepositoryCacheRepository) InvalidateAllRepositoryLists(ctx context.Context, userID int64) error {
	query := `DELETE FROM repository_list_cache WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate all repository lists: %w", err)
	}
	return nil
}

func (r *RepositoryCacheRepository) GetRepositoryAnalysis(ctx context.Context, githubID int64) (*models.RepositoryAnalysis, error) {
	query := `
		SELECT
			full_name,
			languages,
			dependencies,
			key_files,
			commit_count,
			contributor_count
		FROM repository_analysis_cache
		WHERE github_id = $1
		  AND expires_at > NOW()
		LIMIT 1
	`

	var fullName string
	var languagesJSON, dependenciesJSON, keyFilesJSON []byte
	var commitCount, contributorCount int

	err := r.db.QueryRowContext(ctx, query, githubID).Scan(
		&fullName,
		&languagesJSON,
		&dependenciesJSON,
		&keyFilesJSON,
		&commitCount,
		&contributorCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cached repository analysis: %w", err)
	}

	var languages map[string]int
	var dependencies map[string][]string
	var keyFiles map[string]string

	if err := json.Unmarshal(languagesJSON, &languages); err != nil {
		return nil, nil
	}
	if err := json.Unmarshal(dependenciesJSON, &dependencies); err != nil {
		return nil, nil
	}
	if err := json.Unmarshal(keyFilesJSON, &keyFiles); err != nil {
		return nil, nil
	}

	return &models.RepositoryAnalysis{
		Languages:        languages,
		Dependencies:     dependencies,
		KeyFiles:         keyFiles,
		CommitCount:      commitCount,
		ContributorCount: contributorCount,
	}, nil
}

func (r *RepositoryCacheRepository) SetRepositoryAnalysis(ctx context.Context, githubID int64, fullName string, analysis *models.RepositoryAnalysis) error {
	languagesJSON, err := json.Marshal(analysis.Languages)
	if err != nil {
		return fmt.Errorf("failed to marshal languages: %w", err)
	}

	dependenciesJSON, err := json.Marshal(analysis.Dependencies)
	if err != nil {
		return fmt.Errorf("failed to marshal dependencies: %w", err)
	}

	keyFilesJSON, err := json.Marshal(analysis.KeyFiles)
	if err != nil {
		return fmt.Errorf("failed to marshal key files: %w", err)
	}

	query := `
		INSERT INTO repository_analysis_cache
			(github_id, full_name, languages, dependencies, key_files, commit_count, contributor_count, expires_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, NOW() + INTERVAL '7 days')
		ON CONFLICT (github_id) DO UPDATE
		SET
			full_name = EXCLUDED.full_name,
			languages = EXCLUDED.languages,
			dependencies = EXCLUDED.dependencies,
			key_files = EXCLUDED.key_files,
			commit_count = EXCLUDED.commit_count,
			contributor_count = EXCLUDED.contributor_count,
			analyzed_at = NOW(),
			expires_at = NOW() + INTERVAL '7 days'
	`

	_, err = r.db.ExecContext(ctx, query, githubID, fullName, languagesJSON, dependenciesJSON, keyFilesJSON, analysis.CommitCount, analysis.ContributorCount)
	if err != nil {
		return fmt.Errorf("failed to set repository analysis cache: %w", err)
	}

	return nil
}

func (r *RepositoryCacheRepository) InvalidateRepositoryAnalysis(ctx context.Context, githubID int64) error {
	query := `DELETE FROM repository_analysis_cache WHERE github_id = $1`
	_, err := r.db.ExecContext(ctx, query, githubID)
	if err != nil {
		return fmt.Errorf("failed to invalidate repository analysis: %w", err)
	}
	return nil
}

func (r *RepositoryCacheRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM repository_list_cache WHERE expires_at > NOW()) AS cached_lists,
			(SELECT COUNT(*) FROM repository_list_cache WHERE expires_at <= NOW()) AS expired_lists,
			(SELECT COUNT(*) FROM repository_analysis_cache WHERE expires_at > NOW()) AS cached_analyses,
			(SELECT COUNT(*) FROM repository_analysis_cache WHERE expires_at <= NOW()) AS expired_analyses
	`

	var cachedLists, expiredLists, cachedAnalyses, expiredAnalyses int64
	err := r.db.QueryRowContext(ctx, query).Scan(&cachedLists, &expiredLists, &cachedAnalyses, &expiredAnalyses)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}

	return map[string]interface{}{
		"cached_repo_lists":     cachedLists,
		"expired_repo_lists":    expiredLists,
		"cached_analyses":       cachedAnalyses,
		"expired_analyses":      expiredAnalyses,
		"total_cached_entries":  cachedLists + cachedAnalyses,
		"total_expired_entries": expiredLists + expiredAnalyses,
	}, nil
}
