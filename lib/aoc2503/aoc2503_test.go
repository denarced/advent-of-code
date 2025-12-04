package aoc2503

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestDeriveMaxJoltageSum(t *testing.T) {
	readLines := func(req *require.Assertions) []string {
		lines, err := shared.ReadLinesFromFile("testdata/in.txt")
		req.NoError(err, "failed to read test data lines")
		return lines
	}

	t.Run("2", func(t *testing.T) {
		shared.InitNullLogging()
		req := require.New(t)
		// EXERCISE
		req.Equal(int64(357), DeriveMaxJoltageSum(readLines(req), 2))
	})

	t.Run("12", func(t *testing.T) {
		shared.InitNullLogging()
		req := require.New(t)
		// EXERCISE
		req.Equal(int64(3121910778619), DeriveMaxJoltageSum(readLines(req), 12))
	})
}
