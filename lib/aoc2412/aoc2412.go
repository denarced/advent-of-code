package aoc2412

import (
	"fmt"
	"slices"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveTotalPrice(lines []string, discount bool) int {
	if len(lines) == 0 {
		return 0
	}
	brd := shared.NewBoard(lines)
	totalArea := brd.CountArea()
	shared.Logger.Info("Derive total price.", "discount", discount, "total area", totalArea)
	areas := findAreas(brd)
	totalPrice := 0
	for _, each := range areas {
		price := deriveAreaPrice(each, discount)
		shared.Logger.Info("Derived area price.", "price", price, "size", len(each), "area", each)
		totalPrice += price
	}
	return totalPrice
}

func findAreas(brd *shared.Board) [][]shared.Loc {
	totalArea := brd.CountArea()
	claimedLocs := map[shared.Loc]int{}
	areas := [][]shared.Loc{}
	for len(claimedLocs) < totalArea {
		plot := crawl(brd, findUnclaimed(claimedLocs, brd))
		areas = append(areas, plot)
		for _, loc := range plot {
			claimedLocs[loc] = 0
		}
	}
	return areas
}

func deriveAreaPrice(locs []shared.Loc, discount bool) int {
	if discount {
		return len(locs) * countFenceSides(locs)
	}

	mapped := map[shared.Loc]int{}
	for _, each := range locs {
		mapped[each] = 0
	}
	totalPerimeter := 0
	for _, each := range locs {
		count := countNeighbours(mapped, each)
		perimeter := shared.Abs(count - 4)
		totalPerimeter += perimeter
	}

	return totalPerimeter * len(locs)
}

func findUnclaimed(claimed map[shared.Loc]int, brd *shared.Board) shared.Loc {
	var unclaimed shared.Loc
	brd.Iter(func(loc shared.Loc, _ rune) bool {
		if _, ok := claimed[loc]; !ok {
			unclaimed = loc
			return false
		}
		return true
	})
	return unclaimed
}

func crawl(brd *shared.Board, loc shared.Loc) []shared.Loc {
	mapped := map[shared.Loc]int{loc: 0}
	{
		c, ok := brd.Get(loc)
		if !ok {
			panic(fmt.Sprintf("Crawl failed. Loc should never be invalid. Loc: %v.", loc))
		}
		last := []shared.Loc{loc}
		for len(last) > 0 {
			nextLast := []shared.Loc{}
			for _, latest := range last {
				for _, each := range brd.NextTo(latest, c) {
					if _, exist := mapped[each]; exist {
						continue
					}
					mapped[each] = 0
					nextLast = append(nextLast, each)
				}
			}
			last = nextLast
		}
	}

	crawled := []shared.Loc{}
	for loc := range mapped {
		crawled = append(crawled, loc)
	}
	return crawled
}

func countNeighbours(area map[shared.Loc]int, loc shared.Loc) int {
	count := 0
	for _, each := range []shared.Loc{loc.Delta(shared.Loc{X: 1, Y: 0}),
		loc.Delta(shared.Loc{X: 0, Y: -1}),
		loc.Delta(shared.Loc{X: -1, Y: 0}),
		loc.Delta(shared.Loc{X: 0, Y: 1})} {
		if _, exist := area[each]; exist {
			count++
		}
	}
	return count
}

type fence = shared.Pair[shared.Loc]

type fatFence struct {
	f   fence
	loc shared.Loc
}

func (v fatFence) ToString() string {
	return fmt.Sprintf("%s->%s (%s)", v.f.First.ToString(), v.f.Second.ToString(), v.loc.ToString())
}

func newFence(first, second shared.Loc) fence {
	// TODO Prevent illegal pairs.
	return shared.NewPair(first, second)
}

func newFatFence(base, first, second shared.Loc) fatFence {
	return fatFence{
		f:   newFence(first, second),
		loc: base,
	}
}

func countFenceSides(contiguous []shared.Loc) int {
	if len(contiguous) == 0 {
		return 0
	}
	fatFamilies := sortFences(deriveFatties(contiguous))
	sides := 0
	for _, sorted := range fatFamilies {
		firstFence := sorted[0]
		shared.Logger.Debug("First fence.", "fence", firstFence)
		currentDirection := deriveDirectioner(firstFence)
		// shared.Logger.Debug("Direction.", "dir", currentDirection)
		for _, each := range sorted[1:] {
			nextDirection := deriveDirectioner(each)
			if currentDirection != nextDirection {
				sides++
			}
			currentDirection = nextDirection
		}
		if deriveDirectioner(firstFence) != currentDirection {
			sides++
		}
	}
	return sides
}

func deriveDirectioner(fattie fatFence) int {
	if fattie.f.First.X == fattie.f.Second.X {
		// North or south.
		return 1
	}
	// West or east.
	return 2
}

func deriveFatties(contiguous []shared.Loc) []fatFence {
	shared.Logger.Info("Derive fat fences.", "location count", len(contiguous))
	mapped := map[fence][]shared.Loc{}
	for _, block := range contiguous {
		for _, each := range createSurroundingFences(block) {
			norm := clockwiseFence(each)
			mapped[norm] = append(mapped[norm], block)
		}
	}

	fatties := []fatFence{}
	for aFence, locs := range mapped {
		if len(locs) == 1 {
			fatties = append(fatties, fatFence{
				f:   aFence,
				loc: locs[0],
			})
		}
	}

	stringified := []string{}
	for _, each := range fatties {
		stringified = append(stringified, each.ToString())
	}
	shared.Logger.Debug("Derived fat fences.", "count", len(fatties), "fat fences", stringified)
	return fatties
}

func createSurroundingFences(loc shared.Loc) []fence {
	nw := shared.Loc{X: loc.X, Y: loc.Y + 1}
	ne := shared.Loc{X: loc.X + 1, Y: loc.Y + 1}
	se := shared.Loc{X: loc.X + 1, Y: loc.Y}
	// Clockwise.
	return []fence{
		newFence(ne, se),
		newFence(se, loc),
		newFence(loc, nw),
		newFence(nw, ne),
	}
}

func clockwiseFence(f fence) fence {
	dx := f.Second.X - f.First.X
	dy := f.Second.Y - f.First.Y
	if dx == 0 {
		if dy > 0 {
			return newFence(f.Second, f.First)
		}
		return f
	}
	if dx > 0 {
		return newFence(f.Second, f.First)
	}
	return f
}

func sortFences(unsorted []fatFence) [][]fatFence {
	fatFenceLines := [][]fatFence{}
	if len(unsorted) == 0 {
		return fatFenceLines
	}
	shared.Logger.Debug("Sort fences.", "count", len(unsorted), "unsorted", unsorted)
	for len(unsorted) > 0 {
		first := unsorted[0]
		unsorted = slices.Delete(unsorted, 0, 1)
		shared.Logger.Debug("Start next fence set.", "first", first.ToString())
		previous := first
		nestedFatties := []fatFence{first}
		for {
			current, i := findNextFence(unsorted, previous)
			if i < 0 {
				msg := "Failed to find next"
				shared.Logger.Error(msg, "previous", previous)
				panic(msg)
			}
			shared.Logger.Debug("Next in fence set found.", "index", i, "next", current.ToString())
			nestedFatties = append(nestedFatties, current)
			unsorted = slices.Delete(unsorted, i, i+1)
			if current.f.Second == first.f.First {
				break
			}
			previous = current
		}
		shared.Logger.Info("Nested fatties.", "fatties", nestedFatties)
		fatFenceLines = append(fatFenceLines, nestedFatties)
	}
	return fatFenceLines
}

type indexedFatFence struct {
	index int
	f     fatFence
}

func findNextFence(fatties []fatFence, previous fatFence) (fatFence, int) {
	candidates := []indexedFatFence{}
	end := previous.f.Second
	for i, each := range fatties {
		if end == each.f.First && each.f.Second != previous.f.First {
			candidates = append(candidates, indexedFatFence{index: i, f: each})
		}
		if end == each.f.Second && each.f.First != previous.f.First {
			candidates = append(
				candidates,
				indexedFatFence{f: newFatFence(each.loc, each.f.Second, each.f.First), index: i},
			)
		}
	}
	if len(candidates) == 0 {
		return fatFence{}, -1
	}
	if len(candidates) == 1 {
		last := candidates[0]
		return last.f, last.index
	}
	sharedLocs := []indexedFatFence{}
	for _, each := range candidates {
		if each.f.loc == previous.loc {
			sharedLocs = append(sharedLocs, each)
		}
	}
	if len(sharedLocs) == 0 {
		shared.Logger.Error(
			"Can't find a fat fence that shared location with previous.",
			"previous",
			previous,
			"candidates",
			candidates,
		)
		panic("No shared locs..")
	}
	if len(sharedLocs) > 1 {
		shared.Logger.Error(
			"Multiple fat fences to choose from. Can't decide.",
			"previous",
			previous,
			"candidates with shared location",
			sharedLocs,
		)
		panic("Still multiple fat fences to choose from.")
	}
	last := sharedLocs[0]
	return last.f, last.index
}
