package models

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var Origin = NewPoint(0, 0)

type Load struct {
	Id    int64
	Start Point
	End   Point
}

type Point struct {
	X  float64
	Y  float64
	Id int64
}

var PtId atomic.Int64

func NewPoint(x float64, y float64) Point {
	return Point{
		X:  x,
		Y:  y,
		Id: PtId.Add(1),
	}
}
func NewLoads(evalFile string) []Load {
	f, err := os.Open(evalFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()
	s := bufio.NewScanner(f)
	i := 0
	loads := make([]Load, 0)
	s.Scan() // Consume column header
	s.Text() // Consume column header
	for s.Scan() {
		line := s.Text()
		load := ParseLoad(line)
		loads = append(loads, load)
		i++
	}
	return loads
}

func ParseLoad(s string) Load {
	fields := strings.Fields(s)
	if len(fields) != 3 {
		panic(fmt.Sprintf("received bad point line %s", s))
	}
	id, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Cannot load id %d %s", id, err))
	}
	p1 := ParsePoint(fields[1])
	p2 := ParsePoint(fields[2])

	return Load{
		Id:    id,
		Start: p1,
		End:   p2,
	}
}

func ParsePoint(s string) Point {
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	xy := strings.Split(s, ",")
	if len(xy) != 2 {
		panic(fmt.Sprintf("Point without two values %s", s))
	}
	x, err := strconv.ParseFloat(xy[0], 64)
	if err != nil {
		panic(fmt.Sprintf("Could not convert %s to float", xy[0]))
	}
	y, err := strconv.ParseFloat(xy[1], 64)
	if err != nil {
		panic(fmt.Sprintf("Could not convert %s to float", xy[1]))
	}
	return NewPoint(x, y)
}

func PrintRoutes(routes [][]int64) {
	for _, route := range routes {
		routeStr := make([]string, len(route))
		for i, route := range route {
			routeStr[i] = strconv.Itoa(int(route))
		}
		fmt.Printf("[%s]\n", strings.Join(routeStr, ","))
	}
}
