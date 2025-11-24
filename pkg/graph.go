package tsp_solver

// graph is an undirected weighted graph struct
type Graph struct {
	Nodes map[int]struct{}
	Edges map[int][]Edge
}

type Edge struct {
	To     int
	Weight float64
}

func NewGraph(nCount int) *Graph {
	return &Graph{
		Nodes: make(map[int]struct{}, 0),
		Edges: make(map[int][]Edge, nCount),
	}
}

func (g *Graph) AddEdge(from, to int, weight float64) {
	// edges are always stored from the larger node to the smaller node
	f := max(from, to)
	t := min(from, to)

	// minimal allocation to log node existence
	g.Nodes[f] = struct{}{}
	g.Nodes[t] = struct{}{}

	// add the edge
	g.Edges[f] = append(g.Edges[f], Edge{To: t, Weight: weight})
}

func (g *Graph) NodeCount() int {
	return len(g.Nodes)
}

func (g *Graph) EdgeCount() int {
	count := 0
	for _, edges := range g.Edges {
		count += len(edges)
	}
	return count
}
