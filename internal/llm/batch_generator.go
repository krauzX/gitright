package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/krauzx/gitright/internal/models"
)

type BatchProfileRequest struct {
	Username         string
	Bio              string
	Location         string
	Company          string
	TargetRole       string
	ToneOfVoice      string
	EmphasizedSkills []string
	Projects         []models.RepositoryAnalysis
}

type BatchProfileResponse struct {
	ProfilePitch     string               `json:"profile_pitch"`
	ProjectSummaries []ProjectSummaryData `json:"project_summaries"`
	ExtractedSkills  []string             `json:"extracted_skills"`
	Confidence       float64              `json:"confidence"`
}

type ProjectSummaryData struct {
	ProjectName string   `json:"project_name"`
	Summary     string   `json:"summary"`
	Skills      []string `json:"skills"`
}

func (cg *ContentGenerator) GenerateBatchedProfile(ctx context.Context, apiKey string, req BatchProfileRequest) (*BatchProfileResponse, error) {
	tempClient, err := cg.createClientWithAPIKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	systemInstruction := buildBatchedSystemInstruction(req)
	userPrompt := buildBatchedUserPrompt(req)

	responseText, err := tempClient.GenerateStructuredContent(ctx, systemInstruction, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	// Extract JSON from response (handles markdown blocks and extra text)
	jsonStr := extractJSON(responseText)
	if jsonStr == "" {
		slog.Error("No JSON found in response", "response", responseText)
		return nil, fmt.Errorf("invalid response: no JSON content found")
	}

	var response BatchProfileResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		slog.Error("Failed to parse JSON", "error", err, "json", jsonStr[:min(500, len(jsonStr))])
		return nil, fmt.Errorf("invalid response format: %w", err)
	}

	// Validate response
	if response.ProfilePitch == "" {
		return nil, fmt.Errorf("missing profile_pitch in response")
	}
	if len(response.ProjectSummaries) == 0 {
		return nil, fmt.Errorf("missing project_summaries in response")
	}

	return &response, nil
}

func (cg *ContentGenerator) createClientWithAPIKey(apiKey string) (*GeminiClient, error) {
	config := cg.client.config
	config.APIKey = apiKey
	return NewGeminiClient(config)
}

//  creates comprehensive system instruction for batch generation
func buildBatchedSystemInstruction(req BatchProfileRequest) string {
	var sb strings.Builder

	sb.WriteString("You are an expert GitHub Profile Writer creating professional, compelling README profiles.\n\n")

	sb.WriteString(fmt.Sprintf("TARGET ROLE: %s\n", req.TargetRole))
	sb.WriteString(fmt.Sprintf("TONE: %s\n", req.ToneOfVoice))

	if len(req.EmphasizedSkills) > 0 {
		sb.WriteString(fmt.Sprintf("EMPHASIZE SKILLS: %s\n", strings.Join(req.EmphasizedSkills, ", ")))
	}

	sb.WriteString("\n=== RESPONSE FORMAT (STRICT JSON) ===\n")
	sb.WriteString("Return ONLY valid JSON with this EXACT structure (no markdown, no explanations):\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"profile_pitch\": \"2-3 compelling paragraphs introducing the developer\",\n")
	sb.WriteString("  \"project_summaries\": [\n")
	sb.WriteString("    {\n")
	sb.WriteString("      \"project_name\": \"exact project name from input\",\n")
	sb.WriteString("      \"summary\": \"2-3 sentences highlighting impact and technical sophistication\",\n")
	sb.WriteString("      \"skills\": [\"skill1\", \"skill2\", ...]\n")
	sb.WriteString("    }\n")
	sb.WriteString("  ],\n")
	sb.WriteString("  \"extracted_skills\": [\"all unique skills from all projects\"],\n")
	sb.WriteString("  \"confidence\": 0.9\n")
	sb.WriteString("}\n\n")

	sb.WriteString("CRITICAL RULES:\n")
	sb.WriteString("1. Output ONLY the JSON object - no markdown code blocks\n")
	sb.WriteString("2. Match project_name EXACTLY to input names\n")
	sb.WriteString("3. Write in FIRST PERSON (I/my) to show ownership\n")
	sb.WriteString("4. Emphasize quantifiable achievements and technical complexity\n")
	sb.WriteString("5. Keep summaries concise but impactful (2-3 sentences each)\n")
	sb.WriteString("6. Extract 10-15 most relevant skills across all projects\n")

	return sb.String()
}

// buildBatchedUserPrompt creates comprehensive user prompt with ALL project data
func buildBatchedUserPrompt(req BatchProfileRequest) string {
	var sb strings.Builder

	sb.WriteString("=== DEVELOPER PROFILE ===\n")
	sb.WriteString(fmt.Sprintf("Username: %s\n", req.Username))
	if req.Bio != "" {
		sb.WriteString(fmt.Sprintf("Bio: %s\n", req.Bio))
	}
	if req.Location != "" {
		sb.WriteString(fmt.Sprintf("Location: %s\n", req.Location))
	}
	if req.Company != "" {
		sb.WriteString(fmt.Sprintf("Company: %s\n", req.Company))
	}
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("=== %d PROJECTS TO ANALYZE ===\n\n", len(req.Projects)))

	for i, project := range req.Projects {
		sb.WriteString(fmt.Sprintf("--- PROJECT %d: %s ---\n", i+1, project.Repository.Name))

		if project.Repository.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", project.Repository.Description))
		}

		sb.WriteString(fmt.Sprintf("Stars: %d | Forks: %d\n",
			project.Repository.StargazersCount,
			project.Repository.ForksCount))

		if len(project.Languages) > 0 {
			sb.WriteString("Languages: ")
			langs := make([]string, 0, len(project.Languages))
			for lang := range project.Languages {
				langs = append(langs, lang)
			}
			sb.WriteString(strings.Join(langs, ", "))
			sb.WriteString("\n")
		}

		if len(project.Dependencies) > 0 {
			sb.WriteString("Dependencies:\n")
			for ecosystem, deps := range project.Dependencies {
				if len(deps) > 0 {
					// Limit to first 5 deps per ecosystem to avoid token bloat
					depsSlice := deps
					if len(deps) > 5 {
						depsSlice = deps[:5]
					}
					sb.WriteString(fmt.Sprintf("  - %s: %s\n", ecosystem, strings.Join(depsSlice, ", ")))
				}
			}
		}

		if len(project.Repository.Topics) > 0 {
			sb.WriteString(fmt.Sprintf("Topics: %s\n", strings.Join(project.Repository.Topics, ", ")))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("Generate the complete profile response in the specified JSON format.\n")

	return sb.String()
}

// extractJSON extracts JSON from Gemini response (handles markdown blocks)
func extractJSON(response string) string {
	// Try to find JSON in markdown code block
	jsonBlockPattern := regexp.MustCompile("(?s)```(?:json)?\\s*({.*?})\\s*```")
	if matches := jsonBlockPattern.FindStringSubmatch(response); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try to find raw JSON object
	jsonPattern := regexp.MustCompile("(?s)({\\s*\"[^\"]+\".*})")
	if matches := jsonPattern.FindStringSubmatch(response); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Return as-is if already clean
	if strings.HasPrefix(strings.TrimSpace(response), "{") {
		return strings.TrimSpace(response)
	}

	return ""
}
