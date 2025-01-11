package calc

import (
	"math"
	"math/rand"
	"sort"
	"trucks/internal/models"
)

func Distance(p1 models.Point, p2 models.Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func RouteDistance(route []int64, loads []models.Load) float64 {
	// Returns distance an ordered series or loads costs
	d := 0.0
	lastPoint := models.Origin
	m := make(map[int64]models.Load)
	for _, load := range loads {
		m[load.Id] = load
	}
	for _, loadIdx := range route {
		load := m[loadIdx]
		d += Distance(lastPoint, load.Start)
		d += Distance(load.Start, load.End)
		lastPoint = load.End
	}
	d += Distance(lastPoint, models.Origin)
	return d
}

func Route(loads []models.Load, r float64) [][]int64 {
	const max = 12 * 60

	t := 0.0
	newLoads := make([]models.Load, len(loads))
	copy(newLoads, loads)
	if len(newLoads) < 1 {
		panic("No loads")
	}
	nextIdx := RandomIndex(newLoads)
	lastLoad := newLoads[nextIdx]
	lastLoadEnd := lastLoad.End
	routes := make([][]int64, 0)
	route := make([]int64, 0)
	route = append(route, lastLoad.Id)
	t += Distance(models.Origin, lastLoad.Start)
	t += Distance(lastLoad.Start, lastLoad.End)
	newLoads = remove(newLoads, nextIdx)
	futureLoads := make([]models.Load, 0)
	for t < max {
		nextIdx = NextIndex(lastLoadEnd, newLoads, r)
		nextLoad := newLoads[nextIdx]
		loadD := Distance(lastLoadEnd, nextLoad.Start) + Distance(nextLoad.Start, nextLoad.End)
		proposedD := loadD + Distance(models.Origin, nextLoad.End)
		if t+proposedD < max {
			t += loadD
			route = append(route, nextLoad.Id)
			lastLoadEnd = nextLoad.End
		} else {
			futureLoads = append(futureLoads, nextLoad)
		}
		newLoads = remove(newLoads, nextIdx)
		if len(newLoads) == 0 {
			routes = append(routes, route)
			if len(futureLoads) == 0 {
				break
			}
			route = make([]int64, 0)
			t = 0.0
			lastLoadEnd = models.Origin
			newLoads = futureLoads
			futureLoads = make([]models.Load, 0)
		}
	}

	return routes
}

func NextIndex(pt models.Point, loads []models.Load, r float64) int {
	// This is our greedy function. We want to consistently create a long route at a minimal distance.
	// We also will add randomness so that the iterations cover more of the tree
	if len(loads) == 1 {
		return 0
	}
	var i int
	f := rand.Float64()
	// if f < 1 {
	// 	i = WeightedDraw(pt, loads) //rand.Intn(len(loads))
	// 	//fmt.Printf("Weighted draw %d\n", i)
	// } else
	if f < r/2 {
		i = RandomIndex(loads)
	} else if f >= r/2 && f < r+r/2 {
		l := len(loads) / 2
		if l < 1 {
			l = 1
		}
		n := rand.Intn(l)
		i = NthClosestIdx(pt, loads, n)

	} else {
		i = ClosestIdx(pt, loads)
	}
	return i
}

// this does not appear to be working
// inspired by https://cybernetist.com/2019/01/24/random-weighted-draws-in-go/
// func WeightedDraw(pt models.Point, loads []models.Load) int {
// 	w := make([]float64, len(loads))
// 	for i, load := range loads {
// 		w[i] = 1.0/Distance(pt, load.Start) // was hoping to do distance weighted
// 	}
// 	cdf := make([]float64, len(w))
// 	floats.CumSum(cdf, w)
// 	val := distuv.UnitUniform.Rand() * cdf[len(cdf)-1]
// 	return sort.Search(len(cdf), func(i int) bool { return cdf[i] > val })
// }

func RandomIndex(loads []models.Load) int {
	return rand.Intn(len(loads))
}

func NthClosestIdx(pt models.Point, loads []models.Load, n int) int {
	sort.Slice(loads, func(i, j int) bool {
		return Distance(pt, loads[i].Start) > Distance(pt, loads[j].Start)
	})
	if n < len(loads) {
		return n
	} else {
		return len(loads)
	}
}

func ClosestIdx(pt models.Point, loads []models.Load) int {
	closest := math.MaxFloat64
	closestIdx := 0
	for i, load := range loads {
		if d := Distance(pt, load.Start); d < closest {
			closest = d
			closestIdx = i
		}
	}
	return closestIdx
}

func Cost(routes [][]int64, loads []models.Load) float64 {
	c := 500.0 * float64(len(routes))
	for _, r := range routes {
		c += RouteDistance(r, loads)
	}
	return c
}

// Unordered slice removal
// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(s []models.Load, i int) []models.Load {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
