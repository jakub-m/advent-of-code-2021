package main

import (
	"advent"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type line struct {
	start point
	end   point
}

type point struct {
	x, y int
}

type grid map[point]int

func CalcVents(r io.Reader) (int, error) {
	textLines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, err
	}

	grid := grid{}
	lines, err := parseLines(textLines)
	if err != nil {
		return 0, err
	}
	for _, line := range lines {
		// if !(line.isHorizontal() || line.isVertical()) {
		// 	continue
		// }
		// fmt.Println(line)
		line.iterPoints(func(p point) {
			// fmt.Println(p)
			grid[p] += 1
		})
	}

	return grid.countOverlap(), nil
}

func parseLines(textLines []string) ([]line, error) {
	lines := []line{}
	re := regexp.MustCompile(`(\d+),(\d+) -> (\d+),(\d+)`)
	for _, t := range textLines {
		m := re.FindStringSubmatch(t)
		if len(m) != 5 {
			return nil, fmt.Errorf("regex: %s", t)
		}
		line := line{
			start: point{
				x: strToInt(m[1]),
				y: strToInt(m[2]),
			},
			end: point{
				x: strToInt(m[3]),
				y: strToInt(m[4]),
			},
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func strToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int(i)
}

func (l line) isHorizontal() bool {
	return l.start.y == l.end.y
}

func (l line) isVertical() bool {
	return l.start.x == l.end.x
}

func (l line) iterPoints(fn func(point)) {
	dx := l.end.x - l.start.x
	dy := l.end.y - l.start.y
	dx = sign0(dx)
	dy = sign0(dy)
	dp := point{dx, dy}

	p := l.start
	for {
		fn(p)
		if p == l.end {
			break
		}
		p = p.add(dp)
	}
}

func sign0(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	default:
		return 0
	}
}

func (p point) add(d point) point {
	return point{
		x: p.x + d.x,
		y: p.y + d.y,
	}
}

func (g grid) countOverlap() int {
	overlapCount := 0
	for _, v := range g {
		if v >= 2 {
			overlapCount++
		}
	}
	return overlapCount
}
