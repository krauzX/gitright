package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/krauzx/gitright/internal/repository"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

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

			tokenString := parts[1]

			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
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
			c.Set("jwt_jti", claims.JTI)
			if claims.ExpiresAt != nil {
				c.Set("jwt_exp", claims.ExpiresAt.Unix())
			}

			return next(c)
		}
	}
}

func GenerateJWT(userID int64, username, secret string, expiresIn time.Duration) (string, error) {
	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", fmt.Errorf("failed to generate JTI: %w", err)
	}

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		JTI:      hex.EncodeToString(jtiBytes),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
