package aoc2507

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountSplits(t *testing.T) {
	readLines := func(req *require.Assertions) []string {
		lines, err := inr.ReadPath("testdata/in.txt")
		req.NoError(err, "read test data")
		return lines
	}

	t.Run("splits", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		req.Equal(21, CountSplits(readLines(req)))
	})
	t.Run("timelines", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		req.Equal(40, CountTimelines(readLines(req)))
	})
}
