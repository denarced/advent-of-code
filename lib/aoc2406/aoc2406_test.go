package aoc2406

import (
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestToInts(t *testing.T) {
	run := func(name string, s []string, expected []int, errMessage string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual, err := shared.ToInts(s)
			if errMessage == "" {
				req.Nil(err)
				req.Equal(expected, actual)
			} else {
				req.ErrorContains(err, errMessage)
			}
		})
	}

	run("empty", []string{}, nil, "")
	run("happy path", []string{"-1", "0", "1"}, []int{-1, 0, 1}, "")
	run("failure", []string{"e"}, nil, "invalid syntax")
}

func TestCountDistinctPositions(t *testing.T) {
	run := func(name string, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := CountDistinctPositions(lines)
			req.Equal(expected, actual)
		})
	}

	run("empty", []string{}, 0)
	run("example", advent06Lines(), 41)
}

func advent06Lines() []string {
	padded := []string{
		//        0 1 2 3 4 5 6 7 8 9
		/* 9 */ " . . . . # . . . . .",
		/* 8 */ " . . . . . . . . . #",
		/* 7 */ " . . . . . . . . . .",
		/* 6 */ " . . # . . . . . . .",
		/* 5 */ " . . . . . . . # . .",
		/* 4 */ " . . . . . . . . . .",
		/* 3 */ " . # . o ^ . . . . .",
		/* 2 */ " . . . . . . o o # .",
		/* 1 */ " # o . o . . . . . .",
		/* 0 */ " . . . . . . # o . .",
	}
	return shared.StripPadding(padded)
}

func TestCountBlocksForIndefiniteLoops(t *testing.T) {
	run := func(name string, lines []string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			expected := extractExpected(lines)
			actual := CountBlocksForIndefiniteLoops(lines)
			shared.DiffLocSets(t, expected, actual)
		})
	}

	run("06 example", advent06Lines())
	run("straight line", straight06Lines())
	run("square", square06Lines())
	run("four square trap", fourSquareTrap06Lines())
	run("two square trap", twoSquareTrap06Lines())
	run("immediate block", immediateBlock06Lines())
	run("straight line no loops", straightLineNoLoops06Lines())
	run("big", big06Lines())
	run("spiral", spiral06Lines())
	run("inwards spiral", inwardsSpiral06Lines())
	run("snake", snake06Lines())
	run("crossing", crossing06Lines())
	run("precursor", precursor06Lines())
}

func straight06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3
			/* 0 */ " . . # . ",
			/* 1 */ " . . . # ",
			/* 2 */ " . . . . ",
			/* 3 */ " . o ^ . ",
			/* 4 */ " . . # . ",
		},
	)
}

func extractExpected(lines []string) *shared.Set[shared.Loc] {
	locations := []shared.Loc{}
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			r := shared.Abs(y - len(lines) + 1)
			if lines[r][x] == 'o' {
				locations = append(locations, shared.Loc{X: x, Y: y})
			}
		}
	}
	return shared.NewSet(locations)
}

func square06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " . . . . . . . . ",
			/* 1 */ " . # . . o . . . ",
			/* 2 */ " . o . o . o . # ",
			/* 3 */ " . o . # . . . . ",
			/* 4 */ " . o . # . . . . ",
			/* 5 */ " . . . # . . # . ",
			/* 6 */ " # ^ . . . . . . ",
			/* 7 */ " . . # . # . . . ",
		},
	)
}

func fourSquareTrap06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " . . . . . . . . ",
			/* 1 */ " . . # . o . . . ",
			/* 2 */ " . . . . . . # . ",
			/* 3 */ " . . . # . . . . ",
			/* 4 */ " . . ^ . . # . . ",
			/* 5 */ " . . . . . . . . ",
			/* 6 */ " . . . . . . . . ",
			/* 7 */ " . . . . . . . . ",
		},
	)
}

func twoSquareTrap06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " . . . . . . . . ",
			/* 1 */ " . . # . . o . . ",
			/* 2 */ " . . . . . . # . ",
			/* 3 */ " . . . . # . . . ",
			/* 4 */ " . . ^ . . # . . ",
			/* 5 */ " . . . . . . . . ",
			/* 6 */ " . . . . . . . . ",
			/* 7 */ " . . . . . . . . ",
		},
	)
}

func immediateBlock06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2
			/* 0 */ " . o . ",
			/* 1 */ " # ^ # ",
			/* 2 */ " . # . ",
		},
	)
}

func straightLineNoLoops06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2
			/* 0 */ " . . . ",
			/* 1 */ " . # . ",
			/* 2 */ " . . . ",
			/* 3 */ " . ^ . ",
		},
	)
}

func big06Lines() []string {
	return shared.StripPadding(
		[]string{
			//                             1 1 1
			//         0 1 2 3 4 5 6 7 8 9 0 1 2
			/*  0 */ " . # . . . . . . . . . . . ",
			/*  1 */ " . o . . . # . . . . . . . ",
			/*  2 */ " . . . . . . . . o o . . # ",
			/*  3 */ " o . . . . o # . . . . . . ",
			/*  4 */ " . o . . # o o . o # . . . ",
			/*  5 */ " . . . . . . . . # . . . . ",
			/*  6 */ " . . . . . . . . . . . . . ",
			/*  7 */ " . . . . . # . . . . . . . ",
			/*  8 */ " . . . . . . . # . . . o . ",
			/*  9 */ " . . . . . . . . . . . . . ",
			/* 10 */ " # . . o . o . . . . . . . ",
			/* 11 */ " . . . . . . . . . . . # . ",
			/* 12 */ " . . . . . . . . . . . . . ",
			/* 13 */ " . . . . ^ . . . . . . . . ",
			/* 14 */ " . . . . . . . . . . . . . ",
		},
	)
}

func TestBoardCopy(t *testing.T) {
	shared.InitTestLogging(t)

	// EXERCISE
	lines := []string{
		".#..",
		"..#.",
		".^..",
	}
	orig := newFatBoard(vector{loc: shared.Loc{}, dir: shared.RealNorth}, shared.NewBoard(lines))
	copied := orig.copy()
	copied.move(copied.deriveNextLocation()) // Curr.loc and visited modified.
	blockLoc := shared.Loc{X: 0, Y: 1}
	copied.nestedBrd.Set(blockLoc, '#')

	req := require.New(t)
	// VERIFY
	req.Equal(vector{dir: shared.RealNorth}, orig.curr, "curr")
	req.Equal('.', orig.nestedBrd.GetOrDie(blockLoc), "changed block")
	req.Equal(
		shared.NewSet([]vector{{dir: shared.RealNorth}}),
		orig.visited,
		"visited")
}

func spiral06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " . . # . . . . . ",
			/* 1 */ " . . . . . o . # ",
			/* 2 */ " . . o # . . . . ",
			/* 3 */ " . . . . . # . . ",
			/* 4 */ " . . . ^ . . . . ",
			/* 5 */ " . # o . . . . . ",
			/* 6 */ " . . . . # . o . ",
			/* 7 */ " . . . . . . . . ",
		},
	)
}

func inwardsSpiral06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " # . . . o . . . ",
			/* 1 */ " . . . . . . . # ",
			/* 2 */ " . . # . o . . . ",
			/* 3 */ " . . . . . # . . ",
			/* 4 */ " . . . # . . . . ",
			/* 5 */ " . . . . # . . . ",
			/* 6 */ " . # . . . . . . ",
			/* 7 */ " ^ . . . . . # . ",
		},
	)
}

func snake06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " # . . o . . . . ",
			/* 1 */ " . . . . . . # . ",
			/* 2 */ " . . # o . . . . ",
			/* 3 */ " . # . . . . o # ",
			/* 4 */ " . . . . . # o . ",
			/* 5 */ " . . # . . . . . ",
			/* 6 */ " . . . . . . # . ",
			/* 7 */ " ^ . . . . . . . ",
		},
	)
}

func crossing06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " . . . . # . . . ",
			/* 1 */ " . # . . . . # . ",
			/* 2 */ " # . . o . . o # ",
			/* 3 */ " . . . . . # o . ",
			/* 4 */ " . . . . . . o . ",
			/* 5 */ " . . . . . . . . ",
			/* 6 */ " o . . o . o . . ",
			/* 7 */ " . . . . ^ . # . ",
		},
	)
}

func precursor06Lines() []string {
	return shared.StripPadding(
		[]string{
			//        0 1 2 3 4 5 6 7
			/* 0 */ " . # . o . . . . ",
			/* 1 */ " . . . . . . . # ",
			/* 2 */ " . . . o . . . # ",
			/* 3 */ " . . . . . . . # ",
			/* 4 */ " . . . . . . . . ",
			/* 5 */ " . . . . . . . . ",
			/* 6 */ " . . # . . . . . ",
			/* 7 */ " . ^ . . . . # . ",
		},
	)
}

func TestBoardPrint(t *testing.T) {
	run := func(name string, brd *fatBoard, expected string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := brd.print()
			req.Equal(expected, actual)
		})
	}

	run(
		"start",
		newTestBoard(
			shared.Loc{X: 1, Y: 0},
			nil,
			[]string{
				" # ",
				"   ",
			},
		),
		strings.Join(
			[]string{
				" # ",
				" * ",
			},
			"\n")+"\n")
	run(
		"n",
		newTestBoard(
			shared.Loc{X: 2, Y: 0},
			[]vector{
				{loc: shared.Loc{X: 1, Y: 0}, dir: shared.RealNorth},
				{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealNorth},
				{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealEast},
				{loc: shared.Loc{X: 2, Y: 1}, dir: shared.RealEast},
				{loc: shared.Loc{X: 2, Y: 1}, dir: shared.RealSouth},
				{loc: shared.Loc{X: 2, Y: 0}, dir: shared.RealSouth},
			},
			[]string{
				" #  ",
				"   #",
				"    ",
			},
		),
		strings.Join(
			[]string{
				//       0123
				/* 2 */ " #  ",
				/* 1 */ " ++#",
				/* 0 */ " |* ",
			},
			"\n")+"\n")
}

func newTestBoard(
	curr shared.Loc,
	visited []vector,
	lines []string,
) *fatBoard {
	fatBoard := newFatBoard(vector{loc: curr, dir: shared.RealNorth}, shared.NewBoard(lines))
	if visited != nil {
		fatBoard.visited = shared.NewSet(visited)
	}
	return fatBoard
}
