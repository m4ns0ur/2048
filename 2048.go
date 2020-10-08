package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const boardLen = 4

var boardStartY = 0
var gameFieldEndY = 0

func main() {
	err := termbox.Init()
	panicError(err)
	defer termbox.Close()

	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Sync()
	board := initBoard(boardLen)
	drawGameField(board)
	startGame(board)
}

func drawGameField(board [][]int) {
	putNextNumber(board)
	putNextNumber(board)
	boardStartY = printTerminal(0, 0, []string{"Game 2048", ""})
	boardEndY := drawBoard(0, boardStartY, board)
	gameFieldEndY = printTerminal(0, boardEndY, []string{"Esc ←↑↓→", ""})
}

func initBoard(len int) [][]int {
	var board = make([][]int, len)
	for i := range board {
		board[i] = make([]int, len)
	}
	return board
}

func drawBoard(startX, startY int, board [][]int) int {
	strs := []string{}
	for _, row := range board {
		var str string
		for _, cell := range row {
			str += fmt.Sprintf("%4d", cell)
		}
		str = strings.Replace(str, "0", ".", -1)
		strs = append(strs, str, "")
	}

	printTerminal(0, boardStartY, strs)
	termbox.Flush()
	return startY + len(strs)
}

func putNextNumber(board [][]int) (loser bool) {
	emptyCells := findEmtpyCells(board)
	if len(emptyCells) <= 0 {
		gameOver()
		return true
	}
	rndSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSrc)
	emptyCell := emptyCells[rnd.Intn(len(emptyCells))]
	board[emptyCell/len(board)][emptyCell%len(board)] = 2
	return false
}

func findEmtpyCells(board [][]int) []int {
	var emptyCells = []int{}
	for i, row := range board {
		for j, cell := range row {
			if cell == 0 {
				emptyCells = append(emptyCells, i*len(board)+j)
			}
		}
	}
	return emptyCells
}

func startGame(board [][]int) {
	end := false
loop:
	for {
		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break loop
			case termbox.KeyArrowDown:
				if !end {
					board = rotateBoard(board, false)
					board = slideLeft(board)
					board = rotateBoard(board, true)
					end = checkAndRefreshBoard(board)
				}
			case termbox.KeyArrowLeft:
				if !end {
					board = slideLeft(board)
					end = checkAndRefreshBoard(board)
				}
			case termbox.KeyArrowRight:
				if !end {
					board = rotateBoard(board, true)
					board = rotateBoard(board, true)
					board = slideLeft(board)
					board = rotateBoard(board, false)
					board = rotateBoard(board, false)
					end = checkAndRefreshBoard(board)
				}
			case termbox.KeyArrowUp:
				if !end {
					board = rotateBoard(board, true)
					board = slideLeft(board)
					board = rotateBoard(board, false)
					end = checkAndRefreshBoard(board)
				}
			}
		case termbox.EventResize:
			if !end {
				drawBoard(0, boardStartY, board)
			}
		case termbox.EventError:
			panicError(event.Err)
		}
	}
}

func checkAndRefreshBoard(board [][]int) (end bool) {
	if checkWinner(board) {
		return true
	}
	if putNextNumber(board) {
		return true
	}
	drawBoard(0, boardStartY, board)
	return false
}

func slideLeft(board [][]int) [][]int {
	for _, row := range board {
		stopMerge := 0
		for j := 1; j < len(row); j++ {
			if row[j] != 0 {
				for k := j; k > stopMerge; k-- {
					if row[k-1] == 0 {
						row[k-1] = row[k]
						row[k] = 0
					} else if row[k-1] == row[k] {
						row[k-1] += row[k]
						row[k] = 0
						stopMerge = k
						break
					} else {
						break
					}
				}
			}
		}
	}
	return board
}

func checkWinner(board [][]int) (winner bool) {
	for _, row := range board {
		for _, cell := range row {
			if cell == 2048 {
				gameWin()
				return true
			}
		}
	}
	return false
}

func gameOver() {
	printTerminal(0, gameFieldEndY, []string{"Game Over!"})
	termbox.SetCursor(0, gameFieldEndY+1)
	termbox.Flush()
}

func gameWin() {
	printTerminal(0, gameFieldEndY, []string{"You Won!"})
	termbox.SetCursor(0, gameFieldEndY+1)
	termbox.Flush()
}

func panicError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func rotateBoard(board [][]int, counterClockWise bool) [][]int {
	var rotatedBoard = make([][]int, len(board))
	for i, row := range board {
		rotatedBoard[i] = make([]int, len(row))
		for j := range row {
			if counterClockWise {
				rotatedBoard[i][j] = board[j][len(board)-i-1]
			} else {
				rotatedBoard[i][j] = board[len(board)-j-1][i]
			}
		}
	}
	return rotatedBoard
}

func printTerminal(startX, startY int, strs []string) int {
	for y, str := range strs {
		for x, ch := range str {
			termbox.SetCell(startX+x, startY+y, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
	return startY + len(strs)
}
