package aoc2319

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumRatings(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err, "failed to read test data")

	// EXERCISE
	sum := SumRatings(lines)

	// VERIFY
	req.Equal(19_114, sum)
}

func TestParseWorkflows(t *testing.T) {
	req := require.New(t)

	// EXERCISE
	workflows := parseWorkflows([]string{
		"px{a<2006:qkq,m>2090:A,rfg}",
	})

	// VERIFY
	req.Equal(
		map[string]workflow{
			"px": {
				name: "px",
				specs: []spec{
					{attr: "a", dest: "qkq", value: 2006, less: true},
					{attr: "m", dest: "A", value: 2090},
					{dest: "rfg", endComplete: true},
				},
			},
		},
		workflows)
}

func TestParseParts(t *testing.T) {
	req := require.New(t)

	// EXERCISE
	parts := parseParts([]string{
		"{x=787,m=2655,a=1222,s=2876}",
	})

	// VERIFY
	req.Equal(
		[]part{
			map[string]int{"x": 787, "m": 2655, "a": 1222, "s": 2876},
		},
		parts)
}
