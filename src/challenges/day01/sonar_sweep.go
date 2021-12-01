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

func ezMode(input []int) {
	increased := 0
	for i := 0; i < len(input)-1; i++ {
		if input[i] < input[i+1] {
			increased++
		}
	}
	fmt.Printf("  ezMode: %d\n", increased)
}

func hardMode(input []int) {
	increased := 0
	for i := 0; i < len(input)-3; i++ {
		a := input[i] + input[i+1] + input[i+2]
		b := input[i+1] + input[i+2] + input[i+3]
		if a < b {
			increased++
		}
	}
	fmt.Printf("  hardMode: %d\n", increased)
}

func Go() {
	fmt.Println("Sonar Sweep:")
	input := readInput("./src/challenges/day01/sonar_sweep.txt")
	ezMode(input)
	hardMode(input)
}
