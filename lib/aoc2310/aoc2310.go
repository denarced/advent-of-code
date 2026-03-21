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
	case '|', '-', 'S':
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

func FindCrackCount(lines []string) int {
	shared.Logger.Info("Find crack count.", "line count", len(lines))
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	startDirs := findDirections(brd, start)
	chosenDir := startDirs[0]
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Start with board.", "loc", start, "dir", chosenDir)
	}
	aWalker := walker{loc: start, dir: chosenDir}
	walls := gent.NewSet(start)
	var balance int
	for {
		next := step(brd, aWalker)
		if next.loc == start {
			next.dir = chosenDir
			break
		}
		balance += deriveTurn(aWalker.dir, next.dir)
		walls.Add(next.loc)
		aWalker = next
	}
	// Negative balance: turned more left than right so counter clockwise, and therefore inside of
	// shape is on the left.
	if balance == 0 {
		panic("balance can't be zero, impossible")
	}
	shared.Logger.Info("First round done.", "balance", balance)
	aWalker = walker{loc: start, dir: chosenDir}
	squeezed := gent.NewSet[shared.Loc]()
	excluded := gent.NewSet[shared.Loc]()
	callFillUps := func(w walker, dir shared.Direction) {
		for _, each := range []struct {
			set *gent.Set[shared.Loc]
			bal int
		}{
			{squeezed, balance},
			{excluded, -balance},
		} {
			fillUp(brd, walls, each.set, w, dir, each.bal)
		}
	}
	callFillUps(aWalker, aWalker.dir)
	for {
		next := step(brd, aWalker)
		if next.loc == start {
			break
		}
		callFillUps(next, aWalker.dir)
		aWalker = next
	}
	if containCommon(squeezed, excluded) {
		panic("common in squeezed and excluded")
	}
	if containCommon(squeezed, walls) {
		panic("common between squeezed and walls")
	}
	if containCommon(excluded, walls) {
		panic("common in walls and excluded")
	}

	squeezeCount := squeezed.Len()
	wallCount := walls.Len()
	excludedCount := excluded.Len()
	total := squeezeCount + wallCount + excludedCount
	shared.Logger.Info(
		"Got crack count.",
		"count", squeezeCount,
		"wall count", wallCount,
		"excluded count", excludedCount,
		"sum count", total,
		"total count", brd.GetArea(),
	)
	if total != brd.GetArea() {
		panic("count mismatch between board and counts")
	}

	return squeezed.Len()
}

func deriveTurn(before, after shared.Direction) int {
	deriveDirValue := func(left, right shared.Direction) int {
		switch after {
		case left:
			return -1
		case right:
			return 1
		default:
			return 0
		}
	}
	switch before {
	case shared.RealEast:
		return deriveDirValue(shared.RealNorth, shared.RealSouth)
	case shared.RealSouth:
		return deriveDirValue(shared.RealEast, shared.RealWest)
	case shared.RealWest:
		return deriveDirValue(shared.RealSouth, shared.RealNorth)
	case shared.RealNorth:
		return deriveDirValue(shared.RealWest, shared.RealEast)
	default:
		panic("invalid before dir")
	}
}

func fillUp(
	brd *shared.Board,
	walls, squuezed *gent.Set[shared.Loc],
	aWalker walker,
	origignalDir shared.Direction,
	balance int,
) {
	turns := []shared.Direction{aWalker.dir}
	if aWalker.dir != origignalDir {
		turns = append(turns, origignalDir)
	}

	var dirs []shared.Direction
	for _, each := range turns {
		dir := gent.Tri(balance > 0, each.TurnRealRight(), each.TurnRealLeft())
		dirs = append(dirs, dir)
	}

	var firstCandidates []shared.Loc
	for _, each := range dirs {
		firstCandidates = append(firstCandidates, aWalker.loc.Delta(shared.Loc(each)))
	}

	var stack []shared.Loc
	for _, each := range firstCandidates {
		if walls.Has(each) || squuezed.Has(each) {
			continue
		}
		if _, ok := brd.Get(each); !ok {
			continue
		}
		squuezed.Add(each)
		stack = append(stack, each)
	}
	for len(stack) > 0 {
		block := stack[0]
		if len(stack) == 1 {
			stack = nil
		} else {
			stack = stack[1:]
		}
		for _, each := range shared.RealPrimaryDirections {
			cand := block.Delta(shared.Loc(each))
			if _, ok := brd.Get(cand); !ok {
				continue
			}
			if walls.Has(cand) || squuezed.Has(cand) {
				continue
			}
			squuezed.Add(cand)
			stack = append(stack, cand)
		}
	}
}

func containCommon(a, b *gent.Set[shared.Loc]) (common bool) {
	a.ForEach(func(loc shared.Loc, stop func()) {
		if b.Has(loc) {
			common = true
			stop()
		}
	})
	return
}
