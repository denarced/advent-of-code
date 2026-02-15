package aoc2511

import (
	"fmt"
	"slices"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

type counter struct {
	count int
}

func newCounter() *counter {
	return &counter{}
}

func (v *counter) increment() {
	v.count++
}

func CountPaths(lines []string, start string) int {
	shared.Logger.Info("Count paths.", "start", start)
	points := parseTree(lines)
	if start == "you" {
		you := points[start]
		aCounter := newCounter()
		comb(you, points["out"], true, nil, aCounter)
		shared.Logger.Info("Count done.", "count", aCounter.count)
		return aCounter.count
	}

	below := points["fft"]
	above := find(below, "dac", false)
	if above == nil {
		above = find(below, "dac", true)
		swap := above
		above = below
		below = swap
	}
	if above == nil || below == nil {
		panic("dac not found")
	}
	svr := findPoint("svr", above, false)
	out := findPoint("out", below, true)

	var counts []int
	var count int
	for _, each := range []struct {
		from, to *Point
		dig      bool
	}{
		{above, svr, false},
		{below, out, true},
		{above, below, true},
	} {
		for _, p := range points {
			p.foundByDevil = false
			p.foundByAsshole = false
			p.passedByDevil = false
			p.passedByAsshole = false
		}
		count = 0
		var allowed *gent.Set[*Point]
		if each.dig {
			web(each.from, each.to)
		} else {
			web(each.to, each.from)
		}
		allowed = FilterDoubles(points)
		shared.Logger.Debug("Allowed.", "allowed", allowed.Count())
		aCounter := &counter{}
		comb(each.from, each.to, each.dig, allowed, aCounter)
		counts = append(counts, aCounter.count)
	}

	shared.Logger.Info("All counts.", "counts", counts)
	for i, each := range counts {
		if each <= 0 {
			panic(fmt.Sprintf("count %d <= 0: %d", i, each))
		}
	}
	count = 0
	for i, each := range counts {
		if i == 0 {
			count = each
			continue
		}
		count *= each
	}
	shared.Logger.Info("Count done.", "count", count)
	return count
}

func find(start *Point, target string, dig bool) *Point {
	next := gent.Tri(dig, start.kids, start.parents)
	for name, sub := range next {
		if name == target {
			return sub
		}
		if res := find(sub, target, dig); res != nil {
			return res
		}
	}
	return nil
}

func comb(
	current, target *Point,
	dig bool,
	allowed *gent.Set[*Point],
	aCounter *counter,
) {
	next := gent.Tri(dig, current.kids, current.parents)
	for _, sub := range next {
		if sub == target {
			aCounter.increment()
			continue
		}
		if allowed != nil && !allowed.Has(sub) {
			continue
		}
		comb(sub, target, dig, allowed, aCounter)
	}
}

func parseLine(line string) (string, []string) {
	pieces := strings.Split(strings.TrimSpace(line), ":")
	return strings.TrimSpace(pieces[0]), strings.Fields(pieces[1])
}

type Point struct {
	parents         map[string]*Point
	kids            map[string]*Point
	foundByDevil    bool
	foundByAsshole  bool
	passedByDevil   bool
	passedByAsshole bool
}

func FilterDoubles(points map[string]*Point) *gent.Set[*Point] {
	doubles := gent.NewSet[*Point]()
	for _, each := range points {
		if each.foundByAsshole && each.foundByDevil {
			doubles.Add(each)
		}
	}
	return doubles
}

func parseTree(lines []string) map[string]*Point {
	points := map[string]*Point{}
	for _, each := range lines {
		from, to := parseLine(each)
		parent := points[from]
		if parent == nil {
			parent = &Point{
				parents: map[string]*Point{},
				kids:    map[string]*Point{},
			}
			points[from] = parent
		}
		for _, name := range to {
			child := points[name]
			if child == nil {
				child = &Point{
					parents: map[string]*Point{from: parent},
					kids:    map[string]*Point{},
				}
				points[name] = child
			} else {
				if child.parents[from] != nil {
					panic("shouldn't have this parent yet")
				}
				child.parents[from] = parent
			}
			if parent.kids[name] != nil {
				panic("shouldn't exist yet")
			}
			parent.kids[name] = child
		}
	}
	return points
}

type trail struct {
	p      *Point
	crumbs []*Point
}

func webIt(points map[string]*Point) {
	below := points["fft"]
	above := find(below, "dac", false)
	if above == nil {
		above = find(below, "dac", true)
		swap := above
		above = below
		below = swap
	}
	if above == nil || below == nil {
		panic("dac not found")
	}
	web(above, below)
}

func web(above, below *Point) {
	uppers := []trail{{p: below}}
	downers := []trail{{p: above}}
	for {
		if len(uppers) == 0 && len(downers) == 0 {
			break
		}
		uppers = processLayer(
			uppers,
			func(p *Point) map[string]*Point { return p.parents },
			func(p *Point) bool { return p.foundByAsshole },
			func(p *Point) { p.foundByDevil = true },
			func(p *Point) bool { return p.passedByDevil },
			func(p *Point) { p.passedByDevil = true },
		)
		downers = processLayer(
			downers,
			func(p *Point) map[string]*Point { return p.kids },
			func(p *Point) bool { return p.foundByDevil },
			func(p *Point) { p.foundByAsshole = true },
			func(p *Point) bool { return p.passedByAsshole },
			func(p *Point) { p.passedByAsshole = true },
		)
	}
}

func sortedKeys[T any](m map[string]T) []string {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, len(m))
	var i int
	for each := range m {
		keys[i] = each
		i++
	}
	slices.Sort(keys)
	return keys
}

func processLayer(
	trails []trail,
	getNext func(p *Point) map[string]*Point,
	checkFound func(p *Point) bool,
	markFound func(p *Point),
	checkPassed func(p *Point) bool,
	markPassed func(p *Point),
) []trail {
	var next []trail
	for _, each := range trails {
		var ignited bool
		for _, key := range sortedKeys(getNext(each.p)) {
			p := getNext(each.p)[key]
			markFound(p)
			if checkFound(p) {
				if !ignited {
					igniteTrail(each.crumbs)
					ignited = true
				}
			}
			if checkPassed(p) {
				continue
			}
			markPassed(p)
			crumbs := append(each.crumbs, p)
			next = append(next, trail{p: p, crumbs: crumbs})
		}
	}
	return next
}

func igniteTrail(crumbs []*Point) {
	for _, each := range crumbs {
		each.foundByDevil = true
		each.foundByAsshole = true
	}
}

func findPoint(name string, from *Point, dig bool) *Point {
	shared.Logger.Info("Find point.", "name", name)
	point := find(from, name, dig)
	if point == nil {
		panic("svr not found")
	}
	return point
}
