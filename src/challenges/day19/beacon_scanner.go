package beaconScanner

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	hs "github.com/emirpasic/gods/sets/hashset"
	"strings"
)

type vec3 [3]int

type beaconScanner struct {
	id, orientation int
	pos             vec3
	beacons         []vec3
}

func readInput(fileName string) []beaconScanner {
	var arr []beaconScanner
	ReadInput(fileName, func(ioScanner *bufio.Scanner) {
		var scanner beaconScanner
		for ioScanner.Scan() {
			line := ioScanner.Text()
			if line == "" {
				arr = append(arr, scanner)
			} else if line[1] == '-' {
				id := line[len("--- scanner ") : len(line)-len(" ---")]
				scanner = beaconScanner{id: Atoi(id), orientation: -1}
			} else {
				split := strings.Split(line, ",")
				un := Atoi(split[0])
				deux := Atoi(split[1])
				trois := Atoi(split[2])
				scanner.beacons = append(scanner.beacons, vec3{un, deux, trois})
			}
		}
	})
	return arr
}

var orientations = [][]int{
	{1, 1, 1},
	{-1, 1, 1},
	{1, -1, 1},
	{-1, -1, 1},
	{1, 1, -1},
	{-1, 1, -1},
	{1, -1, -1},
	{-1, -1, -1},
}

func compare(scanner beaconScanner, scanner2 beaconScanner) bool {
	return false
}

func orientScanners(scanners []beaconScanner) []beaconScanner {
	known := hs.New(0)
	scanners[0].pos = vec3{0, 0, 0}
	for known.Size() != len(scanners) {
		iter := hs.New()
		for _, k := range known.Values() {
			for id, scanner := range scanners {
				if !known.Contains(id) && compare(scanners[k.(int)], scanner) {
					iter.Add(id)
				}
			}
		}
		known.Add(iter.Values()...)
	}
	return scanners
}

func ezMode(input []beaconScanner, ch chan<- int) {
	ch <- 0
}

func hardMode(input []beaconScanner, ch chan<- int) {
	ch <- 0
}

func Go(fileName string, ch chan string) {
	input := orientScanners(readInput(fileName))
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Beacon Scanner")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
