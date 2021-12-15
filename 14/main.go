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

	head := stringToList(lines[0])

	// single iteration of the plymer
	for i := 0; i < iters; i++ {
		dirty := make(map[*node]bool)
		for _, line := range lines[2:] {
			extendList(head, line, dirty)
		}
		// fmt.Printf("%s\n", head.FullString())
	}

	counts := getCounts(head)
	minCount, maxCount := counts[0], counts[len(counts)-1]
	return maxCount - minCount, nil
}

type node struct {
	val  rune
	next *node
}

func (n *node) FullString() string {
	s := ""
	for ; n != nil; n = n.next {
		s += fmt.Sprintf("%c", n.val)
	}
	return s
}

func (n *node) insert(value rune) *node {
	oldNext := n.next
	newNext := &node{value, oldNext}
	n.next = newNext
	return newNext
}

func stringToList(template string) *node {
	var head *node
	var tail *node
	for _, r := range template {
		newTail := &node{r, nil}
		if head == nil {
			head = newTail
		} else {
			tail.next = newTail
		}
		tail = newTail
	}
	return head
}

func extendList(root *node, operation string, dirty map[*node]bool) {
	runes := []rune(operation)
	prev, next, toInsert := runes[0], runes[1], runes[6]

	cursor := root
	for cursor != nil {
		// fmt.Printf("%s , %s\n", cursor.FullString(), operation)
		if cursor.val == prev && cursor.next != nil && cursor.next.val == next && dirty[cursor] == false && dirty[cursor.next] == false {
			newNode := cursor.insert(rune(toInsert))
			dirty[newNode] = true
			cursor = newNode
		}
		cursor = cursor.next
	}
}

func getCounts(cursor *node) []int {
	stats := make(map[rune]int)
	for cursor != nil {
		stats[cursor.val]++
		cursor = cursor.next
	}

	counts := []int{}
	for _, count := range stats {
		counts = append(counts, count)
	}
	sort.Ints(counts)
	return counts
}
