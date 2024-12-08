package aoc2024

import (
	"fmt"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestToColumns(t *testing.T) {
	run := func(name string, lines, expectedLeft, expectedRight []string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			left, right := shared.ToColumns(lines)
			req.Equal(expectedLeft, left)
			req.Equal(expectedRight, right)
		})
	}

	run("empty", []string{}, nil, nil)
	run("space", []string{"abc efg"}, []string{"abc"}, []string{"efg"})
	run("two spaces", []string{"313  666"}, []string{"313"}, []string{"666"})
}

func TestToInts(t *testing.T) {
	run := func(name string, s []string, expected []int, errMessage string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual, err := ToInts(s)
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

func TestMultiply(t *testing.T) {
	run := func(text string, logic bool, expected int) {
		name := fmt.Sprintf("%slogic: %s", shared.Or(logic, "", "!"), text)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, Multiply(text, logic))
		})
	}

	run("empty", false, 0)
	run("empty", true, 0)
	run("#mul(2,3)", false, 6)
	run("--mul(3,4)mul(6,2)", false, 24)
	run("##mul(a,3)-mul(,3)", false, 0)
	run("mul(2,3)do()mul(3,4)don't()mul(5,2)", true, 2*3+3*4)
	run("don't()mulmulmul(2,3)mul(23,3)do()mul(3,4)", true, 3*4)
}

func TestCountInTable(t *testing.T) {
	run := func(name string, table []string, word string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			count := CountInTable(table, word)
			req.Equal(expected, count)
		})
	}

	run("empty", []string{}, "JOB", 0)
	run("one horizontal", []string{"..XMAS.."}, "XMAS", 1)
	run(
		"happy path",
		[]string{
			"XS..S..S..S.SX",
			"S....A.A.A...S",
			"......MMM.....",
			"....SAMXMAS...",
			"......MMM.....",
			".....A.A.A....",
			"S...S..S..S..S",
			"XS..........SX",
		},
		"XMAS",
		8,
	)
}

func TestCountWordCrosses(t *testing.T) {
	run := func(name string, table []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, CountWordCrosses(table, "MAS"))
		})
	}

	run("empty", []string{}, 0)
	run(
		"half x",
		[]string{
			"......",
			".M....",
			"..A...",
			".M.S..",
		},
		0,
	)
	run(
		"happy path",
		[]string{
			".M.S.M",
			"..A.A.",
			".M.S.M",
			"......",
		},
		2,
	)
}

func TestSumCorrectMiddlePageNumbers(t *testing.T) {
	shared.InitTestLogging(t)
	// 143 is from the problem description.
	require.Equal(t, 143, SumCorrectMiddlePageNumbers(advent05Lines()))
}

func advent05Lines() []string {
	// Example values from problem description.
	return []string{
		"47|53",
		"97|13",
		"97|61",
		"97|47",
		"75|29",
		"61|13",
		"75|53",
		"29|13",
		"97|29",
		"53|29",
		"61|53",
		"97|53",
		"61|29",
		"47|13",
		"75|47",
		"97|75",
		"47|61",
		"75|61",
		"47|29",
		"75|13",
		"53|13",
		"75,47,61,53,29",
		"97,61,53,29,13",
		"75,29,13",
		"75,97,47,61,53",
		"61,13,29",
		"97,13,75,29,47",
	}
}

func TestSumIncorrectMiddlePageNumbers(t *testing.T) {
	shared.InitTestLogging(t)
	// 123 is from the problem description.
	require.Equal(t, 123, SumIncorrectMiddlePageNumbers(advent05Lines()))
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
		/* 0 */ " . . . . # . . . . .",
		/* 1 */ " . . . . . . . . . #",
		/* 2 */ " . . . . . . . . . .",
		/* 3 */ " . . # . . . . . . .",
		/* 4 */ " . . . . . . . # . .",
		/* 5 */ " . . . . . . . . . .",
		/* 6 */ " . # . o ^ . . . . .",
		/* 7 */ " . . . . . . o o # .",
		/* 8 */ " # o . o . . . . . .",
		/* 9 */ " . . . . . . # o . .",
	}
	return stripPadding(padded)
}

func stripPadding(lines []string) []string {
	stripped := make([]string, 0, len(lines))
	for _, each := range lines {
		stripped = append(stripped, strings.ReplaceAll(each, " ", ""))
	}
	return stripped
}

func TestCountBlocksForIndefiniteLoops(t *testing.T) {
	run := func(name string, lines []string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			expected := extractExpected(lines)
			actual := CountBlocksForIndefiniteLoops(lines)
			diffLocationSets(t, expected, actual)
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
	return stripPadding(
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

func extractExpected(lines []string) *shared.Set[location] {
	locations := []location{}
	for r := 0; r < len(lines); r++ {
		for c := 0; c < len(lines[r]); c++ {
			if lines[r][c] == 'o' {
				locations = append(locations, location{r, c})
			}
		}
	}
	return shared.NewSet(locations)
}

func square06Lines() []string {
	return stripPadding(
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
	return stripPadding(
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
	return stripPadding(
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
	return stripPadding(
		[]string{
			//        0 1 2
			/* 0 */ " . o . ",
			/* 1 */ " # ^ # ",
			/* 2 */ " . # . ",
		},
	)
}

func straightLineNoLoops06Lines() []string {
	return stripPadding(
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
	return stripPadding(
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

func diffLocationSets(t *testing.T, expected, actual *shared.Set[location]) {
	stringify := func(l location) string {
		return l.toString()
	}
	require.ElementsMatch(
		t,
		mapValues(expected.ToSlice(), stringify),
		mapValues(actual.ToSlice(), stringify),
	)
}

func TestBoardCopy(t *testing.T) {
	shared.InitTestLogging(t)

	// EXERCISE
	orig := newBoard(location{}, []location{}, 10, 11)
	copied := orig.copy()
	copied.move(copied.deriveNextLocation()) // Curr.loc and visited modified.
	copied.blocks.Add(location{6, 7})
	copied.width = 20
	copied.height = 21

	req := require.New(t)
	// VERIFY
	req.Equal(vector{dir: dirNorth}, orig.curr, "curr")
	req.Equal(shared.NewSet([]location{}), orig.blocks, "blocks")
	req.Equal(shared.NewSet([]vector{{dir: dirNorth}}), orig.visited, "visited")
	req.Equal(10, orig.width, "width")
	req.Equal(11, orig.height, "height")
}

func spiral06Lines() []string {
	return stripPadding(
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
	return stripPadding(
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
	return stripPadding(
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
	return stripPadding(
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
	return stripPadding(
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
	run := func(name string, brd *board, expected string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := brd.print()
			req.Equal(expected, actual)
		})
	}

	run(
		"start",
		newTestBoard(location{1, 1}, []location{{0, 1}}, nil, 3, 2),
		strings.Join(
			[]string{
				" # ",
				" * ",
			},
			"\n")+"\n")
	run(
		"n",
		newTestBoard(
			location{2, 2},
			[]location{{0, 1}, {1, 3}},
			[]vector{
				{loc: location{2, 1}, dir: dirNorth},
				{loc: location{1, 1}, dir: dirNorth},
				{loc: location{1, 1}, dir: dirEast},
				{loc: location{1, 2}, dir: dirEast},
				{loc: location{1, 2}, dir: dirSouth},
				{loc: location{2, 2}, dir: dirSouth},
			},
			4,
			3),
		strings.Join(
			[]string{
				//       0123
				/* 0 */ " #  ",
				/* 1 */ " ++#",
				/* 2 */ " |* ",
			},
			"\n")+"\n")
}

func newTestBoard(curr location, blocks []location, visited []vector, w, h int) *board {
	brd := newBoard(curr, blocks, w, h)
	if visited != nil {
		brd.visited = shared.NewSet(visited)
	}
	return brd
}
