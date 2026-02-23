package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/krauzx/gitright/internal/middleware"
	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService   *services.AuthService
	githubBaseURL string
	clientID      string
	redirectURI   string
	frontendURL   string
	jwtSecret     string
	scopes        []string
}

func NewAuthHandler(
	authService *services.AuthService,
	githubBaseURL, clientID, redirectURI, frontendURL, jwtSecret string,
	scopes []string,
) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		githubBaseURL: githubBaseURL,
		clientID:      clientID,
		redirectURI:   redirectURI,
		frontendURL:   frontendURL,
		jwtSecret:     jwtSecret,
		scopes:        scopes,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	state, err := h.authService.GenerateOAuthState(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate auth state")
	}

	authURL := h.githubBaseURL + "/login/oauth/authorize" +
		"?client_id=" + h.clientID +
		"&redirect_uri=" + h.redirectURI +
		"&scope=" + strings.Join(h.scopes, ",") +
		"&state=" + state

	return c.JSON(http.StatusOK, map[string]interface{}{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *AuthHandler) Callback(c echo.Context) error {
	ctx := c.Request().Context()

	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing code or state")
	}

	if err := h.authService.ValidateOAuthState(ctx, state); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid or expired state: %v", err))
	}

	user, _, err := h.authService.HandleCallback(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Authentication failed: %v", err))
	}

	// 24-hour JWT — use this as the Bearer token for all subsequent requests.
	jwtToken, err := middleware.GenerateJWT(user.ID, user.Username, h.jwtSecret, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	if err := h.authService.DeleteOAuthState(ctx, state); err != nil {
		c.Logger().Warn("Failed to delete OAuth state after successful auth: ", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": jwtToken,
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	return c.JSON(http.StatusOK, user)
}

// Logout revokes the JWT by recording its JTI in the blocklist until it would
// have naturally expired.
func (h *AuthHandler) Logout(c echo.Context) error {
	ctx := c.Request().Context()

	jti, ok := c.Get("jwt_jti").(string)
	if !ok || jti == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	exp, ok := c.Get("jwt_exp").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := h.authService.RevokeSession(ctx, jti, time.Unix(exp, 0)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke session")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
