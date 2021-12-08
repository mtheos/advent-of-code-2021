package theTreacheryofWhales

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readInput(fileName string) []int {
	file, err := os.Open(fileName)
	MaybePanic(err)
	defer func(file *os.File) {
		err := file.Close()
		MaybePanic(err)
	}(file)

	scanner := bufio.NewScanner(file)
	var arr []int
	scanner.Scan()
	line := scanner.Text()
	split := strings.Split(line, ",")
	for _, crabs := range split {
		age, err := strconv.Atoi(crabs)
		MaybePanic(err)
		arr = append(arr, age)
	}
	return arr
}

func ezMode(input []int, move func(int, int) int, ch chan<- int) {
	minFuel := math.MaxInt32
	minDelta, maxDelta := input[0], input[len(input)-1]
	for target := minDelta; target <= maxDelta; target++ {
		fuel := 0
		for _, delta := range input {
			fuel += move(delta, target)
		}
		if fuel < minFuel {
			minFuel = fuel
		}
	}
	ch <- minFuel
}

func basicMove(delta int, targetDelta int) int {
	return int(math.Abs(float64(delta - targetDelta)))
}

func crabMove(delta int, targetDelta int) int {
	n := int(math.Abs(float64(delta - targetDelta)))
	return n * (n + 1) / 2
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	sort.Ints(input)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, basicMove, ezChan)
	go ezMode(input, crabMove, hardChan)

	ch <- fmt.Sprintln("The Treachery of Whales")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
