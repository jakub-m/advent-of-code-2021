package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc0(t *testing.T) {
	in := `11111
	19991
	19191
	19991
	11111`
	v, err := Calc(strings.NewReader(in), 2)
	assert.NoError(t, err)
	assert.Equal(t, 9, v)
}

func TestCalc1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f, 100)
	assert.NoError(t, err)
	assert.Equal(t, 1656, v)
}
