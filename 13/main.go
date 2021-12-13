package main

import (
	"advent"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
)

func Calc(r io.Reader) (string, error) {
	g, folds := readOrigami(r)

	for i, fold := range folds {
		if i == 1 {
			fmt.Printf("number of dots: %d: %d\n", i, len(g))
		}
		// fmt.Printf("fold:\n%s\n\n", g)
		newGrid := grid{}
		if fold.dir == foldDirX {
			for _, p := range g {
				x := p.x
				if x > fold.pos {
					x = fold.pos - (x - fold.pos)
				}
				q := pos{x: x, y: p.y}
				newGrid = append(newGrid, q)
			}

		} else if fold.dir == foldDirY {
			for _, p := range g {
				y := p.y
				if y > fold.pos {
					y = fold.pos - (y - fold.pos)
				}
				q := pos{x: p.x, y: y}
				newGrid = append(newGrid, q)
			}
		} else {
			panic(fmt.Sprintf("fold %v", fold))
		}
		g = newGrid
	}

	return g.String(), nil
}

var reFold = regexp.MustCompile(`fold along (x|y)=(\d+)`)

func readOrigami(r io.Reader) (grid, []fold) {
	lines, err := advent.ReadLinesTrim(r)
	advent.PanicErr(err)

	positions := []pos{}
	folds := []fold{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		posInts, err := advent.ReadIntsFromStringSep(line, ",")
		if err == nil {
			advent.Assertf(len(posInts) == 2, "line: %s", line)
			positions = append(positions, pos{x: posInts[0], y: posInts[1]})
			continue
		}
		m := reFold.FindStringSubmatch(line)
		if m == nil {
			panic(fmt.Sprintf("invalind input: `%s`", line))
		}
		v, err := strconv.Atoi(m[2])
		advent.PanicErr(err)
		folds = append(folds, fold{v, m[1]})
	}

	return grid(positions), folds
}

const (
	foldDirX = "x"
	foldDirY = "y"
)

type fold struct {
	pos int
	dir string
}

type pos struct {
	x, y int
}

type grid []pos

func (g grid) String() string {
	maxX, maxY := 0, 0
	for _, p := range g {
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	rows := [][]int{}

	for iRow := 0; iRow < maxY+1; iRow++ {
		row := []int{}
		for _, p := range g {
			if p.y == iRow {
				row = append(row, p.x)
			}
		}
		sort.Ints(row)
		rows = append(rows, row)
	}

	s := ""
	for _, row := range rows {
		if len(row) != 0 {
			max := advent.MaxInt(row)
			fields := make([]bool, max+1)
			for _, v := range row {
				fields[v] = true
			}
			for _, v := range fields {
				if v {
					s += "#"
				} else {
					s += " "
				}
			}
		}
		s += "\n"
	}
	return s
}
