package main

import (
	"advent"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func Calc(r io.Reader) (int, error) {
	grid, folds := readOrigami(r)

	fmt.Printf("%v\n\n%v\n", grid, folds)

	return 0, nil
}

var reFold = regexp.MustCompile(`fold along (x|y)=(\d+)`)

func readOrigami(r io.Reader) (grid, []fold) {
	lines, err := advent.ReadLinesTrim(r)
	advent.PanicErr(err)

	rows := [][]int{}
	folds := []fold{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		row, err := advent.ReadIntsFromStringSep(line, ",")
		if err == nil {
			rows = append(rows, row)
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

	return grid(rows), folds
}

const (
	foldDirX = "x"
	foldDirY = "y"
)

type fold struct {
	pos int
	dir string
}

type grid [][]int
