package extendedPolymerization

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"math"
	"strings"
)

func readInput(fileName string) (string, map[uint16]uint8) {
	var polymer string
	transformations := make(map[uint16]uint8)
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		scanner.Scan()
		polymer = scanner.Text()
		scanner.Scan()
		scanner.Text() // blank line
		for scanner.Scan() {
			split := strings.Split(scanner.Text(), " -> ")
			// This is as obtuse as I could be while being reasonable
			pair := uint16(split[0][0])<<8 + uint16(split[0][1])
			transformations[pair] = split[1][0]
		}
	})
	return polymer, transformations
}

func countPolymers(chain string) map[uint16]int64 {
	polymers := make(map[uint16]int64)
	for i := 0; i < len(chain)-1; i++ {
		pair := uint16(chain[i])<<8 + uint16(chain[i+1])
		polymers[pair]++
	}
	return polymers
}

// Not sure if this actually keeps track of enough information to recreate the polymer /shrug
func polymerizeButWithDP(polymers map[uint16]int64, transformations map[uint16]uint8) map[uint16]int64 {
	nextPolymers := make(map[uint16]int64)
	for pair, v := range polymers {
		monomer, ok := transformations[pair]
		if !ok {
			panic("I don't know how to deal with rejection :(")
		}
		prePair := uint16(uint8(pair>>8))<<8 + uint16(monomer)
		postPair := uint16(monomer)<<8 + uint16(uint8(pair))
		nextPolymers[prePair] += v
		nextPolymers[postPair] += v
	}
	return nextPolymers
}

func doubleCountMonomersButFixItToo(polymers map[uint16]int64) map[uint8]int64 {
	monomers := make(map[uint8]int64)
	for k, v := range polymers {
		monomers[uint8(k>>8)] += v
		monomers[uint8(k)] += v
	}
	for k, _ := range monomers {
		monomers[k] /= 2
	}
	return monomers
}

func ezMode(chain string, transformations map[uint16]uint8, steps int, ch chan<- int64) {
	polymers := countPolymers(chain)
	for ; steps > 0; steps-- {
		polymers = polymerizeButWithDP(polymers, transformations)
	}
	monomers := doubleCountMonomersButFixItToo(polymers)
	// first and last monomers aren't double counted
	monomers[chain[0]]++
	monomers[chain[len(chain)-1]]++
	var max int64 = 0
	var min int64 = math.MaxInt64
	for _, v := range monomers {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	ch <- max - min
}

func Go(fileName string, ch chan string) {
	polymer, transformations := readInput(fileName)
	ezChan := make(chan int64)
	hardChan := make(chan int64)

	go ezMode(polymer, transformations, 10, ezChan)
	go ezMode(polymer, transformations, 40, hardChan)

	ch <- fmt.Sprintln("Extended Polymerization")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
