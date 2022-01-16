package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
)

func step(lines []string) ([]string, int) {
	other := make([][]rune, len(lines))
	moves := 0
	for i, line := range lines {
		other[i] = make([]rune, len(line))
		for j, c := range line {
			prev := j-1
			if prev < 0 {
				prev = len(line)-1
			}
			next := j+1
			if next >= len(line) {
				next = 0
			}
			if line[prev] == '>' && line[j] == '.' {
				other[i][j] = '>'
				moves++
			} else if line[j] == '>' && line[next] == '.' {
				other[i][j] = '.'
			} else {
				other[i][j] = c
			}
		}
	}

	output := make([]string, len(lines))
	for i, line := range other {
		runes := make([]rune, len(line))
		prev := i-1
		if prev < 0 {
			prev = len(lines)-1
		}
		next := i+1
		if next >= len(lines) {
			next = 0
		}
		for j, c := range line {
			if other[prev][j] == 'v' && other[i][j] == '.' {
				runes[j] = 'v'
				moves++
			} else if other[i][j] == 'v' && other[next][j] == '.' {
				runes[j] = '.'
			} else {
				runes[j] = c
			}
		}
		output[i] = string(runes)
	}
	return output, moves
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
	var lines []string
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	t := 0
	moves := 1
	for moves > 0 {
		lines, moves = step(lines)
		t++
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}