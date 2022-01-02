package main

import (
	"advent"
	"fmt"
	"io"
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
	f, err := os.Open("24/input2")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	is, err := newInstructionsetReader(f)
	if err != nil {
		panic(err)
	}

	var int1to9 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	_ = int1to9
	// var int1 = []int{5}

	// digitsInput := [][]int{
	// 	// {1},     // 1
	// 	// {1},     // 2
	// 	// {1},     // 3
	// 	// {1},     // 4
	// 	// {1},     // 5
	// 	// {1},     // 6
	// 	// {1},     // 7
	// 	// {1},     // 8
	// 	// {1},     // 9
	// 	// {1},     // 10
	// 	// {1},     // 11
	// 	// {1},     // 12
	// 	// {1},     // 13
	// 	int1to9, // 14
	// }
	// _ = digitsInput

	// inspect

	noInspect := func(ins instruction, result state) {}
	_ = noInspect

	for d := 1; d <= 9; d++ {
		for zIn := 1; zIn < 26; zIn++ {
			result := is.execRegZ([]int{d}, zIn, noInspect)
			zOut := result.reg[operRegZ]
			if intIn(zOut, []int{7, 8, 9, 10, 11, 12, 13, 14, 15}) {
				//fmt.Printf("GOOD z:%d %v\n", zIn, d)
				fmt.Printf("zIn,d,zOut\t%d\t%d\t%d\n", zIn, d, zOut)
			}
		}

	}

	// iterDigits(digitsInput, func(d []int) {
	// 	for zIn := 1; zIn < 26; zIn++ {
	// 		result := is.execRegZ(d, zIn, noInspect)
	// 		zOut := result.reg[operRegZ]
	// 		if zOut == 0 {
	// 			//fmt.Printf("GOOD z:%d %v\n", zIn, d)
	// 			fmt.Printf("zIn,d,zOut\t%d\t%d\t%d\n", zIn, zOut, d)
	// 		}
	// 	}

	// 	//fmt.Println(s, d)
	// })
}

func newInstructionsetReader(r io.Reader) (instructionset, error) {
	instructions, err := readInstructions(r)
	if err != nil {
		return nil, err
	}
	return instructionset(instructions), nil
}

type instructionset []instruction

func (is instructionset) exec(input []int, inspect func(instruction, state)) state {
	inputCopy := make([]int, len(input))
	copy(inputCopy, input)
	s := state{input: inputCopy}
	return applyInstructions(s, is, inspect)
}

func (is instructionset) execRegZ(input []int, regZ int, inspect func(instruction, state)) state {
	inputCopy := make([]int, len(input))
	copy(inputCopy, input)
	s := state{input: inputCopy}
	s.reg[operRegZ] = regZ
	return applyInstructions(s, is, inspect)
}

func Calc(r io.Reader) (int, error) {
	// instructions, err := readInstructions(r)
	// if err != nil {
	// 	return 0, err
	// }

	// for i := endInt; i >= startInt; i-- {
	// 	s := state{input: intToDigitsBase8plus1(i)}
	// 	s = applyInstructions(s, instructions)
	// 	isValid := s.reg[operRegZ] == 0
	// 	if isValid {
	// 		fmt.Printf("%v\tz %d\t%v\tin %d\t%v\n", isValid, s.reg[operRegZ], intToDigitsBase8plus1(s.reg[operRegZ]), i, intToDigitsBase8plus1(i))
	// 	}
	// }

	return 0, nil
}

func applyInstructions(s state, instructions []instruction, inspect func(instruction, state)) state {
	for i := range instructions {
		s = eval(s, instructions[i], inspect)
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
	for i, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		p, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		p.lineNumber = i + 1
		p.line = line
		ii = append(ii, p)
	}
	return ii, nil
}

func parseLine(line string) (instruction, error) {
	in := instruction{}
	parts := strings.Split(trimLine(line), " ")
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

func trimLine(line string) string {
	if i := strings.Index(line, "#"); i >= 0 {
		line = line[0:i]
	}
	return strings.Trim(line, " \n")
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

func eval(s state, ins instruction, inspect func(instruction, state)) state {
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

	result := instrTable[ins.id](s, ins.op1, val1, val2)
	inspect(ins, result)

	return result
}

type instruction struct {
	id   instructionId
	op1  operand
	op2  operand
	val2 int
	// debug
	lineNumber int
	line       string
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
	// debug
	lastInput int
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
	s.lastInput = inp
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

func iterDigits(ranges [][]int, call func([]int)) {
	var rec func(ranges [][]int, call func([]int), soFar []int)

	rec = func(ranges [][]int, call func([]int), soFar []int) {
		if len(ranges) == 0 {
			call(soFar)
			return
		}

		for _, d := range ranges[0] {
			rec(ranges[1:], call, append(soFar, d))
		}
	}

	rec(ranges, call, []int{})
}

func intIn(needle int, haystack []int) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}
	return false
}
