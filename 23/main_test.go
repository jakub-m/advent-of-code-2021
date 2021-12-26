package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	v, err := Calc(initialSituation1())
	assert.NoError(t, err)
	assert.Equal(t, 44169, v)
}

func TestCalc2(t *testing.T) {
	v, err := Calc(initialSituation2())
	assert.NoError(t, err)
	assert.Equal(t, 17400, v)
}
