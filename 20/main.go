package main

import (
	"advent"
	"fmt"
	"io"
)

func Calc(r io.Reader) (int, error) {
	mask, board := parseInput(r)
	board = enhance(board, mask)
	board = enhance(board, mask)
	return countLitPixels(board), nil
}

func parseInput(r io.Reader) (mask, grid) {
	lines, err := advent.ReadLinesTrim(r)
	advent.PanicErr(err)
	mask := readMask(lines[0])
	grid := readGrid(lines[2:])
	return mask, grid
}

const (
	MaskLen = 512
)

type mask [MaskLen]bool

type grid map[advent.Pos]bool

func readMask(line string) mask {
	m := mask{}
	for i, c := range line {
		if c == '.' {
			m[i] = false
		} else if c == '#' {
			m[i] = true
		} else {
			panic(fmt.Sprintf("invalid character '%c'", c))
		}
	}
	return m
}

func readGrid(lines []string) grid {
	g := make(map[advent.Pos]bool)
	for iRow, line := range lines {
		for iCol, c := range line {
			p := advent.Pos{X: iCol, Y: iRow}
			if c == '.' {
				g[p] = false
			} else if c == '#' {
				g[p] = true
			} else {
				panic(fmt.Sprintf("invalid character '%c'", c))
			}
		}
	}
	return g
}

func enhance(g grid, m mask) grid {
	minX, minY, maxX, maxY := getMinMaxCoords(g)

	valueAt := func(p advent.Pos, offsetX, offsetY int) int {
		q := advent.Pos{X: p.X + offsetX, Y: p.Y + offsetY}
		if g[q] {
			return 1
		} else {
			return 0
		}
	}

	newGrid := make(grid)
	for y := minY - 1; y <= maxY+1; y++ {
		for x := minX - 1; x <= maxX+1; x++ {
			p := advent.Pos{X: x, Y: y}
			index := 0
			index |= valueAt(p, -1, -1)
			index <<= 1
			index |= valueAt(p, 0, -1)
			index <<= 1
			index |= valueAt(p, 1, -1)
			index <<= 1
			index |= valueAt(p, -1, 0)
			index <<= 1
			index |= valueAt(p, 0, 0)
			index <<= 1
			index |= valueAt(p, 1, 0)
			index <<= 1
			index |= valueAt(p, -1, 1)
			index <<= 1
			index |= valueAt(p, 0, 1)
			index <<= 1
			index |= valueAt(p, 1, 1)
			newGrid[p] = m[index]
		}
	}
	return newGrid
}

func getMinMaxCoords(g grid) (int, int, int, int) {
	xx := []int{}
	yy := []int{}
	for k, v := range g {
		if v {
			xx = append(xx, k.X)
			yy = append(yy, k.Y)
		}
	}
	minX := advent.MinInt(xx)
	minY := advent.MinInt(yy)
	maxX := advent.MaxInt(xx)
	maxY := advent.MaxInt(yy)
	return minX, minY, maxX, maxY
}

func countLitPixels(g grid) int {
	c := 0
	for _, v := range g {
		if v {
			c++
		}
	}
	return c
}
