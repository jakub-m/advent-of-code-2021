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
		return minCost // TODO here terminal might not be true if amphoid is in the first of two fields of its destinatin room
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

type situation [burrowSize]fieldState

func (s situation) String() string {
	t := "#############\n"
	t += fmt.Sprintf("#%s%s.%s.%s.%s.%s%s#\n", s[roomLeft1], s[roomLeft0], s[hallAB], s[hallBC], s[hallCD], s[roomRight0], s[roomRight1])
	t += fmt.Sprintf("###%s#%s#%s#%s###\n", s[roomA0], s[roomB0], s[roomC0], s[roomD0])
	t += fmt.Sprintf("  #%s#%s#%s#%s#  \n", s[roomA1], s[roomB1], s[roomC1], s[roomD1])
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
	s[roomA0] = amphipodB
	s[roomA1] = amphipodA
	s[roomB0] = amphipodC
	s[roomB1] = amphipodD
	s[roomC0] = amphipodB
	s[roomC1] = amphipodC
	s[roomD0] = amphipodD
	s[roomD1] = amphipodA
	return s
}

func (s situation) nextSituationsWithCosts() []situationWithCost {
	next := []situationWithCost{}

	if t := s[roomLeft1]; t != emptyField {
		if o := s[roomLeft0]; o == emptyField {
			s2 := s
			s2[roomLeft1] = emptyField
			s2[roomLeft0] = t
			sc := situationWithCost{s2, t.movementCost(), false}
			next = append(next, sc)
		}
	}

	// todo here roomLeft0

	if t := s[roomA0]; t != emptyField {
		if o := s[roomLeft0]; o == emptyField {
			s2 := s
			s2[roomA0] = emptyField
			s2[roomLeft0] = t
			sc := situationWithCost{s2, 2 * t.movementCost(), false}
			next = append(next, sc)
		}

	}

	// todo move inside destination
	// todo return 0 on final position

	// 	if o := s.hallAB; o == emptyField {
	// 		s2 := s
	// 		s2.roomLeft[0] = emptyField
	// 		s2.hallAB = t
	// 		sc := situationWithCost{s2, 2 * t.movementCost(), false}
	// 		next = append(next, sc)
	// 	}
	// 	if o := s.roomA[0]; o == emptyField && t == amphipodA {
	// 		s2 := s
	// 		s2.roomLeft[0] = emptyField
	// 		s.roomA[0] = t
	// 		sc := situationWithCost{s2, 2 * t.movementCost(), true}
	// 		next = append(next, sc)
	// 	}
	// }

	// here
	// return next
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

type burrowIndex int

const (
	roomLeft0 = iota
	roomLeft1
	roomA0
	roomA1
	roomB0
	roomB1
	roomC0
	roomC1
	roomD0
	roomD1
	roomRight0
	roomRight1
	hallAB
	hallBC
	hallCD
	burrowSize // must be last const
)
