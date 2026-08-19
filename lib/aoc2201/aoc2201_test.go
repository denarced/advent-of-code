package aoc2201

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveHeaviest(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err, "read test data")

	req.Equal(24_000, DeriveHeaviest(lines))
}
