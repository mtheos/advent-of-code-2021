package extendedPolymerization

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func readInput(fileName string) (string, map[string]uint8) {
	file, err := os.Open(fileName)
	MaybePanic(err)
	defer func(file *os.File) {
		err := file.Close()
		MaybePanic(err)
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	polymer := scanner.Text()
	transformations := make(map[string]uint8)
	scanner.Scan()
	scanner.Text() // blank line
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " -> ")
		transformations[split[0]] = split[1][0]
	}
	return polymer, transformations
}

func countPolymers(chain string) map[string]int {
	polymers := make(map[string]int)
	for i := 0; i < len(chain)-1; i++ {
		pair := chain[i : i+2]
		polymers[pair]++
	}
	return polymers
}

// Not sure if this actually keeps track of enough information to recreate the polymer /shrug
func polymerizeButWithDP(polymers map[string]int, transformations map[string]uint8) map[string]int {
	nextPolymers := make(map[string]int)
	for k, v := range polymers {
		monomer, ok := transformations[k]
		if !ok {
			panic("I don't know how to deal with rejection :(")
		}
		prePair := string([]byte{k[0], monomer})
		postPair := string([]byte{monomer, k[1]})
		nextPolymers[prePair] += v
		nextPolymers[postPair] += v
	}
	return nextPolymers
}

func doubleCountMonomersButFixItToo(polymers map[string]int) map[uint8]int {
	monomers := make(map[uint8]int)
	for k, v := range polymers {
		monomers[k[0]] += v
		monomers[k[1]] += v
	}
	for k, _ := range monomers {
		monomers[k] /= 2
	}
	return monomers
}

func polymerize(chain string, transformations map[string]uint8) string {
	sb := strings.Builder{}
	for i := 0; i < len(chain)-1; i++ {
		pair := chain[i : i+2]
		monomer, ok := transformations[pair]
		if !ok {
			panic("I don't know how to deal with rejection :(")
		}
		sb.WriteByte(chain[i])
		sb.WriteByte(monomer)
	}
	sb.WriteByte(chain[len(chain)-1]) // last char
	return sb.String()
}

func countMonomers(chain string) map[rune]int {
	monomers := make(map[rune]int)
	for _, monomer := range chain {
		monomers[monomer]++
	}
	return monomers
}

func ezMode(chain string, transformations map[string]uint8, steps int, ch chan<- int) {
	for ; steps > 0; steps-- {
		chain = polymerize(chain, transformations)
	}
	monomers := countMonomers(chain)
	max := 0
	min := math.MaxInt32
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

func hardMode(chain string, transformations map[string]uint8, steps int, ch chan<- int) {
	polymers := countPolymers(chain)
	for ; steps > 0; steps-- {
		polymers = polymerizeButWithDP(polymers, transformations)
	}
	monomers := doubleCountMonomersButFixItToo(polymers)
	// first and last monomers aren't double counted
	monomers[chain[0]]++
	monomers[chain[len(chain)-1]]++
	max := 0
	min := math.MaxInt64
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
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(polymer, transformations, 10, ezChan)
	go hardMode(polymer, transformations, 40, hardChan)

	ch <- fmt.Sprintln("Extended Polymerization")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
