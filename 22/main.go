package main

import (
	"advent"
	"fmt"
	"io"
	"regexp"
	"sort"
)

func Calc(r io.Reader) (int, error) {
	commands, err := parseInput(r)
	if err != nil {
		return 0, err
	}

	rulers := rulers{[]int{}, []int{}, []int{}}
	for _, comm := range commands {
		rulers.x = updateRuler(rulers.x, comm.cuboid.xx)
		rulers.y = updateRuler(rulers.y, comm.cuboid.yy)
		rulers.z = updateRuler(rulers.z, comm.cuboid.zz)
	}
	rulers.x = advent.UniqInt(rulers.x)
	rulers.y = advent.UniqInt(rulers.y)
	rulers.z = advent.UniqInt(rulers.z)
	sort.Ints(rulers.x)
	sort.Ints(rulers.y)
	sort.Ints(rulers.z)

	commands = splitCuboidsInCommands(commands, func(c cuboid) []cuboid { return explodeAlongDims(c, rulers) })

	advent.Println("rulers", rulers)
	for _, c := range commands {
		advent.Println(c)
	}

	return 0, nil
}

func updateRuler(ruler []int, intRange intRange) []int {
	ruler = append(ruler, intRange.start, intRange.end)
	return ruler
}

func parseInput(r io.Reader) ([]command, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(on|off) x=(-?\d+)\.\.(-?\d+),y=(-?\d+)\.\.(-?\d+),z=(-?\d+)\.\.(-?\d+)`)
	commands := []command{}
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("could not parse " + line)
		}
		c := command{}

		if m[1] == "on" {
			c.onOff = true
		} else if m[1] == "off" {
			c.onOff = false
		} else {
			panic("line" + line)
		}
		c.cuboid.xx.start = advent.Atoi(m[2])
		c.cuboid.xx.end = advent.Atoi(m[3]) + 1
		c.cuboid.yy.start = advent.Atoi(m[4])
		c.cuboid.yy.end = advent.Atoi(m[5]) + 1
		c.cuboid.zz.start = advent.Atoi(m[6])
		c.cuboid.zz.end = advent.Atoi(m[7]) + 1
		commands = append(commands, c)
	}
	return commands, nil
}

type command struct {
	onOff  bool
	cuboid cuboid
}

func (c command) String() string {
	strOnOff := "off"
	if c.onOff {
		strOnOff = "on"
	}
	return fmt.Sprintf("%s %s", strOnOff, c.cuboid)
}

type cuboid struct {
	xx intRange
	yy intRange
	zz intRange
}

func (c cuboid) String() string {
	return fmt.Sprintf("%s,%s,%s", c.xx, c.yy, c.zz)
}

type xyz struct {
	x, y, z int
}

type intRange struct {
	start int // inclusive
	end   int // exclusive
}

func (r intRange) String() string {
	return fmt.Sprintf("(%d,%d)", r.start, r.end)
}

func trim(c cuboid) cuboid {
	if c.xx.start < -50 {
		c.xx.start = -50
	}
	if c.yy.start < -50 {
		c.yy.start = -50
	}
	if c.zz.start < -50 {
		c.zz.start = -50
	}

	if c.xx.end > 50 {
		c.xx.end = 51
	}
	if c.yy.end > 50 {
		c.yy.end = 51
	}
	if c.zz.end > 50 {
		c.zz.end = 51
	}
	return c
}

func sortedCoords(rr ...intRange) []int {
	vals := []int{}
	for _, r := range rr {
		vals = append(vals, r.start, r.end)
	}
	sort.Ints(vals)
	return vals
}

type rulers struct {
	x, y, z []int
}

func splitCuboidsInCommands(commands []command, fn func(cuboid) []cuboid) []command {
	newCommands := []command{}
	for i, comm := range commands {
		fmt.Printf("%.1f%%    \r", 100.0*float32(i)/float32(len(commands)))
		for _, cub := range fn(comm.cuboid) {
			newCommands = append(newCommands, command{
				onOff:  comm.onOff,
				cuboid: cub,
			})
		}
	}
	return newCommands
}

const (
	dimensionX = 1
	dimensionY = 2
	dimensionZ = 3
)

func explodeAlongDims(c cuboid, allRulers rulers) []cuboid {
	smallerRulers := findSmallerRulers(allRulers, c)

	explodedAlongX := explodeAlongDim(c, smallerRulers, dimensionX)

	explodedAlongY := []cuboid{}
	for _, cx := range explodedAlongX {
		explodedAlongY = append(explodedAlongY, explodeAlongDim(cx, smallerRulers, dimensionY)...)
	}

	explodedAlongZ := []cuboid{}
	for _, cy := range explodedAlongY {
		explodedAlongZ = append(explodedAlongZ, explodeAlongDim(cy, smallerRulers, dimensionZ)...)
	}
	return explodedAlongZ
}

func findSmallerRulers(allRulers rulers, c cuboid) rulers {
	indexXStart := findIndexOf(allRulers.x, c.xx.start)
	indexXEnd := findIndexOf(allRulers.x, c.xx.end)
	indexYStart := findIndexOf(allRulers.y, c.yy.start)
	indexYEnd := findIndexOf(allRulers.y, c.yy.end)
	indexZStart := findIndexOf(allRulers.z, c.zz.start)
	indexZEnd := findIndexOf(allRulers.z, c.zz.end)

	return rulers{
		x: allRulers.x[indexXStart : indexXEnd+1],
		y: allRulers.y[indexYStart : indexYEnd+1],
		z: allRulers.z[indexZStart : indexZEnd+1],
	}
}

func findIndexOf(sortedHaystack []int, needle int) int {
	n := len(sortedHaystack)
	k := sort.Search(n, func(i int) bool { return sortedHaystack[i] >= needle })
	if k == n {
		panic(fmt.Sprintf("failed to find %d in %v", needle, sortedHaystack))
	}
	return k
}

func explodeAlongDim(c cuboid, rulers rulers, dim int) []cuboid {
	var ticks []int
	var cuboidDim intRange
	switch dim {
	case dimensionX:
		ticks = rulers.x
		cuboidDim = c.xx
	case dimensionY:
		ticks = rulers.y
		cuboidDim = c.yy
	case dimensionZ:
		ticks = rulers.z
		cuboidDim = c.zz
	default:
		panic(fmt.Sprint("bad dim:", dim))
	}
	exploded := []cuboid{}
	for i := 0; i < len(ticks)-1; i++ {
		r0, r1 := ticks[i], ticks[i+1]
		// fmt.Println("r0", r0, "r1", r1, "cuboidDim", cuboidDim)
		if r0 >= cuboidDim.start && r1 <= cuboidDim.end {
			smallerCuboid := c

			switch dim {
			case dimensionX:
				smallerCuboid.xx.start = r0
				smallerCuboid.xx.end = r1
			case dimensionY:
				smallerCuboid.yy.start = r0
				smallerCuboid.yy.end = r1
			case dimensionZ:
				smallerCuboid.zz.start = r0
				smallerCuboid.zz.end = r1
			default:
				panic(fmt.Sprint("bad dim:", dim))
			}

			exploded = append(exploded, smallerCuboid)
		}
	}
	if len(exploded) < 1 {
		panic(fmt.Sprintf("bug. zero exploded cuboids. cuboid %s, dim %d, cuboidDim: %v, ticks: %v", c, dim, cuboidDim, ticks))
	}
	return exploded
}
