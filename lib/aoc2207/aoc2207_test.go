package aoc2207

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "read test data")

    total,toDelete := SumRecursiveDirSize(lines)
	req.Equal(95_437, total, "total size")
    req.Equal(24_933_642, toDelete, "size to delete")
}
