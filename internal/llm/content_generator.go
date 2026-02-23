package llm

import (
	"github.com/krauzx/gitright/internal/config"
)

type ContentGenerator struct {
	client *GeminiClient
}

func NewContentGenerator(cfg config.GoogleAIConfig) (*ContentGenerator, error) {
	client, err := NewGeminiClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ContentGenerator{client: client}, nil
}
