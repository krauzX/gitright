package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/krauzx/gitright/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (github_id, username, email, avatar_url, bio, location, company, blog, access_token, refresh_token, token_expires_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		user.GitHubID, user.Username, user.Email, user.AvatarURL, user.Bio,
		user.Location, user.Company, user.Blog, user.AccessToken, user.RefreshToken,
		user.TokenExpiresAt, time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, github_id, username, email, avatar_url, bio, location, company, blog,
		       access_token, refresh_token, token_expires_at, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.Email, &user.AvatarURL,
		&user.Bio, &user.Location, &user.Company, &user.Blog, &user.AccessToken,
		&user.RefreshToken, &user.TokenExpiresAt, &user.CreatedAt, &user.UpdatedAt,
		&user.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *UserRepository) GetByGitHubID(ctx context.Context, githubID int64) (*models.User, error) {
	query := `
		SELECT id, github_id, username, email, avatar_url, bio, location, company, blog,
		       access_token, refresh_token, token_expires_at, created_at, updated_at, last_login_at
		FROM users
		WHERE github_id = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, githubID).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.Email, &user.AvatarURL,
		&user.Bio, &user.Location, &user.Company, &user.Blog, &user.AccessToken,
		&user.RefreshToken, &user.TokenExpiresAt, &user.CreatedAt, &user.UpdatedAt,
		&user.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, avatar_url = $3, bio = $4, location = $5,
		    company = $6, blog = $7, access_token = $8, refresh_token = $9,
		    token_expires_at = $10, updated_at = $11, last_login_at = $12
		WHERE id = $13
	`
	_, err := r.db.ExecContext(
		ctx, query,
		user.Username, user.Email, user.AvatarURL, user.Bio, user.Location,
		user.Company, user.Blog, user.AccessToken, user.RefreshToken,
		user.TokenExpiresAt, time.Now(), time.Now(), user.ID,
	)
	return err
}
