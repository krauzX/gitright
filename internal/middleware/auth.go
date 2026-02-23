package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/krauzx/gitright/internal/repository"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	JTI       string `json:"jti"`
	ExpiresAt int64  `json:"exp"`
}

// AuthMiddleware validates JWT tokens, checks the JTI blocklist, and populates
// request context with user data for downstream handlers.
func AuthMiddleware(secret string, userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.ErrUnauthorized
			}

			token := parts[1]

			claims, err := validateJWT(token, secret)
			if err != nil {
				return echo.ErrUnauthorized
			}

			if time.Now().Unix() > claims.ExpiresAt {
				return echo.ErrUnauthorized
			}

			ctx := c.Request().Context()

			if claims.JTI != "" {
				revoked, err := sessionRepo.IsTokenRevoked(ctx, claims.JTI)
				if err != nil || revoked {
					return echo.ErrUnauthorized
				}
			}

			user, err := userRepo.GetByID(ctx, claims.UserID)
			if err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("user", user)
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("access_token", user.AccessToken)
			c.Set("jwt_jti", claims.JTI)
			c.Set("jwt_exp", claims.ExpiresAt)

			return next(c)
		}
	}
}

// GenerateJWT creates a signed HS256 JWT for the given user with a unique JTI.
func GenerateJWT(userID int64, username, secret string, expiresIn time.Duration) (string, error) {
	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", fmt.Errorf("failed to generate JTI: %w", err)
	}
	jti := hex.EncodeToString(jtiBytes)

	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		JTI:       jti,
		ExpiresAt: time.Now().Add(expiresIn).Unix(),
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	message := header + "." + payload
	return message + "." + createSignature(message, secret), nil
}

func validateJWT(token, secret string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	if parts[2] != createSignature(message, secret) {
		return nil, fmt.Errorf("invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func createSignature(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
