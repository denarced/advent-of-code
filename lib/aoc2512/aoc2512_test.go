package aoc2512

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func _TestCountFittingRegions(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	// EXERCISE
	count, err := CountFittingRegions(lines)
	req.NoError(err)
	req.Equal(2, count)
}

func TestParseLines(t *testing.T) {
	req := require.New(t)
	presents, regions, err := parseLines([]string{
		"0:",
		"#.",
		".#",
		"",
		"1:",
		".#",
		"##",
		"",
		"1x2: 0 1",
		"2x1: 2 0",
	})
	req.NoError(err)
	req.Equal(
		[]present{
			{
				table: []string{
					"#.",
					".#",
				},
				spots: 2,
			},
			{
				table: []string{
					".#",
					"##",
				},
				spots: 3,
			},
		},
		presents)
	req.Equal(
		[]region{
			{
				width:  1,
				height: 2,
				counts: []int{0, 1},
			},
			{
				width:  2,
				height: 1,
				counts: []int{2, 0},
			},
		},
		regions)
}
