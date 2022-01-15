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

type Limits [3][2]int

func (b Limits) intersect(o *Limits) *Limits {
	ix := Limits{}
	for dim := 0; dim < 3; dim++ {
		ix[dim][0] = util.Max(b[dim][0], o[dim][0])
		ix[dim][1] = util.Min(b[dim][1], o[dim][1])
		if ix[dim][0] > ix[dim][1] {
			return nil
		}
	}
	return &ix
}

type Instruction struct {
	onoff string
	limits Limits
}

type Box struct {
	limits Limits
	children []*Box
}

func (g Box) getArea() int {
	area := 1
	for dim := 0; dim < 3; dim++ {
		area *= g.limits[dim][1] - g.limits[dim][0] + 1
	}
	for _, child := range g.children {
		area -= child.getArea()
	}
	return area
}

func (g *Box) turnOff(limits *Limits) {
	if ix := g.limits.intersect(limits); ix != nil {
		for _, child := range g.children {
			child.turnOff(ix)
		}
		g.children = append(g.children, &Box{limits: *ix})
	}
}

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
	var instructions []Instruction
	re := regexp.MustCompile(
		`^(on|off) x=(-?[0-9].*)\.\.(-?[0-9].*),y=(-?[0-9].*)\.\.(-?[0-9].*),z=(-?[0-9].*)\.\.(-?[0-9].*)$`)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		groups := re.FindStringSubmatch(line)
		if groups == nil {
			log.Panicln(line)
		}
		instr := Instruction{onoff: groups[1]}
		for dim := 0; dim<3; dim++ {
			for i := 0; i<2; i++ {
				if val, err := strconv.Atoi(groups[2+2*dim+i]); err != nil {
					log.Panicln(err)
				} else {
					instr.limits[dim][i] = val
				}
			}
		}
		instructions = append(instructions, instr)
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	limits50 := &Limits{{-50, 50}, {-50, 50}, {-50, 50}}
	var onboxes []*Box
	for _, instr := range instructions {
		for _, box := range onboxes {
			box.turnOff(&instr.limits)
		}
		if instr.onoff == "on" {
			if ix := instr.limits.intersect(limits50); ix != nil {
				onboxes = append(onboxes, &Box{*ix, nil})
			}
		}
	}
	t := 0
	for _, box := range onboxes {
		t += box.getArea()
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	onboxes = nil
	for _, instr := range instructions {
		for _, box := range onboxes {
			box.turnOff(&instr.limits)
		}
		if instr.onoff == "on" {
			onboxes = append(onboxes, &Box{instr.limits, nil})
		}
	}
	t = 0
	for _, box := range onboxes {
		t += box.getArea()
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}