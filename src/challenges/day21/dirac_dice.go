package diracDice

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
)

type dice interface {
	roll() int
	rollN(n int) []int
	rollNSum(n int) int
	rolled() int
}

type player struct {
	position, score int
}

func (p *player) moveAndScore(n int) {
	p.position = p.position + n
	if p.position > 10 {
		p.position -= p.position - p.position%10
	}
	if p.position == 0 {
		p.position = 10
	}
	p.score += p.position
}

type deterministicDie struct {
	lastRoll, minFace, maxFace, numRolls int
}

func (d deterministicDie) rolled() int {
	return d.numRolls
}

func (d *deterministicDie) roll() int {
	d.numRolls++
	if d.lastRoll < d.minFace {
		d.lastRoll = d.minFace - 1
	}
	d.lastRoll = d.lastRoll + 1
	if d.lastRoll > d.maxFace {
		d.lastRoll = d.minFace
	}
	return d.lastRoll
}

func (d *deterministicDie) rollN(n int) []int {
	return []int{d.roll(), d.roll(), d.roll()}
}

func (d *deterministicDie) rollNSum(n int) int {
	sum := 0
	for _, v := range d.rollN(n) {
		sum += v
	}
	return sum
}

func readInput(fileName string) (player, player) {
	player1, player2 := player{}, player{}
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		const match = "Player x starting position: "
		scanner.Scan()
		line := scanner.Text()
		player1.position = Atoi(line[len(match):])
		scanner.Scan()
		line = scanner.Text()
		player2.position = Atoi(line[len(match):])
	})
	return player1, player2
}

func (p player) isWinner() bool {
	return p.score >= 1000
}

func getLoser(player1, player2 player) player {
	if player1.isWinner() {
		return player2
	}
	return player1
}

func ezMode(player1, player2 player, ch chan<- int) {
	p1 := true
	var die dice
	die = &deterministicDie{minFace: 1, maxFace: 100}
	for !player1.isWinner() && !player2.isWinner() {
		result := die.rollNSum(3)
		if p1 {
			player1.moveAndScore(result)
		} else {
			player2.moveAndScore(result)
		}
		p1 = !p1
	}
	ch <- getLoser(player1, player2).score * die.rolled()
}

func hardMode(player1, player2 player, ch chan<- int) {
	ch <- 0
}

func Go(fileName string, ch chan string) {
	player1, player2 := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(player1, player2, ezChan)
	go hardMode(player1, player2, hardChan)

	ch <- fmt.Sprintln("Dirac Dice")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
