package aoc2218

import (
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveSurfaceArea(lines []string) int {
	shared.Logger.Info("Derive surface area.", "line count", len(lines))
	trips := parseLines(lines)
	cube := createIntCube(trips)
	var result int
	for _, each := range trips {
		result += cube.countFreeSides(each)
	}
	shared.Logger.Info("Surface area derived.", "area", result)
	return result
}

func parseLines(lines []string) [][3]int {
	trips := make([][3]int, 0, len(lines))
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed != "" {
			pieces := strings.Split(trimmed, ",")
			if len(pieces) != 3 {
				panic("invalid line (piece count != 3): " + trimmed)
			}
			var trip [3]int
			for i, piece := range pieces {
				parsed, err := strconv.ParseInt(piece, 10, 64)
				if err != nil {
					panic(err)
				}
				trip[i] = int(parsed)
			}
			trips = append(trips, trip)
		}
	}
	return trips
}

type intCube struct {
	matrix [][][]int
}

func createIntCube(trips [][3]int) *intCube {
	maxPuppets := make([]int, 3)
	for _, each := range trips {
		for i, val := range each {
			maxPuppets[i] = max(maxPuppets[i], val)
		}
	}
	shared.Logger.Debug("Got cube dimensions.", "dimensions", maxPuppets)
	matrix := make([][][]int, maxPuppets[0]+2)
	for i := range matrix {
		matrix[i] = make([][]int, maxPuppets[1]+2)
		for j := range matrix[i] {
			matrix[i][j] = make([]int, maxPuppets[2]+2)
		}
	}
	for _, each := range trips {
		matrix[each[0]][each[1]][each[2]] = 1
	}
	return &intCube{matrix: matrix}
}

func (v *intCube) countFreeSides(trip [3]int) int {
	alt := func(i, mod int) [3]int {
		cand := trip
		cand[i] += mod
		return cand
	}
	check := func(cand [3]int) int {
		for _, each := range cand {
			if each < 0 {
				return 1
			}
		}
		if v.matrix[cand[0]][cand[1]][cand[2]] == 0 {
			return 1
		}
		return 0
	}

	var count int
	for i := range 3 {
		for _, mod := range []int{-1, 1} {
			count += check(alt(i, mod))
		}
	}
	return count
}
