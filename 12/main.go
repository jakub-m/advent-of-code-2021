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

	var traverseRec func(string, map[string]int, []string) int

	traverseRec = func(node string, visitCount map[string]int, path []string) int {
		// fmt.Printf("rec %s -> %v\n", node, graph[node])
		limit, hasLimit := visitLimit(node, visitCount)
		if hasLimit && visitCount[node] >= limit {
			return 0
		}

		if node == end {
			// path = append(path, end)
			// fmt.Printf("END %s\n", strings.Join(path, ","))
			return 1
		}

		if hasLimit {
			visitCount[node]++
		}

		c := 0
		for _, neigh := range graph[node] {
			c += traverseRec(neigh, visitCount, append(path, node))
		}
		if hasLimit {
			visitCount[node]--
		}
		return c
	}
	visitCount := make(map[string]int)
	return traverseRec(start, visitCount, []string{}), nil
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

func visitLimit(node string, visitCount map[string]int) (limit int, hasLimit bool) {
	if strings.ToUpper(node) == node {
		return 0, false
	}

	if node == start || node == end {
		return 1, true
	}

	for _, c := range visitCount {
		if c == 2 {
			return 1, true
		}
	}

	return 2, true

	//return 1, true

	// if len(node) == 1 {
	// 	return 1, true
	// }

	// if len(node) == 2 {
	// 	return 2, true
	// }

	// panic(fmt.Sprintf("visitLimit %s", node))
}
