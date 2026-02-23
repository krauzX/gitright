package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/services"
	"github.com/labstack/echo/v4"
)

type ProgressUpdate struct {
	Stage    string  `json:"stage"`
	Progress float64 `json:"progress"`
	Message  string  `json:"message"`
	Error    string  `json:"error,omitempty"`
}

type WebSocketHandler struct {
	profileService *services.ProfileService
	upgrader       websocket.Upgrader
}

func NewWebSocketHandler(profileService *services.ProfileService, allowedOrigins []string) *WebSocketHandler {
	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = struct{}{}
	}

	return &WebSocketHandler{
		profileService: profileService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return false
				}
				_, allowed := originSet[origin]
				return allowed
			},
		},
	}
}

func (h *WebSocketHandler) HandleProfileGeneration(c echo.Context) error {
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		slog.Error("Failed to upgrade to WebSocket", "error", err)
		return err
	}
	defer ws.Close()

	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		h.sendError(ws, "User not found in context")
		return nil
	}

	var req models.ContentGenerationRequest
	if err := ws.ReadJSON(&req); err != nil {
		h.sendError(ws, "Invalid request format")
		return nil
	}

	ctx := c.Request().Context()

	progressCh := make(chan ProgressUpdate, 10)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case update, open := <-progressCh:
				if !open {
					return
				}
				if err := ws.WriteJSON(update); err != nil {
					slog.Error("Failed to send progress update", "error", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	response, err := h.generateWithProgress(ctx, &req, user, progressCh)
	close(progressCh)
	<-done

	if err != nil {
		h.sendError(ws, fmt.Sprintf("Profile generation failed: %v", err))
		return nil
	}

	if err := ws.WriteJSON(map[string]interface{}{
		"stage":    "complete",
		"progress": 1.0,
		"message":  "Profile generated successfully",
		"result":   response,
	}); err != nil {
		slog.Error("Failed to send final result", "error", err)
	}

	return nil
}

func (h *WebSocketHandler) generateWithProgress(
	ctx context.Context,
	req *models.ContentGenerationRequest,
	user *models.User,
	progressCh chan<- ProgressUpdate,
) (*models.ContentGenerationResponse, error) {
	send := func(stage string, pct float64, msg string) {
		select {
		case progressCh <- ProgressUpdate{Stage: stage, Progress: pct, Message: msg}:
		default:
		}
	}

	send("init", 0.0, "Starting profile generation...")
	send("analyzing", 0.1, fmt.Sprintf("Preparing %d project(s) for analysis...", len(req.Projects)))
	send("generating", 0.4, "Sending to AI — this may take up to 30 seconds...")

	response, err := h.profileService.GenerateProfile(ctx, req, user)
	if err != nil {
		return nil, err
	}

	send("finalizing", 0.95, "Assembling final profile...")

	return response, nil
}

func (h *WebSocketHandler) sendError(ws *websocket.Conn, message string) {
	update := ProgressUpdate{Stage: "error", Message: message, Error: message}
	if err := ws.WriteJSON(update); err != nil {
		slog.Error("Failed to send error via WebSocket", "error", err)
	}
}
