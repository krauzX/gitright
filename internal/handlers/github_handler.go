package handlers

import (
	"net/http"

	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/services"
	"github.com/labstack/echo/v4"
)

type GitHubHandler struct {
	githubService *services.GitHubService
}

func NewGitHubHandler(githubService *services.GitHubService) *GitHubHandler {
	return &GitHubHandler{githubService: githubService}
}

func (h *GitHubHandler) ListRepositories(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	includePrivate := c.QueryParam("include_private") == "true"

	repos, err := h.githubService.ListUserRepositories(ctx, user.ID, user.AccessToken, includePrivate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch repositories")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"repositories": repos,
		"count":        len(repos),
	})
}

func (h *GitHubHandler) BatchAnalyze(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var req struct {
		Repositories []string `json:"repositories"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if len(req.Repositories) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one repository is required")
	}

	if len(req.Repositories) > 10 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maximum 10 repositories allowed")
	}

	results, err := h.githubService.BatchAnalyzeRepositories(ctx, user.AccessToken, req.Repositories)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to analyze repositories")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"analyses": results,
	})
}

func (h *GitHubHandler) ClearCache(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := h.githubService.ClearUserCache(ctx, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to clear cache")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Cache cleared successfully",
	})
}
