package aoc2321

import (
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func CountRangeFromLines(lines []string, stepCount int) int {
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	brd.Set(start, '.')
	walkers := gent.NewSet(start)
	for range stepCount {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Start round.", "walker count", walkers.Count())
		}
		next := gent.NewSet[shared.Loc]()
		walkers.ForEachAll(func(loc shared.Loc) {
			for _, each := range brd.NextTo(loc, '.', false) {
				next.Add(each)
			}
		})
		walkers = next
	}
	result := walkers.Count()
	shared.Logger.Info("Range counted.", "result", result)
	return result
}
