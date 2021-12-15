package chiton

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	cost, x, y int
}

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

func initVisited(dimY, dimX int) [][]bool {
	costs := make([][]bool, dimY)
	for i := range costs {
		costs[i] = make([]bool, dimX)
	}
	for y := 0; y < dimY; y++ {
		for x := 0; x < dimX; x++ {
			costs[y][x] = false
		}
	}
	return costs
}

func linearTimeHeapQIChooseYou(q []Point) (Point, []Point) {
	lowest := q[0]
	slice := 0
	for i, p := range q {
		if p.cost < lowest.cost {
			lowest = p
			slice = i
		}
	}
	return lowest, append(q[:slice], q[slice+1:]...)
}

func ezMode(cave [][]int, ch chan<- int) {
	var stepLUT = [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}
	visited := initVisited(len(cave), len(cave[0]))
	q := make([]Point, 0)
	p := Point{0, 0, 0}
	q = append(q, p)
	for len(q) > 0 {
		p, q = linearTimeHeapQIChooseYou(q)
		if p.y == len(cave)-1 && p.x == len(cave[0])-1 {
			break
		}
		if visited[p.y][p.x] {
			continue
		}
		visited[p.y][p.x] = true
		for _, s := range stepLUT {
			x, y := s[0], s[1]
			if InBounds(p.y+y, len(cave)) && InBounds(p.x+x, len(cave[0])) {
				q = append(q, Point{p.cost + cave[p.y+y][p.x+x], p.x + x, p.y + y})
			}
		}
	}
	ch <- p.cost
}

func Go(fileName string, ch chan string) {
	cave := readInput(fileName)
	realCave := theRealCave(cave)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(cave, ezChan)
	go ezMode(realCave, hardChan)

	ch <- fmt.Sprintln("Chiton")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
