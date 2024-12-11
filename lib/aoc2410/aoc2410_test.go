package aoc2410

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestDeriveSumOfTrailheadScores(t *testing.T) {
	run := func(name string, lines []string, ratings bool, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := DeriveSumOfTrailheadScores(lines, ratings)
			req.Equal(expected, actual)
		})
	}

	run("empty wo ratings", []string{}, false, 0)
	{
		lines := []string{
			"0342",
			"1251",
			"8967",
			"6798",
		}
		run("one wo ratings", lines, false, 1)
		run("one with ratings", lines, true, 1)
	}
	{
		lines := []string{
			"0123",
			"9234",
			"8765",
			"9876",
		}
		run("two wo ratings", lines, false, 2)
		run("two with ratings", lines, true, 15)
	}
	{
		lines := []string{
			"89010123",
			"78121874",
			"87430965",
			"96549874",
			"45678903",
			"32019012",
			"01329801",
			"10456732",
		}
		run("example wo ratings", lines, false, 36)
		run("example with ratings", lines, true, 81)
	}
}
