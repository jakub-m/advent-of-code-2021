package main

import (
	"advent"
	"fmt"
	"io"
	"strings"
)

const (
	height = 5
	width  = 5
)

type bingoInput struct {
	numbers []int
	boards  []*board
}

type board struct {
	fields [height][width]field
}

type field struct {
	val    int
	marked bool
}

func (b board) sumUnmarked() int {
	sum := 0
	for _, row := range b.fields {
		for _, f := range row {
			if !f.marked {
				sum += f.val
			}
		}
	}
	return sum
}

func (b *board) mark(val int) {
	for i, row := range b.fields {
		for j, f := range row {
			if f.val == val {
				// fmt.Printf("marked! %d\n", val)
				// fmt.Printf("%v\n", b.fields[i][j].marked)
				b.fields[i][j].marked = true
				// fmt.Printf("%v\n", b.fields[i][j].marked)
				// row[j] = field{val: f.val, marked: true}
			}
		}
	}
	// fmt.Println(b)
}

func (b board) isWinner() bool {
	for iRow := 0; iRow < height; iRow++ {
		wins := true
		for iCol := 0; iCol < width; iCol++ {
			f := b.fields[iRow][iCol]
			wins = wins && f.marked
		}
		if wins {
			return true
		}
	}

	for iCol := 0; iCol < height; iCol++ {
		wins := true
		for iRow := 0; iRow < width; iRow++ {
			f := b.fields[iRow][iCol]
			wins = wins && f.marked
		}
		if wins {
			return true
		}
	}

	return false
}

func (b board) String() string {
	rows := []string{}
	for iRow := 0; iRow < height; iRow++ {
		row := ""
		for iCol := 0; iCol < width; iCol++ {
			f := b.fields[iRow][iCol]
			m := " "
			if f.marked {
				m = "*"
			}
			row = row + fmt.Sprintf("%.2d%s ", f.val, m)
		}
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n")
}

func GetBingoScore(r io.Reader) int {
	bingoInput, err := loadBingoInput(r)
	advent.PanicErr(err)

	winner, lastNumber, ok := findWinner(bingoInput.boards, bingoInput.numbers)
	advent.Assertf(ok, "no winner")
	return lastNumber * winner.sumUnmarked()
}

func loadBingoInput(r io.Reader) (bingoInput, error) {
	null := bingoInput{}
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return null, err
	}
	numbers, err := advent.ReadIntsFromStringSep(lines[0], ",")
	if err != nil {
		return null, err
	}

	boards, err := readBoards(lines[2:])
	if err != nil {
		return null, err
	}

	return bingoInput{boards: boards, numbers: numbers}, nil
}

func readBoards(lines []string) ([]*board, error) {
	boards := []*board{}
	lines = advent.TrimNextEmptyLines(lines)
	for len(lines) > 0 {
		board, err := readBoard(lines[:5])
		if err != nil {
			return nil, err
		}
		boards = append(boards, &board)
		lines = lines[5:]
		lines = advent.TrimNextEmptyLines(lines)
	}

	return boards, nil
}

func readBoard(lines []string) (board, error) {
	null := board{}
	b := board{}
	if len(lines) != 5 {
		return null, fmt.Errorf("too much lines %d", len(lines))
	}
	for iRow, line := range lines {
		numbers, err := advent.ReadIntsFromStringSep(squashSpaces(line), " ")
		if err != nil {
			return null, err
		}
		for iCol, num := range numbers {
			b.fields[iRow][iCol] = field{num, false}
		}
	}
	return b, nil
}

func findWinner(boards []*board, numbers []int) (*board, int, bool) {
	for _, num := range numbers {
		for _, board := range boards {
			board.mark(num)
			// fmt.Printf("now\n%s\n\n", board)
			if board.isWinner() {
				return board, num, true
			}
		}
	}
	return nil, 0, false
}

func squashSpaces(s string) string {
	s = strings.ReplaceAll(s, "   ", " ")
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}
