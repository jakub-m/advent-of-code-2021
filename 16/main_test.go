package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	in := "110100101111111000101000"
	//in := "110 100 10111 11110 00101 000"
	v, err := Calc(in)
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}
