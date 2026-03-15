package aoc2309

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestSumExtrapolatedValues(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := shared.ReadLinesFromFile("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	require.Equal(t, 114, SumExtrapolatedValues(lines))
}

func TestExtrapolate(t *testing.T) {
	run := func(ints []int, expected int) {
		t.Run(fmt.Sprintf("%v -> %d", ints, expected), func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, extrapolate([]int{1, 2, 3}))
		})
	}

	run([]int{1, 2, 3}, 4)
}
