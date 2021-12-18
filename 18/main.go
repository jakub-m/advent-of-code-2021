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
			if ok := shouldExplode(n); ok {
				leftValue := n.left.(*valueNode)
				if nodeBefore := findFirstValueBefore(leftValue, root); nodeBefore != nil {
					advent.Println("node before", nodeBefore)
					nodeBefore.val += leftValue.val
				}
				rightValue := n.right.(*valueNode)
				if nodeAfter := findFirstValueAfter(rightValue, root); nodeAfter != nil {
					nodeAfter.val += rightValue.val
				}
				n.getParent().(*binaryNode).replaceWithZero(n)
				return true
			}
			// if ok := shouldExplode(n.left); ok {
			// 	advent.Println("should explode (left)", n.left)
			// 	explode(n.left)
			// 	return true
			// }
			// if ok := shouldExplode(n.right); ok {
			// 	advent.Println("should explode (right)", n.right)
			// 	explode(n.right)
			// 	n.right = &valueNode{val: 0, level: n.level + 1, parent: n}
			// 	return true
			// }
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

func (n *binaryNode) replaceWithZero(ref node) {
	zero := valueNode{val: 0, level: ref.getLevel(), parent: n}
	if n.left == ref {
		n.left = &zero
	} else if n.right == ref {
		n.right = &zero
	} else {
		panic("incorrect ref, neither left nor right")
	}

}

func firstRest(s string) (string, string) {
	return s[:1], s[1:]
}

func shouldExplode(n node) bool {
	b, ok := n.(*binaryNode)
	if !ok {
		return false
	}
	_, ok = b.left.(*valueNode)
	if !ok {
		return false
	}
	_, ok = b.right.(*valueNode)
	if !ok {
		return false
	}
	return b.level > 4
}

func findFirstValueBefore(ref *valueNode, root node) *valueNode {
	values := flattenValues(root)
	for i, v := range values {
		if v == ref && i > 0 {
			return values[i-1]
		}
	}
	return nil
}

func findFirstValueAfter(ref *valueNode, root node) *valueNode {
	values := flattenValues(root)
	for i, v := range values {
		if v == ref && i < (len(values)-1) {
			return values[i+1]
		}
	}
	return nil
}

func flattenValues(root node) []*valueNode {
	values := []*valueNode{}
	var rec func(n node)
	rec = func(n node) {
		switch n := n.(type) {
		case *binaryNode:
			rec(n.left)
			rec(n.right)
		case *valueNode:
			values = append(values, n)
		}
	}
	rec(root)
	return values
}
