package aoc2406

import (
	"fmt"
	"strings"
	"sync"

	"github.com/denarced/advent-of-code/shared"
)

type vector struct {
	loc shared.Loc
	dir shared.Direction
}

type fatBoard struct {
	curr      vector
	visited   *shared.Set[vector]
	stepped   *shared.Set[shared.Loc]
	nestedBrd *shared.Board
}

func newFatBoard(init vector, nestedBrd *shared.Board) *fatBoard {
	return &fatBoard{
		curr:      init,
		visited:   shared.NewSet([]vector{init}),
		stepped:   shared.NewSet([]shared.Loc{init.loc}),
		nestedBrd: nestedBrd,
	}
}

func (v *fatBoard) deriveNextLocation() shared.Loc {
	return shared.Loc{
		X: v.curr.loc.X + v.curr.dir.X,
		Y: v.curr.loc.Y + v.curr.dir.Y,
	}
}

func (v *fatBoard) deriveNextVector() vector {
	return vector{
		dir: v.curr.dir,
		loc: shared.Loc{
			X: v.curr.loc.X + v.curr.dir.X,
			Y: v.curr.loc.Y + v.curr.dir.Y,
		},
	}
}

func (v *fatBoard) isInside(l shared.Loc) bool {
	_, ok := v.nestedBrd.Get(l)
	return ok
}

func (v *fatBoard) turn() {
	v.visited.Add(v.curr)
	v.stepped.Add(v.curr.loc)
	v.curr.dir = v.curr.dir.TurnRealRight()
}

func (v *fatBoard) isBlock(l shared.Loc) bool {
	return v.nestedBrd.GetOrDie(l) == '#'
}

func (v *fatBoard) move(l shared.Loc) {
	v.curr.loc = l
	v.visited.Add(v.curr)
	v.stepped.Add(l)
}

func findIndefiniteBlock(
	brd *fatBoard,
	possible shared.Loc,
	wg *sync.WaitGroup,
	locCh chan<- shared.Loc,
) {
	defer wg.Done()
	brd.nestedBrd.Set(possible, '#')
	shared.Logger.Debug("Find-indef.", "curr", brd.curr)
	turnCount := 0
	for i := 0; i < 100*shared.Max(brd.nestedBrd.GetWidth(), brd.nestedBrd.GetHeight()); i++ {
		if brd.visited.Has(brd.deriveNextVector()) {
			shared.Logger.Info("Found block that causes indefinite loop.", "block", possible)
			locCh <- possible
			return
		}

		next := brd.deriveNextLocation()
		if !brd.isInside(next) {
			return
		}
		if brd.isBlock(next) {
			brd.turn()
			turnCount++
			if turnCount <= 4 {
				continue
			}
			locCh <- possible
			return
		}
		turnCount = 0
		brd.move(next)
	}
	panic("indef loop")
}

func (v *fatBoard) copy() *fatBoard {
	return newFatBoard(v.curr, v.nestedBrd.Copy())
}

func (v *fatBoard) deriveVisitedCount() int {
	return v.stepped.Count()
}

func (v *fatBoard) print() string {
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
		v.nestedBrd.Set(loc, deriveDirCharacter(dirs))
	}
	v.nestedBrd.Set(v.curr.loc, '*')
	return strings.Join(v.nestedBrd.GetLines(), "\n") + "\n"
}

func CountDistinctPositions(lines []string) int {
	if len(lines) == 0 {
		return 0
	}
	shared.Logger.Info("Count distinct positions.", "line count", len(lines))
	brd := shared.NewBoard(lines)
	fatBrd := newFatBoard(vector{loc: brd.FindOrDie('^'), dir: shared.RealNorth}, brd)
	counter := 7_000
	for {
		next := fatBrd.deriveNextLocation()
		if !fatBrd.isInside(next) {
			break
		}
		if fatBrd.isBlock(next) {
			fatBrd.turn()
			continue
		}
		shared.Logger.Debug("Step.", "previous", fatBrd.curr, "next", next)
		shared.Logger.Debug(
			"Move.",
			"step",
			fmt.Sprintf("%s -> %s", fatBrd.curr.loc.ToString(), next.ToString()),
		)
		fatBrd.move(next)
		counter--
		if counter < 0 {
			panic("This loop is clearly eternal.")
		}
	}
	return fatBrd.deriveVisitedCount()
}

func CountBlocksForIndefiniteLoops(lines []string) *shared.Set[shared.Loc] {
	if len(lines) == 0 {
		return shared.NewSet([]shared.Loc{})
	}
	shared.Logger.Info("Derive infinite loop locations.")
	nested := shared.NewBoard(lines)
	brd := newFatBoard(vector{loc: nested.FindOrDie('^'), dir: shared.RealNorth}, nested)
	indefLocations := shared.NewSet([]shared.Loc{})
	locCh := make(chan shared.Loc)
	finishedCh := make(chan int)
	go func() {
		for each := range locCh {
			indefLocations.Add(each)
		}
		finishedCh <- 0
	}()
	var wg sync.WaitGroup
	for {
		next := brd.deriveNextLocation()
		shared.Logger.Debug("Now.", "location", brd.curr, "next", next)
		// Can't place a block on path already travelled, it would invalidate the past.
		if brd.isInside(next) && !brd.stepped.Has(next) {
			wg.Add(1)
			go findIndefiniteBlock(brd.copy(), next, &wg, locCh)
		}

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
	wg.Wait()
	close(locCh)
	<-finishedCh
	shared.Logger.Info("Indefinite blocks.", "locations", indefLocations)
	return indefLocations
}

func deriveDirCharacter(dirs []shared.Direction) rune {
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
		return '-'
	}
	if !hor && ver {
		return '|'
	}
	if hor && ver {
		return '+'
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
