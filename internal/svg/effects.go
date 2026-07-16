package svg

import (
	"fmt"
	"io"

	svg "github.com/ajstarks/svgo"
)

type Canvas struct {
	*svg.SVG
	w io.Writer
}

func NewCanvas(w io.Writer, width, height int) *Canvas {
	c := svg.New(w)
	c.Startview(width, height, 0, 0, width, height)
	return &Canvas{SVG: c, w: w}
}

func (c *Canvas) WriteDefs(t *Theme) {
	c.SVG.Def()

	c.SVG.LinearGradient("accentGrad", 0, 0, 100, 0, []svg.Offcolor{
		{Offset: 0, Color: t.AccentStart, Opacity: 1},
		{Offset: 50, Color: t.AccentMid, Opacity: 1},
		{Offset: 100, Color: t.AccentEnd, Opacity: 1},
	})

	c.SVG.LinearGradient("textGrad", 0, 0, 100, 0, []svg.Offcolor{
		{Offset: 0, Color: t.AccentStart, Opacity: 1},
		{Offset: 50, Color: t.AccentMid, Opacity: 1},
		{Offset: 100, Color: t.AccentEnd, Opacity: 1},
	})

	c.SVG.LinearGradient("asciiGrad", 0, 0, 0, 100, []svg.Offcolor{
		{Offset: 0, Color: t.AccentMid, Opacity: 1},
		{Offset: 100, Color: t.AccentEnd, Opacity: 1},
	})

	c.SVG.LinearGradient("pillGrad", 0, 0, 100, 0, []svg.Offcolor{
		{Offset: 0, Color: t.AccentStart, Opacity: 0.15},
		{Offset: 100, Color: t.AccentEnd, Opacity: 0.15},
	})

	c.SVG.LinearGradient("pillStroke", 0, 0, 100, 0, []svg.Offcolor{
		{Offset: 0, Color: t.AccentStart, Opacity: 0.6},
		{Offset: 100, Color: t.AccentEnd, Opacity: 0.6},
	})

	c.SVG.LinearGradient("borderGrad", 0, 0, 100, 100, []svg.Offcolor{
		{Offset: 0, Color: t.AccentStart, Opacity: 0.4},
		{Offset: 50, Color: t.AccentMid, Opacity: 0.2},
		{Offset: 100, Color: t.AccentEnd, Opacity: 0.4},
	})

	c.SVG.LinearGradient("termGrad", 0, 0, 0, 100, []svg.Offcolor{
		{Offset: 0, Color: t.BgSecondary, Opacity: 0.9},
		{Offset: 100, Color: t.BgPrimary, Opacity: 0.95},
	})

	c.SVG.RadialGradient("bgGlow", 50, 50, 50, 50, 50, []svg.Offcolor{
		{Offset: 0, Color: t.AccentStart, Opacity: 0.12},
		{Offset: 100, Color: t.AccentStart, Opacity: 0},
	})

	c.SVG.RadialGradient("bgGlow2", 50, 50, 50, 50, 50, []svg.Offcolor{
		{Offset: 0, Color: t.AccentMid, Opacity: 0.08},
		{Offset: 100, Color: t.AccentMid, Opacity: 0},
	})

	c.SVG.RadialGradient("bgGlow3", 50, 50, 50, 50, 50, []svg.Offcolor{
		{Offset: 0, Color: t.AccentEnd, Opacity: 0.06},
		{Offset: 100, Color: t.AccentEnd, Opacity: 0},
	})

	c.SVG.LinearGradient("glowGrad", 0, 0, 100, 100, []svg.Offcolor{
		{Offset: 0, Color: t.AccentMid, Opacity: 0.3},
		{Offset: 100, Color: t.AccentEnd, Opacity: 0.1},
	})

	c.SVG.LinearGradient("scanGrad", 0, 0, 0, 100, []svg.Offcolor{
		{Offset: 0, Color: t.AccentMid, Opacity: 0},
		{Offset: 45, Color: t.AccentMid, Opacity: 0.05},
		{Offset: 50, Color: "#A5F3FC", Opacity: 0.65},
		{Offset: 55, Color: t.AccentMid, Opacity: 0.05},
		{Offset: 100, Color: t.AccentEnd, Opacity: 0},
	})

	c.SVG.Filter("glow", `-20% -20% 140% 140%`)
	c.SVG.FeGaussianBlur(svg.Filterspec{In: "SourceGraphic", Result: "blur"}, 4, 4)
	c.SVG.FeMerge([]string{"blur", "SourceGraphic"})
	c.SVG.Fend()

	c.SVG.Filter("softGlow", `-10% -10% 120% 120%`)
	c.SVG.FeGaussianBlur(svg.Filterspec{In: "SourceGraphic", Result: "blur"}, 2, 2)
	c.SVG.FeMerge([]string{"blur", "SourceGraphic"})
	c.SVG.Fend()

	c.SVG.Filter("glassBlur", `-5% -5% 110% 110%`)
	c.SVG.FeGaussianBlur(svg.Filterspec{In: "SourceGraphic", Result: "blur"}, 6, 6)
	c.SVG.Fend()

	c.SVG.Pattern("scanlines", 0, 0, 4, 4, "user")
	c.SVG.Rect(0, 0, 4, 1, `fill="#7DD3FC"`, `opacity="0.05"`)
	c.SVG.DefEnd()

	c.SVG.DefEnd()
}

func (c *Canvas) DrawBackground(t *Theme, w, h int) {
	c.SVG.Rect(0, 0, w, h, fmt.Sprintf(`fill="%s"`, t.BgPrimary))

	c.SVG.Circle(200, 150, 300, `fill="url(#bgGlow)"`)
	c.SVG.Circle(900, 450, 350, `fill="url(#bgGlow2)"`)
	c.SVG.Circle(600, 300, 250, `fill="url(#bgGlow3)"`)

	for x := 0; x < w; x += 40 {
		c.SVG.Line(x, 0, x, h, fmt.Sprintf(`stroke="%s"`, t.GridLine), `stroke-width="0.5"`)
	}
	for y := 0; y < h; y += 40 {
		c.SVG.Line(0, y, w, y, fmt.Sprintf(`stroke="%s"`, t.GridLine), `stroke-width="0.5"`)
	}
}


