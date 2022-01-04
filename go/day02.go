package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"./util"
)

type Instr struct {
	dirn string
	steps int
}

func main() {
	var filename string = os.Args[1]
	t0 := time.Now()
	fmt.Printf("Reading from %s,", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var instrs []Instr
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var dir string
		var steps int
		if _, err := fmt.Sscanf(scanner.Text(), "%s %d", &dir, &steps); err != nil {
			log.Fatal(err)
		} else {
			instrs = append(instrs, Instr{dir, steps})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	pos := util.Point2D{}
	for _, instr := range instrs {
		if instr.dirn == "forward" {
			pos.X += instr.steps
		} else if instr.dirn == "down" {
			pos.Y += instr.steps
		} else if instr.dirn == "up" {
			pos.Y -= instr.steps
		}
	}
	fmt.Printf("1: %v, %v\n", pos.X*pos.Y, time.Since(t1))

	t2 := time.Now()
	aim := 0
	pos = util.Point2D{}
	for _, inst := range instrs {
		if inst.dirn == "forward" {
			pos.X += inst.steps
			pos.Y += aim * inst.steps
		} else if inst.dirn == "down" {
			aim += inst.steps
		} else if inst.dirn == "up" {
			aim -= inst.steps
		}
	}

	fmt.Printf("2: %v, %v\n", pos.X*pos.Y, time.Since(t2))
}
