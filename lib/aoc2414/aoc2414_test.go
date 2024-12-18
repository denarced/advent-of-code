package aoc2414

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveSafetyFactor(t *testing.T) {
	run := func(name string, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, DeriveSafetyFactor(lines, 11, 7, 100))
		})
	}

	run("empty", []string{}, 0)
	run(
		"example",
		[]string{
			"p=0,4 v=3,-3",
			"p=6,3 v=-1,-3",
			"p=10,3 v=-1,2",
			"p=2,0 v=2,-1",
			"p=0,0 v=1,3",
			"p=3,0 v=-2,-2",
			"p=7,6 v=-1,-3",
			"p=3,0 v=-1,-2",
			"p=9,3 v=2,3",
			"p=7,3 v=-1,2",
			"p=2,4 v=2,-3",
			"p=9,5 v=-3,-3",
		},
		12)
}

func TestDeriveCoordinates(t *testing.T) {
	run := func(specs []int, width, height, steps, expectedX, expectedY int) {
		name := fmt.Sprintf("%v %dx%d %d steps", specs, width, height, steps)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)

			// EXERCISE
			x, y := deriveCoordinates(specs, width, height, steps)

			// VERIFY
			ass := assert.New(t)
			ass.Equal(expectedX, x, "X")
			ass.Equal(expectedY, y, "Y")
		})
	}

	run([]int{0, 0, 0, 0}, 5, 5, 0, 0, 0)
	run(
		[]int{2, 2, 2, 2},
		11, 7,
		1,
		4, 4)
	run(
		[]int{2, 2, 11, 7},
		11, 7,
		1,
		2, 2)
	run(
		[]int{2, 2, -3, -3},
		11, 7,
		1,
		10, 6)
}
