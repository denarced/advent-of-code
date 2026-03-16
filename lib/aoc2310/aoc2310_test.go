package aoc2310

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountSteps(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err)
	req.Equal(4, CountSteps(lines))
}

func TestFindDirections(t *testing.T) {
	run := func(name string, lines []string, expected []shared.Direction) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			brd := shared.NewBoard(lines)
			start := brd.FindOrDie('S')
			req.Equal("S", string(brd.GetOrDie(start)))
			dirs := findDirections(brd, start)
			req.ElementsMatch(expected, dirs)
		})
	}

	run(
		"-",
		[]string{
			"..-..",
			".-S-.",
			"..-..",
		},
		[]shared.Direction{shared.RealWest, shared.RealEast})
	run(
		"|",
		[]string{
			"..|..",
			".|S|.",
			"..|..",
		},
		[]shared.Direction{shared.RealNorth, shared.RealSouth})
	run(
		"L",
		[]string{
			"..L..",
			".LSL.",
			"..L..",
		},
		[]shared.Direction{shared.RealWest, shared.RealSouth})
	run(
		"J",
		[]string{
			"..J..",
			".JSJ.",
			"..J..",
		},
		[]shared.Direction{shared.RealEast, shared.RealSouth})
	run(
		"7",
		[]string{
			"..7..",
			".7S7.",
			"..7..",
		},
		[]shared.Direction{shared.RealNorth, shared.RealEast})
	run(
		"F",
		[]string{
			"..F..",
			".FSF.",
			"..F..",
		},
		[]shared.Direction{shared.RealNorth, shared.RealWest})
}

func TestStep(t *testing.T) {
	run := func(name string, aWalker walker, expected walker) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			brd := shared.NewBoard([]string{
				"F---7",
				"L7S-J",
				"-J...",
			})

			// EXERCISE
			result := step(brd, aWalker)

			// VERIFY
			req.Equal(expected, result)
		})
	}

	var tests = []struct {
		name     string
		current  walker
		expected walker
	}{
		{
			"S -> -",
			walker{loc: shared.Loc{X: 2, Y: 1}, dir: shared.RealEast},
			walker{loc: shared.Loc{X: 3, Y: 1}, dir: shared.RealEast},
		},
		{
			"F -> L",
			walker{loc: shared.Loc{X: 0, Y: 2}, dir: shared.RealSouth},
			walker{loc: shared.Loc{X: 0, Y: 1}, dir: shared.RealEast},
		},
		{
			"L -> 7",
			walker{loc: shared.Loc{X: 0, Y: 1}, dir: shared.RealEast},
			walker{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealSouth},
		},
		{
			"7 -> J",
			walker{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealSouth},
			walker{loc: shared.Loc{X: 1}, dir: shared.RealWest},
		},
		{
			"- -> F",
			walker{loc: shared.Loc{X: 1, Y: 2}, dir: shared.RealWest},
			walker{loc: shared.Loc{Y: 2}, dir: shared.RealSouth},
		},
	}
	for _, tt := range tests {
		run(tt.name, tt.current, tt.expected)
	}
}
