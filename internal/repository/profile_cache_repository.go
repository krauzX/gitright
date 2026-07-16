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


func GetCacheKey(username, targetRole, toneOfVoice string, projectCount int) string {
	return fmt.Sprintf("profile:v3:%s:%s:%s:%d", username, targetRole, toneOfVoice, projectCount)
}
