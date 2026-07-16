package aoc2323

import (
	"fmt"
	"slices"

	"github.com/denarced/advent-of-code/shared"
)

const (
	charPath      = '.'
	charForest    = '#'
	charHillEast  = '>'
	charHillSouth = 'v'
	charHillWest  = '<'
	charHillNorth = '^'
)

func findOnRow(brd *shared.Board, y int, c rune) shared.Loc {
	for x := 0; x < brd.GetWidth(); x++ {
		loc := shared.Loc{X: x, Y: y}
		if brd.GetOrDie(loc) == c {
			return loc
		}
	}
	panic(fmt.Sprintf("failed to find \"%s\"", string(c)))
}

func FindLongestPath(lines []string) int {
	brd := shared.NewBoard(lines)
	start := findOnRow(brd, brd.GetHeight()-1, charPath)
	end := findOnRow(brd, 0, charPath)
	shared.Logger.Info("Derive longest path.", "start", start, "end", end)
	heads := []*shared.Link[shared.Loc]{shared.AddLink(nil, start)}
	var finishes []*shared.Link[shared.Loc]
	roundCount := 0
	var finished []int
	for len(heads) > 0 && roundCount < 1_000_000 {
		roundCount++
		for i, head := range heads {
			if head.Item == end {
				finished = append(finished, i)
				continue
			}
			nextSteps, _ := deriveNextSteps(brd, head, nil, false)
			for j := 1; j < len(nextSteps); j++ {
				heads = append(heads, shared.AddLink(head, nextSteps[j]))
			}
			heads[i] = shared.AddLink(head, nextSteps[0])
		}
		slices.Reverse(finished)
		for _, i := range finished {
			finishes = append(finishes, heads[i])
			heads = append(heads[:i], heads[i+1:]...)
		}
		finished = finished[:0]
	}
	if len(heads) > 0 {
		panic("exploration loop stopped due to safety limit, not because it finished properly")
	}
	shared.Logger.Info(
		"Exploration finished.",
		"path count", len(finishes),
		"round count", roundCount)
	maximum := -1
	for _, each := range finishes {
		count := countSteps(each)
		shared.Logger.Info("Candidate path counted.", "step count", count)
		maximum = max(maximum, count)
	}
	return maximum
}

func countSteps[T any](link *shared.Link[T]) int {
	var count int
	for l := link.Parent; l != nil; l = l.Parent {
		count++
	}
	return count
}

func deriveNextSteps(
	brd *shared.Board,
	link *shared.Link[shared.Loc],
	dirs []shared.Direction,
	climbUphill bool,
) (nextSteps []shared.Loc, count int) {
	if dirs == nil {
		dirs = shared.RealPrimaryDirections
	}
	for _, each := range dirs {
		moved := link.Item.Delta(shared.Loc(each))
		c, ok := brd.Get(moved)
		if !ok || c == charForest {
			continue
		}
		if c == charPath || climbUphill || validateHillMove(c, each) {
			nextSteps = append(nextSteps, moved)
		}
	}
	if link.Parent != nil {
		for i := range nextSteps {
			if link.Parent.Item == nextSteps[i] {
				nextSteps = append(nextSteps[:i], nextSteps[i+1:]...)
				break
			}
		}
	}
	count = len(nextSteps)
	kept := nextSteps[:0]
	for _, each := range nextSteps {
		keep := true
		for l := link; l != nil; l = l.Parent {
			if l.Item == each {
				keep = false
				break
			}
		}
		if keep {
			kept = append(kept, each)
		}
	}
	nextSteps = kept
	return
}

func validateHillMove(c rune, dir shared.Direction) bool {
	switch c {
	case charHillEast:
		return dir == shared.RealEast
	case charHillSouth:
		return dir == shared.RealSouth
	case charHillWest:
		return dir == shared.RealWest
	case charHillNorth:
		return dir == shared.RealNorth
	default:
		panic("unknown character for hill: " + string(c))
	}
}

type vertice struct {
	loc shared.Loc
}
type edge struct {
	fromIndex int
	toIndex   int
	length    int
	startDir  shared.Direction
}

type graph struct {
	edges    []edge
	vertices []vertice
}

func digFirst(link *shared.Link[shared.Loc]) shared.Loc {
	for l := link; l != nil; l = l.Parent {
		if l.Parent == nil {
			return l.Item
		}
	}
	panic("impossible: failed to dig first loc")
}

func makeEdge(aGraph graph, first, second shared.Loc, rat *shared.Link[shared.Loc]) edge {
	locs := []shared.Loc{first, second}
	sortLocs(locs)
	var dir shared.Direction
	if locs[0] == first {
		dir = extractDirection(rat, true)
	} else {
		dir = extractDirection(rat, false)
	}
	return edge{
		fromIndex: slices.IndexFunc(aGraph.vertices, func(aVertice vertice) bool {
			return aVertice.loc == locs[0]
		}),
		toIndex: slices.IndexFunc(aGraph.vertices, func(aVertice vertice) bool {
			return aVertice.loc == locs[1]
		}),
		length:   countSteps(rat),
		startDir: dir,
	}
}

func dive(
	aGraph graph,
	vToEdges [][]int,
	index int,
	used []int,
	edges []int,
	endIndex int,
	done func([]int),
) {
	if index == endIndex {
		done(edges)
		return
	}
	options := vToEdges[index]
	for _, edgeIndex := range options {
		dest := aGraph.edges[edgeIndex].toIndex
		if index == dest {
			dest = aGraph.edges[edgeIndex].fromIndex
		}
		// This is a lot faster than using an array to store used indexes or a map. With the former
		// and real puzzle input duration was ~3s. With gent.Set it was ~2.5s. This way it's ~0.8s.
		if used[dest] != 0 {
			continue
		}
		used[dest] = 1
		dive(aGraph, vToEdges, dest, used, append(edges, edgeIndex), endIndex, done)
		used[dest] = 0
	}
}

func FindLongestPathWithGraph(lines []string) int {
	aGraph := parseGraph(lines)
	vToEdges := make([][]int, len(aGraph.vertices))
	for i := range aGraph.vertices {
		for j, anEdge := range aGraph.edges {
			if anEdge.fromIndex == i || anEdge.toIndex == i {
				vToEdges[i] = append(vToEdges[i], j)
			}
		}
	}
	var maximum int
	done := func(edgePerm []int) {
		var total int
		for _, i := range edgePerm {
			total += aGraph.edges[i].length
		}
		maximum = max(maximum, total)
	}
	used := make([]int, len(aGraph.vertices))
	used[0] = 1
	dive(aGraph, vToEdges, 0, used, nil, 1, done)
	shared.Logger.Info("Max length derived.", "length", maximum)
	return maximum
}

func parseGraph(lines []string) (result graph) {
	brd := shared.NewBoard(lines)
	start := findOnRow(brd, brd.GetHeight()-1, charPath)
	end := findOnRow(brd, 0, charPath)
	result.vertices = []vertice{{loc: start}, {loc: end}}
	rats := []*shared.Link[shared.Loc]{shared.AddLink(nil, start)}
	dirs := []shared.Direction{shared.RealSouth}
	addEdge := func(firstRat *shared.Link[shared.Loc]) {
		anEdge := makeEdge(result, digFirst(firstRat), firstRat.Item, firstRat)
		if !slices.Contains(result.edges, anEdge) {
			shared.Logger.Debug(
				"Add edge.",
				"from", result.vertices[anEdge.fromIndex],
				"to", result.vertices[anEdge.toIndex],
				"length", anEdge.length)
			result.edges = append(result.edges, anEdge)
		}
	}
	for len(rats) > 0 {
		firstRat := rats[0]
		rats = rats[1:]
		startDir := dirs[0]
		dirs = dirs[1:]
		for {
			if firstRat.Parent != nil && firstRat.Item == end {
				addEdge(firstRat)
				break
			}
			var nextSteps []shared.Loc
			var nextStepOptionCount int
			if firstRat.Parent == nil {
				nextSteps, nextStepOptionCount = deriveNextSteps(
					brd,
					firstRat,
					[]shared.Direction{startDir},
					true)
			} else {
				nextSteps, nextStepOptionCount = deriveNextSteps(brd, firstRat, nil, true)
			}
			if len(nextSteps) == 0 {
				break
			}
			if len(nextSteps) == 1 && nextStepOptionCount == 1 {
				firstRat = shared.AddLink(firstRat, nextSteps[0])
				continue
			}
			if !slices.ContainsFunc(result.vertices, func(aVertice vertice) bool {
				return aVertice.loc == firstRat.Item
			}) {
				result.vertices = append(result.vertices, vertice{loc: firstRat.Item})
				for _, each := range nextSteps {
					rats = append(rats, shared.AddLink(nil, firstRat.Item))
					dirs = append(dirs, deriveDir(firstRat.Item, each))
				}
			}
			addEdge(firstRat)
			break
		}
	}
	shared.Logger.Info(
		"Graph parsed.",
		"edges", len(result.edges),
		"vertices", len(result.vertices))
	if shared.IsDebugEnabled() {
		for i, each := range result.edges {
			shared.Logger.Debug("Edge.", "i", i, "e", each)
		}
		for i, each := range result.vertices {
			shared.Logger.Debug("Vertice.", "i", i, "e", each)
		}
	}
	return
}

func sortLocs(locs []shared.Loc) {
	slices.SortFunc(locs, compareLocs)
}

func compareLocs(a, b shared.Loc) int {
	y := b.Y - a.Y
	if y != 0 {
		return y
	}
	return a.X - b.X
}

func deriveDir(a, b shared.Loc) shared.Direction {
	if a.X != b.X {
		if a.X < b.X {
			return shared.RealEast
		}
		return shared.RealWest
	}
	if a.Y < b.Y {
		return shared.RealNorth
	}
	return shared.RealSouth
}

func extractDirection(link *shared.Link[shared.Loc], first bool) shared.Direction {
	var prev, next shared.Loc
	for l := link; l != nil && l.Parent != nil; l = l.Parent {
		prev = l.Parent.Item
		next = l.Item
		if !first {
			break
		}
	}
	if first {
		return deriveDir(prev, next)
	}
	return deriveDir(next, prev)
}
