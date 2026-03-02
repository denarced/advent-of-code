package aoc2303

import (
	"strconv"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

type hit struct {
	num      int
	from, to shared.Loc
}

func SumPartNumbers(lines []string) int {
	var sum int
	feedPartNumbers(lines, func(aHit hit) {
		shared.Logger.Debug("Check if part number found.", "hit", aHit)
		if isSurroundedBySymbol(lines, aHit) {
			shared.Logger.Info("Part number found.", "number", aHit.num, "loc", aHit.from)
			sum += aHit.num
		}
	})
	return sum
}

func SumGearRatios(lines []string) int {
	var sum int
	feedGears(lines, func(loc shared.Loc) {
		numbers := deriveAdjacentNumbers(lines, loc)
		if len(numbers) != 2 {
			return
		}
		shared.Logger.Info("Pair found.", "numbers", numbers)
		sum += numbers[0] * numbers[1]
	})
	return sum
}

func feedGears(lines []string, cb func(shared.Loc)) {
	for y, line := range lines {
		for x, c := range line {
			if c == '*' {
				cb(shared.Loc{X: x, Y: y})
			}
		}
	}
}

func deriveAdjacentNumbers(lines []string, loc shared.Loc) []int {
	width, height := len(lines[0]), len(lines)
	locations := gent.NewSet[shared.Loc]()
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			deltaLoc := loc.Delta(shared.Loc{X: x, Y: y})
			if isLocValid(deltaLoc, width, height) {
				locations.Add(deltaLoc)
			}
		}
	}
	pickNext := func() shared.Loc {
		var picked shared.Loc
		locations.ForEach(func(each shared.Loc, stop func()) {
			picked = each
			stop()
		})
		return picked
	}
	findLast := func(start, delta shared.Loc) shared.Loc {
		curr := start
		for {
			next := curr.Delta(delta)
			if !isLocValid(next, width, height) {
				break
			}
			if !isDigit(getRune(lines, next)) {
				break
			}
			curr = next
		}
		return curr
	}
	var numbers []int
	for locations.Count() > 0 {
		picked := pickNext()
		if !isDigit(getRune(lines, picked)) {
			locations.Remove(picked)
			continue
		}
		left := findLast(picked, shared.Loc{X: -1})
		right := findLast(picked, shared.Loc{X: 1})
		runes := make([]rune, 0, right.X-left.X+1)
		for x := left.X; x <= right.X; x++ {
			curr := shared.Loc{X: x, Y: picked.Y}
			locations.Remove(curr)
			runes = append(runes, getRune(lines, curr))
		}
		numbers = append(numbers, toInt(runes))
	}
	return numbers
}

func getRune(lines []string, loc shared.Loc) rune {
	return []rune(lines[loc.Y])[loc.X]
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func feedPartNumbers(lines []string, cb func(aHit hit)) {
	report := func(num int, from, to shared.Loc) {
		cb(hit{num: num, from: from, to: to})
	}
	for y, line := range lines {
		var st state = &inactiveState{}
		for x, c := range line {
			st = st.handle(c, shared.Loc{X: x, Y: y}, report)
		}
		st.finish(report)
	}
}

func isLocValid(loc shared.Loc, width, height int) bool {
	if loc.X < 0 || loc.Y < 0 {
		return false
	}
	if loc.X >= width || loc.Y >= height {
		return false
	}
	return true
}

func isSurroundedBySymbol(lines []string, aHit hit) bool {
	width, height := len(lines[0]), len(lines)
	check := func(base, delta shared.Loc) bool {
		loc := base.Delta(delta)
		if !isLocValid(loc, width, height) {
			return false
		}
		return isSymbol(getRune(lines, loc))
	}

	// Check left side of surrounding cells.
	fromDeltas := []shared.Loc{
		{X: -1},
		{X: -1, Y: -1},
		{X: -1, Y: 1},
	}
	for _, each := range fromDeltas {
		if check(aHit.from, each) {
			return true
		}
	}

	// Check right side of surrounding cells.
	toDeltas := []shared.Loc{
		{X: 1, Y: -1},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
	}
	for _, each := range toDeltas {
		if check(aHit.to, each) {
			return true
		}
	}

	// Check cells above and below between "from" and "to".
	current := aHit.from
	for isLocValid(current, width, height) {
		for _, each := range []shared.Loc{{Y: -1}, {Y: 1}} {
			if check(current, each) {
				return true
			}
		}
		if current == aHit.to {
			break
		}
		current = current.Delta(shared.Loc{X: 1})
	}
	return false
}

func isSymbol(c rune) bool {
	if c == '.' {
		return false
	}
	return !isDigit(c)
}

type stateCallback = func(int, shared.Loc, shared.Loc)
type state interface {
	handle(c rune, loc shared.Loc, cb stateCallback) state
	finish(cb stateCallback)
}

type inactiveState struct{}

func (v *inactiveState) handle(c rune, loc shared.Loc, _ func(int, shared.Loc, shared.Loc)) state {
	if '0' <= c && c <= '9' {
		shared.Logger.Debug("Start active state.", "c", string(c), "loc", loc)
		active := &activeState{
			from:  loc,
			to:    loc,
			chars: []rune{c},
		}
		return active
	}
	return v
}

func (*inactiveState) finish(_ stateCallback) {}

type activeState struct {
	from, to shared.Loc
	chars    []rune
}

func (v *activeState) handle(c rune, loc shared.Loc, cb func(int, shared.Loc, shared.Loc)) state {
	if '0' <= c && c <= '9' {
		v.chars = append(v.chars, c)
		v.to = loc
		return v
	}
	shared.Logger.Debug("Ache from active state.")
	cb(toInt(v.chars), v.from, v.to)
	return &inactiveState{}
}

func toInt(chars []rune) int {
	s := string(chars)
	num, err := strconv.Atoi(s)
	if err != nil {
		shared.Logger.Error("Failed to convert to int.", "chars", chars)
		panic(err)
	}
	return num
}

func (v *activeState) finish(cb stateCallback) {
	shared.Logger.Debug("Ache from finish.")
	cb(toInt(v.chars), v.from, v.to)
}
