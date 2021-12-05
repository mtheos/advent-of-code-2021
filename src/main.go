package main

import (
	day01 "advent-of-code-2021/src/challenges/day01"
	day02 "advent-of-code-2021/src/challenges/day02"
	day03 "advent-of-code-2021/src/challenges/day03"
	day04 "advent-of-code-2021/src/challenges/day04"
	day05 "advent-of-code-2021/src/challenges/day05"
	"fmt"
)

type challenge struct {
	run   func(string, chan string)
	ch    chan string
	input string
}

func printResults(chal challenge) {
	for line := range chal.ch {
		fmt.Print(line)
	}
}

func main() {
	var challenges []challenge
	// Sonar Sweep
	challenges = append(challenges, challenge{
		run:   day01.Go,
		ch:    make(chan string),
		input: "./src/challenges/day01/input.txt",
	})
	// Dive!
	challenges = append(challenges, challenge{
		run:   day02.Go,
		ch:    make(chan string),
		input: "./src/challenges/day02/input.txt",
	})
	// Binary Diagnostic
	challenges = append(challenges, challenge{
		run:   day03.Go,
		ch:    make(chan string),
		input: "./src/challenges/day03/input.txt",
	})
	// Giant Squid
	challenges = append(challenges, challenge{
		run:   day04.Go,
		ch:    make(chan string),
		input: "./src/challenges/day04/input.txt",
	})
	// Hydrothermal Venture
	challenges = append(challenges, challenge{
		run:   day05.Go,
		ch:    make(chan string),
		input: "./src/challenges/day05/input.txt",
	})

	for _, chal := range challenges {
		go chal.run(chal.input, chal.ch)
	}
	for _, chal := range challenges {
		printResults(chal)
	}
}
