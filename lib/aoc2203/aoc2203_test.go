package aoc2203

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumPriorities(t *testing.T) {
	run := func(groupSearch bool, expected int) {
		name := "compartments"
		if groupSearch {
			name = "group search"
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err)

			req.Equal(expected, SumPriorities(lines, groupSearch))
		})
	}

	run(false, 157)
	run(true, 70)
}
