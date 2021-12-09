package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input1 = `2199943210
3987894921
9856789892
8767896789
9899965678`

func TestCalc1(t *testing.T) {
	v, err := Calc(strings.NewReader(input1))
	assert.NoError(t, err)
	assert.Equal(t, 3*9*14*9, v)
}

func TestCalc2(t *testing.T) {
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 566, v)
}
