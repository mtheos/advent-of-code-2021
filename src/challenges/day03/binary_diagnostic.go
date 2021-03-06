package binaryDiagnostic

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"strconv"
)

type intPredicate func(i int) bool

func gteZeroPredicate(count int) bool {
	return count >= 0
}

func (predicate intPredicate) mapTrueFalse(truthy uint8, falsy uint8) func(int) uint8 {
	return func(count int) uint8 {
		if predicate(count) {
			return truthy
		} else {
			return falsy
		}
	}
}

func readInput(fileName string) []string {
	var arr []string
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			arr = append(arr, line)
		}
	})
	return arr
}

func countNthBit(input []string, nth int) int {
	count := 0
	for _, byte_ := range input {
		if byte_[nth] == '1' {
			count++
		} else {
			count--
		}
	}
	return count
}

func countAllBits(input []string) []int {
	numBits := len(input[0])
	bits := make([]int, numBits)
	currentBit := 0
	for currentBit < numBits {
		bits[currentBit] = countNthBit(input, currentBit)
		currentBit++
	}
	return bits
}

func getGammaAndEpsilon(bits []int, numBits int) (string, string) {
	gamma := ""
	epsilon := ""
	for i := 0; i < numBits; i++ {
		if bits[i] >= 0 {
			gamma += "1"
			epsilon += "0"
		} else {
			gamma += "0"
			epsilon += "1"
		}
	}
	return gamma, epsilon
}

func whittleDown(input []string, mapper func(int) uint8) string {
	currentBit := 0
	for len(input) > 1 {
		count := countNthBit(input, currentBit)
		var tmp []string
		match := mapper(count)
		for _, byte_ := range input {
			if byte_[currentBit] == match {
				tmp = append(tmp, byte_)
			}
		}
		input = tmp
		currentBit++
	}
	return input[0]
}

func ezMode(input []string, ch chan<- int) {
	numBits := len(input[0])
	bits := countAllBits(input)
	gamma, epsilon := getGammaAndEpsilon(bits, numBits)
	g, err := strconv.ParseInt(gamma, 2, 32)
	MaybePanic(err)
	e, err := strconv.ParseInt(epsilon, 2, 32)
	MaybePanic(err)
	ch <- int(g * e)
}

func hardMode(input []string, ch chan<- int) {
	trueFalseMapper := intPredicate(gteZeroPredicate).mapTrueFalse
	o2Mapper := trueFalseMapper('1', '0')
	co2Mapper := trueFalseMapper('0', '1')
	o2 := whittleDown(input, o2Mapper)
	co2 := whittleDown(input, co2Mapper)

	o, err := strconv.ParseInt(o2, 2, 32)
	MaybePanic(err)
	c, err := strconv.ParseInt(co2, 2, 32)
	MaybePanic(err)

	ch <- int(o * c)
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)

	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Binary Diagnostic")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
