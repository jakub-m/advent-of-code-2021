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
	return 0, nil
	//m := getMinimumEnergy(burrowWithAmphoids, make(map[situation]bool))
	//return m, nil
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
	t += fmt.Sprintf("#%c%c.%c.%c.%c.%c%c#\n", s.roomLeft[1], s.roomLeft[0], s.hallAB, s.hallBC, s.hallCD, s.roomRight[0], s.roomRight[1])
	t += fmt.Sprintf("###%c#%c#%c#%c###\n", s.roomA[0], s.roomB[0], s.roomC[0], s.roomD[0])
	t += fmt.Sprintf("  #%c#%c#%c#%c###\n", s.roomA[1], s.roomB[1], s.roomC[1], s.roomD[1])
	t += fmt.Sprintf("  #########  ", s.roomA[1], s.roomB[1], s.roomC[1], s.roomD[1])
	return t
}

type fieldState rune

const (
	emptyField fieldState = '.'
	amphipodA             = "A"
	amphipodB             = "B"
	amphipodC             = "C"
	amphipodD             = "D"
)

func initialBurrowWithAmphoids() situation {
	s := situation{}
	s.roomLeft[1] = emptyField
	s.
	return s
}

func (sit situation) nextSituationsWithCosts() []situationWithCost {
	panic("todo")
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