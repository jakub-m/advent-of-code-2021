package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMagnitude(t *testing.T) {
	// v, err := Calc(strings.NewReader(""))
	in := "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]"
	n, err := parse(in)
	assert.NoError(t, err)
	assert.Equal(t, in, n.String())
	assert.Equal(t, 3488, n.magnitude())
}

func TestReduce(t *testing.T) {
	in := "[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]"
	n, err := parse(in)
	assert.NoError(t, err)
	n = reduce(n)
	assert.NoError(t, err)
	assert.Equal(t, "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]", n.String())
}

// [[1,2],[[3,4],5]] becomes 143.
// [[[[0,7],4],[[7,8],[6,0]]],[8,1]] becomes 1384.
// [[[[1,1],[2,2]],[3,3]],[4,4]] becomes 445.
// [[[[3,0],[5,3]],[4,4]],[5,5]] becomes 791.
// [[[[5,0],[7,4]],[5,5]],[6,6]] becomes 1137.
// [[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]] becomes 3488.
