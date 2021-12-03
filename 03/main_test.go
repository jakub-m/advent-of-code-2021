package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input = `
00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010
`

func TestCalc2(t *testing.T) {
	t.Skip()
	tcs := []struct {
		input    string
		expected int
	}{
		{
			`
		1
		1
		0
		`, 1,
		},
		{
			`
		01
		01
		10
		`, 1,
		},
		{
			`
		10
		10
		10
		`, 2,
		},
		{
			`
		110
		010
		110
		`, 6,
		}, {
			`
00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010
`, 198,
		},
	}

	for _, tc := range tcs {
		p := Calc(strings.NewReader(tc.input))
		assert.Equal(t, tc.expected, p, "input: \n%s", tc.input)
	}
}

func TestCalc1(t *testing.T) {
	p := Calc(strings.NewReader(input))
	assert.Equal(t, 198, p)
}

func TestCalc(t *testing.T) {
	r, err := os.Open("input")
	assert.NoError(t, err)
	p := Calc(r)
	assert.Equal(t, 4174964, p)

}
