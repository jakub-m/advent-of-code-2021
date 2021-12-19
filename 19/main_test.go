package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotations(t *testing.T) {
	t.Skip()
	in := trim(`
-1,-1,1
-2,-2,2
-3,-3,3
-2,-3,1
5,6,-4
8,0,7`)
	scanners, err := readScanners(strings.NewReader(in))
	assert.NoError(t, err)
	fmt.Println(scanners)
	ref := scanners[0]
	for _, r := range rotations {
		transformed := ref.transform(r)
		fmt.Println(transformed)
	}
	assert.Equal(t, 1, 2)
}

func TestSimple(t *testing.T) {
	f, err := os.Open("inputIdent1")
	assert.NoError(t, err)
	v, err := Calc(f, 6)
	assert.NoError(t, err)
	assert.Equal(t, 6, v)
}

func trim(s string) string {
	return strings.Trim(s, "\n")
}

func TestInput1Piecewise(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	scanners, err := readScanners(f)
	assert.NoError(t, err)
	assert.True(t, isAligned(scanners[0], scanners[1]), "0 and 1")
	// assert.True(t, isAligned(scanners[1], scanners[4]), "1 and 4")
	assert.True(t, false)
}

const defaultThreshold = 12

func isAligned(some, other scanner) bool {
	_, err := alignScanner(other, []scanner{some}, defaultThreshold)
	return err == nil
}

func TestInput1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f, 12)
	assert.NoError(t, err)
	assert.Equal(t, 79, v)
}
