package aoc2504

import (
	"os"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountRolls(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	f, err := os.Open("testdata/in.txt")
	req.NoError(err, "open file")
	defer f.Close()
	lines, err := shared.ReadLines(f)
	req.NoError(err, "read lines")

	ass := assert.New(t)
	ass.Equal(13, CountRolls(lines, 1))
	ass.Equal(43, CountRolls(lines, -1))
}
