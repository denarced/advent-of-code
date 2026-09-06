package aoc2215

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestFindNoGoForBeacons(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	req.Equal(26, FindNoGoForBeacons(lines, 10))
}

func TestDeriveRange(t *testing.T) {
	run := func(centerY, overlap, expectedFrom, expectedTo int) {
		name := fmt.Sprintf("centerY:%d overlap:%d", centerY, overlap)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			from, to := deriveRange(centerY, overlap)

			// VERIFY
			req.Equal([2]int{from, to}, [2]int{expectedFrom, expectedTo})
		})
	}

	run(10, 0, 10, 10)
	run(10, 1, 9, 11)
	run(10, 2, 8, 12)
}
