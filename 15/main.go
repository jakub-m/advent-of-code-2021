package main

import (
	"advent"
	"io"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

var startPos = advent.Pos{0, 0}

func Calc(r io.Reader) (int, error) {
	riskGrid, err := advent.LoadGrid(r)
	// advent.Printf("%s\n", riskGrid)

	riskGrid = extendRiskGrid(riskGrid)

	advent.Printf("%s\n", riskGrid)

	endPos := findEndPos(riskGrid)
	if err != nil {
		return 0, err
	}

	unvisited := make(map[advent.Pos]bool)
	distances := make(map[advent.Pos]int)
	for pos := range riskGrid {
		distances[pos] = MaxInt / 2 // make room for calculations
		unvisited[pos] = true
	}
	distances[startPos] = 0
	delete(unvisited, startPos)

	//advent.Printf("grid:\n%v\n", riskGrid)

	current := startPos
	for {
		advent.Printf("len(unvisited) = %d\n", len(unvisited))
		// advent.Printf("current = %v\n", current)
		neighbors := riskGrid.GetNeighbors4(current)
		for _, neighbour := range neighbors {
			if !unvisited[neighbour] {
				continue
			}
			newTentativeDistance := distances[current] + riskGrid[neighbour]
			distances[neighbour] = advent.MinInt([]int{
				newTentativeDistance,
				distances[neighbour],
			})
		}
		delete(unvisited, current)

		nextPos := advent.Pos{X: -1, Y: -1}
		nextPosDist := MaxInt

		if _, ok := unvisited[endPos]; !ok {
			break
		}

		for u := range unvisited {
			ud := distances[u]
			// fmt.Printf("u %v ud %d\n", u, ud)
			if ud < nextPosDist {
				nextPosDist = ud
				nextPos = u
			}
		}

		// advent.Printf("nextPos: %v, unvisited: %v\n", nextPos, unvisited[nextPos])

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
		advent.Printf("%v = %d\n", p, r)
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
				currDist = neighDist
			}
		}
	}

	return path
}

const Extension = 5

func extendRiskGrid(original advent.GridInt) advent.GridInt {
	extended := make(advent.GridInt)
	origEnd := original.FindEndPos()
	for origPos, origVal := range original {
		for i := 0; i < Extension; i++ {
			for k := 0; k < Extension; k++ {
				p := advent.Pos{
					X: (origEnd.X+1)*i + origPos.X,
					Y: (origEnd.Y+1)*k + origPos.Y,
				}
				dv := i + k //advent.MaxInt([]int{i, k})
				extended[p] = (origVal-1+dv)%9 + 1
			}
		}
	}
	return extended
	// fmt.Printf("%s\n", extended)
	// panic("done)")
	// return extended
}
