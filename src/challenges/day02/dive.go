package dive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var arr []movement
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		direction := mapDirection(words[0][0])
		num, err := strconv.Atoi(words[1])
		if err != nil {
			panic(err)
		}
		arr = append(arr, movement{direction: direction, steps: num})
	}
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

func Go(ch chan string) {
	input := readInput("./src/challenges/day02/dive.txt")

	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
