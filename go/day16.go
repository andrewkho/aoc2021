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

type Packet interface  {
	getVersionSum() int
	getValue() int
}

type Literal struct {
	v int
	t int
	val int
}

func (lit Literal) getVersionSum() int {
	return lit.v
}

func (lit Literal) getValue() int {
	return lit.val
}

type Operator struct {
	v int
	t int
	subpackets []Packet
}

func (op Operator) getVersionSum() int {
	s := op.v
	for _, packet := range op.subpackets {
		s += packet.getVersionSum()
	}
	return s
}

func (op Operator) getValue() int {
	switch op.t {
		case 0: {
			s := 0
			for _, packet := range op.subpackets {
				s += packet.getValue()
			}
			return s
		} 
		case 1: {
			s := 1
			for _, packet := range op.subpackets {
				s *= packet.getValue()
			}
			return s
		}
		case 2: {
			s := 1<<63-1
			for _, packet := range op.subpackets {
				v := packet.getValue()
				if v < s {
					s = v
				}
			}
			return s
		} 
		case 3: {
			s := 0
			for _, packet := range op.subpackets {
				v := packet.getValue()
				if v > s {
					s = v
				}
			}
			return s
		}
		case 5, 6, 7: {
			v0 := op.subpackets[0].getValue()
			v1 := op.subpackets[1].getValue()
			sel := false
			switch op.t {
			case 5: sel = v0 > v1
			case 6: sel = v0 < v1
			case 7: sel = v0 == v1
			}
			if sel {
				return 1
			} else {
				return 0
			}
		}
	default: log.Fatalln(op.t)
	}
	return -1
}

func parseBin(binstr string) int {
	// If you want to compile this on a 32-bit system you
	// need to handle this properly
	val, err := strconv.ParseInt(binstr, 2, 64)
	if err != nil {
		log.Panicln(err, binstr)
	}
	return int(val)
}

func parse(bits string, i int) (Packet, int) {
	v := parseBin(bits[i:i+3])
	i += 3
	t := parseBin(bits[i:i+3])
	i += 3

	if t == 4 {
		vals := ""
		for true {
			leading := parseBin(bits[i:i+1])
			vals += bits[i+1:i+5]
			i += 5
			if leading == 0 {
				break
			}
		}
		val := parseBin(vals)
		return Literal{v: v, t: t, val: val}, i
	} else {
		I := parseBin(bits[i:i+1])
		i += 1
		var subpackets []Packet
		if I == 0 {
			L := parseBin(bits[i:i+15])
			i += 15
			fin := i+int(L)
			for i < fin {
				var x Packet
				x, i = parse(bits, i)
				subpackets = append(subpackets, x)
			}
		} else {
			L := parseBin(bits[i:i+11])
			i += 11
			for j := 0; j<int(L); j++ {
				var x Packet
				x, i = parse(bits, i)
				subpackets = append(subpackets, x)
			}
		}
		return Operator{v: v, t: t, subpackets: subpackets}, i
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
	scanner.Scan()
	input := scanner.Text()
	var bits string
	for _, c := range input {
		hex, err := strconv.ParseUint(string(c), 16, 4)
		if err != nil {
			log.Panicln(err)
		}
		bits += fmt.Sprintf("%04b", hex)
	}

	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	parsed, _ := parse(bits, 0)
	t := parsed.getVersionSum()
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	t = parsed.getValue()
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}
