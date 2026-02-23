package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/krauzx/gitright/internal/models"
)

type ProfileCacheRepository struct {
	db *sql.DB
}

func NewProfileCacheRepository(db *sql.DB) *ProfileCacheRepository {
	return &ProfileCacheRepository{db: db}
}

func (r *ProfileCacheRepository) Get(ctx context.Context, cacheKey string) (*models.ContentGenerationResponse, error) {
	query := `
		SELECT content
		FROM generated_profiles
		WHERE cache_key = $1
		  AND expires_at > NOW()
		LIMIT 1
	`

	var contentJSON string
	err := r.db.QueryRowContext(ctx, query, cacheKey).Scan(&contentJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cached profile: %w", err)
	}

	var response models.ContentGenerationResponse
	if err := json.Unmarshal([]byte(contentJSON), &response); err != nil {
		return nil, nil
	}

	go r.updateCacheStats(context.Background(), cacheKey)

	return &response, nil
}

func (r *ProfileCacheRepository) Set(ctx context.Context, userID, configID int64, cacheKey string, response *models.ContentGenerationResponse, ttl time.Duration) error {
	contentJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	expiresAt := time.Now().Add(ttl)

	query := `
		INSERT INTO generated_profiles
			(user_id, config_id, content, markdown_preview, cache_key, expires_at, version)
		VALUES
			($1, $2, $3, $4, $5, $6, 1)
		ON CONFLICT (cache_key) DO UPDATE
		SET
			content = EXCLUDED.content,
			markdown_preview = EXCLUDED.markdown_preview,
			expires_at = EXCLUDED.expires_at,
			last_accessed_at = NOW()
	`

	_, err = r.db.ExecContext(ctx, query, userID, configID, string(contentJSON), response.Markdown, cacheKey, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to set cached profile: %w", err)
	}

	return nil
}

func (r *ProfileCacheRepository) Invalidate(ctx context.Context, cacheKey string) error {
	query := `DELETE FROM generated_profiles WHERE cache_key = $1`
	_, err := r.db.ExecContext(ctx, query, cacheKey)
	if err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}
	return nil
}

func (r *ProfileCacheRepository) InvalidateByUserID(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM generated_profiles
		WHERE user_id = $1
		  AND cache_key IS NOT NULL
		  AND NOT deployed
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate user cache: %w", err)
	}
	return nil
}

func (r *ProfileCacheRepository) updateCacheStats(ctx context.Context, cacheKey string) {
	query := `
		UPDATE generated_profiles
		SET
			cache_hit_count = cache_hit_count + 1,
			last_accessed_at = NOW()
		WHERE cache_key = $1
	`
	_, _ = r.db.ExecContext(ctx, query, cacheKey)
}

func (r *ProfileCacheRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE cache_key IS NOT NULL AND expires_at > NOW()) AS cached_profiles,
			COUNT(*) FILTER (WHERE cache_key IS NOT NULL AND expires_at <= NOW()) AS expired_profiles,
			COUNT(*) FILTER (WHERE deployed = TRUE) AS deployed_profiles,
			COALESCE(AVG(cache_hit_count) FILTER (WHERE cache_key IS NOT NULL), 0) AS avg_cache_hits,
			COALESCE(MAX(cache_hit_count) FILTER (WHERE cache_key IS NOT NULL), 0) AS max_cache_hits
		FROM generated_profiles
	`

	var cachedProfiles, expiredProfiles, deployedProfiles, maxCacheHits int64
	var avgCacheHits float64

	err := r.db.QueryRowContext(ctx, query).Scan(
		&cachedProfiles,
		&expiredProfiles,
		&deployedProfiles,
		&avgCacheHits,
		&maxCacheHits,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}

	var hitRate float64
	if cachedProfiles > 0 {
		hitRate = avgCacheHits / float64(cachedProfiles) * 100
	}

	return map[string]interface{}{
		"cached_profiles":   cachedProfiles,
		"expired_profiles":  expiredProfiles,
		"deployed_profiles": deployedProfiles,
		"avg_cache_hits":    avgCacheHits,
		"max_cache_hits":    maxCacheHits,
		"cache_hit_rate":    hitRate,
	}, nil
}

func GetCacheKey(username, targetRole, toneOfVoice string, projectCount int) string {
	return fmt.Sprintf("profile:v3:%s:%s:%s:%d", username, targetRole, toneOfVoice, projectCount)
}
