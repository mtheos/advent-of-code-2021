package giantSquid

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Board []Row
type Row []int

func deepCopy(boards []Board) []Board {
	var boardsCpy []Board
	for _, board := range boards {
		boardCpy := make(Board, 5)
		for i, row := range board {
			boardCpy[i] = make(Row, 5)
			for k, val := range row {
				boardCpy[i][k] = val
			}
		}
		boardsCpy = append(boardsCpy, boardCpy)
	}
	return boardsCpy
}

func readInput(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var arr []string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 { // drop empty lines
			arr = append(arr, strings.TrimSpace(line))
		}
	}
	return arr
}

func parseMoves(input string) []int {
	split := strings.Split(input, ",")
	moves := make([]int, len(split))
	for i, val := range split {
		conv, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		moves[i] = conv
	}
	return moves
}

func parseBoards(input []string) []Board {
	var boards []Board
	for i := 0; i < len(input); i += 5 {
		board := make(Board, 5)
		for j := 0; j < 5; j++ {
			board[j] = make(Row, 5)
			clean := strings.ReplaceAll(input[i+j], "  ", " ")
			split := strings.Split(clean, " ")
			for k, val := range split {
				conv, err := strconv.Atoi(val)
				if err != nil {
					panic(err)
				}
				board[j][k] = conv
			}
		}
		boards = append(boards, board)
	}
	return boards
}

func isRowWon(row Row, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, cell := range row {
		if cell != -1 {
			ch <- false
			return
		}
	}
	ch <- true
}

func isColWon(board Board, col int, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for row := 0; row < 5; row++ {
		if board[row][col] != -1 {
			ch <- false
			return
		}
	}
	ch <- true
}

func isWon(board Board) bool {
	ch := make(chan bool, 10)
	var wg sync.WaitGroup

	for _, row := range board {
		wg.Add(1)
		go isRowWon(row, ch, &wg)
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go isColWon(board, i, ch, &wg)
	}

	wg.Wait()
	close(ch)
	for b := range ch {
		if b {
			return true
		}
	}
	return false
}

func playRow(move int, row Row, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, cell := range row {
		if cell == move {
			row[i] = -1
		}
	}
}

func playBoard(move int, board Board) bool {
	var wg sync.WaitGroup
	for _, row := range board {
		wg.Add(1)
		go playRow(move, row, &wg)
	}
	wg.Wait()
	return isWon(board)
}

func scoreBoard(lastMove int, board Board) int {
	var sum int
	for _, row := range board {
		for _, cell := range row {
			if cell != -1 {
				sum += cell
			}
		}
	}
	return lastMove * sum
}

func ezMode(moves []int, boards []Board, ch chan<- int) {
	var lastMove int
	var winningBoard Board
out:
	for _, move := range moves {
		for _, board := range boards {
			isWinner := playBoard(move, board)
			if isWinner {
				lastMove, winningBoard = move, board
				break out
			}
		}
	}
	ch <- scoreBoard(lastMove, winningBoard)
}

func hardMode(moves []int, boards []Board, ch chan<- int) {
	var lastMove int
	var lastWinner Board
	for _, move := range moves {
		for i := 0; i < len(boards); i++ {
			board := boards[i]
			isWinner := playBoard(move, board)
			if isWinner {
				lastMove, lastWinner = move, board
				boards = append(boards[:i], boards[i+1:]...)
				i--
			}
		}
	}
	ch <- scoreBoard(lastMove, lastWinner)
}

func Go(fileName string, ch chan string) {
	input := readInput(fileName)
	moves, boards := parseMoves(input[0]), parseBoards(input[1:])
	boardsCpy := deepCopy(boards)
	ezChan := make(chan int)
	hardChan := make(chan int)

	go ezMode(moves, boards, ezChan)
	go hardMode(moves, boardsCpy, hardChan)

	ch <- fmt.Sprintln("Giant Squid")
	ch <- fmt.Sprintf("  ezMode: %d\n", <-ezChan)
	ch <- fmt.Sprintf("  hardMode: %d\n", <-hardChan)
	close(ch)
}
