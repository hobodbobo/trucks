package main

import (
	"math"
	"os"
	"time"

	"trucks/internal/calc"
	"trucks/internal/models"
)

func main() {
	filesToEval := os.Args[1:]
	if len(filesToEval) == 0 {
		// I chose to panic overall everywhere. Expect happy path valid input.
		panic("Enter a shipping manifest to evaluate")
	}
	for _, evalFile := range filesToEval {
		//fmt.Printf("Routing shipments for %s\n", evalFile)

		loads := models.NewLoads(evalFile)

		// fmt.Printf("%d loads to ship\n", len(loads))
		// fmt.Printf("%d points created\n", models.PtId.Load())

		// time keeping
		const secondsToRun = 29.5
		timeout := time.After(time.Duration(secondsToRun*1000.0) * time.Millisecond)
		ticker := time.NewTicker(5 * time.Second)
		start := time.Now()
		defer ticker.Stop()

		// store the current min cost and route
		minC := math.MaxFloat64
		var lowestCostRoutes [][]int64
		count := 0
	outerLoop:
		for {
			select {
			case <-timeout:
				// out of time
				break outerLoop
			case <-ticker.C:
				// for feedback while waiting
				//fmt.Printf("%d seconds: %d iterations run\n", i, count)
			default:
				randomNess := time.Since(start).Seconds() / secondsToRun * 0.1
				routes := calc.Route(loads, randomNess)
				c := calc.Cost(routes, loads)
				if c < minC {
					minC = c
					lowestCostRoutes = routes
				}
				count++
			}
		}

		// print each route of loads
		models.PrintRoutes(lowestCostRoutes)

		//fmt.Printf("%d total loads assigned with %d drivers\n", total, len(lowestCostRoutes))
		//fmt.Printf("total cost %f", calc.Cost(lowestCostRoutes, loads))
	}

}
