package main

import (
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
	assert.Equal(t, 15, v)
}
