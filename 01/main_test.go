package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	lines := `
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
	lines = strings.Trim(lines, "\n")
	r := strings.NewReader(lines)
	c, err := CountIncreases(r)
	assert.NoError(t, err)
	assert.Equal(t, 7, c)
}
