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

	if false {
		fmt.Println("DEBUG")
		initialSituation := initialSituationDebug()
		fmt.Println(initialSituation)
		fmt.Println("===========================")
		fmt.Println()

		for _, sc := range initialSituation.nextSituationsWithCosts() {
			fmt.Println(sc.situation)
			fmt.Println()
		}
		return

	}

	m, err := Calc(initialSituation1())
	if err != nil {
		panic(err)
	}
	fmt.Println("RESULT", m)
}

const enableAStar = false
const enableBackHops = true

func Calc(initialSituation situation) (int, error) {
	fmt.Println(initialSituation)
	backHops := make(map[situation][]situation) // only to extract the path

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
			advent.Println(current)
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

				if enableBackHops {
					if _, ok := backHops[sc.situation]; !ok {
						backHops[sc.situation] = []situation{}
					}
					backHops[sc.situation] = append(backHops[sc.situation], current)

				}
			}
		}
		visited[current] = true

		if visited[terminalSituation] {
			break
		}
	}

	m := distance[terminalSituation]

	if m > 0 && enableBackHops {
		fmt.Println("get path")
		prevDistance := 0
		path := getPath(initialSituation, terminalSituation, distance, backHops)
		for _, p := range path {
			fmt.Println(distance[p]-prevDistance, distance[p])
			fmt.Println(p)
			fmt.Println()
			prevDistance = distance[p]
		}
	}
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

type fieldState uint8

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

func initialSituationDebug() situation {
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

	s[roomD0] = amphipodA
	s[roomD2] = amphipodC
	s[roomD3] = amphipodA

	s[roomRight0] = amphipodD

	return s
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

var roomIndicesA = []burrowIndex{roomA0, roomA1, roomA2, roomA3}
var roomIndicesB = []burrowIndex{roomB0, roomB1, roomB2, roomB3}
var roomIndicesC = []burrowIndex{roomC0, roomC1, roomC2, roomC3}
var roomIndicesD = []burrowIndex{roomD0, roomD1, roomD2, roomD3}

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
	next = append(next, moveSideRoom0(s, roomLeft0, roomLeft1,
		[]burrowIndex{hallAB, hallBC, hallCD, roomRight0},
		amphipodA, roomIndicesA,
	)...)

	// roomA0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodA, roomA0,
		[]burrowIndex{roomLeft0},
		[]burrowIndex{hallAB, hallBC, hallCD, roomRight0},
		[]burrowIndex{roomA1, roomA2, roomA3},
	)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodA, roomA1, roomA0, roomA2, roomA3)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodA, roomA2, roomA1, roomA3)...)
	next = append(next, moveFromAmphipodRoom3(s, amphipodA, roomA3, roomA2)...)

	// hallAB
	next = append(next, moveFromHallway(s, hallAB, roomLeft0, hallBC, amphipodA, roomIndicesA, amphipodB, roomIndicesB)...)

	// roomB0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodB, roomB0,
		[]burrowIndex{hallAB, roomLeft0},
		[]burrowIndex{hallBC, hallCD, roomRight0},
		[]burrowIndex{roomB1, roomB2, roomB3},
	)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodB, roomB1, roomB0, roomB2, roomB3)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodB, roomB2, roomB1, roomB3)...)
	next = append(next, moveFromAmphipodRoom3(s, amphipodB, roomB3, roomB2)...)

	// hallBC
	next = append(next, moveFromHallway(s, hallBC, hallAB, hallCD, amphipodB, roomIndicesB, amphipodC, roomIndicesC)...)

	// roomC0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodC, roomC0,
		[]burrowIndex{hallBC, hallAB, roomLeft0},
		[]burrowIndex{hallCD, roomRight0},
		[]burrowIndex{roomC1, roomC2, roomC3},
	)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodC, roomC1, roomC0, roomC2, roomC3)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodC, roomC2, roomC1, roomC3)...)
	next = append(next, moveFromAmphipodRoom3(s, amphipodC, roomC3, roomC2)...)

	// hallCD
	next = append(next, moveFromHallway(s, hallCD, hallBC, roomRight0, amphipodC, roomIndicesC, amphipodD, roomIndicesD)...)

	// roomD0+
	next = append(next, moveFromAmphipodRoom0(s, amphipodD, roomD0,
		[]burrowIndex{hallCD, hallBC, hallAB, roomLeft0},
		[]burrowIndex{roomRight0},
		[]burrowIndex{roomD1, roomD2, roomD3},
	)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodD, roomD1, roomD0, roomD2, roomD3)...)
	next = append(next, moveFromAmphipodRoom12(s, amphipodD, roomD2, roomD1, roomD3)...)
	next = append(next, moveFromAmphipodRoom3(s, amphipodD, roomD3, roomD2)...)

	// roomRight0
	next = append(next, moveSideRoom0(s, roomRight0, roomRight1,
		[]burrowIndex{hallCD, hallBC, hallAB, roomLeft0},
		amphipodD, roomIndicesD,
	)...)

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

type burrowIndex uint8

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

func hasDifferentAmp(s situation, roomOwner fieldState, rooms []burrowIndex) bool {
	for _, i := range rooms {
		if v := s[i]; v != emptyField && v != roomOwner {
			return true
		}
	}
	return false
}

func moveFromAmphipodRoom0(s situation, roomOwner fieldState, start burrowIndex,
	roomsLeft []burrowIndex,
	roomsRight []burrowIndex,
	roomsBelow []burrowIndex) []situationWithCost {

	roomImmediatelyBelow := roomsBelow[0]

	// this condition is to allow stepping out the target room if there are other amphiphods below, to let them pass.
	roomsBelowHaveDifferentAmp := hasDifferentAmp(s, roomOwner, roomsBelow)

	next := []situationWithCost{}

	if s[start] != roomOwner || roomsBelowHaveDifferentAmp {

		//move to the left
		prevCostLeft := 0
		for _, rLeft := range roomsLeft {
			if s2, amp, ok := s.shift(start, rLeft); ok {
				prevCostLeft += 2 * amp.movementCost()
				sc := situationWithCost{s2, prevCostLeft}
				next = append(next, sc)
			} else {
				break
			}
		}

		//move to the right
		prevCostRight := 0
		for _, rRight := range roomsRight {
			if s2, amp, ok := s.shift(start, rRight); ok {
				prevCostRight += 2 * amp.movementCost()
				sc := situationWithCost{s2, prevCostRight}
				next = append(next, sc)
			} else {
				break
			}
		}
	}

	// if s2, amp, ok := s.shift(start, destOutLeft); ok {
	// 	if amp != roomOwner || roomsBelowHaveDifferentAmp {
	// 		sc := situationWithCost{s2, 2 * amp.movementCost()}
	// 		next = append(next, sc)
	// 	}
	// }

	// if s2, amp, ok := s.shift(start, destOutRight); ok {
	// 	if amp != roomOwner || roomsBelowHaveDifferentAmp {
	// 		sc := situationWithCost{s2, 2 * amp.movementCost()}
	// 		next = append(next, sc)
	// 	}
	// }

	if s2, amp, ok := s.shift(start, roomImmediatelyBelow); ok {
		if amp == roomOwner {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	return next
}

func moveFromHallway(s situation, start burrowIndex,
	hallwayIndexLeft, hallwayIndexRight burrowIndex,
	roomOwnerLeft fieldState, roomIndicesLeft []burrowIndex,
	roomOwnerRight fieldState, roomIndicesRight []burrowIndex) []situationWithCost {

	next := []situationWithCost{}

	// if s2, amp, ok := s.shift(start, hallwayIndexLeft); ok {
	// 	if hallwayIndexLeft != roomLeft0 {
	// 		sc := situationWithCost{s2, 2 * amp.movementCost()}
	// 		next = append(next, sc)
	// 	}
	// }

	// if s2, amp, ok := s.shift(start, hallwayIndexRight); ok {
	// 	if hallwayIndexRight != roomRight0 {
	// 		sc := situationWithCost{s2, 2 * amp.movementCost()}
	// 		next = append(next, sc)
	// 	}
	// }

	hasDifferentOnLeft := hasDifferentAmp(s, roomOwnerLeft, roomIndicesLeft)
	if s2, amp, ok := s.shift(start, roomIndicesLeft[0]); ok {
		if amp == roomOwnerLeft && !hasDifferentOnLeft {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}

	hasDifferentOnRight := hasDifferentAmp(s, roomOwnerRight, roomIndicesRight)
	if s2, amp, ok := s.shift(start, roomIndicesRight[0]); ok {
		if amp == roomOwnerRight && !hasDifferentOnRight {
			sc := situationWithCost{s2, 2 * amp.movementCost()}
			next = append(next, sc)
		}
	}
	return next
}

func moveFromAmphipodRoom12(s situation, roomOwner fieldState, start, roomAbove burrowIndex, roomsBelow ...burrowIndex) []situationWithCost {
	next := []situationWithCost{}
	roomImmediatelyBelow := roomsBelow[0]
	roomsBelowHaveDifferentAmp := hasDifferentAmp(s, roomOwner, roomsBelow)

	if s2, amp, ok := s.shift(start, roomAbove); ok {
		if amp != roomOwner || roomsBelowHaveDifferentAmp {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	if s2, amp, ok := s.shift(start, roomImmediatelyBelow); ok {
		if amp == roomOwner && !roomsBelowHaveDifferentAmp {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}

	return next
}

func moveFromAmphipodRoom3(s situation, roomOwner fieldState, start, roomAbove burrowIndex) []situationWithCost {
	next := []situationWithCost{}
	if s2, amp, ok := s.shift(start, roomAbove); ok {
		if amp != roomOwner {
			sc := situationWithCost{s2, amp.movementCost()}
			next = append(next, sc)
		}
	}
	return next
}

func moveSideRoom0(s situation, start,
	roomIndexDeeper burrowIndex, roomsHallway []burrowIndex,
	roomOwner fieldState, amphibiousRooms []burrowIndex,
) []situationWithCost {

	next := []situationWithCost{}
	if s2, amp, ok := s.shift(start, roomIndexDeeper); ok {
		sc := situationWithCost{s2, amp.movementCost()}
		next = append(next, sc)
	}

	prevCostHallway := 0
	for _, rHallway := range roomsHallway {
		if s2, amp, ok := s.shift(start, rHallway); ok {
			prevCostHallway += 2 * amp.movementCost()
			sc := situationWithCost{s2, prevCostHallway}
			next = append(next, sc)
		} else {
			break
		}
	}

	// if s2, amp, ok := s.shift(start, roomIndexHallway); ok {
	// 	sc := situationWithCost{s2, 2 * amp.movementCost()}
	// 	next = append(next, sc)
	// }

	hasDifferentInRoom := hasDifferentAmp(s, roomOwner, amphibiousRooms)
	if s2, amp, ok := s.shift(start, amphibiousRooms[0]); ok {
		// step in only if i am the owner and there are no other amphiboius
		if amp == roomOwner && !hasDifferentInRoom {
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
	if enableAStar {
		// A* , we know which results are "closer" to the final solution.
		ci := countInPlace(h[i].situation)
		cj := countInPlace(h[j].situation)

		if ci == cj {
			return h[i].cost < h[j].cost
		} else {
			return ci > cj
		}
	} else {
		// dijkstra
		return h[i].cost < h[j].cost
	}
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

func countInPlace(s situation) int {
	c := 0

	if s[roomA0] == amphipodA {
		c += amphipodA.movementCost()
	}
	if s[roomA1] == amphipodA {
		c += amphipodA.movementCost()
	}
	if s[roomA2] == amphipodA {
		c += amphipodA.movementCost()
	}
	if s[roomA3] == amphipodA {
		c += amphipodA.movementCost()
	}

	if s[roomB0] == amphipodB {
		c += amphipodB.movementCost()
	}
	if s[roomB1] == amphipodB {
		c += amphipodB.movementCost()
	}
	if s[roomB2] == amphipodB {
		c += amphipodB.movementCost()
	}
	if s[roomB3] == amphipodB {
		c += amphipodB.movementCost()
	}

	if s[roomC0] == amphipodC {
		c += amphipodC.movementCost()
	}
	if s[roomC1] == amphipodC {
		c += amphipodC.movementCost()
	}
	if s[roomC2] == amphipodC {
		c += amphipodC.movementCost()
	}
	if s[roomC3] == amphipodC {
		c += amphipodC.movementCost()
	}

	if s[roomD0] == amphipodD {
		c += amphipodD.movementCost()
	}
	if s[roomD1] == amphipodD {
		c += amphipodD.movementCost()
	}
	if s[roomD2] == amphipodD {
		c += amphipodD.movementCost()
	}
	if s[roomD3] == amphipodD {
		c += amphipodD.movementCost()
	}

	return c
}

func getPath(initialSituation, terminalSituation situation, distances map[situation]int, backHops map[situation][]situation) []situation {
	current := terminalSituation

	path := []situation{current}
	for current != initialSituation {
		var sWithMinCost situation
		for i, s := range backHops[current] {
			if i == 0 {
				sWithMinCost = s
			}

			if distances[s] < distances[sWithMinCost] {
				sWithMinCost = s
			}
		}

		path = append([]situation{sWithMinCost}, path...)
		current = sWithMinCost
	}

	return path
}
