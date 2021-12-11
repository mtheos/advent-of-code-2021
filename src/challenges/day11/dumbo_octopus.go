package dumboOctopus

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"os"
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
		row := scanner.Text()
		arr = append(arr, []int{})
		for _, spot := range row {
			arr[len(arr)-1] = append(arr[len(arr)-1], int(spot-'0'))
		}
	}
	return arr
}

func deepCopy(input [][]int) [][]int {
	var inputCpy [][]int
	for _, row := range input {
		inputCpy = append(inputCpy, []int{})
		for _, cell := range row {
			inputCpy[len(inputCpy)-1] = append(inputCpy[len(inputCpy)-1], cell)
		}
	}
	return inputCpy
}

func energise(input [][]int) {
	for x := range input {
		for y := range input[x] {
			input[x][y]++
		}
	}
}

func overflowEnergy(x, y int, input [][]int) {
	var stepLUT = [][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
	for _, step := range stepLUT {
		dx, dy := step[0], step[1]
		if InBounds(x+dx, len(input)) && InBounds(y+dy, len(input[0])) {
			if input[x+dx][y+dy] != 0 {
				input[x+dx][y+dy]++
			}
		}
	}
}

func flashMob(input [][]int) int {
	flashes := 0
	for {
		changed := false
		for x := range input {
			for y := range input[x] {
				if input[x][y] > 9 {
					flash(x, y, input)
					flashes++
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	return flashes
}

func flash(x, y int, input [][]int) {
	input[x][y] = 0
	overflowEnergy(x, y, input)
}

func ezMode(input [][]int, steps int, ch chan<- int) {
	flashes := 0
	for ; steps > 0; steps-- {
		energise(input)
		flashes += flashMob(input)
	}
	ch <- flashes
}

func hardMode(input [][]int, ch chan<- int) {
	steps := 1
	octopi := 0
	for _, row := range input {
		octopi += len(row)
	}
	for ; ; steps++ {
		energise(input)
		flashes := flashMob(input)
		if flashes == octopi {
			break
		}
	}
	ch <- steps
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	inputCpy := deepCopy(input)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, 100, ezChan)
	go hardMode(inputCpy, hardChan)

	ch <- fmt.Sprintln("Dumbo Octopus")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
