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
			nextSteps := deriveNextSteps(brd, head)
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
		// The first location is ignored because we're counting steps, not locations.
		count := -1
		for l := each; l != nil; l = l.Parent {
			count++
		}
		shared.Logger.Info("Candidate path counted.", "step count", count)
		maximum = max(maximum, count)
	}
	return maximum
}

func deriveNextSteps(brd *shared.Board, link *shared.Link[shared.Loc]) (nextSteps []shared.Loc) {
	for _, each := range shared.RealPrimaryDirections {
		moved := link.Item.Delta(shared.Loc(each))
		c, ok := brd.Get(moved)
		if !ok || c == charForest {
			continue
		}
		if c == charPath || validateHillMove(c, each) {
			nextSteps = append(nextSteps, moved)
		}
	}
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
