package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input = `
forward 5
down 5
forward 8
up 3
down 8
forward 2
`

func TestCalcPos(t *testing.T) {
	p, err := CalculatePosition(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, 150, p)
}

func TestCalcPos2(t *testing.T) {
	f, err := os.Open("input")
	assert.NoError(t, err)
	p, err := CalculatePosition(f)
	assert.NoError(t, err)
	assert.Equal(t, 1250395, p)
}

func TestCalcAim(t *testing.T) {
	p, err := CalculateAim(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, 900, p)
}

func TestCalcAim2(t *testing.T) {
	f, err := os.Open("input")
	assert.NoError(t, err)
	p, err := CalculateAim(f)
	assert.NoError(t, err)
	assert.Equal(t, 1451210346, p)
}
