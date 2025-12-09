package aoc2508

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

type pointDistance struct {
	points   [2]int
	distance float64
}

func CountCircuits(lines []string, limit int) int {
	if len(lines) < 2 {
		return 1
	}
	points := parseLines(lines)
	distances := measureDistances(points)
	pointDistances := sortDistances(distances)
	if limit < len(pointDistances) {
		pointDistances = pointDistances[:limit]
	}
	shared.Logger.Debug("Distances.", "distances", pointDistances)
	pointToCircuit := gatherCircuits(pointDistances)
	shared.Logger.Debug("Circuits formed.", "pointToCircuit", pointToCircuit)
	sizes := deriveCircuitSizes(pointToCircuit)
	shared.Logger.Debug("Sizes derived.", "sizes", sizes)
	if len(sizes) > 3 {
		sizes = sizes[:3]
	}
	result := 1
	for _, each := range sizes {
		result *= each
	}
	return result
}

type point [3]int

func (p point) String() string {
	return fmt.Sprintf("%dx%dx%d", p[0], p[1], p[2])
}

func measureDistance(a, b point) float64 {
	x := float64(abs(a[0] - b[0]))
	y := float64(abs(a[1] - b[1]))
	z := float64(abs(a[2] - b[2]))
	ho := hypo(x, y)
	return hypo(ho, z)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func hypo(a, b float64) float64 {
	return math.Sqrt(a*a + b*b)
}

func parseLines(lines []string) (points []point) {
	for _, each := range lines {
		pieces := strings.Split(each, ",")
		points = append(points, point{
			toIntOrDie(pieces[0]),
			toIntOrDie(pieces[1]),
			toIntOrDie(pieces[2]),
		})
	}
	return
}

func toIntOrDie(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func measureDistances(points []point) map[[2]int]float64 {
	distances := map[[2]int]float64{}
	for i := 0; i < len(points)-1; i++ {
		for j := i + 1; j < len(points); j++ {
			dis := measureDistance(points[i], points[j])
			distances[[2]int{i, j}] = dis
		}
	}
	return distances
}

func sortDistances(distances map[[2]int]float64) []pointDistance {
	pointDistances := make([]pointDistance, 0, len(distances))
	for k, v := range distances {
		pointDistances = append(pointDistances, pointDistance{points: k, distance: v})
	}
	slices.SortFunc(pointDistances, func(a, b pointDistance) int {
		if a.distance < b.distance {
			return -1
		}
		if b.distance < a.distance {
			return 1
		}
		return 0
	})
	return pointDistances
}

func gatherCircuits(pointDistances []pointDistance) map[int]int {
	pointToCircuit := map[int]int{}
	circuitIndex := -1
	for _, each := range pointDistances {
		aCircuit, aOk := pointToCircuit[each.points[0]]
		bCircuit, bOk := pointToCircuit[each.points[1]]
		if !aOk && !bOk {
			// New circuit.
			circuitIndex++
			pointToCircuit[each.points[0]] = circuitIndex
			pointToCircuit[each.points[1]] = circuitIndex
		} else if aOk && !bOk {
			// Add b to a's circuit.
			pointToCircuit[each.points[1]] = aCircuit
		} else if !aOk && bOk {
			// Add a to b's circuit.
			pointToCircuit[each.points[0]] = bCircuit
		} else if aCircuit != bCircuit {
			// Merge circuits, move points in b circuit to a circuit.
			var movedPoints []int
			for p, c := range pointToCircuit {
				if c == bCircuit {
					movedPoints = append(movedPoints, p)
				}
			}
			for _, p := range movedPoints {
				pointToCircuit[p] = aCircuit
			}
		}
	}
	return pointToCircuit
}

func deriveCircuitSizes(pointToCircuit map[int]int) []int {
	circuitToPoints := map[int]int{}
	for _, circuit := range pointToCircuit {
		p := circuitToPoints[circuit]
		p++
		circuitToPoints[circuit] = p
	}
	shared.Logger.Debug("Circuit points gathered.", "circuit to points", circuitToPoints)
	sizes := make([]int, len(circuitToPoints))
	var i int
	for _, count := range circuitToPoints {
		sizes[i] = count
		i++
	}
	slices.SortFunc(sizes, func(a, b int) int {
		return b - a
	})
	return sizes
}
