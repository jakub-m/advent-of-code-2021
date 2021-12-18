package main

import (
	"advent"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func CalcSumMagnitude(r io.Reader) (int, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}

	current, err := parse(lines[0])
	if err != nil {
		return 0, fmt.Errorf("line %s :%s", lines[0], err)
	}
	for _, line := range lines[1:] {
		p, err := parse(line)
		if err != nil {
			return 0, fmt.Errorf("line %s :%s", line, err)
		}
		newCurrent := &binaryNode{
			left:  current,
			right: p,
		}
		current = newCurrent
		// fixParentsAndLevels(current)
	}

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
		return &bin, rest, nil
	} else {
		ch, rest := numberAndRest(input)
		val, err := strconv.Atoi(ch)
		if err != nil {
			return nil, "", err
		}
		n := valueNode{val: val, level: level + 1, parent: curr}
		return &n, rest, nil
	}
}

func numberAndRest(in string) (string, string) {
	re := regexp.MustCompile(`(\d+)(.*)`)
	m := re.FindStringSubmatch(in)
	return m[1], m[2]
}

func reduceAll(root node) {
	for {
		if applied := reduce(root); !applied {
			break
		}
	}
}

func reduce(root node) bool {
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
				zero := valueNode{val: 0, level: n.getLevel(), parent: n.getParent()}
				n.getParent().(*binaryNode).replaceWith(n, &zero)
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
			if n.val >= 10 {
				split := &binaryNode{
					level:  n.level,
					parent: n.getParent(),
				}
				split.left = &valueNode{
					val:    (n.val / 2),
					level:  split.level + 1,
					parent: split,
				}
				split.right = &valueNode{
					val:    (n.val/2 + n.val%2),
					level:  split.level + 1,
					parent: split,
				}
				n.getParent().(*binaryNode).replaceWith(n, split)
				return true
			}
			return false
		}
		panic(fmt.Sprintf("invalid type %T %s", n, n))
	}
	return rec(root)
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

func (n *binaryNode) replaceWith(ref node, newNode node) {
	if n.left == ref {
		n.left = newNode
	} else if n.right == ref {
		n.right = newNode
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
