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
)

type Card struct {
	board [][]int
	marks [][]bool
	rowsum []int
	colsum []int
}

func NewCard() Card {
	card := Card{}
	card.board = make([][]int, 5)
	card.marks = make([][]bool, 5)
	for i := 0; i<5; i++ {
		card.board[i] = make([]int, 5)
		card.marks[i] = make([]bool, 5)
	}
	card.rowsum = make([]int, 5)
	card.colsum = make([]int, 5)

	return card
}

func (card *Card) reset() {
	for i := 0; i<5; i++ {
		card.rowsum[i] = 0
		for j := 0; j<5; j++ {
			card.marks[i][j] = false 
			// This is done 5x because lazy
			card.colsum[j] = 0 
		}
	}
}

func (card *Card) mark(v int) (bool, bool) {
	for i := 0; i<5; i++ {
		for j := 0; j<5; j++ {
			if card.board[i][j] == v {
				card.marks[i][j] = true
				card.rowsum[i]++
				card.colsum[j]++
				if card.rowsum[i] == 5 || card.colsum[j] == 5 {
					return true, true
				}
				return true, false
			}
		}
	}
	return false, false
}
func (card *Card) unmarked_sum() int {
	t := 0
	for i := 0; i<5; i++ {
		for j := 0; j<5; j++ {
			if !card.marks[i][j] {
				t += card.board[i][j]
			}
		}
	}
	return t
}

func get_ints(line string, sep string) []int {
	vals := strings.Split(line, sep)
	ints := make([]int, len(vals))
	for i, val := range vals {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Panic(err)
		}
		ints[i] = v
	}
	return ints
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
			numbers = get_ints(line, ",")
		} else if strings.TrimSpace(line) == "" {
			curcard = NewCard()
			currow = 0
			cards = append(cards, curcard)
		} else {
			row := curcard.board[currow]
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
	for _, num := range numbers {
		for _, card := range cards {
			if _, winner := card.mark(num); winner {
				winning_sum = card.unmarked_sum()*num
				break
			}
		}
		if winning_sum >= 0 {
			break
		}
	}
	fmt.Printf("1: %v, %v\n", winning_sum, time.Since(t1))

	t2 := time.Now()
	for _, card := range cards {
		card.reset()
	}
	won := make([]bool, len(cards))
	remain := len(cards)
	winning_sum = -1
	var num int
	for _, num = range numbers {
		for i, card := range cards {
			if won[i] {
				continue
			}
			if _, winner := card.mark(num); winner {
				won[i] = true
				remain--
				if remain == 0 {
					winning_sum = card.unmarked_sum()*num
					continue
				}
			}
		}
		if winning_sum >= 0 {
			break
		}
	}
	fmt.Printf("2: %v, %v, %v\n", winning_sum, num, time.Since(t2))
}