package aoc2402

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestCountSafe(t *testing.T) {
	run := func(name string, levels [][]int, dampener bool, expected int) {
		prefix := "wo dampener"
		if dampener {
			prefix = "with dampener"
		}
		t.Run(fmt.Sprintf("%s - %s", prefix, name), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := CountSafe(levels, dampener)
			req.Equal(expected, actual)
		})
	}

	// Without dampener.
	run("empty", [][]int{}, false, 0)
	run("too short", [][]int{{1}}, false, 0)
	run("minimal asc safe", [][]int{{1, 2}}, false, 1)
	run("minimal desc safe", [][]int{{2, 1}}, false, 1)
	run("asc to desc", [][]int{{1, 2, 1}}, false, 0)
	run("too little diff", [][]int{{1, 2, 2}}, false, 0)
	run("too big diff", [][]int{{1, 2, 3, 7}}, false, 0)
	run("2 safe, 1 unsafe", [][]int{{1, 4, 6, 7}, {9, 5, 4, 3, 2}, {100, 97, 95, 94}}, false, 2)

	// With dampener.
	run("empty", [][]int{}, true, 0)
	run("too short", [][]int{{1}}, true, 0)
	run("minimal asc safe", [][]int{{1, 2}}, true, 1)
	run("minimal desc safe", [][]int{{2, 1}}, true, 1)
	run("asc to desc", [][]int{{1, 2, 1}}, true, 1)
	run("too little diff", [][]int{{1, 2, 2}}, true, 1)
	run("too big diff", [][]int{{1, 2, 3, 7}}, true, 1)
	run("2 safe, 1 unsafe", [][]int{{1, 4, 6, 7}, {9, 5, 4, 3, 2}, {100, 97, 95, 94}}, true, 3)
	run("long asc to desc", [][]int{{4, 7, 4, 1}}, true, 1)
	run("fix before unsafe pair", [][]int{{4, 7, 4, 1}}, true, 1)
	run("drop too big diff", [][]int{{5, 9, 8}}, true, 1)
}
