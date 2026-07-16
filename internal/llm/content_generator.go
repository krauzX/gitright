package llm

import (
	"log/slog"

	"github.com/krauzx/gitright/internal/config"
)

type ContentGenerator struct {
	client *GeminiClient
}

func NewContentGenerator(cfg config.GoogleAIConfig) *ContentGenerator {
	if cfg.APIKey == "" {
		slog.Info("No Gemini API key configured - BYOK mode enabled (users provide their own keys)")
		return &ContentGenerator{client: nil}
	}

	client, err := NewGeminiClient(cfg)
	if err != nil {
		slog.Warn("Failed to create default Gemini client, BYOK mode only", "error", err)
		return &ContentGenerator{client: nil}
	}

	return &ContentGenerator{client: client}
}
