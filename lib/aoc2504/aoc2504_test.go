package aoc2504

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountRolls(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := shared.ReadLinesFromFile("testdata/in.txt")
	req.NoError(err, "read lines")

	ass := assert.New(t)
	ass.Equal(13, CountRolls(lines, 1))
	ass.Equal(43, CountRolls(lines, -1))
}
