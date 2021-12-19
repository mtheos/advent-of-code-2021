package snailfish

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"math"
)

type snailNumber []int32

const (
	LBRACKET = 90 - '[' // -ve sentinel representation
	RBRACKET = 90 - ']'
)

func isBracket(c int32) bool {
	return c == LBRACKET || c == RBRACKET
}

func isNotBracket(c int32) bool {
	return !isBracket(c)
}

func append_(first snailNumber, second snailNumber) snailNumber {
	return append(append(append(snailNumber{LBRACKET}, first...), second...), RBRACKET)
}

func replace(sn snailNumber, start int, count int, replacements ...int32) snailNumber {
	return append(append(append(snailNumber{}, sn[:start]...), replacements...), sn[start+count:]...)
}

func readInput(fileName string) []snailNumber {
	var arr []snailNumber
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			arr = append(arr, parseNumber(scanner.Text()))
		}
	})
	return arr
}

func parseNumber(s string) snailNumber {
	var sn snailNumber
	for _, c := range s {
		var x int32
		switch c {
		case ',':
			continue
		case '[':
			x = LBRACKET
		case ']':
			x = RBRACKET
		default:
			x = c - '0'
		}
		sn = append(sn, x)
	}
	return sn
}

func seek(sn snailNumber, pos, step int) int {
	for ; pos >= 0 && pos < len(sn); pos += step {
		if isNotBracket(sn[pos]) {
			return pos
		}
	}
	return -1
}

func explode(sn snailNumber, pos int) snailNumber {
	left := seek(sn, pos-1, -1)
	right := seek(sn, pos+2, 1)
	if left != -1 {
		sn[left] += sn[pos]
	}
	if right != -1 {
		sn[right] += sn[pos+1]
	}
	return replace(sn, pos-1, 4, 0)
}

func split(sn snailNumber, pos int) snailNumber {
	n := float64(sn[pos]) / 2
	floor, ceil := int32(math.Floor(n)), int32(math.Ceil(n))
	return replace(sn, pos, 1, LBRACKET, floor, ceil, RBRACKET)
}

func reduceHelper(sn snailNumber) (snailNumber, bool) {
	depth := 0
	for pos := 0; pos < len(sn); pos++ {
		if sn[pos] == LBRACKET {
			depth++
		} else if sn[pos] == RBRACKET {
			depth--
		} else if depth == 5 {
			return explode(sn, pos), false
		}
	}
	for pos := 0; pos < len(sn); pos++ {
		if sn[pos] > 9 && isNotBracket(sn[pos]) {
			return split(sn, pos), false
		}
	}
	return sn, true
}

func reduce(sn snailNumber) snailNumber {
	for done := false; !done; sn, done = reduceHelper(sn) {
		// This line is left intentionally blank.
	}
	return sn
}

func magHelper(sn snailNumber) (snailNumber, bool) {
	for pos := 0; pos < len(sn)-1; pos++ {
		if isNotBracket(sn[pos]) && isNotBracket(sn[pos+1]) {
			s2 := append(snailNumber{}, sn[:pos-1]...)
			s2 = append(s2, sn[pos]*3+sn[pos+1]*2)
			s2 = append(s2, sn[pos+3:]...)
			return s2, false
		}
	}
	return sn, true
}

func magnitude(sn snailNumber, ch chan<- int) int {
	for done := false; !done; sn, done = magHelper(sn) {
		// This line is left intentionally blank.
	}
	if ch != nil {
		ch <- int(sn[0])
	}
	return int(sn[0])
}

func ezMode(input []snailNumber, ch chan<- int) {
	result := input[0]
	for _, sn := range input[1:] {
		result = reduce(append_(result, sn))
	}
	ch <- magnitude(result, nil)
}

func hardMode(input []snailNumber, ch chan<- int) {
	iter := 0
	maxMag := 0.0
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input); j++ {
			if i != j {
				iter++
				maxMag = math.Max(maxMag, float64(magnitude(reduce(append_(input[i], input[j])), nil)))
			}
		}
	}
	ch <- int(maxMag)
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Snailfish")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
