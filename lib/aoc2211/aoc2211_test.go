package aoc2211

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveMonkeyBusiness(t *testing.T) {
	run := func(roundCount int, worries bool, expected int) {
		worryStr := "no worries"
		if worries {
			worryStr = "worries"
		}
		name := fmt.Sprintf("%d-%s-%d", roundCount, worryStr, expected)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "read test data")

			req.Equal(expected, DeriveMonkeyBusiness(lines, roundCount, worries))
		})
	}

	run(20, false, 10_605)
	run(10_000, true, 2_713_310_158)
}
