package aoc2311

import (
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumDistances(lines []string) int {
	shared.Logger.Info("Sum distances between galaxies.", "line count", len(lines))
	brd := shared.NewBoard(expand(lines))
	var locations []shared.Loc
	brd.Iter(func(loc shared.Loc, c rune) bool {
		if c == '#' {
			locations = append(locations, loc)
		}
		return true
	})
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

func expand(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	expanded := make([]string, 0, len(lines))
	for _, each := range lines {
		if isFilledWith(each, '.') {
			expanded = append(expanded, each)
		}
		expanded = append(expanded, each)
	}
	var columns []int
mainLoop:
	for i := range expanded[0] {
		for j := range expanded {
			if expanded[j][i] != '.' {
				continue mainLoop
			}
		}
		columns = append(columns, i)
	}
	shared.Logger.Debug("columns", "columns", columns)
	for i := len(columns) - 1; i >= 0; i-- {
		for j, each := range expanded {
			col := columns[i]
			expanded[j] = each[:col] + "." + each[col:]
		}
	}
	return expanded
}

func isFilledWith(s string, c rune) bool {
	for _, each := range s {
		if each != c {
			return false
		}
	}
	return true
}

func measureDistance(alpha, omega shared.Loc) int {
	return shared.Abs(alpha.X-omega.X) + shared.Abs(alpha.Y-omega.Y)
}
