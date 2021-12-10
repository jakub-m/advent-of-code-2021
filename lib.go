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
