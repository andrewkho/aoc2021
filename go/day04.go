package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"./util"
)

type Card struct {
	Board util.Array2D
	Marks [][]bool
	RowSum []int
	ColSum []int
}

func NewCard() Card {
	card := Card{}
	card.Board = util.New2DZeros(5, 5)

	card.Marks = make([][]bool, 5)
	for i := 0; i<5; i++ {
		card.Board[i] = make([]int, 5)
		card.Marks[i] = make([]bool, 5)
	}
	card.RowSum = make([]int, 5)
	card.ColSum = make([]int, 5)

	return card
}

func (card *Card) Reset() {
	for i := 0; i<5; i++ {
		card.RowSum[i] = 0
		for j := 0; j<5; j++ {
			card.Marks[i][j] = false 
			// This is done 5x because lazy
			card.ColSum[j] = 0 
		}
	}
}

func (card *Card) Mark(v int) (bool, bool) {
	for i := 0; i<5; i++ {
		for j := 0; j<5; j++ {
			if card.Board[i][j] == v {
				card.Marks[i][j] = true
				card.RowSum[i]++
				card.ColSum[j]++
				if card.RowSum[i] == 5 || card.ColSum[j] == 5 {
					return true, true
				}
				return true, false
			}
		}
	}
	return false, false
}
func (card *Card) UnmarkedSum() int {
	t := 0
	for i := 0; i<5; i++ {
		for j := 0; j<5; j++ {
			if !card.Marks[i][j] {
				t += card.Board[i][j]
			}
		}
	}
	return t
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
		log.Fatal(err)
	}
	defer file.Close()

	var numbers []int
	var cards []Card
	scanner := bufio.NewScanner(file)
	var curcard Card
	currow := 0
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if i == 0 {
			numbers = util.GetInts(line, ",")
		} else if strings.TrimSpace(line) == "" {
			curcard = NewCard()
			currow = 0
			cards = append(cards, curcard)
		} else {
			row := curcard.Board[currow]
			for ri := 0; ri<5; ri++ { 
				j := ri * 3
				k := j+2
				if k > len(line) {
					k = len(line)
				}
				row[ri], err = strconv.Atoi(strings.TrimSpace(line[j:k]))
				if err != nil {
					log.Panic(err)
				}
			}
			currow++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	winning_sum := -1

	Bingo:
		for _, num := range numbers {
			for _, card := range cards {
				if _, winner := card.Mark(num); winner {
					winning_sum = card.UnmarkedSum()*num
					break Bingo
				}
			}
		}
	fmt.Printf("1: %v, %v\n", winning_sum, time.Since(t1))

	t2 := time.Now()
	for _, card := range cards {
		card.Reset()
	}
	won := make([]bool, len(cards))
	remain := len(cards)
	winning_sum = -1
	var num int

	LastBingo:
		for _, num = range numbers {
			for i, card := range cards {
				if won[i] {
					continue
				}
				if _, winner := card.Mark(num); winner {
					won[i] = true
					remain--
					if remain == 0 {
						winning_sum = card.UnmarkedSum()*num
						break LastBingo
					}
				}
			}
		}
	fmt.Printf("2: %v, %v, %v\n", winning_sum, num, time.Since(t2))
}
