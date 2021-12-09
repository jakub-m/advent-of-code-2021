package main

import (
	"advent"
	"fmt"
	"io"
	"strconv"
)

type coord struct {
	x, y int
}

type landscape map[coord]int

func Calc(r io.Reader) (int, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}

	landscape := make(landscape)
	for iRow, row := range lines {
		for iCol, char := range row {
			height, err := strconv.Atoi(fmt.Sprintf("%c", char))
			if err != nil {
				return 0, err
			}
			c := coord{iCol, iRow}
			// fmt.Printf("%v=%d\n", c, height)
			landscape[c] = height
		}
	}

	score := 0
	for coord, height := range landscape {
		neighbors := landscape.getNeighbors(coord)
		heights := landscape.getMany(neighbors)
		// fmt.Printf("%v , %v, %v\n", coord, neighbors, heights)
		if height < advent.MinInt(heights) {
			score += height + 1
		}
	}

	return score, nil
}

func (d landscape) getMany(coords []coord) []int {
	values := []int{}
	for _, c := range coords {
		if val, ok := d[c]; ok {
			values = append(values, val)
		} else {
			panic(fmt.Sprintf("No such coord: %+v", c))
		}
	}
	return values
}

func (d landscape) getNeighbors(c coord) []coord {
	neighbors := []coord{}

	for _, c := range []coord{
		{x: c.x + 1, y: c.y},
		{x: c.x - 1, y: c.y},
		{x: c.x, y: c.y + 1},
		{x: c.x, y: c.y - 1},
	} {
		if _, ok := d[c]; ok {
			neighbors = append(neighbors, c)
		}
	}
	return neighbors
}
