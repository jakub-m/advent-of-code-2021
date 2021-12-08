package main

import (
	"advent"
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
	score := 0
	input := parseLine(line)
	in := input.checkSegments
	// digit 1
	if selected := selectSegments(in, func(s segment) bool {
		return len(s) == 2
	}); len(selected) > 0 {
		score += len(selected)
	}
	// digit 4
	if selected := selectSegments(in, func(s segment) bool {
		return len(s) == 4
	}); len(selected) > 0 {
		score += len(selected)
	}
	// digit 7
	if selected := selectSegments(in, func(s segment) bool {
		return len(s) == 3
	}); len(selected) > 0 {
		score += len(selected)
	}
	// digit 8
	if selected := selectSegments(in, func(s segment) bool {
		return len(s) == 7
	}); len(selected) > 0 {
		score += len(selected)
	}
	return score
}

type parsedInput struct {
	allSegments   []segment
	checkSegments []segment
}

type segment string

func parseLine(line string) parsedInput {
	// todo sort strings
	parts := strings.Split(line, " | ")
	advent.Assertf(len(parts) == 2, "%s", line)
	return parsedInput{
		allSegments:   parseSegments(parts[0]),
		checkSegments: parseSegments(parts[1]),
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

func parseSegments(input string) []segment {
	segments := []segment{}
	for _, s := range strings.Split(input, " ") {
		segments = append(segments, segment(s))
	}
	return segments
}
