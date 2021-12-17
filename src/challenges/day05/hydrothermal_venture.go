package hydrothermalVenture

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"strings"
)

type seafloor [][]int
type line struct {
	start, end point
}
type point struct {
	x, y int
}

func readInput(fileName string) []line {
	var arr []line
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			in := scanner.Text()
			split := strings.Split(in, " -> ")
			start := parsePoint(split[0])
			end := parsePoint(split[1])
			arr = append(arr, line{start, end})
		}
	})
	return arr
}

func parsePoint(input string) point {
	split := strings.Split(input, ",")
	x := Atoi(split[0])
	y := Atoi(split[1])
	return point{x, y}
}

func removeDiagonalLines(lines []line) []line {
	var notDiag []line
	for _, l := range lines {
		if l.start.x == l.end.x || l.start.y == l.end.y {
			notDiag = append(notDiag, l)
		}
	}
	return notDiag
}

func getDims(lines []line) (int, int) {
	var maxX, maxY int
	// Get dims, our space is from (0,0) -> (n,m)
	for _, l := range lines {
		if l.start.x > maxX {
			maxX = l.start.x
		}
		if l.start.y > maxY {
			maxY = l.start.y
		}
		if l.end.x > maxX {
			maxX = l.end.x
		}
		if l.end.y > maxY {
			maxY = l.end.y
		}
	}
	// zero index
	maxX++
	maxY++
	return maxX, maxY
}

func createSeafloor(maxX, maxY int) seafloor {
	floor := make(seafloor, maxX)
	for i := 0; i < len(floor); i++ {
		floor[i] = make([]int, maxY)
	}
	return floor
}

func orderPointsIncreasing(start point, end point) (point, point) {
	if start.y == end.y && start.x > end.x {
		return end, start // vertical
	}
	if start.x == end.x && start.y < end.y {
		return end, start // horizontal
	}
	if start.x > end.x {
		return end, start // order by x coord increasing
	}
	return start, end
}

func gradient(start point, end point) int {
	if start.x == end.x || start.y == end.y {
		return -1 // vertical and horizontal lines are treated as -ve
	}
	return (end.y - start.y) / (end.x - start.x)
}

func countDangerous(floor seafloor) int {
	var count int
	for _, row := range floor {
		for _, point := range row {
			if point > 1 {
				count++
			}
		}
	}
	return count
}

func plotLine(l line, floor seafloor) {
	start, end := orderPointsIncreasing(l.start, l.end)
	gradient := gradient(start, end)
	vx, vy := start.x, start.y
	var dy int
	if gradient >= 0 {
		dy = 1
	} else {
		dy = -1
	}
	for true {
		floor[vx][vy] += 1
		if vx == end.x && vy == end.y {
			break
		}
		if vx != end.x {
			vx++ // dx is always 1
		}
		if vy != end.y {
			vy += dy
		}
	}
}

func ezMode(input []line, floor seafloor, ch chan<- int) {
	for _, l := range input {
		plotLine(l, floor)
	}
	ch <- countDangerous(floor)
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	maxX, maxY := getDims(input)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(removeDiagonalLines(input), createSeafloor(maxX, maxY), ezChan)
	go ezMode(input, createSeafloor(maxX, maxY), hardChan)

	ch <- fmt.Sprintln("Hydrothermal Venture")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
