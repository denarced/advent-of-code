package aoc2208

import (
	"path/filepath"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveVisibleTreeCount(t *testing.T) {
	run := func(filen string, expected int) {
		t.Run(filen, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath(filepath.Join("testdata", filen))
			req.NoError(err, "failed to read test data")

			// EXERCISE & VERIFY
			req.Equal(expected, DeriveVisibleTreeCount(lines))
		})
	}

	run("in.txt", 21)
	run("in001.txt", 20)
	run("in002.txt", 23)
	run("in003.txt", 18)
	run("in004.txt", 19)
	run("in005.txt", 24)
}

func TestDeriveBestSpot(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	// EXERCISE & VERIFY
	req.Equal(8, DeriveBestSpot(lines))
}
