package handlers

import (
	"net/http"

	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/services"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) Generate(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var req models.ContentGenerationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.profileService.GenerateProfile(ctx, &req, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) Deploy(c echo.Context) error {
	ctx := c.Request().Context()

	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var req models.ContentGenerationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.profileService.GenerateProfile(ctx, &req, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := h.profileService.DeployProfile(ctx, user.AccessToken, user.Username, response.Markdown, response.BannerSVG); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Profile deployed successfully",
		"url":     "https://github.com/" + user.Username,
	})
}


