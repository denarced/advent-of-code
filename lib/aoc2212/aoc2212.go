package aoc2212

import (
	"github.com/denarced/advent-of-code/shared"
)

const (
	StartToEnd HikeStrategy = iota
	LowToEnd
)

type HikeStrategy int

func DeriveFewestSteps(lines []string, hikeDir HikeStrategy) int {
	switch hikeDir {
	case StartToEnd:
		return deriveFewestStepsFromStartToEnd(lines)
	case LowToEnd:
		return deriveFewestStepsFromLowToEnd(lines)
	default:
		panic("unknown HikeStrategy")
	}
}

func deriveFewestStepsFromStartToEnd(lines []string) int {
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	end := brd.FindOrDie('E')
	brd.Set(start, 'a')
	brd.Set(end, 'z')
	nodes := map[shared.Loc]*node{}
	nodes[start] = &node{distance: 0, loc: start}
	shared.Logger.Info(
		"Derive fewest steps from start to end.",
		"start", start,
		"width", brd.GetWidth(),
		"height", brd.GetHeight())
	brd.Iter(func(loc shared.Loc, _ rune) bool {
		if loc != start {
			nodes[loc] = &node{distance: -1, loc: loc}
		}
		return true
	})

	for {
		nod := pickNode(nodes)
		if nod == nil {
			break
		}
		fromElevation := brd.GetOrDie(nod.loc)
		for _, dir := range shared.RealPrimaryDirections {
			loc := nod.loc.Delta(shared.Loc(dir))
			toElevation, ok := brd.Get(loc)
			if !ok || isTooHigh(fromElevation, toElevation) {
				continue
			}
			nearLoc := nodes[loc]
			if nearLoc.visited {
				continue
			}
			if nearLoc.distance < 0 || nod.distance < nearLoc.distance {
				nearLoc.distance = nod.distance + 1
			}
		}
		nod.visited = true
	}

	return nodes[end].distance
}

type node struct {
	distance int
	visited  bool
	loc      shared.Loc
}

func pickNode(nodes map[shared.Loc]*node) *node {
	var candidate *node
	for _, each := range nodes {
		if each.visited || each.distance < 0 {
			continue
		}
		if candidate == nil || each.distance < candidate.distance {
			candidate = each
		}
	}
	return candidate
}

func isTooHigh(from, to rune) bool {
	if from < 'a' || 'z' < from {
		panic("from is evil")
	}
	if to < 'a' || 'z' < to {
		panic("to is evil")
	}
	return int(to-from) > 1
}

func deriveFewestStepsFromLowToEnd(lines []string) int {
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('E')
	brd.Set(start, 'z')
	brd.Set(brd.FindOrDie('S'), 'a')
	nodes := map[shared.Loc]*node{}
	nodes[start] = &node{distance: 0, loc: start}
	shared.Logger.Info(
		"Derive fewest steps from low to end.",
		"start", start,
		"width", brd.GetWidth(),
		"height", brd.GetHeight())
	brd.Iter(func(loc shared.Loc, _ rune) bool {
		if loc != start {
			nodes[loc] = &node{distance: -1, loc: loc}
		}
		return true
	})

	for {
		nod := pickNode(nodes)
		if nod == nil {
			break
		}
		fromElevation := brd.GetOrDie(nod.loc)
		if fromElevation == 'a' {
			nod.visited = true
			continue
		}
		for _, dir := range shared.RealPrimaryDirections {
			loc := nod.loc.Delta(shared.Loc(dir))
			toElevation, ok := brd.Get(loc)
			if !ok || isTooHigh(toElevation, fromElevation) {
				continue
			}
			nearLoc := nodes[loc]
			if nearLoc.visited {
				continue
			}
			if nearLoc.distance < 0 || nod.distance < nearLoc.distance {
				nearLoc.distance = nod.distance + 1
			}
		}
		nod.visited = true
	}

	shared.Logger.Info("Review all 'a' nodes to find smallest distance.")
	smallestDistance := -1
	for loc, each := range nodes {
		c, ok := brd.Get(loc)
		if !ok || c != 'a' || each.distance < 0 {
			continue
		}
		if smallestDistance < 0 || each.distance < smallestDistance {
			shared.Logger.Info(
				"New smallest distance found.",
				"loc", each.loc,
				"distance", each.distance)
			smallestDistance = each.distance
		}
	}
	shared.Logger.Info("Smallest distance found.", "distance", smallestDistance)
	return smallestDistance
}
