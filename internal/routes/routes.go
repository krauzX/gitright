package routes

import (
	"time"

	"github.com/krauzx/gitright/internal/handlers"
	"github.com/krauzx/gitright/internal/middleware"
	"github.com/krauzx/gitright/internal/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

func RegisterRoutes(
	e *echo.Echo,
	authHandler *handlers.AuthHandler,
	githubHandler *handlers.GitHubHandler,
	profileHandler *handlers.ProfileHandler,
	autoImportHandler *handlers.AutoImportHandler,
	graphHandler *handlers.GraphHandler,
	bannerHandler *handlers.BannerHandler,
	healthHandler *handlers.HealthHandler,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	jwtSecret string,
) {
	e.GET("/health", healthHandler.Health)
	e.GET("/health/ready", healthHandler.Ready)
	e.GET("/health/live", healthHandler.Live)

	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	loginLimiter := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      rate.Limit(1),
		Burst:     3,
		ExpiresIn: 10 * time.Minute,
	})
	auth.GET("/login", authHandler.Login, middleware.RateLimiter(loginLimiter))
	auth.GET("/callback", authHandler.Callback)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret, userRepo, sessionRepo))

	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/me", authHandler.Me)

	gh := protected.Group("/github")
	gh.GET("/repositories", githubHandler.ListRepositories)
	gh.POST("/repositories/batch-analyze", githubHandler.BatchAnalyze)
	gh.DELETE("/cache", githubHandler.ClearCache)

	profile := protected.Group("/profile")
	profile.POST("/generate", profileHandler.Generate)
	profile.POST("/deploy", profileHandler.Deploy)
	profile.POST("/auto-import", autoImportHandler.AutoImport)
	profile.GET("/banner", bannerHandler.GenerateBanner)
	profile.GET("/timeline", bannerHandler.GenerateTimeline)

	graph := protected.Group("/graph")
	graph.GET("/telemetry", graphHandler.GenerateTelemetryGraph)
	graph.GET("/contributions", graphHandler.GenerateContributionGraph)
}
