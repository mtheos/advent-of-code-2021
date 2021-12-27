package smokeBasin

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"sort"
)

type point struct {
	x, y int
}

func readInput(fileName string) [][]int {
	var arr [][]int
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			row := scanner.Text()
			arr = append(arr, []int{})
			for _, spot := range row {
				arr[len(arr)-1] = append(arr[len(arr)-1], int(spot-'0'))
			}
		}
	})
	return arr
}

func placeAndSort(largestBasins *[3]int, size int) {
	if size > largestBasins[0] {
		largestBasins[0] = size
	}
	sort.Ints(largestBasins[:])
}

func sizeBasinHelper(x, y int, input [][]int, visited [][]bool) int {
	if !InBounds(x, len(input)) || !InBounds(y, len(input[x])) {
		return 0
	}
	if visited[x][y] {
		return 0
	}
	if input[x][y] == 9 {
		return 0
	}
	visited[x][y] = true
	return 1 + sizeBasinHelper(x+1, y, input, visited) + sizeBasinHelper(x-1, y, input, visited) + sizeBasinHelper(x, y+1, input, visited) + sizeBasinHelper(x, y-1, input, visited)
}

func sizeBasin(x, y int, input [][]int, visited [][]bool, resultChan chan int) {
	resultChan <- sizeBasinHelper(x, y, input, visited)
}

func isLowestPoint(x, y int, basin [][]int) bool {
	stepLUT := []int{-1, 1}
	for _, dx := range stepLUT {
		if InBounds(x+dx, len(basin)) && basin[x+dx][y] <= basin[x][y] {
			return false
		}
	}
	for _, dy := range stepLUT {
		if InBounds(y+dy, len(basin[x])) && basin[x][y+dy] <= basin[x][y] {
			return false
		}
	}
	return true
}

func ezMode(input [][]int, ch chan<- int) {
	var lowestPoints []int
	for x := range input {
		for y := range input[x] {
			if isLowestPoint(x, y, input) {
				lowestPoints = append(lowestPoints, input[x][y])
			}
		}
	}
	risk := 0
	for _, height := range lowestPoints {
		risk += height + 1
	}
	ch <- risk
}

func hardMode(input [][]int, ch chan<- int) {
	var visited [][]bool
	var lowestPoints []point
	for x := range input {
		visited = append(visited, []bool{})
		for y := range input[x] {
			visited[x] = append(visited[x], false)
			if isLowestPoint(x, y, input) {
				lowestPoints = append(lowestPoints, point{x, y})
			}
		}
	}

	var largestBasins = [3]int{}
	resultChan := make(chan int)
	results := 0
	for _, point := range lowestPoints {
		go sizeBasin(point.x, point.y, input, visited, resultChan)
		results++
	}
	for ; results > 0; results-- {
		placeAndSort(&largestBasins, <-resultChan)
	}
	size := 1
	for _, basinSize := range largestBasins {
		size *= basinSize
	}
	ch <- size
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Smoke Basin")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
