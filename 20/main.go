package main

import (
	"advent"
	"fmt"
	"io"
)

func Calc(r io.Reader) (int, error) {
	mask, board := parseInput(r)
	// fmt.Println(board)
	board = enhance(board, mask)
	// fmt.Println(board)
	board = enhance(board, mask)
	// fmt.Println(board)
	return board.countLitPixels(), nil
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

type grid struct {
	fields       map[advent.Pos]bool
	defaultValue bool
}

func newGrid() grid {
	return grid{
		fields:       make(map[advent.Pos]bool),
		defaultValue: false,
	}
}

func (g grid) get(p advent.Pos) bool {
	if v, ok := g.fields[p]; ok {
		return v
	} else {
		return g.defaultValue
	}
}

func (g *grid) set(p advent.Pos, v bool) {
	g.fields[p] = v
}

func (g grid) String() string {
	margin := 3
	s := ""
	minX, minY, maxX, maxY := g.getMinMaxCoords()
	for y := minY - margin; y <= maxY+margin; y++ {
		for x := minX - margin; x <= maxX+margin; x++ {
			if g.get(advent.Pos{X: x, Y: y}) {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

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
	g := newGrid()
	for iRow, line := range lines {
		for iCol, c := range line {
			p := advent.Pos{X: iCol, Y: iRow}
			if c == '.' {
				g.set(p, false)
			} else if c == '#' {
				g.set(p, true)
			} else {
				panic(fmt.Sprintf("invalid character '%c'", c))
			}
		}
	}
	return g
}

func enhance(g grid, m mask) grid {
	minX, minY, maxX, maxY := g.getMinMaxCoords()

	valueAt := func(p advent.Pos, offsetX, offsetY int) int {
		q := advent.Pos{X: p.X + offsetX, Y: p.Y + offsetY}
		if g.get(q) {
			return 1
		} else {
			return 0
		}
	}

	newGrid := newGrid()
	newGrid.defaultValue = g.defaultValue
	if m[0] && !m[511] {
		// flip infinity
		newGrid.defaultValue = !newGrid.defaultValue
	}
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
			newGrid.set(p, m[index])
		}
	}
	return newGrid
}

func (g grid) getMinMaxCoords() (int, int, int, int) {
	xx := []int{}
	yy := []int{}
	for k, v := range g.fields {
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

func (g grid) countLitPixels() int {
	if g.defaultValue {
		panic("cannot count pixels if defaultValue is true")
	}
	c := 0
	for _, v := range g.fields {
		if v {
			c++
		}
	}
	return c
}
