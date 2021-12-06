package giantSquid

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	scanner := bufio.NewScanner(file)
	var arr []int
	scanner.Scan()
	line := scanner.Text()
	split := strings.Split(line, ",")
	for _, fishy := range split {
		age, err := strconv.Atoi(fishy)
		if err != nil {
			panic(err)
		}
		arr = append(arr, age)
	}
	return arr
}

func ezMode(input []int, period int, ch chan<- int) {
	for ; period > 0; period-- {
		lastLen := len(input)
		for i := 0; i < lastLen; i++ {
			if input[i] == 0 {
				input[i] = 6
				input = append(input, 8)
			} else {
				input[i]--
			}
		}
	}
	ch <- len(input)
}

func hardMode(input []int, period int, ch chan<- int) {
	fishies := make([]int, 9)
	for _, age := range input {
		fishies[age]++
	}
	for ; period > 0; period-- {
		reproduced := fishies[0]
		for i := 0; i < len(fishies)-1; i++ {
			fishies[i] = fishies[i+1]
		}
		fishies[8] = reproduced
		fishies[6] += reproduced
	}
	sum := 0
	for _, count := range fishies {
		sum += count
	}
	ch <- sum
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, 80, ezChan)
	go hardMode(input, 256, hardChan) // WTB more memory...a lot more memory

	ch <- fmt.Sprintln("Lanternfish")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
