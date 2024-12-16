package aoc2411

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountStones(t *testing.T) {
	run := func(name string, blinkCount int, stones []int, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := CountStones(blinkCount, stones)
			req.Equal(expected, actual)
		})
	}

	run("1 blink", 1, []int{0, 1, 10, 99, 999}, 7)
	run("2 blinks", 2, []int{22, 44}, 4)

	// 6|125
	// 5|253000
	// 4|253 0
	// 3|512072 1
	// 2|512 72 2024
	// 1|1036288 7 2 20 24
	// 0|2097446912 x x 2 0 2 4

	// 6| 17
	// 5| 1                7
	// 4| 2024             14168
	// 3| 20  24           28676032
	// 2| 2    0 2    4    2867  6032
	// 1| 4048 1 4048 8096 28 67 60 32
	// 0| 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2
	run("6 blinks", 6, []int{125, 17}, 22)
	run("6 125", 6, []int{125}, 7)
	run("6 17", 6, []int{17}, 15)
	run("25 blinks", 25, []int{125, 17}, 55_312)
}

func TestSplitStone(t *testing.T) {
	run := func(stone, first, second int, ok bool) {
		t.Run(strconv.Itoa(stone), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			f, s, cloned := splitStone(stone)
			req.Equal(first, f)
			req.Equal(second, s)
			req.Equal(ok, cloned)
		})
	}

	run(1, 0, 0, false)
	run(2, 0, 0, false)
	run(10, 1, 0, true)
	run(22, 2, 2, true)
	run(3344, 33, 44, true)
	run(33445566, 3344, 5566, true)
}

func TestCountNodesWithMap(t *testing.T) {
	run := func(stone, blinks int, expected map[spottedStone]int) {
		t.Run(fmt.Sprintf("%d-%d", stone, blinks), func(t *testing.T) {
			shared.InitTestLogging(t)

			// EXERCISE
			tree := buildPartialTree(stone, blinks)
			stoneToCount := map[spottedStone]int{}
			countNodes(tree, stoneToCount)

			// VERIFY
			for spotted, expectedCount := range expected {
				assert.Equalf(t, expectedCount, stoneToCount[spotted], "boo: %v", spotted)
			}
		})
	}

	// 12
	// 1       2
	// 2024    4048
	// 20  20  40  48
	// 2 0 2 0 4 0 4 8
	run(12, 4, map[spottedStone]int{
		{value: 20, spots: 1}: 2,
	})
	// 6| 17
	// 5| 1 7
	// 4| 2024 14168
	// 3| 20 24 28676032
	// 2| 2 0 2 4 2867 6032
	// 1| 4048 1 4048 8096 28 67 60 32
	// 0| 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2
	run(17, 6, map[spottedStone]int{
		{value: 20, spots: 3}: 3,
	})

	// Initial arrangement:
	// 0| 17
	// 1| 1 7
	// 2| 2024
	// 3| 20
	// 4| 2
	// 5| 4048
	// 6| 40  48
	// 7| 4 0 4 8
	// 8| 8096 1 8096 16192
	// 9| 80 96 2024 80 96 32772608
	// a| 8 0 9 6 20 24 8 0 9 6 3277 2608
	run(17, 10, map[spottedStone]int{
		{value: 4048, spots: 5}: 12,
	})
}

func TestBuildPartialTree(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	// EXERCISE
	actual := buildPartialTree(125, 6)

	// VERIFY
	req.Equal(512_072, actual.value)
	req.Equal(3, actual.blinks)

	secLevelLeft := actual.kids[0]
	req.Equal(512, secLevelLeft.value)
	req.Equal(2, secLevelLeft.blinks)

	secLevelRight := actual.kids[1]
	req.Equal(72, secLevelRight.value)
	req.Equal(2, secLevelRight.blinks)

	firstLevelLeft := secLevelLeft.kids[0]
	req.Equal(1_036_288, firstLevelLeft.value)
	req.Equal(1, firstLevelLeft.blinks)
	req.Len(secLevelLeft.kids, 1)

	req.Len(secLevelRight.kids, 2)
	firstLevelFirstLeft := secLevelRight.kids[0]
	req.Equal(7, firstLevelFirstLeft.value)
	req.Equal(1, firstLevelFirstLeft.blinks)
	firstLevelFirstRight := secLevelRight.kids[1]
	req.Equal(2, firstLevelFirstRight.value)
	req.Equal(1, firstLevelFirstRight.blinks)

	req.Equal(1_036_288*2024, firstLevelLeft.kids[0].value)
	req.Equal(0, firstLevelLeft.kids[0].blinks)

	req.Equal(7*2024, firstLevelFirstLeft.kids[0].value)
	req.Equal(0, firstLevelFirstLeft.kids[0].blinks)

	req.Equal(2*2024, firstLevelFirstRight.kids[0].value)
	req.Equal(0, firstLevelFirstRight.kids[0].blinks)

	m := map[spottedStone]int{}
	req.Equal(3, countNodes(actual, m))
	req.Equal(
		map[spottedStone]int{
			{value: 2, spots: 1}:         1,
			{value: 7, spots: 1}:         1,
			{value: 1_036_288, spots: 1}: 1,
			{value: 72, spots: 2}:        2,
			{value: 512, spots: 2}:       1,
			{value: 512_072, spots: 3}:   3,
		},
		m)
}
