package aoc2217

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeasureTowerHeight(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	req.Equal(3068, MeasureTowerHeight(lines))
}

func TestMovePiece(t *testing.T) {
	run := func(piece, expected []shared.Loc, expectedOk bool) {
		shared.InitTestLogging(t)
		ass := assert.New(t)

		// EXERCISE
		dropped, ok := movePiece(nil, piece, shared.Loc{Y: -1})

		// VERIFY
		ass.Equal(expectedOk, ok)
		ass.Equal(expected, dropped)
	}

	run(
		[]shared.Loc{
			{X: 2, Y: 11},
			{X: 4, Y: 11},
			{X: 3, Y: 10},
			{X: 2, Y: 9},
			{X: 4, Y: 9},
		},
		[]shared.Loc{
			{X: 2, Y: 10},
			{X: 4, Y: 10},
			{X: 3, Y: 9},
			{X: 2, Y: 8},
			{X: 4, Y: 8},
		},
		true)
	run(
		[]shared.Loc{{X: 0, Y: 0}, {X: 1, Y: 0}},
		[]shared.Loc{{X: 0, Y: 0}, {X: 1, Y: 0}},
		false)
}
