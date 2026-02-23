package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/krauzx/gitright/internal/github"
	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/repository"
)

type AuthService struct {
	githubClient *github.Client
	userRepo     *repository.UserRepository
	sessionRepo  *repository.SessionRepository
}

func NewAuthService(
	githubClient *github.Client,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
) *AuthService {
	return &AuthService{
		githubClient: githubClient,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
	}
}

func (s *AuthService) GenerateOAuthState(ctx context.Context) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	state := base64.URLEncoding.EncodeToString(b)
	expiresAt := time.Now().Add(10 * time.Minute)

	if err := s.sessionRepo.CreateOAuthState(ctx, state, expiresAt); err != nil {
		return "", fmt.Errorf("failed to store state: %w", err)
	}

	return state, nil
}

func (s *AuthService) ValidateOAuthState(ctx context.Context, state string) error {
	valid, err := s.sessionRepo.ValidateOAuthState(ctx, state)
	if err != nil {
		return fmt.Errorf("state validation error: %w", err)
	}
	if !valid {
		return fmt.Errorf("state not found or expired")
	}
	return nil
}

func (s *AuthService) DeleteOAuthState(ctx context.Context, state string) error {
	return s.sessionRepo.DeleteOAuthState(ctx, state)
}

func (s *AuthService) HandleCallback(ctx context.Context, code string) (*models.User, string, error) {
	token, err := s.githubClient.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code: %w", err)
	}

	githubUser, err := s.githubClient.GetUser(ctx, token.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}

	emails, err := s.githubClient.GetUserEmails(ctx, token.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user emails: %w", err)
	}

	var primaryEmail string
	for _, email := range emails {
		if email.GetPrimary() && email.GetVerified() {
			primaryEmail = email.GetEmail()
			break
		}
	}

	existingUser, err := s.userRepo.GetByGitHubID(ctx, githubUser.GetID())
	if err == nil && existingUser != nil {
		existingUser.Username = githubUser.GetLogin()
		existingUser.Email = primaryEmail
		existingUser.AvatarURL = githubUser.GetAvatarURL()
		existingUser.Bio = githubUser.GetBio()
		existingUser.Location = githubUser.GetLocation()
		existingUser.Company = githubUser.GetCompany()
		existingUser.Blog = githubUser.GetBlog()
		existingUser.AccessToken = token.AccessToken
		existingUser.RefreshToken = token.RefreshToken
		existingUser.TokenExpiresAt = token.Expiry
		existingUser.LastLoginAt = time.Now()

		if err := s.userRepo.Update(ctx, existingUser); err != nil {
			return nil, "", fmt.Errorf("failed to update user: %w", err)
		}
		return existingUser, token.AccessToken, nil
	}

	newUser := &models.User{
		GitHubID:       githubUser.GetID(),
		Username:       githubUser.GetLogin(),
		Email:          primaryEmail,
		AvatarURL:      githubUser.GetAvatarURL(),
		Bio:            githubUser.GetBio(),
		Location:       githubUser.GetLocation(),
		Company:        githubUser.GetCompany(),
		Blog:           githubUser.GetBlog(),
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenExpiresAt: token.Expiry,
		LastLoginAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, token.AccessToken, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, jti string, expiresAt time.Time) error {
	return s.sessionRepo.RevokeToken(ctx, jti, expiresAt)
}
