package main

import (
	"advent"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func noInspect(ins instruction, result state) {}

type product struct {
	Section, ZIn, Digit, ZOut int
}

func (p product) String() string {
	return fmt.Sprintf("%d\t%d\t%d\t%d", p.Section, p.ZIn, p.Digit, p.ZOut)
}

func main() {
	programFileName := "24/input2"

	if false {
		f, err := os.Open(programFileName)
		advent.PanicErr(err)
		ins, err := newInstructionsetReader(f)
		advent.PanicErr(err)
		//s := ins.exec([]int{3, 9, 9, 9, 9, 6, 9, 8, 7, 9, 9, 4, 2, 9}, noInspect)
		s := ins.exec([]int{1, 8, 1, 1, 6, 1, 2, 1, 1, 3, 4, 1, 1, 7}, noInspect)
		fmt.Println(s)
		if s.reg[operRegZ] == 0 {
			fmt.Println("GOOD")
		}
	}

	var products []product
	_ = products

	if false {
		//sectionLimit, goal := 6, 9921230
		//sectionLimit, goal := 7, 209893599
		sectionLimit, goal := 8, 159169942
		//programFileName = "24/input3"
		products = getProductsBySection(programFileName, sectionLimit)
		fmt.Println("len(products) =", len(products))
		fmt.Println("last product", products[len(products)-1])
		products = pruneUseless(products, goal)
		fmt.Println("pruned len(products) =", len(products))

		// storeProductsAsGob(products, "products.gob")
		path := findPathLeadingToNumber(products, goal)
		fmt.Println(path)

	}

	if true {
		products = getProductsBySection(programFileName, 0)
		//products := loadProductsFromGob("products.gob")
		fmt.Println("loaded", len(products))
		products = pruneUseless(products, 0)
		fmt.Println("pruned len(products) =", len(products))
		path := findPathLeadingToNumber(products, 0)
		fmt.Println(path)
	}
}

func storeProductsAsGob(p []product, fileName string) {
	fmt.Println("storing", len(p), "to", fileName)
	gobOut, err := os.Create(fileName)
	advent.PanicErr(err)
	defer gobOut.Close()
	encoder := gob.NewEncoder(gobOut)
	if err := encoder.Encode(p); err != nil {
		panic(err)
	}
	fmt.Println("done storing")
}

func loadProductsFromGob(fileName string) []product {
	fmt.Println("loading from", fileName)
	gobIn, err := os.Open(fileName)
	advent.PanicErr(err)
	defer gobIn.Close()
	decoder := gob.NewDecoder(gobIn)
	p := []product{}
	err = decoder.Decode(&p)
	advent.PanicErr(err)
	fmt.Println("loaded")
	return p
}

func getProductsBySection(input string, limit int) []product {
	f, err := os.Open(input)
	advent.PanicErr(err)
	defer f.Close()

	allInstructions, err := newInstructionsetReader(f)
	advent.PanicErr(err)

	sections := splitInstructionsByInp(allInstructions)

	if limit > 0 {
		sections = sections[:limit] // TODO remove
	}
	products := []product{}

	sectionZIn := make(map[int]bool)
	sectionZIn[0] = true
	for iSection := range sections {
		sectionZOut := make(map[int]bool)
		for digit := 1; digit <= 9; digit++ {
			for zIn := range sectionZIn {
				result := sections[iSection].execRegZ([]int{digit}, zIn, noInspect)
				zOut := result.reg[operRegZ]
				// fmt.Println("zIn", zIn, "digit", digit, "zOut", zOut)
				// if !sectionZOut[zOut] || true {
				products = append(products, product{iSection, zIn, digit, zOut})
				// }
				sectionZOut[zOut] = true
				//fmt.Printf("%d\t%d\t%d\t%d\n", iSection, zIn, digit, zOut)
			}
		}
		sectionZIn = sectionZOut
		fmt.Println(iSection, len(sectionZOut))
	}

	return products
}

func newInstructionsetReader(r io.Reader) (instructionset, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return nil, err
	}

	instructions, err := parseInstructionsFromLines(lines)
	if err != nil {
		return nil, err
	}
	return instructionset(instructions), nil
}

func splitInstructionsByInp(is instructionset) []instructionset {
	split := []instructionset{}
	current := instructionset{}
	update := func() {
		if len(current) > 0 {
			split = append(split, current)
		}
		current = instructionset{}
	}
	for _, ins := range is {
		if ins.id == instInp {
			update()
		}
		current = append(current, ins)
	}
	update()
	return split
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
	if i.op2 == operVal {
		return fmt.Sprintf("%s %s %d", i.id, i.op1, i.val2)
	} else {
		return fmt.Sprintf("%s %s %s", i.id, i.op1, i.op2)
	}
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

func findPathLeadingToNumber(products []product, goal int) []int {
	fmt.Println("find path for", len(products), "products")

	mapSectionZinDigitZout := make(map[int]map[int]map[int]int)

	for _, p := range products {
		if _, ok := mapSectionZinDigitZout[p.Section]; !ok {
			mapSectionZinDigitZout[p.Section] = make(map[int]map[int]int)
		}
		mapZinDigitZout := mapSectionZinDigitZout[p.Section]
		if _, ok := mapZinDigitZout[p.ZIn]; !ok {
			mapZinDigitZout[p.ZIn] = make(map[int]int)
		}
		mapDigitZout := mapZinDigitZout[p.ZIn]
		mapDigitZout[p.Digit] = p.ZOut
	}

	// TODO optimize - prune dead ends
	// TODO optimize - drop main table after reindexing

	fmt.Println("start rec")

	// type cacheKey struct {
	// 	section, zIn, targetZOut int
	// }
	// cache := make(map[cacheKey][]int)
	// var cacheHit, cacheMiss *int = new(int), new(int)
	// *cacheHit, *cacheMiss = 0, 0

	var rec func(section int, zIn int, targetZout int) []int
	rec = func(section int, zIn int, targetZOut int) (result []int) {
		// ck := cacheKey{section, zIn, targetZOut}
		// if v, ok := cache[ck]; ok {
		// 	*cacheHit++
		// 	return v
		// } else {
		// 	*cacheMiss++
		// }
		// defer func() { cache[ck] = result }()

		// fmt.Println(section, zIn)
		mapDigitZout := mapSectionZinDigitZout[section][zIn]
		if section == len(mapSectionZinDigitZout)-1 {
			// last section, so we can check the target
			partialPaths := make(map[int][]int)
			for digit, zOut := range mapDigitZout {
				if zOut == targetZOut {
					partialPaths[digit] = []int{digit}
				}
			}
			if len(partialPaths) == 0 {
				return []int{}
			} else {
				// fmt.Printf("findPath: got %d terminal partial paths: %v\n", len(partialPaths), partialPaths)
				//selectedDigit := getMaxIntKey(partialPaths)
				selectedDigit := getMinIntKey(partialPaths)
				return partialPaths[selectedDigit]
			}
		} else {
			// not the last section, DP
			partialPaths := make(map[int][]int)
			for digit, zOut := range mapDigitZout {
				p := rec(section+1, zOut, targetZOut)
				if len(p) > 0 {
					partialPaths[digit] = p
				}
			}
			if len(partialPaths) == 0 {
				return []int{}
			} else {
				maxDigit := getMinIntKey(partialPaths)
				return append([]int{maxDigit}, partialPaths[maxDigit]...)
			}
		}
	}

	// defer func() {
	// 	fmt.Println("cache hit rate", 100*(*cacheHit)/(*cacheHit+*cacheMiss+1), "%")
	// }()
	return rec(0, 0, goal)
}

func getMaxIntKey(m map[int][]int) int {
	maxDigit := 0
	for d := range m {
		if d > maxDigit {
			maxDigit = d
		}
	}
	advent.Assertf(maxDigit > 0, "expected max int to be larger than 0")
	return maxDigit
}

func getMinIntKey(m map[int][]int) int {
	minDigit := 999
	for d := range m {
		if d < minDigit {
			minDigit = d
		}
	}
	advent.Assertf(minDigit >= 1 && minDigit <= 9, "expected max int to be larger than 0")
	return minDigit
}

func pruneUseless(products []product, goal int) []product {
	// parse product in reverse order sinde they are ordered by section

	currentSection := products[len(products)-1].Section

	desiredZOut := make(map[int]bool)
	desiredZOut[goal] = true
	newDesiredZout := make(map[int]bool)
	pruned := []product{}
	for i := len(products) - 1; i >= 0; i-- {
		p := products[i]
		if p.Section == currentSection {
			if _, ok := desiredZOut[p.ZOut]; ok {
				pruned = append(pruned, p)
				newDesiredZout[p.ZIn] = true
			}
		} else if p.Section < currentSection {
			currentSection = p.Section
			desiredZOut = newDesiredZout
			newDesiredZout = make(map[int]bool)
			fmt.Println("prune, now section", currentSection)
		} else {
			panic("assumed that proudct are sorder by section")
		}
	}
	return pruned
}
