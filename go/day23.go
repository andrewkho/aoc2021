package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"./util"
)

type Pos struct {
	col int
	row int
}

type Board struct {
	cols [11][5]rune
	score int
	depth int
}

func (board *Board) winner() bool {
	for i := 2; i<9; i+= 2 {
		exp := rune('A' + i/2 - 1)
		for _, c := range board.cols[i][1:board.depth] {
			if c != exp {
				return false
			}
		}
	}

	return true
}

func (board *Board) move(from Pos, to Pos, undo int) {
	c := board.cols[from.col][from.row]
	t := board.cols[to.col][to.row]
	if t != '.' || !strings.ContainsRune("ABCD", c) {
		log.Panicln(board, from, to)
	}
	board.cols[to.col][to.row] = c
	board.cols[from.col][from.row] = t
	energy := int(math.Pow(10, float64(c - 'A')))
	board.score += undo*energy*(from.row + to.row + util.Abs(from.col - to.col))
}

func (board *Board) getPosFromCol(col int, strt int) Pos {
	row := 0
	if col == 2 || col == 4 || col == 6 || col == 8 {
		for ; row < board.depth; row++ {
			if board.cols[col][row] >= 'A' {
				row -= (1-strt)
				break
			}
		}
	} 	
	if row == board.depth {
		row--
	}
	return Pos{col, row}
}

func (board *Board) moveIsFinal(start Pos, end Pos) bool {
	if end.row == 0 {
		return false
	}

	c := board.cols[start.col][start.row]
	if end.col != board.runeToCol(c) {
		return false
	}

	for j := end.row+1; j<board.depth; j++ {
		if board.cols[end.col][j] != c {
			return false
		}
	}
	return true
}

func (board *Board) runeToCol(c rune) int {
	return 2 + 2*int(c - 'A')
}

func (board *Board) moveIsValid(start Pos, end Pos) bool{
	if start.row == 0 && end.row == 0 {
		return false
	}
	if end.col == 2 || end.col == 4 || end.col == 6 || end.col == 8 {
		return false
	}

	return true
}

func (board *Board) getValidMoves() (valid [][2]Pos) {
	for starti := 0; starti<11; starti++ {
		start := board.getPosFromCol(starti, 1)
		c := board.cols[start.col][start.row]
		if c == '.' {
			continue
		}
		if start.row > 0 && start.col == board.runeToCol(c) {
			ok := true
			for j := start.row+1; j<board.depth; j++ {
				if board.cols[start.col][j] != c {
					ok = false
					break
				}
			}
			if ok {
				continue
			}
		}

		startValids := [][2]Pos{}
		for endi := starti + 1; endi < 11; endi++ {
			end := board.getPosFromCol(endi, 0)
			if board.cols[end.col][end.row] != '.' {
				break
			}
			if board.moveIsFinal(start, end) {
				return [][2]Pos{{start, end}}
			} else if board.moveIsValid(start, end) {
				startValids = append(startValids, [2]Pos{start, end})
			}
		}
		for endi := starti - 1; endi >= 0; endi-- {
			end := board.getPosFromCol(endi, 0)
			if board.cols[end.col][end.row] != '.' {
				break
			}
			if board.moveIsFinal(start, end) {
				return [][2]Pos{{start, end}}
			} else if board.moveIsValid(start, end) {
				startValids = append(startValids, [2]Pos{start, end})
			}
		}
		for _, v := range startValids {
			valid = append(valid, v)
		}
	}
	return valid
}

func solve(board Board) int {
	best := 1<<63-1
	seen := make(map[[11][5]rune]int)

	var recurse func()
	recurse = func() {
		if v, ok := seen[board.cols]; ok {
			if v <= board.score {
				return
			}
		}
		seen[board.cols] = board.score
		if board.winner() {
			if board.score < best {
				best = board.score
			}
		}
		validMoves := board.getValidMoves()
		for _, move := range validMoves {
			board.move(move[0], move[1], 1)
			recurse()
			board.move(move[1], move[0], -1)
		}
	}

	recurse()

	return best
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(string(debug.Stack()))
		}
	}()

	var filename string = os.Args[1]
	t0 := time.Now()
	fmt.Printf("Reading from %s,", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var boardP1, boardP2 Board
	boardP1.depth = 3
	boardP2.depth = 5
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		for col, c := range line[1:] {
			if strings.ContainsRune("ABCD.", c) {
				boardP2.cols[col][i-1] = c
				if i < 3 {
					boardP1.cols[col][i-1] = c
				} else if i > 4 {
					boardP1.cols[col][i-3] = c
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	t := solve(boardP1)
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	t = solve(boardP2)
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}