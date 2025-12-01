package aoc2501

import (
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestSolvePassword(t *testing.T) {
	t.Run("with test data", func(t *testing.T) {
		lines, err := gent.ReadLines("testdata/in.txt")
		require.NoError(t, err)
		lines = gent.Map(lines, func(s string) string {
			return strings.TrimSpace(s)
		})
		lines = gent.Filter(lines, func(s string) bool {
			return s != ""
		})

		t.Run("final zeroes", func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, 3, SolvePassword(lines, false))
		})
		t.Run("all zeros", func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, 6, SolvePassword(lines, true))
		})
	})

	var tests = []struct {
		name     string
		expected int
	}{
		{"R49", 0}, {"R149", 1}, {"R249", 2},
		{"R50", 1}, {"R150", 2}, {"R250", 3},
		{"R51", 1}, {"R151", 2}, {"R251", 3},

		{"L49", 0}, {"L149", 1}, {"L249", 2},
		{"L50", 1}, {"L150", 2}, {"L250", 3},
		{"L51", 1}, {"L151", 2}, {"L251", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, tt.expected, SolvePassword([]string{tt.name}, true))
		})
	}
}

func TestSplitValues(t *testing.T) {
	var tests = []struct {
		name string
		in   []string
		out  []rotation
	}{
		{"empty", []string{}, nil},
		{"L0", []string{"L0"}, []rotation{{m: -1, values: []int{0}}}},
		{"R1", []string{"R1"}, []rotation{{m: 1, values: []int{1}}}},
		{"L99", []string{"L99"}, []rotation{{m: -1, values: []int{99}}}},
		{"R100", []string{"R100"}, []rotation{{m: 1, values: []int{100}}}},
		{"L101", []string{"L101"}, []rotation{{m: -1, values: []int{100, 1}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.InitTestLogging(t)
			actual := splitValues(tt.in)
			require.Equal(t, tt.out, actual)
		})
	}
}
