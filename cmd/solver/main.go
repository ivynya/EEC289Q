package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	tsp_solver "github.com/ivynya/EEC289Q/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func loadGraph(inputFile string) (*tsp_solver.Graph, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	// split file content by newlines
	var lines []string
	start := 0
	b := data
	for i := range b {
		if b[i] == '\n' {
			lines = append(lines, string(b[start:i]))
			start = i + 1
		}
	}
	// append the final segment (may be empty if file ends with a newline)
	lines = append(lines, string(b[start:]))

	// parse node count from the first line
	nodeCountClaim, err := strconv.ParseInt(lines[0], 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing node count")
		return nil, err
	}

	// build adjacency map for the graph: graph[from][to] = weight
	graph := tsp_solver.NewGraph(int(nodeCountClaim))

	for i, line := range lines {
		// skip the first two lines (metadata)
		if i < 2 {
			continue
		}
		// skip empty lines
		if line == "" {
			continue
		}

		to := 0
		from := 0
		weight := 0.0
		n, _ := fmt.Sscanf(line, "%d %d %f", &to, &from, &weight)
		if n != 3 {
			continue
		}
		graph.AddEdge(from, to, weight)
	}

	// warn user if their graph format may be incorrect
	if nodeCountClaim != int64(graph.NodeCount()) {
		log.Warn().Msgf("Warning: claimed node count %d mismatch actual %d", nodeCountClaim, graph.NodeCount())
	}

	return graph, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: solver <input_file>")
		return
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// performance measurement
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		log.Info().Msgf("Execution time: %s", elapsed)
	}()

	// get flags from input
	cpuFlag := flag.Int("cpu", -1, "Max CPU to use (default=-1 : all)")
	timeFlag := flag.Int("time", 59, "Time limit in seconds (default=59)")
	flag.Parse()
	if *cpuFlag <= 0 {
		*cpuFlag = runtime.NumCPU()
	}

	// read input file name to graph struct
	inputFile := flag.Args()[0]
	graph, err := loadGraph(inputFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load graph from file")
		return
	}

	// print summary of loaded graph
	ms_loaded := time.Since(startTime).Milliseconds()
	log.Info().Msgf("Loaded %d nodes / %d edges in %dms", graph.NodeCount(), graph.EdgeCount(), ms_loaded)

	// solve TSP and print result
	log.Info().Msgf("Using %d CPUs and %d seconds time limit", *cpuFlag, *timeFlag)
	path, dist, cycles := tsp_solver.SolveTSP(graph, *cpuFlag, *timeFlag)
	log.Info().Msgf("Best path found with cost %.4f (visited %d cycles): %v", dist, cycles, path)
}
