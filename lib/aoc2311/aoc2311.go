package aoc2311

import (
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumDistances(lines []string, multiplier int) int {
	shared.Logger.Info(
		"Sum distances between galaxies.",
		"line count", len(lines),
		"multiplier", multiplier)
	brd := shared.NewBoard(lines)
	locations := expandLocations(extractGalaxyLocations(brd), multiplier)

	var pairs []gent.Pair[shared.Loc, shared.Loc]
	for i := range locations {
		for j := i + 1; j < len(locations); j++ {
			pairs = append(pairs, gent.NewPair(locations[i], locations[j]))
		}
	}

	var sum int
	for _, each := range pairs {
		sum += measureDistance(each.First, each.Second)
	}
	shared.Logger.Info("Distance sum calculate.", "sum", sum)
	return sum
}

func measureDistance(alpha, omega shared.Loc) int {
	return shared.Abs(alpha.X-omega.X) + shared.Abs(alpha.Y-omega.Y)
}

func expandLocations(locations []shared.Loc, multiplier int) []shared.Loc {
	maximum := shared.Loc{X: -1, Y: -1}
	xValues := gent.NewSet[int]()
	yValues := gent.NewSet[int]()
	for _, each := range locations {
		maximum.X = max(maximum.X, each.X)
		maximum.Y = max(maximum.Y, each.Y)
		xValues.Add(each.X)
		yValues.Add(each.Y)
	}

	findSpace := func(high int, set *gent.Set[int]) []int {
		var empty []int
		for i := 0; i <= high; i++ {
			if !set.Has(i) {
				empty = append(empty, i)
			}
		}
		return empty
	}
	rows := findSpace(maximum.Y, yValues)
	columns := findSpace(maximum.X, xValues)

	expanded := make([]shared.Loc, len(locations))
	for i, each := range locations {
		expanded[i] = shared.Loc{
			X: countBelow(each.X, columns)*(multiplier-1) + each.X,
			Y: countBelow(each.Y, rows)*(multiplier-1) + each.Y,
		}
	}
	return expanded
}

func countBelow(n int, empty []int) int {
	var count int
	for _, each := range empty {
		if each >= n {
			break
		}
		count++
	}
	return count
}

func extractGalaxyLocations(brd *shared.Board) []shared.Loc {
	var locations []shared.Loc
	brd.Iter(func(loc shared.Loc, c rune) bool {
		if c == '#' {
			locations = append(locations, loc)
		}
		return true
	})
	return locations
}
