package aoc2408

import (
	"fmt"

	"github.com/denarced/advent-of-code/shared"
)

type board struct {
	lines []string
	maxX  int
	maxY  int
}

func newBoard(lines []string) *board {
	brd := &board{lines: lines}
	if len(lines) == 0 {
		return brd
	}
	brd.maxY = len(lines) - 1
	brd.maxX = len([]rune(lines[0])) - 1
	shared.Logger.Info("Board max coordinates.", "x", brd.maxX, "y", brd.maxY)
	return brd
}

// Loc is in proper x-y coordinates.
//
// 2|
// 1|
// 0|
// ..---
// ..012
type iterCb func(loc shared.Loc, c rune)

func (v *board) iter(cb iterCb) {
	lineCount := len(v.lines)
	for y := 0; y < lineCount; y++ {
		line := v.lines[shared.Abs(y-lineCount+1)]
		runes := []rune(line)
		for x := 0; x < len(runes); x++ {
			cb(shared.Loc{X: x, Y: y}, runes[x])
		}
	}
}

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
	brd := newBoard(lines)
	freqToAntennas := map[rune][]shared.Loc{}
	brd.iter(func(loc shared.Loc, c rune) {
		if isAntenna(c) {
			existing, ok := freqToAntennas[c]
			if ok {
				freqToAntennas[c] = append(existing, loc)
			} else {
				freqToAntennas[c] = []shared.Loc{loc}
			}
		}
	})
	for freq, locations := range freqToAntennas {
		shared.Logger.Info("Derive antinodes.", "frequency", string(freq), "count", len(locations))
		shared.Logger.Debug("Derive antinodes.", "requency", string(freq), "antennas", locations)
		for _, perm := range createPermutations(locations) {
			for _, each := range deriveAntinodes(
				perm[0],
				perm[1],
				brd.maxX,
				brd.maxY,
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
