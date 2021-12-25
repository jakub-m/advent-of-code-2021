package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	v, err := Calc()
	assert.NoError(t, err)
	assert.Equal(t, 12521, v)
}
