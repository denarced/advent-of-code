package aoc2318

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDig(t *testing.T) {
	run := func(name string, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			dug := Dig(lines)

			req.Equal(expected, dug)
		})
	}

	run("in.txt", gent.OrPanic2(inr.ReadPath("testdata/in.txt"))("read test data"), 62)
	run(
		"square",
		[]string{
			"R 2 ()",
			"D 2 ()",
			"L 2 ()",
			"U 2 ()",
		},
		9)
	run(
		//   ###
		//   # #
		// ### ###
		// #     #
		// ### ###
		//   # #
		//   ###
		"cross",
		[]string{
			"R 2 ()",
			"D 2 ()",
			"R 2 ()",
			"D 2 ()",
			"L 2 ()",
			"D 2 ()",
			"L 2 ()",
			"U 2 ()",
			"L 2 ()",
			"U 2 ()",
			"R 2 ()",
			"U 2 ()",
		},
		33)
	run(
		// ###       |
		// # # ###   |
		// # ### #   |
		// #     #   |
		// #  ####   |
		// #  #      |
		// #  ###### |
		// ###     # |
		//   ### ### |
		//     ###   |
		"duck",
		[]string{
			"R 2 ()",
			"D 2 ()",
			"R 2 ()",
			"U 1 ()",
			"R 2 ()",
			"D 3 ()",
			"L 3 ()",
			"D 2 ()",
			"R 5 ()",
			"D 2 ()",
			"L 2 ()",
			"D 1 ()",
			"L 2 ()",
			"U 1 ()",
			"L 2 ()",
			"U 1 ()",
			"L 2 ()",
			"U 7 ()",
		},
		90-2-5-10-3-2-6)
	run(
		// ######    |
		// #    ##   |
		// # ### ##  |
		// # # ## ## |
		// ###  #### |
		"republican",
		[]string{
			"R 5 ()",
			"D 1 ()",
			"R 1 ()",
			"D 1 ()",
			"R 1 ()",
			"D 1 ()",
			"R 1 ()",
			"D 1 ()",
			"L 3 ()",
			"U 1 ()",
			"L 1 ()",
			"U 1 ()",
			"L 2 ()",
			"D 2 ()",
			"L 2 ()",
			"U 4 ()",
		},
		36)
}

func TestToIndex(t *testing.T) {
	var tests = []struct {
		index    int
		length   int
		expected int
	}{
		{-4, 3, 2},
		{-3, 3, 0},
		{-2, 3, 1},
		{-1, 3, 2},
		{0, 3, 0},
		{1, 3, 1},
		{2, 3, 2},
		{3, 3, 0},
		{4, 3, 1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			require.Equal(t, tt.expected, toIndex(tt.index, tt.length))
		})
	}
}

func TestDeriveTurn(t *testing.T) {
	tests := []struct {
		name     string
		from     shared.Direction
		to       shared.Direction
		expected int
	}{
		{"east to north is left", shared.RealEast, shared.RealNorth, -1},
		{"east to south is right", shared.RealEast, shared.RealSouth, 1},
		{"east to east is straight", shared.RealEast, shared.RealEast, 0},
		{"east to west is straight", shared.RealEast, shared.RealWest, 0},
		{"south to east is left", shared.RealSouth, shared.RealEast, -1},
		{"south to west is right", shared.RealSouth, shared.RealWest, 1},
		{"south to south is straight", shared.RealSouth, shared.RealSouth, 0},
		{"south to north is straight", shared.RealSouth, shared.RealNorth, 0},
		{"west to north is right", shared.RealWest, shared.RealNorth, 1},
		{"west to south is left", shared.RealWest, shared.RealSouth, -1},
		{"west to west is straight", shared.RealWest, shared.RealWest, 0},
		{"west to east is straight", shared.RealWest, shared.RealEast, 0},
		{"north to west is left", shared.RealNorth, shared.RealWest, -1},
		{"north to east is right", shared.RealNorth, shared.RealEast, 1},
		{"north to north is straight", shared.RealNorth, shared.RealNorth, 0},
		{"north to south is straight", shared.RealNorth, shared.RealSouth, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instructions := []instruction{{dir: tt.from}, {dir: tt.to}}
			assert.Equal(t, tt.expected, deriveTurn(instructions, 0))
		})
	}
}

func TestDeriveSide(t *testing.T) {
	run := func(first, second shared.Loc, clockwise bool, expected shared.Loc) {
		name := fmt.Sprintf("%v -> %v (%s)", first, second, gent.Tri(clockwise, "⤾", "⤿"))
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			side := deriveSide(first, second, clockwise)

			// VERIFY
			req.Equal(expected, side)
		})
	}

	run(shared.Loc{}, shared.Loc{X: 1}, false, shared.Loc{Y: 1})
	run(shared.Loc{}, shared.Loc{X: 1}, true, shared.Loc{Y: -1})
	run(shared.Loc{}, shared.Loc{Y: 1}, false, shared.Loc{X: -1})
	run(shared.Loc{}, shared.Loc{Y: 1}, true, shared.Loc{X: 1})
	run(shared.Loc{}, shared.Loc{X: -1}, false, shared.Loc{Y: -1})
	run(shared.Loc{}, shared.Loc{X: -1}, true, shared.Loc{Y: 1})
	run(shared.Loc{}, shared.Loc{Y: -1}, false, shared.Loc{X: 1})
	run(shared.Loc{}, shared.Loc{Y: -1}, true, shared.Loc{X: -1})
}
