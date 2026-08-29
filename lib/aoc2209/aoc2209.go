package aoc2209

import (
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func CountVisitedPositions(lines []string, tailSize int) int {
	steps := parseLines(lines)
	snake := make([]shared.Loc, tailSize+1)
	tails := gent.NewSet(snake[len(snake)-1])
	for _, each := range steps {
		for range each.count {
			prevHead := snake[0]
			snake[0] = snake[0].Delta(each.dir)
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Move head.", "from", prevHead, "to", snake[0])
			}
			for i := 0; i < len(snake)-1; i++ {
				head := snake[i]
				tail := snake[i+1]
				tailDelta, ok := deriveTailDelta(head, tail)
				if ok {
					prevTail := tail
					tail = tail.Delta(tailDelta)
					snake[i+1] = tail
					if shared.IsDebugEnabled() {
						shared.Logger.Debug("Move tail.", "i", i, "from", prevTail, "to", tail)
					}
				}
			}
			tails.Add(snake[len(snake)-1])
		}
	}
	return tails.Count()
}

type step struct {
	dir   shared.Loc
	count int
}

func parseLines(lines []string) []step {
	steps := make([]step, 0, len(lines))
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		pieces := strings.Fields(trimmed)
		if len(pieces) != 2 {
			shared.Logger.Error("Invalid line.", "line", each)
			panic("invalid line")
		}
		dir := toDirection(pieces[0])
		count := gent.OrPanic2(strconv.Atoi(pieces[1]))("failed to parse step count")
		steps = append(steps, step{
			dir:   dir,
			count: count,
		})
	}
	return steps
}

func toDirection(s string) shared.Loc {
	switch strings.ToLower(s)[0] {
	case 'r':
		return shared.Loc(shared.RealEast)
	case 'd':
		return shared.Loc(shared.RealSouth)
	case 'l':
		return shared.Loc(shared.RealWest)
	case 'u':
		return shared.Loc(shared.RealNorth)
	default:
		shared.Logger.Error("Invalid direction string.", "string", s)
		panic("invalid direction string")
	}
}

func deriveTailDelta(head, tail shared.Loc) (loc shared.Loc, move bool) {
	xDelta := head.X - tail.X
	yDelta := head.Y - tail.Y
	if shared.Abs(xDelta) <= 1 && shared.Abs(yDelta) <= 1 {
		return
	}
	move = true
	if xDelta == 0 {
		loc = shared.Loc{X: 0, Y: yDelta / shared.Abs(yDelta)}
		return
	}
	if yDelta == 0 {
		loc = shared.Loc{X: xDelta / shared.Abs(xDelta), Y: 0}
		return
	}
	xDelta /= shared.Abs(xDelta)
	yDelta /= shared.Abs(yDelta)
	loc = shared.Loc{X: xDelta, Y: yDelta}
	return
}
