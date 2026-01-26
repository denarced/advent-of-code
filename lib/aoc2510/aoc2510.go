package aoc2510

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const globalSafety = 10000000

var globalButtonIndexes = map[int][]int{}

func DeriveFewestClicks(lines []string, indicator bool, version int) int {
	shared.Logger.Info("Derive fewest clicks.", "machine count", len(lines), "indicator", indicator)
	var clicks int
	defer func() {
		shared.Logger.Info("Fewest clicks derived.", "count", clicks)
	}()
	deriveFunc := deriveFewestStateClicks
	if !indicator {
		deriveFunc = func(mach Machine) int {
			return deriveFewestJoltageClicks(mach, version)
		}
	}
	var wg sync.WaitGroup
	var done int
	machines := ParseMachines(lines)
	for _, each := range machines {
		shared.Logger.Info("Derive fewest clicks for a machine.", "machine", each)
		wg.Add(1)
		go func(mach Machine) {
			defer wg.Done()
			alpha := time.Now()
			n := deriveFunc(mach)
			done++
			shared.Logger.Info(
				"Clicks derived.",
				"done",
				done,
				"clicks",
				n,
				"duration",
				time.Since(alpha),
				"machine",
				mach,
			)
			clicks += n
		}(each)
	}
	wg.Wait()
	return clicks
}

type Machine struct {
	TargetState []bool
	Buttons     [][]int
	Joltages    []int
	id          int
}

func ParseMachines(lines []string) []Machine {
	machines := make([]Machine, len(lines))
	for i, line := range lines {
		mach := parseMachine(line)
		mach.id = i
		machines[i] = mach
	}
	return machines
}

func parseMachine(s string) Machine {
	pieces := strings.Fields(s)
	stateString := pieces[0]
	joltageString := pieces[len(pieces)-1]
	buttons := make([][]int, len(pieces)-2)
	for i, buttonString := range pieces[1 : len(pieces)-1] {
		buttons[i] = parseInts(buttonString)
	}
	targetState := parseState(stateString)
	return Machine{
		TargetState: targetState,
		Buttons:     buttons,
		Joltages:    parseInts(joltageString),
	}
}

func parseState(s string) []bool {
	state := make([]bool, len(s)-2)
	for i, each := range s[1 : len(s)-1] {
		state[i] = each == '#'
	}
	return state
}

func parseInts(s string) []int {
	var ints []int
	var current *int
	for _, each := range s {
		value := int(each - '0')
		digit := unicode.IsDigit(each)
		if current != nil {
			if digit {
				*current *= 10
				*current += value
			} else {
				ints = append(ints, *current)
				current = nil
			}
		} else {
			if digit {
				current = &value
			}
		}
	}
	if current != nil {
		ints = append(ints, *current)
	}
	return ints
}

func deriveFewestStateClicks(mach Machine) int {
	states := [][]bool{make([]bool, len(mach.TargetState))}
	seen := gent.NewSet(toNumericState(states[0]))
	for count := range 10 {
		var nextStates [][]bool
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Round.", "#", count, "state count", len(states))
		}
		for _, st := range states {
			for _, button := range mach.Buttons {
				state := click(st, button)
				if slices.Equal(state, mach.TargetState) {
					return count + 1
				}
				if seen.Add(toNumericState(state)) {
					nextStates = append(nextStates, state)
				}
			}
		}
		if len(nextStates) == 0 {
			shared.Logger.Warn("Zero states left, fundamentally broken.")
			return -1
		}
		states = nextStates
	}
	return -1
}

func click(state []bool, button []int) []bool {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Click.", "state", state, "button", button)
	}
	var next []bool
	next = append(next, state...)
	for _, each := range button {
		next[each] = !next[each]
	}
	return next
}

func increment(joltages []int, button []int) []int {
	var next []int
	next = append(next, joltages...)
	for _, each := range button {
		next[each]++
	}
	return next
}

func decrement(joltages []int, button []int) []int {
	var next []int
	next = append(next, joltages...)
	for _, each := range button {
		next[each]--
	}
	return next
}

func decInPlace(joltages []int, button []int) {
	for _, each := range button {
		joltages[each]--
	}
}

func decInPlaceN(joltages []int, button []int, n int) {
	for _, each := range button {
		joltages[each] -= n
	}
}

func incInPlace(joltages []int, button []int) {
	for _, each := range button {
		joltages[each]++
	}
}

func toNumericState(state []bool) int64 {
	var v int64 = 1
	var n int64
	for i := len(state) - 1; i >= 0; i-- {
		if state[i] {
			n += v
		}
		v *= 2
	}
	return n
}

func deriveFewestJoltageClicks(mach Machine, version int) int {
	switch version {
	case 1:
		// Starting with zero joltages go through all possible permutations while holding all of
		// them in memory, clicking all buttons once, and then adding another permutation. The
		// amount of permutations quickly explodes.
		return deriveFewestJoltageClicks1(mach)
	case 2:
		// Almost the same as #1 but sort buttons (biggest to smallest), start by clicking the
		// biggest button as many times as possible, and then going through all possible
		// permutations. At least initiallly there's less permutations than #1 because you start
		// from the end and most attempts exceed joltages because the initial clicks almost maxed
		// out joltages.
		return deriveFewestJoltageClicks2(mach)
	case 3:
		// Depth first attempt to overcome the problem with storing all permutation in memory.
		// Otherwise the basic approach is the same as in #2 but recursion is used to essentially
		// store all permutations in the stack because joltages are never copied. Process starts
		// with the target joltages and then each click reduces the values. Since buttons are sorted
		// in descending order (biggest first) the first case of zero joltage is considered the
		// result: minimum amount of clicks.
		return deriveFewestJoltageClicks3(mach)
	case 4:
		// Slight improvement on #3: only relevant buttons are clicked during any given recursion of
		// tryClick. That is, those buttons that reduce remaining >0 joltage values.
		return deriveFewestJoltageClicks4(mach)
	case 5:
		// Not sure why this was created. It goes through all possible button combinations in the
		// same way that #3 and #4 do but it doesn't stop when a combination of clicks is found, nor
		// is there any mechanism to stop descending when known click count is exceeded.
		return deriveFewestJoltageClicks5(mach)
	case 6:
		return deriveFewestJoltageClicks6(mach)
	case 7:
		return deriveFewestJoltageClicks7(mach)
	case 8:
		return deriveFewestJoltageClicks8(mach)
	case 9:
		return deriveFewestJoltageClicks9(mach)
	default:
		panic(fmt.Sprint("no such version:", version))
	}
}

func deriveFewestJoltageClicks1(mach Machine) int {
	count := 0
	zeroState := make([]int, len(mach.Joltages))
	joltages := [][]int{zeroState}
	seen := gent.NewSet(hash(joltages[0]))
	for ; count < 100; count++ {
		var nextStates [][]int
		shared.Logger.Debug("Round.", "#", count, "state count", len(joltages))
		for _, jolt := range joltages {
			for _, button := range mach.Buttons {
				jolted := increment(jolt, button)
				if slices.Equal(jolted, mach.Joltages) {
					return count + 1
				}
				if !more(jolted, mach.Joltages) && seen.Add(hash(jolted)) {
					nextStates = append(nextStates, jolted)
				}
			}
		}
		if len(nextStates) == 0 {
			shared.Logger.Warn("Zero states left, fundamentally broken.")
			return -1
		}
		joltages = nextStates
	}
	return -1
}

func hash(s []int) uint64 {
	var h uint64 = 17
	for _, v := range s {
		h = h*31 + uint64(v)
	}
	return h
}

func more(current, target []int) bool {
	for i := range len(current) {
		if current[i] > target[i] {
			return true
		}
	}
	return false
}

func deriveFewestJoltageClicks2(mach Machine) int {
	sortedButtons := sortButtons(mach.Buttons)
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Button sorted.", "buttons", sortedButtons)
	}
	zeroState := make([]int, len(mach.Joltages))
	initialJoltage, initialCount := deriveInitialJoltage(zeroState, mach.Joltages, sortedButtons[0])
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Initial joltage derived.", "joltage", initialJoltage)
	}
	if slices.Equal(initialJoltage, mach.Joltages) {
		return initialCount
	}
	joltages := [][]int{initialJoltage}
	for ; initialCount >= 0; initialCount-- {
		shared.Logger.Info("Round with initial count.", "initial", initialCount)
		seen := gent.NewSet(hash(joltages[0]))
		for count := initialCount; count < 100; count++ {
			var nextStates [][]int
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Round.", "#", count, "state count", len(joltages))
			}
			for _, jolt := range joltages {
				for _, button := range sortedButtons {
					jolted := increment(jolt, button)
					if slices.Equal(jolted, mach.Joltages) {
						return count + 1
					}
					if !more(jolted, mach.Joltages) && seen.Add(hash(jolted)) {
						nextStates = append(nextStates, jolted)
					}
				}
			}
			if len(nextStates) == 0 {
				break
			}
			joltages = nextStates
		}
		initialJoltage = decrement(initialJoltage, sortedButtons[0])
		joltages = [][]int{initialJoltage}
	}
	return -1
}

func sortButtons(buttons [][]int) [][]int {
	sorted := make([][]int, len(buttons))
	copy(sorted, buttons)
	slices.SortFunc(sorted, func(a, b []int) int {
		return len(b) - len(a)
	})
	return sorted
}

func deriveInitialJoltage(joltage, target, button []int) ([]int, int) {
	var clicks int
	for {
		jolted := increment(joltage, button)
		if more(jolted, target) {
			break
		}
		joltage = jolted
		clicks++
	}
	return joltage, clicks
}

func deriveFewestJoltageClicks3(mach Machine) int {
	sortedButtons := sortButtons(mach.Buttons)
	seen := gent.NewSet[uint64]()
	return tryClick(seen, sortedButtons, mach.Joltages, 0, false)
}

func deriveFewestJoltageClicks4(mach Machine) int {
	sortedButtons := sortButtons(mach.Buttons)
	seen := gent.NewSet[uint64]()
	return tryClick(seen, sortedButtons, mach.Joltages, 0, true)
}

func tryClick(
	seen *gent.Set[uint64],
	buttons [][]int,
	joltage []int,
	clicks int,
	optimize bool,
) int {
	if seen.Has(hash(joltage)) {
		return -1
	}
	if seen.Len() > 4_000_000 {
		count := 1_000_000
		seen.ForEach(func(each uint64, stop func()) {
			seen.Remove(each)
			count--
			if count <= 0 {
				stop()
			}
		})
	}
	for _, i := range optimizeButtons(joltage, buttons, optimize) {
		decInPlace(joltage, buttons[i])
		neg, zeroes := extractSpecs(joltage)
		var res int
		if !neg && !zeroes {
			res = tryClick(seen, buttons, joltage, clicks+1, optimize)
			seen.Add(hash(joltage))
		}
		incInPlace(joltage, buttons[i])
		if neg {
			continue
		}
		if zeroes {
			return clicks + 1
		}
		if res > 0 {
			return res
		}
	}
	return -1
}

func extractSpecs(ints []int) (negative, zeros bool) {
	zeros = true
	for _, each := range ints {
		if each < 0 {
			negative = true
			zeros = false
		} else if each > 0 {
			zeros = false
		}
	}
	return
}

func optimizeButtons(joltage []int, buttons [][]int, optimize bool) []int {
	if !optimize || !hasZero(joltage) {
		size := len(buttons)
		indexes, ok := globalButtonIndexes[size]
		if ok {
			return indexes
		}
		indexes = make([]int, 0, size)
		for i := range buttons {
			indexes = append(indexes, i)
		}
		globalButtonIndexes[len(buttons)] = indexes
		return indexes
	}
	var indexes []int
	var needed []int
	for i, value := range joltage {
		if value > 0 {
			needed = append(needed, i)
		}
	}
	for i, each := range buttons {
		if overlap(each, needed) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func overlap(a, b []int) bool {
	m := map[int]int{}
	for _, each := range a {
		m[each] = 1
	}
	for _, each := range b {
		if m[each] == 1 {
			return true
		}
	}
	return false
}

func hasZero(ints []int) bool {
	for i := range ints {
		if ints[i] <= 0 {
			return true
		}
	}
	return false
}

func isZeroes(ints []int) bool {
	for i := range ints {
		if ints[i] != 0 {
			return false
		}
	}
	return true
}

func deriveFewestJoltageClicks5(mach Machine) int {
	// Derive possible button combinations. Sort buttons into descending order, from the most
	// powerful (affects most joltages) to the least. Then pick smallest joltage relevant for the
	// button, click the button that many times to guarantee that no joltage flips over to the
	// negative side, and move on to the next button. That which is most powerful when measured by
	// the number of remaining joltages it has an effect to. Then again, pick the minimum remaining
	// joltage, click that many times, etc.
	_, high := deriveClickLimits(mach)
	sortedButtons := sortButtons(mach.Buttons)
	mach.Buttons = sortedButtons
	clicks := make([]int, len(mach.Buttons))
	cand := make([]int, len(mach.Joltages))
	minClicks := math.MaxInt
	// TODO Misleading variable name.
	totalClicks := countSum(mach.Joltages)
	iter(clicks, high, len(mach.Buttons), 0, totalClicks, func() bool {
		copy(cand, mach.Joltages)
		for i, each := range clicks {
			if each <= 0 {
				continue
			}
			button := mach.Buttons[i]
			for range each {
				decInPlace(cand, button)
			}
		}
		if shared.IsDebugEnabled() {
			shared.Logger.Debug(
				"Prepared candidate.",
				"cand",
				cand,
				"clicks",
				clicks,
			)
		}
		if isZeroes(cand) {
			minClicks = min(minClicks, countSum(clicks))
		}
		return false
	})

	return minClicks
}

func countSum(ints []int) int {
	var sum int
	for _, each := range ints {
		sum += each
	}
	return sum
}

func deriveClickLimits(mach Machine) (low int, high int) {
	minJoltage, maxJoltage := math.MaxInt, math.MinInt
	for _, each := range mach.Joltages {
		minJoltage = min(each, minJoltage)
		maxJoltage = max(each, maxJoltage)
	}

	minPower, maxPower := math.MaxInt, math.MinInt
	for _, each := range mach.Buttons {
		length := len(each)
		minPower = min(minPower, length)
		maxPower = max(maxPower, length)
	}
	low = int(math.Ceil(float64(minJoltage) / float64(maxPower)))
	high = int(math.Ceil(float64(maxJoltage) / float64(minPower)))
	return
}

func iter(container []int, high, length, index, maxClicks int, doneFunc func() bool) bool {
	for i := high; i >= 0; i-- {
		remaining := maxClicks - i
		if remaining < 0 {
			continue
		}
		container[index] = i
		if index+1 >= length {
			if doneFunc() {
				return true
			}
		} else {
			if iter(container, high, length, index+1, remaining, doneFunc) {
				return true
			}
		}
	}
	return false
}

type counter struct {
	m map[int]int
}

func newCounter() *counter {
	return &counter{m: map[int]int{}}
}

func (v *counter) add(ints []int) {
	for _, each := range ints {
		value := v.m[each]
		v.m[each] = value + 1
	}
}

func (v *counter) remove(ints []int) {
	for _, each := range ints {
		value := v.m[each] - 1
		if value <= 0 {
			if value < 0 {
				panic("value in counter should never become negative")
			}
			delete(v.m, each)
		} else {
			v.m[each] = value
		}
	}
}

func (v *counter) size() int {
	size := len(v.m)
	return size
}

func initDeriveButtonCombinations(mach Machine, cb func(combination []int)) {
	deriveButtonCombinations(
		mach,
		make([]bool, len(mach.Buttons)),
		0,
		newCounter(),
		gent.NewSet[uint64](),
		cb)
}

func deriveButtonCombinations(
	mach Machine,
	buttons []bool,
	buttonIndex int,
	aCounter *counter,
	seen *gent.Set[uint64],
	cb func(combination []int),
) {
	currentButton := mach.Buttons[buttonIndex]
	canClick := canDec(mach.Joltages, currentButton)
	nextButtonIndex := buttonIndex + 1

	for _, adding := range []bool{true, false} {
		buttons[buttonIndex] = adding
		if canClick {
			if adding {
				decInPlace(mach.Joltages, currentButton)
			} else {
				incInPlace(mach.Joltages, currentButton)
			}
		}
		if adding {
			verifyNotNegative(mach.Joltages)
			aCounter.add(currentButton)
		} else {
			aCounter.remove(currentButton)
		}
		if aCounter.size() == len(mach.Joltages) {
			trueIndexes := extractTrueIndexes(buttons)
			if seen.Add(hash(trueIndexes)) {
				if shared.IsDebugEnabled() {
					shared.Logger.Debug("Call back.", "indexes", trueIndexes)
				}
				cb(trueIndexes)
			}
		}
		if !isZeroes(mach.Joltages) && nextButtonIndex < len(buttons) {
			deriveButtonCombinations(
				mach,
				buttons,
				nextButtonIndex,
				aCounter,
				seen,
				cb,
			)
		}
	}
}

func extractTrueIndexes(s []bool) []int {
	var trues []int
	for i, each := range s {
		if each {
			trues = append(trues, i)
		}
	}
	return trues
}

func minRelevantJoltage(joltages, button []int) int {
	minimus := math.MaxInt
	for i, each := range joltages {
		if slices.Contains(button, i) {
			minimus = min(minimus, each)
			if minimus == 0 {
				return minimus
			}
		}
	}
	return minimus
}

func canDec(joltages, button []int) bool {
	maximus := -1
	minimus := math.MaxInt
	for i, each := range joltages {
		if slices.Contains(button, i) {
			minimus = min(minimus, each)
			maximus = max(maximus, each)
		}
	}
	return maximus > 0 && minimus > 0
}

func verifyNotNegative(ints []int) {
	for _, each := range ints {
		if each < 0 {
			panic("should be impossible to have negative values")
		}
	}
}

func deriveFewestJoltageClicks6(mach Machine) int {
	mach.Buttons = sortButtons(mach.Buttons)
	combos := [][]int{}
	initDeriveButtonCombinations(mach, func(combo []int) {
		copied := make([]int, len(combo))
		copy(copied, combo)
		combos = append(combos, copied)
	})
	minClickCount := math.MaxInt
	minJolt := minJoltage(mach.Joltages)
	for _, each := range combos {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Try combo.", "combo", each)
		}
		prev := minClickCount
		dive(&mach, 0, &minClickCount, each, 0)
		if prev != minClickCount {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug(
					"New min.",
					"combo", each,
					"min", minClickCount)
			}
		}
		if minJolt == minClickCount {
			break
		}
	}
	if shared.IsDebugEnabled() {
		shared.Logger.Debug(
			"Fewest derived.",
			"button combo count", len(combos),
			"joltages", mach.Joltages,
			"min count", minClickCount)
	}
	return minClickCount
}

func dive(
	mach *Machine,
	clickCount int,
	minClickCount *int,
	buttonIndexes []int,
	currentButtonIndex int,
) {
	if currentButtonIndex >= len(buttonIndexes) {
		return
	}
	if isJoltageIndexMissing(mach, buttonIndexes[currentButtonIndex:]) {
		return
	}
	currentButton := mach.Buttons[buttonIndexes[currentButtonIndex]]
	count := minRelevantJoltage(mach.Joltages, currentButton)
	last := currentButtonIndex == len(buttonIndexes)-1
	if shared.IsDebugEnabled() {
		shared.Logger.Debug(
			"Start diving.",
			"click count", clickCount,
			"current min", *minClickCount,
			"button index", currentButtonIndex,
			"dec count", count,
			"joltage", mach.Joltages)
	}
	if clickCount+count >= *minClickCount {
		if last {
			return
		}
		count = *minClickCount - clickCount - 1
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Need to decrease dec count.", "new count", count)
		}
	}
	if count <= 0 {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Zero count.", "index", currentButtonIndex)
		}
		dive(mach, clickCount, minClickCount, buttonIndexes, currentButtonIndex+1)
		return
	}
	for range count {
		decInPlace(mach.Joltages, currentButton)
		clickCount++
	}
	for range count {
		if isZeroes(mach.Joltages) {
			if clickCount < *minClickCount {
				if shared.IsDebugEnabled() {
					shared.Logger.Debug(
						"Match found.",
						"clickCount", clickCount,
						"min", *minClickCount)
				}
				*minClickCount = clickCount
			}
		} else {
			dive(mach, clickCount, minClickCount, buttonIndexes, currentButtonIndex+1)
		}
		incInPlace(mach.Joltages, currentButton)
		clickCount--
	}
}

func minJoltage(joltages []int) int {
	minimus := math.MaxInt
	for _, each := range joltages {
		minimus = min(minimus, each)
	}
	return minimus
}

func isJoltageIndexMissing(mach *Machine, buttonIndexes []int) bool {
	for _, buttonIndex := range buttonIndexes {
		button := mach.Buttons[buttonIndex]
		for _, each := range button {
			if mach.Joltages[each] > 0 {
				return false
			}
		}
	}
	return true
}

func deriveFewestJoltageClicks7(mach Machine) int {
	shared.Logger.Info("Derive fewest clicks for joltages.", "machine", mach)
	// There's multiple worlds, overlord, and angels. Each of the angels represent a joltage and
	// they are responsible for finding button combinations with which they can fulfill they're
	// joltage. Each world is a combination of harmonious clicks where all angels are satisfied with
	// all of them.
	maxMultiplier := deriveMostPowerfulButton(mach.Buttons)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lord, angels := setupOverlord(ctx, mach, maxMultiplier)
	seed, others, angelIndex := pickSeed(angels)
	safety := globalSafety
	for safety > 0 && seed != nil {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Start with seed.", "seed", seed, "angel", angelIndex)
		}
		safety--
		candidate := make([]int, len(mach.Joltages))
		copy(candidate, mach.Joltages)
		reduceJoltage(mach.Buttons, candidate, seed)
		less := isLess(candidate)
		zeroes := isZeroes(candidate)
		if less {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Lost - overflow.", "joltage", candidate)
			}
		}
		if !less && !zeroes {
			interested := filterInterestedAngels(others, candidate)
			if len(interested) > 0 {
				others[0].find(candidate, interested, 0, [][]step{seed})
			}
		}
		if zeroes {
			if countSteps(seed) < lord.minCount {
				lord.addWorld([][]step{seed})
			} else {
				if shared.IsDebugEnabled() {
					shared.Logger.Debug("Lost - better match.", "joltage", candidate)
				}
			}
		}

		seed, others, angelIndex = pickSeed(angels)
	}
	if safety <= 0 {
		panic("safety<=0")
	}
	return lord.minCount
}

func reduceJoltage(buttons [][]int, jolts []int, steps []step) {
	for _, each := range steps {
		button := buttons[each.buttonID]
		for _, i := range button {
			jolts[i] -= each.count
		}
	}
}

func pickSeed(angels []*joltaAngel) ([]step, []*joltaAngel, int) {
	for i, each := range angels {
		if each.done {
			continue
		}
		if seed := each.getNextSeed(); seed != nil {
			return seed, exclude(angels, i), i
		}
	}
	return nil, nil, -1
}

func exclude(angels []*joltaAngel, index int) []*joltaAngel {
	result := make([]*joltaAngel, 0, len(angels)-1)
	if index > 0 {
		result = append(result, angels[:index]...)
	}
	if index < len(angels) {
		result = append(result, angels[index+1:]...)
	}
	return result
}

func setupOverlord(
	ctx context.Context,
	mach Machine,
	maxMultiplier int,
) (*overlord, []*joltaAngel) {
	lord := newOverlord(mach.Buttons, mach.Joltages)
	var angels []*joltaAngel
	adjustJoltage := func(joltages []int, steps []step, mod int) {
		for _, each := range steps {
			button := mach.Buttons[each.buttonID]
			for _, i := range button {
				joltages[i] += (mod * each.count)
			}
		}
	}
	reduceJoltage := func(joltages []int, steps []step) {
		adjustJoltage(joltages, steps, -1)
	}
	increaseJoltage := func(joltages []int, steps []step) {
		adjustJoltage(joltages, steps, 1)
	}
	angelMaximums := deriveAngelMaximums(mach)
	for i := range mach.Joltages {
		angels = append(
			angels,
			newJoltaAngel(
				ctx,
				i,
				angelMaximums[i],
				lord.getButtons(i),
				lord.getButtonJoltageIDs(i),
				reduceJoltage,
				increaseJoltage,
				lord,
				maxMultiplier,
			),
		)
	}
	return lord, angels
}

type overlord struct {
	buttons  [][]int
	joltages []int
	worlds   [][][]step
	minCount int
}

func (v *overlord) getButtons(joltaID int) []int {
	var result []int
	for i, joltaIDs := range v.buttons {
		if slices.Contains(joltaIDs, joltaID) {
			result = append(result, i)
		}
	}
	return result
}

func (v *overlord) getButtonJoltageIDs(joltaID int) [][]int {
	var result [][]int
	for _, joltaIDs := range v.buttons {
		if slices.Contains(joltaIDs, joltaID) {
			result = append(result, joltaIDs)
		}
	}
	return result
}

func newOverlord(buttons [][]int, joltages []int) *overlord {
	copied := make([]int, len(joltages))
	copy(copied, joltages)
	return &overlord{
		buttons:  buttons,
		joltages: copied,
		minCount: math.MaxInt,
	}
}

func (v *overlord) addWorld(stepLists [][]step) {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Add world.", "steps", stepLists)
	}
	joltages := make([]int, len(v.joltages))
	for _, nested := range stepLists {
		for _, each := range nested {
			button := v.buttons[each.buttonID]
			for _, index := range button {
				joltages[index] += each.count
			}
		}
	}
	if !slices.Equal(v.joltages, joltages) {
		shared.Logger.Error("Different joltages.", "original", v.joltages, "derived", joltages)
		panic("fundamentally broken: different joltages")
	}
	count := countStepLists(stepLists)
	if v.minCount <= count {
		panic("it shouldn't be possible to add a world with >= current minimimum")
	}
	v.worlds = append(v.worlds, stepLists)
	v.minCount = count
	shared.Logger.Info("World added.", "count", count, "steps", stepLists)
}

type joltaAngel struct {
	ctx              context.Context
	joltageID        int
	value            int
	buttons          []int
	buttonJoltageIDs [][]int
	ch               <-chan []step
	done             bool
	reduceJoltage    func(joltages []int, steps []step)
	increaseJoltage  func(joltages []int, steps []step)
	lord             *overlord
	maxMultiplier    int
}

//revive:disable-next-line:argument-limit
func newJoltaAngel(
	ctx context.Context,
	joltageID, value int,
	buttons []int,
	buttonJoltageIDs [][]int,
	reduceJoltage func([]int, []step),
	increaseJoltage func([]int, []step),
	lord *overlord,
	maxMultiplier int,
) *joltaAngel {
	return &joltaAngel{
		ctx:              ctx,
		joltageID:        joltageID,
		value:            value,
		buttons:          buttons,
		buttonJoltageIDs: buttonJoltageIDs,
		reduceJoltage:    reduceJoltage,
		increaseJoltage:  increaseJoltage,
		lord:             lord,
		maxMultiplier:    maxMultiplier,
	}
}

type step struct {
	buttonID int
	count    int
}

func (v *joltaAngel) getNextSeed() []step {
	if v.done {
		panic("shouldn't call getNextSeed when it's done")
	}
	if v.ch == nil {
		v.ch = createStepFeeder(v.ctx, v.value, v.buttons, nil, nil)
	}
	next, ok := <-v.ch
	if !ok {
		v.ch = nil
		v.done = true
	}
	return next
}

func (v *joltaAngel) isIntested(joltagees []int) bool {
	for i, each := range joltagees {
		if i == v.joltageID && each > 0 {
			return true
		}
	}
	return false
}

func createStepFeeder(
	ctx context.Context,
	value int,
	buttons []int,
	joltages []int,
	buttonJoltageIDs [][]int,
) <-chan []step {
	if isLess(joltages) {
		panic("picnic")
	}
	ch := make(chan []step)
	state := make([]int, len(buttons))
	go feedSteps(ctx, value, 0, buttons, state, 0, ch, joltages, buttonJoltageIDs)
	return ch
}

//revive:disable-next-line:argument-limit
func feedSteps(
	ctx context.Context,
	total, accrued int,
	buttons, state []int,
	index int,
	ch chan []step,
	joltages []int,
	buttonJoltageIDs [][]int,
) {
	defer func() {
		if index == 0 {
			close(ch)
		}
	}()
	high := total - accrued
	last := index == len(state)-1
	if high > 0 && joltages != nil {
		high = min(high, deriveMaxClickCount(joltages, buttonJoltageIDs[index]))
	}
	if shared.IsDebugEnabled() {
		shared.Logger.Debug(
			"Feed.",
			"high", high,
			"total", total,
			"accrued", accrued,
			"index", index,
			"joltages", joltages,
		)
	}
	for i := high; i >= 0; i-- {
		state[index] = i
		localAccrued := accrued + i
		if last {
			if localAccrued == total {
				select {
				case ch <- createSteps(buttons, state):
				case <-ctx.Done():
					return
				}
			}
		} else {
			feedSteps(
				ctx,
				total,
				localAccrued,
				buttons,
				state,
				index+1,
				ch,
				joltages,
				buttonJoltageIDs)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}
}

func createSteps(buttons, state []int) []step {
	steps := make([]step, 0, len(state))
	for i, each := range state {
		if each > 0 {
			steps = append(steps, step{
				buttonID: buttons[i],
				count:    each,
			})
		}
	}
	sortSteps(steps)
	return steps
}

//revive:disable-next-line:function-length
//revive:disable-next-line:cyclomatic
//revive:disable-next-line:cognitive-complexity
func (v *joltaAngel) find(joltages []int, angels []*joltaAngel, index int, steps [][]step) {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug(
			"Find.",
			"joltages", joltages,
			"angel ID", v.joltageID,
			"steps", steps)
	}
	if !v.isIntested(joltages) {
		if index+1 < len(angels) {
			angels[index+1].find(joltages, angels, index+1, steps)
		}
		return
	}
	safety := globalSafety
	baseCount := countStepLists(steps)
	available := v.lord.minCount - baseCount
	joltSum := countSum(joltages)
	if joltSum > 0 {
		minRequired := joltSum / v.maxMultiplier
		if minRequired > available {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug(
					"Lost - not even theoretically possible to beat previous minimum.",
					"min required", minRequired,
					"available", available)
			}
			return
		}
	}
	ch := createStepFeeder(
		context.TODO(),
		joltages[v.joltageID],
		v.buttons,
		joltages,
		v.buttonJoltageIDs,
	)
	for safety > 0 {
		safety--
		newSteps, ok := <-ch
		if !ok {
			return
		}
		stepCount := countSteps(newSteps)
		if stepCount > available {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug(
					"Lost - steps > available.",
					"step count",
					stepCount,
					"available",
					available,
				)
			}
			continue
		}
		currentTotalSteps := baseCount + stepCount
		if currentTotalSteps >= v.lord.minCount {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug(
					"Lost - total steps > lord.minCount.",
					"current", currentTotalSteps,
					"lord.minCount", v.lord.minCount)
			}
			continue
		}
		v.reduceJoltage(joltages, newSteps)
		if isZeroes(joltages) {
			v.lord.addWorld(append(steps, newSteps))
			v.increaseJoltage(joltages, newSteps)
			continue
		}
		if isLess(joltages) {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug(
					"Lost - negative joltage.",
					"joltages", joltages,
					"steps", newSteps)
			}
			v.increaseJoltage(joltages, newSteps)
			continue
		}
		nextIndex := index + 1
		if nextIndex < len(angels) {
			angels[nextIndex].find(joltages, angels, nextIndex, append(steps, newSteps))
		}
	}
	if safety <= 0 {
		panic("safety count in find reached zero")
	}
}

func sortSteps(steps []step) {
	slices.SortFunc(steps, func(a, b step) int {
		button := a.buttonID - b.buttonID
		if button != 0 {
			return button
		}
		return a.count - b.count
	})
}

func countStepLists(steps [][]step) int {
	var count int
	for _, each := range steps {
		count += countSteps(each)
	}
	return count
}

func countSteps(steps []step) int {
	var count int
	for _, each := range steps {
		count += each.count
	}
	return count
}

func isLess(joltages []int) bool {
	for i := range joltages {
		if joltages[i] < 0 {
			return true
		}
	}
	return false
}

func filterInterestedAngels(angels []*joltaAngel, joltages []int) []*joltaAngel {
	var interested []*joltaAngel
	for _, each := range angels {
		if each.isIntested(joltages) {
			interested = append(interested, each)
		}
	}
	return interested
}

func deriveMostPowerfulButton(buttons [][]int) int {
	most := math.MinInt
	for _, each := range buttons {
		most = max(most, len(each))
	}
	return most
}

func deriveAngelMaximums(mach Machine) []int {
	idToIDs := make(map[int]map[int]int, len(mach.Joltages))
	for i := range mach.Joltages {
		idToIDs[i] = map[int]int{}
	}
	for _, button := range mach.Buttons {
		for _, id := range button {
			for _, each := range button {
				idToIDs[id][each] = 0
			}
		}
	}
	maxValues := make([]int, len(mach.Joltages))
	for id, ids := range idToIDs {
		minValue := math.MaxInt
		for joltID := range ids {
			minValue = min(minValue, mach.Joltages[joltID])
		}
		maxValues[id] = minValue
	}
	return maxValues
}

func deriveMaxClickCount(joltages []int, joltageIDs []int) int {
	minimus := math.MaxInt
	for _, each := range joltageIDs {
		minimus = min(minimus, joltages[each])
	}
	return minimus
}

//revive:disable-next-line:function-length
func deriveFewestJoltageClicks8(mach Machine) int {
	mach.Buttons = sortButtons(mach.Buttons)
	tbl := newTable(mach)
	joltageSum := sum(mach.Joltages)
	minJoltageValue := minJoltage(mach.Joltages)
	clickIndexes := make([]int, 0, joltageSum)
	index := 0
	minCount := math.MaxInt
	back := false
	for roundIndex := range math.MaxInt {
		if index >= len(clickIndexes) {
			clickIndexes = append(clickIndexes, -1)
		}
		lastButtonIndex := clickIndexes[index]
		if back && lastButtonIndex >= 0 {
			tbl.undo(lastButtonIndex)
		}
		if index >= minCount || tbl.stuck {
			index--
			if index < 0 {
				panic("Impossible for index<0")
			}
			back = true
			continue
		}
		back = false
		buttonIndex, found := tbl.getClickable(lastButtonIndex)
		if found {
			clickIndexes[index] = buttonIndex
			tbl.click(buttonIndex)
			index++
			continue
		}
		clickIndexes[index] = -1
		if tbl.done() {
			shared.Logger.Info("Found a match.",
				"count", index,
				"click indexes", clickIndexes,
				"round", roundIndex)
			minCount = min(minCount, index)
			if minCount == minJoltageValue {
				break
			}
		}
		index--
		if index < 0 {
			break
		}
		back = true
	}
	shared.Logger.Info("Derived fewest clicks.",
		"click indexes", clickIndexes,
		"count", minCount)
	return minCount
}

func sum(ints []int) int {
	var total int
	for _, each := range ints {
		total += each
	}
	return total
}

type clickEffect struct {
	joltageIndexes []int
	buttonIndexes  []int
}

type table struct {
	clicks           []int
	joltages         []int
	buttons          [][]int
	joltageToButtons [][]int
	effects          []clickEffect
	stuck            bool
}

func newTable(mach Machine) *table {
	clicks := make([]int, len(mach.Buttons))
	for i, button := range mach.Buttons {
		minJolt := math.MaxInt
		for _, joltIndex := range button {
			minJolt = min(minJolt, mach.Joltages[joltIndex])
		}
		clicks[i] = minJolt
	}
	joltageToButtons := make([][]int, len(mach.Joltages))
	for joltageIndex := range mach.Joltages {
		for buttonIndex, button := range mach.Buttons {
			if slices.Contains(button, joltageIndex) {
				joltageToButtons[joltageIndex] = append(joltageToButtons[joltageIndex], buttonIndex)
			}
		}
	}
	effects := make([]clickEffect, len(mach.Buttons))
	tbl := &table{
		clicks:           clicks,
		joltages:         mach.Joltages,
		buttons:          mach.Buttons,
		joltageToButtons: joltageToButtons,
	}
	for i := range mach.Buttons {
		effects[i] = tbl.deriveClickEffects(i)
	}
	tbl.effects = effects
	return tbl
}

func (v *table) getFirstClickable() (int, bool) {
	return v.getClickable(-1)
}

func (v *table) getClickable(after int) (int, bool) {
	for i := after + 1; i < len(v.clicks); i++ {
		each := v.clicks[i]
		if each > 0 {
			return i, true
		}
	}
	return -1, false
}

func (v *table) click(index int) {
	effect := v.effects[index]
	atRisk := false
	for _, each := range effect.joltageIndexes {
		v.joltages[each]--
		if !v.stuck && v.joltages[each] == 0 {
			atRisk = true
		}
	}
	var atRiskButtonIndexes []int
	for _, i := range effect.buttonIndexes {
		zeroClicks := v.clicks[i] == 0
		// Since the button already had zero possible clicks, there's no point in calculating the
		// click count again because it's impossible to change that by click any other button.
		if zeroClicks {
			continue
		}
		btn := v.buttons[i]
		value := math.MaxInt
		positives := false
		for _, j := range btn {
			val := v.joltages[j]
			if val > 0 {
				positives = true
			}
			value = min(value, val)
		}
		v.clicks[i] = value
		if atRisk && positives && value == 0 {
			atRiskButtonIndexes = append(atRiskButtonIndexes, i)
		}
	}
	if !atRisk || len(atRiskButtonIndexes) == 0 {
		return
	}
	v.stuck = v.isStuck(atRiskButtonIndexes)
}

func (v *table) isStuck(atRiskButtonIndexes []int) bool {
	positiveJoltageIndexes := map[int]int{}
	for _, buttonIndex := range atRiskButtonIndexes {
		for _, joltageIndex := range v.buttons[buttonIndex] {
			if v.joltages[joltageIndex] > 0 {
				positiveJoltageIndexes[joltageIndex] = 0
			}
		}
	}
	if len(positiveJoltageIndexes) == 0 {
		return false
	}

outerLoop:

	for joltageIndex := range positiveJoltageIndexes {
		for _, buttonIndex := range v.joltageToButtons[joltageIndex] {
			if v.clicks[buttonIndex] > 0 {
				continue outerLoop
			}
		}
		return true
	}
	return false
}

func (v *table) undo(index int) {
	effect := v.effects[index]
	riskChange := false
	for _, each := range effect.joltageIndexes {
		v.joltages[each]++
		if v.stuck && v.joltages[each] == 1 {
			riskChange = true
		}
	}
	var buttonIndexes []int
	for _, i := range effect.buttonIndexes {
		btn := v.buttons[i]
		value := math.MaxInt
		for _, j := range btn {
			value = min(value, v.joltages[j])
		}
		v.clicks[i] = value
		if v.stuck && riskChange && value == 1 {
			buttonIndexes = append(buttonIndexes, i)
		}
	}
	if !riskChange || len(buttonIndexes) == 0 {
		return
	}
	v.stuck = v.isStuck(buttonIndexes)
}

func (v *table) deriveClickEffects(index int) clickEffect {
	reducedJoltageIndexes := v.buttons[index]
	handled := gent.NewSet[int]()
	affectedButtons := make([]int, 0, len(v.buttons))
	for _, each := range reducedJoltageIndexes {
		buttons := v.joltageToButtons[each]
		for _, button := range buttons {
			if !handled.Add(button) {
				continue
			}
			affectedButtons = append(affectedButtons, button)
		}
	}
	return clickEffect{
		joltageIndexes: reducedJoltageIndexes,
		buttonIndexes:  affectedButtons,
	}
}

func (v *table) done() bool {
	return !slices.ContainsFunc(v.joltages, func(i int) bool { return i != 0 })
}

func shiftButtons(buttons [][]int, count int) [][]int {
	shifted := make([][]int, len(buttons))
	for i := range buttons {
		from := (i + count) % len(buttons)
		shifted[i] = buttons[from]
	}
	return shifted
}

type button struct {
	id             int
	joltageIndexes []int
	clicks         int
}

type demon struct {
	joltageID      int
	joltage        int
	buttons        []*button
	lockedToButton int
}

func (v *demon) extractExclusiveButtons() []int {
	var ids []int
	if len(v.buttons) == 1 {
		ids = append(ids, v.buttons[0].id)
	}
	return ids
}

func setupDemons(mach Machine) []*demon {
	demons := deriveDemons(mach)
	var buttonIDs []int
	for _, each := range demons {
		ids := each.extractExclusiveButtons()
		if len(ids) > 0 {
			buttonIDs = append(buttonIDs, ids...)
		}
	}
	for len(buttonIDs) > 0 {
		buttonIDs = deriveButtonClickCounts(mach.Buttons, demons, buttonIDs)
	}
	// Simple demons first.
	slices.SortFunc(demons, func(a, b *demon) int {
		first := deriveJoltageComplexity(a.joltage, len(a.buttons))
		second := deriveJoltageComplexity(b.joltage, len(b.buttons))
		return first - second
	})
	// Locked demons first.
	slices.SortStableFunc(demons, func(a, b *demon) int {
		if a.lockedToButton >= 0 && b.lockedToButton < 0 {
			return -1
		}
		if a.lockedToButton < 0 && b.lockedToButton >= 0 {
			return 1
		}
		return 0
	})
	return demons
}

//revive:disable-next-line:cognitive-complexity
//revive:disable-next-line:function-length
//revive:disable-next-line:cyclomatic
func deriveButtonClickCounts(buttons [][]int, demons []*demon, buttonIDs []int) []int {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Derive button click counts.", "button IDs", buttonIDs)
	}
	nextDemonIDs := gent.NewSet[int]()
buttonIDLoop:
	for _, buttonID := range buttonIDs {
		buttonDemonIDs := buttons[buttonID]
		if shared.IsDebugEnabled() {
			shared.Logger.Debug(
				"Button's demon IDs.",
				"button ID", buttonID,
				"demon IDs", buttonDemonIDs)
		}
		for _, demonID := range buttonDemonIDs {
			demon := demons[demonID]
			if demon.lockedToButton < 0 && len(demon.buttons) == 1 {
				demon.buttons[0].clicks = demon.joltage
				demon.lockedToButton = buttonID
				nextDemonIDs.Add(demonID)
				if shared.IsDebugEnabled() {
					shared.Logger.Debug("Demon and button locked.")
				}
				continue buttonIDLoop
			}
		}
		// Is this button solely responsible for a specific demon? If it is, lock it to it. The
		// button points to demons and those demons are pointed to by other buttons. If for any
		// given demon all other buttons pointing to it have been locked, then this buttons is
		// solely responsible for taking care of shooting it down. In that case we need to
		// investigate click counts from the other locked buttons, and reduce the sum from the
		// demon's joltage, and lock this button for that.
		targetDemonID := whichDemonIsButtonSolelyResponsibleFor(buttonID, demons, buttonDemonIDs)
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Target demon derived.", "ID", targetDemonID)
		}
		if targetDemonID < 0 {
			continue
		}
		lockedDemon := demons[targetDemonID]
		existingClicks := countExistingLockedClicks(lockedDemon)
		if shared.IsDebugEnabled() {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug(
					"Existing clicks derived.",
					"existing", existingClicks,
					"joltage", lockedDemon.joltage)
			}
		}
		if existingClicks >= lockedDemon.joltage {
			continue
		}
		var lockedButton *button
		for _, btn := range lockedDemon.buttons {
			if btn.id == buttonID {
				lockedButton = btn
				break
			}
		}
		lockedButton.clicks = lockedDemon.joltage - existingClicks
		lockedDemon.lockedToButton = lockedButton.id
		for _, each := range buttonDemonIDs {
			if each != targetDemonID {
				nextDemonIDs.Add(each)
			}
		}
	}
	if nextDemonIDs.Len() == 0 {
		return nil
	}
	nextButtonIDSet := gent.NewSet[int]()
	nextDemonIDs.ForEachAll(func(nextDemonID int) {
		affectedDemonIDs := gent.NewSet[int]()
		for _, btn := range demons[nextDemonID].buttons {
			for _, demonID := range buttons[btn.id] {
				if demons[demonID].lockedToButton < 0 {
					affectedDemonIDs.Add(demonID)
				}
			}
		}
		affectedDemonIDs.ForEachAll(func(demonID int) {
			for _, btn := range demons[demonID].buttons {
				if btn.clicks < 0 {
					nextButtonIDSet.Add(btn.id)
				}
			}
		})
	})
	if nextButtonIDSet.Len() == 0 {
		return nil
	}
	ids := nextButtonIDSet.ToSlice()
	// Just to have deterministic behavior.
	slices.Sort(ids)
	return ids
}

func deriveDemons(mach Machine) []*demon {
	demons := make([]*demon, len(mach.Joltages))
	buttons := make(map[int]*button, len(mach.Buttons))
	joltageIDToButtonIDs := make(map[int][]int, len(mach.Joltages))
	for i := range joltageIDToButtonIDs {
		joltageIDToButtonIDs[i] = []int{}
	}
	for i, joltageIndexes := range mach.Buttons {
		buttons[i] = &button{
			id:             i,
			joltageIndexes: joltageIndexes,
			clicks:         -1,
		}
		for _, joltIndex := range joltageIndexes {
			joltageIDToButtonIDs[joltIndex] = append(joltageIDToButtonIDs[joltIndex], i)
		}
	}
	for i := range mach.Joltages {
		var btns []*button
		for _, each := range joltageIDToButtonIDs[i] {
			btns = append(btns, buttons[each])
		}
		demons[i] = &demon{
			joltageID:      i,
			joltage:        mach.Joltages[i],
			buttons:        btns,
			lockedToButton: -1,
		}
	}
	return demons
}

func deriveJoltageComplexity(joltage, count int) int {
	return nk(joltage+count-1, count-1)
}

func nk(n, k int) int {
	nom := 1
	for i := n; i >= (n - k + 1); i-- {
		nom *= i
	}
	dem := 1
	for i := k; i >= 1; i-- {
		dem *= i
	}
	return nom / dem
}

func whichDemonIsButtonSolelyResponsibleFor(
	buttonID int,
	demons []*demon,
	buttonDemonIDs []int,
) int {
	for _, each := range buttonDemonIDs {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Consider.", "each", each)
		}
		if isButtonSolelyResponsibleForDemon(buttonID, demons[each]) {
			return each
		}
	}
	return -1
}

func isButtonSolelyResponsibleForDemon(buttonID int, aDemon *demon) bool {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug(
			"Is button solely responsible for given demon?",
			"button ID", buttonID,
			"demon", *aDemon)
	}
	// If all buttons pointing to the demon have been locked, then this button is solely responsible
	// for the demon.
	for _, btn := range aDemon.buttons {
		if btn.id != buttonID && btn.clicks < 0 {
			return false
		}
	}
	return true
}

func countExistingLockedClicks(aDemon *demon) int {
	var clicks int
	for _, btn := range aDemon.buttons {
		if btn.clicks >= 0 {
			clicks += btn.clicks
		}
	}
	return clicks
}

func permute[T any](items []T, yield func([]T) bool) bool {
	n := len(items)
	if n == 0 {
		return yield([]T{})
	}

	a := append([]T(nil), items...)
	var walk func(int) bool
	walk = func(i int) bool {
		if i == n-1 {
			p := append([]T(nil), a...)
			return yield(p)
		}

		for j := i; j < n; j++ {
			a[i], a[j] = a[j], a[i]
			if !walk(i + 1) {
				return false
			}
			a[i], a[j] = a[j], a[i]
		}

		return true
	}
	return walk(0)
}

func generateButtonCombinations(buttons [][]int, joltageCount int, cb func([]int)) {
	indexes := make([]int, len(buttons))
	for i := range buttons {
		indexes[i] = i
	}
	slices.SortStableFunc(indexes, func(i, j int) int {
		return len(buttons[j]) - len(buttons[i])
	})
	generateCombinations(indexes, func(combination []int) {
		if coversAllJoltages(buttons, combination, joltageCount) {
			cb(combination)
		}
	})
}

func generateCombinations[T any](input []T, cb func(combination []T)) {
	n := len(input)
	if n == 0 {
		return
	}

	total := int(math.Pow(2, float64(n)))
	for i := range total {
		var combo []T
		for j := range n {
			if (i>>j)&1 == 1 {
				combo = append(combo, input[j])
			}
		}
		cb(combo)
	}
}

func coversAllJoltages(buttons [][]int, indexes []int, joltageCount int) bool {
	m := make(map[int]int, joltageCount)
	for _, i := range indexes {
		btn := buttons[i]
		for _, j := range btn {
			m[j] = 0
		}
		if len(m) >= joltageCount {
			if len(m) > joltageCount {
				panic("impossible: len(m) > joltageCount")
			}
			return true
		}
	}
	return false
}

type clickCounter struct {
	buttonIDToCount map[int]int
	count           int
}

type hop struct {
	// Total is the total overlall reduction in joltages.
	total int
	// ClickIndex is the index of the clicked button. Always >=0.
	clickIndex int
	// unclickIndex is the index of the unclicked button. If <0, no unclick to perform.
	unclickIndex int
}

type blockedButton struct {
	buttonID              int
	joltageIndexes        []int
	blockedJoltageIndexes []int
}

func deriveBestDoubleClicks(joltages []int, buttons [][]int, counter *clickCounter) []hop {
	var hops []hop
	// Filter blocked buttons.
	// For the blocked buttons, filter those that can be enabled by unclicking.
	// For those that remain, filter the combinations that have total effect >=0.
	for _, blocked := range filterBlockedButtons(joltages, buttons) {
		for _, unblockingID := range extractUnblockingButtonIDs(buttons, counter, blocked) {
			total := len(blocked.joltageIndexes) - len(buttons[unblockingID])
			if total < 0 {
				continue
			}
			hops = append(hops, hop{
				total:        total,
				clickIndex:   blocked.buttonID,
				unclickIndex: unblockingID,
			})
		}
	}
	return hops
}

func sortHops(hops []hop) {
	// Stable just to have something deterministic to test.
	slices.SortStableFunc(hops, func(a, b hop) int {
		return b.total - a.total
	})
}

func filterBlockedButtons(joltages []int, buttons [][]int) []blockedButton {
	var blocked []blockedButton
	for i, btn := range buttons {
		var blockedIndexes []int
		for _, j := range btn {
			if joltages[j] == 0 {
				blockedIndexes = append(blockedIndexes, j)
			}
		}
		if len(blockedIndexes) > 0 {
			blocked = append(blocked, blockedButton{
				buttonID:              i,
				joltageIndexes:        btn,
				blockedJoltageIndexes: blockedIndexes,
			})
		}
	}
	return blocked
}

func extractUnblockingButtonIDs(buttons [][]int, counter *clickCounter, btn blockedButton) []int {
	var ids []int
	for buttonID, clickCount := range counter.buttonIDToCount {
		if buttonID == btn.buttonID || clickCount <= 0 {
			continue
		}
		button := buttons[buttonID]
		if isSuperset(button, btn.blockedJoltageIndexes) {
			ids = append(ids, buttonID)
		}
	}
	return ids
}

func isSuperset(container, nested []int) bool {
	containerSize, nestedSize := len(container), len(nested)
	if nestedSize > containerSize {
		return false
	}
	if containerSize == 0 || nestedSize == 0 {
		panic("neither in isSuperset should be empty")
	}
	u := union(container, nested)
	return len(u) == len(nested)
}

func union(a, b []int) []int {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	kid, dad := a, b
	if len(a) > len(b) {
		dad, kid = b, a
	}
	m := make(map[int]struct{}, len(dad))
	for i := range dad {
		m[dad[i]] = struct{}{}
	}
	var common []int
	for i := range kid {
		if _, ok := m[kid[i]]; ok {
			common = append(common, kid[i])
		}
	}
	return common
}

func deriveFewestJoltageClicks9(mach Machine) int {
	hive := createHiveMind(mach)
	return hive.derive()
}

type valueRestrictor struct {
	minimum    int
	maximum    int
	multiplier int
}

type expressionRestrictor struct {
	sum           int
	buttonIndexes []int
}

func (v *expressionRestrictor) isValid(shelf *clickShelf) bool {
	var total int
	nilsExist := false
	for _, i := range v.buttonIndexes {
		if shelf.clicks[i] == nil {
			nilsExist = true
			continue
		}
		total += *shelf.clicks[i]
	}
	if total > v.sum {
		return false
	}
	return nilsExist || total == v.sum
}

type derivedValue struct {
	i      int
	clicks int
}

func toCommaSeparatedString(s []*int) string {
	if s == nil {
		return "nil"
	}
	res := "["
	for i, value := range s {
		if i != 0 {
			res += ","
		}
		if value == nil {
			res += "nil"
		} else {
			res += strconv.Itoa(*value)
		}
	}
	res += "]"
	return res
}

//revive:disable-next-line:cyclomatic
func (v *expressionRestrictor) derive(shelf *clickShelf) (values []derivedValue) {
	buttonCount := len(v.buttonIndexes)
	nilIndexes := make([]int, 0, 5)
	var currentTotal int
	for _, i := range v.buttonIndexes {
		if shelf.clicks[i] == nil {
			nilIndexes = append(nilIndexes, i)
		} else {
			currentTotal += *shelf.clicks[i]
		}
	}

	// Case: this is the only button that affects its joltage.
	if buttonCount == 1 && len(nilIndexes) == 1 {
		val := derivedValue{
			i:      nilIndexes[0],
			clicks: v.sum,
		}
		values = append(values, val)
		return
	}

	// Case: if current total is v.sum, the rest are zeroes.
	if buttonCount > 1 && currentTotal == v.sum {
		for _, each := range nilIndexes {
			val := derivedValue{
				i:      each,
				clicks: 0,
			}
			values = append(values, val)
		}
		return
	}

	if buttonCount <= 1 {
		return
	}

	// Case: this is the only remaining button to affect its joltage.
	remaining := v.sum - currentTotal
	if len(nilIndexes) == 1 && remaining >= 0 {
		val := derivedValue{
			i:      nilIndexes[0],
			clicks: remaining,
		}
		values = append(values, val)
	}
	return
}

func (v *expressionRestrictor) fillButtonRates(rates []int) {
	mul := len(v.buttonIndexes) - 1
	for _, buttonIndex := range v.buttonIndexes {
		current := rates[buttonIndex]
		rate := mul * v.sum
		if current == 0 || (rate < current && rate != 0) {
			shared.Logger.Debug("Set rate.", "button index", buttonIndex, "rate", rate)
			rates[buttonIndex] = rate
		}
	}
}

type hiveMind struct {
	machineID             int
	restrictors           []valueRestrictor
	expressionRestrictors []expressionRestrictor
	joltages              []int
	buttons               [][]int
	buttonRates           []int
}

func createHiveMind(mach Machine) *hiveMind {
	restrictors := make([]valueRestrictor, len(mach.Buttons))
	joltToButtonIndex := map[int][]int{}
	withoutMinimum := make([]int, len(mach.Buttons))
	for i, button := range mach.Buttons {
		withoutMinimum[i] = i
		least := math.MaxInt
		for _, j := range button {
			least = min(least, mach.Joltages[j])
			joltToButtonIndex[j] = append(joltToButtonIndex[j], i)
		}
		restrictors[i] = valueRestrictor{maximum: least, multiplier: len(button)}
	}

	var expressionRestrictors []expressionRestrictor
	for i, value := range mach.Joltages {
		buttonIndexes := joltToButtonIndex[i]
		rest := expressionRestrictor{sum: value, buttonIndexes: buttonIndexes}
		expressionRestrictors = append(expressionRestrictors, rest)
	}

	foundSubsets := true
	for foundSubsets {
		foundSubsets = findSubsets(expressionRestrictors)
	}

	buttonRates := make([]int, len(mach.Buttons))
	for _, each := range expressionRestrictors {
		each.fillButtonRates(buttonRates)
		for _, i := range each.buttonIndexes {
			if restrictors[i].maximum > each.sum {
				restrictors[i].maximum = each.sum
			}
		}
	}

	for i := range restrictors {
		var minimum int
		for _, joltIndex := range mach.Buttons[i] {
			expressor := expressionRestrictors[joltIndex]
			indexes := expressor.buttonIndexes
			// if len(indexes) <= 1 {
			// 	continue
			// }
			var total int
			for _, each := range indexes {
				if each != i {
					total += restrictors[each].maximum
				}
			}
			remainder := expressor.sum - total
			if remainder >= 0 && remainder > minimum {
				minimum = remainder
			}
		}
		restrictors[i].minimum = minimum
	}

	aHive := &hiveMind{
		restrictors:           restrictors,
		expressionRestrictors: expressionRestrictors,
		joltages:              mach.Joltages,
		buttons:               mach.Buttons,
		buttonRates:           buttonRates,
		machineID:             mach.id,
	}
	shared.Logger.Info("HiveMind created.",
		"hive", *aHive)
	return aHive
}

type clickShelf struct {
	clicks        []*int
	clickSum      int
	currentSum    int
	targetSum     int
	buttonWeights []int
}

func createClickShelf(size int, joltages []int, buttons [][]int) *clickShelf {
	clicks := make([]*int, size)
	var targetSum int
	for _, each := range joltages {
		targetSum += each
	}
	buttonWeights := make([]int, len(buttons))
	for i, each := range buttons {
		buttonWeights[i] = len(each)
	}
	return &clickShelf{clicks: clicks, targetSum: targetSum, buttonWeights: buttonWeights}
}

func (v *clickShelf) setp(index int, value *int) {
	var alpha, omega int
	if v.clicks[index] != nil {
		alpha = *v.clicks[index]
	}
	if value != nil {
		omega = *value
	}
	delta := omega - alpha
	v.clickSum += delta
	v.currentSum += v.buttonWeights[index] * delta
	v.clicks[index] = value
}

func (v *clickShelf) set(index, value int) {
	v.setp(index, &value)
}

func (v *clickShelf) replace(clicks []*int) {
	v.clicks = clicks
	var sum int
	for i, each := range clicks {
		if each != nil {
			sum += v.buttonWeights[i] * *each
		}
	}
	v.currentSum = sum
	v.clickSum = sumClicks(clicks)
}

func (v *clickShelf) isSumCorrect() bool {
	return v.currentSum == v.targetSum
}

func (v *hiveMind) derive() int {
	shared.Logger.Info("HiveMind: start derive.")
	shelf := createClickShelf(len(v.restrictors), v.joltages, v.buttons)
	v.deriveRestrictorsUntilEnd(shelf)
	if result, ok := v.deriveResult(shelf); ok {
		return result
	}
	sortedButtonIndexes := v.sortButtons()
	shared.Logger.Info("Start diving.",
		"button indexes", sortedButtonIndexes,
		"clicks", shelf.clicks)
	fewest := math.MaxInt
	v.dive(sortedButtonIndexes, 0, shelf, &fewest)
	return fewest
}

func (v *hiveMind) sortButtons() []int {
	indexes := make([]int, len(v.buttonRates))
	for i := range len(indexes) {
		indexes[i] = i
	}
	slices.SortStableFunc(indexes, func(a, b int) int {
		return v.buttonRates[a] - v.buttonRates[b]
	})
	return indexes
}

//revive:disable-next-line:cyclomatic
func (v *hiveMind) dive(
	buttonIndexes []int,
	currentButtonIndex int,
	shelf *clickShelf,
	fewestClicks *int,
) {
	realButtonIndex := buttonIndexes[currentButtonIndex]
	if shelf.clicks[realButtonIndex] != nil {
		if (currentButtonIndex + 1) >= len(shelf.clicks) {
			panic("boo boo")
		}
		v.dive(buttonIndexes, currentButtonIndex+1, shelf, fewestClicks)
		return
	}
	remaining := *fewestClicks - shelf.clickSum
	if remaining <= 0 {
		return
	}
	top := v.restrictors[realButtonIndex].maximum
	if top >= remaining {
		top = remaining - 1
		if top < 0 {
			return
		}
	}
	for count := top; count >= v.restrictors[realButtonIndex].minimum; count-- {
		previous := shelf.clicks[realButtonIndex]
		shelf.setp(realButtonIndex, &count)
		if !v.validate(shelf) {
			shelf.setp(realButtonIndex, previous)
			continue
		}
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Try combo.", "before derive", toCommaSeparatedString(shelf.clicks))
		}
		originalClicks := copyClicks(shelf.clicks, realButtonIndex, previous)
		v.deriveRestrictorsUntilEnd(shelf)
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Try combo.",
				"after derive", toCommaSeparatedString(shelf.clicks),
				"click count", shelf.clickSum)
		}
		if shelf.clickSum >= *fewestClicks {
			shelf.replace(originalClicks)
			continue
		}
		if v.validate(shelf) {
			if result, ok := v.deriveResult(shelf); ok {
				if result < *fewestClicks {
					*fewestClicks = result
				}
			} else if (currentButtonIndex + 1) < len(shelf.clicks) {
				v.dive(buttonIndexes, currentButtonIndex+1, shelf, fewestClicks)
			}
		}
		shelf.replace(originalClicks)
	}
}

func (v *hiveMind) deriveRestrictorsUntilEnd(shelf *clickShelf) {
	derived := true
	for derived {
		derived = v.deriveRestrictors(shelf)
	}
}

func copyClicks(clicks []*int, index int, value *int) []*int {
	result := make([]*int, len(clicks))
	for i, each := range clicks {
		if i == index {
			result[i] = copyIntp(value)
		} else {
			result[i] = copyIntp(each)
		}
	}
	return result
}

func copyIntp(i *int) *int {
	if i == nil {
		return i
	}
	value := *i
	return &value
}

func (v *hiveMind) validate(shelf *clickShelf) bool {
	for _, each := range v.expressionRestrictors {
		if !each.isValid(shelf) {
			return false
		}
	}
	return true
}

func (v *hiveMind) deriveRestrictors(shelf *clickShelf) bool {
	success := false
	for _, each := range v.expressionRestrictors {
		derivedValues := each.derive(shelf)
		success = len(derivedValues) > 0
		for i := range derivedValues {
			derivedVal := derivedValues[i]
			shelf.setp(derivedVal.i, &derivedVal.clicks)
		}
	}
	return success
}

func (v *hiveMind) deriveResult(shelf *clickShelf) (int, bool) {
	if !shelf.isSumCorrect() {
		return 0, false
	}
	for _, each := range shelf.clicks {
		if each == nil {
			return 0, false
		}
	}
	for _, each := range v.expressionRestrictors {
		if !each.isValid(shelf) {
			return 0, false
		}
	}
	joltages := make([]int, len(v.joltages))
	copy(joltages, v.joltages)
	var total int
	for i, each := range shelf.clicks {
		decInPlaceN(joltages, v.buttons[i], *each)
		total += *each
	}
	if isZeroes(joltages) {
		shared.Logger.Info(
			"Found.",
			"clicks", toCommaSeparatedString(shelf.clicks),
			"sum", shelf.clickSum,
			"machine ID", v.machineID,
		)
		return total, true
	}
	return 0, false
}

func sumClicks(clicks []*int) int {
	var sum int
	for _, each := range clicks {
		if each != nil {
			sum += *each
		}
	}
	return sum
}

func findSubsets(restrictors []expressionRestrictor) bool {
	found := false
	for i := range restrictors {
		for j := range restrictors {
			if i == j {
				continue
			}
			parent := restrictors[i]
			child := restrictors[j]
			if len(parent.buttonIndexes) <= len(child.buttonIndexes) {
				continue
			}
			rest := splitSubset(parent.buttonIndexes, child.buttonIndexes)
			if len(rest) == 0 {
				continue
			}
			parent.buttonIndexes = rest
			parent.sum = parent.sum - child.sum
			restrictors[i] = parent
			found = true
		}
	}
	return found
}

func splitSubset(parent, child []int) []int {
	var without []int
	for i := range child {
		index := slices.Index(parent, child[i])
		if index < 0 {
			return nil
		}
		without = append(without, index)
	}
	result := make([]int, 0, len(parent)-len(child))
	for i := range parent {
		if slices.Contains(without, i) {
			continue
		}
		result = append(result, parent[i])
	}
	return result
}
