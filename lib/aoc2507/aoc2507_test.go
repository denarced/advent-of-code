package aoc2507

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountSplits(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "read test data")
	require.Equal(t, 21, CountSplits(lines))
}
