package advent

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

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

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Assert(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}
