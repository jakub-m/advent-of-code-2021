package main

import (
	"advent"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func main() {
	Calc(strings.NewReader(""), []int{})
}

func Calc(r io.Reader, input []int) (int, state, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return 0, state{}, err
	}

	instructions, err := parseInstructionsFromLines(lines)
	if err != nil {
		return 0, state{}, err
	}

	s := state{input: input}
	for i := range instructions {
		s = eval(s, instructions[i])
	}

	return int(s.reg[operRegZ]), s, nil
}

func parseInstructionsFromLines(lines []string) ([]instruction, error) {
	ii := []instruction{}
	for _, line := range lines {
		p, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		ii = append(ii, p)
	}
	return ii, nil
}

func parseLine(line string) (instruction, error) {
	in := instruction{}
	parts := strings.Split(strings.Trim(line, " \n"), " ")
	var err error
	if parts[0] == "inp" {
		in.id = instInp
		if op, ok := operandFromString(parts[1]); ok {
			in.op1 = op
		} else {
			return in, fmt.Errorf("bad line: %s", line)
		}
	} else {
		switch parts[0] {
		case "mul":
			in.id = instMul
		case "add":
			in.id = instAdd
		case "mod":
			in.id = instMod
		case "eql":
			in.id = instEql
		case "div":
			in.id = instDiv
		default:
			return in, fmt.Errorf("bad line: %s", line)
		}

		in, err = updateInstructionFromTwoStrings(in, parts[1], parts[2])
		if err != nil {
			return in, err
		}
	}

	return in, nil
}

func updateInstructionFromTwoStrings(in instruction, opStr1, opStr2 string) (instruction, error) {
	if op, ok := operandFromString(opStr1); ok {
		in.op1 = op
	} else {
		return in, fmt.Errorf("bad op1 %s", opStr1)
	}
	if op, ok := operandFromString(opStr2); ok {
		in.op2 = op
	} else if val, ok := valueFromString(opStr2); ok {
		in.op2 = operVal
		in.val2 = val
	} else {
		return in, fmt.Errorf("bad op2 %s", opStr2)
	}
	return in, nil
}

func operandFromString(opid string) (operand, bool) {
	switch opid {
	case "w":
		return operRegW, true
	case "x":
		return operRegX, true
	case "y":
		return operRegY, true
	case "z":
		return operRegZ, true
	default:
		return operInvalid, false
	}
}

func valueFromString(s string) (int, bool) {
	v, err := strconv.Atoi(s)
	return v, err == nil
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

	fmt.Printf("%s %s(%d) %s(%d)\n", ins, ins.op1, val1, ins.op2, val2)

	return instrTable[ins.id](s, ins.op1, val1, val2)
}

type instruction struct {
	id   instructionId
	op1  operand
	op2  operand
	val2 int
}

func (i instruction) String() string {
	return fmt.Sprintf("%s %s %s %d", i.id, i.op1, i.op2, i.val2)
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
	operInvalid
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
	fmt.Printf("inp %s (%d)\n", dest, inp)
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
