package main

import (
	"fmt"
	"io"
	"strconv"
)

func Calc(r io.Reader) (int, error) {
	return 0, nil
}

func parse(str string) (node, error) {
	n, rest, err := parseRec(str)
	if err != nil {
		return nil, err
	}
	if rest != "" {
		return nil, fmt.Errorf("incomplete parsing: %s", rest)
	}
	return n, nil
}

func parseRec(input string) (node, string, error) {
	ch, rest := firstRest(input)
	if ch == "[" {
		nodeLeft, rest, err := parseRec(rest)
		if err != nil {
			return nil, "", err
		}
		ch, rest := firstRest(rest)
		if ch != "," {
			return nil, "", fmt.Errorf("expected , got %s", ch)
		}
		nodeRight, rest, err := parseRec(rest)
		if err != nil {
			return nil, "", err
		}
		ch, rest = firstRest(rest)
		if ch != "]" {
			return nil, "", fmt.Errorf("expected ] got %s", ch)
		}
		n := binaryNode{nodeLeft, nodeRight}
		return &n, rest, nil
	} else {
		val, err := strconv.Atoi(ch)
		if err != nil {
			return nil, "", err
		}
		n := valueNode{val}
		return &n, rest, nil
	}
}

type node interface {
	magnitude() int
	String() string
}

type valueNode struct {
	val int
}

func (n valueNode) String() string {
	return fmt.Sprint(n.val)
}

func (n valueNode) magnitude() int {
	return n.val
}

type binaryNode struct {
	left  node
	right node
}

func (n binaryNode) String() string {
	return fmt.Sprintf("[%s,%s]", n.left, n.right)
}

func (n binaryNode) magnitude() int {
	return 3*n.left.magnitude() + 2*n.right.magnitude()
}

func firstRest(s string) (string, string) {
	return s[:1], s[1:]
}
