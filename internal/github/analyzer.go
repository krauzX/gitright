package github

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/krauzx/gitright/internal/models"
)

type Analyzer struct {
	client *Client
}

func NewAnalyzer(client *Client) *Analyzer {
	return &Analyzer{client: client}
}

func (a *Analyzer) AnalyzeRepository(ctx context.Context, token, owner, repo string) (*models.RepositoryAnalysis, error) {
	repository, err := a.client.GetRepository(ctx, token, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	languages, err := a.client.GetRepositoryLanguages(ctx, token, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	files, err := a.listAllFiles(ctx, token, owner, repo, "")
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	keyFiles, err := a.fetchKeyFiles(ctx, token, owner, repo, files)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch key files: %w", err)
	}

	dependencies := a.extractDependencies(keyFiles)

	commitCount, err := a.client.GetCommitCount(ctx, token, owner, repo)
	if err != nil {
		commitCount = 0
	}

	contributorCount, err := a.client.GetContributorCount(ctx, token, owner, repo)
	if err != nil {
		contributorCount = 0
	}

	return &models.RepositoryAnalysis{
		Repository:       a.convertRepository(repository),
		Languages:        languages,
		Files:            files,
		Dependencies:     dependencies,
		KeyFiles:         keyFiles,
		CommitCount:      commitCount,
		ContributorCount: contributorCount,
	}, nil
}

func (a *Analyzer) listAllFiles(ctx context.Context, token, owner, repo, path string) ([]string, error) {
	contents, err := a.client.ListRepositoryContents(ctx, token, owner, repo, path)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, content := range contents {
		if content.Type != nil && *content.Type == "file" {
			files = append(files, *content.Path)
		}
	}

	return files, nil
}

// fetchKeyFiles grabs the content of well-known dependency and config files.
func (a *Analyzer) fetchKeyFiles(ctx context.Context, token, owner, repo string, files []string) (map[string]string, error) {
	keyFilePatterns := []string{
		"package.json", "package-lock.json", "requirements.txt", "Pipfile",
		"pyproject.toml", "go.mod", "go.sum", "Cargo.toml", "Gemfile",
		"pom.xml", "build.gradle", "composer.json", "Dockerfile",
		".dockerignore", "docker-compose.yml", "README.md",
		"tsconfig.json", "vite.config.ts", "webpack.config.js",
	}

	keyFiles := make(map[string]string)

	for _, file := range files {
		filename := filepath.Base(file)
		for _, pattern := range keyFilePatterns {
			if filename == pattern {
				content, err := a.client.GetRepositoryContent(ctx, token, owner, repo, file)
				if err != nil {
					continue
				}
				keyFiles[file] = content
				break
			}
		}
	}

	return keyFiles, nil
}

func (a *Analyzer) extractDependencies(keyFiles map[string]string) map[string][]string {
	dependencies := make(map[string][]string)

	for path, content := range keyFiles {
		filename := filepath.Base(path)

		switch filename {
		case "package.json":
			if deps := a.extractNpmDependencies(content); len(deps) > 0 {
				dependencies["npm"] = deps
			}
		case "requirements.txt":
			if deps := a.extractPipDependencies(content); len(deps) > 0 {
				dependencies["pip"] = deps
			}
		case "go.mod":
			if deps := a.extractGoModDependencies(content); len(deps) > 0 {
				dependencies["go"] = deps
			}
		case "Cargo.toml":
			if deps := a.extractCargoDependencies(content); len(deps) > 0 {
				dependencies["cargo"] = deps
			}
		case "Gemfile":
			if deps := a.extractGemDependencies(content); len(deps) > 0 {
				dependencies["gem"] = deps
			}
		}
	}

	return dependencies
}

func (a *Analyzer) extractNpmDependencies(content string) []string {
	var packageJSON struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal([]byte(content), &packageJSON); err != nil {
		return nil
	}

	var deps []string
	for dep := range packageJSON.Dependencies {
		deps = append(deps, dep)
	}
	for dep := range packageJSON.DevDependencies {
		deps = append(deps, dep)
	}

	return deps
}

func (a *Analyzer) extractPipDependencies(content string) []string {
	lines := strings.Split(content, "\n")
	var deps []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == '=' || r == '>' || r == '<' || r == '~' || r == '!'
		})
		if len(parts) > 0 {
			deps = append(deps, strings.TrimSpace(parts[0]))
		}
	}

	return deps
}

func (a *Analyzer) extractGoModDependencies(content string) []string {
	lines := strings.Split(content, "\n")
	var deps []string
	inRequireBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}

		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		if inRequireBlock || strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				dep := strings.TrimPrefix(parts[0], "require ")
				deps = append(deps, dep)
			}
		}
	}

	return deps
}

func (a *Analyzer) extractCargoDependencies(content string) []string {
	lines := strings.Split(content, "\n")
	var deps []string
	inDepsSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[dependencies]") {
			inDepsSection = true
			continue
		}

		if strings.HasPrefix(line, "[") {
			inDepsSection = false
		}

		if inDepsSection && strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) > 0 {
				deps = append(deps, strings.TrimSpace(parts[0]))
			}
		}
	}

	return deps
}

func (a *Analyzer) extractGemDependencies(content string) []string {
	lines := strings.Split(content, "\n")
	var deps []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "gem ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, strings.Trim(parts[1], "'\""))
			}
		}
	}

	return deps
}

func (a *Analyzer) convertRepository(repo *github.Repository) *models.Repository {
	return &models.Repository{
		ID:              repo.GetID(),
		GitHubID:        repo.GetID(),
		Name:            repo.GetName(),
		FullName:        repo.GetFullName(),
		Description:     repo.GetDescription(),
		Private:         repo.GetPrivate(),
		Fork:            repo.GetFork(),
		Language:        repo.GetLanguage(),
		StargazersCount: repo.GetStargazersCount(),
		ForksCount:      repo.GetForksCount(),
		OpenIssuesCount: repo.GetOpenIssuesCount(),
		DefaultBranch:   repo.GetDefaultBranch(),
		Topics:          repo.Topics,
		HTMLURL:         repo.GetHTMLURL(),
		CloneURL:        repo.GetCloneURL(),
		CreatedAt:       repo.GetCreatedAt().Time,
		UpdatedAt:       repo.GetUpdatedAt().Time,
		PushedAt:        repo.GetPushedAt().Time,
	}
}
