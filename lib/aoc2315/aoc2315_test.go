package aoc2315

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumHashes(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err)

	// EXERCISE & VERIFY
	req.Equal(1320, SumHashes(lines))
}

func TestHash(t *testing.T) {
	require.Equal(t, 30, hash("rn=1"))
}

func TestParseLines(t *testing.T) {
	require.Equal(
		t,
		[]string{" ", "ab", "cd"},
		parseLines([]string{"", ", ,", ",,ab,,cd,,"}))
}
