package handlers

import (
	"net/http"

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

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	includePrivate := c.QueryParam("include_private") == "true"

	repos, err := h.githubService.ListUserRepositories(ctx, userID, accessToken, includePrivate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch repositories")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"repositories": repos,
		"count":        len(repos),
	})
}

func (h *GitHubHandler) GetRepository(c echo.Context) error {
	ctx := c.Request().Context()

	accessToken, ok := c.Get("access_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Owner and repo are required")
	}

	repository, err := h.githubService.GetRepository(ctx, accessToken, owner, repo)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Repository not found")
	}

	return c.JSON(http.StatusOK, repository)
}

func (h *GitHubHandler) AnalyzeRepository(c echo.Context) error {
	ctx := c.Request().Context()

	accessToken, ok := c.Get("access_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Owner and repo are required")
	}

	analysis, err := h.githubService.AnalyzeRepository(ctx, accessToken, owner, repo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to analyze repository")
	}

	return c.JSON(http.StatusOK, analysis)
}

func (h *GitHubHandler) BatchAnalyze(c echo.Context) error {
	ctx := c.Request().Context()

	accessToken, ok := c.Get("access_token").(string)
	if !ok {
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

	results, err := h.githubService.BatchAnalyzeRepositories(ctx, accessToken, req.Repositories)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to analyze repositories")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"analyses": results,
	})
}

func (h *GitHubHandler) ClearCache(c echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := h.githubService.ClearUserCache(ctx, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to clear cache")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Cache cleared successfully",
	})
}
