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

	scores := []int{}
	for _, line := range lines {
		s, missingClosure := getScoreIfLineCorrupted(line)
		if s != 0 {
			continue
		}
		s2 := scoreForMissingClosure(missingClosure)
		scores = append(scores, s2)
		fmt.Printf("score %d, line %s + %s\n", s2, line, runesAsString(missingClosure))
	}
	score := advent.MedianInt(scores)
	return score, nil
}

func getScoreIfLineCorrupted(line string) (int, []rune) {
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
				return meta.score, []rune{}
			}
		}
	}
	return 0, expectedClosing
}

func scoreForMissingClosure(chars []rune) int {
	score := 0
	for _, c := range chars {
		comp := characterMeta[c].complementary
		sc := characterMeta[comp].score
		score = score*5 + sc
		fmt.Printf("%s %c %d %d\n", runesAsString(chars), c, sc, score)
	}
	return score
}

var characterMeta map[rune]meta

type meta struct {
	complementary rune
	score         int
	isOpening     bool
}

func init() {
	characterMeta = make(map[rune]meta)
	characterMeta['('] = meta{')', 1, true}
	characterMeta['['] = meta{']', 2, true}
	characterMeta['{'] = meta{'}', 3, true}
	characterMeta['<'] = meta{'>', 4, true}
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
