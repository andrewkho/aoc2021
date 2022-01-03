package util

import (
	"log"
	"strconv"
	"strings"
)

func GetInts(line string, sep string) []int {
	vals := strings.Split(line, sep)
	ints := make([]int, len(vals))
	for i, val := range vals {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Panic(err)
		}
		ints[i] = v
	}
	return ints
}