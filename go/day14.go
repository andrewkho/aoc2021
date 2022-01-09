package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
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
		log.Panicln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var init string
	splits := make(map[string]string)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if i == 0 {
			init = line
		} else if i == 1 {
			continue
		} else {
			pair := strings.Split(line, " -> ")
			splits[pair[0]] = pair[1]
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	step := func(in map[string]int, counts map[string]int) map[string]int {
		out := make(map[string]int)
		for k, v := range in {
			if m, ok := splits[k]; ok {
				out[string(k[0])+m] += v
				out[m+string(k[1])] += v
				counts[m] += v
			} else {
				out[k] += v
			}
		}
		return out
	}

	t1 := time.Now()
	counts := make(map[string]int)
	pairs := make(map[string]int)
	for i := 0; i<len(init)-1; i++ {
		counts[string(init[i])]++
		pairs[init[i:i+2]]++
	}
	counts[string(init[len(init)-1])]++
	steps := 10
	for i := 0; i<steps; i++ {
		pairs = step(pairs, counts)
	}
	min, max := 1<<63-1, 0
	for _, v := range counts {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	fmt.Printf("1: %v, %v\n", max-min, time.Since(t1))

	t2 := time.Now()
	for i := steps; i<40; i++ {
		pairs = step(pairs, counts)
	}
	min, max = 1<<63-1, 0
	for _, v := range counts {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	fmt.Printf("2: %v, %v\n", max-min, time.Since(t2))
}
