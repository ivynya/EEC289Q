package tsp_solver

// graph is an undirected weighted graph struct
type Graph struct {
	edges map[int]map[int]float64
	nodes []int
}

func NewGraph(nCount int) *Graph {
	g := Graph{
		edges: make(map[int]map[int]float64, nCount),
		nodes: make([]int, 0, nCount),
	}
	return &g
}

func (g *Graph) AddEdge(from, to int, weight float64) {
	// edges are always stored from the larger node to the smaller node
	f := max(from, to)
	t := min(from, to)

	if g.edges[f] == nil {
		g.edges[f] = make(map[int]float64)
		g.nodes = append(g.nodes, f)
	}
	if g.edges[t] == nil {
		g.edges[t] = make(map[int]float64)
		g.nodes = append(g.nodes, t)
	}
	g.edges[f][t] = weight
	g.edges[t][f] = weight
}

func (g *Graph) NodeCount() int {
	return len(g.nodes)
}

func (g *Graph) EdgeCount() int {
	count := 0
	for _, edges := range g.edges {
		count += len(edges)
	}
	return count
}
