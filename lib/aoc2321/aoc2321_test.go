package aoc2321

import (
	"strconv"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountRangeFromLines(t *testing.T) {
	run := func(stepCount, expected int) {
		t.Run(strconv.Itoa(stepCount), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "failed to read test data")

			// EXERCISE & VERIFY
			req.Equal(expected, CountRangeFromLines(lines, stepCount))
		})
	}

	run(1, 2)
	run(2, 4)
	run(6, 16)
}
