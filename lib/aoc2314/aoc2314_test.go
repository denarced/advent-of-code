package aoc2314

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountTotalLoad(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err)
	req.Equal(136, CountTotalLoad(lines, 0))
}

func BenchmarkCountTotalLoad(b *testing.B) {
	shared.InitNullLogging()
	req := require.New(b)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err)

	for range b.N {
		CountTotalLoad(lines, 0)
	}
}

func TestMoveRocks(t *testing.T) {
	run := func(from, to []string, direction shared.Direction) {
		t.Run(fmt.Sprint(direction), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			brd := shared.NewBoard(from)
			// EXERCISE
			moveRocks(brd, direction)

			// VERIFY
			req.Equal(to, brd.GetLines())
		})
	}

	run(
		[]string{
			".OO.",
			"O..#",
			"#...",
			"OO..",
		},
		[]string{
			"OOO.",
			".O.#",
			"#...",
			"O...",
		},
		shared.RealNorth)
	run(
		[]string{
			".OO.",
			"O..#",
			"#...",
			"OO..",
		},
		[]string{
			"..OO",
			"..O#",
			"#...",
			"..OO",
		},
		shared.RealEast)
	run(
		[]string{
			".OO.",
			"O..#",
			"#...",
			"OO..",
		},
		[]string{
			"....",
			"O..#",
			"#O..",
			"OOO.",
		},
		shared.RealSouth)
	run(
		[]string{
			".OO.",
			"O..#",
			"#...",
			"OO..",
		},
		[]string{
			"OO..",
			"O..#",
			"#...",
			"OO..",
		},
		shared.RealWest)
}
