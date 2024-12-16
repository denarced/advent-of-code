package aoc2412

import (
	"fmt"
	"sort"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDerivePrice(t *testing.T) {
	run := func(name string, lines []string, discount bool, expected int) {
		suffix := " without discount"
		if discount {
			suffix = " with discount"
		}
		t.Run(name+suffix, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, DeriveTotalPrice(lines, discount))
		})
	}

	run("empty", []string{}, false, 0)
	run("empty", []string{}, true, 0)
	run("one", []string{"A"}, false, 4)
	run("one", []string{"A"}, true, 4)
	run("2x1", []string{"bb"}, false, 12)
	run(
		"1x1 + T + 1x1",
		[]string{
			"xxx",
			"yxy",
		},
		false,
		4+4+(4*10),
	)

	exampleOne := []string{
		"AAAA",
		"BBCD",
		"BBCC",
		"EEEC",
	}
	run("example 1", exampleOne, false, 140)
	run("example 1", exampleOne, true, 80)

	exampleTwo := []string{
		"OOOOO",
		"OXOXO",
		"OOOOO",
		"OXOXO",
		"OOOOO",
	}
	run("example 2", exampleTwo, false, 772)
	run("example 2", exampleTwo, true, 436)

	exampleE := []string{
		"EEEEE",
		"EXXXX",
		"EEEEE",
		"EXXXX",
		"EEEEE",
	}
	run("example e", exampleE, true, 236)

	run(
		"corners",
		[]string{
			"AB",
			"BA",
		},
		true,
		16)
	run(
		"large example",
		[]string{
			"AAAAAA",
			"AAABBA",
			"AAABBA",
			"ABBAAA",
			"ABBAAA",
			"AAAAAA",
		},
		true,
		// 16 + 16 + 28*(4+8)
		// = 16 + 16 + 336
		// = 368
		4*4+4*4+28*(4+8))
}

func TestCountNeighbours(t *testing.T) {
	run := func(name string, mapped map[shared.Loc]int, loc shared.Loc, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, countNeighbours(mapped, loc))
		})
	}

	run("None", map[shared.Loc]int{{X: 1, Y: 1}: 0}, shared.Loc{X: 1, Y: 1}, 0)
	run(
		"3",
		map[shared.Loc]int{
			{X: 1, Y: 1}: 0,
			{X: 0, Y: 0}: 0,
			{X: 1, Y: 0}: 0,
			{X: 2, Y: 0}: 0,
		},
		shared.Loc{X: 1, Y: 0},
		3,
	)
}

func TestCountFencesSides(t *testing.T) {
	run := func(name string, area []shared.Loc, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, countFenceSides(area))
		})
	}

	run("1x1", []shared.Loc{{X: 1, Y: 1}}, 4)
	run("1x4", []shared.Loc{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}, {X: 2, Y: 5}}, 4)
	run(
		"T",
		toLocs([]string{
			"AAA",
			".A.",
		}, 'A'),
		8)
	run(
		"+",
		toLocs(
			[]string{
				" * ",
				"***",
				" * ",
			},
			'*'),
		12)
	run(
		"O",
		toLocs(
			[]string{
				"+++",
				"+.+",
				"+++",
			},
			'+'),
		8)
}

func toLocs(lines []string, id rune) []shared.Loc {
	brd := shared.NewBoard(lines)
	locs := []shared.Loc{}
	brd.Iter(func(l shared.Loc, c rune) bool {
		if id == c {
			locs = append(locs, l)
		}
		return true
	})
	return locs
}

func makeLoc(s string) shared.Loc {
	return shared.ParseLoc(s)
}

func TestDeriveFatties(t *testing.T) {
	run := func(name string, locs []shared.Loc, expected []fatFence) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			// EXERCISE
			actual := deriveFatties(locs)

			// VERIFY
			ass := assert.New(t)
			// ass.Equal(drawFatties(expected), drawFatties(actual))
			ass.ElementsMatch(stringifyFatties(expected), stringifyFatties(actual))
		})
	}

	run("empty", []shared.Loc{}, []fatFence{})

	{
		fences := []fatFence{
			{f: newFence(makeLoc("1x2"), makeLoc("1x1")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("1x3"), makeLoc("1x2")), loc: makeLoc("1x2")},
			{f: newFence(makeLoc("1x4"), makeLoc("1x3")), loc: makeLoc("1x3")},
			{f: newFence(makeLoc("2x4"), makeLoc("1x4")), loc: makeLoc("1x3")},
			{f: newFence(makeLoc("2x4"), makeLoc("2x3")), loc: makeLoc("1x3")},
			{f: newFence(makeLoc("2x3"), makeLoc("2x2")), loc: makeLoc("1x2")},
			{f: newFence(makeLoc("2x2"), makeLoc("2x1")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("2x1"), makeLoc("1x1")), loc: makeLoc("1x1")},
		}
		run(
			"1x3",
			toLocs(
				[]string{
					"...",
					".v.",
					".v.",
					".v.",
					"...",
				},
				'v'),
			fences,
		)
	}

	{
		fatties := []fatFence{
			// Outer left.
			{f: newFence(makeLoc("0x1"), makeLoc("0x0")), loc: makeLoc("0x0")},
			{f: newFence(makeLoc("0x2"), makeLoc("0x1")), loc: makeLoc("0x1")},
			{f: newFence(makeLoc("0x3"), makeLoc("0x2")), loc: makeLoc("0x2")},

			// Outer top.
			{f: newFence(makeLoc("1x3"), makeLoc("0x3")), loc: makeLoc("0x2")},
			{f: newFence(makeLoc("2x3"), makeLoc("1x3")), loc: makeLoc("1x2")},
			{f: newFence(makeLoc("3x3"), makeLoc("2x3")), loc: makeLoc("2x2")},

			// Outer right.
			{f: newFence(makeLoc("3x2"), makeLoc("3x3")), loc: makeLoc("2x2")},
			{f: newFence(makeLoc("3x1"), makeLoc("3x2")), loc: makeLoc("2x1")},
			{f: newFence(makeLoc("3x0"), makeLoc("3x1")), loc: makeLoc("2x0")},

			// Outer bottom.
			{f: newFence(makeLoc("3x0"), makeLoc("2x0")), loc: makeLoc("2x0")},
			{f: newFence(makeLoc("2x0"), makeLoc("1x0")), loc: makeLoc("1x0")},
			{f: newFence(makeLoc("1x0"), makeLoc("0x0")), loc: makeLoc("0x0")},

			// Inner.
			{f: newFence(makeLoc("1x1"), makeLoc("1x2")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("1x2"), makeLoc("2x2")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("2x2"), makeLoc("2x1")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("2x1"), makeLoc("1x1")), loc: makeLoc("1x1")},
		}
		run(
			"O",
			toLocs(
				[]string{
					"xxx",
					"x.x",
					"xxx",
				},
				'x'),
			fatties)
	}
}

func stringifyFatties(fatties []fatFence) []string {
	strs := make([]string, 0, len(fatties))
	for _, each := range fatties {
		each = normalizeFatFence(each)
		strs = append(
			strs,
			fmt.Sprintf(
				"%s -> %s",
				each.f.First.ToString(),
				each.f.Second.ToString()))
	}
	return strs
}

func normalizeFatFence(f fatFence) fatFence {
	rev := newFatFence(f.loc, f.f.Second, f.f.First)
	if f.f.First.X == f.f.Second.X {
		if f.f.First.Y < f.f.Second.Y {
			return f
		}
		return rev
	}
	if f.f.First.X < f.f.Second.X {
		return f
	}
	return rev
}

func TestFindArea(t *testing.T) {
	run := func(name string, lines []string, expected [][]shared.Loc) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			brd := shared.NewBoard(lines)
			compareLocationTables(
				t,
				expected,
				findAreas(brd))
		})
	}

	run(
		"4 * 1x1",
		[]string{
			"ab",
			"ba",
		},
		[][]shared.Loc{
			{{X: 0, Y: 0}},
			{{X: 1, Y: 0}},
			{{X: 0, Y: 1}},
			{{X: 1, Y: 1}},
		})
	run(
		"O",
		[]string{
			"eee",
			"ege",
			"eee",
		},
		[][]shared.Loc{
			toLocs(
				[]string{
					"vvv",
					"v.v",
					"vvv",
				},
				'v'),
			toLocs(
				[]string{
					"...",
					".a.",
					"...",
				},
				'a'),
		})
}

func compareLocationTables(t *testing.T, a, b [][]shared.Loc) {
	stringify := func(table [][]shared.Loc) [][]string {
		s := [][]string{}
		for _, each := range table {
			line := shared.MapValues(
				each,
				func(l shared.Loc) string {
					return l.ToString()
				})
			sort.Strings(line)
			s = append(s, line)
		}
		return s
	}
	require.ElementsMatch(
		t,
		stringify(a),
		stringify(b))
}

func TestSortFences(t *testing.T) {
	run := func(name string, unsorted []fatFence, expected [][]fatFence) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			// EXERCISE
			actual := sortFences(unsorted)
			// VERIFY
			req.ElementsMatch(
				shared.MapValues(expected, stringifyFatties),
				shared.MapValues(actual, stringifyFatties),
			)
		})
	}

	run(
		"1x1",
		[]fatFence{
			{f: newFence(makeLoc("0x0"), makeLoc("0x1")), loc: makeLoc("0x0")},
			{f: newFence(makeLoc("0x1"), makeLoc("1x1")), loc: makeLoc("0x0")},
			{f: newFence(makeLoc("1x1"), makeLoc("1x0")), loc: makeLoc("0x0")},
			{f: newFence(makeLoc("1x0"), makeLoc("0x0")), loc: makeLoc("0x0")},
		},
		[][]fatFence{
			{
				{f: newFence(makeLoc("0x0"), makeLoc("0x1")), loc: makeLoc("0x0")},
				{f: newFence(makeLoc("0x1"), makeLoc("1x1")), loc: makeLoc("0x0")},
				{f: newFence(makeLoc("1x1"), makeLoc("1x0")), loc: makeLoc("0x0")},
				{f: newFence(makeLoc("1x0"), makeLoc("0x0")), loc: makeLoc("0x0")},
			},
		})
	{
		outer := []fatFence{
			{f: newFence(makeLoc("0x0"), makeLoc("0x1")), loc: makeLoc("0x0")},
			{f: newFence(makeLoc("0x1"), makeLoc("0x2")), loc: makeLoc("0x1")},
			{f: newFence(makeLoc("0x2"), makeLoc("0x3")), loc: makeLoc("0x2")},

			{f: newFence(makeLoc("0x3"), makeLoc("1x3")), loc: makeLoc("0x2")},
			{f: newFence(makeLoc("1x3"), makeLoc("2x3")), loc: makeLoc("1x2")},
			{f: newFence(makeLoc("2x3"), makeLoc("3x3")), loc: makeLoc("2x2")},

			{f: newFence(makeLoc("3x3"), makeLoc("3x2")), loc: makeLoc("2x2")},
			{f: newFence(makeLoc("3x2"), makeLoc("3x1")), loc: makeLoc("2x1")},
			{f: newFence(makeLoc("3x1"), makeLoc("3x0")), loc: makeLoc("2x0")},

			{f: newFence(makeLoc("2x0"), makeLoc("3x0")), loc: makeLoc("2x0")},
			{f: newFence(makeLoc("1x0"), makeLoc("2x0")), loc: makeLoc("1x0")},
			{f: newFence(makeLoc("0x0"), makeLoc("1x0")), loc: makeLoc("0x0")},
		}
		inner := []fatFence{
			{f: newFence(makeLoc("1x1"), makeLoc("1x2")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("1x2"), makeLoc("2x2")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("2x2"), makeLoc("2x1")), loc: makeLoc("1x1")},
			{f: newFence(makeLoc("2x1"), makeLoc("1x1")), loc: makeLoc("1x1")},
		}
		combined := []fatFence{}
		combined = append(combined, outer...)
		combined = append(combined, inner...)
		run(
			"O",
			combined,
			[][]fatFence{
				inner,
				outer,
			})
	}
}
