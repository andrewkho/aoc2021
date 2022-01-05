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
	var input util.Array2D
	for scanner.Scan() {
		row := util.GetInts(scanner.Text(), "")
		input = append(input, row)
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	step := func() int {
		// will this work w/dfs?
		var stack []util.Point2D
		for i := 0; i < len(input); i++ {
			for j := 0; j < len(input[i]); j++ {
				input[i][j]++
				if input[i][j] > 9 {
					stack = append(stack, util.Point2D{Y: i, X: j})
				}
			}
		}
		flashed := make([][]bool, len(input))
		for i := 0; i < len(input); i++ {
			flashed[i] = make([]bool, len(input[i]))
		}
		for len(stack) > 0 {
			n := len(stack)
			p := stack[n-1]
			stack = stack[:n-1]
			if flashed[p.Y][p.X] {
				continue
			}
			flashed[p.Y][p.X] = true
			for dy := -1; dy < 2; dy++ {
				for dx := -1; dx < 2; dx++ {
					p2 := util.Point2D{Y: p.Y + dy, X: p.X + dx}
					if p2.X < 0 || p2.X >= len(input[0]) || p2.Y < 0 || p2.Y >= len(input) {
						continue
					}
					if flashed[p2.Y][p2.X] {
						continue
					}
					input[p2.Y][p2.X]++
					if input[p2.Y][p2.X] > 9 {
						stack = append(stack, p2)
					}
				}
			}
		}

		flashes := 0
		for i := 0; i < len(input); i++ {
			for j := 0; j < len(input[i]); j++ {
				if flashed[i][j] {
					flashes++
					input[i][j] = 0
				}
			}
		}

		return flashes
	}

	t1 := time.Now()
	t := 0
	steps := 100
	for i := 0; i<steps; i++ {
		t += step()
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	steps++
	for step() != 100 {
		steps++
	}
	fmt.Printf("2: %v, %v\n", steps, time.Since(t2))
}
