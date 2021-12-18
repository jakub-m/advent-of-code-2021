package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMagnitude(t *testing.T) {
	in := "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]"
	n, err := parse(in)
	assert.NoError(t, err)
	assert.Equal(t, in, n.String())
	assert.Equal(t, 3488, n.magnitude())
}

func TestReduce(t *testing.T) {
	tcs := []struct {
		in, out string
	}{
		{
			in:  "[[[[[9,8],1],2],3],4]",
			out: "[[[[0,9],2],3],4]",
		},
		{
			in:  "[7,[6,[5,[4,[3,2]]]]]",
			out: "[7,[6,[5,[7,0]]]]",
		},
		{
			in:  "[[6,[5,[4,[3,2]]]],1]",
			out: "[[6,[5,[7,0]]],3]",
		},
		{
			in:  "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
			out: "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
		},
		{
			in:  "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
			out: "[[3,[2,[8,0]]],[9,[5,[7,0]]]]",
		},
	}
	for _, tc := range tcs {
		n, err := parse(tc.in)
		assert.NoError(t, err)
		reduce(n)
		assert.NoError(t, err)
		assert.Equal(t, tc.out, n.String())
	}
}

func TestSplit(t *testing.T) {
	// after addition: [[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]
	// after explode:  [[[[0,7],4],[7,[[8,4],9]]],[1,1]]
	// after explode:  [[[[0,7],4],[15,[0,13]]],[1,1]]
	// after split:    [[[[0,7],4],[[7,8],[0,13]]],[1,1]]
	// after split:    [[[[0,7],4],[[7,8],[0,[6,7]]]],[1,1]]
	// after explode:  [[[[0,7],4],[[7,8],[6,0]]],[8,1]]

	tcs := []struct {
		in, out string
	}{
		{
			in:  "[15,1]",
			out: "[[7,8],1]",
		},
		{
			in:  "[[[[0,7],4],[15,[0,13]]],[1,1]]",
			out: "[[[[0,7],4],[[7,8],[0,13]]],[1,1]]",
		},
		{
			in:  "[[[[0,7],4],[[7,8],[0,[6,7]]]],[1,1]]",
			out: "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
	}
	for _, tc := range tcs {
		n, err := parse(tc.in)
		assert.NoError(t, err)
		reduce(n)
		assert.NoError(t, err)
		assert.Equal(t, tc.out, n.String())
	}

}

func TestReduceAll(t *testing.T) {
	in := "[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]"
	out := "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]"
	n, err := parse(in)
	assert.NoError(t, err)
	reduceAll(n)
	assert.NoError(t, err)
	assert.Equal(t, out, n.String())
}

func TestInput1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := CalcSum(f)
	assert.NoError(t, err)
	assert.Equal(t, "[[[[1,1],[2,2]],[3,3]],[4,4]]", v.String())
}

func TestInput2(t *testing.T) {
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := CalcSum(f)
	assert.NoError(t, err)
	assert.Equal(t, "[[[[3,0],[5,3]],[4,4]],[5,5]]", v.String())
}

func TestInput3(t *testing.T) {
	f, err := os.Open("input3")
	assert.NoError(t, err)
	v, err := CalcSum(f)
	assert.NoError(t, err)
	assert.Equal(t, "[[[[5,0],[7,4]],[5,5]],[6,6]]", v.String())
}

func TestSum1(t *testing.T) {
	testSum(t,
		"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
		"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
		"[[[[4,0],[5,4]],[[7,7],[6,0]]],[[8,[7,7]],[[7,9],[5,0]]]]",
	)
}

func testSum(t *testing.T, left, right, expected string) {
	v, err := CalcSum(strings.NewReader(fmt.Sprintf("%s\n%s", left, right)))
	assert.NoError(t, err)
	assert.Equal(t, expected, v.String())
}

func TestInput10(t *testing.T) {
	f, err := os.Open("input10")
	assert.NoError(t, err)
	v, err := CalcSum(f)
	assert.NoError(t, err)
	assert.Equal(t, "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]", v.String())
}

// [[1,2],[[3,4],5]] becomes 143.
// [[[[0,7],4],[[7,8],[6,0]]],[8,1]] becomes 1384.
// [[[[1,1],[2,2]],[3,3]],[4,4]] becomes 445.
// [[[[3,0],[5,3]],[4,4]],[5,5]] becomes 791.
// [[[[5,0],[7,4]],[5,5]],[6,6]] becomes 1137.
// [[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]] becomes 3488.
