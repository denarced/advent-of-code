package aoc2201

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveHeaviest(t *testing.T) {
	run := func(limit, expected int) {
		t.Run(fmt.Sprintf("%d-%d", limit, expected), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
			req.NoError(err, "read test data")

			req.Equal(expected, DeriveHeaviest(lines, limit))
		})
	}
	run(1, 24_000)
	run(3, 45_000)
}
