package aoc2321

import (
	"strconv"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	plotChar = '.'
	rockChar = '#'
)

func CountRangeFromLines(lines []string, stepCount int, infinite bool) int {
	shared.Logger.Info(
		"Count range.",
		"line count", len(lines),
		"step count", stepCount,
		"infinite", infinite)
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	brd.Set(start, plotChar)
	walkers := walk(
		map[shared.Loc]*gent.Set[shared.Loc]{start: gent.NewSet(shared.Loc{})},
		stepCount,
		infinite,
		brd)
	var result int
	for _, sections := range walkers {
		result += sections.Count()
	}
	shared.Logger.Info("Range counted.", "result", result)
	return result
}

func walk(
	walkers map[shared.Loc]*gent.Set[shared.Loc],
	stepCount int,
	infinite bool,
	brd *shared.Board,
) map[shared.Loc]*gent.Set[shared.Loc] {
	for range stepCount {
		next := deriveNextWalkers(walkers, brd, infinite)
		if len(next) == 0 {
			break
		}
		walkers = next
	}
	return walkers
}

func wrap(
	brd *shared.Board,
	infinite bool,
	loc shared.Loc,
	sections *gent.Set[shared.Loc],
) (shared.Loc, *gent.Set[shared.Loc]) {
	if !infinite {
		return loc, sections
	}
	var shiftX, shiftY int
	if loc.X < 0 {
		claim(loc.X == -1, "x != -1")
		loc.X += brd.GetWidth()
		shiftX--
	} else if loc.X >= brd.GetWidth() {
		claim(loc.X == brd.GetWidth(), "x != width")
		loc.X -= brd.GetWidth()
		shiftX++
	} else if loc.Y < 0 {
		claim(loc.Y == -1, "y != -1")
		loc.Y += brd.GetHeight()
		shiftY--
	} else if loc.Y >= brd.GetHeight() {
		claim(loc.Y == brd.GetHeight(), "y != height")
		loc.Y -= brd.GetHeight()
		shiftY++
	}
	if shiftX != 0 || shiftY != 0 {
		shifted := gent.NewSet[shared.Loc]()
		sections.ForEachAll(func(loc shared.Loc) {
			adjacent := loc.Delta(shared.Loc{X: shiftX, Y: shiftY})
			shifted.Add(adjacent)
		})
		sections = shifted
	}
	return loc, sections
}

func claim(val bool, message string) {
	if !val {
		panic(message)
	}
}

func deriveNextWalkers(
	walkers map[shared.Loc]*gent.Set[shared.Loc],
	brd *shared.Board,
	infinite bool,
) map[shared.Loc]*gent.Set[shared.Loc] {
	next := map[shared.Loc]*gent.Set[shared.Loc]{}
	for loc, sections := range walkers {
		for _, dir := range shared.RealPrimaryDirections {
			adjacent, shifted := wrap(
				brd,
				infinite,
				loc.Delta(shared.Loc(dir)),
				sections)
			cell, ok := brd.Get(adjacent)
			if !ok || cell != plotChar {
				continue
			}
			existing := next[adjacent]
			if existing == nil {
				existing = gent.NewSet[shared.Loc]()
				next[adjacent] = existing
			}
			shifted.ForEachAll(func(section shared.Loc) {
				existing.Add(section)
			})
		}
	}
	return next
}

type deriveSectionFunc func(shared.Loc) shared.Loc

func createDeriveSection(size int) deriveSectionFunc {
	alignNegative := func(i int) int {
		if i >= 0 {
			return i
		}
		return i + 1 - size
	}
	return func(loc shared.Loc) shared.Loc {
		loc.X = alignNegative(loc.X)
		loc.Y = alignNegative(loc.Y)
		return shared.Loc{X: loc.X / size, Y: loc.Y / size}
	}
}

type diamond struct {
	deriveSection deriveSectionFunc
	radius        int
	size          int
	start         shared.Loc
}

func newDiamond(start shared.Loc, stepCount, size int) *diamond {
	point := start.Delta(shared.Loc{X: stepCount})
	deriveSection := createDeriveSection(size)
	pointSection := deriveSection(point)
	radius := max(0, pointSection.X-2)
	d := diamond{
		deriveSection: deriveSection,
		radius:        radius,
		size:          size,
		start:         start,
	}
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("New diamond.", "d", d, "section", pointSection, "point", point)
	}
	return &d
}

func (v *diamond) countTotal(even, odd int) int {
	if v.radius < 1 {
		return 0
	}
	// Given radius
	//     1: 1 even, 4 odd
	//     2: 9 even, 4 odd
	//     3: 9 even, 16 odd
	// So given radius r, section count is r^2 + (r+1)^2. Even section count is always odd. Here
	// "even section" means that sum of its absolute coordinates is even. E.g. sum of 0,0 is even.
	count := func(i int) int {
		i *= i
		if i%2 == 0 {
			return i * odd
		}
		return i * even
	}
	return count(v.radius) + count(v.radius+1)
}

func (v *diamond) isForbidden(loc shared.Loc) bool {
	section := v.deriveSection(loc)
	return shared.Abs(section.X)+shared.Abs(section.Y) <= v.radius
}

func (v *diamond) stepOutside() ([]shared.Loc, int) {
	stepCount := v.size/2 + 1 + v.radius*v.size
	starts := make([]shared.Loc, 0, 4*(v.radius+1))
	starts = append(starts, v.start.Delta(shared.Loc{Y: stepCount}))

	hor := v.size/2 + 1
	ver := stepCount - hor
	for ver > -stepCount {
		starts = append(starts, v.start.Delta(shared.Loc{X: hor, Y: ver}))
		ver -= v.size
		hor = stepCount - shared.Abs(ver)
	}

	starts = append(starts, v.start.Delta(shared.Loc{Y: -stepCount}))

	hor = v.size/2 + 1
	ver = -stepCount + hor
	for ver < stepCount {
		starts = append(starts, v.start.Delta(shared.Loc{X: -hor, Y: ver}))
		ver += v.size
		hor = stepCount - shared.Abs(ver)
	}
	return starts, stepCount
}

func CountInfiniteRange(lines []string, stepCount int) int {
	shared.Logger.Info("Count infinite range.", "step count", stepCount)
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	brd.Set(start, plotChar)

	aDiamond := newDiamond(start, stepCount, brd.GetWidth())
	var calculatedTotal int
	if aDiamond.radius > 0 {
		calculatedTotal = calculateWithinDiamondRadius(brd, stepCount, start, aDiamond)
	}

	starts, stepsToStarts := aDiamond.stepOutside()
	startSets := make([][]shared.Loc, len(starts))
	for i, each := range starts {
		startSets[i] = []shared.Loc{each}
	}
	remainingSteps := stepCount - stepsToStarts
	if remainingSteps%2 != 0 {
		startSets = expandStartSets(brd, start, starts, startSets)
		remainingSteps--
	}
	var discoveredTotal int
	resolveUniqueCount := func(i int) int {
		current := discoverPlots(brd, startSets[i], remainingSteps, aDiamond.isForbidden)
		next := discoverPlots(
			brd,
			startSets[(i+1)%len(startSets)],
			remainingSteps,
			aDiamond.isForbidden,
		)
		uniqueCount := countUnique(current, next)
		discoveredTotal += uniqueCount
		return uniqueCount
	}
	feed := newQuarterFeed(len(startSets))
	for indexFunc := feed.create(); indexFunc != nil; indexFunc = feed.create() {
		var dawn, dusk []int
		var convergedSize int
		quarterStart, quarterEnd := indexFunc()
		for quarterStart >= 0 {
			if convergedSize > 0 {
				if quarterStart < quarterEnd {
					discoveredTotal += (quarterEnd - quarterStart + 1) * convergedSize
					break
				}
				discoveredTotal += convergedSize
				if quarterEnd >= 0 {
					discoveredTotal += convergedSize
				}
			} else {
				dawn = append(dawn, resolveUniqueCount(quarterStart))
				if quarterEnd >= 0 {
					dusk = append(dusk, resolveUniqueCount(quarterEnd))
					if found, ok := isConverging(dawn, dusk); ok {
						convergedSize = found
					}
				}
			}
			quarterStart, quarterEnd = indexFunc()
		}
	}
	total := discoveredTotal + calculatedTotal
	shared.Logger.Info(
		"Range count done.",
		"total", total,
		"calculated", calculatedTotal,
		"discovered", discoveredTotal)
	return total
}

func invadeBoard(brd *shared.Board, start shared.Loc, even bool, quitFunc func(int) []int) int {
	section := createDeriveSection(brd.GetWidth())(start)
	minLoc, maxLoc := deriveBoundsForSection(brd.GetWidth(), section)
	soldiers := []shared.Loc{start}
	for i := 0; ; i++ {
		if sizes := quitFunc(len(soldiers)); sizes != nil {
			var soldierCount int
			if even && i%2 == 0 || !even && i%2 != 0 {
				soldierCount = sizes[1]
			} else {
				soldierCount = sizes[0]
			}
			shared.Logger.Info(
				"Board invaded.",
				"start", start,
				"even", even,
				"sizes", sizes,
				"soldier count", soldierCount)
			return soldierCount
		}
		soldiers = deriveNext(soldiers, brd, minLoc, maxLoc)
	}
}

func deriveNext(
	soldiers []shared.Loc,
	brd *shared.Board,
	minLoc, maxLoc shared.Loc,
) []shared.Loc {
	next := gent.NewSet[shared.Loc]()
	for _, each := range soldiers {
		for _, dir := range shared.RealPrimaryDirections {
			adjacent := each.Delta(shared.Loc(dir))
			if adjacent.X < minLoc.X || adjacent.Y < minLoc.Y || maxLoc.X < adjacent.X ||
				maxLoc.Y < adjacent.Y {
				continue
			}
			if brd.GetRelative(adjacent) != plotChar {
				continue
			}
			next.Add(adjacent)
		}
	}
	return next.ToSlice()
}

func deriveBoundsForSection(size int, section shared.Loc) (minimum, maximum shared.Loc) {
	minimum.X = section.X * size
	minimum.Y = section.Y * size
	maximum.X = minimum.X + size - 1
	maximum.Y = minimum.Y + size - 1
	return
}

func createRepeatMonitor() func(int) []int {
	bufferSize := 10
	var bufferIndex int
	buffer := make([]int, bufferSize)
	var pushCount int
	translateIndex := func(i int) int {
		return i % bufferSize
	}
	return func(size int) []int {
		if size <= 0 {
			panic("impossible size argument: " + strconv.Itoa(size))
		}
		if pushCount > 1000 {
			shared.Logger.Error("Probably infinite loop ongoing.", "buffer", buffer)
			panic("infinite loop")
		}
		buffer[translateIndex(bufferIndex)] = size
		bufferIndex++
		pushCount++
		if pushCount < bufferSize {
			return nil
		}
		if buffer[translateIndex(bufferIndex-1)] != buffer[translateIndex(bufferIndex-3)] {
			return nil
		}
		if buffer[translateIndex(bufferIndex-1)] != buffer[translateIndex(bufferIndex-9)] {
			return nil
		}
		expectedPair := make([]int, 2)
		for i := bufferIndex - 1; i >= bufferIndex-10; i -= 2 {
			pair := []int{buffer[translateIndex(i-1)], buffer[translateIndex(i)]}
			if expectedPair[0] == 0 {
				copy(expectedPair, pair)
				continue
			}
			for j := range expectedPair {
				if expectedPair[j] != pair[j] {
					return nil
				}
			}
		}
		return expectedPair
	}
}

func directionsExcept(dir shared.Direction) []shared.Direction {
	remaining := make([]shared.Direction, 3)
	var i int
	for _, each := range shared.RealPrimaryDirections {
		if each == dir {
			continue
		}
		remaining[i] = each
		i++
	}
	return remaining
}

func stepToDirections(from shared.Loc, directions ...shared.Direction) []shared.Loc {
	if len(directions) == 0 {
		return nil
	}
	result := make([]shared.Loc, len(directions))
	var i int
	for _, each := range directions {
		result[i] = from.Delta(shared.Loc(each))
		i++
	}
	return result
}

func discoverPlots(
	brd *shared.Board,
	soldiers []shared.Loc,
	stepCount int,
	isForbidden func(shared.Loc) bool,
) *gent.Set[shared.Loc] {
	current := gent.NewSet(soldiers...)
	for range stepCount {
		current = discoverNext(current, brd, isForbidden)
	}
	return current
}

func discoverNext(
	soldiers *gent.Set[shared.Loc],
	brd *shared.Board,
	isForbidden func(shared.Loc) bool,
) *gent.Set[shared.Loc] {
	next := gent.NewSet[shared.Loc]()
	soldiers.ForEachAll(func(each shared.Loc) {
		for _, dir := range shared.RealPrimaryDirections {
			adjacent := each.Delta(shared.Loc(dir))
			if brd.GetRelative(adjacent) != plotChar {
				continue
			}
			if isForbidden(adjacent) {
				continue
			}
			next.Add(adjacent)
		}
	})
	return next
}

func countUnique(a, b *gent.Set[shared.Loc]) int {
	var count int
	a.ForEachAll(func(each shared.Loc) {
		if !b.Contains(each) {
			count++
		}
	})
	return count
}

type quarterFeed interface {
	create() func() (int, int)
}

type properQuarterFeed struct {
	period  int
	size    int
	current int
}

func newQuarterFeed(size int) quarterFeed {
	if size < 4 {
		panic("size must be >=4 in newQuarterFeed")
	}
	var period int
	if size%4 == 0 {
		period = size / 4
	} else {
		period = size/4 + 1
	}
	if size < 4*12 {
		return &singleQuarterFeed{size: size, period: period}
	}
	return &properQuarterFeed{size: size, period: period}
}

func (v *properQuarterFeed) create() func() (int, int) {
	low := v.current * v.period
	high := min(v.size-1, (v.current+1)*v.period-1)
	if high < low {
		return nil
	}
	v.current++
	i := 0
	return func() (int, int) {
		head := low + i
		tail := high - i
		i++
		if tail < head {
			return -1, -1
		}
		if head == tail {
			tail = -1
		}
		return head, tail
	}
}

type singleQuarterFeed struct {
	period  int
	size    int
	current int
}

func (v *singleQuarterFeed) create() func() (int, int) {
	low := v.current * v.period
	high := min(v.size-1, (v.current+1)*v.period-1)
	if high < low {
		return nil
	}
	i := low
	v.current++
	return func() (int, int) {
		if high < i {
			return -1, -1
		}
		returned := i
		i++
		return returned, -1
	}
}

func isConverging(dawn, dusk []int) (int, bool) {
	if len(dawn) < 4 || len(dusk) < 4 {
		return 0, false
	}
	expected := dawn[len(dawn)-1]
	for _, each := range [][]int{dawn, dusk} {
		for i := range 2 {
			if each[len(each)-1-i] != expected {
				return 0, false
			}
		}
	}
	shared.Logger.Info(
		"Found converged value.",
		"dawn size", len(dawn),
		"dusk size", len(dusk),
		"value", expected)
	return expected, true
}

func calculateWithinDiamondRadius(
	brd *shared.Board,
	stepCount int,
	start shared.Loc,
	aDiamond *diamond,
) int {
	even := stepCount%2 == 0
	evenBoardCount := invadeBoard(brd, start, even, createRepeatMonitor())
	stepsToAdjacent := brd.GetWidth()/2 + 1
	adjacentBoardStart := start.Delta(shared.Loc{X: stepsToAdjacent})

	even = (stepCount-stepsToAdjacent)%2 == 0
	oddBoardCount := invadeBoard(brd, adjacentBoardStart, even, createRepeatMonitor())
	return aDiamond.countTotal(evenBoardCount, oddBoardCount)
}

func expandStartSets(
	brd *shared.Board,
	start shared.Loc,
	starts []shared.Loc,
	startSets [][]shared.Loc,
) [][]shared.Loc {
	expandedStartSets := make([][]shared.Loc, 0, len(startSets))
	for _, aStart := range starts {
		expanded := make([]shared.Loc, 0, 3)
		var excluded shared.Direction
		if start.X < aStart.X {
			excluded = shared.RealWest
		} else if aStart.X < start.X {
			excluded = shared.RealEast
		} else {
			if start.Y < aStart.Y {
				excluded = shared.RealSouth
			} else {
				excluded = shared.RealNorth
			}
		}
		for _, loc := range stepToDirections(aStart, directionsExcept(excluded)...) {
			if brd.GetRelative(loc) == plotChar {
				expanded = append(expanded, loc)
			}
		}
		expandedStartSets = append(expandedStartSets, expanded)
	}
	shared.Logger.Info(
		"Expanded starts added.",
		"start set count", len(startSets),
		"expanded set count", len(expandedStartSets))
	return expandedStartSets
}
