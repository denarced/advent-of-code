package aoc2415

import (
	"fmt"
	"slices"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func CountCoordinateSum(lines []string, doubled bool) int {
	if len(lines) == 0 {
		return 0
	}
	boardLines, directions := splitLines(lines)
	if doubled {
		boardLines = double(boardLines)
	}
	brd := shared.NewBoard(boardLines)
	walk(brd, directions, doubled)
	return countGps(brd, doubled)
}

func double(lines []string) []string {
	doubled := make([]string, 0, len(lines))
	for _, each := range lines {
		line := []rune{}
		for _, c := range each {
			switch c {
			case '@':
				line = append(line, []rune{'@', '.'}...)
			case 'O':
				line = append(line, []rune{'[', ']'}...)
			default:
				line = append(line, []rune{c, c}...)
			}
		}
		doubled = append(doubled, string(line))
	}
	return doubled
}

func countGps(brd *shared.Board, doubled bool) int {
	total := 0
	ch := 'O'
	if doubled {
		ch = '['
	}
	brd.Iter(func(loc shared.Loc, c rune) bool {
		if c == ch {
			yDistance := brd.GetHeight() - 1 - loc.Y
			shared.Logger.Debug("Add to total.", "X", loc.X, "Y distance", yDistance)
			total += loc.X + 100*yDistance
		}
		return true
	})
	return total
}

func findRobot(brd *shared.Board) shared.Loc {
	var robotLoc shared.Loc
	found := false
	brd.Iter(func(loc shared.Loc, c rune) bool {
		if c == '@' {
			robotLoc = loc
			found = true
			return false
		}
		return true
	})
	if !found {
		panic("robo not found")
	}
	return robotLoc
}

func splitLines(lines []string) (board []string, directions []rune) {
	onBoard := true
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		if onBoard {
			if !isDirection(rune(each[0])) {
				board = append(board, trimmed)
				continue
			}
			onBoard = !onBoard
		}
		directions = append(directions, []rune(each)...)
	}
	return
}

func isDirection(c rune) bool {
	return c == '<' || c == '^' || c == '>' || c == 'v'
}

func charToLoc(c rune) shared.Loc {
	if c == '<' {
		return shared.Loc{X: -1, Y: 0}
	}
	if c == '^' {
		return shared.Loc{X: 0, Y: 1}
	}
	if c == '>' {
		return shared.Loc{X: 1, Y: 0}
	}
	if c == 'v' {
		return shared.Loc{X: 0, Y: -1}
	}
	panic(fmt.Sprintf("Unknown direction: %c.", c))
}

// revive:disable-next-line:cognitive-complexity
func walk(brd *shared.Board, directions []rune, doubled bool) {
	robotLoc := findRobot(brd)
	for _, d := range directions {
		shared.Logger.Info("Move robot.", "direction", d, "robot", robotLoc)
		loc := charToLoc(d)
		to := robotLoc.Delta(loc)
		c, ok := brd.Get(to)
		shared.Logger.Debug(
			"About to move robot.",
			"from", robotLoc,
			"to", to,
			"c", string(c),
			"ok", ok,
			"delta", loc,
			"dir", string(d),
		)
		if !ok {
			continue
		}
		if c == '.' {
			swap(brd, robotLoc, to)
			robotLoc = to
			continue
		}
		if c == '#' {
			continue
		}
		if isBox(c, doubled) {
			if !doubled {
				robotLoc = moveRobot(brd, robotLoc, d)
			} else {
				moves := deriveMovedBoxes(brd, robotLoc, d)
				if moves != nil {
					robotLoc = moveRobotAndBoxes(brd, moves)
				}
			}
		}
	}
}

func moveRobotAndBoxes(brd *shared.Board, moves []shared.Pair[shared.Loc]) shared.Loc {
	var robotLoc shared.Loc
	locToC := map[shared.Loc]rune{}
	for _, move := range moves {
		c, _ := brd.Get(move.First)
		locToC[move.First] = c
		brd.Set(move.First, '.')
		if c == '@' {
			robotLoc = move.Second
		}
	}
	for _, move := range moves {
		brd.Set(move.Second, locToC[move.First])
	}
	return robotLoc
}

func moveRobot(brd *shared.Board, robotLoc shared.Loc, direction rune) (rLoc shared.Loc) {
	rLoc = robotLoc
	loc := charToLoc(direction)
	empty, found := findEmpty(brd, rLoc.Delta(loc).Delta(loc), loc)
	emptyC, emptyOk := brd.Get(empty)
	shared.Logger.Debug("Tried to find empty.", "empty", string(emptyC), "emptyOk", emptyOk)
	if !found {
		return
	}
	rev := loc.Rev()
	for empty != rLoc {
		next := empty.Delta(rev)
		swap(brd, empty, next)
		if next == rLoc {
			rLoc = empty
			break
		}
		empty = empty.Delta(rev)
	}
	return
}

func deriveMovedBoxes(
	brd *shared.Board,
	robot shared.Loc,
	direction rune,
) []shared.Pair[shared.Loc] {
	var pairs []shared.Pair[shared.Loc]
	delta := charToLoc(direction)
	start := robot.Delta(delta)
	shared.Logger.Debug("Derive moved boxes.", "robot", robot, "start location", start)
	boxLayers := findBoxLayers(brd, start, delta)
	candBoxLines := deriveBoxLocations(boxLayers, delta, 1)
	if checkOverlap(brd, candBoxLines, delta) {
		return nil
	}
	pairs = append(pairs, createMovePair(robot, delta))
	for _, layer := range boxLayers {
		for _, each := range layer {
			pairs = append(pairs, createMovePair(each, delta))
		}
	}
	return pairs
}

func createMovePair(loc shared.Loc, delta shared.Loc) shared.Pair[shared.Loc] {
	return shared.NewPair(loc, loc.Delta(delta))
}

func findBoxLayers(brd *shared.Board, start, delta shared.Loc) [][]shared.Loc {
	shared.Logger.Debug("Find box layers.", "start", start, "delta", delta)
	first := deriveBoxPair(brd, start)
	layers := [][]shared.Loc{first}
	current := first
	for {
		unique := shared.FilterValues(
			deriveUniqueDeltaLocations(current, delta),
			func(loc shared.Loc) bool {
				c, ok := brd.Get(loc)
				if !ok {
					return false
				}
				return isBox(c, true)
			})
		if len(unique) == 0 {
			break
		}
		candidate := shared.NewSet([]shared.Loc{})
		for _, aUnique := range unique {
			for _, member := range deriveBoxPair(brd, aUnique) {
				candidate.Add(member)
			}
		}
		if candidate.Count() == 0 {
			break
		}
		added := candidate.ToSlice()
		slices.SortFunc(added, func(a, b shared.Loc) int {
			if a.X < b.X {
				return -1
			}
			return 1
		})
		layers = append(layers, added)
		current = added
	}
	return layers
}

func deriveUniqueDeltaLocations(locs []shared.Loc, delta shared.Loc) []shared.Loc {
	unique := shared.NewSet([]shared.Loc{})
	for _, each := range locs {
		unique.Add(each.Delta(delta))
	}
	for _, each := range locs {
		unique.Remove(each)
	}
	result := unique.ToSlice()
	shared.Logger.Debug("Unique delta locations.", "locations", result)
	return result
}

func deriveBoxPair(brd *shared.Board, loc shared.Loc) []shared.Loc {
	shared.Logger.Debug("Derive box pair.", "location", loc)
	c, ok := brd.Get(loc)
	if !ok {
		panic("First loc should be OK.")
	}
	if c == '[' {
		right := loc.Delta(shared.Loc{X: 1, Y: 0})
		c, ok = brd.Get(right)
		if !ok {
			panic("Right should be OK.")
		}
		if c != ']' {
			panic("Right should be ].")
		}
		return []shared.Loc{loc, right}
	}
	if c != ']' {
		panic("Initial right should be ].")
	}
	left := loc.Delta(shared.Loc{X: -1, Y: 0})
	if !ok {
		panic("Left should be OK.")
	}
	if c, ok := brd.Get(left); !ok || c != '[' {
		panic("Left should be OK and [.")
	}
	return []shared.Loc{left, loc}
}

func deriveBoxLocations(boxLayers [][]shared.Loc, delta shared.Loc, steps int) [][]shared.Loc {
	movedLayers := make([][]shared.Loc, 0, len(boxLayers))
	for _, layer := range boxLayers {
		moved := make([]shared.Loc, 0, len(layer))
		for _, each := range layer {
			moved = append(moved, each.Delta(shared.Loc{X: steps * delta.X, Y: steps * delta.Y}))
		}
		movedLayers = append(movedLayers, moved)
	}
	return movedLayers
}

func checkOverlap(brd *shared.Board, layers [][]shared.Loc, delta shared.Loc) bool {
	for _, each := range deriveFront(layers, delta) {
		c, ok := brd.Get(each)
		if !ok {
			shared.Logger.Warn(
				"Stepping outside of board should be impossible.",
				"boxes",
				each,
				"location",
				each,
			)
			return true
		}
		if c != '.' {
			return true
		}
	}
	return false
}

func deriveFront(layers [][]shared.Loc, delta shared.Loc) []shared.Loc {
	allLocs := shared.NewSet([]shared.Loc{})
	for _, outer := range layers {
		for _, inner := range outer {
			allLocs.Add(inner)
		}
	}
	front := []shared.Loc{}
	for _, outer := range layers {
		for _, inner := range outer {
			if !allLocs.Has(inner.Delta(delta)) {
				front = append(front, inner)
			}
		}
	}
	return front
}

func isBox(c rune, doubled bool) bool {
	if doubled {
		return c == '[' || c == ']'
	}
	return c == 'O'
}

func swap(brd *shared.Board, a, b shared.Loc) {
	first, firstOk := brd.Get(a)
	second, secondOk := brd.Get(b)
	if !firstOk || !secondOk {
		shared.Logger.Error(
			"Illegal state when swapping.",
			"a",
			a,
			"b",
			b,
			"first ok",
			firstOk,
			"second ok",
			secondOk,
		)
		panic("Illegal state when swapping.")
	}
	shared.Logger.Debug("Swap", "a", a, "second", string(second), "b", b, "first", string(first))
	brd.Set(b, first)
	brd.Set(a, second)
}

func findEmpty(brd *shared.Board, start shared.Loc, delta shared.Loc) (shared.Loc, bool) {
	curr := start
	for {
		c, ok := brd.Get(curr)
		if !ok {
			shared.Logger.Error(
				"Find empty ran into something impossible.",
				"start",
				start,
				"current",
				curr,
			)
			panic("Borked! Can't find empty.")
		}
		switch c {
		case 'O':
			curr = curr.Delta(delta)
		case '#':
			return shared.Loc{}, false
		case '.':
			return curr, true
		default:
			shared.Logger.Error("Unknown character.", "c", string(c), "current", curr)
			panic("Unknown character.")
		}
	}
}
