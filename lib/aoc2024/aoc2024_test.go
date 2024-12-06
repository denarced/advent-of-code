package aoc2024

import (
	"fmt"
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

func TestAdvent01Distance(t *testing.T) {
	run := func(name string, left, right []int, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := Advent01Distance(left, right)
			req.Equal(expected, actual)
		})
	}

	run("empty", []int{}, []int{}, 0)
	run("happy path", []int{11, 7, 1}, []int{2, 11, 5}, 3)
}

func TestAdvent01Similarity(t *testing.T) {
	run := func(name string, left, right []int, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			similarity := Advent01Similarity(left, right)
			req.Equal(expected, similarity)
		})
	}

	run("empty", []int{}, []int{}, 0)
	run("empty right", []int{1, 2, 3}, []int{}, 0)
	run("empty left", []int{}, []int{1, 2, 2, 3}, 0)
	run("no shared", []int{1, 2, 2, 3, 3, 3}, []int{6, 6, 6}, 0)
	run("one shared", []int{3, 1, 3}, []int{3, 4, 5}, 6)
	run("multiple shared", []int{2, 3, 4}, []int{2, 2, 3, 3, 3, 5}, 4+9)
}

func TestToInts(t *testing.T) {
	run := func(name string, strings []string, expected []int, errMessage string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual, err := ToInts(strings)
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

func TestToIntTable(t *testing.T) {
	run := func(name string, lines []string, expected [][]int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := ToIntTable(lines)
			req.Equal(expected, actual)
		})
	}

	run("nil", nil, nil)
	run("empty", []string{}, nil)
	run("1x1", []string{" 33 "}, [][]int{{33}})
	run("3x2", []string{"1 2", "2 3", "3 4"}, [][]int{{1, 2}, {2, 3}, {3, 4}})
}

func TestCountSafe(t *testing.T) {
	run := func(name string, levels [][]int, dampener bool, expected int) {
		prefix := "wo dampener"
		if dampener {
			prefix = "with dampener"
		}
		t.Run(fmt.Sprintf("%s - %s", prefix, name), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := CountSafe(levels, dampener)
			req.Equal(expected, actual)
		})
	}

	// Without dampener.
	run("empty", [][]int{}, false, 0)
	run("too short", [][]int{{1}}, false, 0)
	run("minimal asc safe", [][]int{{1, 2}}, false, 1)
	run("minimal desc safe", [][]int{{2, 1}}, false, 1)
	run("asc to desc", [][]int{{1, 2, 1}}, false, 0)
	run("too little diff", [][]int{{1, 2, 2}}, false, 0)
	run("too big diff", [][]int{{1, 2, 3, 7}}, false, 0)
	run("2 safe, 1 unsafe", [][]int{{1, 4, 6, 7}, {9, 5, 4, 3, 2}, {100, 97, 95, 94}}, false, 2)

	// With dampener.
	run("empty", [][]int{}, true, 0)
	run("too short", [][]int{{1}}, true, 0)
	run("minimal asc safe", [][]int{{1, 2}}, true, 1)
	run("minimal desc safe", [][]int{{2, 1}}, true, 1)
	run("asc to desc", [][]int{{1, 2, 1}}, true, 1)
	run("too little diff", [][]int{{1, 2, 2}}, true, 1)
	run("too big diff", [][]int{{1, 2, 3, 7}}, true, 1)
	run("2 safe, 1 unsafe", [][]int{{1, 4, 6, 7}, {9, 5, 4, 3, 2}, {100, 97, 95, 94}}, true, 3)
	run("long asc to desc", [][]int{{4, 7, 4, 1}}, true, 1)
	run("fix before unsafe pair", [][]int{{4, 7, 4, 1}}, true, 1)
	run("drop too big diff", [][]int{{5, 9, 8}}, true, 1)
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
	run(
		"example",
		[]string{
			// 23456789
			"....#.....", // 0
			".........#", // 1
			"..........", // 2
			"..#.......", // 3
			".......#..", // 4
			"..........", // 5
			".#..^.....", // 6
			"........#.", // 7
			"#.........", // 8
			"......#...", // 9
		},
		41)
}
