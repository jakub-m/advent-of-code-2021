package main

import (
	"advent"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	r, err := os.Open("input1")
	assert.NoError(t, err)
	advent.ReadLinesTrim(r)
}
