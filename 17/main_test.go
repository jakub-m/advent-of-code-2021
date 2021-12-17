package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	// target area: x=20..30, y=-10..-5
	v, err := Calc(20, 30, -10, -5)
	assert.NoError(t, err)
	assert.Equal(t, 45, v)
}

func TestCalc2(t *testing.T) {
	// target area: x=235..259, y=-118..-62
	v, err := Calc(235, 259, -118, -62)
	assert.NoError(t, err)
	assert.Equal(t, 6903, v)
}
