package main

const (
	maxDim = 1000
)

func Calc(xStart, xEndIncl, yStart, yEndIncl int) (int, error) {
	box := rect{
		xStart: xStart,
		xEnd:   xEndIncl + 1,
		yStart: yStart,
		yEnd:   yEndIncl + 1,
	}

	maxY := 0
	for x := 0; x < box.xEnd; x++ {
		for y := 0; y < maxDim; y++ {
			p := xy{x, y}
			if topY, ok := simulate(p, box); ok {
				if topY > maxY {
					maxY = topY
				}
			}
		}
	}
	return maxY, nil
}

func simulate(vector xy, target rect) (int, bool) {
	pos := xy{0, 0}
	maxY, wasHit := 0, false

	for !(wasHit || (pos.y <= target.yEnd && pos.y < target.yStart)) {
		// advent.Printf("pos %v vector %v\n", pos, vector)
		if pos.y > maxY {
			maxY = pos.y
		}
		if pos.x >= target.xStart && pos.x < target.xEnd && pos.y >= target.yStart && pos.y < target.yEnd {
			wasHit = true
		}
		pos.x += vector.x
		pos.y += vector.y

		if vector.x > 0 {
			vector.x--
		} else if vector.x < 0 {
			vector.x++
		}
		vector.y--
	}
	return maxY, wasHit
}

type rect struct {
	xStart, xEnd, yStart, yEnd int
}

type xy struct {
	x, y int
}
