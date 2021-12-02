package advent

import (
	"bufio"
	"io"
	"strings"
)

func ParseLinesReader(r io.Reader, parseLine func(string) (interface{}, error)) ([]interface{}, error) {
	scanner := bufio.NewScanner(r)
	var values []interface{}
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.Trim(t, " \n")
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
