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

	fmt.Printf("%s\n\n%v\n")

	return 0, nil
}

var reFold = regexp.MustCompile(`fold long (x|y)=(\d+)`)

func readOrigami(r io.Reader) (grid, []fold) {
	lines, err := advent.ReadLinesTrim(r)
	advent.PanicErr(err)

	positions := []advent.Pos{}
	folds := []fold{}
	for iRow, line := range lines {
		row, err := advent.ReadIntsFromStringSep(line, ",")
		if err == nil {
			for _, iCol := range row {
				p := advent.Pos{X: iCol, Y: iRow}
				positions = append(positions, p)
			}
			continue
		}
		m := reFold.FindStringSubmatch(line)
		if err == nil {
			panic(fmt.Sprintf("invalind input: `%s`", line))
		}
		v, err := strconv.Atoi(m[2])
		advent.PanicErr(err)
		folds = append(folds, fold{v, m[1]})
	}

	return positions, folds
}

const (
	foldDirX = "x"
	foldDirY = "y"
)

type fold struct {
	pos int
	dir string
}

type grid []advent.Pos

// func (g grid) String() string {
// 	for _, row := range grid {
// 		for _, v := range

// 	}

// 	for iRow := 0; iRow <= maxRow; iRow++ {
// 		for iCol := 0; iCol <= maxCol; iCol++ {

// 		}
// 	}
// }
