package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
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

	t1 := time.Now()
	t := 0
	for i, line := range lines {
		for j, c := range line {
			bc := byte(c)
			if i > 0 && lines[i-1][j] <= bc {
				continue
			} else if i < len(lines) - 1 && lines[i+1][j] <= bc {
				continue
			} else if j > 0 && lines[i][j-1] <= bc {
				continue
			} else if j < len(line) - 1 && lines[i][j+1] <= bc {
				continue
			}
			v, err := strconv.Atoi(string(c))
			if err != nil {
				log.Panicln(err)
			}
			t += 1+v
		}
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	seen := make([][]bool, len(lines))
	for i := 0; i<len(lines); i++ {
		seen[i] = make([]bool, len(lines[0]))
	}
	var regions util.IntMinHeap
	for i := 0; i<len(lines); i++ {
		for j := 0; j<len(lines[0]); j++ {
			if seen[i][j] || lines[i][j] == byte('9') {
				continue
			}
			// do dfs
			seen[i][j] = true
			stack := []util.Point2D{{Y: i, X: j}}
			curReg := 0
			for len(stack) > 0 {
				p := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				curReg++

				for _, d := range [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
					p2 := util.Point2D{X: p.X + d[0], Y: p.Y + d[1]}
					if p2.X < 0 || p2.X >= len(seen[0]) || p2.Y < 0 || p2.Y >= len(seen) {
						continue
					}
					if seen[p2.Y][p2.X] || lines[p2.Y][p2.X] == byte('9') {
						continue
					}
					seen[p2.Y][p2.X] = true
					stack = append(stack, p2)
				}
			}
			heap.Push(&regions, curReg)
			if len(regions) > 3 {
				heap.Pop(&regions)
			}
		}
	}
	t = 1
	for i := 0; i<len(regions); i++ {
		t *= regions[i]
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
