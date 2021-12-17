package giantSquid

import (
	. "advent-of-code-2021/src/utils"
	"bufio"
	"fmt"
	"strings"
	"sync"
)

type board []row
type row []int

func deepCopy(boards []board) []board {
	var boardsCpy []board
	for _, gameBoard := range boards {
		boardCpy := make(board, 5)
		for i, gameRow := range gameBoard {
			boardCpy[i] = make(row, 5)
			for k, val := range gameRow {
				boardCpy[i][k] = val
			}
		}
		boardsCpy = append(boardsCpy, boardCpy)
	}
	return boardsCpy
}

func readInput(fileName string) []string {
	var arr []string
	ReadInput(fileName, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 0 { // drop empty lines
				arr = append(arr, strings.TrimSpace(line))
			}
		}
	})
	return arr
}

func parseMoves(input string) []int {
	split := strings.Split(input, ",")
	moves := make([]int, len(split))
	for i, val := range split {
		conv := Atoi(val)
		moves[i] = conv
	}
	return moves
}

func parseBoards(input []string) []board {
	var boards []board
	for i := 0; i < len(input); i += 5 {
		gameBoard := make(board, 5)
		for j := 0; j < 5; j++ {
			gameBoard[j] = make(row, 5)
			clean := strings.ReplaceAll(input[i+j], "  ", " ")
			split := strings.Split(clean, " ")
			for k, val := range split {
				conv := Atoi(val)
				gameBoard[j][k] = conv
			}
		}
		boards = append(boards, gameBoard)
	}
	return boards
}

func isRowWon(gameRow row, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, cell := range gameRow {
		if cell != -1 {
			ch <- false
			return
		}
	}
	ch <- true
}

func isColWon(gameBoard board, col int, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for gameRow := 0; gameRow < 5; gameRow++ {
		if gameBoard[gameRow][col] != -1 {
			ch <- false
			return
		}
	}
	ch <- true
}

func isWon(gameBoard board) bool {
	ch := make(chan bool, 10)
	var wg sync.WaitGroup

	for _, gameRow := range gameBoard {
		wg.Add(1)
		go isRowWon(gameRow, ch, &wg)
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go isColWon(gameBoard, i, ch, &wg)
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

func playRow(move int, gameRow row, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, cell := range gameRow {
		if cell == move {
			gameRow[i] = -1
		}
	}
}

func playBoard(move int, gameBoard board) bool {
	var wg sync.WaitGroup
	for _, gameRow := range gameBoard {
		wg.Add(1)
		go playRow(move, gameRow, &wg)
	}
	wg.Wait()
	return isWon(gameBoard)
}

func scoreBoard(lastMove int, gameBoard board) int {
	var sum int
	for _, gameRow := range gameBoard {
		for _, cell := range gameRow {
			if cell != -1 {
				sum += cell
			}
		}
	}
	return lastMove * sum
}

func ezMode(moves []int, boards []board, ch chan<- int) {
	var lastMove int
	var winningBoard board
out:
	for _, move := range moves {
		for _, gameBoard := range boards {
			isWinner := playBoard(move, gameBoard)
			if isWinner {
				lastMove, winningBoard = move, gameBoard
				break out
			}
		}
	}
	ch <- scoreBoard(lastMove, winningBoard)
}

func hardMode(moves []int, boards []board, ch chan<- int) {
	var lastMove int
	var lastWinner board
	for _, move := range moves {
		for i := 0; i < len(boards); i++ {
			gameBoard := boards[i]
			isWinner := playBoard(move, gameBoard)
			if isWinner {
				lastMove, lastWinner = move, gameBoard
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
