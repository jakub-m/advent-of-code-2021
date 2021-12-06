package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	v, err := Calc(strings.NewReader(""))
	assert.NoError(t, err)
	assert.Equal(t, 0, v)
}
