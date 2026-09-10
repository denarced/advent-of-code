package aoc2215

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindNoGoForBeacons(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	req.Equal(26, FindNoGoForBeacons(lines, 10))
}

func TestDeriveRange(t *testing.T) {
	run := func(centerY, overlap, expectedFrom, expectedTo int) {
		name := fmt.Sprintf("centerY:%d overlap:%d", centerY, overlap)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			from, to := deriveRange(centerY, overlap)

			// VERIFY
			req.Equal([2]int{from, to}, [2]int{expectedFrom, expectedTo})
		})
	}

	run(10, 0, 10, 10)
	run(10, 1, 9, 11)
	run(10, 2, 8, 12)
}

func assertRangeSet(ass *assert.Assertions, rSet *rangeSet, expected [][2]int) {
	ass.Equal(len(expected), len(rSet.ranges), "range set length mismatch")
	maxIndex := min(len(expected), len(rSet.ranges))
	for i := range maxIndex {
		ass.Equal(expected[i], rSet.ranges[i])
	}
}

func TestRangeSetAdd(t *testing.T) {
	run := func(name string, added [][2]int, expected [][2]int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			ass := assert.New(t)

			// EXERCISE
			rSet := new(rangeSet)
			for _, each := range added {
				rSet.add(each)
			}

			// VERIFY
			assertRangeSet(ass, rSet, expected)
		})
	}

	run("no join", [][2]int{{1, 1}}, [][2]int{{1, 1}})
	run(
		"1-2+4-5",
		[][2]int{
			{1, 2},
			{4, 5},
		},
		[][2]int{
			{1, 2},
			{4, 5},
		})
	run(
		"1-2+3-4",
		[][2]int{
			{1, 2},
			{3, 4},
		},
		[][2]int{{1, 4}})
	run(
		"5-7+3-6",
		[][2]int{
			{5, 7},
			{3, 6},
		},
		[][2]int{{3, 7}})
	run(
		"0-3+4-6+7-20",
		[][2]int{
			{0, 3},
			{4, 6},
			{7, 20},
		},
		[][2]int{{0, 20}})
	run(
		"7-9+3-5+1-1+2-2+6-6",
		[][2]int{
			{7, 9},
			{3, 5},
			{1, 1},
			{2, 2},
			{6, 6},
		},
		[][2]int{{1, 9}})
	run(
		"2-2+2-2",
		[][2]int{
			{2, 2},
			{2, 2},
		},
		[][2]int{{2, 2}})
	run(
		"3-5+1-7",
		[][2]int{
			{3, 5},
			{1, 7},
		},
		[][2]int{{1, 7}})
	run(
		"1-7+1-2",
		[][2]int{
			{1, 7},
			{1, 2},
		},
		[][2]int{{1, 7}})
}

func TestRangeSetRemove(t *testing.T) {
	run := func(name string, added [][2]int, removed int, expected [][2]int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			ass := assert.New(t)

			rSet := new(rangeSet)
			for _, each := range added {
				rSet.add(each)
			}
			// EXERCISE
			rSet.remove(removed)

			// VERIFY
			assertRangeSet(ass, rSet, expected)
		})
	}

	run("-1", nil, 1, nil)
	run("+3,-3", [][2]int{{3, 3}}, 3, nil)

	runBase := func(removed int, expected [][2]int) {
		base := [][2]int{
			{1, 2},
			{4, 6},
			{8, 8},
			{10, 20},
		}
		name := fmt.Sprintf("base:-%d", removed)
		run(name, base, removed, gent.Tri(expected == nil, base, expected))
	}
	runBase(
		8,
		[][2]int{
			{1, 2},
			{4, 6},
			{10, 20},
		})
	runBase(0, nil)
	runBase(21, nil)
	runBase(
		4,
		[][2]int{
			{1, 2},
			{5, 6},
			{8, 8},
			{10, 20},
		})
	runBase(
		20,
		[][2]int{
			{1, 2},
			{4, 6},
			{8, 8},
			{10, 19},
		})
	runBase(
		5,
		[][2]int{
			{1, 2},
			{4, 4},
			{6, 6},
			{8, 8},
			{10, 20},
		})
}

func TestFindBeaconFrequency(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	// EXERCISE & VERIFY
	req.Equal(56_000_011, FindBeaconFrequency(lines, 20))
}

func TestDeriveDiamond(t *testing.T) {
	run := func(sens sensor, maxCoord int, expected [][3]int) {
		name := fmt.Sprintf("%v -> %v with max %d", sens.loc, sens.closest, maxCoord)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			var expectedIndex int
			deriveDiamond(sens, maxCoord, func(y, xFrom, xTo int) {
				req.Equal([3]int{y, xFrom, xTo}, expected[expectedIndex])
				expectedIndex++
			})
			req.Equal(len(expected), expectedIndex, "expected index should match length")
		})
	}

	run(
		sensor{
			loc:     shared.Loc{X: 1, Y: 1},
			closest: shared.Loc{X: 1, Y: 2},
		},
		2,
		[][3]int{
			{0, 1, 1},
			{1, 0, 2},
			{2, 1, 1},
		})
	run(
		sensor{
			loc:     shared.Loc{X: 1, Y: 1},
			closest: shared.Loc{X: 1, Y: 3},
		},
		2,
		[][3]int{
			{0, 0, 2},
			{1, 0, 2},
			{2, 0, 2},
		})
	run(
		sensor{
			loc:     shared.Loc{X: -1, Y: 1},
			closest: shared.Loc{X: -2, Y: 1},
		},
		2,
		[][3]int{
			{1, 0, 0},
		})
	run(
		sensor{
			loc:     shared.Loc{X: 3, Y: 1},
			closest: shared.Loc{X: 4, Y: 1},
		},
		2,
		[][3]int{
			{1, 2, 2},
		})
	run(
		sensor{
			loc:     shared.Loc{X: 1, Y: -1},
			closest: shared.Loc{X: 0, Y: -1},
		},
		2,
		[][3]int{
			{0, 1, 1},
		})
}
