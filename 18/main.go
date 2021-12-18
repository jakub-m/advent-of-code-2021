package main

import (
	"advent"
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
	// rec return flag if an action was applied.
	advent.Println("reduce", root)
	var rec func(node) bool
	rec = func(n node) bool {
		// fmt.Println(n.getLevel(), n, n.getParent())
		switch n := n.(type) {
		case *binaryNode:
			if _, _, ok := shouldExplode(n.left); ok {
				advent.Println("should explode (left)", n.left)
				explode(n.left)
				return true
			}
			if _, _, ok := shouldExplode(n.right); ok {
				advent.Println("should explode (right)", n.right)
				explode(n.right)
				n.right = &valueNode{val: 0, level: n.level + 1, parent: n}
				return true
			}
			if applied := rec(n.left); applied {
				return applied
			}
			if applied := rec(n.right); applied {
				return applied
			}
			return false
		case *valueNode:
			return false
			// TODO split here
		}
		panic(fmt.Sprintf("%T %s", n, n))
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

func explode(n node) {
	lv, rv, ok := shouldExplode(n)
	if !ok {
		panic(fmt.Sprintf("should not explode %s", n))
	}
	p := n.getParent()

	advent.Println("find up other left branch")
	if u, ok := findUp(p, hasLeftBranchOtherThan(p)); ok {
		r := findRightMostValueNode(u)
		advent.Println("found value on left", r)
		r.val += lv
	}
	advent.Println("find up other right branch")
	if u, ok := findUp(p, hasRightBranchOtherThan(p)); ok {
		l := findLeftMostValue(u)
		advent.Println("found value on right", l)
		l.val += rv
	}
}

func shouldExplode(n node) (int, int, bool) {
	b, ok := n.(*binaryNode)
	if !ok {
		return 0, 0, false
	}
	ln, ok := b.left.(*valueNode)
	if !ok {
		return 0, 0, false
	}
	rn, ok := b.right.(*valueNode)
	if !ok {
		return 0, 0, false
	}
	if b.level >= 4 {
		return ln.val, rn.val, true
	} else {
		return 0, 0, false
	}
}

func findUp(root node, cond func(node) bool) (*binaryNode, bool) {
	advent.Println("findUp", root)
	if root == nil {
		return nil, false
	}
	b, ok := root.(*binaryNode)
	if !ok {
		return nil, false
	}
	if ok := cond(b); ok {
		advent.Println("findUp OK", root)
		return b, true
	} else {
		return findUp(b.parent, cond)
	}
}

func hasLeftBranchOtherThan(ref node) func(node) bool {
	return func(n node) bool {
		b, ok := n.(*binaryNode)
		if !ok {
			return false
		}
		return b.left != ref
	}
}

func hasRightBranchOtherThan(ref node) func(node) bool {
	return func(n node) bool {
		b, ok := n.(*binaryNode)
		if !ok {
			return false
		}
		return b.right != ref
	}
}

func findRightMostValueNode(b *binaryNode) *valueNode {
	switch v := b.right.(type) {
	case *binaryNode:
		return findRightMostValueNode(v)
	case *valueNode:
		return v
	}
	return nil
}

func findLeftMostValue(b *binaryNode) *valueNode {
	switch v := b.left.(type) {
	case *binaryNode:
		return findLeftMostValue(v)
	case *valueNode:
		return v
	}
	return nil
}
