package aoc2317

import (
	"math"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func DeriveLeastHeatLoss(lines []string, minJump, maxJump int) int {
	brd := shared.NewBoard(lines)
	shared.Logger.Info(
		"Derive least heat loss.",
		"width", brd.GetWidth(),
		"height", brd.GetHeight())
	return findMinimumHeat(brd, shared.Loc{X: brd.GetWidth() - 1}, minJump, maxJump)
}

type hop struct {
	loc shared.Loc
	dir shared.Direction
}

type runner struct {
	ID        int
	latestHop hop
	sum       int
	tail      *shared.Link[shared.Loc]
}

//revive:disable-next-line:function-length
func findMinimumHeat(brd *shared.Board, target shared.Loc, minJump, maxJump int) int {
	shared.Logger.Info(
		"Start running.",
		"target", target.ToString(),
		"min jump", minJump,
		"max jump", maxJump)
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
			tail:      shared.AddLink(nil, firstCell),
		},
		{
			ID:        createRunnerID(),
			latestHop: hop{loc: firstCell, dir: shared.RealSouth},
			sum:       0,
			tail:      shared.AddLink(nil, firstCell),
		},
	}
	runnerCount := 1

	minHeat := math.MaxInt
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
		}
	}

	// Cut total duration of aoc-2023-17 to about half by not constantly creating new slices of
	// hops.
	hops := make([]hop, 0, 7)
	for len(runners) > 0 {
		aRunner := runners[0]
		runners = runners[1:]

		if shouldStop(aRunner, hasher, target, minHeat) {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Should stop.", "tail", stringifyRoute(aRunner.tail))
			}
			continue
		}
		hops = hops[:0]
		hops = filterHops(
			hasher,
			aRunner,
			brd,
			deriveNextHops(brd, aRunner, hops, minJump, maxJump))
		for _, each := range hops {
			_, sum := countBetween(aRunner.latestHop.loc, each.loc, brd)
			r := &runner{
				ID:        aRunner.ID,
				latestHop: each,
				sum:       sum + aRunner.sum,
				tail:      shared.AddLink(aRunner.tail, each.loc),
			}
			if r.sum > minHeat {
				continue
			}
			if each.loc == target {
				finishRace(r)
				continue
			}
			if aRunner.latestHop.dir == each.dir {
				distance := 1
				sum := aRunner.sum
				doInBetween(
					aRunner.latestHop.loc,
					each.loc,
					func(loc shared.Loc) {
						sum += brd.GetIntOrDie(loc)
						hasher.set(
							hop{
								loc: loc,
								dir: each.dir,
							},
							distance,
							sum,
						)
					})
				continue
			}
			if !hasher.set(r.latestHop, 0, r.sum) {
				continue
			}
			r.ID = createRunnerID()
			runners = append(runners, r)
			runnerCount++
		}
	}
	shared.Logger.Info("Done running.", "runner count", runnerCount, "min heat", minHeat)
	return minHeat
}

func doInBetween(from, to shared.Loc, callback func(loc shared.Loc)) {
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
		callback(loc)
	}
}

func measureDistance(a, b shared.Loc) int {
	x := shared.Abs(a.X - b.X)
	y := shared.Abs(a.Y - b.Y)
	return x + y
}

func deriveNextHops(brd *shared.Board, aRunner *runner, hops []hop, minJump, maxJump int) []hop {
	var straightCandidate *shared.Loc
	for stepCount := minJump; stepCount <= maxJump; stepCount++ {
		delta := shared.Loc(aRunner.latestHop.dir)
		if delta.X != 0 {
			delta.X *= stepCount
		} else {
			delta.Y *= stepCount
		}
		straight := aRunner.latestHop.loc.Delta(delta)
		if stepCount > 0 {
			_, ok := brd.Get(straight)
			if !ok {
				break
			}
			straightCandidate = &straight
		}

		add := func(loc shared.Loc, dir shared.Direction) {
			if _, ok := brd.Get(loc); ok {
				hops = append(hops, hop{loc: loc, dir: dir})
			}
		}

		leftDir := aRunner.latestHop.dir.TurnRealLeft()
		add(straight, leftDir)

		rightDir := aRunner.latestHop.dir.TurnRealRight()
		add(straight, rightDir)
	}
	if straightCandidate != nil {
		hops = append(hops, hop{
			loc: *straightCandidate,
			dir: aRunner.latestHop.dir,
		})
	}
	return hops
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

func shouldStop(aRunner *runner, hasher *sumHasher, target shared.Loc, minHeat int) bool {
	if aRunner.sum+measureDistance(aRunner.latestHop.loc, target) > minHeat {
		return true
	}
	if hasher.isOver(aRunner) {
		return true
	}
	return false
}

func filterHops(
	hasher *sumHasher,
	aRunner *runner,
	brd *shared.Board,
	candidates []hop,
) []hop {
	for i := len(candidates) - 1; i >= 0; i-- {
		each := candidates[i]
		arrowCount, sum := countBetween(aRunner.latestHop.loc, each.loc, brd)
		if aRunner.latestHop.dir != each.dir {
			arrowCount = 0
		}
		if isGoingOverboard(brd, each) ||
			hasher.isOverWithDetails(sum, each.loc, each.dir, arrowCount) {
			candidates = append(candidates[:i], candidates[i+1:]...)
		}
	}
	return candidates
}

type record struct {
	dir      shared.Direction
	velocity int
	sum      int
}

type sumHasher struct {
	// Avoid expensive operations by storing coordinates separately here,
	// those that have records in recs.
	s     [][]int
	recs  [][]record
	width int
}

func newSumHasher(width, height int) *sumHasher {
	s := make([][]int, height)
	for i := range height {
		s[i] = make([]int, width)
	}
	recs := make([][]record, width*height)
	return &sumHasher{s: s, recs: recs, width: width}
}

func (v *sumHasher) deriveRecordIndex(loc shared.Loc) int {
	return loc.Y*v.width + loc.X
}

func (v *sumHasher) isOver(r *runner) bool {
	if v.s[r.latestHop.loc.Y][r.latestHop.loc.X] == 0 {
		return false
	}
	return v.isOverWithDetails(r.sum, r.latestHop.loc, r.latestHop.dir, 0)
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
	recIndex := v.deriveRecordIndex(loc)
	locRecs := v.recs[recIndex]
	if len(locRecs) == 0 {
		return false
	}
	for _, each := range locRecs {
		if each.dir == dir {
			if sum > each.sum && arrow >= each.velocity {
				return true
			}
		}
	}
	return false
}

func (v *sumHasher) set(aHop hop, arrow, sum int) bool {
	recIndex := v.deriveRecordIndex(aHop.loc)
	locRecs := v.recs[recIndex]
	index := -1
	for i, each := range locRecs {
		if each.dir == aHop.dir && arrow == each.velocity {
			if each.sum < sum {
				return false
			}
			index = i
		}
	}
	if index >= 0 {
		locRecs[index].sum = sum
		return true
	}

	v.recs[recIndex] = append(v.recs[recIndex], record{
		dir:      aHop.dir,
		velocity: arrow,
		sum:      sum,
	})
	v.s[aHop.loc.Y][aHop.loc.X] = 1
	return true
}

func countBetween(from, to shared.Loc, brd *shared.Board) (count int, sum int) {
	xd := to.X - from.X
	if xd != 0 {
		xd /= shared.Abs(xd)
		for x := from.X + xd; x != to.X+xd; x += xd {
			loc := shared.Loc{X: x, Y: from.Y}
			i := brd.GetIntOrDie(loc)
			count++
			sum += i
		}
		return
	}

	yd := to.Y - from.Y
	if yd != 0 {
		yd /= shared.Abs(yd)
		for y := from.Y + yd; y != to.Y+yd; y += yd {
			loc := shared.Loc{X: from.X, Y: y}
			i := brd.GetIntOrDie(loc)
			count++
			sum += i
		}
		return
	}
	return
}

func isGoingOverboard(brd *shared.Board, h hop) bool {
	if h.loc.X == 0 && h.dir == shared.RealWest {
		return true
	}
	if h.loc.Y == 0 && h.dir == shared.RealSouth {
		return true
	}
	if h.loc.X == brd.GetWidth()-1 && h.dir == shared.RealEast {
		return true
	}
	if h.loc.Y == brd.GetHeight()-1 && h.dir == shared.RealNorth {
		return true
	}
	return false
}
