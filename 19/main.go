package main

import (
	"advent"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const (
	X = 0
	Y = 1
	Z = 2
)

func Calc(r io.Reader) (int, error) {
	scanners, err := readScanners(r)
	if err != nil {
		return 0, err
	}
	alignedScanners := []scanner{scanners[0]}
	for _, scanner := range scanners[1:] {
		aligned, err := alignScanner(scanner, alignedScanners)
		if err != nil {
			return 0, err
		}
		alignedScanners = append(alignedScanners, aligned)
	}

	return countAllBeacons(alignedScanners), nil
}

func readScanners(r io.Reader) ([]scanner, error) {
	lines, err := advent.ReadLinesTrim(r)
	if err != nil {
		return nil, err
	}
	scanners := []scanner{}
	var currentScanner scanner = nil
	for _, line := range lines {
		if strings.HasPrefix(line, "--") {
			if currentScanner != nil {
				scanners = append(scanners, currentScanner)
			}
			currentScanner = scanner{}
		} else if line == "" {
			continue
		} else {
			beacon, err := lineToBeacon(line)
			if err != nil {
				return nil, err
			}
			currentScanner = append(currentScanner, beacon)
		}
	}
	scanners = append(scanners, currentScanner)
	return scanners, nil
}

func lineToBeacon(line string) (beacon, error) {
	b := beacon{}
	for i, p := range strings.Split(line, ",") {
		val, err := strconv.Atoi(p)
		if err != nil {
			return b, err
		}
		b[i] = val
	}
	return b, nil
}

type scanner []beacon
type beacon [3]int
type transformation func(beacon) beacon

func (b beacon) String() string {
	return fmt.Sprintf("[%d,%d,%d]", b[X], b[Y], b[Z])
}

func (s scanner) String() string {
	str := []string{}
	for _, b := range s {
		str = append(str, b.String())
	}
	sort.Strings(str)
	return strings.Join(str, " ")
}

func (s scanner) transform(tran transformation) scanner {
	newScanner := []beacon{}
	for _, b := range s {
		newScanner = append(newScanner, tran(b))
	}
	return newScanner
}

var rotations []transformation

func rotX(b beacon) beacon {
	return beacon{b[X], b[Z], -b[Y]}
}

func rotY(b beacon) beacon {
	return beacon{b[Z], b[Y], -b[X]}
}

func rotZ(b beacon) beacon {
	return beacon{b[Y], -b[X], b[Z]}
}

func ident(b beacon) beacon {
	return b
}

func chain(transformations ...transformation) transformation {
	return func(b beacon) beacon {
		for _, t := range transformations {
			b = t(b)
		}
		return b
	}
}

// a,  b, c,
// b, -a, c,

func init() {
	for _, x := range getRotations(rotX) {
		for _, y := range getRotations(rotY) {
			for _, z := range getRotations(rotZ) {
				r := chain(x, y, z)
				rotations = append(rotations, r)
			}
		}
	}
}

func getRotations(t transformation) []transformation {
	return []transformation{
		ident,
		t,
		chain(t, t),
		chain(t, t, t),
	}
}

func alignScanner(candidate scanner, alignedScanners []scanner) (scanner, error) {
	for _, aligned := range alignedScanners {
		for _, tran := range rotations {
			transformedCandidate := candidate.transform(tran)
			for _, refBeacon := range aligned {
				offsetTran := getOffsetTran(refBeacon)
				offsetCandidate := transformedCandidate.transform(offsetTran)

				if scannerOverlap(offsetCandidate, aligned) {
					return offsetCandidate, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("cannot align scanner %v", candidate)
}

func getOffsetTran(ref beacon) transformation {
	return func(b beacon) beacon {
		return beacon{
			b[X] - ref[X],
			b[Y] - ref[Y],
			b[Z] - ref[Z],
		}
	}
}

func scannerOverlap(some, other scanner) bool {
	m := make(map[beacon]int)
	for _, b := range some {
		m[b]++
	}
	for _, b := range other {
		m[b]++
	}
	countOverlapping := 0
	for _, v := range m {
		if v == 2 {
			countOverlapping++
		}
	}
	return countOverlapping >= 12
}

func countAllBeacons(scanners []scanner) int {
	m := make(map[beacon]bool)
	for _, scanner := range scanners {
		for _, b := range scanner {
			m[b] = true
		}
	}
	return len(m)
}
