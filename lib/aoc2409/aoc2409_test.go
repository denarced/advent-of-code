package aoc2409

import (
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestCountChecksum(t *testing.T) {
	run := func(name string, s string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			req.Equal(expected, CountChecksum(s))
		})
	}

	run("empty", "", 0)
	run("example", "2333133121414131402", 1928)
	// 122333
	// 0..111
	// 01.11
	// 0111
	run("123", "123", sum("0111"))
	// 001...2
	// 0012
	run("20131", "20131", sum("0012"))
	// 0..111....22222
	// 02.111....2222
	// 022111....222
	// 0221112...22
	// 02211122..2
	// 022111222
	run("12345", "12345", sum("022111222"))
}

func sum(s string) int {
	total := 0
	ints, err := shared.ToInts(strings.Split(s, ""))
	shared.Die(err, "sum -> ToInts")
	for i, each := range ints {
		total += i * each
	}
	return total
}

func TestCountDefragmentedChecksum(t *testing.T) {
	run := func(s string, expected int) {
		t.Run(s, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := CountDefragmentedChecksum(s)
			req.Equal(expected, actual)
		})
	}

	run("", sum(""))
	// 0..111
	run("123", sum("000111"))
	// 000..1
	// 0001
	run("321", sum("0001"))
	// 0........1112222
	// 02222....111
	// 02222111
	run("18304", sum("02222111"))
	run("2333133121414131402", 2858)
}
