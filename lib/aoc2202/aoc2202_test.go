package aoc2202

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveTotalScore(t *testing.T) {
	run := func(declareEnd bool, expected int) {
		name := "what to play"
		if declareEnd {
			name = "declare end"
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "failed to read test data")

			req.Equal(expected, DeriveTotalScore(lines, declareEnd))
		})
	}
	run(false, 15)
	run(true, 12)
}
