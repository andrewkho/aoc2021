package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var filename string = os.Args[1]
	t0 := time.Now()
	fmt.Printf("Reading from %s,", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	width := len(lines[0])
	bits := make([]int, width)
	for _, line := range lines {
		for i, c := range line {
			if c == '1' {
				bits[i] += 1
			}
		}
	}
	g, e := 0, 0
	for i, c := range bits {
		if c > len(lines) / 2 {
			g += 1 << (width-i-1)
		} else {
			e += 1 << (width-i-1)
		}
	}
	fmt.Printf("1: %v, %v\n", g*e, time.Since(t1))

	t2 := time.Now()
	o2, err := computeVals(lines, false)
	co2, err := computeVals(lines, true)

	fmt.Printf("2: %v, %v, %v, %v\n", o2*co2, o2, co2, time.Since(t2))
}

func computeVals(lines []string, co2 bool) (int, error) {
	includes := map[int]bool{}
	for i := 0; i<len(lines); i++ {
		includes[i] = true
	}

	i := 0
	for len(includes) > 1 {
		b := 0
		for k := range includes {
			if lines[k][i] == '1' {
				b++
			}
		}
		c := '1'
		if b >= len(includes) / 2 {
			c = '0'
		}
		for k := range includes {
			if (!co2 && rune(lines[k][i]) == c) || (co2 && rune(lines[k][i]) != c) {
				delete(includes, k)
			}	
		}
		i += 1
	} 

	if len(includes) != 1 {
		return 0, errors.New(fmt.Sprintf("%v", includes))
	}

	v := 0
	for k := range includes {
		w := len(lines[k])
		for i, c := range lines[k] {
			if c == '1' {
				v += 1 << (w - i - 1)
			}
		}
		break
	}
	return v, nil
}