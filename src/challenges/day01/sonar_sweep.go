package sonarSweep

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInput(fileName string) []int {
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

	var arr []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		arr = append(arr, num)
	}
	return arr
}

func ezMode(input []int, ch chan<- int) {
	increased := 0
	for i := 0; i < len(input)-1; i++ {
		if input[i] < input[i+1] {
			increased++
		}
	}
	ch <- increased
}

func hardMode(input []int, ch chan<- int) {
	increased := 0
	for i := 0; i < len(input)-3; i++ {
		a := input[i] + input[i+1] + input[i+2]
		b := input[i+1] + input[i+2] + input[i+3]
		if a < b {
			increased++
		}
	}
	ch <- increased
}

func Go(ch chan string) {
	input := readInput("./src/challenges/day01/sonar_sweep.txt")

	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
