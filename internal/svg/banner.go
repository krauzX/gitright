package svg

import (
	"fmt"
	"io"
	"strings"
)

const (
	BannerWidth  = 1180
	BannerHeight = 610
)

type BannerData struct {
	Name          string
	Login         string
	Bio           string
	AvatarURL     string
	Location      string
	Company       string
	Website       string
	Email         string
	Skills        []string
	Languages     map[string]int
	Contributions int
	Followers     int
	Repos         int
	PRReviews     int
	Issues        int
	SocialLinks   []SocialLink
	TopRepos      []TopRepoInfo
}

type TopRepoInfo struct {
	Name     string
	Stars    int
	Language string
}

type SocialLink struct {
	Provider string
	URL      string
}

func GenerateBanner(w io.Writer, data *BannerData, dark bool) {
	theme := DarkTheme
	if !dark {
		theme = LightTheme
	}

	c := NewCanvas(w, BannerWidth, BannerHeight)
	c.WriteDefs(theme)
	c.DrawBackground(theme, BannerWidth, BannerHeight)

	c.SVG.Rect(0, 0, BannerWidth, BannerHeight, `rx="18"`, `fill="url(#scanlines)"`)

	drawTitleBar(c, data, theme)
	drawLeftPanel(c, data, theme)
	drawRightPanel(c, data, theme)
	drawBorderFrame(c, BannerWidth, BannerHeight)

	c.SVG.End()
}

func drawTitleBar(c *Canvas, data *BannerData, t *Theme) {
	c.SVG.Rect(3, 3, 1174, 34, `rx="16"`, `fill="#0B1120"`, `fill-opacity="0.85"`)

	dotColors := []string{"#EF4444", "#F59E0B", "#10B981"}
	for i, color := range dotColors {
		c.SVG.Circle(24+i*18, 20, 5, fmt.Sprintf(`fill="%s"`, color))
	}

	title := fmt.Sprintf("%s@devos ~ %s", data.Login, "./profile.sh --live")
	c.SVG.Text(590, 25, title, `text-anchor="middle"`,
		`font-family="'Courier New', Consolas, monospace"`,
		`font-size="12"`, `fill="#64748B"`, `letter-spacing="0.5"`)

	c.SVG.Circle(1122, 20, 4, `fill="#F87171"`)
	c.SVG.Text(1132, 24, "SCANNING",
		`font-family="'Courier New', monospace"`,
		`font-size="10"`, `fill="#F87171"`, `letter-spacing="1"`)
}

func drawLeftPanel(c *Canvas, data *BannerData, t *Theme) {
	c.SVG.Group(`transform="translate(0,38)"`)

	c.SVG.Rect(14, 26, 488, 468, `rx="14"`,
		`fill="#0B1120"`, `fill-opacity="0.35"`,
		`stroke="url(#borderGrad)"`, `stroke-width="1"`, `opacity="0.35"`)

	c.SVG.Text(30, 24, "VISUAL.MAP",
		`font-family="'Courier New', Consolas, monospace"`,
		`font-size="11"`, `fill="#38BDF8"`, `letter-spacing="2"`, `opacity="0.7"`)

	if data.AvatarURL != "" {
		c.SVG.ClipPath("avatarClip")
		c.SVG.Circle(258, 230, 130, `id="avatarClip"`)
		c.SVG.DefEnd()

		c.SVG.Circle(258, 230, 135, `fill="none"`,
			`stroke="url(#borderGrad)"`, `stroke-width="2"`)

		c.SVG.Image(128, 100, 260, 260, data.AvatarURL,
			`clip-path="url(#avatarClip)"`,
			`preserveAspectRatio="xMidYMid slice"`)
	} else {
		drawFallbackASCII(c, data.Login)
	}

	if data.Name != "" {
		c.SVG.Text(258, 395, data.Name,
			`font-family="'Courier New', Consolas, monospace"`,
			`font-size="18"`, `fill="url(#textGrad)"`,
			`font-weight="bold"`, `text-anchor="middle"`,
			`filter="url(#softGlow)"`)
	}

	if data.Bio != "" {
		c.SVG.Text(258, 420, truncateStr(data.Bio, 50),
			`font-family="'Courier New', Consolas, monospace"`,
			`font-size="12"`, `fill="#94A3B8"`, `text-anchor="middle"`)
	}

	statsY := 450
	stats := []struct {
		label string
		value int
	}{
		{"repos", data.Repos},
		{"contribs", data.Contributions},
		{"followers", data.Followers},
	}
	for i, s := range stats {
		x := 70 + i*140
		c.SVG.Text(x, statsY, fmt.Sprintf("%d", s.value),
			`font-size="20"`, `font-family="'Courier New', monospace"`,
			`font-weight="bold"`, `fill="url(#textGrad)"`)
		c.SVG.Text(x, statsY+18, s.label,
			`font-size="10"`, `font-family="'Courier New', monospace"`,
			`fill="#475569"`)
	}

	c.SVG.Gend()
}

func drawFallbackASCII(c *Canvas, login string) {
	ascii := []string{
		"  ████████╗██╗  ██╗███████╗",
		"  ╚══██╔══╝██║  ██║██╔════╝",
		"     ██║   ███████║█████╗  ",
		"     ██║   ██╔══██║██╔══╝  ",
		"     ██║   ██║  ██║███████╗",
		"     ╚═╝   ╚═╝  ╚═╝╚══════╝",
	}
	startY := 160
	for i, line := range ascii {
		c.SVG.Text(80, startY+i*14, line,
			`font-family="'Courier New', Consolas, monospace"`,
			`font-size="10"`, `fill="url(#asciiGrad)"`,
			`letter-spacing="-0.2"`, `xml:space="preserve"`)
	}
}

func drawRightPanel(c *Canvas, data *BannerData, t *Theme) {
	c.SVG.Group(`transform="translate(0,38)"`)

	c.SVG.Rect(508, 10, 655, 500, `rx="14"`,
		`fill="#0B1120"`, `fill-opacity="0.35"`,
		`stroke="url(#borderGrad)"`, `stroke-width="1"`, `opacity="0.35"`)

	c.SVG.Text(524, 24, "SYSTEM.INFO",
		`font-family="'Courier New', Consolas, monospace"`,
		`font-size="11"`, `fill="#38BDF8"`, `letter-spacing="2"`, `opacity="0.7"`)

	y := 50
	drawTerminalLine(c, 520, y, "Subject", data.Name, t)
	y += 22
	drawTerminalLine(c, 520, y, "Login", data.Login, t)
	y += 22
	if data.Bio != "" {
		drawTerminalLine(c, 520, y, "Role", data.Bio, t)
		y += 22
	}
	if data.Location != "" {
		drawTerminalLine(c, 520, y, "Origin", data.Location, t)
		y += 22
	}
	if data.Company != "" {
		drawTerminalLine(c, 520, y, "Company", data.Company, t)
		y += 22
	}

	y += 6
	drawTerminalDivider(c, 520, y, "Skills", t)
	y += 20

	if len(data.Skills) > 0 {
		drawTerminalLine(c, 520, y, "Core.Lang", strings.Join(data.Skills[:min(6, len(data.Skills))], ", "), t)
		y += 22
	}

	if len(data.Languages) > 0 {
		langs := make([]string, 0)
		for lang := range data.Languages {
			langs = append(langs, lang)
		}
		if len(langs) > 4 {
			langs = langs[:4]
		}
		drawTerminalLine(c, 520, y, "Languages", strings.Join(langs, ", "), t)
		y += 22
	}

	y += 6
	drawTerminalDivider(c, 520, y, "Contact", t)
	y += 20

	if data.Email != "" {
		drawTerminalLine(c, 520, y, "Grid.Mail", data.Email, t)
		y += 22
	}
	if data.Website != "" {
		drawTerminalLine(c, 520, y, "Grid.Web", data.Website, t)
		y += 22
	}
	drawTerminalLine(c, 520, y, "Grid.Github", data.Login, t)
	y += 22

	y += 6
	drawTerminalDivider(c, 520, y, "Stats", t)
	y += 20

	drawTerminalLine(c, 520, y, "Contributions", fmt.Sprintf("%d", data.Contributions), t)
	y += 22
	drawTerminalLine(c, 520, y, "Followers", fmt.Sprintf("%d", data.Followers), t)
	y += 22
	drawTerminalLine(c, 520, y, "Repos", fmt.Sprintf("%d", data.Repos), t)
	y += 22
	drawTerminalLine(c, 520, y, "PR Reviews", fmt.Sprintf("%d", data.PRReviews), t)
	y += 22
	drawTerminalLine(c, 520, y, "Issues", fmt.Sprintf("%d", data.Issues), t)
	y += 26

	if len(data.TopRepos) > 0 {
		drawTerminalDivider(c, 520, y, "Top Repos", t)
		y += 20

		for i, repo := range data.TopRepos {
			if i >= 4 {
				break
			}
			repoName := repo.Name
			if len(repoName) > 30 {
				repoName = repoName[:27] + "..."
			}
			lang := repo.Language
			if lang == "" {
				lang = "—"
			}
			drawTerminalLine(c, 520, y, repoName, fmt.Sprintf("★%d · %s", repo.Stars, lang), t)
			y += 22
		}
	}

	c.SVG.Gend()
}

func drawTerminalLine(c *Canvas, x, y int, label, value string, t *Theme) {
	c.SVG.Text(x, y, ". ",
		`font-family="'Courier New', monospace"`, `font-size="15"`, `fill="#475569"`)
	c.SVG.Text(x+18, y, label,
		`font-family="'Courier New', Consolas, monospace"`,
		`font-size="15"`, `fill="#22D3EE"`, `font-weight="bold"`)
	c.SVG.Text(x+120, y, ": ",
		`font-family="'Courier New', monospace"`, `font-size="15"`, `fill="#475569"`)
	c.SVG.Text(x+135, y, truncateStr(value, 35),
		`font-family="'Courier New', Consolas, monospace"`,
		`font-size="15"`, `fill="#E5E7EB"`)
}

func drawTerminalDivider(c *Canvas, x, y int, title string, t *Theme) {
	c.SVG.Text(x, y, fmt.Sprintf("- %s", title),
		`font-family="'Courier New', monospace"`,
		`font-size="15"`, `fill="#10B981"`, `font-weight="bold"`)
	c.SVG.Text(x+80, y, "—————————————————————————————————————",
		`font-family="'Courier New', monospace"`,
		`font-size="15"`, `fill="#475569"`)
}

func drawBorderFrame(c *Canvas, w, h int) {
	c.SVG.Rect(3, 3, w-6, h-6, `rx="16"`,
		`fill="none"`, `stroke="url(#borderGrad)"`, `stroke-width="2"`, `opacity="0.8"`)
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
