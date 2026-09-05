package aoc2214

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	run := func(infiniteFloor bool, expected int) {
		name := "no floor"
		if infiniteFloor {
			name = "infinite floor"
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "failed to read test data")

			req.Equal(expected, MeasureSand(lines, infiniteFloor))

		})
	}

	run(false, 24)
	run(true, 93)
}
