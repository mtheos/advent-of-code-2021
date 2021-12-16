package syntaxScoring

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"math"
	"sort"
)

var syntaxScoring = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var autoCompleteScoring = map[rune]int{
	'(': 1,
	'[': 2,
	'{': 3,
	'<': 4,
}

func readInput(fileName string) []string {
	var arr []string
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			arr = append(arr, line)
		}
	})
	return arr
}

func matchToken(token rune, stack *lls.Stack) bool {
	switch token {
	case '(':
		fallthrough
	case '[':
		fallthrough
	case '{':
		fallthrough
	case '<':
		stack.Push(token)
	case ')':
		fallthrough
	case ']':
		fallthrough
	case '}':
		fallthrough
	case '>':
		val, _ := stack.Pop()
		diff := int(math.Abs(float64(val.(rune) - token)))
		if diff > 2 { // []{}<> have an ascii difference of 2, () has a difference of 1
			return true
		}
	}
	return false
}

func syntaxScore(line string, stack *lls.Stack) int {
	for _, token := range line {
		corrupted := matchToken(token, stack)
		if corrupted {
			val := syntaxScoring[token]
			if val == 0 {
				panic("panic")
			}
			return val
		}
	}
	return 0
}

func autoCompleteScore(stack *lls.Stack) (int, bool) {
	score := 0
	it := stack.Iterator()
	for it.Next() {
		token := it.Value().(rune)
		val := autoCompleteScoring[token]
		if val == 0 {
			panic("panic")
		}
		score = score*5 + val
	}
	return score, false
}

func ezMode(input []string, ch chan<- int) {
	score := 0
	stack := lls.New()
	for _, line := range input {
		stack.Clear()
		s := syntaxScore(line, stack)
		score += s
	}
	ch <- score
}

func hardMode(input []string, ch chan<- int) {
	var scores []int
	stack := lls.New()
	for _, line := range input {
		stack.Clear()
		s := syntaxScore(line, stack)
		if s == 0 && !stack.Empty() {
			s, _ := autoCompleteScore(stack)
			scores = append(scores, s)
		}
	}
	sort.Ints(scores)
	ch <- scores[len(scores)/2]
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(input, ezChan)
	go hardMode(input, hardChan)

	ch <- fmt.Sprintln("Syntax Scoring")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
