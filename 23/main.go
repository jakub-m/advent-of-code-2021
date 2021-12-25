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
	//m := getMinimumEnergy(burrowWithAmphoids, make(map[situation]bool))
	for _, sc := range burrowWithAmphoids.nextSituationsWithCosts() {
		fmt.Println()
		fmt.Println(sc.situation)
	}
	return 0, nil
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
		nextCost := sc.cost + getMinimumEnergy(sc.situation, updatedConsideredStates)
		nextCosts = append(nextCosts, nextCost)
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

	// roomLeft1
	if s2, amp, ok := s.shift(roomLeft1, roomLeft0); ok {
		sc := situationWithCost{s2, amp.movementCost()}
		next = append(next, sc)
	}

	// roomLeft0
	if s2, amp, ok := s.shift(roomLeft0, roomLeft1); ok {
		sc := situationWithCost{s2, amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(roomLeft0, hallAB); ok {
		sc := situationWithCost{s2, 2 * amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(roomLeft0, roomA0); ok {
		if amp == amphipodA {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	// roomA0
	if s2, amp, ok := s.shift(roomA0, roomLeft0); ok {
		if amp != amphipodA {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(roomA0, hallAB); ok {
		if amp != amphipodA {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(roomA0, roomA1); ok {
		if amp == amphipodA {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	// roomA1
	if s2, amp, ok := s.shift(roomA1, roomA0); ok {
		if amp != amphipodA {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	// hallAB
	if s2, amp, ok := s.shift(hallAB, roomLeft0); ok {
		sc := situationWithCost{s2, 2 * amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(hallAB, hallBC); ok {
		sc := situationWithCost{s2, 2 * amp.movementCost()}
		next = append(next, sc)
	}

	if s2, amp, ok := s.shift(hallAB, roomA0); ok {
		if amp == amphipodA {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(hallAB, roomB0); ok {
		if amp == amphipodB {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	// roomB0
	if s2, amp, ok := s.shift(roomB0, hallAB); ok {
		if amp != amphipodB {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(roomB0, hallBC); ok {
		if amp != amphipodB {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(roomB0, roomB1); ok {
		if amp == amphipodB {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	// roomB1
	if s2, amp, ok := s.shift(roomB1, roomB0); ok {
		if amp != amphipodB {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}
	// roomC0
	// roomC1
	// roomD0
	// roomD1
	// roomRight0
	// roomRight1
	// hallBC
	// hallCD

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
