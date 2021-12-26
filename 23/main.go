package main

import (
	"advent"
	"fmt"
	"sort"
)

const (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
)

func main() {
	m, err := Calc()
	if err != nil {
		panic(err)
	}
	fmt.Println("RESULT", m)
}

func Calc() (int, error) {
	//m = getMinimumEnergy(initialSituation, make(map[situation]bool))

	// dijkstra for implicit graph
	visited := make(map[situation]bool)
	distance := make(map[situation]int) // use MaxInt if not in set
	distance[initialSituation] = 0

	backlogList := []situation{initialSituation}
	backlogSet := make(map[situation]bool)
	backlogSet[initialSituation] = true
	iter := 0
	for len(backlogList) > 0 {
		iter++
		var current situation
		current, backlogList = backlogList[0], backlogList[1:]
		backlogSet[current] = false

		var currentDistance int
		if d, ok := distance[current]; ok {
			currentDistance = d
		} else {
			panic(fmt.Sprint("current node has no distance:", current))
		}
		if iter == 1000 {
			advent.Println("len(backlog)", len(backlogList), "currentDistance", currentDistance)
			iter = 0
		}

		for _, sc := range current.nextSituationsWithCosts() {
			if visited[sc.situation] {
				continue
			}
			tentativeDistance := currentDistance + sc.cost
			if d, ok := distance[sc.situation]; ok {
				if tentativeDistance < d {
					distance[sc.situation] = tentativeDistance
				}
			} else {
				distance[sc.situation] = tentativeDistance
			}
			if !backlogSet[sc.situation] {
				backlogList = append(backlogList, sc.situation)
				backlogSet[sc.situation] = true
			}
		}
		visited[current] = true

		if visited[terminalSituation] {
			break
		}

		sort.Slice(backlogList, func(i, j int) bool {
			getDistOrMax := func(s situation) int {
				if d, ok := distance[s]; ok {
					return d
				} else {
					return MaxInt
				}
			}
			return getDistOrMax(backlogList[i]) < getDistOrMax(backlogList[j])
		})
	}

	m := distance[terminalSituation]

	return m, nil
}

func getMinimumEnergy(inputSituation situation, alreadyConsideredStates map[situation]bool) int {
	possibleNextStates := []situationWithCost{}
	for _, sc := range inputSituation.nextSituationsWithCosts() {
		if !alreadyConsideredStates[sc.situation] {
			possibleNextStates = append(possibleNextStates, sc)
		}
	}

	if advent.PrintEnabled {
		fmt.Println("====================================")
		fmt.Println(inputSituation)
		fmt.Println()
		fmt.Println("------------------------------------")
		for _, s := range possibleNextStates {
			fmt.Println("cost", s.cost)
			fmt.Println(s.situation)
			fmt.Println()
		}
	}

	updatedConsideredStates := cloneConsideredStates(alreadyConsideredStates)
	updatedConsideredStates[inputSituation] = true
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

var initialSituation situation

//proper initial sit
func init() {
	s := situation{}
	s[roomA0] = amphipodB
	s[roomA1] = amphipodA
	s[roomB0] = amphipodC
	s[roomB1] = amphipodD
	s[roomC0] = amphipodB
	s[roomC1] = amphipodC
	s[roomD0] = amphipodD
	s[roomD1] = amphipodA
	initialSituation = s
}

// func init() {
// s := situation{}
// s[roomLeft0] = amphipodA
// s[roomA1] = amphipodA
// s[roomB0] = amphipodB
// s[roomB1] = amphipodB
// s[roomC0] = amphipodC
// s[roomC1] = amphipodC
// s[roomD0] = amphipodD
// s[roomD1] = amphipodD
// initialSituation = s
// }

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

func (s situation) nextSituationsWithCosts() []situationWithCost {
	if s == terminalSituation {
		return []situationWithCost{{
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
