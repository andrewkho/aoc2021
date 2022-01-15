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
	var init [2]int
	re := regexp.MustCompile(`.* ([0-9].*)$`)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		groups := re.FindStringSubmatch(line)
		if groups == nil {
			log.Panicln(line)
		}
		pos, err := strconv.Atoi(groups[1])
		if err != nil {
			log.Panicln(err)
		}
		init[i] = pos
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	gameState := GameState{pos: init}
	rolls := 0
	for player := 0; gameState.score[0] < 1000 && gameState.score[1] < 1000; player = 1-player {
		dice := rolls % 100 + 1
		gameState = gameState.move(player, 3*dice+3)
		rolls += 3
	}
	t := util.Min(gameState.score[0], gameState.score[1])*rolls
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	
	var roller [10]int
	for i := 1; i<4; i++ {
		for j := 1; j<4; j++ {
			for k := 1; k<4; k++ {
				roller[i+j+k]++
			}
		}
	}
	var winners [2]int
	states := make(map[GameState]int)
	states[GameState{pos: init}] = 1
	for player := 0; len(states) > 0; player = 1-player{
		newStates := make(map[GameState]int)
		for state, n := range states {
			for i, m := range roller[3:] {
				newState := state.move(player, i+3)
				if newState.score[player] >= 21 {
					winners[player] += n*m
				} else {
					newStates[newState] += n*m
				}
			}
		}
		states = newStates
	}
	t = util.Max(winners[0], winners[1])
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}

type GameState struct {
	pos [2]int
	score [2]int
}

func (state GameState) move(player int, steps int) GameState {
	state.pos[player] += steps
	state.pos[player] = (state.pos[player] - 1) % 10 + 1
	state.score[player] += state.pos[player]
	return state
}
