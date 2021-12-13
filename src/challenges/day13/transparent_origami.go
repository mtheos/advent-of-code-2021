package transparentOrigami

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Page struct {
	limX, limY int
	dots       [][]bool
}

func readInput(fileName string) (Page, []int) {
	file, err := os.Open(fileName)
	MaybePanic(err)
	defer func(file *os.File) {
		err := file.Close()
		MaybePanic(err)
	}(file)

	scanner := bufio.NewScanner(file)
	var page Page
	var instructions []int
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		split := strings.Split(scanner.Text(), ",")
		x, err := strconv.Atoi(split[0])
		MaybePanic(err)
		y, err := strconv.Atoi(split[1])
		MaybePanic(err)
		maybeResize(x, y, &page)
		page.dots[y][x] = true
	}
	for scanner.Scan() {
		line := scanner.Text()[len("fold along "):]
		split := strings.Split(line, "=")
		z, err := strconv.Atoi(split[1])
		MaybePanic(err)
		if split[0] == "y" {
			z = -z
		}
		instructions = append(instructions, z)
	}
	return page, instructions
}

func maybeResize(x int, y int, page *Page) {
	if x+1 > page.limX {
		page.limX = x + 1
	}
	if y+1 > page.limY {
		page.limY = y + 1
	}
	if len(page.dots) < page.limY {
		for j := len(page.dots); j < page.limY; j++ {
			page.dots = append(page.dots, []bool{})
		}
	}
	for i := 0; i < len(page.dots); i++ {
		if len(page.dots[i]) < page.limX {
			for j := len((page.dots)[i]); j < page.limX; j++ {
				page.dots[i] = append(page.dots[i], false)
			}
		}
	}
}

func foldLeft(page Page, foldLine int) {
}

func foldUp(page Page, foldLine int) {
	for y := 0; y < foldLine; y++ {
		for x := 0; x < page.limX; x++ {
			page.dots[y][x] = page.dots[y][x] || page.dots[page.limY-y-1][x]
		}
	}
}

func countDots(page Page) int {
	count := 0
	for _, row := range page.dots {
		for _, dot := range row {
			if dot {
				count++
			}
		}
	}
	return count
}

func ezMode(page Page, instructions []int, ch chan<- int) {
	for _, instruction := range instructions {
		if instruction < 0 { // store y instructions as -ve
			foldUp(page, -instruction)
		} else {
			foldLeft(page, instruction)
		}
	}
	ch <- countDots(page)
}

func hardMode(page Page, instructions []int, ch chan<- int) {
	ch <- 0
}

func Go(fileName string, ch chan string) {
	page, instructions := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(page, instructions[:1], ezChan)
	go hardMode(page, instructions, hardChan)

	ch <- fmt.Sprintln("Transparent Origami")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
