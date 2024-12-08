// Package aoc2024 contains implementation for 2024 Advent of Code solutions.
package aoc2024

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

var (
	mulPattern = regexp.MustCompile(`mul\(\d+,\d+\)|don't\(\)|do\(\)`)
	dirNorth   = direction{0, -1}
	dirEast    = direction{1, 0}
	dirSouth   = direction{0, 1}
	dirWest    = direction{-1, 0}
	directions = []direction{
		dirEast,
		{1, -1},
		dirNorth,
		{-1, -1},
		dirWest,
		{-1, 1},
		dirSouth,
		{1, 1},
	}
)

func Advent01Distance(left, right []int) int {
	shared.Logger.Info(
		"Advent 01: derive distance.",
		"left length",
		len(left),
		"right length",
		len(right),
	)
	slices.Sort(left)
	slices.Sort(right)
	distance := 0
	for i := range len(left) {
		distance += shared.Abs(left[i] - right[i])
		if distance < 0 {
			panic("Distance is <0 (int overflow).")
		}
	}
	shared.Logger.Info("Distance derived.", "distance", distance)
	return distance
}

func Advent01Similarity(left, right []int) int {
	shared.Logger.Info(
		"Advent 01: derive similarity",
		"left length",
		len(left),
		"right length",
		len(right),
	)
	counts := deriveCounts(right)
	similarity := 0
	for _, each := range left {
		c, ok := counts[each]
		if ok {
			similarity += (c * each)
		}
		if similarity < 0 {
			panic("Similarity is <0 (int overflow).")
		}
	}
	shared.Logger.Info("Similarity derived.", "similarity", similarity)
	return similarity
}

func deriveCounts(s []int) map[int]int {
	counts := map[int]int{}
	for _, each := range s {
		curr, ok := counts[each]
		if ok {
			counts[each] = curr + 1
		} else {
			counts[each] = 1
		}
	}
	return counts
}

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

func ToIntTable(s []string) (table [][]int) {
	for _, each := range s {
		cells := strings.Fields(each)
		var row []int
		for _, c := range cells {
			n, err := strconv.Atoi(c)
			if err != nil {
				panic("Invalid number: " + c)
			}
			row = append(row, n)
		}
		if row == nil {
			panic("Empty row")
		}
		table = append(table, row)
	}
	return
}

func CountSafe(levels [][]int, dampener bool) int {
	count := 0
	for _, each := range levels {
		if len(each) < 2 {
			continue
		}
		index := deriveUnsafe(each)
		if index < 0 {
			count++
			continue
		}
		if !dampener {
			continue
		}

		for i := index + 1; i >= 0; i-- {
			trimmed := append([]int{}, each[0:i]...)
			trimmed = append(trimmed, each[i+1:]...)
			index = deriveUnsafe(trimmed)
			if index < 0 {
				shared.Logger.Debug("Safe after trimmed.", "original", each, "trimmed", trimmed)
				break
			}
		}
		if index < 0 {
			count++
		}
	}
	return count
}

// Derive index where unsafe was detected or -1 if levels are safe.
func deriveUnsafe(levels []int) int {
	asc := levels[0] < levels[1]
	for i := range len(levels) - 1 {
		first, second := levels[i], levels[i+1]
		if asc && first > second || !asc && first < second {
			return i
		}
		diff := shared.Abs(first - second)
		if diff < 1 || diff > 3 {
			return i
		}
	}
	return -1
}

func Multiply(text string, logic bool) int {
	total := 0
	skipping := false
	for {
		pair := mulPattern.FindStringIndex(text)
		if pair == nil {
			break
		}
		piece := text[pair[0]:pair[1]]
		text = text[pair[1]:]

		if strings.HasPrefix(piece, "don't") {
			skipping = true
			continue
		}
		if strings.HasPrefix(piece, "do") {
			skipping = false
			continue
		}
		if logic && skipping {
			continue
		}
		if strings.HasPrefix(piece, "mul") {
			a, b := splitMul(piece)
			total += a * b
		}
	}
	return total
}

//revive:disable-next-line:confusing-results
func splitMul(s string) (int, int) {
	separated := s[4 : len(s)-1]
	broken := strings.Split(separated, ",")
	return toInt(broken[0]), toInt(broken[1])
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Failed to convert to int: " + s)
	}
	return i
}

func CountInTable(table []string, word string) int {
	if len(word) == 0 {
		return 0
	}
	count := 0
	for r := 0; r < len(table); r++ {
		for c := 0; c < len(table[r]); c++ {
			if table[r][c] == word[0] {
				count += countWordsAt(table, word, r, c)
			}
		}
	}
	return count
}

func CountWordCrosses(table []string, word string) int {
	locations := findWordLocations(table, word)
	counts := map[location]int{}
	total := 0
	for _, each := range locations {
		if count, ok := counts[each]; ok {
			if count > 1 {
				panic("wtf")
			}
			counts[each]++
			total++
			shared.Logger.Info("Found word cross.", "location", each)
		} else {
			shared.Logger.Debug("Found half of word cross.", "location", each)
			counts[each] = 1
		}
	}
	return total
}

type direction struct {
	x int
	y int
}

func (v direction) turnRight() direction {
	if v == dirNorth {
		return dirEast
	}
	if v == dirEast {
		return dirSouth
	}
	if v == dirSouth {
		return dirWest
	}
	if v == dirWest {
		return dirNorth
	}
	panic("no direction")
}

func countWordsAt(table []string, word string, row, col int) int {
	directions := []direction{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}}
	count := 0
	for _, each := range directions {
		if readTableAt(table, row, col, len(word), each) == word {
			count++
		}
	}
	return count
}

type vector struct {
	loc location
	dir direction
}

type location struct {
	x int
	y int
}

func (v location) toString() string {
	return fmt.Sprintf("%dx%d", v.x, v.y)
}

func findWordLocations(table []string, word string) []location {
	if len(word)%2 != 1 {
		panic("only works with odd length words: 3, 5, 7, ...")
	}
	var locations []location
	mid := len(word) / 2
	for r := 0; r < len(table); r++ {
		for c := 0; c < len(table[r]); c++ {
			for _, each := range directions {
				if each.x == 0 || each.y == 0 {
					continue
				}
				if readTableAt(table, r, c, len(word), each) == word {
					x := r + each.x*mid
					y := c + each.y*mid
					locations = append(locations, location{x, y})
				}
			}
		}
	}
	return locations
}

func readTableAt(table []string, row, col, count int, dir direction) string {
	result := ""
	for range count {
		if row < 0 || row >= len(table) {
			break
		}
		line := table[row]
		if col < 0 || col >= len(line) {
			break
		}
		result += line[col : col+1]
		row += dir.x
		col += dir.y
	}
	return result
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
	blocks  *shared.Set[location]
	visited *shared.Set[vector]
	stepped *shared.Set[location]
	width   int
	height  int
}

func newBoard(init location, blocks []location, width, height int) *board {
	curr := vector{loc: init, dir: dirNorth}
	return &board{
		curr:    curr,
		blocks:  shared.NewSet(blocks),
		visited: shared.NewSet([]vector{curr}),
		stepped: shared.NewSet([]location{curr.loc}),
		width:   width,
		height:  height,
	}
}

func (v *board) deriveNextLocation() location {
	return location{
		x: v.curr.loc.x + v.curr.dir.y,
		y: v.curr.loc.y + v.curr.dir.x,
	}
}

func (v *board) deriveNextVector() vector {
	return vector{
		dir: v.curr.dir,
		loc: location{
			x: v.curr.loc.x + v.curr.dir.y,
			y: v.curr.loc.y + v.curr.dir.x,
		},
	}
}

func (v *board) isInside(l location) bool {
	if l.x < 0 || l.y < 0 {
		return false
	}
	if l.x >= v.height {
		return false
	}
	return l.y < v.width
}

func (v *board) turn() {
	v.visited.Add(v.curr)
	v.stepped.Add(v.curr.loc)
	v.curr.dir = v.curr.dir.turnRight()
}

func (v *board) isBlock(l location) bool {
	return v.blocks.Has(l)
}

func (v *board) move(l location) {
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
	v.blocks.Iter(func(item location) bool {
		setCharacter(lines, item, "#")
		return true
	})
	locToDirs := map[location][]direction{}
	v.visited.Iter(func(v vector) bool {
		l := v.loc
		if dirs, ok := locToDirs[l]; ok {
			locToDirs[l] = append(dirs, v.dir)
		} else {
			locToDirs[l] = []direction{v.dir}
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
			fmt.Sprintf("%s -> %s", brd.curr.loc.toString(), next.toString()),
		)
		brd.move(next)
		counter--
		if counter < 0 {
			panic("This loop is clearly eternal.")
		}
	}
	return brd.deriveVisitedCount()
}

func CountBlocksForIndefiniteLoops(lines []string) *shared.Set[location] {
	if len(lines) == 0 {
		return shared.NewSet([]location{})
	}
	shared.Logger.Info("Derive infinite loop locations.")
	brd := newBoard(
		findCharacter(lines, '^'),
		findLocations(lines, '#'),
		len(lines[0]),
		len(lines))
	indefLocations := shared.NewSet([]location{})
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

func findLocations(lines []string, c byte) []location {
	locs := []location{}
	for r := 0; r < len(lines); r++ {
		for l := 0; l < len(lines[r]); l++ {
			if lines[r][l] == c {
				locs = append(locs, location{r, l})
			}
		}
	}
	return locs
}

func findCharacter(lines []string, char byte) location {
	for r := 0; r < len(lines); r++ {
		for c := 0; c < len(lines[r]); c++ {
			if lines[r][c] == char {
				return location{r, c}
			}
		}
	}
	return location{-1, -1}
}

func setCharacter(lines []string, loc location, char string) {
	if char == "" || len(char) > 1 {
		panic(fmt.Sprintf("Invalid character length(%d): \"%s\".", len(char), char))
	}
	line := lines[loc.x]
	line = line[0:loc.y] + char + line[loc.y+1:]
	lines[loc.x] = line
}

func deriveDirCharacter(dirs []direction) string {
	hor, ver := false, false
	shared.NewSet(dirs).Iter(func(d direction) bool {
		if d.x == 0 && d.y != 0 {
			ver = true
		} else if d.x != 0 && d.y == 0 {
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
