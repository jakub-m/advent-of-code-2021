package main

import (
	"advent"
	"fmt"
	"io"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

var startPos = advent.Pos{0, 0}

func Calc(r io.Reader) (int, error) {
	riskGrid, err := advent.LoadGrid(r)

	endPos := findEndPos(riskGrid)
	if err != nil {
		return 0, err
	}

	unvisited := make(map[advent.Pos]bool)
	distances := make(map[advent.Pos]int)
	for pos := range riskGrid {
		distances[pos] = MaxInt
		unvisited[pos] = true
	}
	distances[startPos] = 0
	delete(unvisited, startPos)

	fmt.Printf("grid:\n%v\n", riskGrid)

	current := startPos
	for {
		fmt.Printf("len(unvisited) = %d\n", len(unvisited))
		fmt.Printf("current = %v\n", current)
		neighbors := riskGrid.GetNeighbors4(current)
		for _, neighbour := range neighbors {
			if !unvisited[neighbour] {
				continue
			}
			newDistThroughCurrent := distances[current] + riskGrid[neighbour]
			distances[neighbour] = advent.MinInt([]int{
				newDistThroughCurrent,
				distances[neighbour],
			})
		}
		delete(unvisited, current)

		nextPos := advent.Pos{X: -1, Y: -1}
		nextPosDist := MaxInt

		for u := range unvisited {
			ud := distances[u]
			// fmt.Printf("u %v ud %d\n", u, ud)
			if ud < nextPosDist {
				nextPosDist = ud
				nextPos = u
			}
		}

		fmt.Printf("nextPos: %v, unvisited: %v\n", nextPos, unvisited[nextPos])

		if _, ok := unvisited[endPos]; !ok {
			break
		}

		advent.Assertf(len(unvisited) > 0, "no more unvisited")
		current = nextPos
	}

	// todo add end field risk
	path := trackPath(distances, startPos, endPos)
	// risks := riskGrid.GetMany(path)

	sumRisk := 0
	for _, p := range path {
		r := riskGrid[p]
		sumRisk += r
		fmt.Printf("%v = %d\n", p, r)
	}
	return sumRisk, nil
}

func findEndPos(grid advent.GridInt) advent.Pos {
	end := startPos
	for p := range grid {
		if p.X > end.X {
			end = p
		} else if p.X == end.X {
			if p.Y > end.Y {
				end = p
			}
		}
	}
	return end
}

func trackPath(distances map[advent.Pos]int, start, end advent.Pos) []advent.Pos {

	distGrid := advent.GridInt(distances)

	path := []advent.Pos{}

	curr := end
	for curr != start {
		path = append(path, curr)
		currDist := distances[curr]
		neighbours := distGrid.GetNeighbors4(curr)
		for _, neigh := range neighbours {
			neighDist := distGrid[neigh]
			if neighDist < currDist {
				curr = neigh
			}
		}
	}

	path = append(path, startPos)
	return path
}
