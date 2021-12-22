package main

import (
	"advent"
	"fmt"
	"io"
	"regexp"
)

func Calc(r io.Reader) (int, error) {
	commands, err := parseInput(r)
	if err != nil {
		return 0, err
	}

	reactor := make(map[xyz]bool)
	for x := -50; x <= 50; x++ {
		for y := -50; y <= 50; y++ {
			for z := -50; z <= 50; z++ {
				reactor[xyz{x, y, z}] = false
			}
		}
	}

	for _, c := range commands {
		cub := trim(c.cuboid)
		for x := cub.xx.startInc; x <= cub.xx.endInc; x++ {
			for y := cub.yy.startInc; y <= cub.yy.endInc; y++ {
				for z := cub.zz.startInc; z <= cub.zz.endInc; z++ {
					p := xyz{x, y, z}
					if _, ok := reactor[p]; ok {
						reactor[p] = c.onOff
					}
				}
			}
		}
	}

	count := 0
	for _, v := range reactor {
		if v {
			count++
		}
	}

	return count, nil
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
		c.cuboid.xx.startInc = advent.Atoi(m[2])
		c.cuboid.xx.endInc = advent.Atoi(m[3])
		c.cuboid.yy.startInc = advent.Atoi(m[4])
		c.cuboid.yy.endInc = advent.Atoi(m[5])
		c.cuboid.zz.startInc = advent.Atoi(m[6])
		c.cuboid.zz.endInc = advent.Atoi(m[7])
		commands = append(commands, c)
	}
	return commands, nil
}

type command struct {
	onOff  bool
	cuboid cuboid
}

type cuboid struct {
	xx intRange
	yy intRange
	zz intRange
}

type xyz struct {
	x, y, z int
}

type intRange struct {
	startInc, endInc int
}

func trim(c cuboid) cuboid {
	if c.xx.startInc < -50 {
		c.xx.startInc = -50
	}
	if c.yy.startInc < -50 {
		c.yy.startInc = -50
	}
	if c.zz.startInc < -50 {
		c.zz.startInc = -50
	}

	if c.xx.endInc > 50 {
		c.xx.endInc = 50
	}
	if c.yy.endInc > 50 {
		c.yy.endInc = 50
	}
	if c.zz.endInc > 50 {
		c.zz.endInc = 50
	}
}
