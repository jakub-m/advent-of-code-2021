package main

import (
	"advent"
	"fmt"
	"io"
)

func Calc(r io.Reader) (int, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}

	score := 0
	for _, line := range lines {
		s := getScoreForLine(line)
		fmt.Printf("score %d, line %s\n", s, line)
		score += s
	}
	return score, nil
}

func getScoreForLine(line string) int {
	expectedClosing := []rune{}

	for _, c := range line {
		meta, ok := characterMeta[c]
		advent.Assertf(ok, fmt.Sprintf("char? %c\n", c))
		// fmt.Printf("char %c\n", c)
		if meta.isOpening {
			comp := meta.complementary
			expectedClosing = append([]rune{comp}, expectedClosing...)
			// fmt.Printf("pushed %s\n", runesAsString(expectedClosing))
		} else {
			expected := expectedClosing[0]
			if expected == c {
				expectedClosing = expectedClosing[1:]
				// fmt.Printf("popped %s\n", runesAsString(expectedClosing))
			} else {
				//fmt.Printf("mismatch expected %c\n", expected)
				//return characterMeta[expected].score
				return meta.score
			}
		}
	}
	return 0 // incomplete, return 0
	//panic(fmt.Sprintf("no score for: %s", line))
}

var characterMeta map[rune]meta

type meta struct {
	complementary rune
	score         int
	isOpening     bool
}

func init() {
	characterMeta = make(map[rune]meta)
	characterMeta['('] = meta{')', 0, true}
	characterMeta['['] = meta{']', 0, true}
	characterMeta['{'] = meta{'}', 0, true}
	characterMeta['<'] = meta{'>', 0, true}
	characterMeta[')'] = meta{'(', 3, false}
	characterMeta[']'] = meta{'[', 57, false}
	characterMeta['}'] = meta{'{', 1197, false}
	characterMeta['>'] = meta{'<', 25137, false}
}

func runesAsString(runes []rune) string {
	s := ""
	for _, r := range runes {
		s += fmt.Sprintf("%c", r)
	}
	return s
}
