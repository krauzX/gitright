package graph

import (
	"github.com/krauzx/gitright/internal/github"
)

func BuildFromProfile(profile *github.DeepProfileData) *Graph {
	g := &Graph{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}

	g.Nodes = append(g.Nodes, Node{
		ID:       profile.Login,
		Label:    profile.Name,
		Category: "profile",
		Size:     25,
		Color:    "#f59e0b",
	})

	langCounts := make(map[string]int)
	for _, repo := range profile.TopRepositories {
		langCounts[repo.Language] += repo.Contributions
	}

	langNodes := make(map[string]bool)
	for lang, count := range langCounts {
		size := 8 + count/10
		if size > 18 {
			size = 18
		}
		g.Nodes = append(g.Nodes, Node{
			ID:       "lang:" + lang,
			Label:    lang,
			Category: categorizeLanguage(lang),
			Size:     size,
		})
		langNodes[lang] = true
		g.Edges = append(g.Edges, Edge{
			Source: profile.Login,
			Target: "lang:" + lang,
			Weight: float64(count) / 100,
		})
	}

	for _, repo := range profile.TopRepositories {
		size := 5 + repo.Stars/100
		if size > 15 {
			size = 15
		}
		g.Nodes = append(g.Nodes, Node{
			ID:       "repo:" + repo.Name,
			Label:    repo.Name,
			Category: "repository",
			Size:     size,
		})

		if langNodes[repo.Language] {
			g.Edges = append(g.Edges, Edge{
				Source: "lang:" + repo.Language,
				Target: "repo:" + repo.Name,
				Weight: float64(repo.Contributions) / 50,
			})
		}
	}

	return g
}

func categorizeLanguage(lang string) string {
	systems := map[string]bool{"go": true, "rust": true, "c": true, "cpp": true}
	data := map[string]bool{"python": true, "r": true, "julia": true}
	frontend := map[string]bool{"javascript": true, "typescript": true, "html": true, "css": true}

	if systems[lang] {
		return "systems"
	}
	if data[lang] {
		return "data"
	}
	if frontend[lang] {
		return "frontend"
	}
	return "other"
}
