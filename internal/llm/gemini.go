package llm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/krauzx/gitright/internal/config"
	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	config config.GoogleAIConfig
}

func NewGeminiClient(cfg config.GoogleAIConfig) (*GeminiClient, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{client: client, config: cfg}, nil
}

func (g *GeminiClient) GenerateContent(ctx context.Context, systemInstruction, userPrompt string) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, g.config.Timeout)
	defer cancel()

	contents := genai.Text(userPrompt)

	cfg := &genai.GenerateContentConfig{
		Temperature:       genai.Ptr(float32(0.7)),
		TopK:              genai.Ptr(float32(40)),
		TopP:              genai.Ptr(float32(0.95)),
		MaxOutputTokens:   int32(8192),
		SystemInstruction: genai.NewContentFromText(systemInstruction, "system"),
	}

	if g.config.UseGrounding {
		cfg.Tools = []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		}
	}

	resp, err := g.client.Models.GenerateContent(ctxWithTimeout, g.config.Model, contents, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned from Gemini")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	if candidate.Content.Parts[0].Text == "" {
		return "", fmt.Errorf("unexpected response type from Gemini")
	}

	return candidate.Content.Parts[0].Text, nil
}

func (g *GeminiClient) GenerateStructuredContent(ctx context.Context, systemInstruction, userPrompt string) (string, error) {
	enhancedInstruction := systemInstruction + "\n\n" +
		"=== CRITICAL OUTPUT RULES ===\n" +
		"1. Output MUST be ONLY valid JSON - nothing else\n" +
		"2. Do NOT wrap JSON in markdown code blocks\n" +
		"3. Do NOT add any explanatory text before or after the JSON\n" +
		"4. Start directly with { and end with }\n" +
		"5. Ensure all strings are properly escaped\n"

	response, err := g.GenerateContent(ctx, enhancedInstruction, userPrompt)
	if err != nil {
		return "", err
	}

	slog.Debug("Gemini structured response", "preview", response[:min(200, len(response))])

	return extractJSON(response), nil
}
