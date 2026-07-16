package graph

import (
	"fmt"
	"math"
	"strings"
)

func RenderSVG(g *Graph, cfg GraphConfig) string {
	if len(g.Nodes) == 0 {
		return emptySVG(cfg)
	}

	positions := forceLayout(g, cfg)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, cfg.Width, cfg.Height, cfg.Width, cfg.Height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="%s"/>`, cfg.Width, cfg.Height, cfg.BgColor))

	sb.WriteString(`<defs>`)
	sb.WriteString(`<filter id="glow" x="-50%" y="-50%" width="200%" height="200%">`)
	sb.WriteString(`<feGaussianBlur stdDeviation="3" result="blur"/>`)
	sb.WriteString(`<feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>`)
	sb.WriteString(`</filter>`)
	sb.WriteString(`</defs>`)

	for _, edge := range g.Edges {
		src, ok1 := positions[edge.Source]
		tgt, ok2 := positions[edge.Target]
		if !ok1 || !ok2 {
			continue
		}
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f" stroke-opacity="0.3"/>`,
			src.X, src.Y, tgt.X, tgt.Y, cfg.NodeColors["edge"], math.Max(0.5, edge.Weight*2)))
	}

	for _, node := range g.Nodes {
		pos, ok := positions[node.ID]
		if !ok {
			continue
		}
		r := float64(node.Size)
		color := cfg.NodeColors[node.Category]
		if color == "" {
			color = cfg.NodeColors["default"]
		}
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" filter="url(#glow)"/>`, pos.X, pos.Y, r, color))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" fill="white" font-size="10" font-family="monospace">%s</text>`, pos.X, pos.Y+r+12, truncateLabel(node.Label, 20)))
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

func emptySVG(cfg GraphConfig) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d"><rect width="%d" height="%d" fill="%s"/><text x="50%%" y="50%%" text-anchor="middle" fill="#666" font-family="monospace" font-size="14">No data</text></svg>`,
		cfg.Width, cfg.Height, cfg.Width, cfg.Height, cfg.Width, cfg.Height, cfg.BgColor)
}

type Position struct {
	X, Y float64
}

func forceLayout(g *Graph, cfg GraphConfig) map[string]Position {
	positions := make(map[string]Position)
	n := len(g.Nodes)
	if n == 0 {
		return positions
	}

	cx, cy := float64(cfg.Width)/2, float64(cfg.Height)/2
	for i, node := range g.Nodes {
		angle := 2 * math.Pi * float64(i) / float64(n)
		radius := math.Min(float64(cfg.Width), float64(cfg.Height)) * 0.35
		positions[node.ID] = Position{
			X: cx + radius*math.Cos(angle),
			Y: cy + radius*math.Sin(angle),
		}
	}

	for iter := 0; iter < 100; iter++ {
		forces := make(map[string]Position)
		for _, node := range g.Nodes {
			forces[node.ID] = Position{0, 0}
		}

		for i := 0; i < len(g.Nodes); i++ {
			for j := i + 1; j < len(g.Nodes); j++ {
				a, b := g.Nodes[i], g.Nodes[j]
				dx := positions[a.ID].X - positions[b.ID].X
				dy := positions[a.ID].Y - positions[b.ID].Y
				dist := math.Max(math.Sqrt(dx*dx+dy*dy), 1)
				force := 5000.0 / (dist * dist)
				fx := force * dx / dist
				fy := force * dy / dist
				fa := forces[a.ID]
				fa.X += fx
				fa.Y += fy
				forces[a.ID] = fa
				fb := forces[b.ID]
				fb.X -= fx
				fb.Y -= fy
				forces[b.ID] = fb
			}
		}

		for _, edge := range g.Edges {
			src, ok1 := positions[edge.Source]
			tgt, ok2 := positions[edge.Target]
			if !ok1 || !ok2 {
				continue
			}
			dx := tgt.X - src.X
			dy := tgt.Y - src.Y
			dist := math.Max(math.Sqrt(dx*dx+dy*dy), 1)
			force := (dist - 100) * 0.01
			fx := force * dx / dist
			fy := force * dy / dist
			fs := forces[edge.Source]
			fs.X += fx
			fs.Y += fy
			forces[edge.Source] = fs
			ft := forces[edge.Target]
			ft.X -= fx
			ft.Y -= fy
			forces[edge.Target] = ft
		}

		for _, node := range g.Nodes {
			f := forces[node.ID]
			p := positions[node.ID]
			p.X += f.X * 0.1
			p.Y += f.Y * 0.1
			p.X = math.Max(float64(cfg.MinNodeSize+20), math.Min(float64(cfg.Width-cfg.MinNodeSize-20), p.X))
			p.Y = math.Max(float64(cfg.MinNodeSize+20), math.Min(float64(cfg.Height-cfg.MinNodeSize-20), p.Y))
			positions[node.ID] = p
		}
	}

	return positions
}

func truncateLabel(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-2] + ".."
	}
	return s
}
