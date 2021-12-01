package sonarSweep

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
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

func ezMode(input []int, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()
	increased := 0
	for i := 0; i < len(input)-1; i++ {
		if input[i] < input[i+1] {
			increased++
		}
	}
	ch <- increased
}

func hardMode(input []int, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()
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

func Go() {
	fmt.Println("Sonar Sweep:")
	input := readInput("./src/challenges/day01/sonar_sweep.txt")

	var wg sync.WaitGroup
	wg.Add(2)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, &wg, ezChan)
	go hardMode(input, &wg, hardChan)

	fmt.Printf("  ezMode: %d\n", <-ezChan)
	fmt.Printf("  hardMode: %d\n", <-hardChan)
	wg.Wait()
}
