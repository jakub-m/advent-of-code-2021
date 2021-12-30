package main

import (
	"advent"
	"fmt"
	"io"
)

func Calc(r io.Reader) (int, error) {
	heard := loadHeard(r)

	var i int
	for {
		advent.Println("i", i)
		advent.Println(heard)

		c := moveHeard(&heard)
		i++
		if c == 0 {
			break
		}
		advent.Println("count", c)
	}
	return i, nil
}

func loadHeard(r io.Reader) heard {
	h := heard{}
	h.cucumbers = make(map[pos]direction)

	lines, err := advent.ReadLinesTrim(r)
	advent.PanicErr(err)
	maxX, maxY := 0, 0
	for y, line := range lines {
		maxY = y
		for x, ch := range line {
			maxX = x
			var dir direction
			switch ch {
			case '>':
				dir = east
			case 'v':
				dir = south
			case '.':
				continue
			default:
				panic(fmt.Sprintf("char? %c", ch))
			}
			h.cucumbers[pos{x, y}] = dir
		}
	}
	h.height = maxY + 1
	h.width = maxX + 1
	return h
}

type heard struct {
	width, height int
	cucumbers     map[pos]direction
}

func (h heard) String() string {
	s := ""
	for y := 0; y < h.height; y++ {
		for x := 0; x < h.width; x++ {
			if d, ok := h.cucumbers[pos{x, y}]; ok {
				s += d.String()
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

type direction uint8

const (
	unknownDirection direction = iota
	south
	east
)

func (d direction) String() string {
	switch d {
	case south:
		return "v"
	case east:
		return ">"
	default:
		return "?"
	}
}

type pos struct {
	x, y int
}

func moveHeard(h *heard) int {
	count := 0
	originalMap := h.cucumbers

	//east
	mapAfterMoveEast := make(map[pos]direction)
	for p, d := range originalMap {
		if d != east {
			mapAfterMoveEast[p] = d
			continue
		}
		pEast := h.east(p)
		if _, ok := originalMap[pEast]; !ok {
			mapAfterMoveEast[pEast] = d
			count++
		} else {
			mapAfterMoveEast[p] = d
		}
	}

	//south
	mapAfterMoveSouth := make(map[pos]direction)
	for p, d := range mapAfterMoveEast {
		if d != south {
			mapAfterMoveSouth[p] = d
			continue
		}
		pSouth := h.south(p)
		if _, ok := mapAfterMoveEast[pSouth]; !ok {
			mapAfterMoveSouth[pSouth] = d
			count++
		} else {
			mapAfterMoveSouth[p] = d
		}
	}

	h.cucumbers = mapAfterMoveSouth
	return count
}

func (h heard) east(p pos) pos {
	p.x = (p.x + 1) % h.width
	return p
}

func (h heard) south(p pos) pos {
	p.y = (p.y + 1) % h.height
	return p
}
