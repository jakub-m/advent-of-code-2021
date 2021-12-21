package main

import (
	"advent"
	"fmt"
	"strings"
)

// a special value
const startDiceResult = 0
const winThreshold = 21
const totalFields = 10

func Calc(pos1, pos2 uint8) (int, error) {

	start := universeState{
		a: state{scoreSoFar: 0, field: pos1, combinations: 1},
		b: state{scoreSoFar: 0, field: pos2, combinations: 1},
	}

	backlog := []backlogItem{
		{
			diceResultSum: startDiceResult,
			rollingPlayer: playerNone,
			universeState: start,
		},
	}

	totalUniversesWinsA, totalUniversesWinsB := 0, 0
	for len(backlog) > 0 {
		last := backlog[len(backlog)-1]
		scoreSoFarA := last.universeState.a.scoreSoFar
		scoreSoFarB := last.universeState.b.scoreSoFar
		if advent.PrintEnabled {
			advent.Printf("%s (%d, %d)\n", backlogAsString(backlog), scoreSoFarA, scoreSoFarB)
		}
		if scoreSoFarA >= winThreshold {
			// "a" wins
			totalUniversesWinsA += last.universeState.a.combinations
			// remove the last element from the backlog and add the new one.
			// trim backlog
			backlog = iterateBacklog(backlog)
		} else if scoreSoFarB >= winThreshold {
			// "b" wins
			totalUniversesWinsB += last.universeState.b.combinations
			backlog = iterateBacklog(backlog)
		} else if scoreSoFarA >= winThreshold && scoreSoFarB >= winThreshold {
			panic(fmt.Sprintf("both cannot win: %v", backlog))
		} else {
			// nobody wins, enlarge the backlog
			backlog = enlargeBacklog(backlog)
		}
	}

	if totalUniversesWinsA > totalUniversesWinsB {
		return totalUniversesWinsA, nil
	} else {
		return totalUniversesWinsB, nil
	}
}

type backlogItem struct {
	diceResultSum uint8
	rollingPlayer player
	universeState universeState
}

func (b backlogItem) String() string {
	return fmt.Sprintf("r:%d p:%d %s", b.diceResultSum, b.rollingPlayer, b.universeState)
}

type universeState struct {
	a, b state
}

func (s universeState) String() string {
	return fmt.Sprintf("[a %s, b %s]", s.a, s.b)
}

type state struct {
	scoreSoFar int
	field      uint8
	// not a state "per se", but a product of the combinations that led to this state
	combinations int
}

func (s state) String() string {
	return fmt.Sprintf("s:%d, f:%d", s.scoreSoFar, s.field)
}

type player uint8

func (p player) otherPlayer() player {
	if p == playerA {
		return playerB
	} else if p == playerB {
		return playerA
	} else {
		panic(fmt.Sprint("wrong player?", p))
	}
}

const (
	playerNone player = 0
	playerA    player = 1
	playerB    player = 2
)

func iterateBacklog(backlog []backlogItem) []backlogItem {
	// remove the current element and add a next element. If cannot add next, then remove the previous element and try to add the
	last := backlog[len(backlog)-1]
	trimmedBacklog := backlog[:len(backlog)-1]

	if last.diceResultSum == startDiceResult {
		if len(trimmedBacklog) > 0 {
			panic(fmt.Sprintf("at start but trimmed backlog is not empty: %v", trimmedBacklog))
		}
		return trimmedBacklog
	} else if last.diceResultSum >= 3 && last.diceResultSum < 9 {
		newDiceResultSum := last.diceResultSum + 1
		newRollingPlayer := last.rollingPlayer.otherPlayer()
		newUniverseState := last.universeState

		if newRollingPlayer == playerA {
			newUniverseState.a = stateForNewRoll(newUniverseState.a, newDiceResultSum)
		} else if newRollingPlayer == playerB {
			newUniverseState.b = stateForNewRoll(newUniverseState.b, newDiceResultSum)
		} else {
			panic(fmt.Sprint("rolling player?", newRollingPlayer))
		}

		newHead := backlogItem{
			diceResultSum: newDiceResultSum,
			rollingPlayer: newRollingPlayer,
			universeState: newUniverseState,
		}
		trimmedBacklog = append(trimmedBacklog, newHead)
		return trimmedBacklog

	} else if last.diceResultSum == 9 {
		return iterateBacklog(trimmedBacklog)

	} else {
		panic(fmt.Sprint("unexpected dice result: ", last))
	}
}

func universeStateForNewRoll(u universeState, rollSum uint8, rollingPlayer player) universeState {
	if rollingPlayer == playerA {
		u.a = stateForNewRoll(u.a, rollSum)
		return u
	} else if rollingPlayer == playerB {
		u.b = stateForNewRoll(u.b, rollSum)
		return u
	} else {
		panic(fmt.Sprint("rolling player?", rollingPlayer))
	}
}

func stateForNewRoll(s state, roll uint8) state {
	newField := newFieldForRoll(s.field, roll)
	s.field = newField
	s.scoreSoFar += int(newField)
	s.combinations *= combinationsForRoll(roll)
	return s
}

func newFieldForRoll(field, roll uint8) uint8 {
	s := field + roll
	if s > totalFields {
		s = s - 10
	}

	if s > totalFields {
		panic(fmt.Sprint("field or roll out of bound: ", field, roll))
	}

	return s
}

var combinationsForRollTable []int

func init() {
	combinationsForRollTable = make([]int, 10)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				v := i + j + k
				combinationsForRollTable[v]++
			}
		}
	}
}

func combinationsForRoll(roll uint8) int {
	if roll < 3 || roll > 9 {
		panic(fmt.Sprint("incorrect roll", roll))
	}
	return combinationsForRollTable[roll]
}

func enlargeBacklog(backlog []backlogItem) []backlogItem {
	last := backlog[len(backlog)-1]
	newRoll := uint8(3)

	var newPlayer player
	if len(backlog) == 1 {
		newPlayer = playerA
	} else {
		newPlayer = last.rollingPlayer.otherPlayer()
	}

	newUniverseState := universeStateForNewRoll(last.universeState, newRoll, newPlayer)
	newItem := backlogItem{
		diceResultSum: newRoll,
		rollingPlayer: newPlayer,
		universeState: newUniverseState,
	}
	backlog = append(backlog, newItem)
	return backlog
}

func backlogAsString(backlog []backlogItem) string {
	s := []string{}
	for _, b := range backlog {
		s = append(s, b.String())
	}
	return strings.Join(s, " | ")
}
