package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	v, err := Calc(4, 8)
	assert.NoError(t, err)
	assert.Equal(t, 739785, v)
}

func TestCalc2(t *testing.T) {
	v, err := Calc(4, 3)
	assert.NoError(t, err)
	assert.Equal(t, 734820, v)
}
