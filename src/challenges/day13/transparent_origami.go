package transparentOrigami

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type page struct {
	limX, limY int
	dots       [][]bool
}

func readInput(fileName string) (page, []int) {
	var p page
	var instructions []int
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			split := strings.Split(line, ",")
			x, err := strconv.Atoi(split[0])
			MaybePanic(err)
			y, err := strconv.Atoi(split[1])
			MaybePanic(err)
			maybeResize(x, y, &p)
			p.dots[y][x] = true
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
	})
	return p, instructions
}

func deepCopy(input page) page {
	inputCpy := page{input.limX, input.limY, nil}
	for _, dots := range input.dots {
		inputCpy.dots = append(inputCpy.dots, []bool{})
		for _, dot := range dots {
			inputCpy.dots[len(inputCpy.dots)-1] = append(inputCpy.dots[len(inputCpy.dots)-1], dot)
		}
	}
	return inputCpy
}

func maybeResize(x int, y int, p *page) {
	if x+1 > p.limX {
		p.limX = x + 1
	}
	if y+1 > p.limY {
		p.limY = y + 1
	}
	if len(p.dots) < p.limY {
		for j := len(p.dots); j < p.limY; j++ {
			p.dots = append(p.dots, []bool{})
		}
	}
	for i := 0; i < len(p.dots); i++ {
		if len(p.dots[i]) < p.limX {
			for j := len((p.dots)[i]); j < p.limX; j++ {
				p.dots[i] = append(p.dots[i], false)
			}
		}
	}
}

func foldLeft(p *page, foldLine int) {
	for y := 0; y < p.limY; y++ {
		for x := 1; x <= foldLine; x++ {
			if InBounds(foldLine+x, len(p.dots[y])) {
				p.dots[y][foldLine-x] = p.dots[y][foldLine-x] || p.dots[y][foldLine+x]
			}
		}
	}
	p.limX = foldLine
}

func foldUp(p *page, foldLine int) {
	for y := 1; y <= foldLine; y++ {
		for x := 0; x < p.limX; x++ {
			if InBounds(foldLine+y, len(p.dots)) {
				p.dots[foldLine-y][x] = p.dots[foldLine-y][x] || p.dots[foldLine+y][x]
			}
		}
	}
	p.limY = foldLine
}

func countDots(p page) int {
	count := 0
	for y := 0; y < p.limY; y++ {
		for x := 0; x < p.limX; x++ {
			if p.dots[y][x] {
				count++
			}
		}
	}
	return count
}

func toString(p page) string {
	sb := strings.Builder{}
	for y := 0; y < p.limY; y++ {
		line := "  "
		for x := 0; x < p.limX; x++ {
			if p.dots[y][x] {
				line += "O"
			} else {
				line += " "
			}
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

func ezMode(p page, instructions []int, ch chan<- int) {
	for _, instruction := range instructions {
		if instruction < 0 { // store y instructions as -ve
			foldUp(&p, -instruction)
		} else {
			foldLeft(&p, instruction)
		}
	}
	ch <- countDots(p)
}

func hardMode(p page, instructions []int, ch chan<- string) {
	for _, instruction := range instructions {
		if instruction < 0 { // store y instructions as -ve
			foldUp(&p, -instruction)
		} else {
			foldLeft(&p, instruction)
		}
	}
	ch <- toString(p)
}

func Go(fileName string, ch chan string) {
	p, instructions := readInput(fileName)
	page2 := deepCopy(p)
	ezChan := make(chan int)
	hardChan := make(chan string)

	go ezMode(p, instructions[:1], ezChan)
	go hardMode(page2, instructions, hardChan)

	ch <- fmt.Sprintln("Transparent Origami")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: \n%s\n", <-hardChan)
	close(ch)
}
