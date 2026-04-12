package aoc2314

import (
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
	rocks := make([]shared.Loc, brd.GetArea()/5)
	for x := range brd.GetWidth() {
		for y := brd.GetHeight() - 1; y >= 0; y-- {
			loc := shared.Loc{X: x, Y: y}
			if brd.GetOrDie(loc) == roundRock {
				rocks = append(rocks, loc)
			}
		}
	}
	var moveCount, rockCount int
	for _, each := range rocks {
		var moved bool
		dest := each
		for {
			cand := dest.Delta(shared.Loc(shared.RealNorth))
			if c, ok := brd.Get(cand); ok && c == emptySpace {
				dest = cand
				moved = true
				moveCount++
				continue
			}
			break
		}
		if !moved {
			continue
		}
		rockCount++
		brd.Set(dest, roundRock)
		brd.Set(each, emptySpace)
		shared.Logger.Debug("Move rock.", "from", each, "to", dest)
	}
	shared.Logger.Info("Rocks moved.", "rock count", rockCount, "move count", moveCount)
}
