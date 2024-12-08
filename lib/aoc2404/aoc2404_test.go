package aoc2404

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestCountInTable(t *testing.T) {
	run := func(name string, table []string, word string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			count := CountInTable(table, word)
			req.Equal(expected, count)
		})
	}

	run("empty", []string{}, "JOB", 0)
	run("one horizontal", []string{"..XMAS.."}, "XMAS", 1)
	run(
		"happy path",
		[]string{
			"XS..S..S..S.SX",
			"S....A.A.A...S",
			"......MMM.....",
			"....SAMXMAS...",
			"......MMM.....",
			".....A.A.A....",
			"S...S..S..S..S",
			"XS..........SX",
		},
		"XMAS",
		8,
	)
}

func TestCountWordCrosses(t *testing.T) {
	run := func(name string, table []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, CountWordCrosses(table, "MAS"))
		})
	}

	run("empty", []string{}, 0)
	run(
		"half x",
		[]string{
			"......",
			".M....",
			"..A...",
			".M.S..",
		},
		0,
	)
	run(
		"happy path",
		[]string{
			".M.S.M",
			"..A.A.",
			".M.S.M",
			"......",
		},
		2,
	)
}
