package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateOAuthState(ctx context.Context, state string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (id, state_type, state_value, expires_at)
		VALUES ($1, 'oauth_state', 'valid', $2)
		ON CONFLICT (id) DO UPDATE
		SET expires_at = $2, state_value = 'valid'
	`
	_, err := r.db.ExecContext(ctx, query, "oauth_state:"+state, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create OAuth state: %w", err)
	}
	return nil
}

func (r *SessionRepository) ValidateOAuthState(ctx context.Context, state string) (bool, error) {
	query := `
		SELECT state_value FROM sessions
		WHERE id = $1
		  AND state_type = 'oauth_state'
		  AND expires_at > NOW()
	`
	var value string
	err := r.db.QueryRowContext(ctx, query, "oauth_state:"+state).Scan(&value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to validate OAuth state: %w", err)
	}
	return value == "valid", nil
}

func (r *SessionRepository) DeleteOAuthState(ctx context.Context, state string) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1 AND state_type = 'oauth_state'
	`
	_, err := r.db.ExecContext(ctx, query, "oauth_state:"+state)
	if err != nil {
		return fmt.Errorf("failed to delete OAuth state: %w", err)
	}
	return nil
}

func (r *SessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `SELECT cleanup_expired_data()`)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	return nil
}

func (r *SessionRepository) RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (id, state_type, state_value, expires_at)
		VALUES ($1, 'revoked_token', 'revoked', $2)
		ON CONFLICT (id) DO UPDATE
		SET expires_at = $2, state_value = 'revoked'
	`
	_, err := r.db.ExecContext(ctx, query, "revoked:"+jti, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

func (r *SessionRepository) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	query := `
		SELECT 1 FROM sessions
		WHERE id = $1
		  AND state_type = 'revoked_token'
		  AND expires_at > NOW()
		LIMIT 1
	`
	var dummy int
	err := r.db.QueryRowContext(ctx, query, "revoked:"+jti).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check token revocation: %w", err)
	}
	return true, nil
}

func (r *SessionRepository) GetStats(ctx context.Context) (map[string]int64, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE state_type = 'oauth_state') AS oauth_states_count,
			COUNT(*) FILTER (WHERE expires_at < NOW()) AS expired_count
		FROM sessions
	`
	var oauthStatesCount, expiredCount int64
	err := r.db.QueryRowContext(ctx, query).Scan(&oauthStatesCount, &expiredCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get session stats: %w", err)
	}
	return map[string]int64{
		"oauth_states": oauthStatesCount,
		"expired":      expiredCount,
	}, nil
}
