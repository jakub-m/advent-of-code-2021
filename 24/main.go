package main

import (
	"advent"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	digit1 int = iota
	digit2
	digit3
	digit4
	digit5
	digit6
	digit7
	digit8
	digit9
)

func main() {
	f, err := os.Open("24/input2test")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	instructions, err := readInstructions(f)
	if err != nil {
		panic(err)
	}

	for {
		i := rand.Intn(endInt + 1)
		s := state{input: intToDigitsBase8plus1(i)}
		s = applyInstructions(s, instructions)
		//isValid := s.reg[operRegZ] == 0
		in := intToDigitsBase8plus1(i)
		w := in[digit2]
		x := 1
		y := in[digit2] + 9
		z2 := (in[digit1]+12)*26 + in[digit2] + 9

		w = in[digit3]
		y = in[digit3] + 8
		z3 := z2*26 + in[digit3] + 8

		z := z3
		fmt.Printf("%s\t%v\n", s, in)
		fmt.Printf("w:%d x:%d y:%d z:%d\n", w, x, y, z)
	}

	// v, err := Calc(f)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("RESULT", v)
}

var startInt int
var endInt int

func init() {
	startInt = 0

	e, err := strconv.ParseInt("88888888888888", 9, 64)
	if err != nil {
		panic(err)
	}
	endInt = int(e)
}

func Calc(r io.Reader) (int, error) {
	instructions, err := readInstructions(r)
	if err != nil {
		return 0, err
	}

	for i := endInt; i >= startInt; i-- {
		s := state{input: intToDigitsBase8plus1(i)}
		s = applyInstructions(s, instructions)
		isValid := s.reg[operRegZ] == 0
		if isValid {
			fmt.Printf("%v\tz %d\t%v\tin %d\t%v\n", isValid, s.reg[operRegZ], intToDigitsBase8plus1(s.reg[operRegZ]), i, intToDigitsBase8plus1(i))
		}
	}

	return 0, nil
}

func intToDigitsBase8plus1(in int) []int {
	ints := make([]int, 14)
	s := strconv.FormatInt(int64(in), 9)
	ss := []string{}
	for _, c := range s {
		ss = append(ss, fmt.Sprintf("%c", c))
	}
	for i := 13; i >= 0; i-- {
		if len(ss) > 0 {
			c := ss[len(ss)-1]
			ss = ss[:len(ss)-1]
			v := advent.Atoi(c) + 1
			ints[i] = v
		} else {
			ints[i] = 1
		}
	}

	return ints
}

func applyInstructions(s state, instructions []instruction) state {
	for i := range instructions {
		s = eval(s, instructions[i])
	}
	return s
}

func readInstructions(r io.Reader) ([]instruction, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return nil, err
	}

	instructions, err := parseInstructionsFromLines(lines)
	if err != nil {
		return nil, err
	}
	return instructions, nil
}

func parseInstructionsFromLines(lines []string) ([]instruction, error) {
	ii := []instruction{}
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
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

	// fmt.Printf("%s %s(%d) %s(%d)\n", ins, ins.op1, val1, ins.op2, val2)

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
	//return fmt.Sprintf("w:%d x:%d y:%d z:%d inp:%v", s.reg[operRegW], s.reg[operRegX], s.reg[operRegY], s.reg[operRegZ], s.input)
	return fmt.Sprintf("w:%d x:%d y:%d z:%d", s.reg[operRegW], s.reg[operRegX], s.reg[operRegY], s.reg[operRegZ])
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
	// fmt.Printf("inp %s (%d)\n", dest, inp)
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
