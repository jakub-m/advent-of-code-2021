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
	n, rest, err := parseRec(str, 0, nil)
	if err != nil {
		return nil, err
	}
	if rest != "" {
		return nil, fmt.Errorf("incomplete parsing: %s", rest)
	}
	return n, nil
}

func parseRec(input string, level int, curr node) (node, string, error) {
	ch, rest := firstRest(input)
	if ch == "[" {
		bin := binaryNode{level: level + 1, parent: curr}
		nodeLeft, rest, err := parseRec(rest, level+1, &bin)
		if err != nil {
			return nil, "", err
		}
		ch, rest := firstRest(rest)
		if ch != "," {
			return nil, "", fmt.Errorf("expected , got %s", ch)
		}
		nodeRight, rest, err := parseRec(rest, level+1, &bin)
		if err != nil {
			return nil, "", err
		}
		ch, rest = firstRest(rest)
		if ch != "]" {
			return nil, "", fmt.Errorf("expected ] got %s", ch)
		}
		bin.left = nodeLeft
		bin.right = nodeRight
		// n := binaryNode{left: nodeLeft, right: nodeRight, level: level + 1}
		// return &n, rest, nil
		return &bin, rest, nil
	} else {
		val, err := strconv.Atoi(ch)
		if err != nil {
			return nil, "", err
		}
		n := valueNode{val: val, level: level + 1, parent: curr}
		return &n, rest, nil
	}
}

func reduce(root node) node {
	var rec func(node)
	rec = func(n node) {
		fmt.Println(n.getLevel(), n, n.getParent())
		switch b := n.(type) {
		case *binaryNode:
			rec(b.left)
			rec(b.right)
		case *valueNode:

		// if n.getLevel() == 4 {
		// 	fmt.Print(n)
		default:
			panic(fmt.Sprintf("%T %s", n, n))
		}
	}
	rec(root)
	return root
}

type node interface {
	magnitude() int
	String() string
	getLevel() int
	getParent() node
}

type valueNode struct {
	val    int
	level  int
	parent node
}

func (n valueNode) String() string {
	return fmt.Sprint(n.val)
}

func (n valueNode) magnitude() int {
	return n.val
}

func (n valueNode) getLevel() int {
	return n.level
}

func (n valueNode) getParent() node {
	return n.parent
}

type binaryNode struct {
	left   node
	right  node
	level  int
	parent node
}

func (n binaryNode) String() string {
	return fmt.Sprintf("[%s,%s]", n.left, n.right)
}

func (n binaryNode) magnitude() int {
	return 3*n.left.magnitude() + 2*n.right.magnitude()
}

func (n binaryNode) getLevel() int {
	return n.level
}

func (n binaryNode) getParent() node {
	return n.parent
}

func firstRest(s string) (string, string) {
	return s[:1], s[1:]
}
