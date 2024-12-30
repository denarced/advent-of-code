package aoc2416

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/denarced/advent-of-code/shared"
)

const (
	pointsStep = 1
	pointsTurn = 1_000
)

func CountLowestScore(lines []string, drawWinners bool) (minPoints int, seatCount int) {
	shared.Logger.Info("Count lowest score.", "draw winners", drawWinners, "line count", len(lines))
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	end := brd.FindOrDie('E')
	possibleDirections := derivePossibleDirections(start, end)
	brd.Set(start, '.')
	brd.ReadOnly = true
	breaker := &pointBreak{
		m: map[vector]int{},
	}
	runners, loopCount, winners, minPoints := count(
		breaker,
		start,
		end,
		brd,
		possibleDirections,
		999_999_999_999,
	)
	shared.Logger.Info(
		"Midpoint numbers.",
		"runner count",
		runners.count(),
		"loop count",
		loopCount,
		"winner count",
		len(winners),
		"breaker size",
		len(breaker.m),
		"min score",
		minPoints,
	)
	runners, loopCount, winners, minPoints = count(
		breaker,
		start,
		end,
		brd,
		possibleDirections,
		minPoints,
	)
	shared.Logger.Info(
		"End numbers.",
		"runner count",
		runners.count(),
		"loop count",
		loopCount,
		"winner count",
		len(winners),
		"breaker size",
		len(breaker.m),
		"min score",
		minPoints,
	)
	bestSeats := shared.NewSet([]shared.Loc{})
	for _, winner := range winners {
		if drawWinners {
			draw(lines, winner.steps)
		}
		for step := winner.steps; step != nil; step = step.parent {
			bestSeats.Add(step.item)
		}
	}
	seatCount = bestSeats.Count()
	shared.Logger.Info("Seats counted.", "count", seatCount)
	return minPoints, seatCount
}

//revive:disable-next-line:function-result-limit
func count(
	breaker *pointBreak,
	start, end shared.Loc,
	brd *shared.Board,
	possibleDirections map[shared.Direction][]shared.Direction,
	minPoints int,
) (*queue[runner], int, []runner, int) {
	runners := &queue[runner]{
		s: []runner{{
			vec:   vector{loc: start, dir: shared.RealEast},
			steps: addLink(nil, start),
		}},
	}
	includeAllWinners := len(breaker.m) > 0
	if includeAllWinners {
		breaker.loose = true
	}
	counter := 0
	winners := []runner{}
	for i := 0; i < 1_000_000_000 && !runners.isEmpty(); i++ {
		counter++
		r := runners.pop()
		if includeAllWinners && r.points > minPoints ||
			!includeAllWinners && r.points >= minPoints {
			continue
		}
		if breaker.isHigher(r.vec, r.points) {
			continue
		}
		if r.vec.loc == end {
			if minPoints == r.points {
				winners = append(winners, r)
			} else if minPoints > r.points {
				winners = []runner{r}
			}
			minPoints = shared.Min(minPoints, r.points)
			continue
		}
		possibleVectors := getPossibleVectors(brd, r.vec, possibleDirections)
		for _, each := range possibleVectors {
			copied := r
			copied.points += derivePoints(copied.vec.dir, each.dir)
			nextLoc := each.loc.Delta(shared.Loc(each.dir))
			copied.vec.loc = nextLoc
			copied.steps = addLink(copied.steps, nextLoc)
			copied.vec.dir = each.dir
			runners.add(copied)
		}
	}
	return runners, counter, winners, minPoints
}

type vector struct {
	loc shared.Loc
	dir shared.Direction
}

type runner struct {
	points int
	vec    vector
	steps  *link[shared.Loc]
}

func getPossibleVectors(
	brd *shared.Board,
	vec vector,
	possibleDirections map[shared.Direction][]shared.Direction,
) []vector {
	vecs := []vector{}
	for _, d := range possibleDirections[vec.dir] {
		loc := vec.loc.Delta(shared.Loc(d))
		c := brd.GetOrDie(loc)
		if c == '.' || c == 'E' {
			vecs = append(vecs, vector{loc: vec.loc, dir: d})
		}
	}
	return vecs
}

func getPossibleDirections(all []shared.Direction, dir shared.Direction) []shared.Direction {
	rev := shared.Direction{X: -dir.X, Y: -dir.Y}
	return shared.FilterValues(all, func(dir shared.Direction) bool {
		return dir != rev
	})
}

func derivePoints(previous, current shared.Direction) int {
	if previous == current {
		return pointsStep
	}
	return pointsTurn + pointsStep
}

type queue[T any] struct {
	s []T
}

func (v *queue[T]) pop() T {
	end := len(v.s) - 1
	item := v.s[end]
	v.s = v.s[:end]
	return item
}

func (v *queue[T]) add(item T) {
	v.s = append(v.s, item)
}

func (v *queue[T]) isEmpty() bool {
	return len(v.s) == 0
}

func (v *queue[T]) count() int {
	return len(v.s)
}

type pointBreak struct {
	m     map[vector]int
	loose bool
}

func (v *pointBreak) isHigher(vec vector, points int) bool {
	min, exists := v.m[vec]
	if !exists || (!v.loose && min > points) || (v.loose && min >= points) {
		if !exists || min > points {
			v.m[vec] = points
		}
		return false
	}
	return true
}

func draw(lines []string, step *link[shared.Loc]) {
	nanos := time.Now().UnixNano()
	dirp := fmt.Sprintf("/tmp/aoc16/%d", nanos)
	err := os.MkdirAll(dirp, 0755)
	shared.Die(err, "Failed to create draw dir.")
	brd := shared.NewBoard(append([]string{}, lines...))
	for ; step != nil; step = step.parent {
		brd.Set(step.item, 'O')
	}
	content := strings.Join(brd.GetLines(), "\n") + "\n"
	filep := filepath.Join(dirp, "board.txt")
	err = os.WriteFile(filep, []byte(content), 0644)
	shared.Die(err, "Failed to write board.txt.")
	shared.Logger.Info("Winner drawn.", "filepath", filep)
}

func sortDirections(start, end shared.Loc) []shared.Direction {
	x := toZeroOrOne(end.X, start.X)
	y := toZeroOrOne(end.Y, start.Y)
	all := append([]shared.Direction{}, shared.RealPrimaryDirections...)
	slices.SortFunc(all, func(a, b shared.Direction) int {
		grade := func(dir shared.Direction) int {
			points := 0
			if dir.X != x {
				points++
			}
			if dir.Y != y {
				points++
			}
			return points
		}
		return grade(a) - grade(b)
	})
	return all
}

func toZeroOrOne(first, second int) int {
	difference := first - second
	if difference == 0 {
		return difference
	}
	return difference / shared.Abs(difference)
}

func derivePossibleDirections(
	start, end shared.Loc,
) map[shared.Direction][]shared.Direction {
	// At least usually it'll go roughly to the right direction.
	// Reduces duration with real input from 24s to 18s.
	all := sortDirections(start, end)
	return map[shared.Direction][]shared.Direction{
		shared.RealEast:  getPossibleDirections(all, shared.RealEast),
		shared.RealSouth: getPossibleDirections(all, shared.RealSouth),
		shared.RealWest:  getPossibleDirections(all, shared.RealWest),
		shared.RealNorth: getPossibleDirections(all, shared.RealNorth),
	}
}

type link[T any] struct {
	item   T
	parent *link[T]
}

func addLink[T any](parent *link[T], item T) *link[T] {
	return &link[T]{
		item:   item,
		parent: parent,
	}
}
