package main

import (
	"advent"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc0(t *testing.T) {
	input := `
on x=0..10,y=0..11,z=0..10
on x=0..5,y=0..12,z=0..10
on x=0..2,y=0..13,z=0..10
`
	v, err := Calc(strings.NewReader(strings.Trim(input, "\n")))
	assert.NoError(t, err)
	assert.Equal(t, 0, v)
}

func TestCalc1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 590784, v)
}

func TestCalc2(t *testing.T) {
	advent.PrintEnabled = false
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 543306, v)
}

func TestExplode(t *testing.T) {
	// fmt.Println(len(c))
	// for _, c := range c {
	// 	fmt.Println(c)
	// }
}
