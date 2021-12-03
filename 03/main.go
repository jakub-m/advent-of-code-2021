package main

import (
	"advent"
	"fmt"
	"io"
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
	gammaRate := 0
	epsilonRate := 0
	fmt.Printf("nBits=%d, nSamples=%d, %+v\n", nBits, nSamples, countsForOne)
	for i := 0; i < nBits; i++ {
		gammaRate = gammaRate << 1
		epsilonRate = epsilonRate << 1
		nOnes, ok := countsForOne[i]
		advent.Assert(ok, "i: %d", i)
		nZeros := nSamples - nOnes
		advent.Assert(nOnes != nZeros, "ones and zeros equal, i=%d, nBits=%d, %+v, ", i, nBits, countsForOne)
		if nOnes > nZeros {
			gammaRate = gammaRate | 1
		} else {
			epsilonRate = epsilonRate | 1
		}
	}

	return gammaRate * epsilonRate
}
