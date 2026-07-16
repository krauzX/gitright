package graph

import (
	"fmt"
	"strings"

	"github.com/krauzx/gitright/internal/github"
)

func RenderContributionSVG(weeks []github.ContributionWeek, cfg GraphConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, cfg.Width, cfg.Height, cfg.Width, cfg.Height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="%s"/>`, cfg.Width, cfg.Height, cfg.BgColor))

	cellSize := 12
	padding := 20
	startX := padding
	startY := padding

	for weekIdx, week := range weeks {
		for dayIdx, day := range week.ContributionDays {
			x := float64(startX + weekIdx*(cellSize+2))
			y := float64(startY + dayIdx*(cellSize+2))
			color := contributionColor(day.ContributionCount)
			sb.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%d" height="%d" rx="2" fill="%s"/>`, x, y, cellSize, cellSize, color))
		}
	}

	sb.WriteString(fmt.Sprintf(`<text x="50%%" y="%d" text-anchor="middle" fill="#666" font-family="monospace" font-size="10">GitHub Contributions</text>`, cfg.Height-padding))

	sb.WriteString(`</svg>`)
	return sb.String()
}

func contributionColor(count int) string {
	switch {
	case count == 0:
		return "#161b22"
	case count <= 3:
		return "#0e4429"
	case count <= 6:
		return "#006d32"
	case count <= 9:
		return "#26a641"
	default:
		return "#39d353"
	}
}
