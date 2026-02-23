package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	db HealthChecker
}

type HealthChecker interface {
	Ping() error
}

func NewHealthHandler(db HealthChecker) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c echo.Context) error {
	health := map[string]interface{}{
		"status": "healthy",
		"services": map[string]string{
			"database": "healthy",
		},
	}

	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			health["status"] = "unhealthy"
			health["services"].(map[string]string)["database"] = "unhealthy"
		}
	}

	statusCode := http.StatusOK
	if health["status"] == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, health)
}

func (h *HealthHandler) Ready(c echo.Context) error {
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
				"reason": "database unavailable",
			})
		}
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
}

func (h *HealthHandler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "alive"})
}
