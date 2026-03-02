package aoc2303

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumPartNumbers(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	sum := SumPartNumbers(lines)
	req.Equal(4361, sum)
}

func TestSumPartNumbersBorders(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines := []string{
		"1/10.2/20",
		".........",
		"4/40.8/80",
	}
	sum := SumPartNumbers(lines)
	req.Equal(165, sum)
}

func TestSumGearRatios(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	// EXERCISE
	sum := SumGearRatios(lines)
	// VERIFY
	req.Equal(467835, sum)
}

func TestIsLocValid(t *testing.T) {
	run := func(name string, loc shared.Loc, width, height int, expected bool) {
		t.Run(name, func(t *testing.T) {
			req := require.New(t)
			req.Equal(expected, isLocValid(loc, width, height))
		})
	}
	run("x:-1", shared.Loc{X: -1}, 1, 1, false)
	run("y:-1", shared.Loc{Y: -1}, 1, 1, false)
	run("0,0", shared.Loc{}, 1, 1, true)
	run("x>width", shared.Loc{X: 2}, 2, 9, false)
	run("y>height", shared.Loc{Y: 2}, 9, 2, false)
}

func TestIsSymbol(t *testing.T) {
	for _, each := range []struct {
		c        rune
		expected bool
	}{
		{'.', false},
		{'0', false},
		{'9', false},
		{' ', true},
		{'#', true},
	} {
		t.Run(string(each.c), func(t *testing.T) {
			require.Equal(t, each.expected, isSymbol(each.c))
		})
	}
}

func TestIsSurroundedBySymbol(t *testing.T) {
	for _, each := range []struct {
		name     string
		expected bool
		lines    []string
		hit      hit
	}{
		{
			name:     "no symbol",
			expected: false,
			lines: []string{
				".....",
				".123.",
				".....",
			},
			hit: hit{
				num:  123,
				from: shared.Loc{X: 1, Y: 1},
				to:   shared.Loc{X: 3, Y: 1},
			},
		},
		{
			name:     "no symbol + squeezed",
			expected: false,
			lines: []string{
				"123",
			},
			hit: hit{
				num:  123,
				from: shared.Loc{X: 0},
				to:   shared.Loc{X: 2},
			},
		},
		{
			name:     "west",
			expected: true,
			lines: []string{
				".....",
				"#123.",
				".....",
			},
			hit: hit{
				num:  123,
				from: shared.Loc{X: 1, Y: 1},
				to:   shared.Loc{X: 3, Y: 1},
			},
		},
		{
			name:     "north-west",
			expected: true,
			lines: []string{
				"*....",
				".123.",
				".....",
			},
			hit: hit{
				num:  123,
				from: shared.Loc{X: 1, Y: 1},
				to:   shared.Loc{X: 3, Y: 1},
			},
		},
		{
			name:     "from-north",
			expected: true,
			lines: []string{
				"./...",
				".123.",
				".....",
			},
			hit: hit{
				num:  123,
				from: shared.Loc{X: 1, Y: 1},
				to:   shared.Loc{X: 3, Y: 1},
			},
		},
		{
			name:     "from-south",
			expected: true,
			lines: []string{
				".....",
				".123.",
				".*...",
			},
			hit: hit{
				num:  123,
				from: shared.Loc{X: 1, Y: 1},
				to:   shared.Loc{X: 3, Y: 1},
			},
		},
		{
			name:     "to-south-east",
			expected: true,
			lines: []string{
				".....",
				".123.",
				"...._",
			},
			hit: hit{
				num:  123,
				from: shared.Loc{X: 1, Y: 1},
				to:   shared.Loc{X: 3, Y: 1},
			},
		},
		{
			name:     "consecutive numbers",
			expected: true,
			lines: []string{
				".......",
				"*12.18.",
				"...._..",
			},
			hit: hit{
				num:  18,
				from: shared.Loc{X: 4, Y: 1},
				to:   shared.Loc{X: 5, Y: 1},
			},
		},
	} {
		t.Run(each.name, func(t *testing.T) {
			require.Equal(t, each.expected, isSurroundedBySymbol(each.lines, each.hit))
		})
	}
}

func TestFeedGears(t *testing.T) {
	lines := []string{
		"*......",
		".....*.",
		"......*",
	}
	var locations []shared.Loc
	cb := func(loc shared.Loc) {
		locations = append(locations, loc)
	}
	// EXERCISE
	feedGears(lines, cb)

	// VERIFY
	require.Equal(
		t,
		[]shared.Loc{
			{},
			{X: 5, Y: 1},
			{X: 6, Y: 2},
		},
		locations,
	)
}

func TestDeriveAdjacentNumbers(t *testing.T) {
	shared.InitTestLogging(t)
	lines := []string{
		"..678..",
		".21*...",
		"..3.57.",
	}
	// EXERCISE
	numbers := deriveAdjacentNumbers(lines, shared.Loc{X: 3, Y: 1})

	// VERIFY
	require.ElementsMatch(t, []int{678, 21, 3, 57}, numbers)
}
