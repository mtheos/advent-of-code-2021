package main

import (
	day01 "advent-of-code-2021/src/challenges/day01"
	day02 "advent-of-code-2021/src/challenges/day02"
	day03 "advent-of-code-2021/src/challenges/day03"
	day04 "advent-of-code-2021/src/challenges/day04"
	day05 "advent-of-code-2021/src/challenges/day05"
	day06 "advent-of-code-2021/src/challenges/day06"
	day07 "advent-of-code-2021/src/challenges/day07"
	day08 "advent-of-code-2021/src/challenges/day08"
	day09 "advent-of-code-2021/src/challenges/day09"
	day10 "advent-of-code-2021/src/challenges/day10"
	day11 "advent-of-code-2021/src/challenges/day11"
	day12 "advent-of-code-2021/src/challenges/day12"
	day13 "advent-of-code-2021/src/challenges/day13"
	day14 "advent-of-code-2021/src/challenges/day14"
	day15 "advent-of-code-2021/src/challenges/day15"
	day16 "advent-of-code-2021/src/challenges/day16"
	day17 "advent-of-code-2021/src/challenges/day17"
	day18 "advent-of-code-2021/src/challenges/day18"
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
	challenges := createChallenges()
	//challenges = challenges[len(challenges)-1:]
	for _, chal := range challenges {
		go chal.run(chal.input, chal.ch)
	}
	for _, chal := range challenges {
		printResults(chal)
	}
}

func createChallenges() []challenge {
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
	// Lanternfish
	challenges = append(challenges, challenge{
		run:   day06.Go,
		ch:    make(chan string),
		input: "./src/challenges/day06/input.txt",
	})
	// The Treachery of Whales
	challenges = append(challenges, challenge{
		run:   day07.Go,
		ch:    make(chan string),
		input: "./src/challenges/day07/input.txt",
	})
	// Seven Segment Search
	challenges = append(challenges, challenge{
		run:   day08.Go,
		ch:    make(chan string),
		input: "./src/challenges/day08/input.txt",
	})
	// Smoke Basin
	challenges = append(challenges, challenge{
		run:   day09.Go,
		ch:    make(chan string),
		input: "./src/challenges/day09/input.txt",
	})
	// Syntax Scoring
	challenges = append(challenges, challenge{
		run:   day10.Go,
		ch:    make(chan string),
		input: "./src/challenges/day10/input.txt",
	})
	// Dumbo Octopus
	challenges = append(challenges, challenge{
		run:   day11.Go,
		ch:    make(chan string),
		input: "./src/challenges/day11/input.txt",
	})
	// Passage Pathing
	challenges = append(challenges, challenge{
		run:   day12.Go,
		ch:    make(chan string),
		input: "./src/challenges/day12/input.txt",
	})
	// Transparent Origami
	challenges = append(challenges, challenge{
		run:   day13.Go,
		ch:    make(chan string),
		input: "./src/challenges/day13/input.txt",
	})
	// Extended Polymerization
	challenges = append(challenges, challenge{
		run:   day14.Go,
		ch:    make(chan string),
		input: "./src/challenges/day14/input.txt",
	})
	// Chiton
	challenges = append(challenges, challenge{
		run:   day15.Go,
		ch:    make(chan string),
		input: "./src/challenges/day15/input.txt",
	})
	// Packet Decoder
	challenges = append(challenges, challenge{
		run:   day16.Go,
		ch:    make(chan string),
		input: "./src/challenges/day16/input.txt",
	})
	// Trick Shot
	challenges = append(challenges, challenge{
		run:   day17.Go,
		ch:    make(chan string),
		input: "./src/challenges/day17/input.txt",
	})
	// Snailfish
	challenges = append(challenges, challenge{
		run:   day18.Go,
		ch:    make(chan string),
		input: "./src/challenges/day18/input.txt",
	})
	return challenges
}
