package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/krauzx/gitright/internal/github"
)

type AutoImportService struct{}

func NewAutoImportService() *AutoImportService {
	return &AutoImportService{}
}

type TopRepoEntry struct {
	Name          string `json:"name"`
	Stars         int    `json:"stars"`
	Language      string `json:"language"`
	Contributions int    `json:"contributions"`
}

type AutoImportData struct {
	Name               string            `json:"name"`
	Login              string            `json:"login"`
	Bio                string            `json:"bio"`
	AvatarURL          string            `json:"avatar_url"`
	Company            string            `json:"company"`
	Location           string            `json:"location"`
	Email              string            `json:"email"`
	WebsiteURL         string            `json:"website_url"`
	TwitterUsername    string            `json:"twitter_username"`
	Followers          int               `json:"followers"`
	TotalContributions int               `json:"total_contributions"`
	IssueContributions int               `json:"issue_contributions"`
	ReviewContributions int              `json:"review_contributions"`
	Skills             []string          `json:"skills"`
	LanguageBreakdown  map[string]int    `json:"language_breakdown"`
	TopRepositories    []TopRepoEntry    `json:"top_repositories"`
	SocialLinks        []SocialLinkEntry `json:"social_links"`
}

type SocialLinkEntry struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
	Username string `json:"username"`
}

func (s *AutoImportService) FetchDeepProfile(ctx context.Context, token string, login string) (*github.DeepProfileData, error) {
	client := github.NewGraphQLClient(token)

	user, err := client.GetDeepUserProfile(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}

	stats, err := client.GetUserContributionStats(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contribution stats: %w", err)
	}

	socialAccounts, err := client.GetUserSocialAccounts(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch social accounts: %w", err)
	}

	return user.ToDeepProfileData(stats, socialAccounts), nil
}

func (s *AutoImportService) FetchAndExtract(ctx context.Context, token string, login string) (*AutoImportData, error) {
	client := github.NewGraphQLClient(token)

	user, err := client.GetDeepUserProfile(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}

	stats, err := client.GetUserContributionStats(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contribution stats: %w", err)
	}

	socialAccounts, err := client.GetUserSocialAccounts(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch social accounts: %w", err)
	}

	skills := extractSkills(stats)
	langBreakdown := buildLanguageBreakdown(stats)
	topRepos := buildTopRepos(stats)
	socialLinks := buildSocialLinks(socialAccounts, user)

	data := &AutoImportData{
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
		LanguageBreakdown:  langBreakdown,
		TopRepositories:    topRepos,
		SocialLinks:        socialLinks,
	}

	return data, nil
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

	sort.Strings(skills)
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

func buildTopRepos(stats *github.ContributionStats) []TopRepoEntry {
	repos := make([]TopRepoEntry, 0, len(stats.TopRepositories))
	for _, r := range stats.TopRepositories {
		repos = append(repos, TopRepoEntry{
			Name:          r.Name,
			Stars:         r.Stars,
			Language:      r.Language,
			Contributions: r.Contributions,
		})
	}
	return repos
}

func buildSocialLinks(accounts []github.SocialAccount, user *github.DeepUserProfile) []SocialLinkEntry {
	var links []SocialLinkEntry

	for _, acc := range accounts {
		links = append(links, SocialLinkEntry{
			Provider: acc.Provider,
			URL:      acc.URL,
			Username: acc.Username,
		})
	}

	if user.WebsiteURL != "" {
		hasWebsite := false
		for _, l := range links {
			if l.URL == user.WebsiteURL {
				hasWebsite = true
				break
			}
		}
		if !hasWebsite {
			links = append(links, SocialLinkEntry{
				Provider: "WEBSITE",
				URL:      user.WebsiteURL,
				Username: user.Login,
			})
		}
	}

	return links
}
