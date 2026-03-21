package aoc2311

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestSumDistances(t *testing.T) {
	run := func(expected, multiplier int) {
		t.Run(strconv.Itoa(multiplier), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err)
			req.Equal(expected, SumDistances(lines, multiplier))
		})
	}
	run(374, 2)
	run(1030, 10)
	run(8410, 100)
}

func TestExpandLocations(t *testing.T) {
	shared.InitTestLogging(t)

	// EXERCISE
	expanded := expandLocations(
		extractGalaxyLocations(
			shared.NewBoard([]string{
				".#....",
				"......",
				"....#.",
			})),
		2)

	toString := func(v []shared.Loc) []string {
		strs := gent.Map(v, func(loc shared.Loc) string {
			return fmt.Sprint(loc)
		})
		slices.Sort(strs)
		return strs
	}
	// VERIFY
	require.Equal(
		t,
		toString(
			[]shared.Loc{
				{X: 2, Y: 3},
				{X: 7, Y: 0},
			}),
		toString(expanded),
	)
}

func TestCountBelow(t *testing.T) {
	run := func(expected, value int, values []int) {
		name := fmt.Sprintf(
			"%d < %s",
			value,
			strings.Join(
				gent.Map(values, func(i int) string {
					return strconv.Itoa(i)
				}),
				","))
		t.Run(name, func(t *testing.T) {
			require.Equal(t, expected, countBelow(value, values))
		})
	}

	run(0, 0, []int{1, 3, 5})
	run(1, 2, []int{1, 3, 5})
	run(2, 4, []int{1, 3, 5})
	run(3, 6, []int{1, 3, 5})
}
