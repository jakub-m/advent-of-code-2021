package advent

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func ReadLinesTrim(r io.Reader) ([]string, error) {
	lines := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.Trim(t, " \t\n")
		lines = append(lines, t)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil

}

func ParseLinesReader(r io.Reader, parseLine func(string) (interface{}, error)) ([]interface{}, error) {
	scanner := bufio.NewScanner(r)
	var values []interface{}
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.Trim(t, " \t\n")
		if t == "" {
			continue
		}
		val, err := parseLine(t)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

func ReadIntsFromStringSep(input string, sep string) ([]int, error) {
	values := []int{}
	for _, s := range strings.Split(input, sep) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

func MinInt(ints []int) int {
	min := ints[0]
	for _, n := range ints {
		if n < min {
			min = n
		}
	}
	return min
}

func MaxInt(ints []int) int {
	Assertf(len(ints) > 0, "MaxInt: length of the input is 0")
	max := ints[0]
	for _, n := range ints {
		if n > max {
			max = n
		}
	}
	return max
}

func AbsInt(n int) int {
	if n >= 0 {
		return n
	} else {
		return -1 * n
	}
}

func TrimNextEmptyLines(lines []string) []string {
	for i, line := range lines {
		line = strings.Trim(line, " \n\t")
		if line != "" {
			return lines[i:]
		}
	}
	return []string{}
}

func MedianInt(values []int) int {
	copied := make([]int, len(values))
	copy(copied, values)
	sort.Slice(copied, func(i, j int) bool { return copied[i] < copied[j] })
	n := len(copied)
	Assertf(n%2 == 1, "n is even %d", n)
	return copied[n/2]
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Assertf(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}

type Pos struct {
	X, Y int
}

type GridInt map[Pos]int

func LoadGrid(r io.Reader) (GridInt, error) {
	lines, err := ReadLinesTrim(r)
	if err != nil {
		return nil, err
	}

	landscape := make(GridInt)
	for iRow, row := range lines {
		for iCol, char := range row {
			val, err := strconv.Atoi(fmt.Sprintf("%c", char))
			if err != nil {
				return nil, err
			}
			c := Pos{iCol, iRow}
			landscape[c] = val
		}
	}

	return landscape, nil
}

func (d GridInt) GetMany(coords []Pos) []int {
	values := []int{}
	for _, c := range coords {
		if val, ok := d[c]; ok {
			values = append(values, val)
		} else {
			panic(fmt.Sprintf("No such coord: %+v", c))
		}
	}
	return values
}

func (d GridInt) GetNeighbors4(c Pos) []Pos {
	neighbors := []Pos{}

	for _, c := range []Pos{
		{X: c.X + 1, Y: c.Y},
		{X: c.X - 1, Y: c.Y},
		{X: c.X, Y: c.Y + 1},
		{X: c.X, Y: c.Y - 1},
	} {
		if _, ok := d[c]; ok {
			neighbors = append(neighbors, c)
		}
	}
	return neighbors
}

func (d GridInt) GetNeighbors8(c Pos) []Pos {
	neighbors := []Pos{}

	for _, c := range []Pos{
		{X: c.X + 1, Y: c.Y},
		{X: c.X - 1, Y: c.Y},
		{X: c.X, Y: c.Y + 1},
		{X: c.X, Y: c.Y - 1},

		{X: c.X - 1, Y: c.Y - 1},
		{X: c.X + 1, Y: c.Y - 1},
		{X: c.X - 1, Y: c.Y + 1},
		{X: c.X + 1, Y: c.Y + 1},
	} {
		if _, ok := d[c]; ok {
			neighbors = append(neighbors, c)
		}
	}
	return neighbors
}

func (d GridInt) FindEndPos() Pos {
	end := Pos{0, 0}
	for p := range d {
		if p.X > end.X {
			end = p
		} else if p.X == end.X {
			if p.Y > end.Y {
				end = p
			}
		}
	}
	return end
}

func (d GridInt) String() string {
	end := d.FindEndPos()

	s := ""
	for y := 0; y <= end.Y; y++ {
		for x := 0; x <= end.X; x++ {
			v := d[Pos{X: x, Y: y}]
			if v < 0 || v > 9 {
				s += "!"
			} else {
				s += fmt.Sprintf("%d", v)
			}
		}
		s += "\n"
	}
	return s

}

var PrintEnabled = true

func Printf(f string, args ...interface{}) {
	if PrintEnabled {
		fmt.Printf(f, args...)
	}
}

func Println(args ...interface{}) {
	if PrintEnabled {
		fmt.Println(args...)
	}
}

func Atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic("Atoi: " + s)
	}
	return v
}

func UniqInt(values []int) []int {
	m := make(map[int]bool)
	uniq := []int{}
	for _, v := range values {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = true
		uniq = append(uniq, v)
	}
	return uniq
}
