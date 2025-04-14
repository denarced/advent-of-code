package aoc2411

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

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

func TestWalkIntoStones(t *testing.T) {
	run := func(name string, values []int, blinks, expected int) {
		if name == "" {
			plural := ""
			if blinks != 1 {
				plural = "s"
			}
			name = fmt.Sprintf("%d blink%s %s", blinks, plural, joinInts(values, ", "))
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := CountStones(values, blinks)
			req.Equal(big.NewInt(int64(expected)), actual)
		})
	}

	run("empty", []int{}, 0, 0)
	run("nil", nil, 0, 0)
	run("", []int{99}, 0, 1)
	run("", []int{22, 33, 44, 55, 66, 77, 88}, 0, 7)
	run("", []int{0}, 1, 1)
	run("", []int{11}, 1, 2)
	run("", []int{4401}, 2, 3)
	// 6. 4
	// 5. 8096
	// 4. 80 96
	// 3. 8 0 9 6
	// 2. 16192 1 18216 12144
	// 1. 32772608 2024 36869184 24579456
	// 0. 3277 2608 20 24 3686 9184 2457 9456
	run("", []int{4}, 6, 8)
	// 0. 17 125
	// 1. 1 7 253000
	// 2. 2024 14168 253 0
	// 3. 20 24 28676032 512072 1
	// 4. 2    0 2    4    2867  6032  512     72  2024
	// 5. 4048  1    4048  8096  28  67  60  32  1036288    7     2    20  24
	// 6. 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2 2097446912 14168 4048 2 0 2 4
	run("", []int{17, 125}, 6, 22)
	run("", []int{125}, 6, 7)
	// 0. 17
	// 1. 1                7
	// 2. 2024             14168
	// 3. 20  24           28676032
	// 4. 2    0 2    4    2867  6032
	// 5. 4048 1 4048 8096 28 67 60 32
	// 6. 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2
	run("", []int{17}, 6, 15)
}

func joinInts(s []int, sep string) string {
	strs := gent.Map(
		s,
		func(i int) string {
			return strconv.Itoa(i)
		})
	return strings.Join(strs, sep)
}
