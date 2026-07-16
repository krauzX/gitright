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

	c.SVG.Rect(0, 0, BannerWidth, BannerHeight, `rx="18"`, `fill="url(#bgGlow)"`)
	c.SVG.Rect(0, 0, BannerWidth, BannerHeight, `rx="18"`, `fill="url(#scanlines)"`)

	drawTitleBar(c, data, theme)
	drawLeftPanel(c, data, theme)
	drawRightPanel(c, data, theme)
	drawScanlineSweep(c, BannerWidth, BannerHeight)
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

	c.SVG.Gend()
}

func drawRightPanel(c *Canvas, data *BannerData, t *Theme) {
	c.SVG.Group(`transform="translate(0,38)"`)

	c.SVG.Rect(508, 10, 655, 500, `rx="14"`,
		`fill="#0B1120"`, `fill-opacity="0.35"`,
		`stroke="url(#borderGrad)"`, `stroke-width="1"`, `opacity="0.35"`)

	c.SVG.Text(524, 24, "SYSTEM.INFO",
		`font-family="'Courier New', Consolas, monospace"`,
		`font-size="11"`, `fill="#38BDF8"`, `letter-spacing="2"`, `opacity="0.7"`)

	lines := buildTerminalLines(data)
	for i, line := range lines {
		y := 42 + i*22
		clipID := fmt.Sprintf("lc%d", i)
		c.SVG.ClipPath(clipID)
		c.SVG.Rect(500, y-18, 0, 24, fmt.Sprintf(`id="%s"`, clipID))
		c.SVG.DefEnd()

		c.SVG.Group(fmt.Sprintf(`clip-path="url(#%s)"`, clipID))
		c.SVG.Text(520, 0, "", `fill="#dbeafe"`)
		drawTerminalLineRaw(c, 520, y, line, t)
		c.SVG.Gend()
	}

	c.SVG.Rect(522, 491, 9, 16, `fill="#22D3EE"`, `opacity="0"`)

	c.SVG.Gend()
}

func buildTerminalLines(data *BannerData) []string {
	var lines []string

	lines = append(lines, fmt.Sprintf(".head %s%s", data.Login+"@devos", " -——————————————————————————————————————————-—-"))
	lines = append(lines, fmt.Sprintf(". Subject: %s %s", padDot(22), data.Name))
	lines = append(lines, fmt.Sprintf(". Role: %s %s", padDot(23), data.Bio))

	if data.Location != "" {
		lines = append(lines, fmt.Sprintf(". Origin: %s %s", padDot(21), data.Location))
	}
	if data.Company != "" {
		lines = append(lines, fmt.Sprintf(". Company: %s %s", padDot(20), data.Company))
	}

	lines = append(lines, "")

	if len(data.Skills) > 0 {
		lines = append(lines, fmt.Sprintf(". Core.Lang: %s %s", padDot(18), strings.Join(data.Skills[:min(4, len(data.Skills))], ", ")))
	}
	if len(data.Languages) > 0 {
		langs := make([]string, 0)
		for lang := range data.Languages {
			langs = append(langs, lang)
		}
		if len(langs) > 4 {
			langs = langs[:4]
		}
		lines = append(lines, fmt.Sprintf(". Languages: %s %s", padDot(18), strings.Join(langs, ", ")))
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf(".accent - Contact %s", "————————————————————————————————————————————-—-"))

	if data.Email != "" {
		lines = append(lines, fmt.Sprintf(". Grid.Mail: %s %s", padDot(15), data.Email))
	}
	if data.Website != "" {
		lines = append(lines, fmt.Sprintf(". Grid.Portfolio: %s %s", padDot(12), data.Website))
	}
	lines = append(lines, fmt.Sprintf(". Grid.Github: %s %s", padDot(14), data.Login))

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf(".accent - Live Stats %s", "————————————————————————————————————————————-—-"))
	lines = append(lines, fmt.Sprintf(". Contributions: %d  Followers: %d  Repos: %d", data.Contributions, data.Followers, data.Repos))

	return lines
}

func padDot(n int) string {
	return strings.Repeat(".", n)
}

func drawTerminalLineRaw(c *Canvas, x, y int, line string, t *Theme) {
	if strings.HasPrefix(line, ".head ") {
		parts := strings.SplitN(line, " ", 2)
		c.SVG.Text(x, y, parts[0][1:],
			`font-family="'Courier New', Consolas, monospace"`,
			`font-size="17"`, `fill="#7C3AED"`, `font-weight="bold"`)
		if len(parts) > 1 {
			c.SVG.Text(x+200, y, parts[1],
				`font-family="'Courier New', monospace"`,
				`font-size="15"`, `fill="#475569"`)
		}
		return
	}

	if strings.HasPrefix(line, ".accent ") {
		parts := strings.SplitN(line, " ", 2)
		c.SVG.Text(x, y, parts[0][1:],
			`font-family="'Courier New', monospace"`,
			`font-size="15"`, `fill="#10B981"`, `font-weight="bold"`)
		if len(parts) > 1 {
			c.SVG.Text(x+80, y, parts[1],
				`font-family="'Courier New', monospace"`,
				`font-size="15"`, `fill="#475569"`)
		}
		return
	}

	parts := strings.SplitN(line, ": ", 2)
	if len(parts) == 2 {
		label := parts[0]
		value := parts[1]

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
	} else {
		c.SVG.Text(x, y, line,
			`font-family="'Courier New', monospace"`,
			`font-size="15"`, `fill="#E5E7EB"`)
	}
}

func drawScanlineSweep(c *Canvas, w, h int) {
	c.SVG.Rect(0, -70, w, 70, `fill="url(#scanGrad)"`, `opacity="0.7"`)
	c.SVG.Animate("#scanSweep", "y", -70, h+10, 4.2, 0)
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
