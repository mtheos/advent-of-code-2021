package trenchMap

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	console "github.com/atomicgo/cursor"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
)

type enhancerMap [512]uint8

type trenchMap struct {
	points   map[point]uint8
	infinity uint8
	A, B     point
}

type point struct {
	x, y int
}

const (
	dark  = 0
	light = 1
)

const fps = 30
const displayPadding = 3

func (trench trenchMap) getPixel(p point) uint8 {
	if p.x < trench.A.x || p.x > trench.B.x || p.y < trench.A.y || p.y > trench.B.y {
		return trench.infinity
	}
	if t, exists := trench.points[p]; exists && t == light {
		return light
	}
	return dark
}

func (trench *trenchMap) maybeAdjustBounds(p point) {
	trench.A.x = int(math.Min(float64(trench.A.x), float64(p.x)))
	trench.B.x = int(math.Max(float64(trench.B.x), float64(p.x)))
	trench.A.y = int(math.Min(float64(trench.A.y), float64(p.y)))
	trench.B.y = int(math.Max(float64(trench.B.y), float64(p.y)))
}

func (trench *trenchMap) matchBounds(old trenchMap) {
	trench.A = old.A
	trench.B = old.B
}

func (trench *trenchMap) setInfinity(old trenchMap, enhancer enhancerMap) {
	if old.infinity == dark {
		trench.infinity = enhancer[0]
	} else {
		trench.infinity = enhancer[(2<<8)-1]
	}
}

func clear() {
	cmd := exec.Command("cmd", "/c", "cls") // Windows console
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	MaybePanic(err)
}

func delimit(p point, m *point) uint8 {
	if m == nil {
		return ' '
	}
	if p.x == m.x-1 && p.y >= m.y-1 && p.y <= m.y+1 {
		return '['
	} else if p.x == m.x+2 && p.y >= m.y-1 && p.y <= m.y+1 {
		return ']'
	} else if p.x == m.x && p.y == m.y {
		return '['
	} else if p.x == m.x+1 && p.y == m.y {
		return ']'
	}
	return ' '
}

func buildLine(trench *trenchMap, y int, m *point, sb *strings.Builder) {
	if y < trench.A.y-displayPadding || y > trench.B.y+displayPadding {
		return
	}
	for x := trench.A.x - displayPadding; x <= trench.B.x+displayPadding; x++ {
		p := point{x, y}
		b := uint8('.')
		sb.WriteByte(delimit(p, m))
		if trench.getPixel(p) == light {
			b = '#'
		}
		sb.WriteByte(b)
	}
	sb.WriteByte(delimit(point{trench.B.x + displayPadding + 1, y}, m))
}

func visualEyes(oldTrench *trenchMap, newTrench *trenchMap, m *point) {
	sb := strings.Builder{}
	lines := 0
	for y := oldTrench.A.y - displayPadding; y <= oldTrench.B.y+displayPadding; y++ {
		buildLine(oldTrench, y, m, &sb)
		if newTrench != nil {
			sb.WriteString("  ")
			buildLine(newTrench, y, m, &sb)
		}
		sb.WriteByte('\n')
		lines++
	}
	console.HorizontalAbsolute(0)
	console.ClearLinesUp(lines + 1)
	fmt.Printf("%s\n", sb.String())
	tpf := time.Duration(1000 / fps)
	time.Sleep(1000 * 1000 * tpf)
}

func readInput(fileName string) (trenchMap, enhancerMap) {
	var enhancer enhancerMap
	trench := trenchMap{points: make(map[point]uint8)}
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		scanner.Scan()
		for i, c := range scanner.Text() {
			if c == '#' {
				enhancer[i] = light
			}
		}
		scanner.Scan() // blank line
		y := 0
		for scanner.Scan() {
			trench.B.y = y
			for x, c := range scanner.Text() {
				if c == '#' {
					trench.points[point{x, y}] = light
					trench.B.x = x
				}
			}
			y++
		}
	})
	trench.infinity = dark
	return trench, enhancer
}

func enhancePixel(trench trenchMap, p point, enhancer enhancerMap) uint8 {
	var image [9]uint8
	idx := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			image[idx] = trench.getPixel(point{p.x + j, p.y + i})
			idx++
		}
	}
	offset := 0
	for i, x := range image {
		offset += int(x) << (8 - i)
	}
	return enhancer[offset]
}

func enhance(trench trenchMap, enhancer enhancerMap) trenchMap {
	newTrench := trenchMap{points: make(map[point]uint8)}
	newTrench.matchBounds(trench)
	newTrench.setInfinity(trench, enhancer)
	for y := trench.A.y - 2; y <= trench.B.y+2; y++ {
		for x := trench.A.x - 2; x <= trench.B.x+2; x++ {
			p := point{x, y}
			//visualEyes(&trench, &newTrench, &p)
			if enhancePixel(trench, p, enhancer) == light {
				newTrench.points[p] = light
				newTrench.maybeAdjustBounds(p)
			}
		}
	}
	return newTrench
}

func ezMode(trench trenchMap, enhancer enhancerMap, count int, ch chan<- int) {
	//clear()
	trench.maybeAdjustBounds(trench.A)
	trench.maybeAdjustBounds(trench.B)
	for i := 0; i < count; i++ {
		trench = enhance(trench, enhancer)
	}
	if trench.infinity == light {
		panic("WTF")
	}
	ch <- len(trench.points)
}

func Go(fileName string, ch chan string) {
	trench, enhancer := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(trench, enhancer, 2, ezChan)
	go ezMode(trench, enhancer, 50, hardChan)

	ch <- fmt.Sprintln("Trench Map")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
