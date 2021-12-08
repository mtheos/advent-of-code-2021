package hydrothermalVenture

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Seafloor [][]int
type Line struct {
	start, end Point
}
type Point struct {
	x, y int
}

func readInput(fileName string) []Line {
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

	scanner := bufio.NewScanner(file)
	var arr []Line
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " -> ")
		start := parsePoint(split[0])
		end := parsePoint(split[1])
		arr = append(arr, Line{start, end})
	}
	return arr
}

func parsePoint(input string) Point {
	split := strings.Split(input, ",")
	x, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	return Point{x, y}
}

func removeDiagonalLines(lines []Line) []Line {
	var notDiag []Line
	for _, line := range lines {
		if line.start.x == line.end.x || line.start.y == line.end.y {
			notDiag = append(notDiag, line)
		}
	}
	return notDiag
}

func getDims(lines []Line) (int, int) {
	var maxX, maxY int
	// Get dims, our space is from (0,0) -> (n,m)
	for _, line := range lines {
		if line.start.x > maxX {
			maxX = line.start.x
		}
		if line.start.y > maxY {
			maxY = line.start.y
		}
		if line.end.x > maxX {
			maxX = line.end.x
		}
		if line.end.y > maxY {
			maxY = line.end.y
		}
	}
	// zero index
	maxX++
	maxY++
	return maxX, maxY
}

func createSeafloor(maxX, maxY int) Seafloor {
	seafloor := make(Seafloor, maxX)
	for i := 0; i < len(seafloor); i++ {
		seafloor[i] = make([]int, maxY)
	}
	return seafloor
}

func orderPointsIncreasing(start Point, end Point) (Point, Point) {
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

func gradient(start Point, end Point) int {
	if start.x == end.x || start.y == end.y {
		return -1 // vertical and horizontal lines are treated as -ve
	}
	return (end.y - start.y) / (end.x - start.x)
}

func countDangerous(seafloor Seafloor) int {
	var count int
	for _, row := range seafloor {
		for _, point := range row {
			if point > 1 {
				count++
			}
		}
	}
	return count
}

func plotLine(line Line, seafloor Seafloor) {
	start, end := orderPointsIncreasing(line.start, line.end)
	gradient := gradient(start, end)
	vx, vy := start.x, start.y
	var dy int
	if gradient >= 0 {
		dy = 1
	} else {
		dy = -1
	}
	for true {
		seafloor[vx][vy] += 1
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

func ezMode(input []Line, seafloor Seafloor, ch chan<- int) {
	for _, line := range input {
		plotLine(line, seafloor)
	}
	ch <- countDangerous(seafloor)
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
