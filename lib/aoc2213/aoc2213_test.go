package aoc2213

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumSortedIndexes(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err, "failed to read test data")

	req.Equal(13, SumSortedIndexes(lines))
}

func TestSumValues(t *testing.T) {
	values := []any{}
	values = append(values, 4)
	values = append(values, []any{1, 2})
	require.Equal(t, 7, SumValues(values))
}

func TestSortAll(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	req.Equal(140, SortAll(lines))
}
