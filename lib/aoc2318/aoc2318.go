package aoc2318

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func Dig(lines []string) int {
	shared.Logger.Info("Start digging.", "line count", len(lines))
	instructions := parseLines(lines)
	clockwise := isClockwise(instructions)
	var loc shared.Loc
	var path []shared.Loc
	// Sanity check. Ultimately movement sum should be zero. That is, U+D == 0 and L+R == 0.
	var balance shared.Loc
	for _, instr := range instructions {
		next := loc.Delta(instr.delta)
		path = append(path, between(loc, next)...)
		loc = next
		balance = balance.Delta(instr.delta)
	}
	if balance != (shared.Loc{}) {
		panic("out of balance")
	}
	stash := gent.NewSet(path...)
	volume := countSpaceInBetween(path, clockwise, stash)
	shared.Logger.Info("Digging done.", "volume", volume)
	return volume
}

func isClockwise(instructions []instruction) bool {
	var right int
	for i := range instructions {
		right += deriveTurn(instructions, i)
	}
	return right > 0
}

func parseLines(lines []string) []instruction {
	if len(lines) == 0 {
		return nil
	}
	instructions := make([]instruction, len(lines))
	for i, each := range lines {
		instructions[i] = parseLine(each)
	}
	return instructions
}

type instruction struct {
	dir       shared.Direction
	stepCount int
	delta     shared.Loc
}

func newInstruction(dir shared.Direction, stepCount int) instruction {
	delta := shared.Loc(dir)
	if delta.X != 0 {
		delta.X *= stepCount
	} else {
		delta.Y *= stepCount
	}
	return instruction{
		dir:       dir,
		stepCount: stepCount,
		delta:     delta,
	}
}

func parseLine(line string) instruction {
	pieces := strings.Fields(line)
	if len(pieces) != 3 {
		shared.Logger.Error("Invalid line, piece count != 3.", "count", len(pieces), line, "line")
		panic("invalid line, piece count != 3")
	}
	dir := toDirection(pieces[0])
	count := gent.OrPanic2(strconv.Atoi(pieces[1]))("invalid count: " + pieces[1])
	return newInstruction(dir, count)
}

func toDirection(s string) shared.Direction {
	switch strings.TrimSpace(s) {
	case "R":
		return shared.RealEast
	case "D":
		return shared.RealSouth
	case "L":
		return shared.RealWest
	case "U":
		return shared.RealNorth
	default:
		panic(fmt.Sprintf("invalid direction: %s", s))
	}
}

func toIndex(i, length int) int {
	i %= length
	if i >= 0 {
		return i
	}
	return i + length
}

func deriveTurn(instructions []instruction, i int) int {
	curr := instructions[toIndex(i, len(instructions))]
	next := instructions[toIndex(i+1, len(instructions))]
	if curr.dir == next.dir || curr.dir.X+next.dir.X == 0 && curr.dir.Y+next.dir.Y == 0 {
		return 0
	}
	pick := func(dir shared.Direction) int {
		if dir.X != 0 {
			return dir.X
		}
		return dir.Y
	}
	if curr.dir.X == 0 {
		if pick(curr.dir) == pick(next.dir) {
			return 1
		}
		return -1
	}
	if pick(curr.dir) != pick(next.dir) {
		return 1
	}
	return -1
}

func between(from, to shared.Loc) []shared.Loc {
	var stepCount int
	deriveUnit := func(start, end int) int {
		diff := end - start
		stepCount = max(stepCount, shared.Abs(diff))
		if diff == 0 {
			return diff
		}
		if diff < 0 {
			return -1
		}
		return 1
	}
	step := shared.Loc{X: deriveUnit(from.X, to.X), Y: deriveUnit(from.Y, to.Y)}
	steps := make([]shared.Loc, 0, stepCount)
	for loc := from; loc != to; loc = loc.Delta(step) {
		steps = append(steps, loc)
	}
	return steps
}

func countSpaceInBetween(path []shared.Loc, clockwise bool, stash *gent.Set[shared.Loc]) int {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Count space.", "clockwise", clockwise)
	}
	for i := range path {
		curr := path[i]
		nextIndex := toIndex(i+1, len(path))
		next := path[nextIndex]
		side := deriveSide(curr, next, clockwise)
		if curr == next {
			panic("broken: sequential locations are identical")
		}
		spawn(stash, side)
	}
	return stash.Len()
}

func deriveSide(first, second shared.Loc, clockwise bool) shared.Loc {
	if first.X == second.X {
		if first.Y < second.Y {
			if clockwise {
				return first.Delta(shared.Loc(shared.RealEast))
			}
			return first.Delta(shared.Loc(shared.RealWest))
		}
		if clockwise {
			return first.Delta(shared.Loc(shared.RealWest))
		}
		return first.Delta(shared.Loc(shared.RealEast))
	}
	if first.X < second.X {
		if clockwise {
			return first.Delta(shared.Loc(shared.RealSouth))
		}
		return first.Delta(shared.Loc(shared.RealNorth))
	}
	if clockwise {
		return first.Delta(shared.Loc(shared.RealNorth))
	}
	return first.Delta(shared.Loc(shared.RealSouth))
}

func spawn(stash *gent.Set[shared.Loc], loc shared.Loc) {
	if stash.Has(loc) {
		return
	}
	stash.Add(loc)
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Add.", "loc", loc)
	}
	for _, dir := range shared.RealPrimaryDirections {
		spawn(stash, loc.Delta(shared.Loc(dir)))
	}
}
