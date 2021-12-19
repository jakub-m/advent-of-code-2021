package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotations(t *testing.T) {
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

func trim(s string) string {
	return strings.Trim(s, "\n")
}
