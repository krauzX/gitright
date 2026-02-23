package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v60/github"
	"github.com/krauzx/gitright/internal/config"
	"golang.org/x/oauth2"
	goauth "golang.org/x/oauth2/github"
)

type Client struct {
	config      *config.GitHubConfig
	oauthConfig *oauth2.Config
}

func NewClient(cfg config.GitHubConfig) *Client {
	return &Client{
		config: &cfg,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURI,
			Scopes:       cfg.Scopes,
			Endpoint:     goauth.Endpoint,
		},
	}
}

func (c *Client) GetAuthorizationURL(state string) string {
	return c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

func (c *Client) NewAuthenticatedClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (c *Client) GetUser(ctx context.Context, token string) (*github.User, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return user, nil
}

func (c *Client) GetUserEmails(ctx context.Context, token string) ([]*github.UserEmail, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	emails, resp, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return emails, nil
}

func (c *Client) ListRepositories(ctx context.Context, token string, includePrivate bool) ([]*github.Repository, error) {
	client := c.NewAuthenticatedClient(ctx, token)

	opts := &github.RepositoryListOptions{
		Visibility: "public",
		Sort:       "updated",
		Direction:  "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	if includePrivate {
		opts.Visibility = "all"
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}
		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func (c *Client) GetRepository(ctx context.Context, token, owner, repo string) (*github.Repository, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	repository, resp, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return repository, nil
}

func (c *Client) GetRepositoryLanguages(ctx context.Context, token, owner, repo string) (map[string]int, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	languages, resp, err := client.Repositories.ListLanguages(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository languages: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return languages, nil
}

func (c *Client) GetRepositoryContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	fileContent, _, resp, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get file content: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if fileContent == nil {
		return "", fmt.Errorf("file content is nil")
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return content, nil
}

func (c *Client) ListRepositoryContents(ctx context.Context, token, owner, repo, path string) ([]*github.RepositoryContent, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	_, directoryContent, resp, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory contents: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return directoryContent, nil
}

func (c *Client) GetCommitCount(ctx context.Context, token, owner, repo string) (int, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	commits, resp, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get commit count: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.LastPage > 0 {
		return resp.LastPage, nil
	}

	return len(commits), nil
}

func (c *Client) GetContributorCount(ctx context.Context, token, owner, repo string) (int, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	contributors, resp, err := client.Repositories.ListContributors(ctx, owner, repo, &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get contributor count: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return len(contributors), nil
}

func (c *Client) CreateOrUpdateFile(ctx context.Context, token, owner, repo, path, message, content, sha string) error {
	client := c.NewAuthenticatedClient(ctx, token)

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(content),
		Branch:  github.String("main"),
	}

	if sha != "" {
		opts.SHA = github.String(sha)
	}

	_, resp, err := client.Repositories.CreateFile(ctx, owner, repo, path, opts)
	if err != nil {
		return fmt.Errorf("failed to create/update file: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GetProfileReadmeSHA(ctx context.Context, token, username string) (string, error) {
	client := c.NewAuthenticatedClient(ctx, token)
	fileContent, _, resp, err := client.Repositories.GetContents(ctx, username, username, "README.md", nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to get README.md: %w", err)
	}

	if fileContent != nil && fileContent.SHA != nil {
		return *fileContent.SHA, nil
	}

	return "", nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (bool, error) {
	_, err := c.GetUser(ctx, token)
	return err == nil, err
}
