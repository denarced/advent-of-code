package aoc2413

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestDeriveFewestTokens(t *testing.T) {
	run := func(name string, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, DeriveFewestTokens(lines, false))
		})
	}

	run("empty", []string{}, 0)
}

func TestParseButton(t *testing.T) {
	shared.InitTestLogging(t)
	btn := parseButton("Button A: X+94, Y-34")
	require.Equal(t, button{name: "A", loc: shared.Loc{X: 94, Y: -34}}, btn)
}

func TestParsePrize(t *testing.T) {
	shared.InitTestLogging(t)
	btn := parsePrize("Prize: X=8400, Y=5400", false)
	require.Equal(t, shared.Loc{X: 8400, Y: 5400}, btn)
}

func TestParseMachines(t *testing.T) {
	shared.InitTestLogging(t)
	machines := parseMachines([]string{
		"Button A: X+94, Y+34",
		"Button B: X+22, Y+67",
		"Prize: X=8400, Y=5400",
		"",
		"Button A: X+26, Y+66",
		"Button B: X+67, Y-21",
		"Prize: X=12748, Y=12176",
	}, false)
	require.Equal(
		t,
		[]machine{
			{
				a:     button{name: "A", loc: shared.Loc{X: 94, Y: 34}},
				b:     button{name: "B", loc: shared.Loc{X: 22, Y: 67}},
				prize: shared.Loc{X: 8400, Y: 5400},
			},
			{
				a:     button{name: "A", loc: shared.Loc{X: 26, Y: 66}},
				b:     button{name: "B", loc: shared.Loc{X: 67, Y: -21}},
				prize: shared.Loc{X: 12748, Y: 12176},
			},
		},
		machines)
}

func TestDeriveCheapest(t *testing.T) {
	run := func(name string, a, b button, prize shared.Loc, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := deriveCheapest(a, b, prize)
			req.Equal(expected, actual)
		})
	}

	run(
		"example 1",
		button{loc: shared.Loc{X: 94, Y: 34}},
		button{loc: shared.Loc{X: 22, Y: 67}},
		shared.Loc{X: 8400, Y: 5400},
		80*3+40)
	run(
		"example 2",
		button{loc: shared.Loc{X: 26, Y: 66}},
		button{loc: shared.Loc{X: 67, Y: 21}},
		shared.Loc{X: 12748, Y: 12176},
		-1)
}
