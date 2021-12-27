// Dijkstra on implicit graph works
// Dynamic programming does not work because the graph is not directed (or I could not make it work)

package main

import (
	"advent"
	"container/heap"
	"fmt"
)

const (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
)

func main() {
	m, err := Calc(initialSituation1())
	if err != nil {
		panic(err)
	}
	fmt.Println("RESULT", m)
}

func Calc(initialSituation situation) (int, error) {
	if advent.PrintEnabled {
		fmt.Println("initial")
		fmt.Println(initialSituation)
		fmt.Println()
		fmt.Println("terminal")
		fmt.Println(terminalSituation)
		fmt.Println()
	}

	// dijkstra for implicit graph
	visited := make(map[situation]bool)
	distance := make(map[situation]int) // use MaxInt if not in set
	distance[initialSituation] = 0

	backlog := &backlogHeap{{initialSituation, 0}}
	backlogSet := make(map[situation]bool)
	iter := 0
	for len(*backlog) > 0 {
		iter++
		backlogHead := heap.Pop(backlog).(situationWithCost)
		current := backlogHead.situation
		backlogSet[current] = false

		var currentDistance int
		if d, ok := distance[current]; ok {
			currentDistance = d
		} else {
			panic(fmt.Sprint("current node has no distance:", current))
		}
		if iter == 10000 {
			advent.Println("len(backlog)", len(*backlog), "currentDistance", currentDistance)
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
				heap.Push(backlog, situationWithCost{sc.situation, distance[sc.situation]})
				backlogSet[sc.situation] = true
			}
		}
		visited[current] = true

		if visited[terminalSituation] {
			break
		}
	}

	m := distance[terminalSituation]

	return m, nil
}

type situation [burrowSize]fieldState

func (s situation) String() string {
	t := "#############\n"
	t += fmt.Sprintf("#%s%s.%s.%s.%s.%s%s#\n", s[roomLeft1], s[roomLeft0], s[hallAB], s[hallBC], s[hallCD], s[roomRight0], s[roomRight1])
	t += fmt.Sprintf("###%s#%s#%s#%s###\n", s[roomA0], s[roomB0], s[roomC0], s[roomD0])
	t += fmt.Sprintf("  #%s#%s#%s#%s#  \n", s[roomA1], s[roomB1], s[roomC1], s[roomD1])
	t += fmt.Sprintf("  #%s#%s#%s#%s#  \n", s[roomA2], s[roomB2], s[roomC2], s[roomD2])
	t += fmt.Sprintf("  #%s#%s#%s#%s#  \n", s[roomA3], s[roomB3], s[roomC3], s[roomD3])
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

func initialSituation1() situation {
	s := situation{}
	s[roomA0] = amphipodB
	s[roomA1] = amphipodD
	s[roomA2] = amphipodD
	s[roomA3] = amphipodA
	s[roomB0] = amphipodC
	s[roomB1] = amphipodC
	s[roomB2] = amphipodB
	s[roomB3] = amphipodD
	s[roomC0] = amphipodB
	s[roomC1] = amphipodB
	s[roomC2] = amphipodA
	s[roomC3] = amphipodC
	s[roomD0] = amphipodD
	s[roomD1] = amphipodA
	s[roomD2] = amphipodC
	s[roomD3] = amphipodA
	return s
}

func initialSituation2() situation {
	s := situation{}
	s[roomA0] = amphipodD
	s[roomA1] = amphipodD
	s[roomA2] = amphipodD
	s[roomA3] = amphipodC
	s[roomB0] = amphipodC
	s[roomB1] = amphipodC
	s[roomB2] = amphipodB
	s[roomB3] = amphipodD
	s[roomC0] = amphipodA
	s[roomC1] = amphipodB
	s[roomC2] = amphipodA
	s[roomC3] = amphipodA
	s[roomD0] = amphipodB
	s[roomD1] = amphipodA
	s[roomD2] = amphipodC
	s[roomD3] = amphipodB
	return s
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
	s[roomA2] = amphipodA
	s[roomA3] = amphipodA
	s[roomB0] = amphipodB
	s[roomB1] = amphipodB
	s[roomB2] = amphipodB
	s[roomB3] = amphipodB
	s[roomC0] = amphipodC
	s[roomC1] = amphipodC
	s[roomC2] = amphipodC
	s[roomC3] = amphipodC
	s[roomD0] = amphipodD
	s[roomD1] = amphipodD
	s[roomD2] = amphipodD
	s[roomD3] = amphipodD
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

	// roomA0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodA, roomA0, roomLeft0, hallAB, roomA1)...)
	next = append(next, moveFromAmphipodRoom12(s, roomA1, roomA0, roomA2, amphipodA)...)
	next = append(next, moveFromAmphipodRoom12(s, roomA2, roomA1, roomA3, amphipodA)...)
	next = append(next, moveFromAmphipodRoom3(s, roomA3, roomA2, amphipodA)...)

	// hallAB
	next = append(next, moveFromHallway(s, hallAB, roomLeft0, hallBC, amphipodA, roomA0, amphipodB, roomB0)...)

	// roomB0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodB, roomB0, hallAB, hallBC, roomB1)...)
	next = append(next, moveFromAmphipodRoom12(s, roomB1, roomB0, roomB2, amphipodB)...)
	next = append(next, moveFromAmphipodRoom12(s, roomB2, roomB1, roomB3, amphipodB)...)
	next = append(next, moveFromAmphipodRoom3(s, roomB3, roomB2, amphipodB)...)

	// hallBC
	next = append(next, moveFromHallway(s, hallBC, hallAB, hallCD, amphipodB, roomB0, amphipodC, roomC0)...)

	// roomC0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodC, roomC0, hallBC, hallCD, roomC1)...)
	next = append(next, moveFromAmphipodRoom12(s, roomC1, roomC0, roomC2, amphipodC)...)
	next = append(next, moveFromAmphipodRoom12(s, roomC2, roomC1, roomC3, amphipodC)...)
	next = append(next, moveFromAmphipodRoom3(s, roomC3, roomC2, amphipodC)...)

	// hallCD
	next = append(next, moveFromHallway(s, hallCD, hallBC, roomRight0, amphipodC, roomC0, amphipodD, roomD0)...)

	// roomD0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodD, roomD0, hallCD, roomRight0, roomD1)...)
	next = append(next, moveFromAmphipodRoom12(s, roomD1, roomD0, roomD2, amphipodD)...)
	next = append(next, moveFromAmphipodRoom12(s, roomD2, roomD1, roomD3, amphipodD)...)
	next = append(next, moveFromAmphipodRoom3(s, roomD3, roomD2, amphipodD)...)

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
	roomA2
	roomA3
	roomB0
	roomB1
	roomB2
	roomB3
	roomC0
	roomC1
	roomC2
	roomC3
	roomD0
	roomD1
	roomD2
	roomD3
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

func moveFromAmphipodRoom12(s situation, start, roomAbove, roomBelow burrowIndex, roomOwner fieldState) []situationWithCost {
	next := []situationWithCost{}

	if s2, amp, ok := s.shift(start, roomAbove); ok {
		if amp != roomOwner {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(start, roomBelow); ok {
		if amp == roomOwner {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	return next
}

func moveFromAmphipodRoom3(s situation, start, roomAbove burrowIndex, roomOwner fieldState) []situationWithCost {
	next := []situationWithCost{}
	if s2, amp, ok := s.shift(start, roomAbove); ok {
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

type backlogHeap []situationWithCost

func (h backlogHeap) Len() int {
	return len(h)
}

func (h backlogHeap) Less(i, j int) bool {
	return h[i].cost < h[j].cost
}

func (h backlogHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *backlogHeap) Push(x interface{}) {
	*h = append(*h, x.(situationWithCost))
}

func (h *backlogHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
