package aoc2509

import (
	"fmt"
	"slices"

	"github.com/denarced/advent-of-code/shared"
)

func calculateArea(a, b [2]int) int {
	w := abs(a[0]-b[0]) + 1
	h := abs(a[1]-b[1]) + 1
	return w * h
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func toCoords(lines []string) ([][2]int, error) {
	coords := make([][2]int, len(lines))
	for i, each := range lines {
		var x, y int
		_, err := fmt.Sscanf(each, "%d,%d", &x, &y)
		if err != nil {
			return nil, err
		}
		point := [2]int{x, y}
		coords[i] = point
	}
	return coords, nil
}

func deriveClockwise(coords [][2]int) bool {
	size := len(coords)
	if size < 4 {
		panic(fmt.Sprintf("not enough coordinates: %d", size))
	}
	turns := 0
	for i := 0; i < size-2; i++ {
		start := coords[i]
		mid := coords[i+1]
		end := coords[0]
		if i+2 < size {
			end = coords[i+2]
		}
		turns += deriveTurn(start, mid, end)
	}
	return turns > 0
}

// DeriveTurn return -1 for left turn, 1 for right.
func deriveTurn(start, mid, end [2]int) int {
	if start[0] == mid[0] {
		if start[1] < mid[1] {
			// to top
			if end[0] < mid[0] {
				// up-left
				return -1
			}
			// up-right
			return 1
		}
		if end[0] < mid[0] {
			// bottom-left
			return 1
		}
		// bottom-right
		return -1
	}
	if start[0] < mid[0] {
		if end[1] > mid[1] {
			// right-up
			return -1
		}
		// right-down
		return 1
	}
	if end[1] > mid[1] {
		// left-up
		return 1
	}
	// left-down
	return -1
}

func screen(s [][2]int) [][2][2]int {
	pairs := make([][2][2]int, len(s))
	for i := range len(s) {
		start := s[i]
		end := s[(i+1)%len(s)]
		pairs[i] = [2][2]int{
			start,
			end,
		}
	}
	return pairs
}

func fillLine(start, end [2]int) [][2]int {
	shared.Logger.Debug("Fill line.", "start", start, "end", end)
	var line [][2]int
	xd := 0
	if end[0] > start[0] {
		xd = 1
	} else if end[0] < start[0] {
		xd = -1
	}
	yd := 0
	if end[1] > start[1] {
		yd = 1
	} else if end[1] < start[1] {
		yd = -1
	}
	x := start[0]
	y := start[1]
	for x != end[0] || y != end[1] {
		line = append(line, [2]int{x, y})
		x += xd
		y += yd
	}
	line = append(line, end)
	shared.Logger.Debug("Line filled.", "line", line)
	return line
}

func DeriveBiggestRectangle(lines []string, redGreen bool) int {
	corners, err := toCoords(lines)
	if err != nil {
		shared.Logger.Error(
			"Failed to convert lines to coordinates.",
			"line count",
			len(lines),
			"err",
			err,
		)
		panic(err)
	}
	return deriveBiggest(corners, redGreen)
}

func modulo(i, n int) int {
	for i < 0 {
		i += n
	}
	return i % n
}

type corner struct {
	prev, curr, next int
}

func screenCorner[t any](s []t, forward bool) func() (corner, bool) {
	if len(s) < 3 {
		panic("len(s) must be >= 3")
	}
	counter := len(s)
	i := 0
	add := 1
	if !forward {
		i = len(s) - 1
		add = -1
	}
	return func() (corner, bool) {
		var corn corner
		if counter == 0 {
			return corn, true
		}
		corn.curr = i
		corn.next = modulo(i+add, len(s))
		corn.prev = modulo(i-add, len(s))
		i += add
		counter--
		return corn, false
	}
}

func findCandidates(coords [][2]int, corn corner) []int {
	prev := coords[corn.prev]
	curr := coords[corn.curr]
	next := coords[corn.next]
	filters := []func([2]int) bool{
		createCoordinateFilter(prev, curr),
		createCoordinateFilter(curr, next),
	}
	matchBothFilters := func(point [2]int) bool {
		for _, each := range filters {
			if !each(point) {
				return false
			}
		}
		return true
	}
	var candidates []int
	for i, each := range coords {
		if i == corn.curr {
			continue
		}
		if matchBothFilters(each) {
			candidates = append(candidates, i)
		}
	}
	return candidates
}

func createCoordinateFilter(from, to [2]int) func([2]int) bool {
	if from[0] == to[0] {
		if from[1] < to[1] {
			// up
			return func(point [2]int) bool {
				return point[0] >= from[0]
			}
		}
		// down
		return func(point [2]int) bool {
			return point[0] <= from[0]
		}
	}
	if from[0] < to[0] {
		// right
		return func(point [2]int) bool {
			return point[1] <= from[1]
		}
	}
	// left
	return func(point [2]int) bool {
		return point[1] >= from[1]
	}
}

func deriveBiggest(coords [][2]int, redGreen bool) int {
	it := screenCorner(coords, deriveClockwise(coords))
	corn, done := it()
	maximum := -1
	for !done {
		candidates := findCandidates(coords, corn)
		for _, candIndex := range candidates {
			candidate := coords[candIndex]
			area := calculateArea(coords[corn.curr], candidate)
			if area <= maximum {
				continue
			}
			if redGreen && fallsIn(coords, corn.curr, candIndex) {
				shared.Logger.Debug("Fell in.", "area", area)
				continue
			}
			maximum = area
		}
		corn, done = it()
	}
	return maximum
}

func fallsIn(coords [][2]int, fromIndex, toIndex int) bool {
	shared.Logger.Debug("Falls in?", "from", fromIndex, "to", toIndex)
	from, to := coords[fromIndex], coords[toIndex]
	right := max(from[0], to[0])
	bottom := min(from[1], to[1])
	left := min(from[0], to[0])
	top := max(from[1], to[1])
	isIn := func(p [2]int) bool {
		return left < p[0] && p[0] < right && bottom < p[1] && p[1] < top
	}
	for i, each := range coords {
		if i == fromIndex || i == toIndex {
			continue
		}
		if isIn(each) {
			shared.Logger.Debug(
				"Failed.",
				"right", right,
				"bottom", bottom,
				"left", left,
				"top", top,
				"each", each)
			return true
		}
	}
	for _, pair := range screen(coords) {
		if slices.ContainsFunc(fillLine(pair[0], pair[1]), isIn) {
			return true
		}
	}
	return false
}
