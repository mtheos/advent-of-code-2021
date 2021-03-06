package sonarSweep

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
)

func readInput(fileName string) []int {
	var arr []int
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			num := Atoi(line)
			arr = append(arr, num)
		}
	})
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

func Go(fileName string, ch chan string) {
	input := readInput(fileName)

	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Sonar Sweep")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
