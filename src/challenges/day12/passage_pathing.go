package passagePathing

import (
	"advent-of-code-2021/src/challenges/day12/graph"
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

func readInput(fileName string) graph.Graph {
	file, err := os.Open(fileName)
	MaybePanic(err)
	defer func(file *os.File) {
		err := file.Close()
		MaybePanic(err)
	}(file)

	scanner := bufio.NewScanner(file)
	g := graph.ArrayGraph{}
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "-")
		from, to := split[0], split[1]
		g.MaybeCreate(from)
		g.MaybeCreate(to)
		g.Connect(from, to, true)
	}
	return &g
}

func alwaysVisitable(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func canVisitSimplex(path string, neighbour string) bool {
	return alwaysVisitable(neighbour) || !strings.Contains(path, neighbour)
}

func canVisitComplex(path string, neighbour string) bool {
	if alwaysVisitable(neighbour) {
		return true
	}
	for _, special := range []string{"start", "end"} {
		if neighbour == special && strings.Contains(path, special) {
			return false
		}
	}
	caves := make(map[string]int)
	caves[neighbour]++
	for _, cave := range strings.Split(path, " -> ") {
		if !alwaysVisitable(cave) {
			if caves[cave] == 2 {
				return false
			}
			caves[cave]++
		}
	}
	foundDouble := false
	for _, count := range caves {
		if count == 2 && foundDouble {
			return false
		}
		if count == 2 {
			foundDouble = true
		}
	}
	return true
}

func findAllPaths(g graph.Graph, canVisit func(path string, neighbour string) bool, node string, path string) []string {
	if node == "end" {
		return []string{path}
	}
	var paths []string
	for _, nIdx := range g.Neighbours(g.Idx(node)) {
		neighbour := g.Name(nIdx)
		if canVisit(path, neighbour) {
			paths = append(paths, findAllPaths(g, canVisit, neighbour, path+" -> "+neighbour)...)
		}
	}
	return paths
}

func findAllPathsFast(g graph.Graph, canVisit func(path string, neighbour string) bool, node string, path string, pathsChan chan []string) {
	if node == "end" {
		pathsChan <- []string{path}
		return
	}
	for _, nIdx := range g.Neighbours(g.Idx(node)) {
		neighbour := g.Name(nIdx)
		if canVisit(path, neighbour) {
			go findAllPathsFast(g, canVisit, neighbour, path+" -> "+neighbour, pathsChan)
		}
	}
}

func ezMode(input graph.Graph, canVisit func(path string, neighbour string) bool, ch chan<- int) {
	paths := findAllPaths(input, canVisit, "start", "start")
	ch <- len(paths)
}

func fastMode(input graph.Graph, canVisit func(path string, neighbour string) bool, ch chan<- int) {
	pathsChan := make(chan []string)
	go findAllPathsFast(input, canVisit, "start", "start", pathsChan)
	count := 0
	oneSecond := time.Duration(1 * 1000 * 1000 * 1000)
	timer := time.NewTimer(oneSecond)
breaker1:
	for {
		select {
		case <-pathsChan:
			count++
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(oneSecond)
		case <-timer.C:
			break breaker1
		}
	}
	ch <- count
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	// No goroutines
	//go ezMode(input, canVisitSimplex, ezChan)
	//go ezMode(input, canVisitComplex, hardChan)

	// Uses a timeout to break when no data is received.
	// This means it could be non-deterministic
	// if a goroutine hangs for too long
	go fastMode(input, canVisitSimplex, ezChan)
	go fastMode(input, canVisitComplex, hardChan)

	ch <- fmt.Sprintln("Passage Pathing")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	// No guarantee which order ez/fast will be written
	// to the channel /shrug
	//ch <- fmt.Sprintf("  ezFastMode: %d\n", <-ezChan)
	//ch <- fmt.Sprintf("  hardFastMode: %d\n", <-hardChan)
	close(ch)
}
