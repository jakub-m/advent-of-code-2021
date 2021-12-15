package main

import (
	"advent"
	"io"
	"sort"
)

func Calc(r io.Reader, iters int) (int, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}

	pairs := stringToMapOfPairs(lines[0])

	// single iteration of the plymer
	for i := 0; i < iters; i++ {
		for _, line := range lines[2:] {
			extendPairMap(pairs, line)
		}
	}

	counts := getPairCounts(pairs)
	minCount, maxCount := counts[0], counts[len(counts)-1]
	return maxCount - minCount, nil
}

func extendPairMap(m map[pair]int, operation string) {
	runes := []rune(operation)
	left, right, toInsert := runes[0], runes[1], runes[6]

	originalPair := pair{left, right}
	leftPair := pair{left, toInsert}
	rightPair := pair{toInsert, right}

	if c := m[originalPair]; c > 0 {
		m[originalPair]--
	} else {
		return
	}
	m[leftPair]++
	m[rightPair]++
}

func getPairCounts(m map[pair]int) []int {
	counts := []int{}
	for _, count := range m {
		counts = append(counts, count)
	}
	sort.Ints(counts)
	return counts
}

type pair struct {
	left, right rune
}

func stringToMapOfPairs(template string) map[pair]int {

}
