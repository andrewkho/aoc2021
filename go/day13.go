package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
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
	re := regexp.MustCompile(`^fold along (x|y)=([0-9].*)$`)
	var dots []util.Point2D
	var instrs [][]string
	pastBreak := false
	var N, M int
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			pastBreak = true
			continue
		}
		if !pastBreak {
			dot := util.GetInts(line, ",")
			p := util.Point2D{Y: dot[1], X: dot[0]}
			if p.Y+1 > N {
				N = p.Y + 1
			}
			if p.X+1 > M {
				M = p.X + 1
			}
			dots = append(dots, p)
		} else {
			matches := re.FindStringSubmatch(line)
			instrs = append(instrs, matches[1:])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	array0 := util.New2DZeros(N, M)
	for _, dot := range dots {
		array0[dot.Y][dot.X] = 1
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	t := 0
	b := fold(array0, instrs[0])
	for _, row := range b {
		for _, v := range row {
			if v > 0 {
				t++
			}
		}
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	for _, instr := range instrs[1:] {
		b = fold(b, instr)
	}
	for _, row := range b {
		for _, v := range row {
			if v > 0 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}

func fold(a util.Array2D, instr []string) util.Array2D {
	dir := instr[0]
	loc, err := strconv.Atoi(instr[1])
	if err != nil {
		log.Panicln(err)
	}

	N, M := len(a), len(a[0])
	if dir == "x" {
		for i := 0; i < N; i++ {
			for dj := 0; dj < M-loc; dj++ {
				a[i][loc-dj] += a[i][loc+dj]
			}
			a[i] = a[i][:loc]
		}
	} else {
		for di := 0; di < N-loc; di++ {
			for j := 0; j < M; j++ {
				a[loc-di][j] += a[loc+di][j]
			}
		}
		a = a[:loc]
	}
	return a
}
