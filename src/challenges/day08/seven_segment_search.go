package sevenSegmentSearch

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

var segmentLUT = []int{
	1: 2,
	7: 3,
	4: 4,
	2: 5,
	3: 5,
	5: 5,
	0: 6,
	6: 6,
	9: 6,
	8: 7,
}

type Display struct {
	inputs  [10]string
	outputs [4]string
}

func readInput(fileName string) []Display {
	file, err := os.Open(fileName)
	MaybePanic(err)
	defer func(file *os.File) {
		err := file.Close()
		MaybePanic(err)
	}(file)

	scanner := bufio.NewScanner(file)
	var arr []Display
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " | ")
		inputs := split[0]
		var display Display
		for i, in := range strings.Split(inputs, " ") {
			display.inputs[i] = in
		}
		outputs := split[1]
		for i, out := range strings.Split(outputs, " ") {
			display.outputs[i] = out
		}
		arr = append(arr, display)
	}
	return arr
}

func putFirstEmpty(arr *[3]string, input string) {
	for i := 0; i < len(arr); i++ {
		if arr[i] == "" {
			arr[i] = input
			return
		}
	}
	panic("panic")
}

func fn(inputs [10]string) ([4]string, [3]string, [3]string) {
	var length5 [3]string
	var length6 [3]string
	var known [4]string

	for _, input := range inputs {
		switch len(input) {
		case segmentLUT[1]:
			known[0] = input
		case segmentLUT[4]:
			known[1] = input
		case segmentLUT[7]:
			known[2] = input
		case segmentLUT[8]:
			known[3] = input
		default:
			if len(input) == 5 {
				putFirstEmpty(&length5, input)
			} else {
				putFirstEmpty(&length6, input)
			}
		}
	}
	return known, length5, length6
}

/*
                 0000    0000
     2  1    2       2  1    2
     2  1    2       2  1    2
         3333            3333
     5       5       5  4    5
     5       5       5  4    5
                         6666

 0000    0000    0000
1            2       2
1            2       2
 3333    3333    3333
     5  4            5
     5  4            5
 6666    6666    6666

 0000    0000    0000
1    2  1       1    2
1    2  1       1    2
         3333    3333
4    5  4    5       5
4    5  4    5       5
 6666    6666    6666
*/

func solveIt(display Display) int {
	var segments [7]string
	known, length5, length6 := fn(display.inputs)
	one, four, seven, eight := known[0], known[1], known[2], known[3]

	fillSegmentZero(seven, one, &segments)
	fillSegmentOne(length5, four, &segments)
	fillSegmentTwo(length6, one, &segments)
	fillSegmentFour(length5, one, &segments)

	zero := ""
	fillSegmentThree(zero, eight, &segments)
	fillSegmentFive(&segments)
	fillSegmentSix(&segments)

	return 0
}

func fillSegmentZero(seven string, one string, segments *[7]string) {
	for _, seg := range seven {
		segment := string(seg)
		if !strings.Contains(one, segment) {
			segments[0] = segment
			return
		}
	}
	panic("panic")
}

func fillSegmentOne(length5 [3]string, four string, segments *[7]string) {
	opts := make(map[string]int)
	for _, number := range length5 {
		for _, seg := range number {
			segment := string(seg)
			opts[segment]++
		}
	}
	for key, segCount := range opts {
		if segCount == 1 && strings.Contains(four, key) {
			segments[1] = key
			return
		}
	}
	panic("panic")
}

func fillSegmentTwo(length6 [3]string, one string, segments *[7]string) {
	for _, number := range length6 {
		for _, seg := range number {
			segment := string(seg)
			if strings.Contains(one, segment) {
				segments[2] = segment
				return
			}
		}
	}
	panic("panic")
}

func fillSegmentThree(zero string, eight string, segments *[7]string) {
	for _, seg := range eight {
		segment := string(seg)
		if !strings.Contains(zero, segment) {
			segments[3] = segment
			return
		}
	}
	panic("panic")
}

func fillSegmentFour(length5 [3]string, one string, segments *[7]string) {
	opts := make(map[string]int)
	for _, number := range length5 {
		for _, seg := range number {
			segment := string(seg)
			opts[segment]++
		}
	}
	for key, segCount := range opts {
		if segCount == 1 && !strings.Contains(one, key) {
			segments[4] = key
			return
		}
	}
	panic("panic")
}

func fillSegmentFive(segments *[7]string) {
	panic("panic")
}

func fillSegmentSix(segments *[7]string) {
	panic("panic")
}

func ezMode(input []Display, ch chan<- int) {
	count := 0
	for _, display := range input {
		for _, output := range display.outputs {
			switch len(output) {
			case segmentLUT[1]:
				fallthrough
			case segmentLUT[4]:
				fallthrough
			case segmentLUT[7]:
				fallthrough
			case segmentLUT[8]:
				count++
			default:
			}
		}
	}
	ch <- count
}

func hardMode(input []Display, ch chan<- int) {
	for _, display := range input {
		solveIt(display)
	}
	ch <- 0
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Seven Segment Search")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
