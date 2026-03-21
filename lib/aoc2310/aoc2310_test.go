package aoc2310

import (
	"path/filepath"
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

func TestMovement(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	brd := shared.NewBoard([]string{
		".FS7.",
		".|.|.",
		".L-J.",
	})

	gatherSteps := func(dir shared.Direction) []shared.Loc {
		start := brd.FindOrDie('S')
		aWalker := walker{
			loc: start,
			dir: dir,
		}
		steps := []shared.Loc{start}
		for {
			next := step(brd, aWalker)
			if next.loc == start {
				break
			}
			steps = append(steps, next.loc)
			aWalker = next
		}
		return steps
	}

	req.Equal(
		[]shared.Loc{
			{X: 2, Y: 2},
			{X: 3, Y: 2},
			{X: 3, Y: 1},
			{X: 3, Y: 0},
			{X: 2, Y: 0},
			{X: 1, Y: 0},
			{X: 1, Y: 1},
			{X: 1, Y: 2},
		},
		gatherSteps(shared.RealEast))
	req.Equal(
		[]shared.Loc{
			{X: 2, Y: 2},
			{X: 1, Y: 2},
			{X: 1, Y: 1},
			{X: 1, Y: 0},
			{X: 2, Y: 0},
			{X: 3, Y: 0},
			{X: 3, Y: 1},
			{X: 3, Y: 2},
		},
		gatherSteps(shared.RealWest))
}

func TestFindCrackCount(t *testing.T) {
	run := func(filen string, expected int) {
		t.Run(filen, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := inr.ReadPath(filepath.Join("testdata", filen))
			req.NoError(err)

			// EXERCISE
			count := FindCrackCount(lines)

			// VERIFY
			if expected >= 0 {
				req.Equal(expected, count)
			}
		})
	}

	run("in2.txt", 4)
	run("in3.txt", 4)
	run("in4.txt", 10)
	run("in5.txt", 8)
	run("in6.txt", 1)
}
