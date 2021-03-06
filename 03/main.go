package main

import (
	"advent"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type bitStats struct {
	countsForOne map[int]int
	nSamples     int
	nBits        int
}

func Calc(r io.Reader) int {
	rawLines, err := advent.ParseLinesReader(r, func(s string) (interface{}, error) {
		s = strings.Trim(s, " \t\n")
		return s, nil
	})
	advent.PanicErr(err)
	lines := []string{}
	for _, r := range rawLines {
		lines = append(lines, r.(string))
	}

	oxygenLines := make([]string, len(lines))
	copy(oxygenLines, lines)

	bitStats := countBitStats(lines)

	for i := 0; i < bitStats.nBits; i++ {
		bitStats = countBitStats(oxygenLines)
		oxygenLines = filterOxygenLines(oxygenLines, bitStats, i)
		fmt.Printf("%d: %+v\n", i, oxygenLines)
		if len(oxygenLines) <= 1 {
			break
		}
	}

	fmt.Println(oxygenLines)
	advent.Assertf(len(oxygenLines) == 1, "%d", len(oxygenLines))
	oxygenRate, err := strconv.ParseInt(oxygenLines[0], 2, 32)
	advent.PanicErr(err)

	scrubberLines := make([]string, len(lines))
	copy(scrubberLines, lines)

	for i := 0; i < bitStats.nBits; i++ {
		bitStats = countBitStats(scrubberLines)
		scrubberLines = filterScrubberLines(scrubberLines, bitStats, i)
		fmt.Printf("%d: %+v\n", i, scrubberLines)
		if len(scrubberLines) <= 1 {
			break
		}
	}

	fmt.Println(scrubberLines)
	advent.Assertf(len(scrubberLines) == 1, "%d", len(scrubberLines))
	scrubberRate, err := strconv.ParseInt(scrubberLines[0], 2, 32)
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
	return int(oxygenRate) * int(scrubberRate)
}

func filterOxygenLines(lines []string, bitStats bitStats, iBit int) []string {
	cOne := bitStats.countsForOne[iBit]
	cZero := bitStats.nSamples - cOne

	filtered := []string{}
	var popularBit int
	if cOne > cZero {
		popularBit = 1
	} else if cOne == cZero {
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

func filterScrubberLines(lines []string, bitStats bitStats, iBit int) []string {
	cOne := bitStats.countsForOne[iBit]
	cZero := bitStats.nSamples - cOne

	filtered := []string{}
	var popularBit int
	if cOne < cZero {
		popularBit = 1
	} else if cOne == cZero {
		popularBit = 0
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

func countBitStats(lines []string) bitStats {
	countsForOne := make(map[int]int)
	for _, line := range lines {
		for i, c := range line {
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
	nBits := len(lines[0])

	return bitStats{
		countsForOne: countsForOne,
		nSamples:     nSamples,
		nBits:        nBits,
	}
}
