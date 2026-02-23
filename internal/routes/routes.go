package routes

import (
	"github.com/krauzx/gitright/internal/handlers"
	"github.com/krauzx/gitright/internal/middleware"
	"github.com/krauzx/gitright/internal/repository"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	authHandler *handlers.AuthHandler,
	githubHandler *handlers.GitHubHandler,
	profileHandler *handlers.ProfileHandler,
	healthHandler *handlers.HealthHandler,
	wsHandler *handlers.WebSocketHandler,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	jwtSecret string,
) {
	e.GET("/health", healthHandler.Health)
	e.GET("/health/ready", healthHandler.Ready)
	e.GET("/health/live", healthHandler.Live)

	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.GET("/login", authHandler.Login)
	auth.GET("/callback", authHandler.Callback)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret, userRepo, sessionRepo))

	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/me", authHandler.Me)

	gh := protected.Group("/github")
	gh.GET("/repositories", githubHandler.ListRepositories)
	gh.GET("/repositories/:owner/:repo", githubHandler.GetRepository)
	gh.GET("/repositories/:owner/:repo/analyze", githubHandler.AnalyzeRepository)
	gh.POST("/repositories/batch-analyze", githubHandler.BatchAnalyze)
	gh.DELETE("/cache", githubHandler.ClearCache)

	profile := protected.Group("/profile")
	profile.POST("/generate", profileHandler.Generate)
	profile.POST("/deploy", profileHandler.Deploy)
	profile.POST("/preview", profileHandler.Preview)
	profile.GET("/ws", wsHandler.HandleProfileGeneration)
}
