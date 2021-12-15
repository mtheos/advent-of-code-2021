package chiton

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func readInput(fileName string) [][]int {
	file, err := os.Open(fileName)
	MaybePanic(err)
	defer func(file *os.File) {
		err := file.Close()
		MaybePanic(err)
	}(file)

	scanner := bufio.NewScanner(file)
	var arr [][]int
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "")
		arr = append(arr, []int{})
		for _, s := range split {
			v, err := strconv.Atoi(s)
			MaybePanic(err)
			arr[len(arr)-1] = append(arr[len(arr)-1], v)
		}
	}
	return arr
}

func theRealCave(cave [][]int) [][]int {
	realCave := make([][]int, len(cave)*5)
	for i := range realCave {
		realCave[i] = make([]int, len(cave[0])*5)
	}
	for y := range realCave {
		for x := range realCave[y] {
			var v int
			if y < len(cave) && x < len(cave[0]) {
				v = cave[y][x]
			} else if y < len(cave) {
				v = realCave[y][x-len(cave[0])]%9 + 1
			} else { // either can handle x & y > len(cave)
				v = realCave[y-len(cave)][x]%9 + 1
			}
			realCave[y][x] = v
		}
	}
	return realCave
}

func initCosts(dimY, dimX int) [][]int {
	costs := make([][]int, dimY)
	for i := range costs {
		costs[i] = make([]int, dimX)
	}

	for y := 0; y < dimY; y++ {
		for x := 0; x < dimX; x++ {
			costs[y][x] = -1
		}
	}
	return costs
}

func lowestCostTo(cave [][]int, costs [][]int, x int, y int) int {
	if costs[y][x] != -1 {
		return costs[y][x]
	} else if x == 0 && y == 0 {
		return 0
	} else if x == 0 {
		return lowestCostTo(cave, costs, x, y-1) + cave[y][x]
	} else if y == 0 {
		return lowestCostTo(cave, costs, x-1, y) + cave[y][x]
	}
	return int(math.Min(float64(lowestCostTo(cave, costs, x-1, y)), float64(lowestCostTo(cave, costs, x, y-1)))) + cave[y][x]
}

func ezMode(cave [][]int, ch chan<- int) {
	costs := initCosts(len(cave), len(cave[0]))
	for y := range cave {
		for x := range cave[y] {
			if costs[y][x] == -1 {
				costs[y][x] = lowestCostTo(cave, costs, x, y)
			} else {
				panic("interesting")
			}
		}
	}
	ch <- costs[len(cave)-1][len(cave[0])-1]
}

func hardMode(cave [][]int, ch chan<- int) {
	costs := initCosts(len(cave), len(cave[0]))
	for y := range cave {
		for x := range cave[y] {
			if costs[y][x] == -1 {
				costs[y][x] = lowestCostTo(cave, costs, x, y)
			} else {
				panic("interesting")
			}
		}
	}
	ch <- costs[len(cave)-1][len(cave[0])-1]
}

func Go(fileName string, ch chan string) {
	cave := readInput(fileName)
	realCave := theRealCave(cave)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(cave, ezChan)
	go hardMode(realCave, hardChan)

	ch <- fmt.Sprintln("Chiton")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
