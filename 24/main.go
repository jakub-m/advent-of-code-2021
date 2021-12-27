package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	Calc(strings.NewReader(""))
}

func Calc(r io.Reader) (int, error) {
	s := state{input: []int{-1}}
	i1 := instruction{
		id:  instInp,
		op1: operRegX,
	}
	i2 := instruction{
		id:   instMul,
		op1:  operRegX,
		op2:  operVal,
		val2: -1,
	}
	s = eval(s, i1)
	s = eval(s, i2)
	fmt.Println(s)

	return 0, nil
}

func eval(s state, ins instruction) state {
	if ins.op1 == operVal {
		panic(fmt.Sprintf("op1 cannot be val: %+v", ins))
	}
	val1 := s.reg[ins.op1]

	val2 := 0
	if ins.op2 == operVal {
		val2 = ins.val2
	} else {
		val2 = s.reg[ins.op2]
	}

	return instrTable[ins.id](s, ins.op1, val1, val2)
}

type instruction struct {
	id   instructionId
	op1  operand
	op2  operand
	val2 int
}

type instructionId uint8

const (
	instInp instructionId = iota
	instAdd
	instMul
	instDiv
	instMod
	instEql
	instructionCount
)

func (i instructionId) String() string {
	switch i {
	case instInp:
		return "inp"
	case instAdd:
		return "add"
	case instMul:
		return "mul"
	case instDiv:
		return "div"
	case instMod:
		return "mod"
	case instEql:
		return "eql"
	default:
		return "?"
	}

}

type operand uint8

const (
	operRegW operand = iota
	operRegX
	operRegY
	operRegZ
	regCount
	operVal
)

func (o operand) String() string {
	switch o {
	case operRegW:
		return "w"
	case operRegX:
		return "x"
	case operRegY:
		return "y"
	case operRegZ:
		return "z"
	case operVal:
		return "val"
	default:
		return "?"
	}
}

type instrFunc func(state, operand, int, int) state

type state struct {
	input []int
	reg   [regCount]int
}

func (s state) String() string {
	return fmt.Sprintf("w:%d x:%d y:%d z:%d inp:%v", s.reg[operRegW], s.reg[operRegX], s.reg[operRegY], s.reg[operRegZ], s.input)
}

var instrTable [instructionCount]instrFunc

func init() {
	instrTable[instInp] = instFuncInp
	instrTable[instAdd] = instFuncAdd
	instrTable[instMul] = instFuncMul
	instrTable[instDiv] = instFuncDiv
	instrTable[instMod] = instFuncMod
	instrTable[instEql] = instFuncEql
}

func instFuncInp(s state, dest operand, val1, val2 int) state {
	inp := s.input[0]
	s.input = s.input[1:]
	s.reg[dest] = inp
	return s
}

func instFuncAdd(s state, dest operand, val1, val2 int) state {
	s.reg[dest] = val1 + val2
	return s
}

func instFuncMul(s state, dest operand, val1, val2 int) state {
	s.reg[dest] = val1 * val2
	return s
}

func instFuncDiv(s state, dest operand, val1, val2 int) state {
	if val2 == 0 {
		panic(fmt.Sprintf("div: op2 is zero. %+v", s))
	}
	s.reg[dest] = val1 / val2
	return s
}

func instFuncMod(s state, dest operand, val1, val2 int) state {
	if val1 < 0 {
		panic(fmt.Sprintf("mod: op1 < 0. %+v", s))
	}
	if val2 < 0 {
		panic(fmt.Sprintf("mod: op1 < 0. %+v", s))
	}
	s.reg[dest] = val1 % val2
	return s
}

func instFuncEql(s state, dest operand, val1, val2 int) state {
	if val1 == val2 {
		s.reg[dest] = 1
	} else {
		s.reg[dest] = 0
	}
	return s
}
