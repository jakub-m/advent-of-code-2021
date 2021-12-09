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

	// score := 0
	// for coord, height := range landscape {
	// 	neighbors := landscape.getNeighbors(coord)
	// 	heights := landscape.getMany(neighbors)
	// 	// fmt.Printf("%v , %v, %v\n", coord, neighbors, heights)
	// 	if height < advent.MinInt(heights) {
	// 		score += height + 1
	// 	}
	// }

	currBasinNumber := 0
	basinAnnotations := make(map[coord]int)

	var fillBasinRec func(start coord, basinNumber int)

	fillBasinRec = func(start coord, basinNumber int) {
		// fmt.Printf("rec %v %d\n", start, basinNumber)
		if height := landscape[start]; height == 9 {
			return
		}
		if num, ok := basinAnnotations[start]; ok {
			if num == basinNumber {
				return
			} else {
				panic(fmt.Sprintf("Field %v already annotated %d, wanted %d", start, num, basinNumber))
			}
		}
		basinAnnotations[start] = basinNumber
		for _, c := range landscape.getNeighbors(start) {
			fillBasinRec(c, basinNumber)
		}
	}

	for coord, height := range landscape {
		if _, ok := basinAnnotations[coord]; height == 9 || ok {
			continue
		}

		currBasinNumber++
		fillBasinRec(coord, currBasinNumber)
	}

	basinSize := make(map[int]int)
	for _, b := range basinAnnotations {
		basinSize[b]++
	}

	// fmt.Printf("%+v\n", basinSize)
	// fmt.Printf("%+v\n", basinAnnotations)

	score := 1
	for _, v := range basinSize {
		score *= v
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
