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

	re := regexp.MustCompile(`^target area: x=(-?[0-9].*)\.\.(-?[0-9].*), y=(-?[0-9].*)\.\.(-?[0-9].*)$`)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	groups := re.FindStringSubmatch(scanner.Text())
	var ints []int
	for _, m := range groups[1:] {
		v, err := strconv.Atoi(m)
		if err != nil {
			log.Fatalln(err)
		}
		ints = append(ints, v)
	}
	x0, x1, y0, y1 := ints[0], ints[1], ints[2], ints[3]

	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))
	fmt.Println(x0, x1, y0, y1)

	step := func(pos *util.Point2D, vel *util.Point2D) {
		pos.X += vel.X
		pos.Y += vel.Y

		if vel.X > 0 {
			vel.X--
		} else if vel.X < 0 {
			vel.X++
		}
		vel.Y--
	}

	trial := func(vel util.Point2D) (bool, int) {
		pos := util.Point2D{X: 0, Y: 0}
		maxY := 0
		for pos.Y >= y0 && pos.X <= x1 {
			if pos.Y > maxY {
				maxY = pos.Y
			}
			if pos.X >= x0 && pos.X <= x1 && pos.Y >= y0 && pos.Y <= y1 {
				return true, maxY
			}
			step(&pos, &vel)
		}
		return false, maxY
	}

	t1 := time.Now()
	t := 0
	for u := 1; u < x1+1; u++ {
		for v := 0; v < -y0; v++ {
			if hit, maxY := trial(util.Point2D{X: u, Y: v}); hit {
				if maxY > t {
					t = maxY
				}
			}
		}

	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	t = 0
	for u := 1; u < x1+1; u++ {
		for v := y0; v < -y0; v++ {
			if hit, _ := trial(util.Point2D{X: u, Y: v}); hit {
				t++
			}
		}

	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
