package aoc2502

import (
	"fmt"
	"os"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	readTestData := func(req *require.Assertions) []string {
		file, err := os.Open("testdata/in.txt")
		req.NoError(err, "failed to open test file")
		defer file.Close()
		lines, err := shared.ReadLines(file)
		req.NoError(err, "failed to read lines")
		return lines
	}

	t.Run("twice", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		lines := readTestData(req)
		req.Equal(int64(1227775554), SumInvalidIDs(lines[0], true))
	})

	t.Run("more", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		lines := readTestData(req)
		req.Equal(int64(4174379265), SumInvalidIDs(lines[0], false))
	})
}

func TestBreaks(t *testing.T) {
	var tests = []struct {
		n        int64
		minSplit int
		maxSplit int
		expected bool
	}{
		{11, 2, 2, true},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%d-%d-%d-%v", tt.n, tt.minSplit, tt.maxSplit, tt.expected)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, tt.expected, breaks(tt.n, tt.minSplit, tt.maxSplit))
		})
	}
}
