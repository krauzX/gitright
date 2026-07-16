package svg

type Theme struct {
	BgPrimary   string
	BgSecondary string
	BgPanel     string
	Border      string
	TextPrimary string
	TextMuted   string
	AccentStart string
	AccentMid   string
	AccentEnd   string
	GridLine    string
}

var DarkTheme = &Theme{
	BgPrimary:   "#030712",
	BgSecondary: "#0F172A",
	BgPanel:     "rgba(15,23,42,0.6)",
	Border:      "rgba(255,255,255,0.08)",
	TextPrimary: "#F8FAFC",
	TextMuted:   "#94A3B8",
	AccentStart: "#7C3AED",
	AccentMid:   "#22D3EE",
	AccentEnd:   "#10B981",
	GridLine:    "rgba(255,255,255,0.03)",
}

var LightTheme = &Theme{
	BgPrimary:   "#FFFFFF",
	BgSecondary: "#F8FAFC",
	BgPanel:     "rgba(248,250,252,0.8)",
	Border:      "rgba(15,23,42,0.08)",
	TextPrimary: "#0F172A",
	TextMuted:   "#475569",
	AccentStart: "#2563EB",
	AccentMid:   "#06B6D4",
	AccentEnd:   "#10B981",
	GridLine:    "rgba(15,23,42,0.04)",
}
