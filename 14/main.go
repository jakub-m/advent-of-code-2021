package main

import (
	"advent"
	"fmt"
	"io"
	"sort"
)

func Calc(r io.Reader, iters int) (int, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}

	template := lines[0]
	pairs := StringToMapOfPairs(template)

	// single iteration of the plymer
	for i := 0; i < iters; i++ {
		// fmt.Printf("%v\n", pairs)
		//newPairs := copyPairs(pairs)
		newPairs := make(map[pair]int)
		for _, line := range lines[2:] {
			extendPairMap(pairs, newPairs, line)
		}
		pairs = newPairs
	}
	// fmt.Printf("%v\n", pairs)

	counts := getElementCounts(pairs, template)
	minCount, maxCount := counts[0], counts[len(counts)-1]
	return maxCount - minCount + 1, nil
}

func extendPairMap(inputMap map[pair]int, outputMap map[pair]int, operation string) {
	runes := []rune(operation)
	left, right, toInsert := runes[0], runes[1], runes[6]

	originalPair := pair{left, right}
	leftPair := pair{left, toInsert}
	rightPair := pair{toInsert, right}

	if c := inputMap[originalPair]; c > 0 {
		//outputMap[originalPair] = 0
		// fmt.Printf("%s, %s\n", leftPair, rightPair)
		outputMap[leftPair] += inputMap[originalPair]
		outputMap[rightPair] += inputMap[originalPair]
	}

	// fmt.Printf("%s: %v\n", operation, outputMap)
}

func getElementCounts(m map[pair]int, template string) []int {
	letterCounts := make(map[rune]int)
	for p, c := range m {
		letterCounts[p.left] += c
		letterCounts[p.right] += c
	}
	r := []rune(template)
	letterCounts[r[0]]--
	letterCounts[r[len(r)-1]]--

	counts := []int{}
	// fmt.Println(letterCounts)
	for _, c := range letterCounts {
		counts = append(counts, c/2)
	}
	sort.Ints(counts)
	return counts
}

type pair struct {
	left, right rune
}

func (p pair) String() string {
	return fmt.Sprintf("%c%c", p.left, p.right)
}

func StringToMapOfPairs(template string) map[pair]int {
	m := make(map[pair]int)
	runes := []rune(template)
	for i := 0; i < len(runes)-1; i++ {
		p := pair{runes[i], runes[i+1]}
		m[p]++
	}
	return m
}
