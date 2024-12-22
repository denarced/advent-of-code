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

func CountLowestScore(lines []string, drawWinners bool) int {
	shared.Logger.Info("Count lowest score.", "draw winners", drawWinners, "line count", len(lines))
	possibleDirections := map[shared.Direction][]shared.Direction{
		shared.RealEast:  getPossibleDirections(shared.RealEast),
		shared.RealSouth: getPossibleDirections(shared.RealSouth),
		shared.RealWest:  getPossibleDirections(shared.RealWest),
		shared.RealNorth: getPossibleDirections(shared.RealNorth),
	}
	brd := shared.NewBoard(lines)
	start := brd.FindOrDie('S')
	breaker := &pointBreak{
		m: map[vector]int{},
	}
	runners := &queue[runner]{
		s: []runner{{
			vec:   vector{loc: start, dir: shared.RealEast},
			steps: []shared.Loc{start},
		}},
	}
	brd.Set(start, '.')
	minPoints := 999_999_999_999
	counter := 0
	winners := []runner{}
	for i := 0; i < 1_000_000_000 && !runners.isEmpty(); i++ {
		counter++
		r := runners.pop()
		if r.points >= minPoints {
			continue
		}
		if breaker.isHigher(r.vec, r.points) {
			continue
		}
		if brd.GetOrDie(r.vec.loc) == 'E' {
			shared.Logger.Info("Finished.", "points", r.points)
			if drawWinners {
				if minPoints == r.points {
					winners = append(winners, r)
				} else if minPoints > r.points {
					winners = []runner{r}
				}
			}
			minPoints = shared.Min(minPoints, r.points)
			continue
		}
		possibleVectors := getPossibleVectors(brd, r.vec, possibleDirections)
		copyCount := 0
		for _, each := range possibleVectors {
			cor := corner{vec: r.vec, turn: each.dir}
			if slices.Contains(r.corners, cor) {
				continue
			}
			copyCount++
			copied := r
			copied.points += derivePoints(copied.vec.dir, each.dir)
			nextLoc := each.loc.Delta(shared.Loc(each.dir))
			copied.vec.loc = nextLoc
			if copyCount > 1 {
				// Having multiple added runners somehow corrupts previously added runner's steps so
				// we need to copy. However, constantly copying the steps is expensive so it's only
				// done when necessary. With the actual input, the difference between copying and
				// not copying is 17s vs 1m05s. With this avoidance, it's 24s.
				copied.steps = append([]shared.Loc{}, copied.steps...)
				copied.steps = append(copied.steps, nextLoc)
			} else {
				copied.steps = append(copied.steps, nextLoc)
			}
			copied.vec.dir = each.dir
			if len(possibleVectors) > 1 {
				copied.corners = append(copied.corners, cor)
			}
			runners.add(copied)
		}
	}
	shared.Logger.Info(
		"Runner count at end.",
		"count",
		runners.count(),
		"counter",
		counter,
		"winner count",
		len(winners),
	)
	for _, each := range winners {
		draw(lines, each.steps)
	}
	return minPoints
}

type vector struct {
	loc shared.Loc
	dir shared.Direction
}

type runner struct {
	points  int
	vec     vector
	corners []corner
	steps   []shared.Loc
}

type corner struct {
	vec  vector
	turn shared.Direction
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

func getPossibleDirections(dir shared.Direction) []shared.Direction {
	all := append([]shared.Direction{}, shared.RealDirections...)
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
	m map[vector]int
}

func (v *pointBreak) isHigher(vec vector, points int) bool {
	min, exists := v.m[vec]
	if !exists || min > points {
		v.m[vec] = points
		return false
	}
	return true
}

func draw(lines []string, steps []shared.Loc) {
	nanos := time.Now().UnixNano()
	dirp := fmt.Sprintf("/tmp/aoc16/%d", nanos)
	err := os.MkdirAll(dirp, 0755)
	shared.Die(err, "Failed to create draw dir.")
	brd := shared.NewBoard(append([]string{}, lines...))
	for _, each := range steps {
		brd.Set(each, 'O')
	}
	content := strings.Join(brd.GetLines(), "\n") + "\n"
	filep := filepath.Join(dirp, "board.txt")
	err = os.WriteFile(filep, []byte(content), 0644)
	shared.Die(err, "Failed to write board.txt.")
	shared.Logger.Info("Winner drawn.", "filepath", filep)
}
