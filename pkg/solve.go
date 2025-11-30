package tsp_solver

import (
	"context"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// SolveTSP computes an approximate solution to the Traveling Salesperson Problem.
// It uses a parallelized approach with randomized Nearest Neighbor initialization
// followed by 2-Opt local search optimization.
// The function targets a runtime of approximately 1 minute.
func SolveTSP(graph *Graph, maxCPU, maxSeconds int) ([]int, float64, int) {
	nodeCount := graph.NodeCount()
	if nodeCount == 0 {
		return []int{}, 0.0, 0
	}
	if nodeCount == 1 {
		return graph.nodes, 0.0, 1
	}

	// global best solution tracking
	var bestPath []int
	var totalCycles int
	bestCost := math.MaxFloat64
	var mu sync.Mutex

	// setup context with timeout
	timeout := time.Duration(maxSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	numWorkers := min(runtime.NumCPU(), maxCPU)
	var wg sync.WaitGroup

	// per-thread tsp solver
	worker := func(seed int64, workerID int) {
		defer wg.Done()
		rng := rand.New(rand.NewSource(seed))

		for {
			select {
			case <-ctx.Done():
				return
			default:
				startNode := graph.nodes[rng.Intn(nodeCount)]

				path, cost, ok := randomizedNearestNeighbor(startNode, nodeCount, graph, rng)
				// if failed to find a valid tour (e.g. disconnected graph), retry
				if !ok {
					continue
				}
				path, cost = twoOpt(path, cost, graph, ctx)

				mu.Lock()
				totalCycles++
				if cost < bestCost {
					log.Debug().Msgf("New best cost found by worker %d: %.4f", workerID, cost)
					bestCost = cost
					bestPath = make([]int, len(path))
					copy(bestPath, path)
				}
				mu.Unlock()
			}
		}
	}

	for i := range numWorkers {
		wg.Add(1)
		go worker(time.Now().UnixNano()+int64(i), i)
	}
	wg.Wait()

	// If no path found (e.g. disconnected graph), return empty
	if bestPath == nil {
		return []int{}, 0.0, totalCycles
	}
	return bestPath, bestCost, totalCycles
}

// randomizedNearestNeighbor constructs a path using a greedy approach with some randomness.
// At each step, it considers the k nearest unvisited neighbors and picks one randomly.
func randomizedNearestNeighbor(startNode int, numNodes int, graph *Graph, rng *rand.Rand) ([]int, float64, bool) {
	path := make([]int, 0, numNodes)
	visited := make(map[int]bool, numNodes)

	current := startNode
	path = append(path, current)
	visited[current] = true
	totalCost := 0.0

	type candidate struct {
		node int
		dist float64
	}

	for len(path) < numNodes {
		neighbors := graph.edges[current]

		// Find top 3 nearest unvisited neighbors
		var topK []candidate
		k := 3

		for n, w := range neighbors {
			if visited[n] {
				continue
			}

			// Maintain top K smallest distances
			if len(topK) < k {
				topK = append(topK, candidate{n, w})
			} else {
				// Find worst in topK
				maxDistIdx := 0
				for i := 1; i < len(topK); i++ {
					if topK[i].dist > topK[maxDistIdx].dist {
						maxDistIdx = i
					}
				}

				if w < topK[maxDistIdx].dist {
					topK[maxDistIdx] = candidate{n, w}
				}
			}
		}

		if len(topK) == 0 {
			return nil, 0, false // Dead end
		}

		// Pick random from topK
		choice := topK[rng.Intn(len(topK))]

		current = choice.node
		path = append(path, current)
		visited[current] = true
		totalCost += choice.dist
	}

	// Close the loop
	first := path[0]
	last := path[len(path)-1]
	if w, ok := graph.edges[last][first]; ok {
		totalCost += w
	} else {
		return nil, 0, false // Cannot close loop
	}

	return path, totalCost, true
}

// twoOpt improves the path by swapping edges.
func twoOpt(path []int, currentCost float64, graph *Graph, ctx context.Context) ([]int, float64) {
	size := len(path)
	improved := true

	for improved {
		select {
		case <-ctx.Done():
			return path, currentCost
		default:
		}

		improved = false
		for i := range size {
			for j := i + 2; j < size; j++ {
				// Avoid swapping adjacent edges (0,1) and (N-1,0) which is a no-op for cost
				if i == 0 && j == size-1 {
					continue
				}

				u1 := path[i]
				v1 := path[(i+1)%size]
				u2 := path[j]
				v2 := path[(j+1)%size]

				w1 := graph.edges[u1][v1]
				w2 := graph.edges[u2][v2]

				wNew1, ok1 := graph.edges[u1][u2]
				wNew2, ok2 := graph.edges[v1][v2]

				if ok1 && ok2 {
					delta := (wNew1 + wNew2) - (w1 + w2)
					if delta < -1e-9 {
						reverse(path, i+1, j)
						currentCost += delta
						improved = true
					}
				}
			}
		}
	}
	return path, currentCost
}

func reverse(path []int, i, j int) {
	for i < j {
		path[i], path[j] = path[j], path[i]
		i++
		j--
	}
}
