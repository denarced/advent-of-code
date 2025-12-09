package aoc2508

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountCircuits(t *testing.T) {
	readLines := func(req *require.Assertions) []string {
		lines, err := inr.ReadPath("testdata/in.txt")
		req.NoError(err, "read test data")
		return lines
	}

	t.Run("three largest", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		req.Equal(40, CountCircuits(readLines(req), 10))
	})

	t.Run("last x*x", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		req.Equal(25_272, CountCircuits(readLines(req), 0))
	})
}
