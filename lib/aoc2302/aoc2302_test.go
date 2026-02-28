package aoc2302

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveGameCountSum(t *testing.T) {
	shared.InitTestLogging(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req := require.New(t)
	req.NoError(err, "failed to read test data")
	count := DeriveGameCountSum(
		lines,
		map[Kind]int{
			KindBlue:  14,
			KindRed:   12,
			KindGreen: 13,
		})
	require.Equal(t, 8, count)
}

func TestDerivePowerSum(t *testing.T) {
	lines, err := inr.ReadPath("testdata/in.txt")
	req := require.New(t)
	req.NoError(err, "failed to read test data")
	shared.InitTestLogging(t)
	count := DerivePowerSum(lines)
	require.Equal(t, 2286, count)
}
