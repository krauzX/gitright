package services

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/krauzx/gitright/internal/llm"
	"github.com/krauzx/gitright/internal/models"
	"github.com/krauzx/gitright/internal/repository"
)

type ProfileService struct {
	contentGenerator *llm.ContentGenerator
	projectRepo      *repository.ProjectRepository
	githubService    *GitHubService
	profileCacheRepo *repository.ProfileCacheRepository
}

func NewProfileService(
	contentGenerator *llm.ContentGenerator,
	projectRepo *repository.ProjectRepository,
	githubService *GitHubService,
	profileCacheRepo *repository.ProfileCacheRepository,
) *ProfileService {
	return &ProfileService{
		contentGenerator: contentGenerator,
		projectRepo:      projectRepo,
		githubService:    githubService,
		profileCacheRepo: profileCacheRepo,
	}
}

func (s *ProfileService) GenerateProfile(ctx context.Context, req *models.ContentGenerationRequest, user *models.User) (*models.ContentGenerationResponse, error) {
	if req.UserAPIKey == "" {
		return nil, fmt.Errorf("API key required - get free key: https://aistudio.google.com/app/apikey")
	}
	if len(req.Projects) == 0 {
		return nil, fmt.Errorf("at least one project required")
	}

	cacheKey := repository.GetCacheKey(user.Username, req.TargetRole, req.ToneOfVoice, len(req.Projects))
	if cached, err := s.profileCacheRepo.Get(ctx, cacheKey); err == nil && cached != nil {
		return cached, nil
	}

	batchReq := llm.BatchProfileRequest{
		Username:         user.Username,
		Bio:              user.Bio,
		Location:         user.Location,
		Company:          user.Company,
		TargetRole:       req.TargetRole,
		ToneOfVoice:      req.ToneOfVoice,
		EmphasizedSkills: req.EmphasizedSkills,
		Projects:         req.Projects,
	}

	batchResp, err := s.contentGenerator.GenerateBatchedProfile(ctx, req.UserAPIKey, batchReq)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	summaries := make([]models.ProjectSummary, 0, len(batchResp.ProjectSummaries))
	for i, proj := range batchResp.ProjectSummaries {
		if i >= len(req.Projects) {
			break
		}
		summaries = append(summaries, models.ProjectSummary{
			Repository: req.Projects[i].Repository,
			Summary:    proj.Summary,
			TechStack:  proj.Skills,
		})
	}

	config := &models.ProfileConfig{
		TargetRole:     req.TargetRole,
		SkillsEmphasis: req.EmphasizedSkills,
		ToneOfVoice:    req.ToneOfVoice,
		ContactPrefs:   req.ContactPrefs,
	}

	badges := s.buildBadgesFromProjectData(req.Projects, batchResp.ExtractedSkills, req.EmphasizedSkills)
	markdown := s.buildMarkdown(user, req, batchResp.ProfilePitch, summaries, badges, config)

	response := &models.ContentGenerationResponse{
		Markdown:        markdown,
		ExtractedSkills: batchResp.ExtractedSkills,
		SuggestedBadges: badges,
		Confidence:      batchResp.Confidence,
	}

	if err := s.profileCacheRepo.Set(ctx, user.ID, 0, cacheKey, response, 24*time.Hour); err != nil {
		slog.Warn("Failed to cache profile generation result", "username", user.Username, "error", err)
	}
	return response, nil
}

// DeployProfile deploys the profile markdown to GitHub.
func (s *ProfileService) DeployProfile(ctx context.Context, accessToken, username, markdown string) error {
	if err := s.githubService.DeployProfileREADME(ctx, accessToken, username, markdown); err != nil {
		return fmt.Errorf("failed to deploy profile: %w", err)
	}
	return nil
}

// Badge generation

// buildBadgesFromProjectData creates badges sourced from (in priority order):
//  1. EmphasizedSkills from the request
//  2. Actual programming languages found in every RepositoryAnalysis
//  3. Frameworks/libraries inferred from dependency names
//  4. LLM-extracted skills
func (s *ProfileService) buildBadgesFromProjectData(
	projects []models.RepositoryAnalysis,
	llmSkills, emphasizedSkills []string,
) []models.Badge {
	catalog := buildBadgeCatalog()
	badgeMap := make(map[string]models.Badge)

	add := func(key string) {
		if badge, ok := catalog[strings.ToLower(key)]; ok {
			if _, exists := badgeMap[badge.Name]; !exists {
				badgeMap[badge.Name] = badge
			}
		}
	}

	// Priority 1 – user-selected emphasis
	for _, skill := range emphasizedSkills {
		add(skill)
	}

	// Priority 2 – every language detected in every repo
	for _, p := range projects {
		for lang := range p.Languages {
			add(lang)
		}
	}

	// Priority 3 – infer frameworks from dependency names
	for _, p := range projects {
		for _, deps := range p.Dependencies {
			for _, dep := range deps {
				key := strings.ToLower(dep)
				// Normalise scoped npm packages: @angular/core -> angular
				if strings.HasPrefix(key, "@") {
					parts := strings.SplitN(key[1:], "/", 2)
					key = parts[0]
				}
				add(key)
			}
		}
	}

	// Priority 4 – LLM extracted skills fill remaining gaps
	for _, skill := range llmSkills {
		add(skill)
	}

	badges := make([]models.Badge, 0, len(badgeMap))
	for _, b := range badgeMap {
		badges = append(badges, b)
	}
	return badges
}

// buildBadgeCatalog returns a comprehensive technology → Badge map (all keys lowercase).
// Only Name and Color are needed; the URL is constructed dynamically when rendering.
func buildBadgeCatalog() map[string]models.Badge {
	entries := []models.Badge{
		// ---------- Languages ----------
		{Name: "Go", Color: "00ADD8"},
		{Name: "Python", Color: "3776AB"},
		{Name: "JavaScript", Color: "F7DF1E"},
		{Name: "TypeScript", Color: "3178C6"},
		{Name: "Rust", Color: "000000"},
		{Name: "Java", Color: "ED8B00"},
		{Name: "Kotlin", Color: "7F52FF"},
		{Name: "Swift", Color: "FA7343"},
		{Name: "C++", Color: "00599C"},
		{Name: "C", Color: "A8B9CC"},
		{Name: "C#", Color: "239120"},
		{Name: "PHP", Color: "777BB4"},
		{Name: "Ruby", Color: "CC342D"},
		{Name: "Dart", Color: "0175C2"},
		{Name: "Scala", Color: "DC322F"},
		{Name: "Elixir", Color: "4B275F"},
		{Name: "Haskell", Color: "5D4F85"},
		{Name: "Lua", Color: "2C2D72"},
		{Name: "Shell", Color: "4EAA25"},
		{Name: "HTML5", Color: "E34F26"},
		{Name: "CSS3", Color: "1572B6"},
		// ---------- Frontend ----------
		{Name: "React", Color: "61DAFB"},
		{Name: "Vue.js", Color: "4FC08D"},
		{Name: "Angular", Color: "DD0031"},
		{Name: "Svelte", Color: "FF3E00"},
		{Name: "Next.js", Color: "000000"},
		{Name: "Nuxt.js", Color: "00DC82"},
		{Name: "Gatsby", Color: "663399"},
		{Name: "Remix", Color: "000000"},
		{Name: "Astro", Color: "FF5D01"},
		{Name: "TailwindCSS", Color: "06B6D4"},
		{Name: "Vite", Color: "646CFF"},
		// ---------- Backend ----------
		{Name: "Node.js", Color: "339933"},
		{Name: "Express.js", Color: "000000"},
		{Name: "Fastify", Color: "000000"},
		{Name: "NestJS", Color: "E0234E"},
		{Name: "Django", Color: "092E20"},
		{Name: "Flask", Color: "000000"},
		{Name: "FastAPI", Color: "009688"},
		{Name: "Spring Boot", Color: "6DB33F"},
		{Name: "Laravel", Color: "FF2D20"},
		{Name: "Ruby on Rails", Color: "CC0000"},
		{Name: "Fiber", Color: "00ADD8"},
		{Name: "Gin", Color: "00ADD8"},
		{Name: "Echo", Color: "00ADD8"},
		// ---------- Mobile ----------
		{Name: "Flutter", Color: "02569B"},
		{Name: "React Native", Color: "61DAFB"},
		// ---------- Databases ----------
		{Name: "PostgreSQL", Color: "316192"},
		{Name: "MySQL", Color: "00000F"},
		{Name: "MongoDB", Color: "47A248"},
		{Name: "Redis", Color: "DC382D"},
		{Name: "SQLite", Color: "07405E"},
		{Name: "Cassandra", Color: "1287B1"},
		{Name: "Elasticsearch", Color: "005571"},
		{Name: "Supabase", Color: "3ECF8E"},
		{Name: "Firebase", Color: "FFCA28"},
		{Name: "Neon", Color: "00E699"},
		{Name: "Prisma", Color: "2D3748"},
		{Name: "Drizzle", Color: "C5F74F"},
		// ---------- DevOps & Cloud ----------
		{Name: "Docker", Color: "2496ED"},
		{Name: "Kubernetes", Color: "326CE5"},
		{Name: "Terraform", Color: "7B42BC"},
		{Name: "Ansible", Color: "EE0000"},
		{Name: "AWS", Color: "FF9900"},
		{Name: "GCP", Color: "4285F4"},
		{Name: "Azure", Color: "0078D4"},
		{Name: "Vercel", Color: "000000"},
		{Name: "Netlify", Color: "00C7B7"},
		{Name: "Heroku", Color: "430098"},
		{Name: "Nginx", Color: "009639"},
		// ---------- Tooling ----------
		{Name: "Git", Color: "F05032"},
		{Name: "GraphQL", Color: "E10098"},
		{Name: "gRPC", Color: "244C5A"},
		{Name: "Apache Kafka", Color: "231F20"},
		{Name: "RabbitMQ", Color: "FF6600"},
		{Name: "Prometheus", Color: "E6522C"},
		{Name: "Grafana", Color: "F46800"},
		{Name: "Linux", Color: "FCC624"},
		{Name: "OpenAI", Color: "412991"},
	}

	// Build lookup map: every reasonable alias → canonical Badge
	m := make(map[string]models.Badge, len(entries)*2)
	aliases := map[string]string{
		"golang":       "Go",
		"js":           "JavaScript",
		"ts":           "TypeScript",
		"cpp":          "C++",
		"csharp":       "C#",
		"html":         "HTML5",
		"css":          "CSS3",
		"vue":          "Vue.js",
		"next":         "Next.js",
		"nextjs":       "Next.js",
		"nuxt":         "Nuxt.js",
		"express":      "Express.js",
		"nestjs":       "NestJS",
		"@nestjs/core": "NestJS",
		"spring":       "Spring Boot",
		"spring-boot":  "Spring Boot",
		"rails":        "Ruby on Rails",
		"k8s":          "Kubernetes",
		"postgres":     "PostgreSQL",
		"node":         "Node.js",
		"tailwind":     "TailwindCSS",
		"react-native": "React Native",
		"kafka":        "Apache Kafka",
		"grpc":         "gRPC",
		"bash":         "Shell",
	}

	// Index entries by lowercased canonical name
	byName := make(map[string]models.Badge, len(entries))
	for _, e := range entries {
		byName[strings.ToLower(e.Name)] = e
		m[strings.ToLower(e.Name)] = e
	}

	// Add aliases
	for alias, canonical := range aliases {
		if b, ok := byName[strings.ToLower(canonical)]; ok {
			m[strings.ToLower(alias)] = b
		}
	}

	return m
}

// toLogoSlug converts a badge display name to its shields.io simple-icons slug.
func toLogoSlug(name string) string {
	special := map[string]string{
		"C++":           "cplusplus",
		"C#":            "csharp",
		"Vue.js":        "vuedotjs",
		"Next.js":       "nextdotjs",
		"Nuxt.js":       "nuxtdotjs",
		"Express.js":    "express",
		"Node.js":       "nodedotjs",
		"Ruby on Rails": "rubyonrails",
		"React Native":  "react",
		"Spring Boot":   "springboot",
		"Apache Kafka":  "apachekafka",
		"TailwindCSS":   "tailwindcss",
		"HTML5":         "html5",
		"CSS3":          "css3",
		"gRPC":          "grpc",
		"GCP":           "googlecloud",
		"AWS":           "amazonaws",
		"Drizzle":       "drizzle",
		"Neon":          "neon",
		"Vite":          "vite",
		"Astro":         "astro",
		"Gatsby":        "gatsby",
		"OpenAI":        "openai",
		"Prometheus":    "prometheus",
		"Grafana":       "grafana",
		"Elasticsearch": "elasticsearch",
		"Supabase":      "supabase",
		"Firebase":      "firebase",
		"Prisma":        "prisma",
	}
	if slug, ok := special[name]; ok {
		return slug
	}
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "+", "plus")
	s = strings.ReplaceAll(s, "#", "sharp")
	return s
}

// Helpers – aggregate project data

// collectTopLanguages sums bytes per language across all repos, returns top-n names.
func collectTopLanguages(projects []models.RepositoryAnalysis, n int) []string {
	totals := make(map[string]int)
	for _, p := range projects {
		for lang, b := range p.Languages {
			totals[lang] += b
		}
	}
	type pair struct {
		name  string
		bytes int
	}
	sorted := make([]pair, 0, len(totals))
	for name, b := range totals {
		sorted = append(sorted, pair{name, b})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].bytes > sorted[j].bytes })
	result := make([]string, 0, n)
	for i, p := range sorted {
		if i >= n {
			break
		}
		result = append(result, p.name)
	}
	return result
}

// collectAllTopics gathers unique topics across all repos.
func collectAllTopics(projects []models.RepositoryAnalysis) []string {
	seen := make(map[string]bool)
	var topics []string
	for _, p := range projects {
		if p.Repository == nil {
			continue
		}
		for _, t := range p.Repository.Topics {
			if !seen[t] {
				seen[t] = true
				topics = append(topics, t)
			}
		}
	}
	return topics
}

// portfolioURL returns the best available website URL for the user.
func portfolioURL(config *models.ProfileConfig, user *models.User) string {
	if config.ContactPrefs.PersonalWebsite != "" {
		return config.ContactPrefs.PersonalWebsite
	}
	return user.Blog
}

// buildTypingLines creates URL-encoded lines for the readme-typing-svg service.
func buildTypingLines(config *models.ProfileConfig, topLangs, topics []string) []string {
	encode := func(s string) string {
		return strings.ReplaceAll(strings.TrimSpace(s), " ", "+")
	}
	var lines []string

	if config.TargetRole != "" {
		lines = append(lines, encode(config.TargetRole))
	}
	switch len(topLangs) {
	case 0:
		// nothing
	case 1:
		lines = append(lines, topLangs[0]+"+Developer")
	default:
		lines = append(lines, encode(topLangs[0])+" & "+encode(topLangs[1])+" Developer")
	}
	if len(config.SkillsEmphasis) > 0 {
		lines = append(lines, "Expert+in+"+encode(config.SkillsEmphasis[0]))
	}
	for _, t := range topics {
		if len(lines) >= 5 {
			break
		}
		lines = append(lines, encode(t))
	}
	if len(lines) == 0 {
		lines = append(lines, "Software+Developer")
	}
	return lines
}

// Markdown builder

// buildMarkdown assembles the README from all pipeline data — no hardcoded content.
func (s *ProfileService) buildMarkdown(
	user *models.User,
	req *models.ContentGenerationRequest,
	pitch string,
	summaries []models.ProjectSummary,
	badges []models.Badge,
	config *models.ProfileConfig,
) string {
	var md strings.Builder

	username := user.Username
	topLangs := collectTopLanguages(req.Projects, 5)
	allTopics := collectAllTopics(req.Projects)
	siteURL := portfolioURL(config, user)

	// Prefer contact email, fall back to account email
	contactEmail := config.ContactPrefs.Email
	if contactEmail == "" {
		contactEmail = user.Email
	}

	// ── HERO ──────────────────────────────────────────────────────────────
	md.WriteString("<div align=\"center\">\n\n")

	if user.AvatarURL != "" {
		md.WriteString(fmt.Sprintf(
			"<img src=\"%s\" width=\"120\" height=\"120\" style=\"border-radius:50%%\" alt=\"@%s\" />\n\n",
			user.AvatarURL, username,
		))
	}

	md.WriteString(fmt.Sprintf(
		"# Hi, I'm @%s <img src=\"https://raw.githubusercontent.com/MartinHeinz/MartinHeinz/master/wave.gif\" width=\"28px\" />\n\n",
		username,
	))

	if user.Bio != "" {
		md.WriteString(fmt.Sprintf("**%s**\n\n", user.Bio))
	}

	// Location / company metadata
	var meta []string
	if user.Location != "" {
		meta = append(meta, "📍 "+user.Location)
	}
	if user.Company != "" {
		c := user.Company
		if strings.HasPrefix(c, "@") {
			meta = append(meta, fmt.Sprintf("🏢 [%s](https://github.com/%s)", c, strings.TrimPrefix(c, "@")))
		} else {
			meta = append(meta, "🏢 "+c)
		}
	}
	if len(meta) > 0 {
		md.WriteString(strings.Join(meta, " &nbsp;·&nbsp; ") + "\n\n")
	}

	// Typing SVG built from real data
	typingLines := buildTypingLines(config, topLangs, allTopics)
	md.WriteString(fmt.Sprintf(
		"[![Typing SVG](https://readme-typing-svg.demolab.com?font=Fira+Code&size=22&duration=3000&pause=1000&color=2E97F7&center=true&vCenter=true&width=650&height=80&lines=%s)](https://git.io/typing-svg)\n\n",
		strings.Join(typingLines, ";"),
	))

	// Profile counters
	md.WriteString(fmt.Sprintf("![Profile Views](https://komarev.com/ghpvc/?username=%s&label=Profile%%20Views&color=0e75b6&style=flat)\n", username))
	md.WriteString(fmt.Sprintf("[![Followers](https://img.shields.io/github/followers/%s?label=Followers&style=social)](https://github.com/%s?tab=followers)\n", username, username))
	md.WriteString(fmt.Sprintf("[![Stars](https://img.shields.io/github/stars/%s?label=Stars&style=social)](https://github.com/%s)\n\n", username, username))
	md.WriteString("</div>\n\n")

	// ── ABOUT ME ──────────────────────────────────────────────────────────
	md.WriteString("## 👨‍💻 About Me\n\n")
	md.WriteString(pitch)
	md.WriteString("\n\n")

	// Dynamic bullets sourced from real data
	if config.TargetRole != "" {
		md.WriteString(fmt.Sprintf("- 🎯 Growing as a **%s**\n", config.TargetRole))
	}
	if len(summaries) > 0 {
		names := make([]string, 0, min(2, len(summaries)))
		for _, sum := range summaries[:min(2, len(summaries))] {
			if sum.Repository != nil {
				names = append(names, fmt.Sprintf("[%s](%s)", sum.Repository.Name, sum.Repository.HTMLURL))
			}
		}
		if len(names) > 0 {
			md.WriteString(fmt.Sprintf("- 🔭 Currently building **%s**\n", strings.Join(names, " & ")))
		}
	}
	if len(config.SkillsEmphasis) > 0 {
		md.WriteString(fmt.Sprintf("- 🌱 Deepening expertise in **%s**\n",
			strings.Join(config.SkillsEmphasis[:min(2, len(config.SkillsEmphasis))], " & ")))
	} else if len(topLangs) > 0 {
		md.WriteString(fmt.Sprintf("- 🌱 Deepening expertise in **%s**\n", topLangs[0]))
	}
	if len(topLangs) > 0 {
		md.WriteString(fmt.Sprintf("- 💬 Ask me about **%s**\n",
			strings.Join(topLangs[:min(3, len(topLangs))], ", ")))
	}
	if contactEmail != "" {
		md.WriteString(fmt.Sprintf("- 📫 Reach me at **%s**\n", contactEmail))
	}
	md.WriteString("\n")

	// ── CONNECT ───────────────────────────────────────────────────────────
	md.WriteString("## 🌐 Connect\n\n")
	md.WriteString("<div align=\"center\">\n\n")

	if config.ContactPrefs.LinkedIn != "" {
		md.WriteString(fmt.Sprintf("[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](%s)\n", config.ContactPrefs.LinkedIn))
	}
	if config.ContactPrefs.Twitter != "" {
		md.WriteString(fmt.Sprintf("[![Twitter/X](https://img.shields.io/badge/Twitter-000000?style=for-the-badge&logo=x&logoColor=white)](%s)\n", config.ContactPrefs.Twitter))
	}
	if contactEmail != "" {
		md.WriteString(fmt.Sprintf("[![Email](https://img.shields.io/badge/Email-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:%s)\n", contactEmail))
	}
	if siteURL != "" {
		md.WriteString(fmt.Sprintf("[![Website](https://img.shields.io/badge/Website-FF5722?style=for-the-badge&logo=googlechrome&logoColor=white)](%s)\n", siteURL))
	}
	md.WriteString(fmt.Sprintf("[![GitHub](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/%s)\n\n", username))
	md.WriteString("</div>\n\n")

	// ── TECH STACK ────────────────────────────────────────────────────────
	if len(badges) > 0 {
		md.WriteString("## 🛠️ Tech Stack\n\n")
		md.WriteString("<div align=\"center\">\n\n")

		categories := s.organizeBadgesByCategory(badges)
		for _, cat := range []string{"Languages", "Frameworks & Libraries", "Databases", "Tools & Platforms"} {
			catBadges, ok := categories[cat]
			if !ok || len(catBadges) == 0 {
				continue
			}
			md.WriteString(fmt.Sprintf("**%s**\n\n", cat))
			for _, b := range catBadges {
				md.WriteString(fmt.Sprintf(
					"![%s](https://img.shields.io/badge/%s-%s?style=flat-square&logo=%s&logoColor=white) ",
					b.Name,
					strings.ReplaceAll(b.Name, " ", "%20"),
					b.Color,
					toLogoSlug(b.Name),
				))
			}
			md.WriteString("\n\n")
		}
		md.WriteString("</div>\n\n")
	}

	// ── GITHUB STATS ──────────────────────────────────────────────────────
	md.WriteString("## 📊 GitHub Stats\n\n")
	md.WriteString("<div align=\"center\">\n\n")
	md.WriteString(fmt.Sprintf(
		"![%s's stats](https://github-readme-stats.vercel.app/api?username=%s&show_icons=true&count_private=true&theme=tokyonight&hide_border=true)\n",
		username, username,
	))
	md.WriteString(fmt.Sprintf(
		"![Top langs](https://github-readme-stats.vercel.app/api/top-langs/?username=%s&layout=compact&theme=tokyonight&hide_border=true)\n\n",
		username,
	))
	md.WriteString(fmt.Sprintf(
		"![Streak](https://streak-stats.demolab.com?user=%s&theme=tokyonight&hide_border=true)\n\n",
		username,
	))
	md.WriteString(fmt.Sprintf(
		"[![Trophies](https://github-profile-trophy.vercel.app/?username=%s&theme=tokyonight&no-frame=true&margin-w=4)](https://github.com/ryo-ma/github-profile-trophy)\n\n",
		username,
	))
	md.WriteString("</div>\n\n")

	// ── FEATURED PROJECTS ─────────────────────────────────────────────────
	if len(summaries) > 0 {
		md.WriteString("## 🚀 Featured Projects\n\n")

		for i, sum := range summaries {
			if sum.Repository == nil {
				continue
			}
			repo := sum.Repository
			owner := username
			if strings.Contains(repo.FullName, "/") {
				owner = strings.Split(repo.FullName, "/")[0]
			}

			md.WriteString(fmt.Sprintf("### [%s](%s)\n\n", repo.Name, repo.HTMLURL))

			if repo.Description != "" {
				md.WriteString(fmt.Sprintf("> %s\n\n", repo.Description))
			}

			md.WriteString(fmt.Sprintf(
				"[![Repo Card](https://github-readme-stats.vercel.app/api/pin/?username=%s&repo=%s&theme=tokyonight&hide_border=true)](%s)\n\n",
				owner, repo.Name, repo.HTMLURL,
			))

			// LLM summary
			if sum.Summary != "" {
				md.WriteString(sum.Summary + "\n\n")
			}

			// Tech stack from LLM
			if len(sum.TechStack) > 0 {
				md.WriteString("**Tech:** ")
				for _, tech := range sum.TechStack {
					md.WriteString(fmt.Sprintf("`%s` ", tech))
				}
				md.WriteString("\n\n")
			}

			// Real stats from repo + analysis
			var stats []string
			if repo.StargazersCount > 0 {
				stats = append(stats, fmt.Sprintf("⭐ %d stars", repo.StargazersCount))
			}
			if repo.ForksCount > 0 {
				stats = append(stats, fmt.Sprintf("🍴 %d forks", repo.ForksCount))
			}
			// Commit + contributor counts from the RepositoryAnalysis
			if i < len(req.Projects) {
				analysis := req.Projects[i]
				if analysis.CommitCount > 0 {
					stats = append(stats, fmt.Sprintf("📝 %d commits", analysis.CommitCount))
				}
				if analysis.ContributorCount > 0 {
					stats = append(stats, fmt.Sprintf("👥 %d contributors", analysis.ContributorCount))
				}
				// Top 3 languages from actual language map
				topProjLangs := collectTopLanguages([]models.RepositoryAnalysis{analysis}, 3)
				if len(topProjLangs) > 0 {
					stats = append(stats, "🔤 "+strings.Join(topProjLangs, " / "))
				}
			}
			// Topics
			if len(repo.Topics) > 0 {
				shown := repo.Topics[:min(5, len(repo.Topics))]
				stats = append(stats, "🏷️ "+strings.Join(shown, ", "))
			}

			if len(stats) > 0 {
				md.WriteString(strings.Join(stats, " &nbsp;·&nbsp; ") + "\n\n")
			}

			if i < len(summaries)-1 {
				md.WriteString("---\n\n")
			}
		}
	}

	// ── ACTIVITY GRAPH ────────────────────────────────────────────────────
	md.WriteString("## 📈 Contribution Activity\n\n")
	md.WriteString("<div align=\"center\">\n\n")
	md.WriteString(fmt.Sprintf(
		"[![Activity Graph](https://github-readme-activity-graph.vercel.app/graph?username=%s&theme=tokyo-night&hide_border=true)](https://github.com/ashutosh00710/github-readme-activity-graph)\n\n",
		username,
	))
	md.WriteString("</div>\n\n")

	// ── FOOTER ────────────────────────────────────────────────────────────
	md.WriteString("---\n\n")
	md.WriteString("<div align=\"center\">\n\n")
	md.WriteString(fmt.Sprintf(
		"*Generated with [GitRight](https://github.com/%s) · ![](https://komarev.com/ghpvc/?username=%s&style=flat-square)*\n\n",
		username, username,
	))
	md.WriteString("</div>\n")

	return md.String()
}

// Badge categorisation

// organizeBadgesByCategory sorts badges into four display groups.
func (s *ProfileService) organizeBadgesByCategory(badges []models.Badge) map[string][]models.Badge {
	langSet := map[string]bool{
		"Go": true, "Python": true, "JavaScript": true, "TypeScript": true,
		"Rust": true, "Java": true, "Kotlin": true, "Swift": true,
		"C++": true, "C": true, "C#": true, "PHP": true, "Ruby": true,
		"Dart": true, "Scala": true, "Elixir": true, "Haskell": true,
		"Lua": true, "Shell": true, "HTML5": true, "CSS3": true,
	}
	frameworkSet := map[string]bool{
		"React": true, "Vue.js": true, "Angular": true, "Svelte": true,
		"Next.js": true, "Nuxt.js": true, "Gatsby": true, "Remix": true,
		"Astro": true, "TailwindCSS": true, "Vite": true,
		"Node.js": true, "Express.js": true, "Fastify": true, "NestJS": true,
		"Django": true, "Flask": true, "FastAPI": true, "Spring Boot": true,
		"Laravel": true, "Ruby on Rails": true, "Fiber": true, "Gin": true,
		"Echo": true, "Flutter": true, "React Native": true,
	}
	dbSet := map[string]bool{
		"PostgreSQL": true, "MySQL": true, "MongoDB": true, "Redis": true,
		"SQLite": true, "Cassandra": true, "Elasticsearch": true,
		"Supabase": true, "Firebase": true, "Neon": true,
		"Prisma": true, "Drizzle": true,
	}

	cats := map[string][]models.Badge{
		"Languages":              {},
		"Frameworks & Libraries": {},
		"Databases":              {},
		"Tools & Platforms":      {},
	}

	for _, b := range badges {
		switch {
		case langSet[b.Name]:
			cats["Languages"] = append(cats["Languages"], b)
		case frameworkSet[b.Name]:
			cats["Frameworks & Libraries"] = append(cats["Frameworks & Libraries"], b)
		case dbSet[b.Name]:
			cats["Databases"] = append(cats["Databases"], b)
		default:
			cats["Tools & Platforms"] = append(cats["Tools & Platforms"], b)
		}
	}

	// Drop empty buckets
	for k, v := range cats {
		if len(v) == 0 {
			delete(cats, k)
		}
	}
	return cats
}
