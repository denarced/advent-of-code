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

func CountBricksFromLines(lines []string) int {
	bricks := parseLines(lines)
	shared.Logger.Info(
		"Lines parsed, count bricks that can be disintegrated.",
		"brick count",
		len(bricks),
		"line count",
		len(lines),
	)
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

func findSlackers(bricks []brick, byLow, byHigh map[int][]int) *gent.Set[brick] {
	slackers := gent.NewSet[brick]()
	for _, eachHighZ := range sortKeys(byHigh) {
		belowBrickIndexes := byHigh[eachHighZ]
		aboveBrickIndexes := byLow[eachHighZ+1]
		bricksAbove := len(aboveBrickIndexes) > 0
		if shared.IsDebugEnabled() {
			if !bricksAbove {
				shared.Logger.Debug(
					"No bricks above, all slackers with high Z.",
					"z", eachHighZ,
					"count", len(belowBrickIndexes))
			}
		}
		aboveToBelow := map[int][]int{}
		singleParents := gent.NewSet[brick]()
		for _, belowBrickIndex := range belowBrickIndexes {
			// This brick is the focus here: does it support nothing or is it one of many for bricks
			// above that it supports.
			lowBrick := bricks[belowBrickIndex]
			if !bricksAbove {
				if slackers.Add(lowBrick) {
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
				if slackers.Add(lowBrick) {
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
	floor, zCoords := createFloor(bricks)
	for _, z := range zCoords {
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

func createFloor(bricks []brick) ([][]int, []int) {
	zCoords := gent.NewSet[int]()
	for _, each := range bricks {
		zCoords.Add(each.start.z)
	}
	zSorted := zCoords.ToSlice()
	slices.Sort(zSorted)

	maximum := findMaximumXY(bricks)
	maximum.X++
	maximum.Y++
	floor := make([][]int, maximum.X)
	for i := range maximum.X {
		floor[i] = make([]int, maximum.Y)
	}
	return floor, zSorted
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
