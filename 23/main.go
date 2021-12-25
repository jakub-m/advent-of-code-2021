package main

import (
	"advent"
	"fmt"
)

const (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
)

func Calc() (int, error) {
	burrowWithAmphoids := initialBurrowWithAmphoids()
	fmt.Println(burrowWithAmphoids)
	m := getMinimumEnergy(burrowWithAmphoids, make(map[situation]bool))
	return m, nil
}

func getMinimumEnergy(burrowWithAmphoids situation, alreadyConsideredStates map[situation]bool) int {
	possibleNextStates := []situationWithCost{}
	for _, sc := range burrowWithAmphoids.nextSituationsWithCosts() {
		if !alreadyConsideredStates[sc.situation] {
			possibleNextStates = append(possibleNextStates, sc)
		}
	}

	terminal, nonTerminal := []situationWithCost{}, []situationWithCost{}

	for _, sc := range possibleNextStates {
		if sc.isTerminal {
			terminal = append(terminal, sc)
		} else {
			nonTerminal = append(nonTerminal, sc)
		}
	}

	if len(terminal) > 0 {
		minCost := MaxInt
		for _, sc := range terminal {
			if sc.cost < minCost {
				minCost = sc.cost
			}
		}
		return minCost
	} else if len(nonTerminal) > 0 {
		updatedConsideredStates := cloneConsideredStates(alreadyConsideredStates)
		updatedConsideredStates[burrowWithAmphoids] = true
		nextCosts := []int{}
		for _, sc := range nonTerminal {
			nextCost := sc.cost + getMinimumEnergy(sc.situation, updatedConsideredStates)
			nextCosts = append(nextCosts, nextCost)
		}
		return advent.MinInt(nextCosts)
	}

	panic("no terminal and no non-termianl states")
}

type situation struct {
	roomLeft  [2]fieldState
	roomA     [2]fieldState
	hallAB    fieldState
	roomB     [2]fieldState
	hallBC    fieldState
	roomC     [2]fieldState
	hallCD    fieldState
	roomD     [2]fieldState
	roomRight [2]fieldState
}

func (s situation) String() string {
	t := "#############\n"
	t += fmt.Sprintf("#%s%s.%s.%s.%s.%s%s#\n", s.roomLeft[1], s.roomLeft[0], s.hallAB, s.hallBC, s.hallCD, s.roomRight[0], s.roomRight[1])
	t += fmt.Sprintf("###%s#%s#%s#%s###\n", s.roomA[0], s.roomB[0], s.roomC[0], s.roomD[0])
	t += fmt.Sprintf("  #%s#%s#%s#%s#  \n", s.roomA[1], s.roomB[1], s.roomC[1], s.roomD[1])
	t += "  #########  "
	return t
}

type fieldState int

const (
	emptyField fieldState = iota
	amphipodA
	amphipodB
	amphipodC
	amphipodD
)

func (s fieldState) String() string {
	switch s {
	case emptyField:
		return "."
	case amphipodA:
		return "A"
	case amphipodB:
		return "B"
	case amphipodC:
		return "C"
	case amphipodD:
		return "D"
	default:
		return "?"
	}
}

func (s fieldState) movementCost() int {
	switch s {
	case amphipodA:
		return 1
	case amphipodB:
		return 10
	case amphipodC:
		return 100
	case amphipodD:
		return 1000
	default:
		panic(fmt.Sprintf("movementCost? %+v", s))
	}
}

func initialBurrowWithAmphoids() situation {
	s := situation{}
	s.roomA = [2]fieldState{amphipodB, amphipodA}
	s.roomB = [2]fieldState{amphipodC, amphipodD}
	s.roomC = [2]fieldState{amphipodB, amphipodC}
	s.roomD = [2]fieldState{amphipodD, amphipodA}
	return s
}

func (s situation) nextSituationsWithCosts() []situationWithCost {
	next := []situationWithCost{}

	if t := s.roomLeft[1]; t != emptyField {
		if o := s.roomLeft[0]; o == emptyField {
			s2 := s
			s2.roomLeft[1] = emptyField
			s2.roomLeft[0] = t
			sc := situationWithCost{s2, t.movementCost(), false}
			next = append(next, sc)
		}
	}

	if t := s.roomLeft[0]; t != emptyField {
		if o := s.roomLeft[1]; o == emptyField {
			s2 := s
			s2.roomLeft[0] = emptyField
			s2.roomLeft[1] = t
			sc := situationWithCost{s2, t.movementCost(), false}
			next = append(next, sc)
		}
		if o := s.hallAB; o == emptyField {
			s2 := s
			s2.roomLeft[0] = emptyField
			s2.hallAB = t
			sc := situationWithCost{s2, 2 * t.movementCost(), false}
			next = append(next, sc)
		}
		if o := s.roomA[0]; o == emptyField && t == amphipodA {
			s2 := s
			s2.roomLeft[0] = emptyField
			s.roomA[0] = t
			sc := situationWithCost{s2, 2 * t.movementCost(), true}
			next = append(next, sc)
		}
	}

	// here
	return next
}

func cloneConsideredStates(m map[situation]bool) map[situation]bool {
	c := make(map[situation]bool)
	for k, v := range m {
		c[k] = v
	}
	return c
}

type situationWithCost struct {
	situation  situation
	cost       int
	isTerminal bool
}
