package aoc2317

import (
	"encoding/json"
	"math"
	"os"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const maxJump = 3

func DeriveLeastHeatLoss(lines []string) int {
	brd := shared.NewBoard(lines)
	shared.Logger.Info(
		"Derive least heat loss.",
		"width", brd.GetWidth(),
		"height", brd.GetHeight())
	return findMinimumHeat(brd, shared.Loc{X: brd.GetWidth() - 1})
}

type hop struct {
	loc shared.Loc
	dir shared.Direction
}

type runner struct {
	ID        int
	latestHop hop
	sum       int
	arrow     shared.Loc
	tail      *shared.Link[shared.Loc]
}

func (v *runner) setArrow(aHop hop) {
	if v.latestHop.dir == aHop.dir {
		if v.latestHop.dir.X != 0 {
			v.arrow.X++
		} else {
			v.arrow.Y++
		}
		return
	}

	v.arrow = shared.Loc{}
	if aHop.dir.X != 0 {
		v.arrow.X = 1
	} else {
		v.arrow.Y = 1
	}
}

//revive:disable-next-line:function-length,cyclomatic,cognitive-complexity
func findMinimumHeat(brd *shared.Board, target shared.Loc) int {
	shared.Logger.Info("Start running.", "target", target.ToString())
	hasher := newSumHasher(brd.GetWidth(), brd.GetHeight())
	firstCell := shared.Loc{Y: brd.GetHeight() - 1}
	createRunnerID := func() func() int {
		var id int
		return func() int {
			next := id
			id++
			return next
		}
	}()
	runners := []*runner{
		{
			ID:        createRunnerID(),
			latestHop: hop{loc: firstCell, dir: shared.RealEast},
			sum:       0,
			arrow:     shared.Loc{X: 1},
			tail:      shared.AddLink(nil, firstCell),
		},
	}
	runnerCount := 1

	minHeat := math.MaxInt
	solutions := []map[string]any{}
	exportFilep := os.Getenv("data_filep")
	if exportFilep != "" {
		defer func() {
			values := make([][]int, brd.GetHeight())
			row := brd.GetHeight() - 1
			for y := range brd.GetHeight() {
				values[y] = make([]int, brd.GetWidth())
				for x := range brd.GetWidth() {
					values[y][x] = brd.GetIntOrDie(shared.Loc{X: x, Y: row})
				}
				row--
			}
			exports := map[string]any{
				"width":     brd.GetWidth(),
				"height":    brd.GetHeight(),
				"solutions": solutions,
				"values":    values,
			}
			f, err := os.Create(exportFilep)
			if err != nil {
				panic(err)
			}
			b, err := json.Marshal(exports)
			if err != nil {
				panic(err)
			}
			_, err = f.Write(b)
			if err != nil {
				panic(err)
			}
		}()
	}
	finishRace := func(r *runner) {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug(
				"Race finished.",
				"runner ID", r.ID,
				"sum", r.sum)
		}
		if minHeat > r.sum {
			minHeat = r.sum
			shared.Logger.Info(
				"New record achieved.",
				"sum", r.sum,
				"route", stringifyRoute(r.tail))
			if exportFilep != "" {
				m := map[string]any{
					"id":    r.ID,
					"route": revertCoordinates(extractRoute(r.tail), brd.GetHeight()),
					"total": r.sum,
				}
				solutions = append(solutions, m)
			}
		}
	}

	for len(runners) > 0 {
		aRunner := runners[0]
		runners = runners[1:]
		for {
			if shouldStop(aRunner, hasher, target, minHeat) {
				break
			}
			hops := filterHops(hasher, aRunner, brd, deriveNextHops(brd, aRunner))
			if len(hops) == 0 {
				break
			}
			for _, each := range hops[1:] {
				r := &runner{
					ID:        createRunnerID(),
					latestHop: each,
					sum:       aRunner.sum + brd.GetIntOrDie(each.loc),
					arrow:     deriveArrow(each.dir),
					tail:      shared.AddLink(aRunner.tail, each.loc),
				}
				if each.loc == target {
					finishRace(r)
					continue
				}
				runners = append(runners, r)
				runnerCount++
				hasher.set(r.latestHop, r.arrow, r.sum)
			}
			nextHop := hops[0]
			aRunner.setArrow(nextHop)
			aRunner.latestHop = nextHop
			aRunner.sum += brd.GetIntOrDie(nextHop.loc)
			aRunner.tail = shared.AddLink(aRunner.tail, nextHop.loc)
			hasher.set(nextHop, aRunner.arrow, aRunner.sum)

			if aRunner.latestHop.loc == target {
				finishRace(aRunner)
				break
			}
		}
	}
	shared.Logger.Info("Done running.", "runner count", runnerCount, "min heat", minHeat)
	return minHeat
}

func between(from, to shared.Loc) (steps []shared.Loc) {
	if from == to {
		return
	}
	diff := to.X - from.X
	var step shared.Loc
	if diff == 0 {
		step = shared.Loc{Y: gent.Tri((to.Y-from.Y) < 0, -1, 1)}
	} else {
		step = shared.Loc{X: gent.Tri((to.X-from.X) < 0, -1, 1)}
	}
	start := from.Delta(shared.Loc(step))
	for loc := start; loc != to; loc = loc.Delta(step) {
		steps = append(steps, loc)
	}
	return
}

func measureDistance(a, b shared.Loc) int {
	x := shared.Abs(a.X - b.X)
	y := shared.Abs(a.Y - b.Y)
	return x + y
}

func deriveNextHops(brd *shared.Board, aRunner *runner) []hop {
	straight := aRunner.latestHop.loc.Delta(shared.Loc(aRunner.latestHop.dir))
	var hops []hop
	if _, ok := brd.Get(straight); ok {
		if !(aRunner.latestHop.dir.X != 0 && shared.Abs(aRunner.arrow.X) >= maxJump) &&
			!(aRunner.latestHop.dir.Y != 0 && shared.Abs(aRunner.arrow.Y) >= maxJump) {
			hops = append(hops, hop{loc: straight, dir: aRunner.latestHop.dir})
		}
	}

	add := func(loc shared.Loc, dir shared.Direction) {
		if _, ok := brd.Get(loc); ok {
			hops = append(hops, hop{loc: loc, dir: dir})
		}
	}

	leftDir := aRunner.latestHop.dir.TurnRealLeft()
	left := aRunner.latestHop.loc.Delta(shared.Loc(leftDir))
	add(left, leftDir)

	rightDir := aRunner.latestHop.dir.TurnRealRight()
	right := aRunner.latestHop.loc.Delta(shared.Loc(rightDir))
	add(right, rightDir)
	return hops
}

func deriveArrow(dir shared.Direction) shared.Loc {
	if dir.X != 0 {
		return shared.Loc{X: 1}
	}
	return shared.Loc{Y: 1}
}

func stringifyRoute(head *shared.Link[shared.Loc]) string {
	if head == nil {
		return "nil"
	}
	var total int
	for l := head; l != nil; l = l.Parent {
		total++
	}
	links := make([]string, total)
	curr := head
	for i := total - 1; i >= 0; i-- {
		links[i] = curr.Item.ToString()
		curr = curr.Parent
	}
	return strings.Join(links, " -> ")
}

func extractRoute(head *shared.Link[shared.Loc]) [][]int {
	if head == nil {
		return nil
	}
	var total int
	for l := head; l != nil; l = l.Parent {
		total++
	}
	links := make([][]int, total)
	curr := head
	for i := total - 1; i >= 0; i-- {
		links[i] = []int{curr.Item.X, curr.Item.Y}
		curr = curr.Parent
	}
	return links
}

func shouldStop(aRunner *runner, hasher *sumHasher, target shared.Loc, minHeat int) bool {
	if aRunner.sum+measureDistance(aRunner.latestHop.loc, target) > minHeat {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug(
				"Runner can't even theoretically beat the current minimum.",
				"route", stringifyRoute(aRunner.tail))
		}
		return true
	}
	if hasher.isOver(aRunner) {
		if shared.IsDebugEnabled() {
			shared.Logger.Debug(
				"Runner already has too high a sum.",
				"sum", aRunner.sum,
				"route", stringifyRoute(aRunner.tail),
			)
		}
		return true
	}
	return false
}

func filterHops(
	hasher *sumHasher,
	aRunner *runner,
	brd *shared.Board,
	candidates []hop,
) (result []hop) {
	for _, each := range candidates {
		if !hasher.isOverWithDetails(aRunner.sum+brd.GetIntOrDie(each.loc), each.loc, each.dir, 1) {
			result = append(result, each)
		}
	}
	return
}

func revertCoordinates(coords [][]int, height int) [][]int {
	reverted := make([][]int, len(coords))
	for i, row := range coords {
		reverted[i] = make([]int, len(row))
		reverted[i][0] = row[0]
		reverted[i][1] = height - row[1] - 1
	}
	return reverted
}

type sumHasher struct {
	// Avoid unnecessary hashing by storing coordinates separate here, those
	// that have a sum in m.
	s [][]int
	m map[uint64]int
}

func newSumHasher(width, height int) *sumHasher {
	s := make([][]int, height)
	for i := range height {
		s[i] = make([]int, width)
	}
	return &sumHasher{m: map[uint64]int{}, s: s}
}

func (v *sumHasher) isOver(r *runner) bool {
	if v.s[r.latestHop.loc.Y][r.latestHop.loc.X] == 0 {
		return false
	}
	_, current, ok := v.hash(r.latestHop.loc, r.latestHop.dir, v.deriveArrowValue(r.arrow))
	if ok {
		return r.sum > current
	}
	return false
}

func (*sumHasher) deriveArrowValue(arrow shared.Loc) int {
	if arrow.X > 0 {
		return arrow.X
	}
	return arrow.Y
}

func (v *sumHasher) isOverWithDetails(
	sum int,
	loc shared.Loc,
	dir shared.Direction,
	arrow int,
) bool {
	if v.s[loc.Y][loc.X] == 0 {
		return false
	}
	_, current, ok := v.hash(loc, dir, arrow)
	if ok {
		return sum > current
	}
	return false
}

func (v *sumHasher) hash(
	loc shared.Loc,
	dir shared.Direction,
	arrow int,
) (hashSum uint64, sum int, ok bool) {
	hashSum = hashInts(
		[]int{
			loc.X,
			loc.Y,
			dir.X,
			dir.Y,
			arrow,
		})
	sum, ok = v.m[hashSum]
	return
}

func hashInts(arr []int) uint64 {
	// FNV-1a constants
	var hash uint64 = 14695981039346656037
	const prime uint64 = 1099511628211

	for _, v := range arr {
		hash ^= uint64(v)
		hash *= prime
	}
	return hash
}

func (v *sumHasher) set(aHop hop, arrow shared.Loc, sum int) {
	h, current, ok := v.hash(aHop.loc, aHop.dir, v.deriveArrowValue(arrow))
	if ok && current < sum {
		shared.Logger.Error("Attempt to increase sum in hasher.",
			"hop", aHop,
			"arrow", arrow,
			"sum", sum,
			"current sum", current)
		panic("can't increase heat in hasher")
	}
	v.m[h] = sum
	v.s[aHop.loc.Y][aHop.loc.X] = 1
}
