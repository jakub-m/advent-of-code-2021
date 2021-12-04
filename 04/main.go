package main

import (
	"advent"
	"io"
)

const (
	height = 5
	width  = 5
)

type bingoInput struct {
	numbers []int
	boards  []board
}

type board struct {
	fields [height][width]field
}

type field struct {
	val    int
	marked bool
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

func readBoards(lines []string) ([]board, error) {
	boards := []board{}
	lines = advent.TrimNextEmptyLines(lines)
	for len(lines) > 0 {
		board, err := readBoard(lines[:5])
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
		lines = lines[5:]
		lines = advent.TrimNextEmptyLines(lines)
	}

}

func findWinner(boards []board, numbers []int) (board, int, bool) {

}
