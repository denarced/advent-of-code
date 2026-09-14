package aoc2216

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveMaximumPressure(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	req.Equal(1651, DeriveMaximumPressure(lines, 30, true))
}

func TestParseLine(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	parsed := parseLine("Valve AA has flow rate=0; tunnels lead to valves DD, II, BB")

	req.Equal(
		valve{
			name:           "AA",
			flowRate:       0,
			connectedNames: []string{"DD", "II", "BB"},
		},
		*parsed)
}
