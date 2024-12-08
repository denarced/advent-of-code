package aoc2407

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestDeriveCalibrationSum(t *testing.T) {
	run := func(name string, withConcat bool, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, DeriveCalibrationSum(lines, withConcat))
		})
	}

	run("empty", false, []string{}, 0)
	run("example without concat", false, getExampleLines(), 3749)
	run("example with concat", true, getExampleLines(), 11387)
}

func getExampleLines() []string {
	return []string{
		"190: 10 19",
		"3267: 81 40 27",
		"83: 17 5",
		"156: 15 6",
		"7290: 6 8 6 15",
		"161011: 16 10 13",
		"192: 17 8 14",
		"21037: 9 7 18 13",
		"292: 11 6 16 20",
	}
}

func TestGeneratePermutations(t *testing.T) {
	run := func(length, base int, expected [][]int) {
		t.Run(fmt.Sprintf("%d-%d", length, base), func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, generatePermutations(length, base))
		})
	}

	run(1, 2, [][]int{{0}, {1}})
	run(2, 2, [][]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}})
	run(3, 2, [][]int{
		{0, 0, 0},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 1},
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, 0},
		{1, 1, 1},
	})

	run(1, 3, [][]int{{0}, {1}, {2}})
	run(
		2,
		3,
		[][]int{
			{0, 0},
			{0, 1},
			{0, 2},
			{1, 0},
			{1, 1},
			{1, 2},
			{2, 0},
			{2, 1},
			{2, 2},
		})
}

func TestConcat(t *testing.T) {
	shared.InitTestLogging(t)
	require.Equal(t, 149910, concat(1499, 10))
}
