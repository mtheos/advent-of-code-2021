package dive

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"strings"
)

const (
	F byte = iota
	U
	D
)

type movement struct {
	direction byte
	steps     int
}

type position struct {
	horizontal int
	vertical   int
	aim        int
}

func mapDirection(direction byte) byte {
	switch direction {
	case 'f':
		return F
	case 'u':
		return U
	case 'd':
		return D
	default:
		panic("Unknown direction " + string(direction))
	}
}

func readInput(fileName string) []movement {
	var arr []movement
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			words := strings.Fields(scanner.Text())
			direction := mapDirection(words[0][0])
			num := Atoi(words[1])
			arr = append(arr, movement{direction: direction, steps: num})
		}
	})
	return arr
}

func ezMode(input []movement, ch chan<- int) {
	pos := position{}
	for _, move := range input {
		switch move.direction {
		case F:
			pos.horizontal += move.steps
			break
		case U:
			pos.vertical -= move.steps
			break
		case D:
			pos.vertical += move.steps
			break
		default:
			panic("Unmatched case " + string(move.direction))
		}
	}
	ch <- pos.horizontal * pos.vertical
}

func hardMode(input []movement, ch chan<- int) {
	pos := position{}
	for _, move := range input {
		switch move.direction {
		case F:
			pos.horizontal += move.steps
			pos.vertical += pos.aim * move.steps
			break
		case U:
			pos.aim -= move.steps
			break
		case D:
			pos.aim += move.steps
			break
		default:
			panic("Unmatched case " + string(move.direction))
		}
	}
	ch <- pos.horizontal * pos.vertical
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)

	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Dive!")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
