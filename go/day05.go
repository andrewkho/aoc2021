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

type Vent struct {
	Start util.Point2D
	End util.Point2D
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
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var vents []Vent
	nrows := 0
	ncols := 0
	re := regexp.MustCompile(`^([0-9].*),([0-9].*) -> ([0-9].*),([0-9].*)$`)
	for i := 0; scanner.Scan(); i++ {
		match := re.FindStringSubmatch(scanner.Text())
		if len(match) == 0 {
			continue
		}
		ints := make([]int, 4)
		for i, v := range match[1:] {
			ints[i], err = strconv.Atoi(v)
			if err != nil {
				log.Panicln(err)
			}
			if i%2 == 1 && ints[i]+1 > nrows {
				nrows = ints[i]+1
			} else if i%2 == 0 && ints[i]+1 > ncols {
				ncols = ints[i]+1
			} 
		}
		vents = append(vents, Vent{ 
			Start: util.Point2D{X: ints[0], Y: ints[1]},
			End: util.Point2D{X: ints[2], Y: ints[3]},
		})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()

	board := util.New2DZeros(nrows, ncols)
	for _, vent := range vents {
		MarkBoard(board, vent, false)
	}

	t := 0
	for i := 0; i<nrows; i++ {
		for j := 0; j<ncols; j++ {
			if board[i][j] > 1 {
				t++
			}
		}
	}

	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	board = util.New2DZeros(nrows, ncols)
	for _, vent := range vents {
		MarkBoard(board, vent, true)
	}

	t = 0
	for i := 0; i<nrows; i++ {
		for j := 0; j<ncols; j++ {
			if board[i][j] > 1 {
				t++
			}
		}
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}

func MarkBoard(board util.Array2D, vent Vent, diag bool) {
	x0, x1 := vent.Start.X, vent.End.X
	y0, y1 := vent.Start.Y, vent.End.Y

	Nx := util.Abs(x0 - x1) + 1
	Ny := util.Abs(y0 - y1) + 1
	N := util.Max(Nx, Ny)

	if !diag && Nx > 1 && Ny > 1 {
		return
	}

	var dx, dy int
	if x0 < x1 {
		dx = 1
	} else if x0 > x1 {
		dx = -1
	} else {
		dx = 0
	}

	if y0 < y1 {
		dy = 1
	} else if y0 > y1 {
		dy = -1
	} else {
		dy = 0
	}

	for i := 0; i<N; i++ {
		board[y0+i*dy][x0+i*dx] += 1
	}
	
}