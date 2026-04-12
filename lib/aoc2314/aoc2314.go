package aoc2314

import (
	"slices"

	"github.com/denarced/advent-of-code/shared"
)

const (
	roundRock  rune = 'O'
	emptySpace rune = '.'
)

func CountTotalLoad(lines []string) int {
	shared.Logger.Info("Count total load.")
	brd := shared.NewBoard(lines)
	moveRocksNorth(brd)
	var weight int
	shared.Logger.Info("Count weight.")
	brd.Iter(func(loc shared.Loc, c rune) (keepGoing bool) {
		keepGoing = true
		if c != roundRock {
			return
		}
		weight += loc.Y + 1
		return
	})
	shared.Logger.Info("Weight counted.", "weight", weight)
	return weight
}

func moveRocksNorth(brd *shared.Board) {
	shared.Logger.Info("Move rocks north.")
	var rocks []shared.Loc
	brd.Iter(func(loc shared.Loc, c rune) (keepGoing bool) {
		keepGoing = true
		if c != roundRock {
			return
		}
		rocks = append(rocks, loc)
		return
	})
	// Sort so that rocks on the left are first, higher rocks first. So ascending order for X and
	// descending for Y. The former doesn't matter (just need some order), the latter matters.
	slices.SortFunc(rocks, func(a, b shared.Loc) int {
		xDiff := a.X - b.X
		if xDiff != 0 {
			return xDiff
		}
		return b.Y - a.Y
	})
	maxY := brd.GetHeight()
	var moveCount, rockCount int
	for _, each := range rocks {
		var moved bool
		for {
			if each.Y+1 >= maxY {
				break
			}
			newLoc := each.Delta(shared.Loc(shared.RealNorth))
			if brd.GetOrDie(newLoc) != emptySpace {
				break
			}
			moved = true
			brd.Set(newLoc, roundRock)
			brd.Set(each, emptySpace)
			shared.Logger.Debug("Move rock.", "from", each, "to", newLoc)
			each = newLoc
			moveCount++
		}
		if moved {
			rockCount++
		}
	}
	shared.Logger.Info("Rocks moved.", "rock count", rockCount, "move count", moveCount)
}
