package main

import (
	"advent"
	"io"
	"strings"
)

const (
	start = "start"
	end   = "end"
)

func Calc(r io.Reader) (int, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}
	graph := loadGraphFromLines(lines)

	var traverseRec func(string, map[string]bool) int

	traverseRec = func(node string, visited map[string]bool) int {
		//fmt.Printf("rec %s -> %v\n", node, graph[node])
		if visited[node] {
			return 0
		}

		if node == end {
			return 1
		}

		visitOnce := strings.ToLower(node) == node

		if visitOnce {
			visited[node] = true
		}

		c := 0
		for _, neigh := range graph[node] {
			c += traverseRec(neigh, visited)
		}
		if visitOnce {
			visited[node] = false
		}
		return c
	}
	visited := make(map[string]bool)
	return traverseRec(start, visited), nil
}

type graph map[string][]string

func loadGraphFromLines(lines []string) graph {
	graph := make(map[string][]string)

	for _, line := range lines {
		parts := strings.Split(line, "-")
		advent.Assertf(len(parts) == 2, "invalid line: %s", line)

		nodeLeft := parts[0]
		nodeRight := parts[1]

		if neighbours, ok := graph[nodeLeft]; ok {
			graph[nodeLeft] = append(neighbours, nodeRight)
		} else {
			graph[nodeLeft] = []string{nodeRight}
		}

		if neighbours, ok := graph[nodeRight]; ok {
			graph[nodeRight] = append(neighbours, nodeLeft)
		} else {
			graph[nodeRight] = []string{nodeLeft}
		}

	}

	return graph
}
