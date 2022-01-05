package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
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
	var lefts, rights [][]string
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " | ")

		left := make([]string, len(split[0]))
		for i, s := range strings.Split(split[0], " ") {
			r := []rune(s)
			sort.Slice(r, func(i int, j int) bool {return r[i] < r[j]})
			left[i] = string(r)
		}
		lefts = append(lefts, left)

		right := make([]string, len(split[1]))
		for i, s := range strings.Split(split[1], " ") {
			r := []rune(s)
			sort.Slice(r, func(i int, j int) bool {return r[i] < r[j]})
			right[i] = string(r)
		}
		rights = append(rights, right)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	t := 0
	for i, left := range lefts {
		mapping := decode(left, true)
		for _, s := range rights[i] {
			if _, ok := mapping[s]; ok {
				t++
			}
		}
	}
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	t = 0
	for i, left := range lefts {
		mapping := decode(left, false)
		var decoded string
		for _, s := range rights[i] {
			decoded += mapping[s]
		}
		val, err := strconv.Atoi(decoded)
		if err != nil {
			log.Panicln(err)
		}
		t += val
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}

func decode(left []string, part1 bool) map[string]string {
	mapping := make(map[string]string)
	lengths := make(map[int][]string)

	// 1:1 mappings from length
	var one string
	for _, s := range left {
		switch len(s) {
		case 2: 
			mapping[s] = "1"
			one = s
		case 4:
			mapping[s] = "4"
		case 3:
			mapping[s] = "7"
		case 7:
			mapping[s] = "8"
		default:
			lengths[len(s)] = append(lengths[len(s)], s)
		}
	}

	if part1 {
		return mapping
	}

	// 6
	var missing rune
	for _, s := range lengths[6] {
		if !strings.ContainsRune(s, rune(one[0])) {
			missing = rune(one[0])
		} else if !(strings.ContainsRune(s, rune(one[1]))) {
			missing = rune(one[1])
		}
		if missing != 0 {
			mapping[s] = "6"
			break
		}
	}

	// 2, 3, 5
	var missing2 rune
	for _, s := range lengths[5] {
		if strings.ContainsRune(s, rune(one[0])) && strings.ContainsRune(s, rune(one[1])) {
			mapping[s] = "3"
		} else if strings.ContainsRune(s, missing) {
			mapping[s] = "2"
		} else {
			mapping[s] = "5"
			for _, c := range "abcdefg" {
				if c == missing {
					continue
				}
				if !strings.ContainsRune(s, c) {
					missing2 = c
					break
				}
			}
		}
	}

	// 9, 0
	for _, s := range lengths[6] {
		if _, ok := mapping[s]; ok {
			continue
		} else if strings.ContainsRune(s, missing2) {
			mapping[s] = "0"
		} else {
			mapping[s] = "9"
		}
	}

	return mapping
}