package svg

import (
	"fmt"
	"io"
)

type TimelineData struct {
	Name            string
	Login           string
	Repos           []TimelineRepo
	TotalContribs   int
	MemberSince     string
	Languages       []string
}

type TimelineRepo struct {
	Name       string
	Language   string
	Stars      int
	CreatedAt  string
}

func GenerateTimeline(w io.Writer, data *TimelineData, dark bool) {
	theme := DarkTheme
	if !dark {
		theme = LightTheme
	}

	c := NewCanvas(w, 1180, 300)
	c.WriteDefs(theme)
	c.DrawBackground(theme, 1180, 300)

	drawTimelinePanel(c, data, theme)

	c.SVG.End()
}

func drawTimelinePanel(c *Canvas, data *TimelineData, t *Theme) {
	c.SVG.Group(`transform="translate(30, 20)"`)

	c.SVG.Rect(0, 0, 1120, 260, `rx="12"`,
		fmt.Sprintf(`fill="%s"`, t.BgSecondary), `fill-opacity="0.35"`,
		`stroke="url(#borderGrad)"`, `stroke-width="1"`)

	c.SVG.Text(20, 28, "DEVELOPER.JOURNEY",
		`font-family="'Courier New', Consolas, monospace"`,
		`font-size="11"`, fmt.Sprintf(`fill="%s"`, t.AccentMid),
		`letter-spacing="2"`, `opacity="0.7"`)

	drawTimelineStats(c, data, t)

	drawTimelineLine(c, 20, 120, 1080, t)

	drawTimelineRepos(c, data, t)

	c.SVG.Gend()
}

func drawTimelineStats(c *Canvas, data *TimelineData, t *Theme) {
	stats := []struct {
		label string
		value string
	}{
		{"Member Since", data.MemberSince},
		{"Total Contribs", fmt.Sprintf("%d", data.TotalContribs)},
		{"Languages", fmt.Sprintf("%d", len(data.Languages))},
		{"Repos", fmt.Sprintf("%d", len(data.Repos))},
	}

	x := 20
	for _, s := range stats {
		c.SVG.Text(x, 55, s.label,
			`font-size="9"`, `font-family="'SF Mono', monospace"`,
			fmt.Sprintf(`fill="%s"`, t.TextMuted))
		c.SVG.Text(x, 72, s.value,
			`font-size="14"`, `font-family="'SF Pro Display', 'Inter', system-ui, sans-serif"`,
			`font-weight="bold"`, `fill="url(#textGrad)"`)
		x += 260
	}

	if len(data.Languages) > 0 {
		y := 95
		x := 20
		for i, lang := range data.Languages {
			if i >= 8 {
				break
			}
			w := len(lang)*7 + 16
			if x+w > 1080 {
				x = 20
				y += 24
			}
			c.SVG.Roundrect(x, y, w, 18, 9, 9,
				`fill="url(#pillGrad)"`, `stroke="url(#pillStroke)"`, `stroke-width="0.5"`)
			c.SVG.Text(x+w/2, y+13, lang,
				`font-size="8"`, `font-family="'SF Mono', monospace"`,
				fmt.Sprintf(`fill="%s"`, t.AccentMid), `text-anchor="middle"`)
			x += w + 4
		}
	}
}

func drawTimelineLine(c *Canvas, x1, y, x2 int, t *Theme) {
	c.SVG.Line(x1, y, x2, y,
		fmt.Sprintf(`stroke="%s"`, t.Border), `stroke-width="2"`)

	c.SVG.Circle(x1, y, 5, `fill="url(#accentGrad)"`)
	c.SVG.Circle(x2, y, 5, `fill="url(#accentGrad)"`)
}

func drawTimelineRepos(c *Canvas, data *TimelineData, t *Theme) {
	if len(data.Repos) == 0 {
		return
	}

	x := 40
	maxRepos := min(6, len(data.Repos))
	step := 1000 / maxRepos

	for i := 0; i < maxRepos; i++ {
		repo := data.Repos[i]
		pos := x + i*step

		c.SVG.Circle(pos, 120, 4, `fill="url(#accentGrad)"`)

		c.SVG.Line(pos, 124, pos, 155,
			fmt.Sprintf(`stroke="%s"`, t.Border), `stroke-width="1"`)

		name := repo.Name
		if len(name) > 15 {
			name = name[:12] + "..."
		}
		c.SVG.Text(pos, 170, name,
			`font-size="9"`, `font-family="'SF Mono', monospace"`,
			fmt.Sprintf(`fill="%s"`, t.TextPrimary), `text-anchor="middle"`)

		if repo.Language != "" {
			c.SVG.Text(pos, 185, repo.Language,
				`font-size="8"`, `font-family="'SF Mono', monospace"`,
				fmt.Sprintf(`fill="%s"`, t.AccentMid), `text-anchor="middle"`)
		}

		if repo.Stars > 0 {
			c.SVG.Text(pos, 200, fmt.Sprintf("★ %d", repo.Stars),
				`font-size="8"`, `font-family="'SF Mono', monospace"`,
				fmt.Sprintf(`fill="%s"`, t.TextMuted), `text-anchor="middle"`)
		}

		if repo.CreatedAt != "" {
			c.SVG.Text(pos, 215, repo.CreatedAt,
				`font-size="7"`, `font-family="'SF Mono', monospace"`,
				fmt.Sprintf(`fill="%s"`, t.TextMuted), `text-anchor="middle"`, `opacity="0.6"`)
		}
	}
}
