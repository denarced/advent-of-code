package aoc2311

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumDistances(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err)
	req.Equal(374, SumDistances(lines))
}

func TestExpand(t *testing.T) {
	shared.InitTestLogging(t)
	expanded := expand([]string{
		"##...",
		".....",
		"...#.",
	})
	require.Equal(
		t,
		[]string{
			"##.....",
			".......",
			".......",
			"....#..",
		},
		expanded)
}
