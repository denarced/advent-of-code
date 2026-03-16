package aoc2310

import (
	"fmt"
	"slices"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

type walker struct {
	loc shared.Loc
	dir shared.Direction
}

func CountSteps(lines []string) int {
	shared.Logger.Info("Count steps.", "line count", len(lines))
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	startDirs := findDirections(brd, start)
	if len(startDirs) != 2 {
		shared.Logger.Error(
			"Invalid board, should have 2 directions from S.",
			"start dir count",
			len(startDirs),
		)
		panic("invalid board")
	}
	walkers := make([]walker, len(startDirs))
	for i, each := range startDirs {
		walkers[i] = walker{loc: start, dir: each}
	}
	count := 0
	next := make([]walker, len(walkers))
	for {
		for i, each := range walkers {
			next[i] = step(brd, each)
		}
		count++
		if shared.IsDebugEnabled() {
			for i := range walkers {
				shared.Logger.Debug("Moved.", "i", i, "from", walkers[i], "to", next[i])
			}
		}
		locs := gent.Map(next, func(w walker) shared.Loc {
			return w.loc
		})
		if allEqual(locs) || (locs[0] == walkers[1].loc && locs[1] == walkers[0].loc) {
			shared.Logger.Info("Got the count.", "count", count)
			return count
		}
		walkers = next
	}
}

func allEqual[S ~[]T, T comparable](s S) bool {
	if len(s) == 0 {
		return true
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if first != s[i] {
			return false
		}
	}
	return true
}

func findDirections(brd *shared.Board, loc shared.Loc) []shared.Direction {
	allowed := [][]shared.Direction{
		{shared.RealNorth, shared.RealSouth},
		{shared.RealWest, shared.RealEast},
		{shared.RealSouth, shared.RealWest},
		{shared.RealSouth, shared.RealEast},
		{shared.RealEast, shared.RealNorth},
		{shared.RealWest, shared.RealNorth},
	}
	runes := []rune{'|', '-', 'L', 'J', '7', 'F'}
	near := brd.NextTo(loc, runes, false)
	shared.Logger.Debug("Near.", "start", loc, "near", near)
	var valid []shared.Direction
	for _, spot := range near {
		i := slices.Index(runes, brd.GetOrDie(spot))
		for _, each := range allowed[i] {
			cand := loc.Delta(shared.Loc(each))
			if cand == spot {
				valid = append(valid, each)
			}
		}
	}
	return valid
}

func step(brd *shared.Board, aWalker walker) (result walker) {
	nextLoc := aWalker.loc.Delta(shared.Loc(aWalker.dir))
	result.loc = nextLoc
	directions := map[rune]map[shared.Direction]shared.Direction{
		'L': {
			shared.RealSouth: shared.RealEast,
			shared.RealWest:  shared.RealNorth,
		},
		'7': {
			shared.RealEast:  shared.RealSouth,
			shared.RealNorth: shared.RealWest,
		},
		'J': {
			shared.RealSouth: shared.RealWest,
			shared.RealEast:  shared.RealNorth,
		},
		'F': {
			shared.RealWest:  shared.RealSouth,
			shared.RealNorth: shared.RealEast,
		},
	}
	c := brd.GetOrDie(nextLoc)
	switch c {
	case '|', '-':
		result.dir = aWalker.dir
	default:
		matching, ok := directions[c]
		if !ok {
			panic(fmt.Sprintf("no such character in map: %s", string(c)))
		}
		result.dir = matching[aWalker.dir]
	}
	return
}
