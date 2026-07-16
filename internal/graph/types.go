package graph

type Node struct {
	ID       string
	Label    string
	Category string
	Size     int
	Color    string
}

type Edge struct {
	Source string
	Target string
	Weight float64
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

type GraphConfig struct {
	Width       int
	Height      int
	BgColor     string
	NodeColors  map[string]string
	MinNodeSize int
	MaxNodeSize int
}
