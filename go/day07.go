package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"./util"
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
	dists := make(map[int]int)
	mindist := 1<<62
	maxdist := 0
	for scanner.Scan() {
		ints := util.GetInts(scanner.Text(), ",")
		for _, i := range ints {
			dists[i]++
			if i < mindist {
				mindist = i
			}
			if i > maxdist {
				maxdist = i
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	t := 1 << 62
	for i := mindist; i<maxdist; i++ {
		c := 0
		for k, v := range dists {
			c += v*util.Abs(k-i)
		}
		if c < t {
			t = c
		} else {
			break
		}

	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	t = 1 << 62
	for i := mindist; i<maxdist; i++ {
		c := 0
		for k, v := range dists {
			d := util.Abs(k-i)
			c += v*d*(d+1)/2
		}
		if c < t {
			t = c
		} else {
			break
		}
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
