package sevenSegmentSearch

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"strings"
)

type display struct {
	inputs  [10]string
	outputs [4]string
}

const (
	T byte = iota
	TL
	TR
	M
	BL
	BR
	B
)

var segmentCountLUT = []int{
	0: 6,
	1: 2,
	2: 5,
	3: 5,
	4: 4,
	5: 5,
	6: 6,
	7: 3,
	8: 7,
	9: 6,
}

func readInput(fileName string) []display {
	var arr []display
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			split := strings.Split(line, " | ")
			inputs := split[0]
			var segmentDisplay display
			for i, in := range strings.Split(inputs, " ") {
				segmentDisplay.inputs[i] = in
			}
			outputs := split[1]
			for i, out := range strings.Split(outputs, " ") {
				segmentDisplay.outputs[i] = out
			}
			arr = append(arr, segmentDisplay)
		}
	})
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

func parseInput(inputs [10]string) ([4]string, [3]string, [3]string) {
	var length5 [3]string
	var length6 [3]string
	var known [4]string

	for _, input := range inputs {
		switch len(input) {
		case segmentCountLUT[1]:
			known[0] = input
		case segmentCountLUT[4]:
			known[1] = input
		case segmentCountLUT[7]:
			known[2] = input
		case segmentCountLUT[8]:
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

func solveIt(segmentDisplay display, resultChan chan int) {
	var segments [7]string
	known, length5, length6 := parseInput(segmentDisplay.inputs)
	one, four, seven, eight := known[0], known[1], known[2], known[3]

	fillSegmentZero(seven, one, &segments)
	fillSegmentsOneThreeFour(length5, four, &segments)
	fillSegmentsTwoAndFive(length6, one, &segments)
	fillSegmentSix(eight, &segments)
	LUT := createLUT(segments)
	result := 0
	multiplier := 1000
	for _, digit := range segmentDisplay.outputs {
		result += LUT[StrSort(digit)] * multiplier
		multiplier /= 10
	}
	resultChan <- result
}

func createLUT(segments [7]string) map[string]int {
	digits := make(map[string]int)
	digits[StrSort(StrCat(segments[T], segments[TL], segments[TR], segments[BL], segments[BR], segments[B]))] = 0
	digits[StrSort(StrCat(segments[TR], segments[BR]))] = 1
	digits[StrSort(StrCat(segments[T], segments[TR], segments[M], segments[BL], segments[B]))] = 2
	digits[StrSort(StrCat(segments[T], segments[TR], segments[M], segments[BR], segments[B]))] = 3
	digits[StrSort(StrCat(segments[TL], segments[TR], segments[M], segments[BR]))] = 4
	digits[StrSort(StrCat(segments[T], segments[TL], segments[M], segments[BR], segments[B]))] = 5
	digits[StrSort(StrCat(segments[T], segments[TL], segments[M], segments[BL], segments[BR], segments[B]))] = 6
	digits[StrSort(StrCat(segments[T], segments[TR], segments[BR]))] = 7
	digits[StrSort(StrCat(segments[T], segments[TL], segments[TR], segments[M], segments[BL], segments[BR], segments[B]))] = 8
	digits[StrSort(StrCat(segments[T], segments[TL], segments[TR], segments[M], segments[BR], segments[B]))] = 9
	return digits
}

func fillSegmentZero(seven string, one string, segments *[7]string) {
	for _, seg := range seven {
		segment := string(seg)
		if !strings.Contains(one, segment) {
			segments[T] = segment
			return
		}
	}
	if segments[T] == "" {
		panic("panic")
	}
}

func fillSegmentsOneThreeFour(length5 [3]string, four string, segments *[7]string) {
	opts := make(map[string]int)
	for _, number := range length5 {
		for _, seg := range number {
			segment := string(seg)
			opts[segment]++
		}
	}
	for key, segCount := range opts {
		if segCount == 1 && strings.Contains(four, key) {
			segments[TL] = key
		} else if segCount == 3 && strings.Contains(four, key) {
			segments[M] = key
		} else if segCount == 1 && !strings.Contains(four, key) {
			segments[BL] = key
		}
	}
	if segments[TL] == "" || segments[M] == "" || segments[BL] == "" {
		panic("panic")
	}
}

func fillSegmentsTwoAndFive(length6 [3]string, one string, segments *[7]string) {
	opts := make(map[string]int)
	for _, number := range length6 {
		for _, seg := range number {
			segment := string(seg)
			opts[segment]++
		}
	}
	for key, segCount := range opts {
		if segCount == 2 && strings.Contains(one, key) {
			segments[TR] = key
		} else if segCount == 3 && strings.Contains(one, key) {
			segments[BR] = key
		}
	}
	if segments[TR] == "" || segments[BR] == "" {
		panic("panic")
	}
}

func fillSegmentSix(eight string, segments *[7]string) {
	for _, segment := range segments {
		eight = strings.Replace(eight, segment, "", 1)
	}
	segments[B] = eight
	if segments[B] == "" {
		panic("panic")
	}
}

func ezMode(input []display, ch chan<- int) {
	count := 0
	for _, segmentDisplay := range input {
		for _, output := range segmentDisplay.outputs {
			switch len(output) {
			case segmentCountLUT[1]:
				fallthrough
			case segmentCountLUT[4]:
				fallthrough
			case segmentCountLUT[7]:
				fallthrough
			case segmentCountLUT[8]:
				count++
			default:
			}
		}
	}
	ch <- count
}

func hardMode(input []display, ch chan<- int) {
	resultChan := make(chan int)
	for _, segmentDisplay := range input {
		go solveIt(segmentDisplay, resultChan)
	}
	sum := 0
	for i := 0; i < len(input); i++ {
		sum += <-resultChan
	}
	ch <- sum
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
