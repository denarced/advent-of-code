package aoc2324

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func CountIntersections(lines []string, from, to int64) int {
	shared.Logger.Info("Count intersections.", "from", from, "to", to)
	hailstones := parseLines(lines)
	rays := toRays(hailstones)
	var count int
	testArea := &RatSegment{
		minimum: RatCoordinate{big.NewRat(from, 1), big.NewRat(from, 1), nil},
		maximum: RatCoordinate{big.NewRat(to, 1), big.NewRat(to, 1), nil},
	}
	for i := 0; i < len(rays)-1; i++ {
		for j := i + 1; j < len(rays); j++ {
			if intersect(rays, testArea, i, j) {
				count++
			}
		}
	}
	shared.Logger.Info("Intersections counted.", "total", count)
	return count
}

func intersect(rays []ray, segment *RatSegment, i, j int) bool {
	coord := deriveIntersection(rays[i], rays[j])
	if coord == nil {
		return false
	}
	for _, i := range []int{i, j} {
		if !isInFuture(coord, rays[i]) {
			return false
		}
	}
	if isIn(segment, coord) {
		shared.Logger.Debug("Matches.", "i", i, "j", j)
		return true
	}
	return false
}

func toRays(stones []hailstone) []ray {
	rays := make([]ray, len(stones))
	for i, each := range stones {
		rays[i] = toRay(each)
	}
	return rays
}

func toRay(stone hailstone) ray {
	n := new(big.Rat).Quo(stone.position[0], stone.velocity[0])
	n.Mul(n, stone.velocity[1])
	n.Neg(n)
	c := new(big.Rat).Add(stone.position[1], n)
	r := ray{
		a:        new(big.Rat).Quo(stone.velocity[1], stone.velocity[0]),
		c:        c,
		position: stone.position,
		velocity: stone.velocity,
	}
	return r
}

func parseLines(lines []string) []hailstone {
	stones := make([]hailstone, len(lines))
	for i, each := range lines {
		stones[i] = parseLine(each)
	}
	return stones
}

func parseLine(line string) hailstone {
	linePieces := gent.Map(strings.Split(line, "@"), strings.TrimSpace)
	if len(linePieces) != 2 {
		shared.Logger.Error("Invalid line: it should have one @ character.", "line", line)
		panic("invalid line: should have one @ character")
	}
	return hailstone{
		position: parseCoordinate(linePieces[0]),
		velocity: parseCoordinate(linePieces[1]),
	}
}

func parseCoordinate(s string) RatCoordinate {
	values := gent.Map(gent.Map(strings.Split(s, ","), strings.TrimSpace), func(s string) int64 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			shared.Logger.Error("Failed to convert to int.", "err", err)
			panic(err)
		}
		return i
	})
	if len(values) != 3 {
		shared.Logger.Error("Should have 3 values for a coordinate.", "values", values)
		panic("invalid coordinate")
	}
	return RatCoordinate{
		big.NewRat(values[0], 1),
		big.NewRat(values[1], 1),
		big.NewRat(values[2], 1),
	}
}

type hailstone struct {
	position RatCoordinate
	velocity RatCoordinate
}

func (v hailstone) String() string {
	return fmt.Sprintf("%s-%s", v.position.String(), v.velocity.String())
}

type RatCoordinate [3]*big.Rat

func (v *RatCoordinate) String() string {
	if v == nil {
		return "nil"
	}
	pieces := make([]string, len(v))
	for i, each := range v {
		if each == nil {
			pieces[i] = "nil"
		} else {
			pieces[i] = each.FloatString(1)
		}
	}
	return strings.Join(pieces, ",")
}

type ray struct {
	// As in "f(x) = ax^2 + bx + c".
	a, c     *big.Rat
	position RatCoordinate
	velocity RatCoordinate
}

func (v ray) String() string {
	return fmt.Sprintf("%sx+%s", v.a.FloatString(1), v.c.FloatString(1))
}

func deriveIntersection(a, b ray) *RatCoordinate {
	if a.a.Cmp(b.a) == 0 {
		return nil
	}
	multiplier := new(big.Rat).Add(a.a, new(big.Rat).Neg(b.a))
	x := new(big.Rat).Add(b.c, new(big.Rat).Neg(a.c))
	x.Quo(x, multiplier)
	y := new(big.Rat).Mul(x, a.a)
	y.Add(y, a.c)
	return &RatCoordinate{x, y, nil}
}

type RatSegment struct {
	minimum, maximum RatCoordinate
}

func (v *RatSegment) String() string {
	x := fmt.Sprintf("%s-%s", v.minimum[0].FloatString(1), v.maximum[0].FloatString(1))
	y := fmt.Sprintf("%s-%s", v.minimum[1].FloatString(1), v.maximum[1].FloatString(1))
	return fmt.Sprintf("%s,%s", x, y)
}

func overlap(a, b *RatSegment) bool {
	for i := range 2 {
		k := a.minimum[i]
		l := a.maximum[i]
		m := b.minimum[i]
		n := b.maximum[i]
		if k.Cmp(m) < 0 && l.Cmp(m) < 0 {
			return false
		}
		if m.Cmp(k) < 0 && n.Cmp(k) < 0 {
			return false
		}
	}
	return true
}

func isIn(segment *RatSegment, coord *RatCoordinate) bool {
	for i := range 2 {
		if coord[i].Cmp(segment.minimum[i]) < 0 {
			return false
		}
		if coord[i].Cmp(segment.maximum[i]) > 0 {
			return false
		}
	}
	return true
}

func isInFuture(coord *RatCoordinate, aRay ray) (inFuture bool) {
	if aRay.velocity[0].Sign() >= 0 {
		inFuture = coord[0].Cmp(aRay.position[0]) >= 0
	} else {
		inFuture = coord[0].Cmp(aRay.position[0]) < 0
	}
	return
}
