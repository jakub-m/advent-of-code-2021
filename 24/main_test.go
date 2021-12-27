package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	in := `
inp w
add z w
mod z 2
div w 2
add y w
mod y 2
div w 2
add x w
mod x 2
div w 2
mod w 2`
	s, err := calcStateFromInput(strings.NewReader(strings.Trim(in, "\n")), []int{9})
	assert.NoError(t, err)
	assert.Equal(t, [4]int{1, 0, 0, 1}, s.reg)

	s, err = calcStateFromInput(strings.NewReader(strings.Trim(in, "\n")), []int{0})
	assert.NoError(t, err)
	assert.Equal(t, [4]int{0, 0, 0, 0}, s.reg)

	s, err = calcStateFromInput(strings.NewReader(strings.Trim(in, "\n")), []int{1})
	assert.NoError(t, err)
	assert.Equal(t, [4]int{0, 0, 0, 1}, s.reg)

	s, err = calcStateFromInput(strings.NewReader(strings.Trim(in, "\n")), []int{8})
	assert.NoError(t, err)
	assert.Equal(t, [4]int{1, 0, 0, 0}, s.reg)
}

func TestCalc1a(t *testing.T) {
	in := `
inp z
inp x
mul z 3
eql z x
`

	s, err := calcStateFromInput(strings.NewReader(strings.Trim(in, "\n")), []int{6, 2})
	assert.NoError(t, err)
	assert.Equal(t, 0, s.reg[operRegZ])

	s, err = calcStateFromInput(strings.NewReader(strings.Trim(in, "\n")), []int{2, 6})
	assert.NoError(t, err)
	assert.Equal(t, 1, s.reg[operRegZ])
}

// func TestCalc2(t *testing.T) {
// 	f, err := os.Open("input2")
// 	assert.NoError(t, err)
// 	c := 1
// 	v, err := calcStateFromInput(f, []int{c, c, c, c, c, c, c, c, c, c, c, c, c, c})
// 	//v, s, err := Calc(f, []int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9})
// 	isValid := v == 0
// 	assert.NoError(t, err)
// 	fmt.Println(s)
// 	assert.Equal(t, true, isValid)
// }

func TestIntConf(t *testing.T) {
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, intToDigitsBase8plus1(startInt))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1}, intToDigitsBase8plus1(9))
	assert.Equal(t, []int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9}, intToDigitsBase8plus1(endInt))
}

func calcStateFromInput(r io.Reader, input []int) (state, error) {
	instructions, err := readInstructions(r)
	if err != nil {
		return state{}, err
	}
	s := state{input: input}
	s = applyInstructions(s, instructions)
	return s, nil
}
