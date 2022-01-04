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
	ages := make(map[int]int)
	for scanner.Scan() {
		ints := util.GetInts(scanner.Text(), ",")
		for _, v := range ints {
			ages[v]++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	step := func() {
		for k := 0; k < 9; k++ {
			ages[k-1] = ages[k]
		}
		ages[8] = ages[-1]
		ages[6] += ages[-1]
		ages[-1] = 0
	}

	fmt.Printf(" took %v\n", time.Since(t0))
	t1 := time.Now()
	for i := 0; i < 80; i++ {
		step()
	}
	t := 0
	for _, v := range ages {
		t += v
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	for i := 80; i < 256; i++ {
		step()
	}
	t = 0
	for _, v := range ages {
		t += v
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
