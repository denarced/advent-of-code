package aoc2416

import (
	"fmt"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestCountLowestScore(t *testing.T) {
	run := func(name string, lines []string, expectedScore int, expectedSeats int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			score, seats := CountLowestScore(lines, true)
			require.Equal(t, expectedScore, score)
			require.Equal(t, expectedSeats, seats)
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
		2006,
		7)
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
		7_036,
		45)
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
		3007,
		12)
}

func TestGetPossibleVectors(t *testing.T) {
	all := append([]shared.Direction{}, shared.RealPrimaryDirections...)
	possibleDirections := map[shared.Direction][]shared.Direction{
		shared.RealEast:  getPossibleDirections(all, shared.RealEast),
		shared.RealSouth: getPossibleDirections(all, shared.RealSouth),
		shared.RealWest:  getPossibleDirections(all, shared.RealWest),
		shared.RealNorth: getPossibleDirections(all, shared.RealNorth),
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
			all := append([]shared.Direction{}, shared.RealPrimaryDirections...)
			require.ElementsMatch(t, expected, getPossibleDirections(all, dir))
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

// "0x0 -> 1x1" to proper shared.Loc values.
func parseLocPair(move string) (shared.Loc, shared.Loc) {
	fields := strings.Fields(move)
	from := shared.ParseLoc(fields[0])
	to := shared.ParseLoc(fields[2])
	return from, to
}

func TestSortDirections(t *testing.T) {
	run := func(startToEnd string, expected []shared.Direction) {
		start, end := parseLocPair(startToEnd)
		t.Run(fmt.Sprintf("%s -> %s", start.ToString(), end.ToString()), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			actual := sortDirections(start, end)

			// VERIFY
			req.Equal(4, len(actual))
			preferred := actual[:len(expected)]
			req.ElementsMatch(expected, preferred)
		})
	}

	run("0x0 -> 1x0", []shared.Direction{shared.RealEast})
	run("0x0 -> 1x-1", []shared.Direction{shared.RealEast, shared.RealSouth})
	run("0x0 -> 0x-1", []shared.Direction{shared.RealSouth})
	run("0x0 -> -1x-1", []shared.Direction{shared.RealSouth, shared.RealWest})
	run("0x0 -> -1x0", []shared.Direction{shared.RealWest})
	run("0x0 -> -1x1", []shared.Direction{shared.RealNorth, shared.RealWest})
	run("0x0 -> 0x1", []shared.Direction{shared.RealNorth})
	run("0x0 -> 1x1", []shared.Direction{shared.RealNorth, shared.RealEast})
}
