package aoc2304

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestSumPoints(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	require.Equal(t, 13, SumPoints(lines))
}

func TestParseLines(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	cards, err := parseLines([]string{
		"Card 1: 8 4 2 | 16 32",
	})
	req.NoError(err, "failed to parse lines")
	req.Equal(1, len(cards))
	first := cards[0]
	req.Equal(1, first.ID)
	req.True(gent.NewSet(2, 4, 8).Equal(first.winners), "winners")
	req.True(gent.NewSet(16, 32).Equal(first.yours), "yours")
}
