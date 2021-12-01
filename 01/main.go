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
	err := mainerr()
	if err != nil {
		log.Fatal(err)
	}
}

func mainerr() error {
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	//c, err := CountIncreasesReader(f)
	c, err := CountIncreasesWindowedReader(f)
	if err != nil {
		return err
	}
	fmt.Println(c)
	return nil
}

func CountIncreasesReader(r io.Reader) (int, error) {
	values, err := ReadFileAsIntegers(r)
	if err != nil {
		return 0, nil
	}

	increases := CountIncreases(values)
	return increases, nil
}

func CountIncreasesWindowedReader(r io.Reader) (int, error) {
	values, err := ReadFileAsIntegers(r)
	if err != nil {
		return 0, nil
	}

	windowSums := []int{}

	for i := range values {
		j := i + 3
		if j > len(values) {
			break
		}

		window := values[i:j]
		windowSums = append(windowSums, SumInt(window))
	}

	increses := CountIncreases(windowSums)
	return increses, nil
}

func SumInt(values []int) int {
	s := 0
	for _, v := range values {
		s += v
	}
	return s
}

func CountIncreases(values []int) int {
	prev := values[0]
	increases := 0
	for _, val := range values[1:] {
		if val > prev {
			increases++
		}
		prev = val
	}
	return increases
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
