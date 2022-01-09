package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sort"
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
	var ints []int
	dists := make(map[int]int)
	mindist := 1<<63-1
	maxdist := 0
	for scanner.Scan() {
		ints = util.GetInts(scanner.Text(), ",")
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
	sort.Ints(ints)
	med := ints[len(ints)/2]
	t := 0
	for k, v := range dists {
		t += v*util.Abs(k - med)
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	cost := func(p int) int {
		c := 0
		for k, v := range dists {
			d := util.Abs(k-p)
			c += v*d*(d+1)/2
		}
		return c
	}
	srch := func() int {
		l := ints[0]
		r := ints[len(ints)-1]

		for l < r {
			m := l + (r-l)/2
			cl := cost(m-1)
			cm := cost(m)
			if cm < cl { // want to move left up
				l = m
			} else {
				r = m-1
			}
		}
		return l
	}

	t = cost(srch())
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
