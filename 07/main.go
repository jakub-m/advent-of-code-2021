package main

import (
	"advent"
)

func Calc(input string) (int, error) {
	initialPositions, err := advent.ReadIntsFromStringSep(input, ",")
	if err != nil {
		return 0, err
	}

	min, max := advent.MinInt(initialPositions), advent.MaxInt(initialPositions)

	minPositionCost := -1

	for tentativePos := min; tentativePos <= max; tentativePos++ {
		cost := calcMoveCost(initialPositions, tentativePos)
		if tentativePos == min || cost < minPositionCost {
			minPositionCost = cost
		}
	}
	return minPositionCost, nil
}

func calcMoveCost(positions []int, targetPos int) int {
	cost := 0
	for _, pos := range positions {
		n := advent.AbsInt(pos - targetPos)
		cost += (1 + n) * n / 2
	}
	return cost
}
