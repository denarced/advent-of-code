package aoc2214

import (
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func MeasureSand(lines []string) int {
	sandPipe := shared.Loc{X: 500, Y: 0}
	shared.Logger.Info("Let the sand flow.", "line count", len(lines), "sand pipe", sandPipe)
	rocks := parseLines(lines)
	lowestRockY := findLowestRockLevel(rocks)
	isRock := createIsRock(rocks)
	sand := gent.NewSet[shared.Loc]()
primaryLoop:
	for {
		if sand.Contains(sandPipe) {
			break
		}
		grain := sandPipe
		for {
			candidate := grain
			candidate.Y++
			if !sand.Contains(candidate) && !isRock(candidate) {
				grain = candidate
				if grain.Y >= lowestRockY {
					break primaryLoop
				}
				continue
			}
			left := shared.Loc{X: candidate.X - 1, Y: candidate.Y}
			if !sand.Contains(left) && !isRock(left) {
				grain = left
				continue
			}
			right := shared.Loc{X: candidate.X + 1, Y: candidate.Y}
			if !sand.Contains(right) && !isRock(right) {
				grain = right
				continue
			}
			sand.Add(grain)
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Sand grain settled.", "loc", grain)
			}
			break
		}
	}
	return sand.Count()
}

func parseLines(lines []string) [][2]shared.Loc {
	var rocks [][2]shared.Loc
	for _, each := range lines {
		parseLine(each, func(from, to shared.Loc) {
			from, to = sortLocs(from, to)
			rocks = append(rocks, [2]shared.Loc{from, to})
		})
	}
	return rocks
}

func parseLine(line string, consume func(from, to shared.Loc)) {
	trimmed := strings.TrimSpace(line)
	pairs := gent.Map(strings.Split(trimmed, "->"), strings.TrimSpace)
	locs := gent.Map(pairs, func(s string) shared.Loc {
		pieces := strings.Split(s, ",")
		atoi := func(intStr string) int {
			i, err := strconv.Atoi(intStr)
			if err != nil {
				shared.Logger.Error("Failed to convert coordinate int.", "err", err, "str", intStr)
				panic(err)
			}
			return i
		}
		return shared.Loc{
			X: atoi(pieces[0]),
			Y: atoi(pieces[1]),
		}
	})
	for i := 0; i < len(locs)-1; i++ {
		consume(locs[i], locs[i+1])
	}
}

func createIsRock(rocks [][2]shared.Loc) func(shared.Loc) bool {
	return func(candidate shared.Loc) bool {
		for _, pair := range rocks {
			if pair[0].X <= candidate.X && candidate.X <= pair[1].X {
				if pair[0].Y <= candidate.Y && candidate.Y <= pair[1].Y {
					return true
				}
			}
		}
		return false
	}
}

func sortLocs(a, b shared.Loc) (first, second shared.Loc) {
	if a.X == b.X {
		if a.Y < b.Y {
			return a, b
		}
		return b, a
	}
	if a.X < b.X {
		return a, b
	}
	return b, a
}

func findLowestRockLevel(rocks [][2]shared.Loc) int {
	var lowest int
	for _, each := range rocks {
		lowest = max(lowest, each[0].Y)
		lowest = max(lowest, each[1].Y)
	}
	return lowest
}
