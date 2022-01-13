package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"time"

	"./util"
)

type VectorPair [2]util.Vector3D

type ProbeScanner struct {
	probes []util.Vector3D
	rot *util.Matrix3D
	translate *util.Vector3D
	hash map[int][]VectorPair
}

func (ps *ProbeScanner) InitHash() {
	n := len(ps.probes)
	ps.hash = make(map[int][]VectorPair, n*(n-1)/2)
	for i, p0 := range ps.probes {
		for _, p1 := range ps.probes[i+1:] {
			dx := p0.Sub(p1)
			v := dx.Dot(dx)
			ps.hash[v] = append(ps.hash[v], VectorPair{p0, p1})
		}
	}
}

func (ps *ProbeScanner) HashOverlap(o *ProbeScanner) (int, map[util.Vector3D]bool, map[util.Vector3D]bool) {
	overlap := 0
	left := make(map[util.Vector3D]bool)
	right := make(map[util.Vector3D]bool)
	for h1, v1 := range o.hash {
		if v0, ok := ps.hash[h1]; ok {
			overlap += util.Min(len(v0), len(v1))
			for _, el := range v0 {
				left[el[0]] = true
				left[el[1]] = true
			}
			for _, el := range v1 {
				right[el[0]] = true
				right[el[1]] = true
			}
		}
	}
	return overlap, left, right
}

func (ps *ProbeScanner) Union(s1 *ProbeScanner) {
	seen := make(map[util.Vector3D]bool, len(ps.probes))
	for _, probe := range ps.probes {
		seen[probe] = true
	}
	for _, probe := range s1.probes {
		txProbe := s1.translate.Add(s1.rot.Dot(probe))
		if _, ok := seen[txProbe]; !ok {
			ps.probes = append(ps.probes, txProbe)
		}
	}

	for h1, v1 := range s1.hash {
		for _, o := range v1 {
			o2 := VectorPair{ 
				s1.translate.Add(s1.rot.Dot(o[0])), 
				s1.translate.Add(s1.rot.Dot(o[1])),
			}
			ps.hash[h1] = append(ps.hash[h1], o2) 
		}
	}
}

func checkOverlap(leftover map[util.Vector3D]bool, rightover map[util.Vector3D]bool, rot *util.Matrix3D) (bool, util.Vector3D) {
	uniques := make(map[util.Vector3D]int)
	for pt0 := range leftover {
		for pt1 := range rightover {
			dx := pt0.Sub(rot.Dot(pt1))
			uniques[dx]++
		}
	}
	
	for k, v := range uniques {
		if v >= 12 {
			return true, k
		}
	}
	return false, util.Vector3D{}
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
	var scanners []*ProbeScanner
	var ps *ProbeScanner
	re := regexp.MustCompile(`.* scanner ([0-9].*) .*`)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		groups := re.FindStringSubmatch(line)
		if groups != nil {
			eye := util.Matrix3D{}
			for i := 0; i<3; i++ {
				eye[i][i] = 1
			}
			ps = &ProbeScanner{[]util.Vector3D{}, &eye, &util.Vector3D{}, nil}
			scanners = append(scanners, ps)
			continue
		}
		pt := util.Vector3D{}
		for i, v := range util.GetInts(line, ",") {
			pt[i] = v
		}
		ps.probes = append(ps.probes, pt)
	}
	if err := scanner.Err(); err != nil {
		log.Panicln(err)
	}
	fmt.Printf(" took %v\n", time.Since(t0))

	t1 := time.Now()
	for _, s := range scanners {
		s.InitHash()
	}
	ROTS := makeRots()
	s0 := scanners[0]
	remaining := scanners[1:]
	for len(remaining) > 0 {
		s1 := remaining[0]
		remaining = remaining[1:]

		if overlap, overl, overr := s0.HashOverlap(s1); overlap < 66 {
			remaining = append(remaining, s1)
			continue
		} else {
			for _, rot := range ROTS {
				if ok, tx := checkOverlap(overl, overr, rot); ok {
					s1.rot = rot
					s1.translate = &tx
					s0.Union(s1)
					break
				}
			}
		}
	}
	t := len(s0.probes)
	fmt.Printf("1: %v, %v\n", t, time.Since(t1))

	t2 := time.Now()
	t = 0
	for i, s0 := range scanners {
		for _, s1 := range scanners[i+1:] {
			d := 0
			for k := 0; k<3; k++ {
				d += util.Abs(s0.translate[k] - s1.translate[k])
			}
			if d > t {
				t = d
			}
		}
	}
	fmt.Printf("2: %v, %v\n", t, time.Since(t2))
}

func makeRots() []*util.Matrix3D {
	itor := make([][]int, 6)
	for dim := 0; dim<3; dim++ {
		for i, v := range []int{-1, 1} {
			itor[dim*2+i] = []int{dim, v}
		}
	}

	var rots []*util.Matrix3D
	for _, xv := range itor {
		x, xval := xv[0], xv[1]
		for _, yv := range itor {
			y, yval := yv[0], yv[1]
			if y == x {
				continue
			}

			var xvec, yvec util.Vector3D
			xvec[x] = xval
			yvec[y] = yval

			rot := util.Matrix3D{
				xvec,
				yvec,
				xvec.Cross(yvec),
			}
			rots = append(rots, &rot)
		}
	}

	return rots
}
