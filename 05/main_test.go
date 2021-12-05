package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	r, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := CalcVents(r)
	assert.NoError(t, err)
	assert.Equal(t, 12, v)
}

func Test2(t *testing.T) {
	r, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := CalcVents(r)
	assert.NoError(t, err)
	assert.Equal(t, 19172, v)
}
