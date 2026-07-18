package aoc2324

import (
	"math/big"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountIntersections(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	// EXERCISE & VERIFY
	req.Equal(2, CountIntersections(lines, 7, 27))
}

func TestToRay(t *testing.T) {
	run := func(stone hailstone, expected ray) {
		t.Run(stone.String(), func(t *testing.T) {
			shared.InitTestLogging(t)
			ass := assert.New(t)

			// EXERCISE
			actual := toRay(stone)

			// VERIFY
			ass.Equalf(
				0,
				expected.a.Cmp(actual.a),
				"a: expected %s, got %s",
				expected.a.RatString(),
				actual.a.RatString())
			ass.Equalf(
				0,
				expected.c.Cmp(actual.c),
				"b: expected %s, got %s",
				expected.c.RatString(),
				actual.c.RatString())
		})
	}

	newRatInt := func(i int) *big.Rat {
		return big.NewRat(int64(i), 1)
	}
	newRatFlt := func(a, b int) *big.Rat {
		return big.NewRat(int64(a), int64(b))
	}

	tests := []struct {
		px, py int
		vx, vy int
		a, b   *big.Rat
	}{
		// y = x through origo
		{-1, -1, 3, 3, newRatInt(1), newRatInt(0)},
		// y = x through 0,-1
		{-1, -2, 3, 3, newRatInt(1), newRatInt(-1)},
		// same as above but different expression for velocity
		{-1, -2, -2, -2, newRatInt(1), newRatInt(-1)},
		// y = (2/3)x + 1/3
		// y(1) = 2/3 + 1/3 = 1
		// y(-2) = (2/3)*-2 + 1/3
		// y(-2) = -4/3 + 1/3 = -1
		{1, 1, 3, 2, newRatFlt(2, 3), newRatFlt(1, 3)},
		// downward, expressed as up and left
		{3, -17, -1, 8, newRatInt(-8), newRatInt(7)},
		// right and up, expressed as left and down
		{-10, -1, -1, -2, newRatInt(2), newRatInt(19)},
	}
	for _, each := range tests {
		run(
			hailstone{
				position: RatCoordinate{
					big.NewRat(int64(each.px), 1),
					big.NewRat(int64(each.py), 1),
				},
				velocity: RatCoordinate{
					big.NewRat(int64(each.vx), 1),
					big.NewRat(int64(each.vy), 1),
				},
			},
			ray{a: each.a, c: each.b})
	}
}

func TestDeriveIntersection(t *testing.T) {
	run := func(name string, a, b ray, expected *RatCoordinate) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			ass := assert.New(t)

			// EXERCISE
			actual := deriveIntersection(a, b)

			// VERIFY
			if expected == nil {
				ass.Nil(actual, "expected nil coordinate")
				return
			}
			ass.Equalf(
				0,
				expected[0].Cmp(actual[0]),
				"X: expected %s, got %s",
				expected[0].RatString(),
				actual[0].RatString())
			ass.Equalf(
				0,
				expected[1].Cmp(actual[1]),
				"Y: expected %s, got %s",
				expected[1].RatString(),
				actual[1].RatString())
		})
	}

	toRat := func(x, y int) *RatCoordinate {
		return &RatCoordinate{
			big.NewRat(int64(x), 1),
			big.NewRat(int64(y), 1),
			nil,
		}
	}
	run(
		"parallel",
		ray{a: big.NewRat(1, 4), c: new(big.Rat)},
		ray{a: big.NewRat(1, 4), c: new(big.Rat)},
		nil)
	run(
		"pyramid top-right from origo",
		ray{a: big.NewRat(1, 1), c: new(big.Rat)},
		ray{a: big.NewRat(-1, 1), c: big.NewRat(6, 1)},
		toRat(3, 3))
	run(
		"a bit more complex",
		ray{a: big.NewRat(2, 1), c: big.NewRat(1, 1)},
		ray{a: big.NewRat(-1, 1), c: big.NewRat(4, 1)},
		toRat(1, 3))
}

func TestOverlap(t *testing.T) {
	run := func(name string, first, second *RatSegment, expected bool) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			actual := overlap(first, second)

			req.Equal(expected, actual)
		})
	}

	one := big.NewRat(1, 1)
	run(
		"a before b on x, space between",
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(2, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(4, 1), one, nil},
		},
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(6, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(8, 1), one, nil},
		},
		false)
	run(
		"a before b on x, half overlap",
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(2, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(4, 1), one, nil},
		},
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(3, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(5, 1), one, nil},
		},
		true)
	run(
		"a after b on x, half overlap",
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(3, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(5, 1), one, nil},
		},
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(2, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(4, 1), one, nil},
		},
		true)
	run(
		"a after b on x, space between",
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(5, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(7, 1), one, nil},
		},
		&RatSegment{
			minimum: [3]*big.Rat{big.NewRat(2, 1), one, nil},
			maximum: [3]*big.Rat{big.NewRat(4, 1), one, nil},
		},
		false)
	run(
		"a before b on y, space between",
		&RatSegment{
			minimum: [3]*big.Rat{one, big.NewRat(2, 1), nil},
			maximum: [3]*big.Rat{one, big.NewRat(4, 1), nil},
		},
		&RatSegment{
			minimum: [3]*big.Rat{one, big.NewRat(5, 1), nil},
			maximum: [3]*big.Rat{one, big.NewRat(7, 1), nil},
		},
		false)
	run(
		"a before b on y, half overlap",
		&RatSegment{
			minimum: [3]*big.Rat{one, big.NewRat(2, 1), nil},
			maximum: [3]*big.Rat{one, big.NewRat(4, 1), nil},
		},
		&RatSegment{
			minimum: [3]*big.Rat{one, big.NewRat(3, 1), nil},
			maximum: [3]*big.Rat{one, big.NewRat(5, 1), nil},
		},
		true)
	run(
		"a after b on y, space between",
		&RatSegment{
			minimum: [3]*big.Rat{one, big.NewRat(6, 1), nil},
			maximum: [3]*big.Rat{one, big.NewRat(8, 1), nil},
		},
		&RatSegment{
			minimum: [3]*big.Rat{one, big.NewRat(2, 1), nil},
			maximum: [3]*big.Rat{one, big.NewRat(4, 1), nil},
		},
		false)
}

func TestIsIn(t *testing.T) {
	c := map[int]*big.Rat{
		1: big.NewRat(1, 1),
		2: big.NewRat(2, 1),
		3: big.NewRat(3, 1),
		4: big.NewRat(4, 1),
		5: big.NewRat(5, 1),
	}
	segment := &RatSegment{
		minimum: RatCoordinate{c[2], c[2], nil},
		maximum: RatCoordinate{c[4], c[4], nil},
	}
	run := func(name string, coordinate *RatCoordinate, expected bool) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE & VERIFY
			req.Equal(expected, isIn(segment, coordinate))
		})
	}

	run("on left outside", &RatCoordinate{c[1], c[3], nil}, false)
	run("on left border", &RatCoordinate{c[2], c[3], nil}, true)
	run("middle", &RatCoordinate{c[3], c[3], nil}, true)
	run("above outside", &RatCoordinate{c[3], c[5], nil}, false)
	run("on top border", &RatCoordinate{c[3], c[4], nil}, true)
	run("on right border", &RatCoordinate{c[4], c[3], nil}, true)
	run("below outside", &RatCoordinate{c[3], c[1], nil}, false)
	run("on bottom border", &RatCoordinate{c[3], c[2], nil}, true)
}
