package handlers

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"

	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/services"
	"github.com/krauzx/gitright/internal/svg"
	"github.com/labstack/echo/v4"
)

type BannerHandler struct {
	autoImportService *services.AutoImportService
}

func NewBannerHandler(autoImportService *services.AutoImportService) *BannerHandler {
	return &BannerHandler{autoImportService: autoImportService}
}

func (h *BannerHandler) GenerateBanner(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	theme := c.QueryParam("theme")
	dark := theme != "light"

	data, err := h.autoImportService.FetchAndExtract(c.Request().Context(), user.AccessToken, user.Username)
	if err != nil {
		slog.Error("Failed to fetch GitHub data for banner", "error", err, "username", user.Username)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch GitHub data")
	}

	socialLinks := make([]svg.SocialLink, 0, len(data.SocialLinks))
	for _, sl := range data.SocialLinks {
		socialLinks = append(socialLinks, svg.SocialLink{
			Provider: sl.Provider,
			URL:      sl.URL,
		})
	}

	topRepos := make([]svg.TopRepoInfo, 0, len(data.TopRepositories))
	for _, r := range data.TopRepositories {
		topRepos = append(topRepos, svg.TopRepoInfo{
			Name:     r.Name,
			Stars:    r.Stars,
			Language: r.Language,
		})
	}

	bannerData := &svg.BannerData{
		Name:          data.Name,
		Login:         data.Login,
		Bio:           data.Bio,
		AvatarURL:     data.AvatarURL,
		Location:      data.Location,
		Company:       data.Company,
		Website:       data.WebsiteURL,
		Email:         data.Email,
		Skills:        data.Skills,
		Languages:     data.LanguageBreakdown,
		Contributions: data.TotalContributions,
		Followers:     data.Followers,
		Repos:         len(data.TopRepositories),
		PRReviews:     data.ReviewContributions,
		Issues:        data.IssueContributions,
		SocialLinks:   socialLinks,
		TopRepos:      topRepos,
	}

	var buf bytes.Buffer
	svg.GenerateBanner(&buf, bannerData, dark)

	c.Response().Header().Set("Content-Type", "image/svg+xml")
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	return c.String(http.StatusOK, buf.String())
}

func (h *BannerHandler) GenerateTimeline(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	theme := c.QueryParam("theme")
	dark := theme != "light"

	data, err := h.autoImportService.FetchAndExtract(c.Request().Context(), user.AccessToken, user.Username)
	if err != nil {
		slog.Error("Failed to fetch GitHub data for timeline", "error", err, "username", user.Username)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch GitHub data")
	}

	repos := make([]svg.TimelineRepo, 0, len(data.TopRepositories))
	for _, r := range data.TopRepositories {
		repos = append(repos, svg.TimelineRepo{
			Name:      r.Name,
			Language:  r.Language,
			Stars:     r.Stars,
			CreatedAt: "",
		})
	}

	timelineData := &svg.TimelineData{
		Name:          data.Name,
		Login:         data.Login,
		Repos:         repos,
		TotalContribs: data.TotalContributions,
		MemberSince:   time.Now().Format("2006"),
		Languages:     data.Skills,
	}

	var buf bytes.Buffer
	svg.GenerateTimeline(&buf, timelineData, dark)

	c.Response().Header().Set("Content-Type", "image/svg+xml")
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	return c.String(http.StatusOK, buf.String())
}
