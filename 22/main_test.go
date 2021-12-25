package main

import (
	"advent"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc0(t *testing.T) {
	input := `
on x=1..10,y=1..10,z=1..10
on x=3..10,y=3..10,z=3..10
off x=1..5,y=1..5,z=1..5
off x=1..5,y=1..5,z=1..5
`
	v, err := Calc(strings.NewReader(strings.Trim(input, "\n")), false)
	assert.NoError(t, err)
	assert.Equal(t, 0, v)
}

func TestCalc1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f, true)
	assert.NoError(t, err)
	assert.Equal(t, 590784, v)
}

func TestCalc1b(t *testing.T) {
	f, err := os.Open("input1b")
	assert.NoError(t, err)
	v, err := Calc(f, false)
	assert.NoError(t, err)
	assert.Equal(t, 2758514936282235, v)
}
func TestCalc2(t *testing.T) {
	advent.PrintEnabled = false
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f, false)
	assert.NoError(t, err)
	assert.Equal(t, 543306, v)
}
