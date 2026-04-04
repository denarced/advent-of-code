package aoc2312

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumPermutations(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err)

	// EXERCISE & VERIFY
	req.Equal(21, SumPermutations(lines))
}

func TestParseLine(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	// EXERCISE
	row := parseLine("???.### 1,1,3")

	// VERIFY
	req.Equal(
		springRow{
			springs: []condition{
				condUnknown,
				condUnknown,
				condUnknown,
				condOperational,
				condDamaged,
				condDamaged,
				condDamaged,
			},
			groups: []int{1, 1, 3},
		},
		row)
}

func TestCountPermutations(t *testing.T) {
	run := func(line string, expected int) {
		t.Run(line, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			row := parseLine(line)
			// EXERCISE
			count := countPermutations(row)

			// VERIFY
			req.Equal(expected, count)
		})
	}
	run("???.### 1,1,3", 1)
	run(".??..??...?##. 1,1,3", 4)
	run("?###???????? 3,2,1", 10)
}

func BenchmarkCountPermutations(b *testing.B) {
	shared.InitNullLogging()
	for range b.N {
		countPermutations(parseLine("???.### 1,1,3"))
		countPermutations(parseLine(".??..??...?##. 1,1,3"))
		countPermutations(parseLine("?###???????? 3,2,1"))
	}
}

func TestHypothesize(t *testing.T) {
	run := func(s string, expected []string) {
		t.Run(s, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			var variants []spring
			cb := func(s spring) {
				variants = append(variants, s)
			}
			parsed := parseLine(s)
			// EXERCISE
			hypothesize(parsed.springs, parsed.groups, cb)

			// VERIFY
			req.Equal(stringify(expected), stringify(variants))
		})
	}
	run("?##??.?#.????? 5,1,4",
		[]string{
			"#####..#..####",
			"#####..#.####.",
		})
	run("???.### 1,1,3", []string{"#.#.###"})
	run("?????#??????????# 6,1,6", []string{
		"..######.#.######",
		".######..#.######",
		".######.#..######",
		"######...#.######",
		"######..#..######",
		"######.#...######",
	})
}

func stringify[S ~[]T, T any](s S) (result []string) {
	if s == nil {
		return
	}
	result = make([]string, len(s))
	for i, each := range s {
		result[i] = fmt.Sprint(each)
	}
	return
}
