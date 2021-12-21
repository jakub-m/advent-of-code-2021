package main

import (
	"advent"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	advent.PrintEnabled = false
	v, err := Calc(4, 8)
	assert.NoError(t, err)
	assert.Equal(t, 444356092776315, v)
}

func TestCalc2(t *testing.T) {
	v, err := Calc(4, 3)
	assert.NoError(t, err)
	assert.Equal(t, 734820, v)
}
