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
	m := getMinimumEnergy(burrowWithAmphoids, make(map[situation]bool))
	// fmt.Println(burrowWithAmphoids)
	// for _, sc := range burrowWithAmphoids.nextSituationsWithCosts() {
	// 	fmt.Println()
	// 	fmt.Println("=============================")
	// 	fmt.Println(sc.situation)
	// 	fmt.Println("=============================")

	// 	for _, sc := range sc.situation.nextSituationsWithCosts() {
	// 		fmt.Println()
	// 		fmt.Println(sc.situation)
	// 	}

	// }
	return m, nil
}

func getMinimumEnergy(burrowWithAmphoids situation, alreadyConsideredStates map[situation]bool) int {
	possibleNextStates := []situationWithCost{}
	for _, sc := range burrowWithAmphoids.nextSituationsWithCosts() {
		if !alreadyConsideredStates[sc.situation] {
			possibleNextStates = append(possibleNextStates, sc)
		}
	}

	updatedConsideredStates := cloneConsideredStates(alreadyConsideredStates)
	updatedConsideredStates[burrowWithAmphoids] = true
	nextCosts := []int{}
	for _, sc := range possibleNextStates {
		nextMin := getMinimumEnergy(sc.situation, updatedConsideredStates)
		if nextMin != MaxInt {
			nextCost := sc.cost + nextMin
			nextCosts = append(nextCosts, nextCost)
		}
	}
	if len(nextCosts) == 0 {
		return MaxInt
	}
	return advent.MinInt(nextCosts)

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

var terminalSituation situation

func init() {
	s := situation{}
	s[roomA0] = amphipodA
	s[roomA1] = amphipodA
	s[roomB0] = amphipodB
	s[roomB1] = amphipodB
	s[roomC0] = amphipodC
	s[roomC1] = amphipodC
	s[roomD0] = amphipodD
	s[roomD1] = amphipodD
	terminalSituation = s
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
	if s == terminalSituation {
		return []situationWithCost{situationWithCost{
			s, 0,
		}}
	}

	next := []situationWithCost{}

	// roomLeft1
	next = append(next, moveSideRoom1(s, roomLeft1, roomLeft0)...)

	// roomLeft0
	next = append(next, moveSideRoom0(s, roomLeft0, roomLeft1, hallAB, amphipodA, roomA0)...)

	// roomA0
	next = append(next, moveFromAmphipodRoom0(s, amphipodA, roomA0, roomLeft0, hallAB, roomA1)...)

	// roomA1
	next = append(next, moveFromAmphipodRoom1(s, roomA1, roomA0, amphipodA)...)

	// hallAB
	next = append(next, moveFromHallway(s, hallAB, roomLeft0, hallBC, amphipodA, roomA0, amphipodB, roomB0)...)

	// roomB0
	next = append(next, moveFromAmphipodRoom0(s, amphipodB, roomB0, hallAB, hallBC, roomB1)...)

	// roomB1
	next = append(next, moveFromAmphipodRoom1(s, roomB1, roomB0, amphipodB)...)

	// hallBC
	next = append(next, moveFromHallway(s, hallBC, hallAB, hallCD, amphipodB, roomB0, amphipodC, roomC0)...)

	// roomC0
	next = append(next, moveFromAmphipodRoom0(s, amphipodC, roomC0, hallBC, hallCD, roomC1)...)

	// roomC1
	next = append(next, moveFromAmphipodRoom1(s, roomC1, roomC0, amphipodC)...)

	// hallCD
	next = append(next, moveFromHallway(s, hallCD, hallBC, roomRight0, amphipodC, roomC0, amphipodD, roomD0)...)

	// roomD0
	next = append(next, moveFromAmphipodRoom0(s, amphipodD, roomD0, hallCD, roomRight0, roomD1)...)

	// roomD1
	next = append(next, moveFromAmphipodRoom1(s, roomD1, roomD0, amphipodD)...)

	// roomRight0
	next = append(next, moveSideRoom0(s, roomRight0, roomRight1, hallCD, amphipodD, roomD0)...)

	// roomRight1
	next = append(next, moveSideRoom1(s, roomRight1, roomRight0)...)

	return next
}

func (s situation) shift(start, end burrowIndex) (situation, fieldState, bool) {
	this := s[start]
	if this == emptyField {
		return s, this, false
	}
	other := s[end]
	if other != emptyField {
		return s, this, false
	}
	s2 := s
	s2[start], s2[end] = s2[end], s2[start]
	return s2, this, true
}

func cloneConsideredStates(m map[situation]bool) map[situation]bool {
	c := make(map[situation]bool)
	for k, v := range m {
		c[k] = v
	}
	return c
}

type situationWithCost struct {
	situation situation
	cost      int
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

func moveFromAmphipodRoom0(s situation, roomOwner fieldState, start, destOutLeft, destOutRight, room1Index burrowIndex) []situationWithCost {
	next := []situationWithCost{}

	if s2, amp, ok := s.shift(start, destOutLeft); ok {
		if amp != roomOwner || s[room1Index] != roomOwner {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(start, destOutRight); ok {
		if amp != roomOwner || s[room1Index] != roomOwner {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(start, room1Index); ok {
		if amp == roomOwner {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	return next
}

func moveFromHallway(s situation, start burrowIndex,
	hallwayIndexLeft, hallwayIndexRight burrowIndex,
	roomOwnerLeft fieldState, roomIndexLeft burrowIndex,
	roomOwnerRitht fieldState, roomIndexRight burrowIndex) []situationWithCost {

	next := []situationWithCost{}

	if s2, amp, ok := s.shift(start, hallwayIndexLeft); ok {
		sc := situationWithCost{s2, 2 * amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(start, hallwayIndexRight); ok {
		sc := situationWithCost{s2, 2 * amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(start, roomIndexLeft); ok {
		if amp == roomOwnerLeft {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(start, roomIndexRight); ok {
		if amp == roomOwnerRitht {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}
	return next
}

func moveFromAmphipodRoom1(s situation, start, room0Index burrowIndex, roomOwner fieldState) []situationWithCost {
	next := []situationWithCost{}
	if s2, amp, ok := s.shift(start, room0Index); ok {
		if amp != roomOwner {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}
	return next
}

func moveSideRoom0(s situation, start, destSideRoom1, destHallway burrowIndex, roomOwner fieldState, amphibiousRoom burrowIndex) []situationWithCost {
	next := []situationWithCost{}
	if s2, amp, ok := s.shift(start, destSideRoom1); ok {
		sc := situationWithCost{s2, amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(start, destHallway); ok {
		sc := situationWithCost{s2, 2 * amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(start, amphibiousRoom); ok {
		if amp == roomOwner {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	return next
}

func moveSideRoom1(s situation, start, destSideRoom0 burrowIndex) []situationWithCost {
	next := []situationWithCost{}
	if s2, amp, ok := s.shift(start, destSideRoom0); ok {
		sc := situationWithCost{s2, amp.movementCost()}
		next = append(next, sc)
	}
	return next
}
