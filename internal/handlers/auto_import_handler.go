package handlers

import (
	"log/slog"
	"net/http"

	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/services"
	"github.com/labstack/echo/v4"
)

type AutoImportHandler struct {
	autoImportService *services.AutoImportService
}

func NewAutoImportHandler(
	autoImportService *services.AutoImportService,
) *AutoImportHandler {
	return &AutoImportHandler{
		autoImportService: autoImportService,
	}
}

func (h *AutoImportHandler) AutoImport(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	slog.Info("Auto-import started", "username", user.Username)

	data, err := h.autoImportService.FetchAndExtract(c.Request().Context(), user.AccessToken, user.Username)
	if err != nil {
		slog.Error("Failed to fetch GitHub data", "error", err, "username", user.Username)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch GitHub data: "+err.Error())
	}

	slog.Info("Auto-import completed", "username", user.Username, "skills", len(data.Skills), "contributions", data.TotalContributions)

	return c.JSON(http.StatusOK, data)
}
