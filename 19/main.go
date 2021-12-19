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

func Calc(r io.Reader, threshold int) (int, error) {
	scanners, err := readScanners(r)
	if err != nil {
		return 0, err
	}

	referenceScanner := scanners[0]
	otherScanners := scanners[1:]

	vectors := []beacon{}

	for len(otherScanners) > 0 {
		fmt.Println("other scanners", len(otherScanners))
		succeededAlign := false
		for i, candidateScanner := range otherScanners {
			mergedReference, _, vec, ok := alignScanner(referenceScanner, candidateScanner, threshold)
			if !ok {
				continue
			}
			// fmt.Println("vector", vector)
			// fmt.Println("overlap", overlap)
			succeededAlign = true
			referenceScanner = mergedReference
			// fmt.Println("length", len(otherScanners), "i", i)
			if i == len(otherScanners)-1 {
				otherScanners = otherScanners[:i]
			} else {
				otherScanners = append(otherScanners[:i], otherScanners[i+1:]...)
			}
			vectors = append(vectors, vec)
			break
		}
		if !succeededAlign {
			return 0, fmt.Errorf("failed to align")
		}
	}

	return findLargestManhattan(vectors), nil
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
		} else if strings.HasPrefix(line, "#") {
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
	return fmt.Sprintf("%d,%d,%d", b[X], b[Y], b[Z])
}

func (b beacon) negative() beacon {
	return beacon{-b[X], -b[Y], -b[Z]}
}

func (b beacon) add(o beacon) beacon {
	return beacon{
		b[X] + o[X],
		b[Y] + o[Y],
		b[Z] + o[Z],
	}
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
	// returns 64 rotations but 24 would be enough
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

func alignScanner(referenceScanner, candidateScanner scanner, threshold int) (scanner, []beacon, beacon, bool) {
	for _, referenceBeacon := range referenceScanner {
		translateToReferenceBeacon := getOffsetTran(referenceBeacon)
		referenceAtZero := referenceScanner.transform(translateToReferenceBeacon)

		for _, rot := range rotations {
			candidateRotated := candidateScanner.transform(rot)
			for _, candidateBeacon := range candidateRotated {
				translateToCandidateBeacon := getOffsetTran(candidateBeacon)
				candidateAtZero := candidateRotated.transform(translateToCandidateBeacon)
				if overlap, ok := scannerOverlap(referenceAtZero, candidateAtZero, threshold); ok {
					baNeg := referenceBeacon.negative()
					counterTran := getOffsetTran(baNeg)
					// fmt.Println("overlap", scanner(overlap).transform(counterTran))
					// fmt.Println("baAlign", referenceBeacon, "baCand", candidateBeacon)
					vec := referenceBeacon.add(candidateBeacon.negative())
					merged := mergeScanners(candidateAtZero, referenceAtZero).transform(counterTran)
					overlap = scanner(overlap).transform(counterTran)
					// fmt.Println("vec", vec)
					// fmt.Println("overlap", len(overlap), ":", overlap)
					return merged, overlap, vec, true
				}
			}
		}
	}
	return nil, nil, beacon{}, false
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

func scannerOverlap(some, other scanner, threshold int) ([]beacon, bool) {
	m := make(map[beacon]int)
	for _, b := range some {
		m[b]++
	}
	for _, b := range other {
		m[b]++
	}
	overlappingBeacons := []beacon{}
	for b, v := range m {
		if v == 2 {
			overlappingBeacons = append(overlappingBeacons, b)
		}
	}
	return overlappingBeacons, len(overlappingBeacons) >= threshold
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

func mergeScanners(some, other scanner) scanner {
	m := make(map[beacon]bool)
	for _, b := range some {
		m[b] = true
	}
	for _, b := range other {
		m[b] = true
	}

	merged := []beacon{}
	for b := range m {
		merged = append(merged, b)
	}
	return merged
}

func findLargestManhattan(vectors []beacon) int {
	maxManhattan := 0
	for _, some := range vectors {
		for _, other := range vectors {
			m := manhattan(some, other)
			if m > maxManhattan {
				maxManhattan = m
			}
		}
	}
	return maxManhattan
}

func manhattan(some, other beacon) int {
	return advent.AbsInt(some[X]-other[X]) + advent.AbsInt(some[Y]-other[Y]) + advent.AbsInt(some[Z]-other[Z])
}
