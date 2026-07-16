package github

type DeepUserProfile struct {
	Login          string                     `graphql:"login"`
	Name           string                     `graphql:"name"`
	Bio            string                     `graphql:"bio"`
	AvatarURL      string                     `graphql:"avatarUrl(size: 400)"`
	Company        string                     `graphql:"company"`
	Location       string                     `graphql:"location"`
	Email          string                     `graphql:"email"`
	WebsiteURL     string                     `graphql:"websiteUrl"`
	TwitterUsername string                     `graphql:"twitterUsername"`
	Followers      struct{ TotalCount int `graphql:"totalCount"` } `graphql:"followers"`
	Following      struct{ TotalCount int `graphql:"totalCount"` } `graphql:"following"`
	Repositories   struct{ TotalCount int `graphql:"totalCount"` } `graphql:"repositories"`
	Contributions  ContributionsCollection    `graphql:"contributionsCollection"`
	SocialAccounts SocialAccountConnection    `graphql:"socialAccounts(first: 20)"`
}

type ContributionsCollection struct {
	ContributionCalendar                 ContributionCalendar     `graphql:"contributionCalendar"`
	CommitContributionsByRepository     []CommitContributionByRepo `graphql:"commitContributionsByRepository"`
	PullRequestContributionsByRepository []PRContributionByRepo    `graphql:"pullRequestContributionsByRepository"`
	IssueContributions                  struct{ TotalCount int `graphql:"totalCount"` } `graphql:"issueContributions"`
	PullRequestReviewContributions      struct{ TotalCount int `graphql:"totalCount"` } `graphql:"pullRequestReviewContributions"`
}

type ContributionCalendar struct {
	TotalContributions int                `graphql:"totalContributions"`
	Weeks              []ContributionWeek `graphql:"weeks"`
}

type ContributionWeek struct {
	ContributionDays []ContributionDay `graphql:"contributionDays"`
}

type ContributionDay struct {
	Date              string `graphql:"date"`
	ContributionCount int    `graphql:"contributionCount"`
	Color             string `graphql:"color"`
}

type CommitContributionByRepo struct {
	Repository struct {
		NameWithOwner   string `graphql:"nameWithOwner"`
		StargazerCount  int    `graphql:"stargazerCount"`
		PrimaryLanguage struct {
			Name string `graphql:"name"`
		} `graphql:"primaryLanguage"`
	} `graphql:"repository"`
	Contributions struct {
		ContributionCount int `graphql:"totalCount"`
	} `graphql:"contributions"`
}

type PRContributionByRepo struct {
	Repository struct {
		NameWithOwner string `graphql:"nameWithOwner"`
	} `graphql:"repository"`
	Contributions struct {
		ContributionCount int `graphql:"totalCount"`
	} `graphql:"contributions"`
}

type SocialAccountConnection struct {
	Nodes []SocialAccountNode `graphql:"nodes"`
}

type SocialAccountNode struct {
	Provider string `graphql:"provider"`
	URL      string `graphql:"url"`
	DisplayName string `graphql:"displayName"`
}

type SocialAccount struct {
	Provider string
	URL      string
	Username string
}

type ContributionStats struct {
	TotalContributions int
	IssueContributions int
	PRContributions    int
	TopRepositories    []TopRepository
}

type TopRepository struct {
	Name          string
	Stars         int
	Language      string
	Contributions int
}

type DeepProfileData struct {
	Login              string             `json:"login"`
	Name               string             `json:"name"`
	Bio                string             `json:"bio"`
	AvatarURL          string             `json:"avatar_url"`
	Company            string             `json:"company"`
	Location           string             `json:"location"`
	Email              string             `json:"email"`
	WebsiteURL         string             `json:"website_url"`
	TwitterUsername    string             `json:"twitter_username"`
	Followers          int                `json:"followers"`
	Following          int                `json:"following"`
	TotalRepos         int                `json:"total_repos"`
	TotalContributions int                `json:"total_contributions"`
	ContributionWeeks  []ContributionWeek `json:"contribution_weeks"`
	TopRepositories    []TopRepository    `json:"top_repositories"`
	SocialAccounts     []SocialAccount    `json:"social_accounts"`
	PRContributions    int                `json:"pr_contributions"`
	IssueContributions int                `json:"issue_contributions"`
	ReviewContributions int               `json:"review_contributions"`
}

func (p *DeepUserProfile) ToDeepProfileData(stats *ContributionStats, socialAccounts []SocialAccount) *DeepProfileData {
	data := &DeepProfileData{
		Login:              p.Login,
		Name:               p.Name,
		Bio:                p.Bio,
		AvatarURL:          p.AvatarURL,
		Company:            p.Company,
		Location:           p.Location,
		Email:              p.Email,
		WebsiteURL:         p.WebsiteURL,
		TwitterUsername:    p.TwitterUsername,
		Followers:          p.Followers.TotalCount,
		Following:          p.Following.TotalCount,
		TotalRepos:         p.Repositories.TotalCount,
		TotalContributions: p.Contributions.ContributionCalendar.TotalContributions,
		ContributionWeeks:  p.Contributions.ContributionCalendar.Weeks,
		SocialAccounts:     socialAccounts,
	}

	if stats != nil {
		data.TotalContributions = stats.TotalContributions
		data.IssueContributions = stats.IssueContributions
		data.PRContributions = stats.PRContributions
		data.ReviewContributions = stats.PRContributions
		data.TopRepositories = stats.TopRepositories
	}

	return data
}
