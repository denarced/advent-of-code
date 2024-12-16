package aoc2408

import (
	"fmt"

	"github.com/denarced/advent-of-code/shared"
)

func CountUniqueAntinodeLocations(lines []string, resonantHarmonics bool) int {
	shared.Logger.Info(
		"Count unique antinode locations.",
		"line count",
		len(lines),
		"harmonics",
		resonantHarmonics,
	)
	return deriveUniqueAntinodeLocations(lines, resonantHarmonics).Count()
}

func deriveUniqueAntinodeLocations(
	lines []string,
	resonantHarmonics bool,
) *shared.Set[shared.Loc] {
	antinodes := shared.NewSet([]shared.Loc{})
	if len(lines) == 0 {
		return antinodes
	}
	brd := shared.NewBoard(lines)
	freqToAntennas := map[rune][]shared.Loc{}
	brd.Iter(func(loc shared.Loc, c rune) bool {
		if isAntenna(c) {
			existing, ok := freqToAntennas[c]
			if ok {
				freqToAntennas[c] = append(existing, loc)
			} else {
				freqToAntennas[c] = []shared.Loc{loc}
			}
		}
		return true
	})
	for freq, locations := range freqToAntennas {
		shared.Logger.Info("Derive antinodes.", "frequency", string(freq), "count", len(locations))
		shared.Logger.Debug("Derive antinodes.", "requency", string(freq), "antennas", locations)
		for _, perm := range createPermutations(locations) {
			for _, each := range deriveAntinodes(
				perm[0],
				perm[1],
				brd.MaxX,
				brd.MaxY,
				resonantHarmonics) {
				antinodes.Add(each)
			}
		}
	}
	shared.Logger.Debug("Antinodes.", "antinodes", antinodes)
	shared.Logger.Info("Antinodes counted.", "count", antinodes.Count())
	return antinodes
}

func isAntenna(c rune) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	if 'a' <= c && c <= 'z' {
		return true
	}
	if 'A' <= c && c <= 'Z' {
		return true
	}
	return false
}

func createPermutations(locs []shared.Loc) [][]shared.Loc {
	perms := [][]shared.Loc{}
	length := len(locs)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			if i == j {
				continue
			}
			perms = append(perms, []shared.Loc{locs[i], locs[j]})
		}
	}
	return perms
}

func deriveAntinodes(
	a, b shared.Loc,
	maxX, maxY int,
	resonantHarmonics bool,
) []shared.Loc {
	if a == b {
		panic(fmt.Sprintf("Not allowed to have identical locations. Location: %v.", a))
	}
	xDiff := b.X - a.X
	yDiff := b.Y - a.Y
	antinodes := []shared.Loc{}
	if resonantHarmonics {
		antinodes = append(antinodes, a, b)
	}
	// Need to have some kind of a limit to prevent infinite loops.
	for i := 1; i < 1_000; i++ {
		first := shared.Loc{X: a.X - i*xDiff, Y: a.Y - i*yDiff}
		second := shared.Loc{X: b.X + i*xDiff, Y: b.Y + i*yDiff}
		waveNodes := []shared.Loc{}
		for _, each := range []shared.Loc{first, second} {
			if isWithin(0, maxX, each.X) && isWithin(0, maxY, each.Y) {
				waveNodes = append(waveNodes, each)
			}
		}
		if len(waveNodes) == 0 {
			break
		}
		antinodes = append(antinodes, waveNodes...)
		if !resonantHarmonics {
			break
		}
	}
	return antinodes
}

func isWithin(low, high, value int) bool {
	return low <= value && value <= high
}
