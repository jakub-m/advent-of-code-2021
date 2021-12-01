package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input = `
199
200
208
210
200
207
240
269
260
263
`

func TestCountIncreases(t *testing.T) {
	lines := strings.Trim(input, "\n")
	r := strings.NewReader(lines)
	c, err := CountIncreasesReader(r)
	assert.NoError(t, err)
	assert.Equal(t, 7, c)
}

func TestCoundSlidingWindoes(t *testing.T) {
	lines := strings.Trim(input, "\n")
	r := strings.NewReader(lines)
	c, err := CountIncreasesWindowedReader(r)
	assert.NoError(t, err)
	assert.Equal(t, 5, c)
}
