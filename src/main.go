package main

import (
	sonarSweep "advent-of-code-2021/src/challenges/day01"
	dive "advent-of-code-2021/src/challenges/day02"
	"fmt"
)

type challenge struct {
	name string
	ch   chan string
	run  func(chan string)
}

func printResults(chal challenge) {
	fmt.Println(chal.name)
	for line := range chal.ch {
		fmt.Print(line)
	}
}

func main() {
	var challenges []challenge
	challenges = append(challenges, challenge{
		name: "Sonar Sweep",
		run:  sonarSweep.Go,
		ch:   make(chan string),
	})
	challenges = append(challenges, challenge{
		name: "Dive!",
		run:  dive.Go,
		ch:   make(chan string),
	})

	for _, chal := range challenges {
		go chal.run(chal.ch)
	}

	for _, chal := range challenges {
		printResults(chal)
	}
}
