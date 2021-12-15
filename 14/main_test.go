package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1588, v)
}

func TestCalc2(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f, 40)
	assert.NoError(t, err)
	assert.Equal(t, 2188189693529, v)
}

func TestCalc3(t *testing.T) {
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f, 10)
	assert.NoError(t, err)
	assert.Equal(t, 3048, v-2)
}

func TestCalc4(t *testing.T) {
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f, 40)
	assert.NoError(t, err)
	assert.Equal(t, 0, v-2) // idiotic guess.
}
