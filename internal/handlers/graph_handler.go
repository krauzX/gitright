package handlers

import (
	"net/http"

	"github.com/krauzx/gitright/internal/graph"
	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/services"
	"github.com/labstack/echo/v4"
)

type GraphHandler struct {
	autoImportService *services.AutoImportService
}

func NewGraphHandler(autoImportService *services.AutoImportService) *GraphHandler {
	return &GraphHandler{autoImportService: autoImportService}
}

func (h *GraphHandler) GenerateTelemetryGraph(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	profile, err := h.autoImportService.FetchDeepProfile(
		c.Request().Context(),
		user.AccessToken,
		user.Username,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch GitHub data")
	}

	g := graph.BuildFromProfile(profile)
	svg := graph.RenderSVG(g, graph.GraphConfig{
		Width:       800,
		Height:      600,
		BgColor:     "#0d1117",
		NodeColors:  map[string]string{"profile": "#f59e0b", "systems": "#06b6d4", "data": "#10b981", "frontend": "#3b82f6", "other": "#8b5cf6", "repository": "#6366f1", "edge": "#334155", "default": "#71717a"},
		MinNodeSize: 5,
		MaxNodeSize: 25,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"svg":    svg,
		"nodes":  len(g.Nodes),
		"edges":  len(g.Edges),
	})
}

func (h *GraphHandler) GenerateContributionGraph(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	profile, err := h.autoImportService.FetchDeepProfile(
		c.Request().Context(),
		user.AccessToken,
		user.Username,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch GitHub data")
	}

	svg := graph.RenderContributionSVG(profile.ContributionWeeks, graph.GraphConfig{
		Width:   720,
		Height:  120,
		BgColor: "#0d1117",
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"svg":  svg,
		"days": profile.TotalContributions,
	})
}
