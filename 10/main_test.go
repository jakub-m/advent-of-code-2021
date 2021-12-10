package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	t.Skip()
	tcs := []struct {
		input string
		score int
	}{
		// {
		// 	"{([(<{}[<>[]}>{[]{[(<()>",
		// 	57,
		// },
		// {
		// 	"[({(<(())[]>[[{[]{<()<>>", // incomplete, return 0
		// 	0,
		// },
		{"{([(<{}[<>[]}>{[]{[(<()>", 57},  //  - Expected ], but found } instead.
		{"[[<[([]))<([[{}[[()]]]", 57},    //  - Expected ], but found ) instead.
		{"[{[{({}]{}}([{[{{{}}([]", 3},    //  - Expected ), but found ] instead.
		{"[<(<(<(<{}))><([]([]()", 25137}, //  - Expected >, but found ) instead.
		{"<{([([[(<>()){}]>(<<{{", 57},    //  - Expected ], but found > instead.
	}
	for _, tc := range tcs {
		v, err := Calc(strings.NewReader(tc.input))
		assert.NoError(t, err)
		assert.Equal(t, tc.score, v)
	}
}

func TestCalc2(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 26397, v)
}

func TestCalc3(t *testing.T) {
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 341823, v)
}
