package aoc2315

import (
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumHashes(t *testing.T) {
	run := func(f func([]string) int, expected int) {
		shared.InitTestLogging(t)
		req := require.New(t)

		lines, err := inr.ReadPath("testdata/in.txt")
		req.NoError(err)

		// EXERCISE & VERIFY
		req.Equal(expected, f(lines))
	}

	run(SumHashes, 1320)
	run(DeriveFocusingPower, 145)
}

func TestHash(t *testing.T) {
	require.Equal(t, 30, hash("rn=1"))
}

func TestParseLines(t *testing.T) {
	var pieces []string
	parseLines(
		[]string{"", ", ,", ",,ab,,cd,,"},
		func(s string) {
			pieces = append(pieces, s)
		})

	require.Equal(
		t,
		[]string{" ", "ab", "cd"},
		pieces)
}

func TestLensBoxAdd(t *testing.T) {
	box := &lensBox{}
	box.add("hg", 6)
	box.add("au", 2)
	req := require.New(t)
	req.Equal(
		[]lens{
			{
				label: "hg",
				focal: 6,
			},
			{
				label: "au",
				focal: 2,
			},
		},
		box.lenses)
}

func TestLensBoxRemove(t *testing.T) {
	run := func(expectedIndexes []int, labels ...string) {
		t.Run(strings.Join(labels, ","), func(t *testing.T) {
			req := require.New(t)
			box := &lensBox{}
			box.add("fe", 7)
			box.add("hg", 3)
			box.add("ag", 9)

			var expected []lens
			for _, i := range expectedIndexes {
				expected = append(expected, box.lenses[i])
			}

			// EXERCISE
			for _, each := range labels {
				box.remove(each)
			}

			// VERIFY
			req.Equal(expected, box.lenses)
		})
	}

	run([]int{1, 2}, "fe")
	run([]int{0, 2}, "hg")
	run([]int{0, 1}, "ag")
	run([]int{2}, "fe", "hg")
	run(nil, "fe", "hg", "ag")

	require.Equal(t, 2, len(strings.Split("a-", "-")))
}
