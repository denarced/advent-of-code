package aoc2322

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	houndRunning houndState = iota
	houndWaiting
	houndSleeping
)

type houndState int

func (v houndState) String() string {
	switch v {
	case houndRunning:
		return "running"
	case houndWaiting:
		return "waiting"
	case houndSleeping:
		return "sleeping"
	default:
		panic("unhandled enum case")
	}
}

func CountBricksFromLines(lines []string) int {
	bricks := parseLines(lines)
	shared.Logger.Info(
		"Lines parsed, count bricks that can be disintegrated.",
		"brick count", len(bricks),
		"line count", len(lines))
	byLowZ, byHighZ := createSearchIndexes(bricks)
	descend(bricks, byLowZ, byHighZ)
	slackers := findSlackers(bricks, byLowZ, byHighZ)
	total := slackers.Count()
	shared.Logger.Info("Slackers found.", "count", total, "total brick count", len(bricks))
	return total
}

func parseLines(lines []string) []brick {
	bricks := make([]brick, len(lines))
	for i, each := range lines {
		bricks[i] = parseLine(each)
	}
	return bricks
}

type brick struct {
	start, end coordinate
}

func (v brick) String() string {
	return fmt.Sprintf(
		"x(%d-%d),y(%d-%d),z(%d-%d)",
		v.start.x, v.end.x,
		v.start.y, v.end.y,
		v.start.z, v.end.z)
}

type coordinate struct {
	x, y, z int
}

func parseLine(s string) brick {
	pieces := strings.Split(s, "~")
	claimStringCount("parse line", pieces, 2)
	low, high := parseCoordinate(pieces[0]), parseCoordinate(pieces[1])
	claimAscending(low.x, high.x)
	claimAscending(low.y, high.y)
	claimAscending(low.z, high.z)
	return brick{start: low, end: high}
}

func claimAscending(values ...int) {
	var previous *int
	for _, each := range values {
		if previous == nil {
			local := each
			previous = &local
			continue
		}
		if each < *previous {
			panic(fmt.Sprintf("not ascending: %d, %d", *previous, each))
		}
	}
}

func parseCoordinate(s string) coordinate {
	pieces := strings.Split(s, ",")
	claimStringCount("split coordinates", pieces, 3)
	values := make([]int, 3)
	for i := range pieces {
		values[i] = gent.OrPanic2(strconv.Atoi(pieces[i]))("int parse failed")
	}
	return coordinate{
		x: values[0],
		y: values[1],
		z: values[2],
	}
}

func claimStringCount(message string, s []string, count int) {
	if len(s) != count {
		panic(fmt.Sprintf("incorrect count: %d, expected: %d, message: %s", len(s), count, message))
	}
}

func createSearchIndexes(bricks []brick) (byLow, byHigh map[int][]int) {
	byLow = map[int][]int{}
	byHigh = map[int][]int{}
	for i, each := range bricks {
		low := each.start.z
		byLow[low] = append(byLow[low], i)
		high := each.end.z
		byHigh[high] = append(byHigh[high], i)
	}
	return
}

func createAboveToBelowMapping(
	bricks []brick,
	aboveBrickIndexes, belowBrickIndexes []int,
	addSlackerFunc func(aBrick brick) bool,
) map[int][]int {
	aboveToBelow := map[int][]int{}
	for _, belowBrickIndex := range belowBrickIndexes {
		// This brick is the focus here: does it support nothing or is it one of many for bricks
		// above that it supports.
		lowBrick := bricks[belowBrickIndex]
		if len(aboveBrickIndexes) == 0 {
			if addSlackerFunc(lowBrick) {
				shared.Logger.Info(
					"New slacker found.",
					"brick", lowBrick,
					"index", belowBrickIndex)
			}
			continue
		}
		overlapIndexes := findOverlaps(bricks, aboveBrickIndexes, lowBrick)
		// Nothing above.
		if len(overlapIndexes) == 0 {
			if addSlackerFunc(lowBrick) {
				shared.Logger.Info(
					"New slacker found. No bricks above.",
					"brick", lowBrick,
					"index", belowBrickIndex)
			}
			continue
		}
		for _, overlapIndex := range overlapIndexes {
			aboveToBelow[overlapIndex] = append(aboveToBelow[overlapIndex], belowBrickIndex)
		}
	}
	return aboveToBelow
}

func findSlackers(bricks []brick, byLow, byHigh map[int][]int) *gent.Set[brick] {
	slackers := gent.NewSet[brick]()
	for _, eachHighZ := range sortKeys(byHigh) {
		aboveToBelow := createAboveToBelowMapping(
			bricks,
			byLow[eachHighZ+1],
			byHigh[eachHighZ],
			func(aBrick brick) bool {
				return slackers.Add(aBrick)
			},
		)
		singleParents := gent.NewSet[brick]()
		for _, belowIndexes := range aboveToBelow {
			if len(belowIndexes) > 1 {
				for _, each := range belowIndexes {
					if slackers.Add(bricks[each]) {
						shared.Logger.Info(
							"New slacker found. Others would support.",
							"brick", bricks[each],
							"index", each)
					}
				}
			} else if len(belowIndexes) == 1 {
				singleParents.Add(bricks[belowIndexes[0]])
			}
		}
		singleParents.ForEachAll(func(each brick) {
			slackers.Remove(each)
		})
	}
	return slackers
}

func sortKeys(m map[int][]int) []int {
	keys := make([]int, len(m))
	var i int
	for each := range m {
		keys[i] = each
		i++
	}
	slices.Sort(keys)
	return keys
}

func findOverlaps(bricks []brick, indexes []int, base brick) (resultIndexes []int) {
	for _, i := range indexes {
		candidate := bricks[i]
		if isOverlap(base, candidate) {
			resultIndexes = append(resultIndexes, i)
		}
	}
	return
}

func isOverlap(a, b brick) bool {
	return rangesOverlap(a.start.x, a.end.x, b.start.x, b.end.x) &&
		rangesOverlap(a.start.y, a.end.y, b.start.y, b.end.y)
}

func rangesOverlap(alphaFrom, alphaTo, omegaFrom, omegaTo int) bool {
	if alphaTo < alphaFrom {
		panic("alpha to<from")
	}
	if omegaTo < omegaFrom {
		panic("omega to<from")
	}
	if alphaTo < omegaFrom || omegaTo < alphaFrom {
		return false
	}
	return true
}

func descend(bricks []brick, byLowZ, byHighZ map[int][]int) {
	floor := createFloor(bricks)
	for _, z := range extractSortedZ(bricks) {
		brickIndexes := byLowZ[z]
		for _, brickIndex := range brickIndexes {
			candidate := bricks[brickIndex]
			descendTo, shouldDescend := canDescend(floor, candidate)
			if shouldDescend {
				zDiff := candidate.start.z - descendTo
				if shared.IsDebugEnabled() {
					shared.Logger.Debug(
						"Descend brick.",
						"brick", candidate,
						"descend to", descendTo,
						"z diff", zDiff)
				}
				claimGreater("difference of >=1 is expected when descending bricks", zDiff, 0)
				originalStartZ := candidate.start.z
				originalEndZ := candidate.end.z
				candidate.start.z -= zDiff
				candidate.end.z -= zDiff
				claimGreater("candidate start z must be above ground", candidate.start.z, 0)
				bricks[brickIndex] = candidate
				moveBrickInMap(byLowZ, brickIndex, originalStartZ, candidate.start.z)
				moveBrickInMap(byHighZ, brickIndex, originalEndZ, candidate.end.z)
			}
			updateFloor(floor, candidate)
		}
	}
}

func createFloor(bricks []brick) [][]int {
	maximum := findMaximumXY(bricks)
	maximum.X++
	maximum.Y++
	floor := make([][]int, maximum.X)
	for i := range maximum.X {
		floor[i] = make([]int, maximum.Y)
	}
	return floor
}

func extractSortedZ(bricks []brick) []int {
	zCoords := gent.NewSet[int]()
	for _, each := range bricks {
		zCoords.Add(each.start.z)
	}
	zSorted := zCoords.ToSlice()
	slices.Sort(zSorted)
	return zSorted
}

func findMaximumXY(bricks []brick) shared.Loc {
	var maximum shared.Loc
	for _, each := range bricks {
		maximum.X = max(maximum.X, each.end.x)
		maximum.Y = max(maximum.Y, each.end.y)
	}
	return maximum
}

func claimGreater(message string, value, greaterThan int) {
	if value > greaterThan {
		return
	}
	panic(fmt.Sprintf("%s, %d is not >%d", message, value, greaterThan))
}

func canDescend(floor [][]int, br brick) (int, bool) {
	highest := math.MinInt
	for x := br.start.x; x <= br.end.x; x++ {
		for y := br.start.y; y <= br.end.y; y++ {
			peak := floor[x][y]
			if peak >= (br.start.z - 1) {
				return 0, false
			}
			highest = max(highest, peak)
		}
	}
	return highest + 1, true
}

func updateFloor(floor [][]int, br brick) {
	for x := br.start.x; x <= br.end.x; x++ {
		for y := br.start.y; y <= br.end.y; y++ {
			claimGreater("floor should always increase", br.end.z, floor[x][y])
			floor[x][y] = br.end.z
		}
	}
}

func moveBrickInMap(m map[int][]int, brickIndex, fromKey, toKey int) {
	if fromKey == toKey {
		panic("can't have the same keys")
	}
	indexes, ok := m[fromKey]
	if !ok {
		panic(
			fmt.Sprintf("brick index should always exist in map, from: %d, to: %d", fromKey, toKey),
		)
	}
	i := slices.Index(indexes, brickIndex)
	if i < 0 {
		panic("index should be always >=0")
	}
	fromIndexes := append(append([]int(nil), indexes[:i]...), indexes[i+1:]...)
	// Easier to test when result is deterministic so sort.
	slices.Sort(fromIndexes)
	if len(fromIndexes) > 0 {
		m[fromKey] = fromIndexes
	} else {
		delete(m, fromKey)
	}
	toIndexes := append([]int(nil), m[toKey]...)
	toIndexes = append(toIndexes, brickIndex)
	slices.Sort(toIndexes)
	m[toKey] = toIndexes
}

type greyhound struct {
	brickIndex int
	state      houndState
	stepCount  int
}

func newGreyhound(brickIndex int, state houndState) *greyhound {
	return &greyhound{
		brickIndex: brickIndex,
		state:      state,
	}
}

func (v *greyhound) String() string {
	return fmt.Sprintf(
		"{brickIndex:%d state:%s stepCount:%d}",
		v.brickIndex,
		v.state.String(),
		v.stepCount)
}

type runner struct {
	bricks      []brick
	ascendants  [][]int
	descendents [][]int
}

func (v *runner) moveHounds(hounds *[]*greyhound) bool {
	var houndMoved bool
	for _, each := range *hounds {
		if each.state == houndRunning {
			ascendants := v.ascendants[each.brickIndex]
			if len(ascendants) == 0 {
				each.state = houndSleeping
				continue
			}
			houndMoved = true
			for i, aboveIndex := range ascendants {
				if i > 0 {
					hound := newGreyhound(each.brickIndex, houndRunning)
					*hounds = append(*hounds, hound)
					each = hound
				}
				each.brickIndex = aboveIndex
				if len(v.descendents[aboveIndex]) == 1 {
					each.stepCount++
				} else {
					each.state = houndWaiting
				}
			}
		}
	}
	return houndMoved
}

func (v *runner) checkWaiting(hounds []*greyhound) bool {
	var waitChanged bool
	byBrickIndex := map[int][]int{}
	for i, hound := range hounds {
		if hound.state == houndWaiting {
			byBrickIndex[hound.brickIndex] = append(byBrickIndex[hound.brickIndex], i)
		}
	}
	for brickIndex, houndIndexes := range byBrickIndex {
		if len(houndIndexes) < 2 {
			continue
		}
		if len(houndIndexes) < len(v.descendents[brickIndex]) {
			continue
		}
		waitChanged = true
		topDog := hounds[houndIndexes[0]]
		topDog.state = houndRunning
		topDog.stepCount++
		for j, houndIndex := range houndIndexes {
			if j != 0 {
				hounds[houndIndex].state = houndSleeping
			}
		}
	}
	return waitChanged
}

func (v *runner) run(brickIndex int) int {
	hounds := []*greyhound{
		newGreyhound(brickIndex, houndRunning),
	}
	for {
		houndMoved := v.moveHounds(&hounds)
		waitChanged := v.checkWaiting(hounds)
		if !houndMoved && !waitChanged {
			// The hounds that started to wait but never stopped are those that stepped on to a
			// brick that's supported by bricks that are not reachable. Thus they are "discarded" at
			// this time.
			for _, hound := range hounds {
				if hound.state == houndWaiting {
					hound.state = houndSleeping
				}
			}
			break
		}
	}

	var total int
	for _, each := range hounds {
		if each.state != houndSleeping {
			shared.Logger.Error("Found a hound that's not sleeping.", "hound", each.String())
			panic("all hounds should be asleep")
		}
		total += each.stepCount
	}
	return total
}

func KillBricks(lines []string) int {
	bricks := parseLines(lines)
	ascendants := make([][]int, len(bricks))
	descendents := make([][]int, len(bricks))
	shared.Logger.Info(
		"Lines parsed, count destruction.",
		"brick count", len(bricks),
		"line count", len(lines))
	byLowZ, byHighZ := createSearchIndexes(bricks)
	descend(bricks, byLowZ, byHighZ)
	for i, each := range bricks {
		descendents[i] = findOverlaps(bricks, byHighZ[each.start.z-1], each)
		ascendants[i] = findOverlaps(bricks, byLowZ[each.end.z+1], each)
	}
	aRunner := &runner{
		bricks:      bricks,
		ascendants:  ascendants,
		descendents: descendents,
	}
	var total int
	for i := range bricks {
		count := aRunner.run(i)
		total += count
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Run finished.", "i", i, "count", count)
		}
	}
	shared.Logger.Info("Fallen bricks counted.", "count", total)
	return total
}
