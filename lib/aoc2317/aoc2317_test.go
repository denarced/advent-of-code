package aoc2317

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveLeastHeatLoss(t *testing.T) {
	run := func(minJump, maxJump, expected int) {
		name := fmt.Sprintf("%d - %d", minJump, maxJump)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err)

			req.Equal(expected, DeriveLeastHeatLoss(lines, minJump, maxJump))
		})
	}

	run(1, 3, 102)
	run(4, 10, 94)
}

func BenchmarkDeriveLeastHeatLoss(b *testing.B) {
	shared.InitNullLogging()
	lines, _ := inr.ReadPath("testdata/in.txt")

	for range b.N {
		DeriveLeastHeatLoss(lines, 1, 3)
	}
}

func TestBetween(t *testing.T) {
	run := func(from, to shared.Loc, expected []shared.Loc) {
		name := fmt.Sprintf("%v -> %v", from, to)
		t.Run(name, func(t *testing.T) {
			req := require.New(t)
			var steps []shared.Loc
			doInBetween(from, to, func(loc shared.Loc) {
				steps = append(steps, loc)
			})
			req.Equal(expected, steps)
		})
	}

	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 1, Y: 3}, []shared.Loc{{X: 2, Y: 3}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 2, Y: 3}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 3}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 4, Y: 3}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 5, Y: 3}, []shared.Loc{{X: 4, Y: 3}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 5, Y: 3}, []shared.Loc{{X: 4, Y: 3}})

	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 1}, []shared.Loc{{X: 3, Y: 2}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 2}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 4}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 5}, []shared.Loc{{X: 3, Y: 4}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 6}, []shared.Loc{{X: 3, Y: 4}, {X: 3, Y: 5}})
}

func TestDeriveNextHops(t *testing.T) {
	run := func(name string, latest hop, expected []hop) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			brd := shared.NewBoard([]string{
				"012345",
				"012345",
				"012345",
			})
			aRunner := &runner{latestHop: latest}
			hops := make([]hop, 0, 3)
			// EXERCISE
			hops = deriveNextHops(brd, aRunner, hops, 1, 3)

			// VERIFY
			req.Equal(expected, hops)
		})
	}

	run(
		"happy path X",
		hop{loc: shared.Loc{X: 4, Y: 1}, dir: shared.RealEast},
		[]hop{
			{loc: shared.Loc{X: 5, Y: 1}, dir: shared.RealNorth},
			{loc: shared.Loc{X: 5, Y: 1}, dir: shared.RealSouth},
			{loc: shared.Loc{X: 5, Y: 1}, dir: shared.RealEast},
		})
	run(
		"happy path Y",
		hop{loc: shared.Loc{X: 3, Y: 1}, dir: shared.RealNorth},
		[]hop{
			{loc: shared.Loc{X: 3, Y: 2}, dir: shared.RealWest},
			{loc: shared.Loc{X: 3, Y: 2}, dir: shared.RealEast},
			{loc: shared.Loc{X: 3, Y: 2}, dir: shared.RealNorth},
		})
	run("dead end X", hop{loc: shared.Loc{}, dir: shared.RealWest}, []hop{})
	run("dead end Y", hop{loc: shared.Loc{X: 5, Y: 2}, dir: shared.RealEast}, []hop{})
	run(
		"long",
		hop{
			loc: shared.Loc{X: 0, Y: 1},
			dir: shared.RealEast,
		},
		[]hop{
			{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealNorth},
			{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealSouth},

			{loc: shared.Loc{X: 2, Y: 1}, dir: shared.RealNorth},
			{loc: shared.Loc{X: 2, Y: 1}, dir: shared.RealSouth},

			{loc: shared.Loc{X: 3, Y: 1}, dir: shared.RealNorth},
			{loc: shared.Loc{X: 3, Y: 1}, dir: shared.RealSouth},

			{loc: shared.Loc{X: 3, Y: 1}, dir: shared.RealEast},
		})
}
