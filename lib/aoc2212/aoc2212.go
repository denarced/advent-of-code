package aoc2212

import (
	"github.com/denarced/advent-of-code/shared"
)

func DeriveFewestSteps(lines []string) int {
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	end := brd.FindOrDie('E')
	brd.Set(start, 'a')
	brd.Set(end, 'z')
	nodes := map[shared.Loc]*node{}
	nodes[start] = &node{distance: 0, loc: start}
	shared.Logger.Info(
		"Derive fewest steps.",
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
	if to <= from {
		return false
	}
	return int(to-from) > 1
}
