package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Instruction struct {
	dir string
	steps int
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

	var instructions []Instruction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var dir string
		var steps int
		if _, err := fmt.Sscanf(scanner.Text(), "%s %d", &dir, &steps); err != nil {
			log.Fatal(err)
		} else {
			instructions = append(instructions, Instruction{dir, steps})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	pos := Position{0, 0}
	for _, inst := range instructions {
		if inst.dir == "forward" {
			pos.x += inst.steps
		} else if inst.dir == "down" {
			pos.y += inst.steps
		} else if inst.dir == "up" {
			pos.y -= inst.steps
		}
	}
	fmt.Printf("1: %v, %v\n", pos.x*pos.y, time.Since(t1))

	t2 := time.Now()
	aim := 0
	pos = Position{0, 0}
	for _, inst := range instructions {
		if inst.dir == "forward" {
			pos.x += inst.steps
			pos.y += aim * inst.steps
		} else if inst.dir == "down" {
			aim += inst.steps
		} else if inst.dir == "up" {
			aim -= inst.steps
		}
	}

	fmt.Printf("2: %v, %v\n", pos.x*pos.y, time.Since(t2))
}
