package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
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
		log.Panicln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	edges := make(map[string][]string)
	lower := make(map[string]bool)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "-")
		for _, s := range row {
			if s == strings.ToLower(s) {
				lower[s] = true
			}
		}
		edges[row[0]] = append(edges[row[0]], row[1])
		edges[row[1]] = append(edges[row[1]], row[0])
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	visited := make(map[string]int)

	var srch func(node string, maxVisits int) int
	srch = func(node string, maxVisits int) int {
		if node == "end" {
			return 1
		}
		visited[node]++
		paths := 0
		maxCheck := maxVisits
		if maxVisits == 2 {
			for k := range lower {
				if visited[k] == 2 {
					maxCheck = 1
					break
				}
			}
		}
		for _, nei := range edges[node] {
			if lower[nei] && visited[nei] >= maxCheck {
				continue
			}
			paths += srch(nei, maxVisits)
		}
		visited[node]--
		return paths
	}
	t1 := time.Now()
	visited["start"] = 1
	t := srch("start", 1)
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	visited = make(map[string]int)
	visited["start"] = 2
	t = srch("start", 2)
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
