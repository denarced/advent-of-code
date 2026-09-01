package aoc2212

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveFewestSteps(t *testing.T) {
	run := func(hikeDir HikeStrategy, expected int) {
		name := "start-to-end"
		if hikeDir == LowToEnd {
			name = "low-to-end"
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "failed to read test data")

			req.Equal(expected, DeriveFewestSteps(lines, hikeDir))
		})
	}

	run(StartToEnd, 31)
	run(LowToEnd, 29)
}

func TestIsTooHigh(t *testing.T) {
	run := func(from, to rune, expected bool) {
		t.Run(fmt.Sprintf("%c-%c", from, to), func(t *testing.T) {
			require.Equal(t, expected, isTooHigh(from, to))
		})
	}

	run('b', 'a', false)
	run('a', 'a', false)
	run('a', 'b', false)
	run('a', 'c', true)
}
