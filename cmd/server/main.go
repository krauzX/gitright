package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/krauzx/gitright/internal/config"
	"github.com/krauzx/gitright/internal/github"
	"github.com/krauzx/gitright/internal/handlers"
	"github.com/krauzx/gitright/internal/llm"
	"github.com/krauzx/gitright/internal/repository"
	"github.com/krauzx/gitright/internal/routes"
	"github.com/krauzx/gitright/internal/services"
	"github.com/krauzx/gitright/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logLevel := logger.ParseLevel(cfg.LogLevel)
	logger.Setup(logLevel, cfg.Environment)

	slog.Info("Starting GitRight server",
		"version", "1.0.0",
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	profileCacheRepo := repository.NewProfileCacheRepository(db)
	repoCacheRepo := repository.NewRepositoryCacheRepository(db)

	githubClient := github.NewClient(cfg.GitHub)
	githubAnalyzer := github.NewAnalyzer(githubClient)

	contentGenerator, err := llm.NewContentGenerator(cfg.GoogleAI)
	if err != nil {
		slog.Error("Failed to initialize content generator", "error", err)
		os.Exit(1)
	}

	authService := services.NewAuthService(githubClient, userRepo, sessionRepo)
	githubService := services.NewGitHubService(githubClient, githubAnalyzer, repoCacheRepo)
	profileService := services.NewProfileService(contentGenerator, projectRepo, githubService, profileCacheRepo)

	authHandler := handlers.NewAuthHandler(
		authService,
		"https://github.com",
		cfg.GitHub.ClientID,
		cfg.GitHub.RedirectURI,
		cfg.FrontendURL,
		cfg.Session.Secret,
		cfg.GitHub.Scopes,
	)
	githubHandler := handlers.NewGitHubHandler(githubService)
	profileHandler := handlers.NewProfileHandler(profileService)
	healthHandler := handlers.NewHealthHandler(db)
	wsHandler := handlers.NewWebSocketHandler(profileService, cfg.CORS.AllowedOrigins)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(logger.Middleware())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), cfg.HTTPTimeout)
			defer cancel()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	}))
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit.RequestsPerMinute)),
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;",
	}))

	routes.RegisterRoutes(e, authHandler, githubHandler, profileHandler, healthHandler, wsHandler, userRepo, sessionRepo, cfg.Session.Secret)

	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		slog.Info("Server listening", "address", addr)

		if cfg.Security.EnableHTTPS {
			if err := e.StartTLS(addr, cfg.Security.CertFile, cfg.Security.KeyFile); err != nil && err != http.ErrServerClosed {
				slog.Error("Failed to start HTTPS server", "error", err)
				os.Exit(1)
			}
		} else {
			if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
				slog.Error("Failed to start HTTP server", "error", err)
				os.Exit(1)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}
