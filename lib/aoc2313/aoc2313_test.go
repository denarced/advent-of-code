package aoc2313

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumReflections(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err, "failed to read test data")

	// EXERCISE
	sum := SumReflections(lines)

	// VERIFY
	req.Equal(405, sum)
}

func TestDeriveMaxDistance(t *testing.T) {
	run := func(index, count, expected int) {
		name := fmt.Sprintf("%d+%d==%d", index, count, expected)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			maximus := deriveMaxDistance(index, count)

			// VERIFY
			req.Equal(expected, maximus)
		})
	}

	run(0, 2, 1)
	run(1, 3, 1)
	run(1, 4, 2)
}

func TestCheckReflection(t *testing.T) {
	lines := []string{
		".##",
		"#..",
		"#..",
		".##",
	}

	run := func(row, distance, maxDistance int, expected bool) {
		name := fmt.Sprintf("%d-%d-%d", row, distance, maxDistance)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			req.Equal(
				expected,
				checkHorizontalReflection(lines, row, distance, maxDistance))
		})
	}

	run(1, 1, 2, true)
	run(2, 1, 1, false)
}
