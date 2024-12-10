package aoc2408

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestCountUniqueAntinodeLocations(t *testing.T) {
	run := func(name string, lines []string, resonantHarmonics bool, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(
				t,
				expected,
				CountUniqueAntinodeLocations(lines, resonantHarmonics),
			)
		})
	}

	run("empty wo harmonics", []string{}, false, 0)
	run("empty with harmonics", []string{}, true, 0)
	run("example part 1", getExampleLines(), false, 14)
	run("example part 2", getExampleLines(), true, 34)
}

func getExampleLines() []string {
	return shared.StripPadding([]string{
		/* 11 */ ". . . . . . . . . . . .",
		/* 10 */ ". . . . . . . . 0 . . .",
		/*  9 */ ". . . . . 0 . . . . . .",
		/*  8 */ ". . . . . . . 0 . . . .",
		/*  7 */ ". . . . 0 . . . . . . .",
		/*  6 */ ". . . . . . A . . . . .",
		/*  5 */ ". . . . . . . . . . . .",
		/*  4 */ ". . . . . . . . . . . .",
		/*  3 */ ". . . . . . . . A . . .",
		/*  2 */ ". . . . . . . . . A . .",
		/*  1 */ ". . . . . . . . . . . .",
		/*  0 */ ". . . . . . . . . . . .",
		//        0 1 2 3 4 5 6 7 8 9 a b
	})
}

func TestDeriveAntinodes(t *testing.T) {
	run := func(
		name string,
		a,
		b shared.Loc,
		resonantHarmonics bool,
		expected []shared.Loc) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.ElementsMatch(t, expected, deriveAntinodes(a, b, 10, 10, resonantHarmonics))
		})
	}

	run(
		"0x0 1x1",
		shared.ParseLoc("0x0"),
		shared.ParseLoc("1x1"),
		false,
		[]shared.Loc{{X: 2, Y: 2}},
	)
	run(
		"6x4 8x5",
		shared.Loc{X: 6, Y: 4},
		shared.Loc{X: 4, Y: 5},
		false,
		[]shared.Loc{
			{X: 2, Y: 6},
			{X: 8, Y: 3},
		})
	run(
		"3x4 3x5",
		shared.Loc{X: 3, Y: 4},
		shared.Loc{X: 3, Y: 5},
		true,
		[]shared.Loc{
			{X: 3, Y: 0},
			{X: 3, Y: 1},
			{X: 3, Y: 2},
			{X: 3, Y: 3},
			{X: 3, Y: 4},
			{X: 3, Y: 5},
			{X: 3, Y: 6},
			{X: 3, Y: 7},
			{X: 3, Y: 8},
			{X: 3, Y: 9},
			{X: 3, Y: 10},
		})
}

func TestDeriveUniqueAntiNodeLocations(t *testing.T) {
	run := func(
		name string,
		lines []string,
		harmonics bool,
		expected *shared.Set[shared.Loc]) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			shared.DiffLocSets(t, expected, deriveUniqueAntinodeLocations(lines, harmonics))
		})
	}

	run("empty wo harmonics", []string{}, false, shared.NewSet([]shared.Loc{}))
	run(
		"minimal with harmonics",
		shared.StripPadding([]string{
			/* 9 */ ". . . . . . . . . .",
			/* 8 */ ". . . . . . . . . .",
			/* 7 */ ". . . . . . . . . .",
			/* 6 */ ". . . . . . . . . .",
			/* 5 */ ". . . . . . . . . .",
			/* 4 */ ". . 0 . . . . . 2 .",
			/* 3 */ ". 3 . . . . . . . .",
			/* 2 */ "3 . 0 1 . . . 2 . .",
			/* 1 */ ". . . . . . . . . .",
			/* 0 */ "0 . 1 . . . . . . .",
			//       0 1 2 3 4 5 6 7 8 9
		}),
		true,
		shared.NewSet([]shared.Loc{
			{X: 0, Y: 0}, // 0
			{X: 0, Y: 2}, // 3
			{X: 1, Y: 3}, // 3
			{X: 2, Y: 0}, // 0
			{X: 2, Y: 2}, // 0
			{X: 2, Y: 4}, // 3
			{X: 2, Y: 6}, // 0
			{X: 2, Y: 8}, // 0
			{X: 3, Y: 2}, // 3
			{X: 3, Y: 5}, // 3
			{X: 4, Y: 4}, // 0, 1
			{X: 4, Y: 6}, // 3
			{X: 4, Y: 8}, // 0
			{X: 5, Y: 6}, // 1
			{X: 5, Y: 7}, // 3
			{X: 6, Y: 0}, // 2
			{X: 6, Y: 6}, // 0
			{X: 6, Y: 8}, // 1, 3
			{X: 7, Y: 2}, // 2
			{X: 7, Y: 9}, // 3
			{X: 8, Y: 4}, // 2
			{X: 8, Y: 8}, // 0
			{X: 9, Y: 6}, // 2
		}))
	// b # # . . . . # . . . . #
	// a . # . # . . . . 0 . . .
	// 9 . . # . # 0 . . . . # .
	// 8 . . # # . . . 0 . . . .
	// 7 . . . . 0 . . . . # . .
	// 6 . # . . . # A . . . . #
	// 5 . . . # . . # . . . . .
	// 4 # . . . . # . # . . . .
	// 3 . . # . . . . . A . . .
	// 2 . . . . # . . . . A . .
	// 1 . # . . . . . . . . # .
	// 0 . . . # . . . . . . # #
	//   0 1 2 3 4 5 6 7 8 9 a b
	run(
		"example",
		getExampleLines(),
		true,
		shared.NewSet([]shared.Loc{
			{X: 0, Y: 11},
			{X: 0, Y: 4},
			{X: 1, Y: 10},
			{X: 1, Y: 11},
			{X: 1, Y: 1},
			{X: 1, Y: 6},
			{X: 10, Y: 0},
			{X: 10, Y: 1},
			{X: 10, Y: 9},
			{X: 11, Y: 0},
			{X: 11, Y: 11},
			{X: 11, Y: 6},
			{X: 2, Y: 3},
			{X: 2, Y: 8},
			{X: 2, Y: 9},
			{X: 3, Y: 0},
			{X: 3, Y: 10},
			{X: 3, Y: 5},
			{X: 3, Y: 8},
			{X: 4, Y: 2},
			{X: 4, Y: 7},
			{X: 4, Y: 9},
			{X: 5, Y: 4},
			{X: 5, Y: 6},
			{X: 5, Y: 9},
			{X: 6, Y: 11},
			{X: 6, Y: 5},
			{X: 6, Y: 6},
			{X: 7, Y: 4},
			{X: 7, Y: 8},
			{X: 8, Y: 10},
			{X: 8, Y: 3},
			{X: 9, Y: 2},
			{X: 9, Y: 7},
		}))
}

func TestIsAntenna(t *testing.T) {
	run := func(c rune, expected bool) {
		t.Run(string(c), func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, isAntenna(c))
		})
	}

	run('0', true)
	run('9', true)
	run('.', false)
	run('a', true)
	run('z', true)
	run('|', false)
	run('A', true)
	run('Z', true)
}
