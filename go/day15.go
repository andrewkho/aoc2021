package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"./util"
)

type PriorityPos struct {
	p int
	pos util.Point2D
}

type PriorityQueue []PriorityPos
    
func (h PriorityQueue) Len() int           { return len(h) }
func (h PriorityQueue) Less(i, j int) bool { return h[i].p < h[j].p }
func (h PriorityQueue) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PriorityQueue) Push(x interface{}) {
	*h = append(*h, x.(PriorityPos))
}

func (h *PriorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
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
	var risks util.Array2D
	for i := 0; scanner.Scan(); i++ {
		ints := util.GetInts(scanner.Text(), "")
		risks = append(risks, ints)
	}

	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	t := srch(risks)

	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2a := time.Now()
	risks5 := util.New2DZeros(len(risks)*5, len(risks[0])*5)
	for i := 0; i<len(risks5); i++ {
		for j := 0; j<len(risks5[i]); j++ {
			v := risks[i%len(risks)][j%len(risks[0])] + i/len(risks) + j/len(risks[0])
			risks5[j][i] = (v-1) % 9 + 1
		}
	}
	fmt.Printf("t2a: %v\n", time.Since(t2a))
	t2b := time.Now()
	t = srch(risks5)
	fmt.Printf("2: %v, %v\n", t, time.Since(t2b))
}

func srch(risks util.Array2D) int {
	N, M := len(risks), len(risks[0])
	start := util.Point2D{X: 0, Y: 0}
	end := util.Point2D{X: M-1, Y: N-1}

	dist := util.New2DZeros(N, M)
	for i := 0; i<N; i++ {
		for j := 0; j<M; j++ {
			dist[i][j] = 1<<63-1
		}
	}
	dist[start.Y][start.X] = 0
	pq := PriorityQueue{PriorityPos{0, start}}

	for len(pq) > 0 {
		p := heap.Pop(&pq).(PriorityPos)
		if p.pos == end {
			continue
		}
		for _, d := range [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			p2 := util.Point2D{X: p.pos.X+d[0], Y: p.pos.Y+d[1]}
			if p2.X < 0 || p2.X >= M || p2.Y < 0 || p2.Y >= N {
				continue
			}
			d2 := p.p + risks[p2.Y][p2.X]
			if d2 >= dist[p2.Y][p2.X] {
				continue
			}
			dist[p2.Y][p2.X] = d2
			heap.Push(&pq, PriorityPos{p: d2, pos: p2})
		}
	}

	return dist[end.Y][end.X]
}