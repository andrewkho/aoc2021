package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	var filename string = os.Args[1]
	t0 := time.Now()
	fmt.Printf("Reading from %s,", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var depths []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if d, err := strconv.Atoi(scanner.Text()); err != nil {
			log.Fatal(err)
		} else {
			depths = append(depths, d)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" took %v\n", time.Since(t0))
	t1 := time.Now()
	increases := 0
	for i := 1; i < len(depths); i++ {
		if depths[i] > depths[i-1] {
			increases++
		}
	}
	fmt.Printf("1: %v, %v\n", increases, time.Since(t1))

	t2 := time.Now()
	increases = 0
	cur := depths[0] + depths[1] + depths[2]
	prev := cur
	for i := 3; i < len(depths); i++ {
		cur += depths[i] - depths[i-3] 
		if cur > prev {
			increases++
		}
		prev = cur
	}
	fmt.Printf("2: %v, %v\n", increases, time.Since(t2))
}
