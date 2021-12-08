package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc1(t *testing.T) {
	tcs := []struct {
		input string
		score int
	}{
		{"be cfbegad cbdgef fgaecd cgeb fdcge agebfd fecdb fabcd edb | fdgacbe cefdb cefbgd gcbe", 8394},
		{"edbfga begcd cbg gc gcadebf fbgde acbgfd abcde gfcbed gfec | fcgedb cgb dgebacf gc", 9781},
		{"fgaebd cg bdaec gdafb agbcfd gdcbef bgcad gfac gcb cdgabef | cg cg fdcagb cbg", 1197},
		{"fbegcd cbd adcefb dageb afcb bc aefdc ecdab fgdeca fcdbega | efabcd cedba gadfec cb", 9361},
		{"aecbfdg fbg gf bafeg dbefa fcge gcbea fcaegb dgceab fcbdga | gecf egdcabf bgf bfgea", 4873},
		{"fgeab ca afcebg bdacfeg cfaedg gcfdb baec bfadeg bafgc acf | gebdcfa ecba ca fadegcb", 8418},
		{"dbcfg fgd bdegcaf fgec aegbdf ecdfab fbedc dacgb gdcebf gf | cefg dcbef fcge gbcadfe", 4548},
		{"bdfegc cbegaf gecbf dfcage bdacg ed bedf ced adcbefg gebcd | ed bcgafe cdgba cbgef", 1625},
		{"egadfb cdbfeg cegd fecab cgb gbdefca cg fgcdab egfdb bfceg | gbdfcae bgc cg cgb", 8717},
		{"gcafb gcf dcaebfg ecagb gf abcdeg gaef cafbge fdbac fegbdc | fgae cfgab fg bagce", 4315},
	}

	for _, tc := range tcs {
		s := strings.NewReader(tc.input)
		v, err := Calc(s)
		assert.NoError(t, err)
		assert.Equal(t, tc.score, v, "line: %s", tc.input)
	}
}

func TestCalcInput1(t *testing.T) {
	f, err := os.Open("input1")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 61229, v)
}

func TestCalcInput2(t *testing.T) {
	f, err := os.Open("input2")
	assert.NoError(t, err)
	v, err := Calc(f)
	assert.NoError(t, err)
	assert.Equal(t, 1031553, v)
}

func TestCalc3(t *testing.T) {
	input := "acedgfb cdfbe gcdfa fbcad dab cefabd cdfgeb eafb cagedb ab | cdfeb fcadb cdfeb cdbaf"
	v, err := Calc(strings.NewReader(input))
	assert.NoError(t, err)
	assert.Equal(t, 5353, v)
}
