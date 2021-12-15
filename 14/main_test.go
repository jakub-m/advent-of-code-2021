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
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1588, v)
}
