package main

import (
	"advent"
	"fmt"
	"io"
)

func Calc(r io.Reader) (int, error) {
	grid, err := advent.LoadGrid(r)
	if err != nil {
		return 0, err
	}

	for i := 0; i < 1000; i++ {
		flashCount := flashGrid(grid)
		if flashCount == len(grid) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("no solution")
}

func flashGrid(grid advent.GridInt) int {
	flashed := make(map[advent.Pos]bool)

	var flashRec func(advent.Pos)
	flashRec = func(pos advent.Pos) {
		if flashed[pos] {
			return
		}

		val := grid[pos]

		if val > 9 {
			flashed[pos] = true
			for _, p := range grid.GetNeighbors8(pos) {
				grid[p] += 1
				flashRec(p)
			}
		} else {
			return
		}
	}
	for p := range grid {
		flashRec(p)
	}

	flashCount := 0
	for p, v := range flashed {
		if v {
			flashCount++
			grid[p] = 0
		}
	}

	for pos, v := range grid {
		grid[pos] = v + 1
	}

	return flashCount
}
