package aoc2309

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestSumExtrapolatedValues(t *testing.T) {
	run := func(right bool, expected int) {
		name := gent.Tri(right, "right", "left")
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := shared.ReadLinesFromFile("testdata/in.txt")
			req.NoError(err, "failed to read test data")
			req.Equal(expected, SumExtrapolatedValues(lines, right))
		})
	}
	run(true, 114)
	run(false, 2)
}

func TestExtrapolate(t *testing.T) {
	run := func(ints []int, right bool, expected int) {
		t.Run(fmt.Sprintf("%v -> %d", ints, expected), func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, extrapolate([]int{1, 2, 3}, right))
		})
	}

	run([]int{1, 2, 3}, true, 4)
}
