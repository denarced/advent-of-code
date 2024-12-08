// Package aoc2024 contains implementation for 2024 Advent of Code solutions.
package aoc2024

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func ToInts(s []string) (nums []int, err error) {
	shared.Logger.Info("Convert string slice to ints.", "length", len(s))
	for _, each := range s {
		var n int
		n, err = strconv.Atoi(each)
		if err != nil {
			shared.Logger.Error("Failed to convert to int.", "string", each, "err", err)
			return
		}
		nums = append(nums, n)
	}
	return
}

type vector struct {
	loc shared.Location
	dir shared.Direction
}

func SumCorrectMiddlePageNumbers(lines []string) int {
	return sumMiddlePageNumbers(lines, true)
}

func SumIncorrectMiddlePageNumbers(lines []string) int {
	return sumMiddlePageNumbers(lines, false)
}

func sumMiddlePageNumbers(lines []string, correct bool) int {
	rules, pages := toRulesAndPages(lines)
	shared.Logger.Info(
		"Sum middle page numbers.",
		"rule count",
		len(rules),
		"page list count",
		len(pages),
		"correct",
		correct,
	)
	sum := 0
	filtered := filterValues(
		pages,
		func(s []int) bool {
			return isSortedAccordingToRules(rules, s) == correct
		})
	shared.Logger.Info(
		"Page lists filtered.",
		"page list count",
		len(filtered))
	for _, each := range filtered {
		if !correct {
			each = sortWithRules(rules, each)
		}
		middle := each[len(each)/2]
		shared.Logger.Debug("Add middle page to the sum.", "pages", each, "middle", middle)
		sum += middle
	}
	shared.Logger.Info("Middle page numbers summed.", "sum", sum)
	return sum
}

func toRulesAndPages(lines []string) ([][]int, [][]int) {
	toInts := func(s []string) []int {
		var ints []int
		for _, each := range s {
			i, err := strconv.Atoi(each)
			if err != nil {
				panic(err)
			}
			ints = append(ints, i)
		}
		return ints
	}
	contains := func(sub string) func(s string) bool {
		return func(s string) bool {
			return strings.Contains(s, sub)
		}
	}
	split := func(sep string) func(s string) []string {
		return func(s string) []string {
			return strings.Split(s, sep)
		}
	}
	filterAndSplit := func(sep string) [][]int {
		return mapValues(
			mapValues(
				filterValues(lines, contains(sep)),
				split(sep)),
			toInts)
	}
	return filterAndSplit("|"), filterAndSplit(",")
}

func mapValues[T any, U any](s []T, f func(v T) U) []U {
	var result []U
	for _, each := range s {
		result = append(result, f(each))
	}
	return result
}

func filterValues[T any](s []T, f func(v T) bool) []T {
	var result []T
	for _, each := range s {
		if f(each) {
			result = append(result, each)
		}
	}
	return result
}

func isSortedAccordingToRules(rules [][]int, pages []int) bool {
	sorted := sortWithRules(rules, pages)
	return slices.Equal(pages, sorted)
}

func sortWithRules(rules [][]int, pages []int) []int {
	dup := append([]int{}, pages...)
	slices.SortStableFunc(
		dup,
		func(a, b int) int {
			rule := findRelevantRule(rules, a, b)
			if rule == nil {
				return 0
			}
			if rule[0] == a {
				return -1
			}
			return 1
		})
	return dup
}

func findRelevantRule(rules [][]int, a, b int) []int {
	var relevant []int
	for _, each := range rules {
		if a == each[0] && b == each[1] || a == each[1] && b == each[0] {
			return each
		}
	}
	return relevant
}

type board struct {
	curr    vector
	blocks  *shared.Set[shared.Location]
	visited *shared.Set[vector]
	stepped *shared.Set[shared.Location]
	width   int
	height  int
}

func newBoard(init shared.Location, blocks []shared.Location, width, height int) *board {
	curr := vector{loc: init, dir: shared.DirNorth}
	return &board{
		curr:    curr,
		blocks:  shared.NewSet(blocks),
		visited: shared.NewSet([]vector{curr}),
		stepped: shared.NewSet([]shared.Location{curr.loc}),
		width:   width,
		height:  height,
	}
}

func (v *board) deriveNextLocation() shared.Location {
	return shared.Location{
		X: v.curr.loc.X + v.curr.dir.Y,
		Y: v.curr.loc.Y + v.curr.dir.X,
	}
}

func (v *board) deriveNextVector() vector {
	return vector{
		dir: v.curr.dir,
		loc: shared.Location{
			X: v.curr.loc.X + v.curr.dir.Y,
			Y: v.curr.loc.Y + v.curr.dir.X,
		},
	}
}

func (v *board) isInside(l shared.Location) bool {
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

func (v *board) isBlock(l shared.Location) bool {
	return v.blocks.Has(l)
}

func (v *board) move(l shared.Location) {
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
	v.blocks.Iter(func(item shared.Location) bool {
		setCharacter(lines, item, "#")
		return true
	})
	locToDirs := map[shared.Location][]shared.Direction{}
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

func CountBlocksForIndefiniteLoops(lines []string) *shared.Set[shared.Location] {
	if len(lines) == 0 {
		return shared.NewSet([]shared.Location{})
	}
	shared.Logger.Info("Derive infinite loop locations.")
	brd := newBoard(
		findCharacter(lines, '^'),
		findLocations(lines, '#'),
		len(lines[0]),
		len(lines))
	indefLocations := shared.NewSet([]shared.Location{})
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

func findLocations(lines []string, c byte) []shared.Location {
	locs := []shared.Location{}
	for r := 0; r < len(lines); r++ {
		for l := 0; l < len(lines[r]); l++ {
			if lines[r][l] == c {
				locs = append(locs, shared.Location{X: r, Y: l})
			}
		}
	}
	return locs
}

func findCharacter(lines []string, char byte) shared.Location {
	for r := 0; r < len(lines); r++ {
		for c := 0; c < len(lines[r]); c++ {
			if lines[r][c] == char {
				return shared.Location{X: r, Y: c}
			}
		}
	}
	return shared.Location{X: -1, Y: -1}
}

func setCharacter(lines []string, loc shared.Location, char string) {
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
