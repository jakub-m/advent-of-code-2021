package main

import (
	"advent"
	"fmt"
	"io"
	"strings"
)

func Calc(r io.Reader) (int, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}
	totalScore := 0
	for _, line := range lines {
		score := calculateScoreFromString(line)
		totalScore += score
	}
	return totalScore, nil
}

func calculateScoreFromString(line string) int {
	input := parseLine(line)
	digitMap := decodeDigitsFromSegments(input.allSegments)
	score := 0
	for _, s := range input.checkSegments {
		score = score * 10
		score += digitMap.decode(s)
	}
	return score

}

type parsedInput struct {
	allSegments   []segment
	checkSegments []segment
}

type segment string

func (s segment) subCount(other segment) int {
	count := 0
	for _, cs := range s {
		count++
		for _, co := range other {
			if cs == co {
				count--
				break
			}

		}
	}
	return count
}

func parseLine(line string) parsedInput {
	// todo sort strings
	parts := strings.Split(line, " | ")
	advent.Assertf(len(parts) == 2, "%s", line)
	return parsedInput{
		allSegments:   parseSegments(parts[0]),
		checkSegments: parseSegments(parts[1]),
	}
}

func parseSegments(input string) []segment {
	segments := []segment{}
	for _, s := range strings.Split(input, " ") {
		segments = append(segments, segment(s))
	}
	return segments
}

type digitEncoding map[segment]int

func (d digitEncoding) decode(segment segment) int {

	for k, v := range d {
		if k.subCount(segment) == 0 && segment.subCount(k) == 0 {
			return v
		}
	}
	panic(fmt.Sprintf("cannot decode: %s with %+v", segment, d))
}

func decodeDigitsFromSegments(allDigits []segment) digitEncoding {
	digitEncoding := make(map[segment]int)

	d1 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 2
	})

	d4 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 4
	})

	d7 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 3
	})

	d8 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 7
	})

	d3 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 5 && s.subCount(d7) == 2
	})

	d6 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 6 &&
			s.subCount(d3) == 2 &&
			d3.subCount(s) == 1 &&
			d1.subCount(s) == 1
	})

	d5 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 5 && d6.subCount(s) == 1 && d1.subCount(s) == 1
	})

	d9 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 6 && s.subCount(d3) == 1 && d3.subCount(s) == 0
	})

	d2 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 5 && d4.subCount(s) == 2 && s.subCount(d4) == 3
	})

	d0 := selectOnly(allDigits, func(s segment) bool {
		return len(s) == 6 && s.subCount(d7) == 3 && d7.subCount(s) == 0 && d4.subCount(s) == 1
	})

	digitEncoding[d0] = 0
	digitEncoding[d1] = 1
	digitEncoding[d2] = 2
	digitEncoding[d3] = 3
	digitEncoding[d4] = 4
	digitEncoding[d5] = 5
	digitEncoding[d6] = 6
	digitEncoding[d7] = 7
	digitEncoding[d8] = 8
	digitEncoding[d9] = 9

	return digitEncoding
}

func selectOnly(segments []segment, cond func(segment) bool) segment {
	if selected := selectSegments(segments, cond); len(selected) == 1 {
		return selected[0]
	} else {
		panic(fmt.Sprintf("selected: %v", selected))
	}
}

func selectSegments(segments []segment, cond func(segment) bool) []segment {
	selected := []segment{}
	for _, s := range segments {
		if cond(s) {
			selected = append(selected, s)
		}
	}
	return selected
}
