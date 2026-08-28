package aoc2208

import "github.com/denarced/advent-of-code/shared"

func DeriveVisibleTreeCount(lines []string) int {
	brd := shared.NewBoard(lines)
	seen := createSeen(brd, brd.GetWidth(), brd.GetHeight())
	if shared.IsDebugEnabled() {
		shared.Logger.Debug(
			"Seen initialized.",
			"width", brd.GetWidth(),
			"height", brd.GetHeight())
	}
	for y := 1; y < brd.GetHeight()-1; y++ {
		dig(brd, seen, createGenerator(0, 1, brd.GetWidth()-1), createFixedGenerator(y))
		dig(brd, seen, createGenerator(brd.GetWidth()-1, -1, 0), createFixedGenerator(y))
	}
	for x := 1; x < brd.GetWidth()-1; x++ {
		dig(brd, seen, createFixedGenerator(x), createGenerator(0, 1, brd.GetHeight()-1))
		dig(brd, seen, createFixedGenerator(x), createGenerator(brd.GetHeight()-1, -1, 0))
	}
	var total int
	for y := 0; y < brd.GetHeight(); y++ {
		for x := 0; x < brd.GetWidth(); x++ {
			if seen[y][x] {
				total++
			}
		}
	}
	return total
}

type genFunc func() (int, bool)

func createGenerator(alpha, delta, omega int) genFunc {
	return func() (int, bool) {
		next := alpha
		valid := alpha != omega
		alpha += delta
		return next, valid
	}
}

func createFixedGenerator(value int) genFunc {
	return func() (int, bool) {
		return value, true
	}
}

func dig(brd *shared.Board, seen [][]bool, xGen, yGen genFunc) {
	prevLoc := shared.Loc{X: -1, Y: -1}
	for {
		x, xValid := xGen()
		y, yValid := yGen()
		if !xValid || !yValid {
			break
		}
		currLoc := shared.Loc{X: x, Y: y}
		var prevValue, currValue int
		if prevLoc.X >= 0 {
			prevValue = brd.GetIntOrDie(prevLoc)
			currValue = brd.GetIntOrDie(currLoc)
			if shared.IsDebugEnabled() {
				shared.Logger.Debug(
					"Check value.",
					"prev loc", prevLoc,
					"prev value", prevValue,
					"curr loc", currLoc,
					"curr value", currValue)
			}
			if prevValue < currValue {
				seen[y][x] = true
				shared.Logger.Debug("Match.")
			}
		}

		if prevValue == 0 && prevValue == currValue || prevValue < currValue {
			prevLoc = currLoc
		}
	}
}

func createSeen(brd *shared.Board, width, height int) [][]bool {
	seen := make([][]bool, brd.GetHeight())
	for i := range brd.GetHeight() {
		seen[i] = make([]bool, brd.GetWidth())
	}
	if shared.IsDebugEnabled() {
		shared.Logger.Debug(
			"Seen initialized.",
			"width", brd.GetWidth(),
			"height", brd.GetHeight())
	}
	for x := range brd.GetWidth() {
		// Bottom
		seen[0][x] = true
		// Top
		seen[brd.GetHeight()-1][x] = true
	}
	for y := range brd.GetHeight() {
		// Left
		seen[y][0] = true
		// Right
		seen[y][brd.GetWidth()-1] = true
	}
	return seen
}
