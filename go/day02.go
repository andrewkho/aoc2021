package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Instr struct {
	Dirn string
	Steps int
}

type Position struct {
	x int
	y int
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
	pos := Position{0, 0}
	for _, instr := range instrs {
		if instr.Dirn == "forward" {
			pos.x += instr.Steps
		} else if instr.Dirn == "down" {
			pos.y += instr.Steps
		} else if instr.Dirn == "up" {
			pos.y -= instr.Steps
		}
	}
	fmt.Printf("1: %v, %v\n", pos.x*pos.y, time.Since(t1))

	t2 := time.Now()
	aim := 0
	pos = Position{0, 0}
	for _, inst := range instrs {
		if inst.Dirn == "forward" {
			pos.x += inst.Steps
			pos.y += aim * inst.Steps
		} else if inst.Dirn == "down" {
			aim += inst.Steps
		} else if inst.Dirn == "up" {
			aim -= inst.Steps
		}
	}

	fmt.Printf("2: %v, %v\n", pos.x*pos.y, time.Since(t2))
}
