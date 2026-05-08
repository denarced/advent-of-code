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
			require.Equal(t, expected, Dig(lines, false))
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

	// 6 ####
	// 5 #..#
	// 4 ##.#
	// 3 .#.#
	// 2 ##.#
	// 1 #..#
	// 0 ####
	//   0123
	// Problem because code failed to include 1x3.
	run(
		"problem",
		[]string{
			"U 2 ()",
			"R 1 ()",
			"U 2 ()",
			"L 1 ()",
			"U 2 ()",
			"R 3 ()",
			"D 6 ()",
			"L 3 ()",
		},
		27)

	run(
		"3x3 square",
		[]string{
			"R 2 ()",
			"D 2 ()",
			"L 2 ()",
			"U 2 ()",
		},
		9)

	//  1 .##
	//  0 ###
	// -1 #.#
	// -2 ###
	//    012
	run(
		"simple complex",
		[]string{
			"R 1 ()",
			"U 1 ()",
			"R 1 ()",
			"D 3 ()",
			"L 2 ()",
			"U 2 ()",
		},
		11)

	//  1 .####..
	//  0 ##..#..
	// -1 #...###
	// -2 ###...#
	// -3 ..#####
	//    0123456
	run(
		"complex",
		[]string{
			"R 1 ()",
			"U 1 ()",
			"R 3 ()",
			"D 2 ()",
			"R 2 ()",
			"D 2 ()",
			"L 4 ()",
			"U 1 ()",
			"L 2 ()",
			"U 2 ()",
		},
		28)
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

func TestParseLine(t *testing.T) {
	run := func(line string, expected, expectedWithMagic instruction) {
		for _, magic := range []bool{true, false} {
			t.Run(
				fmt.Sprintf("%s - %s", line, gent.Tri(magic, "magic", "no-magic")),
				func(t *testing.T) {
					shared.InitTestLogging(t)
					req := require.New(t)

					// EXERCISE
					parsed := parseLine(line, magic)

					// VERIFY
					if magic {
						req.Equal(expectedWithMagic, parsed)
					} else {
						req.Equal(expected, parsed)
					}
				},
			)
		}
	}

	run(
		"R 2 (#70c710)",
		instruction{
			dir:       shared.RealEast,
			stepCount: 2,
			delta:     shared.Loc{X: 2, Y: 0},
		},
		instruction{
			dir:       shared.RealEast,
			stepCount: 461937,
			delta:     shared.Loc{X: 461937},
		})
}

func TestIsAbove(t *testing.T) {
	var tests = []struct {
		name     string
		dir      shared.Direction
		expected bool
	}{
		{"east", shared.RealEast, true},
		{"west", shared.RealWest, false},
	}
	for _, tt := range tests {
		for _, clockwise := range []bool{false, true} {
			suffix := gent.Tri(clockwise, " clockwise", " counterclockwise")
			if clockwise {
				tt.expected = !tt.expected
			}
			t.Run(tt.name+suffix, func(t *testing.T) {
				require.Equal(t, tt.expected, isAbove(tt.dir, clockwise))
			})
		}
	}
}

func TestGatherRoutes(t *testing.T) {
	run := func(
		name string,
		instructions []instruction,
		clockwise bool,
		expectedRoutes []route,
		expectedXcoords []int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			routes, xCoords := gatherRoutes(instructions, clockwise)

			// VERIFY
			req.Equal(stringify(expectedRoutes), stringify(routes), "routes")
			req.Equal(expectedXcoords, xCoords, "x coordinates")
		})
	}

	run(
		"square",
		[]instruction{
			// 0,0 -> 0,-3
			{dir: shared.RealSouth, stepCount: 3, delta: shared.Loc{Y: -3}},
			// 0,-3 -> -3,-3
			{dir: shared.RealWest, stepCount: 3, delta: shared.Loc{X: -3}},
			// -3,-3 -> -3,0
			{dir: shared.RealNorth, stepCount: 3, delta: shared.Loc{Y: 3}},
			// -3,0 -> 0,0
			{dir: shared.RealEast, stepCount: 3, delta: shared.Loc{X: 3}},
		},
		true,
		[]route{
			{y: -1, above: false},
			{y: -2, above: true},
			{from: -3, to: 0, y: -3, above: true},
			{from: -3, to: -3, y: -1, above: false},
			{from: -3, to: -3, y: -2, above: true},
			{from: -3, to: 0, y: 0, above: false},
		},
		[]int{-3, 0})

	run(
		"rectangle",
		[]instruction{
			// 0,0 -> -3,0
			{dir: shared.RealWest, stepCount: 3, delta: shared.Loc{X: -3}},
			// -3,0 -> -3,-2
			{dir: shared.RealSouth, stepCount: 2, delta: shared.Loc{Y: -2}},
			// -3,-2 -> 1,-2
			{dir: shared.RealEast, stepCount: 4, delta: shared.Loc{X: 4}},
			// 1,-2 -> 1,0
			{dir: shared.RealNorth, stepCount: 2, delta: shared.Loc{Y: 2}},
			// 1,0 -> 0,0
			{dir: shared.RealWest, stepCount: 1, delta: shared.Loc{X: -1}},
		},
		false,
		[]route{
			{from: -3, to: 0, y: 0, above: false},
			{from: -3, to: -3, y: -1},
			{from: -3, to: -3, y: -1, above: true},
			{from: -3, to: 1, y: -2, above: true},
			{from: 1, to: 1, y: -1},
			{from: 1, to: 1, y: -1, above: true},
			{from: 0, to: 1, y: 0, above: false},
		},
		[]int{-3, 0, 1})
}

func TestExpandAllRoutes(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	// 0 t---0-->--t
	// 1 |         |
	// 2 |         |
	// 3 t---------t
	//   43210123456
	routes, xCoords := gatherRoutes(
		[]instruction{
			{
				dir:       shared.RealEast,
				stepCount: 6,
				delta:     shared.Loc{X: 6},
			},
			{
				dir:       shared.RealSouth,
				stepCount: 3,
				delta:     shared.Loc{Y: -3},
			},
			{
				dir:       shared.RealWest,
				stepCount: 10,
				delta:     shared.Loc{X: -10},
			},
			{
				dir:       shared.RealNorth,
				stepCount: 3,
				delta:     shared.Loc{Y: 3},
			},
			{
				dir:       shared.RealEast,
				stepCount: 4,
				delta:     shared.Loc{X: 4},
			},
		},
		true)
	// EXERCISE
	expanded := expandAllRoutes(routes, xCoords)

	// VERIFY
	a := []route{
		{from: 0, to: 0},
		{from: 1, to: 5},
		{from: 6, to: 6},
		{from: 6, to: 6, y: -1},
		{from: 6, to: 6, y: -2, above: true},
		{from: -4, to: -4, y: -3, above: true},
		{from: -3, to: -1, y: -3, above: true},
		{y: -3, above: true},
		{from: 1, to: 5, y: -3, above: true},
		{from: 6, to: 6, y: -3, above: true},
		{from: -4, to: -4, y: -1},
		{from: -4, to: -4, y: -2, above: true},
		{from: -4, to: -4},
		{from: -3, to: -1},
		{from: 0, to: 0},
	}
	req.Equal(
		stringify(a),
		stringify(expanded))
}

func stringify[S ~[]T, T any](s S) []string {
	strs := make([]string, len(s))
	for i, each := range s {
		strs[i] = fmt.Sprint(each)
	}
	return strs
}
