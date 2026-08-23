package aoc2205

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveTopCrates(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty(), inr.NoTrim())
	req.NoError(err, "read test data")

	req.Equal("CMZ", DeriveTopCrates(lines))
}
