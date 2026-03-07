package aoc2306

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiplyCounts(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	req.Equal(288, MultiplyCounts(lines))
}

func TestFindRoots(t *testing.T) {
	neg, pos := findRoots(race{
		time:     15,
		distance: 40,
	})
	ass := assert.New(t)
	delta := 0.000001
	ass.InDelta(3.594875, neg, delta)
	ass.InDelta(11.405124, pos, delta)
}
