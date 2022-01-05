package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

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

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	points := func(c rune) int {
		switch c {
		case ')': return 3
		case ']': return 57
		case '}': return 1197
		case '>': return 25137
		default: return 0
		}
	}

	opening := func(c rune) rune {
		switch c {
		case ')': return '('
		case ']': return '['
		case '}': return '{'
		case '>': return '<'
		default: return rune(0)
		}
	}

	points2 := func(c rune) int {
		switch c {
		case '(': return 1
		case '[': return 2
		case '{': return 3
		case '<': return 4
		default: return 0
		}
	}

	t1 := time.Now()
	t := 0
	var p2_scores []int
Part1:
	for _, line := range lines {
		var stack []rune
		for _, c := range line {
			if strings.ContainsRune("([{<", c) {
				stack = append(stack, c)
			} else {
				o := opening(c)
				n := len(stack)
				if stack[n-1] != o {
					t += points(c)
					continue Part1
				} else {
					stack = stack[:n-1]
				}
			}
		}

		sum := 0
		for j := len(stack) - 1; j >= 0; j-- {
			v := points2(stack[j])
			sum = sum*5 + v
		}
		p2_scores = append(p2_scores, sum)
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	sort.Ints(p2_scores)
	t = p2_scores[len(p2_scores)/2]
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
