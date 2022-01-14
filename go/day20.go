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
		log.Panicln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var code [512]int
	var img util.Array2D
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if i == 0 {
			for j, c := range line {
				if c == '#' {
					code[j] = 1
				} else {
					code[j] = 0
				}
			}
		} else if i == 1 {
			continue
		} else {
			row := make([]int, len(line))
			for j, c := range line {
				if c == '#' {
					row[j] = 1
				} else {
					row[j] = 0
				}
			}
			img = append(img, row)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	for i := 0; i<2; i++ {
		fill := 0
		if i%2 == 1 {
			fill = code[0]
		}
		img = enhance(img, code[:], fill)
	}

	t := 0
	for _, row := range img {
		for _, v := range row {
			t += v
		}
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	for i := 2; i<50; i++ {
		fill := 0
		if i%2 == 1 {
			fill = code[0]
		}
		img = enhance(img, code[:], fill)
	}

	t = 0
	for _, row := range img {
		for _, v := range row {
			t += v
		}
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}

func enhance(img util.Array2D, code []int, fill int) util.Array2D {
	N, M := len(img), len(img[0])
	out := util.New2DZeros(N+2, M+2)

	get := func(i int, j int) int {
		key, idx := 0, 0
		for dy := 1; dy >= -1; dy-- {
			for dx := 1; dx >= -1; dx-- {
				imgi, imgj := i-1+dy, j-1+dx
				var val int
				if imgi < 0 || imgi >= N || imgj < 0 || imgj >= M {
					val = fill
				} else {
					val = img[imgi][imgj]
				}
				key += val<<idx
				idx++
			}
		}
		return code[key]
	}

	for i := 0; i<N+2; i++ {
		for j := 0; j<M+2; j++ {
			out[i][j] = get(i, j)
		}
	}
	return out
}