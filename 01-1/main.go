package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	c, err := CountIncreases(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
}

func CountIncreases(r io.Reader) (int, error) {
	values, err := ReadFileAsIntegers(r)
	if err != nil {
		return 0, nil
	}

	prev := values[0]
	increases := 0
	for _, val := range values[1:] {
		if val > prev {
			increases++
		}
		prev = val
	}
	return increases, nil
}

func ReadFileAsIntegers(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	values := []int{}
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.Trim(t, " \n")
		if t == "" {
			continue
		}
		val, err := strconv.Atoi(t)
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
