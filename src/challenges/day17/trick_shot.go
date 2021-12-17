package trickShot

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"math"
	"strings"
)

type targetArea struct {
	x0, x1, y0, y1 int
}

type projectile struct {
	x, y, vx, vy int
}

func readInput(fileName string) targetArea {
	var target targetArea
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		scanner.Scan()
		line := scanner.Text()[len("target area: "):]
		parts := strings.Split(line, ", ")
		xCoords := strings.Split(parts[0][len("x="):], "..")
		yCoords := strings.Split(parts[1][len("y="):], "..")
		target.x0 = Atoi(xCoords[0])
		target.x1 = Atoi(xCoords[1])
		target.y0 = Atoi(yCoords[0])
		target.y1 = Atoi(yCoords[1])
	})
	return target
}

func step(probe projectile) projectile {
	probe.x += probe.vx
	probe.y += probe.vy
	probe.vy -= 1
	if probe.vx > 0 {
		probe.vx -= 1
	} else if probe.vx < 0 {
		probe.vx += 1
	}
	return probe
}

func completelyMissed(probe projectile, target targetArea) bool {
	if probe.vx == 0 && probe.x < target.x0 {
		return true // never reach
	}
	return probe.x >= target.x1 || probe.y <= target.y0 // impossible to hit from here
}

func hitTarget(probe projectile, target targetArea) bool {
	return probe.x >= target.x0 && probe.x <= target.x1 &&
		probe.y >= target.y0 && probe.y <= target.y1
}

func ezMode(target targetArea, ch chan<- int) {
	// The probe will pass back through 0 with -ve initial velocity, so the velocity
	// that will give the greatest height is 1 less than the lowest depth to the target.
	y0 := float64(-target.y0 - 1)
	ch <- int(y0 * math.Ceil(y0/2))
}

// You could do this mathematically, but it's easier to just simulate it.
func hardMode(target targetArea, ch chan<- int) {
	var probe projectile
	viableTrajectories := 0
	maxY := 0
	for vx := 0; vx <= target.x1; vx++ {
		if (vx+1)*vx/2 < target.x0 {
			continue // never reaches the target
		}
		for vy := target.y0; vy < -target.y0; vy++ {
			probe = projectile{x: 0, y: 0, vx: vx, vy: vy}
			for !completelyMissed(probe, target) {
				probe = step(probe)
				if hitTarget(probe, target) {
					maxY = int(math.Min(float64(maxY), float64(vy)))
					viableTrajectories++
					break
				}
			}
		}
	}
	ch <- viableTrajectories
}

func Go(fileName string, ch chan string) {
	p := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(p, ezChan)
	go hardMode(p, hardChan)

	ch <- fmt.Sprintln("Trick Shot")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
