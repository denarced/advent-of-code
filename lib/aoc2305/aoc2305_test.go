package aoc2305

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveLowestLocation(t *testing.T) {
	run := func(useRange bool, expected int) {
		name := "range"
		if !useRange {
			name = "!range"
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
			req.NoError(err, "failed to read test data")
			req.Equal(expected, DeriveLowestLocation(lines, useRange))
		})
	}
	run(false, 35)
	run(true, 46)
}

func TestSplitToBlocks(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	blocks := splitToBlocks([]string{
		"",
		"seed-to-soil map:",
		"50 98 2",
		"52 50 48",
		"",
		"soil-to-fertilizer map:",
		"0 15 37",
		"37 52 2",
		"39 0 15",
	})
	req.Equal(
		[][]string{
			{
				"seed-to-soil map:",
				"50 98 2",
				"52 50 48",
			}, {
				"soil-to-fertilizer map:",
				"0 15 37",
				"37 52 2",
				"39 0 15",
			},
		},
		blocks)
}

func TestTranslate(t *testing.T) {
	shared.InitTestLogging(t)
	tested := corr{
		specs: []spec{
			{src: 10, dst: 15, size: 2},
			{src: 13, dst: 5, size: 2},
		},
	}
	for _, each := range []struct {
		src      int
		expected int
	}{
		// First spec.
		{10, 15},
		{11, 16},
		// No spec.
		{12, 12},
		// Second spec.
		{13, 5},
		{14, 6},
		// No spec.
		{15, 15},
	} {
		t.Run(fmt.Sprintf("%d -> %d", each.src, each.expected), func(t *testing.T) {
			require.Equal(t, each.expected, tested.translate(each.src))
		})
	}
}

func TestToRanges(t *testing.T) {
	require.Equal(
		t,
		[]intRange{
			{start: 79, end: 92},
			{start: 10, end: 10},
		},
		toRanges([]int{79, 14, 10, 1}, true))
}
