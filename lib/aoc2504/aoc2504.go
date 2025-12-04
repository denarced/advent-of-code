package aoc2504

import "github.com/denarced/advent-of-code/shared"

func CountRolls(lines []string, tries int) int {
	board := shared.NewBoard(lines)
	var count int
	for i := tries; i != 0; i-- {
		prev := count
		removables := deriveMovableRolls(board)
		count += len(removables)
		for _, each := range removables {
			board.Set(each, 'x')
		}
		if prev == count {
			break
		}
	}
	return count
}

func deriveMovableRolls(board *shared.Board) []shared.Loc {
	var locs []shared.Loc
	board.Iter(func(loc shared.Loc, c rune) bool {
		if c != '@' {
			return true
		}
		rollCount := len(board.NextTo(loc, '@', true))
		if rollCount < 4 {
			locs = append(locs, loc)
		}
		return true
	})
	return locs
}
