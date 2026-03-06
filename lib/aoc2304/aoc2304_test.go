package aoc2304

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestSumPoints(t *testing.T) {
	run := func(expected int, spawn bool) {
		name := "spawn"
		if !spawn {
			name = "non-spawn"
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "failed to read test data")
			require.Equal(t, expected, SumPoints(lines, spawn))
		})
	}
	run(13, false)
	run(30, true)
}

func TestParseLines(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	cards, maxID, err := parseLines([]string{
		"Card 1: 8 4 2 | 16 32",
	})
	req.NoError(err, "failed to parse lines")
	req.Equal(1, len(cards), "card count")
	first, ok := cards[1]
	req.True(ok, "first card not found")
	req.Equal(1, first.ID, "first ID")
	req.True(gent.NewSet(2, 4, 8).Equal(first.winners), "winners")
	req.True(gent.NewSet(16, 32).Equal(first.yours), "yours")
	req.Equal(1, maxID, "max ID")
}
