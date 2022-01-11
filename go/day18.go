package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

type Node struct {
	l *Node
	r *Node
	val int
	isLeaf bool
}

func (node *Node) toString() string {
	if node.isLeaf {
		return fmt.Sprintf("%v", node.val)
	} else {
		return fmt.Sprintf("[%v,%v]",
			node.l.toString(), 
			node.r.toString(),
		)
	}
}

func parse(snum string) *Node {
	if val, err := strconv.Atoi(snum); err == nil {
		return &Node{nil, nil, val, true}
	}
	// else parse out substrings
	var stack []int
	var comma int
	var left, right *Node
	for i, c := range snum {
		switch c {
		case '[': {
			stack = append(stack, i)
		}
		case ']': {
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				right = parse(snum[comma+1:i])
			}
		}
		case ',': {
			if len(stack) == 1 {
				left = parse(snum[1:i])
				comma = i
			}
		}
		}
	}

	return &Node{left, right, -1, false}
}

func explode(x *Node) bool {
	type Pair struct {
		node *Node
		depth int
	}
	cur := Pair{x, 0}
	var stack []Pair
	var prevprev, prev *Node
	var rval int
	found := false

	for cur.node != nil || len(stack) > 0 {
		if cur.node != nil {
			stack = append(stack, cur)
			cur = Pair{cur.node.l, cur.depth+1}
		} else {
			cur = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// process cur
			if cur.node.isLeaf {
				if found {
					cur.node.val += rval
					return true
				}
				prevprev = prev
				prev = cur.node
			} else if cur.depth == 4 {
				if cur.node.isLeaf {
					log.Panicln(fmt.Sprintf("%v", cur.node.toString()))
				}
				lval := cur.node.l.val
				rval = cur.node.r.val
				found = true
				cur.node.l, cur.node.r = nil, nil
				cur.node.val, cur.node.isLeaf = 0, true
				if prevprev != nil {
					prevprev.val += lval
				}
			}
			cur = Pair{cur.node.r, cur.depth+1}
		}
	}
	return found
}

func split(cur *Node) bool {
	var stack []*Node
	for cur != nil || len(stack) > 0 {
		if cur != nil {
			stack = append(stack, cur)
			cur = cur.l
		} else {
			cur = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if cur.isLeaf && cur.val > 9 {
				cur.l = &Node{nil, nil, cur.val/2, true}
				cur.r = &Node{nil, nil, (cur.val+1)/2, true}
				cur.val = 0
				cur.isLeaf = false
				return true
			}
			cur = cur.r
		}
	}
	return false
}

func reduce(x *Node) {
	for true {
		if explode(x) {
			continue
		}
		if split(x) {
			continue
		}
		break
	}
}

func magnitude(x *Node) int {
	if x.isLeaf {
		return x.val
	} else {
		return 3*magnitude(x.l) + 2*magnitude(x.r)
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
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	x := parse(lines[0])
	for _, line := range lines[1:] {
		node := parse(line)
		x = &Node{x, node, 0, false}
		reduce(x)
	}
	t := magnitude(x)
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	t = 0
	for i, line := range lines {
		for j, line2 := range lines {
			if i == j {
				continue
			}
			x = &Node{parse(line), parse(line2), 0, false}
			reduce(x)
			m := magnitude(x)
			if m > t {
				t = m
			}
		}
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
