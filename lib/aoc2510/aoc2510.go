package aoc2510

import (
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

func DeriveFewestClicks(lines []string, indicator bool) int {
	shared.Logger.Info("Derive fewest clicks.", "machine count", len(lines), "indicator", indicator)
	var clicks int
	defer func() {
		shared.Logger.Info("Fewest clicks derived.", "count", clicks)
	}()
	deriveFunc := deriveFewestStateClicks
	if !indicator {
		deriveFunc = deriveFewestJoltageClicks
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

func decInPlaceN(joltages []int, button []int, n int) {
	for _, each := range button {
		joltages[each] -= n
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

func isZeroes(ints []int) bool {
	for i := range ints {
		if ints[i] != 0 {
			return false
		}
	}
	return true
}

func deriveFewestJoltageClicks(mach Machine) int {
	return createHiveMind(mach).derive()
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
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Set rate.", "button index", buttonIndex, "rate", rate)
			}
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
	shared.Logger.Info("HiveMind created.", "hive", *aHive)
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
	shared.Logger.Info("HiveMind: start derive.", "machine ID", v.machineID)
	shelf := createClickShelf(len(v.restrictors), v.joltages, v.buttons)
	v.deriveRestrictorsUntilEnd(shelf)
	if result, ok := v.deriveResult(shelf); ok {
		return result
	}
	sortedButtonIndexes := v.sortButtons()
	shared.Logger.Info("Start diving.",
		"button indexes", sortedButtonIndexes,
		"clicks", shelf.clicks,
		"machine ID", v.machineID)
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
