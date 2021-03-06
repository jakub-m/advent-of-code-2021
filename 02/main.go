package main

import (
	"advent"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Command int

const (
	Forward Command = iota
	Up
	Down
)

type Instruction struct {
	Command Command
	Value   int
}

func CalculatePosition(r io.Reader) (int, error) {
	commands, err := parseLinesToCommands(r)
	if err != nil {
		return 0, err
	}

	horizontal, depth := 0, 0

	for _, i := range commands {
		c := i.(Instruction)
		switch c.Command {
		case Forward:
			horizontal += c.Value
		case Up:
			depth -= c.Value
		case Down:
			depth += c.Value
		default:
			panic(c)
		}
	}
	return horizontal * depth, nil
}

func CalculateAim(r io.Reader) (int, error) {
	commands, err := parseLinesToCommands(r)
	if err != nil {
		return 0, err
	}

	horizontal, depth, aim := 0, 0, 0

	for _, i := range commands {
		c := i.(Instruction)
		switch c.Command {
		case Down:
			aim += c.Value
		case Up:
			aim -= c.Value
		case Forward:
			horizontal += c.Value
			depth += aim * c.Value
		default:
			panic(c)
		}
	}
	return horizontal * depth, nil
}

func parseLinesToCommands(r io.Reader) ([]interface{}, error) {
	parseLine := func(line string) (interface{}, error) {
		sp := strings.Split(line, " ")
		if len(sp) != 2 {
			return nil, fmt.Errorf(line)
		}
		c := Instruction{}
		switch sp[0] {
		case "forward":
			c.Command = Forward
		case "up":
			c.Command = Up
		case "down":
			c.Command = Down
		default:
			return nil, fmt.Errorf(line)
		}
		val, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, err
		}
		c.Value = val
		return c, nil
	}

	return advent.ParseLinesReader(r, parseLine)
}
