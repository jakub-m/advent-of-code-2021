package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input1 = `start-A
start-b
A-c
A-b
b-d
A-end
b-end`

func TestCalc1(t *testing.T) {
	v, err := Calc(strings.NewReader(input1))
	assert.NoError(t, err)
	assert.Equal(t, 36, v)
}

var input2 = `dc-end
HN-start
start-kj
dc-start
dc-HN
LN-dc
HN-end
kj-sa
kj-HN
kj-dc`

func TestCalc2(t *testing.T) {
	v, err := Calc(strings.NewReader(input2))
	assert.NoError(t, err)
	assert.Equal(t, 103, v)
}

var input3 = `fs-end
he-DX
fs-he
start-DX
pj-DX
end-zg
zg-sl
zg-pj
pj-he
RW-he
fs-DX
pj-RW
zg-RW
start-pj
he-WI
zg-he
pj-fs
start-RW`

func TestCalc3(t *testing.T) {
	v, err := Calc(strings.NewReader(input3))
	assert.NoError(t, err)
	assert.Equal(t, 3509, v)
}

var input4 = `vp-BY
ui-oo
kk-IY
ij-vp
oo-start
SP-ij
kg-uj
ij-UH
SP-end
oo-IY
SP-kk
SP-vp
ui-ij
UH-ui
ij-IY
start-ui
IY-ui
uj-ui
kk-oo
IY-start
end-vp
uj-UH
ij-kk
UH-end
UH-kk`

func TestCalc4(t *testing.T) {
	v, err := Calc(strings.NewReader(input4))
	assert.NoError(t, err)
	assert.Equal(t, 143562, v)
}
