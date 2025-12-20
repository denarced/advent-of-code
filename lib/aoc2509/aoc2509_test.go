package aoc2509

import (
	"fmt"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveBiggestRectangle(t *testing.T) {
	readLines := func(req *require.Assertions) []string {
		lines, err := inr.ReadPath("testdata/in.txt")
		req.NoError(err, "read test data")
		return lines
	}
	t.Run("!redGreen", func(t *testing.T) {
		shared.InitNullLogging()
		req := require.New(t)
		require.Equal(t, 50, DeriveBiggestRectangle(readLines(req), false))
	})
	t.Run("redGreen", func(t *testing.T) {
		shared.InitNullLogging()
		req := require.New(t)
		require.Equal(t, 24, DeriveBiggestRectangle(readLines(req), true))
	})
}

func TestCalculateArea(t *testing.T) {
	var tests = []struct {
		first    [2]int
		second   [2]int
		expected int
	}{
		{[2]int{2, 2}, [2]int{4, 4}, 9},
		{[2]int{9, 9}, [2]int{9, 9}, 1},
		{[2]int{-2, 0}, [2]int{2, 2}, 15},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%dx%d * %dx%d", tt.first[0], tt.first[1], tt.second[0], tt.second[1])
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.expected, calculateArea(tt.first, tt.second))
		})
	}
}

func TestDeriveClockwise(t *testing.T) {
	var tests = []struct {
		coords   [][2]int
		expected bool
	}{
		{
			[][2]int{
				{0, 0},
				{0, 1},
				{-1, 1},
				{-1, 0},
			},
			false,
		},
		{
			[][2]int{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 0},
			},
			true,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%dx%d -> %v", tt.coords[0][0], tt.coords[0][1], tt.expected)
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.expected, deriveClockwise(tt.coords))
		})
	}
}

func TestDeriveTurn(t *testing.T) {
	var tests = []struct {
		name     string
		lines    []string
		expected int
	}{
		{
			"right-up",
			[]string{
				"   | 2 ",
				"---+01-",
				"   |   ",
			},
			-1,
		},
		{
			"right-down",
			[]string{
				"   |   ",
				"-01+---",
				"  2|   ",
			},
			1,
		},
		{
			"left-up",
			[]string{
				"   |   ",
				"--2+---",
				"  10   ",
			},
			1,
		},
		{
			"left-down",
			[]string{
				"  10   ",
				"--2+---",
				"   |   ",
			},
			-1,
		},
		{
			"up-left",
			[]string{
				"   |   ",
				"---+-21",
				"   |  0",
			},
			-1,
		},
		{
			"up-right",
			[]string{
				"12 |   ",
				"0--+---",
				"   |   ",
			},
			1,
		},
		{
			"down-left",
			[]string{
				"   |  0",
				"---+-21",
				"   |   ",
			},
			1,
		},
		{
			"down-right",
			[]string{
				"   |   ",
				"---+-0-",
				"   | 12",
			},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.InitTestLogging(t)
			coords := asciiToCoordinates(tt.lines)
			require.Equal(t, tt.expected, deriveTurn(coords[0], coords[1], coords[2]))
		})
	}
}

func TestAsciiToCoordinates(t *testing.T) {
	shared.InitTestLogging(t)
	require.Equal(
		t,
		[][2]int{{0, 0}, {3, 0}, {3, 1}, {2, 1}, {2, 2}},
		asciiToCoordinates([]string{
			"     |    ",
			"     | 4  ",
			"     | 32 ",
			"-----0--1-",
			"     |    ",
			"     |    ",
			"     |    ",
		}),
	)
	require.Equal(
		t,
		[][2]int{{0, -1}, {-1, -1}, {-2, -1}, {-2, 0}},
		asciiToCoordinates([]string{
			"   |   ",
			"-3-+---",
			" 210   ",
		}),
	)
}

func asciiToCoordinates(table []string) [][2]int {
	var zeroX, zeroY *int
	for y, line := range table {
		if zeroX != nil && zeroY != nil {
			break
		}
		if zeroX == nil {
			if x := strings.Index(line, "|"); x >= 0 {
				zeroX = &x
			}
		}
		if zeroY == nil {
			if strings.Contains(line, "-") {
				zeroY = &y
			}
		}
	}
	coordinates := make([][2]int, 10)
	var maxCoord int
	for y, line := range table {
		realY := deriveRealY(y, len(table), *zeroY)
		for x, c := range line {
			if '0' <= c && c <= '9' {
				coordIndex := int(c - '0')
				maxCoord = max(maxCoord, coordIndex)
				coordinates[coordIndex] = [2]int{x - *zeroX, realY}
			}
		}
	}
	result := make([][2]int, maxCoord+1)
	copy(result, coordinates[:maxCoord+1])
	return result
}

func deriveRealY(y, size, zeroY int) int {
	result := abs(y - size + 1)
	return result - abs(zeroY-size+1)
}

func TestDeriveRealY(t *testing.T) {
	var tests = []struct {
		y        int
		size     int
		zeroY    int
		expected int
	}{
		{0, 1, 0, 0},
		{0, 2, 1, 1},
		{8, 20, 19, 11},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%+v", tt)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, tt.expected, deriveRealY(tt.y, tt.size, tt.zeroY))
		})
	}
}

func TestFillLine(t *testing.T) {
	var tests = []struct {
		name       string
		start, end [2]int
		expected   [][2]int
	}{
		{
			"right",
			[2]int{0, 0},
			[2]int{2, 0},
			[][2]int{{0, 0}, {1, 0}, {2, 0}},
		},
		{
			"up",
			[2]int{0, 0},
			[2]int{0, 1},
			[][2]int{{0, 0}, {0, 1}},
		},
		{
			"left",
			[2]int{0, 0},
			[2]int{-2, 0},
			[][2]int{{0, 0}, {-1, 0}, {-2, 0}},
		},
		{
			"down",
			[2]int{0, 0},
			[2]int{0, -1},
			[][2]int{{0, 0}, {0, -1}},
		},
		{
			"down-left",
			[2]int{1, 1},
			[2]int{-1, -1},
			[][2]int{{1, 1}, {0, 0}, {-1, -1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := fillLine(tt.start, tt.end)
			require.ElementsMatch(t, tt.expected, line)
		})
	}
}

func TestScreen(t *testing.T) {
	require.Equal(
		t,
		[][2][2]int{
			{{0, 0}, {0, 1}},
			{{0, 1}, {0, 2}},
			{{0, 2}, {0, 0}},
		},
		screen([][2]int{{0, 0}, {0, 1}, {0, 2}}),
	)
}

func TestModulo(t *testing.T) {
	var tests = []struct {
		i, n, expected int
	}{
		{0, 1, 0},
		{1, 1, 0},
		{0, 5, 0},
		{1, 5, 1},
		{4, 5, 4},
		{5, 5, 0},
		{-1, 5, 4},
		{-9, 5, 1},
		{-10, 5, 0},
		{-11, 5, 4},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%d-%d", tt.i, tt.n)
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.expected, modulo(tt.i, tt.n))
		})
	}
}

func TestScreenCorner(t *testing.T) {
	exercise := func(f func() (corner, bool)) []corner {
		trips := []corner{}
		for {
			corn, done := f()
			if done {
				break
			}
			trips = append(trips, corn)
		}
		return trips
	}

	assert.Equal(
		t,
		[]corner{
			{4, 0, 1},
			{0, 1, 2},
			{1, 2, 3},
			{2, 3, 4},
			{3, 4, 0},
		},
		exercise(screenCorner([]int{3, 4, 5, 6, 7}, true)))
	assert.Equal(
		t,
		[]corner{
			{0, 3, 2},
			{3, 2, 1},
			{2, 1, 0},
			{1, 0, 3},
		},
		exercise(screenCorner([]string{"", "", "", ""}, false)))
}

func TestDeriveBiggest(t *testing.T) {
	run := func(name string, coords [][2]int, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, deriveBiggest(coords, true))
		})
	}

	var tests = []struct {
		name     string
		coords   [][2]int
		expected int
	}{
		{
			"square",
			asciiToCoordinates([]string{
				"2 | 1",
				"  |  ",
				"--+--",
				"  |  ",
				"3 | 0",
			}),
			25,
		},
		{
			"L",
			asciiToCoordinates([]string{
				"1 2|  ",
				"   |  ",
				"--3+-4",
				"   |  ",
				"0  | 5",
			}),
			18,
		},
	}
	for _, tt := range tests {
		run(tt.name, tt.coords, tt.expected)
	}
}
