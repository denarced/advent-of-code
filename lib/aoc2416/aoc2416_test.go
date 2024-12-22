package aoc2416

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestCountLowestScore(t *testing.T) {
	run := func(name string, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, CountLowestScore(lines, false))
		})
	}

	run(
		"small",
		[]string{
			"#####",
			"#S..#",
			"###.#",
			"#E..#",
			"#####",
		},
		2006)
	run(
		"example",
		[]string{
			"###############",
			"#.......#....E#",
			"#.#.###.#.###.#",
			"#.....#.#...#.#",
			"#.###.#####.#.#",
			"#.#.#.......#.#",
			"#.#.#####.###.#",
			"#...........#.#",
			"###.#.#####.#.#",
			"#...#.....#.#.#",
			"#.#.#.###.#.#.#",
			"#.....#...#.#.#",
			"#.###.#.#.#.#.#",
			"#S..#.....#...#",
			"###############",
		},
		7_036)
	run(
		"small branching",
		[]string{
			/* 5 */ "########",
			/* 4 */ "###.####",
			/* 3 */ "###....#",
			/* 2 */ "#S..##E#",
			/* 1 */ "###....#",
			/* 0 */ "########",
			//       01234567
		},
		3007)
}

func TestGetPossibleVectors(t *testing.T) {
	possibleDirections := map[shared.Direction][]shared.Direction{
		shared.RealEast:  getPossibleDirections(shared.RealEast),
		shared.RealSouth: getPossibleDirections(shared.RealSouth),
		shared.RealWest:  getPossibleDirections(shared.RealWest),
		shared.RealNorth: getPossibleDirections(shared.RealNorth),
	}
	run := func(name string, lines []string, dir shared.Direction, expected []vector) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)

			brd := shared.NewBoard(lines)
			start := brd.FindOrDie('S')
			vec := vector{loc: start, dir: dir}
			// EXERCISE
			vectors := getPossibleVectors(brd, vec, possibleDirections)

			// VERIFY
			require.ElementsMatch(t, expected, vectors)
		})
	}

	run(
		"only forward",
		[]string{
			"##",
			"S.",
			"##",
		},
		shared.RealEast,
		[]vector{{loc: shared.Loc{X: 0, Y: 1}, dir: shared.RealEast}},
	)
	run(
		"left or right",
		[]string{
			"###",
			".S.",
			"#.#",
		},
		shared.RealNorth,
		[]vector{
			{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealWest},
			{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealEast},
		},
	)
}

func TestGetPossibleDirections(t *testing.T) {
	run := func(dir shared.Direction, expected []shared.Direction) {
		t.Run(fmt.Sprintf("%v", dir), func(t *testing.T) {
			shared.InitTestLogging(t)
			require.ElementsMatch(t, expected, getPossibleDirections(dir))
		})
	}

	run(shared.RealEast, []shared.Direction{shared.RealNorth, shared.RealEast, shared.RealSouth})
	run(shared.RealSouth, []shared.Direction{shared.RealWest, shared.RealEast, shared.RealSouth})
	run(shared.RealWest, []shared.Direction{shared.RealWest, shared.RealSouth, shared.RealNorth})
	run(shared.RealNorth, []shared.Direction{shared.RealWest, shared.RealEast, shared.RealNorth})
}

func TestDerivePoints(t *testing.T) {
	run := func(previous, current shared.Direction, expected int) {
		t.Run(fmt.Sprintf("%v -> %v", previous, current), func(t *testing.T) {
			shared.InitTestLogging(t)
			actual := derivePoints(previous, current)
			require.Equal(t, expected, actual)
		})
	}

	run(shared.RealNorth, shared.RealNorth, pointsStep)
	run(shared.RealNorth, shared.RealEast, pointsTurn+pointsStep)
}
