package aoc2204

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountOverlaps(t *testing.T) {
	run := func(total bool, expected int) {
		name := "partially"
		if total {
			name = "total"
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err)

			req.Equal(expected, CountOverlaps(lines, total))
		})
	}
	run(true, 2)
	run(false, 4)
}

func TestOverlaps(t *testing.T) {
	run := func(v [4]int, expected bool) {
		name := fmt.Sprint(v)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			req.Equal(expected, overlaps(v[0], v[1], v[2], v[3]), "normal")
			req.Equal(expected, overlaps(v[2], v[3], v[0], v[1]), "reverse")
		})
	}

	run([4]int{1, 2, 3, 4}, false)
	run([4]int{1, 4, 2, 3}, true)
	run([4]int{1, 2, 2, 3}, true)
	run([4]int{1, 2, 1, 2}, true)
	run([4]int{1, 3, 2, 4}, true)
}
