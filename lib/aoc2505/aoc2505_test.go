package aoc2505

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountFreshAvailableIngredients(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err)
	req.Equal(3, CountFreshAvailableIngredients(lines))
}

func TestCountFreshIngredients(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err)
	req.Equal(14, CountFreshIngredients(lines))
}

func TestJoinIntRanges(t *testing.T) {
	var tests = []struct {
		first    intRange
		second   intRange
		expected *intRange
	}{
		{intRange{1, 3}, intRange{5, 7}, nil},
		{intRange{1, 3}, intRange{4, 6}, nil},
		{intRange{1, 3}, intRange{3, 5}, &intRange{1, 5}},
		{intRange{1, 3}, intRange{2, 4}, &intRange{1, 4}},
		{intRange{1, 3}, intRange{1, 3}, &intRange{1, 3}},
		{intRange{1, 3}, intRange{0, 2}, &intRange{0, 3}},
		{intRange{1, 3}, intRange{-1, 1}, &intRange{-1, 3}},
		{intRange{1, 3}, intRange{-2, 0}, nil},
		{intRange{1, 3}, intRange{-3, -1}, nil},
	}
	for _, tt := range tests {
		name := fmt.Sprintf(
			"%d..%d+%d..%d",
			tt.first.from,
			tt.first.to,
			tt.second.from,
			tt.second.to,
		)
		t.Run(name, func(t *testing.T) {
			req := require.New(t)
			joined := joinIntRanges(&tt.first, &tt.second)
			if joined == nil {
				req.Nil(tt.expected)
			} else {
				req.Equal(*tt.expected, *joined)
			}
		})
	}
}

func TestMergeIntRanges(t *testing.T) {
	shared.InitTestLogging(t)
	ranges := []intRange{
		{0, 10},
		{40, 50},
		{10, 20},
		{30, 40},
		{20, 30},
		{100, 200},
		{1, 49},
	}
	merged := mergeIntRanges(ranges)
	req := require.New(t)
	req.Equal(
		[]intRange{
			{0, 50},
			{100, 200},
		},
		merged,
	)
}
