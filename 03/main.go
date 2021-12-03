package main

import (
	"advent"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Calc(r io.Reader) int {
	lines, err := advent.ParseLinesReader(r, func(s string) (interface{}, error) {
		s = strings.Trim(s, " \t\n")
		return s, nil
	})
	advent.PanicErr(err)

	countsForOne := make(map[int]int)
	for _, line := range lines {
		for i, c := range line.(string) {
			switch c {
			case '1':
				countsForOne[i]++
			case '0':
				countsForOne[i] += 0 // set value
			default:
				panic(fmt.Sprintf("char: `%c`, line: `%s`", c, line))
			}
		}
	}

	nSamples := len(lines)
	nBits := len(lines[0].(string))
	// epsilonRate := 0

	oxygenLines := []string{}
	for _, line := range lines {
		oxygenLines = append(oxygenLines, line.(string))
	}

	for i := 0; i < nBits; i++ {
		oxygenLines = filterOxygenLines(oxygenLines, countsForOne, nSamples, i)
		fmt.Printf("%d: %+v\n", i, oxygenLines)
		if len(oxygenLines) <= 1 {
			break
		}
	}

	fmt.Println(oxygenLines)
	advent.Assertf(len(oxygenLines) == 1, "%d", len(oxygenLines))
	strconv.ParseInt(oxygenLines[0], 2, 32)
	advent.PanicErr(err)

	// // fmt.Printf("nBits=%d, nSamples=%d, %+v\n", nBits, nSamples, countsForOne)
	// for i := 0; i < nBits; i++ {
	// 	oxygenRate = oxygenRate << 1
	// 	nOnes, ok := countsForOne[i]
	// 	advent.Assertf(ok, "i: %d", i)
	// 	nZeros := nSamples - nOnes
	// 	advent.Assertf(nOnes != nZeros, "ones and zeros equal, i=%d, nBits=%d, %+v, ", i, nBits, countsForOne)
	// 	if nOnes > nZeros {
	// 		oxygenRate = oxygenRate | 1
	// 	}
	// }

	// return oxygenRate * epsilonRate
	return 0
}

func filterOxygenLines(lines []string, countsForOne map[int]int, nSamples int, iBit int) []string {
	cOne := countsForOne[iBit]
	cZero := nSamples - cOne

	filtered := []string{}
	var popularBit int
	if cOne >= cZero {
		popularBit = 1
	} else {
		popularBit = 0
	}

	fmt.Printf("popular bit %d, cOne %d, cZero %d\n", popularBit, cOne, cZero)
	for _, line := range lines {
		if bit := getBitFromString(line, iBit); bit == popularBit {
			filtered = append(filtered, line)
		}
	}
	return filtered
}

func getBitFromString(s string, i int) int {
	switch s[i] {
	case '0':
		return 0
	case '1':
		return 1
	}
	panic(s)
}
