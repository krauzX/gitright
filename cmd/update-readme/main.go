package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krauzx/gitright/internal/github"
	"github.com/krauzx/gitright/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN required")
	}

	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		username = "krauzX"
	}

	client := github.NewGraphQLClient(token)

	user, err := client.GetDeepUserProfile(context.Background(), username)
	if err != nil {
		log.Fatalf("Failed to fetch user: %v", err)
	}

	stats, err := client.GetUserContributionStats(context.Background(), username)
	if err != nil {
		log.Fatalf("Failed to fetch stats: %v", err)
	}

	socialAccounts, err := client.GetUserSocialAccounts(context.Background(), username)
	if err != nil {
		log.Printf("Warning: failed to fetch social accounts: %v", err)
	}

	skills := extractSkills(stats)

	socialLinks := make([]services.SocialLinkEntry, 0)
	for _, acc := range socialAccounts {
		socialLinks = append(socialLinks, services.SocialLinkEntry{
			Provider: acc.Provider,
			URL:      acc.URL,
			Username: acc.Username,
		})
	}

	topRepos := make([]services.TopRepoEntry, 0)
	for _, r := range stats.TopRepositories {
		topRepos = append(topRepos, services.TopRepoEntry{
			Name:          r.Name,
			Stars:         r.Stars,
			Language:      r.Language,
			Contributions: r.Contributions,
		})
	}

	data := &services.AutoImportData{
		Name:               user.Name,
		Login:              user.Login,
		Bio:                user.Bio,
		AvatarURL:          user.AvatarURL,
		Company:            user.Company,
		Location:           user.Location,
		Email:              user.Email,
		WebsiteURL:         user.WebsiteURL,
		TwitterUsername:    user.TwitterUsername,
		Followers:          user.Followers.TotalCount,
		TotalContributions: stats.TotalContributions,
		IssueContributions: stats.IssueContributions,
		ReviewContributions: stats.PRContributions,
		Skills:             skills,
		LanguageBreakdown:  buildLanguageBreakdown(stats),
		TopRepositories:    topRepos,
		SocialLinks:        socialLinks,
	}

	_ = data

	fmt.Println("README update completed")
	fmt.Printf("User: %s (%s)\n", data.Name, data.Login)
	fmt.Printf("Skills: %v\n", data.Skills)
	fmt.Printf("Contributions: %d\n", data.TotalContributions)
}

func extractSkills(stats *github.ContributionStats) []string {
	seen := make(map[string]bool)
	var skills []string
	for _, repo := range stats.TopRepositories {
		if repo.Language != "" && !seen[repo.Language] {
			seen[repo.Language] = true
			skills = append(skills, repo.Language)
		}
	}
	return skills
}

func buildLanguageBreakdown(stats *github.ContributionStats) map[string]int {
	langs := make(map[string]int)
	for _, repo := range stats.TopRepositories {
		if repo.Language != "" {
			langs[repo.Language] += repo.Contributions
		}
	}
	return langs
}
