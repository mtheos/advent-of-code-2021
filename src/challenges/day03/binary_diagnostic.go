package binaryDiagnostic

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInput(fileName string) []string {
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

	var arr []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		arr = append(arr, line)
	}
	return arr
}

func countBits(input []string) (map[int]int, int) {
	numBits := 0
	bits := make(map[int]int)
	for _, byte_ := range input {
		for i, bit := range byte_ {
			numBits = i + 1
			if bit == '1' {
				bits[i]++
			} else {
				bits[i]--
			}
		}
	}
	return bits, numBits
}

func getGammaAndEpsilon(bits map[int]int, numBits int) (string, string) {
	gamma := ""
	epsilon := ""
	for i := 0; i < numBits; i++ {
		bit := bits[i]
		if bit >= 0 {
			gamma += "1"
			epsilon += "0"
		} else {
			gamma += "0"
			epsilon += "1"
		}
	}
	return gamma, epsilon
}

func whittleDown(input []string, predicate func(int) bool) string {
	currentBit := 0
	for len(input) > 1 {
		count := 0
		for _, byte_ := range input {
			if byte_[currentBit] == '1' {
				count++
			} else {
				count--
			}
		}
		var tmp []string
		var val uint8
		if predicate(count) {
			val = '1'
		} else {
			val = '0'
		}
		for _, byte_ := range input {
			if byte_[currentBit] == val {
				tmp = append(tmp, byte_)
			}
		}
		input = tmp
		currentBit++
	}
	return input[0]
}

func ezMode(input []string, ch chan<- int) {
	bits, numBits := countBits(input)
	gamma, epsilon := getGammaAndEpsilon(bits, numBits)
	g, err := strconv.ParseInt(gamma, 2, 32)
	if err != nil {
		panic(err)
	}
	e, err := strconv.ParseInt(epsilon, 2, 32)
	if err != nil {
		panic(err)
	}
	ch <- int(g * e)
}

func hardMode(input []string, ch chan<- int) {
	predicate := func(count int) bool { return count >= 0 }
	o2 := whittleDown(input, predicate)
	co2 := whittleDown(input, func(count int) bool { return !predicate(count) })
	o, err := strconv.ParseInt(o2, 2, 32)
	if err != nil {
		panic(err)
	}
	c, err := strconv.ParseInt(co2, 2, 32)
	if err != nil {
		panic(err)
	}
	ch <- int(o * c)
}

func Go(ch chan string) {
	input := readInput("./src/challenges/day03/binary_diagnostic.txt")

	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
