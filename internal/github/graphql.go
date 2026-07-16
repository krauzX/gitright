package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GraphQLClient struct {
	client *githubv4.Client
}

func NewGraphQLClient(token string) *GraphQLClient {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), src)
	return &GraphQLClient{
		client: githubv4.NewClient(tc),
	}
}

func (g *GraphQLClient) GetDeepUserProfile(ctx context.Context, login string) (*DeepUserProfile, error) {
	var query struct {
		User DeepUserProfile `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	if err := g.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("graphql user query failed: %w", err)
	}

	slog.Debug("Fetched deep profile via GraphQL", "login", login, "contributions", query.User.Contributions.ContributionCalendar.TotalContributions)
	return &query.User, nil
}

func (g *GraphQLClient) GetUserSocialAccounts(ctx context.Context, login string) ([]SocialAccount, error) {
	var query struct {
		User struct {
			SocialAccounts SocialAccountConnection `graphql:"socialAccounts(first: 20)"`
		} `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	if err := g.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("graphql social accounts query failed: %w", err)
	}

	accounts := make([]SocialAccount, 0, len(query.User.SocialAccounts.Nodes))
	for _, node := range query.User.SocialAccounts.Nodes {
		accounts = append(accounts, SocialAccount{
			Provider: node.Provider,
			URL:      node.URL,
			Username: node.DisplayName,
		})
	}

	return accounts, nil
}

func (g *GraphQLClient) GetUserContributionStats(ctx context.Context, login string) (*ContributionStats, error) {
	type contributionStatsQuery struct {
		User struct {
			Contributions ContributionsCollection `graphql:"contributionsCollection"`
		} `graphql:"user(login: $login)"`
	}

	var query contributionStatsQuery

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	if err := g.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("graphql contribution stats query failed: %w", err)
	}

	contribs := query.User.Contributions
	stats := &ContributionStats{
		TotalContributions: contribs.ContributionCalendar.TotalContributions,
		IssueContributions: contribs.IssueContributions.TotalCount,
		PRContributions:    contribs.PullRequestReviewContributions.TotalCount,
		TopRepositories:    make([]TopRepository, 0, len(contribs.CommitContributionsByRepository)),
	}

	for _, repo := range contribs.CommitContributionsByRepository {
		stats.TopRepositories = append(stats.TopRepositories, TopRepository{
			Name:          repo.Repository.NameWithOwner,
			Stars:         repo.Repository.StargazerCount,
			Language:      repo.Repository.PrimaryLanguage.Name,
			Contributions: repo.Contributions.ContributionCount,
		})
	}

	return stats, nil
}
