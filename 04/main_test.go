package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	r, err := os.Open("input1")
	assert.NoError(t, err)
	score := GetBingoScore(r)
	assert.Equal(t, 4512, score)
}
