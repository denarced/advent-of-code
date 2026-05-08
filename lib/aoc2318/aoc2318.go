package aoc2318

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func Dig(lines []string, magic bool) int {
	shared.Logger.Info("Start digging.", "line count", len(lines), "magic", magic)
	instructions := parseLines(lines, magic)
	clockwise := isClockwise(instructions)
	routes, xCoords := gatherRoutes(instructions, clockwise)
	routes = finishRoutes(routes, xCoords)
	shared.Logger.Info("Routes ready.", "count", len(routes))

	pop := createIterator(routes)
	var total int
	for {
		desc := popRoutes(pop, true)
		if len(desc) == 0 {
			panic("0 desc fucks to give")
		}

		asc := popRoutes(pop, false)
		if len(asc) == 0 {
			panic("0 asc fucks to give")
		}

		alpha := desc[0]
		omega := asc[len(asc)-1]
		if alpha.from != omega.from || alpha.to != omega.to {
			shared.Logger.Error(
				"Alpha ain't omegaing.",
				"alpha", alpha,
				"omega", omega,
				"desc", desc,
				"asc", asc)
			panic("alpha ain't omegaing")
		}
		width := alpha.to - alpha.from + 1
		height := alpha.y - omega.y + 1
		area := width * height
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Adding area.",
				"width", width,
				"height", height,
				"area", area,
				"alpha", alpha.String(),
				"omega", omega.String())
		}
		total += area

		_, notEmpty, _ := pop()
		if !notEmpty {
			break
		}
	}
	shared.Logger.Info("Digging done.", "volume", total)
	return total
}

func finishRoutes(routes []route, xCoords []int) []route {
	routes = expandAllRoutes(routes, xCoords)
	slices.SortFunc(routes, func(a, b route) int {
		if a.from != b.from {
			return a.from - b.from
		}
		if a.to != b.to {
			return a.to - b.to
		}
		if a.y != b.y {
			return b.y - a.y
		}
		if !a.above && b.above {
			return -1
		}
		return 0
	})
	return slices.Compact(routes)
}

func isClockwise(instructions []instruction) bool {
	var right int
	for i := range instructions {
		right += deriveTurn(instructions, i)
	}
	return right > 0
}

func parseLines(lines []string, magic bool) []instruction {
	if len(lines) == 0 {
		return nil
	}
	instructions := make([]instruction, len(lines))
	for i, each := range lines {
		instructions[i] = parseLine(each, magic)
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

func parseLine(line string, magic bool) instruction {
	pieces := strings.Fields(line)
	if len(pieces) != 3 {
		shared.Logger.Error("Invalid line, piece count != 3.", "count", len(pieces), line, "line")
		panic("invalid line, piece count != 3")
	}
	dir := toDirection(pieces[0])
	var count int
	if magic {
		// Without quotes and #. "(#70c710)" -> "70c710".
		last := pieces[2][2 : len(pieces[2])-1]
		countPart := last[0 : len(last)-1]
		dirPart := int(last[len(last)-1] - '0')
		dir = []shared.Direction{
			shared.RealEast,
			shared.RealSouth,
			shared.RealWest,
			shared.RealNorth,
		}[dirPart]
		count = int(
			gent.OrPanic2(strconv.ParseInt(countPart, 16, 64))("invalid count: " + last),
		)
	} else {
		count = gent.OrPanic2(strconv.Atoi(pieces[1]))("invalid count: " + pieces[1])
	}
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

type route struct {
	from, to int
	y        int
	above    bool
}

func (v route) String() string {
	above := "↓"
	if v.above {
		above = "↑"
	}
	return fmt.Sprintf("%d->%d y:%d %s", v.from, v.to, v.y, above)
}

func gatherRoutes(instructions []instruction, clockwise bool) ([]route, []int) {
	routes := make([]route, 0, len(instructions)/2)
	var loc shared.Loc
	xCoords := gent.NewSet[int]()
	for i := range instructions {
		instr := instructions[i]
		to := loc.Delta(instr.delta)
		if instr.dir.Y == 0 {
			xCoords.Add(to.X)
			r := route{
				from:  min(loc.X, to.X),
				to:    max(loc.X, to.X),
				y:     loc.Y,
				above: isAbove(instr.dir, clockwise),
			}
			routes = append(routes, r)
		} else {
			top := max(loc.Y, to.Y)
			bottom := min(loc.Y, to.Y)
			if top-bottom > 1 {
				r1 := route{from: loc.X, to: loc.X, y: top - 1, above: false}
				r2 := route{from: loc.X, to: loc.X, y: bottom + 1, above: true}
				if shared.IsDebugEnabled() {
					shared.Logger.Debug(
						"Add pipe.",
						"top", r1.String(),
						"bot", r2.String())
				}
				routes = append(routes, r1)
				routes = append(routes, r2)
			}
		}
		loc = to
	}
	xSlice := xCoords.ToSlice()
	slices.Sort(xSlice)
	return routes, xSlice
}

func isAbove(dir shared.Direction, clockwise bool) bool {
	switch dir {
	case shared.RealEast:
		return !clockwise
	case shared.RealWest:
		return clockwise
	default:
		panic("illegal dir for isAbove")
	}
}

func expandAllRoutes(routes []route, xCoords []int) []route {
	expanded := make([]route, 0, len(routes))
	for _, each := range routes {
		expanded = append(expanded, expandRoute(each, xCoords)...)
	}
	return expanded
}

func expandRoute(r route, xCoords []int) []route {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Expand.", "route", r)
	}
	low, bottom := slices.BinarySearch(xCoords, r.from)
	high, ceiling := slices.BinarySearch(xCoords, r.to)
	if !bottom || !ceiling {
		shared.Logger.Error("No bottom or ceiling.", "route", r)
		panic("no bottom or ceiling")
	}
	routes := []route{{from: r.from, to: r.from, y: r.y, above: r.above}}
	for i := low; i < high; i++ {
		from := xCoords[i] + 1
		to := xCoords[i+1] - 1
		if from <= to {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Add from<=to.", "from", from, "to", to)
			}
			routes = append(
				routes,
				route{
					from:  from,
					to:    to,
					y:     r.y,
					above: r.above,
				})
		}
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Add next.")
		}
		routes = append(
			routes,
			route{
				from:  xCoords[i+1],
				to:    xCoords[i+1],
				y:     r.y,
				above: r.above,
			})
	}
	return routes
}

func equalButY(a, b route) bool {
	return a.from == b.from && a.to == b.to && a.above == b.above
}

type iteratorFunc func() (route, bool, func())

func createIterator(routes []route) iteratorFunc {
	var i int
	return func() (route, bool, func()) {
		if i >= len(routes) {
			return route{}, false, nil
		}
		return routes[i], true, func() { i++ }
	}
}

func popRoutes(pop iteratorFunc, above bool) []route {
	var routes []route
	for {
		r, found, resume := pop()
		if !found {
			break
		}
		if r.above == above {
			break
		}
		if len(routes) > 0 && !equalButY(routes[0], r) {
			shared.Logger.Error("Not in line, not equal.", "routes", routes)
			panic("not in line")
		}
		routes = append(routes, r)
		resume()
	}
	return routes
}
