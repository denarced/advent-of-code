package aoc2406

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

type vector struct {
	loc shared.Loc
	dir shared.Direction
}

type board struct {
	curr    vector
	blocks  *shared.Set[shared.Loc]
	visited *shared.Set[vector]
	stepped *shared.Set[shared.Loc]
	width   int
	height  int
}

func newBoard(init shared.Loc, blocks []shared.Loc, width, height int) *board {
	curr := vector{loc: init, dir: shared.DirNorth}
	return &board{
		curr:    curr,
		blocks:  shared.NewSet(blocks),
		visited: shared.NewSet([]vector{curr}),
		stepped: shared.NewSet([]shared.Loc{curr.loc}),
		width:   width,
		height:  height,
	}
}

func (v *board) deriveNextLocation() shared.Loc {
	return shared.Loc{
		X: v.curr.loc.X + v.curr.dir.Y,
		Y: v.curr.loc.Y + v.curr.dir.X,
	}
}

func (v *board) deriveNextVector() vector {
	return vector{
		dir: v.curr.dir,
		loc: shared.Loc{
			X: v.curr.loc.X + v.curr.dir.Y,
			Y: v.curr.loc.Y + v.curr.dir.X,
		},
	}
}

func (v *board) isInside(l shared.Loc) bool {
	if l.X < 0 || l.Y < 0 {
		return false
	}
	if l.X >= v.height {
		return false
	}
	return l.Y < v.width
}

func (v *board) turn() {
	v.visited.Add(v.curr)
	v.stepped.Add(v.curr.loc)
	v.curr.dir = v.curr.dir.TurnRight()
}

func (v *board) isBlock(l shared.Loc) bool {
	return v.blocks.Has(l)
}

func (v *board) move(l shared.Loc) {
	v.curr.loc = l
	v.visited.Add(v.curr)
	v.stepped.Add(l)
}

func (v *board) findIndefiniteBlock() bool {
	w := v.copy()
	possible := w.deriveNextLocation()
	// Can't place a block on path already travelled, it would invalidate the past.
	if v.stepped.Has(possible) {
		return false
	}
	w.blocks.Add(possible)
	shared.Logger.Debug("Find-indef.", "curr", w.curr)
	turnCount := 0
	for i := 0; i < 100*shared.Max(v.width, v.height); i++ {
		if w.visited.Has(w.deriveNextVector()) {
			shared.Logger.Info("Found block that causes indefinite loop.", "block", possible)
			return true
		}

		next := w.deriveNextLocation()
		if !w.isInside(next) {
			return false
		}
		if w.isBlock(next) {
			w.turn()
			turnCount++
			if turnCount <= 4 {
				continue
			}
			return true
		}
		turnCount = 0
		w.move(next)
	}
	panic("indef loop")
}

func (v *board) copy() *board {
	return &board{
		curr:    v.curr,
		blocks:  v.blocks.Copy(),
		visited: v.visited.Copy(),
		stepped: v.stepped.Copy(),
		width:   v.width,
		height:  v.height,
	}
}

func (v *board) deriveVisitedCount() int {
	return v.stepped.Count()
}

func (v *board) print() string {
	lines := []string{}
	for range v.height {
		lines = append(lines, strings.Repeat(" ", v.width))
	}
	v.blocks.Iter(func(item shared.Loc) bool {
		setCharacter(lines, item, "#")
		return true
	})
	locToDirs := map[shared.Loc][]shared.Direction{}
	v.visited.Iter(func(v vector) bool {
		l := v.loc
		if dirs, ok := locToDirs[l]; ok {
			locToDirs[l] = append(dirs, v.dir)
		} else {
			locToDirs[l] = []shared.Direction{v.dir}
		}
		return true
	})
	for loc, dirs := range locToDirs {
		setCharacter(lines, loc, deriveDirCharacter(dirs))
	}
	setCharacter(lines, v.curr.loc, "*")
	return strings.Join(lines, "\n") + "\n"
}

func CountDistinctPositions(lines []string) int {
	if len(lines) == 0 {
		return 0
	}
	shared.Logger.Info("Count distinct positions.", "line count", len(lines))
	brd := newBoard(
		findCharacter(lines, '^'),
		findLocations(lines, '#'),
		len(lines[0]),
		len(lines))
	counter := 7_000
	for {
		next := brd.deriveNextLocation()
		if !brd.isInside(next) {
			break
		}
		if brd.isBlock(next) {
			brd.turn()
			continue
		}
		shared.Logger.Debug("Step.", "previous", brd.curr, "next", next)
		shared.Logger.Debug(
			"Move.",
			"step",
			fmt.Sprintf("%s -> %s", brd.curr.loc.ToString(), next.ToString()),
		)
		brd.move(next)
		counter--
		if counter < 0 {
			panic("This loop is clearly eternal.")
		}
	}
	return brd.deriveVisitedCount()
}

func CountBlocksForIndefiniteLoops(lines []string) *shared.Set[shared.Loc] {
	if len(lines) == 0 {
		return shared.NewSet([]shared.Loc{})
	}
	shared.Logger.Info("Derive infinite loop locations.")
	brd := newBoard(
		findCharacter(lines, '^'),
		findLocations(lines, '#'),
		len(lines[0]),
		len(lines))
	indefLocations := shared.NewSet([]shared.Loc{})
	for {
		shared.Logger.Debug("Now.", "location", brd.curr)
		if brd.findIndefiniteBlock() {
			next := brd.deriveNextLocation()
			indefLocations.Add(next)
		}

		next := brd.deriveNextLocation()
		if !brd.isInside(next) {
			shared.Logger.Info("Guard left the area.")
			break
		}
		if brd.isBlock(next) {
			brd.turn()
			continue
		}
		brd.move(next)
	}
	shared.Logger.Info("Indefinite blocks.", "locations", indefLocations)
	return indefLocations
}

func findLocations(lines []string, c byte) []shared.Loc {
	locs := []shared.Loc{}
	for r := 0; r < len(lines); r++ {
		for l := 0; l < len(lines[r]); l++ {
			if lines[r][l] == c {
				locs = append(locs, shared.Loc{X: r, Y: l})
			}
		}
	}
	return locs
}

func findCharacter(lines []string, char byte) shared.Loc {
	for r := 0; r < len(lines); r++ {
		for c := 0; c < len(lines[r]); c++ {
			if lines[r][c] == char {
				return shared.Loc{X: r, Y: c}
			}
		}
	}
	return shared.Loc{X: -1, Y: -1}
}

func setCharacter(lines []string, loc shared.Loc, char string) {
	if char == "" || len(char) > 1 {
		panic(fmt.Sprintf("Invalid character length(%d): \"%s\".", len(char), char))
	}
	line := lines[loc.X]
	line = line[0:loc.Y] + char + line[loc.Y+1:]
	lines[loc.X] = line
}

func deriveDirCharacter(dirs []shared.Direction) string {
	hor, ver := false, false
	shared.NewSet(dirs).Iter(func(d shared.Direction) bool {
		if d.X == 0 && d.Y != 0 {
			ver = true
		} else if d.X != 0 && d.Y == 0 {
			hor = true
		} else {
			panic(fmt.Sprintf("Only horizontal or vertical are allowed: %v.", d))
		}
		return true
	})
	if hor && !ver {
		return "-"
	}
	if !hor && ver {
		return "|"
	}
	if hor && ver {
		return "+"
	}
	panic(
		fmt.Sprintf(
			"Illegal state: should be hor and/or ver. Hor: %t. Ver: %t. Dirs: %v.",
			hor,
			ver,
			dirs,
		),
	)
}
