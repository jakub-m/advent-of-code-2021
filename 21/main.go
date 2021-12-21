package main

import "fmt"

func Calc(pos1, pos2 int) (int, error) {
	player1 := position(pos1)
	player2 := position(pos2)
	score1, score2 := 0, 0
	roll := 1
	losingScore := 0
	nRolls := 0
	for {
		roll1, sum := rollTimes(roll, 3)
		fmt.Println("roll state", roll1, "sum", sum)
		nRolls += 3
		roll = roll1
		player1.move(sum)
		score1 += int(player1)
		fmt.Println("player1 pos", player1, "score", score1)
		if score1 >= 1000 {
			losingScore = score2
			break
		}

		roll2, sum := rollTimes(roll, 3)
		fmt.Println("roll state", roll1, "sum", sum)
		nRolls += 3
		roll = roll2
		player2.move(sum)
		score2 += int(player2)
		fmt.Println("player2 pos", player2, "score", score2)
		if score2 >= 1000 {
			losingScore = score1
			break
		}
	}
	return losingScore * nRolls, nil
}

func rollTimes(roll, times int) (int, int) {
	s := 0
	for i := 0; i < times; i++ {
		s += roll
		roll++
		if roll > 100 {
			roll = 1
		}
	}
	return roll, s
}

type position int

const ringSize = 10

func (p *position) move(n int) {
	pi := int(*p)
	*p = position((pi-1+n)%ringSize + 1)
}
