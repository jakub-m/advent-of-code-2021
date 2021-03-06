package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 315, v)
}

func TestCalc0(t *testing.T) {
	t.Skip()
	f, err := os.Open("input0")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 8-1, v)
}

func TestCalc2(t *testing.T) {
	// advent.PrintEnabled = false
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 3012, v)
}
