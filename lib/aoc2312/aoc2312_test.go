package aoc2312

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestSumPermutations(t *testing.T) {
	run := func(expected, mul int) {
		t.Run(fmt.Sprint(mul), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err)

			// EXERCISE & VERIFY
			req.Equal(expected, SumPermutations(lines, mul))
		})
	}

	run(21, 1)
	run(525152, 5)
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
	run := func(line string, expected, mul int) {
		t.Run(line, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			row := parseLine(line)
			// EXERCISE
			count := countPermutations(row, mul)

			// VERIFY
			req.Equal(expected, count)
		})
	}
	run("???.### 1,1,3", 1, 1)
	run("???.### 1,1,3", 1, 5)
	run(".??..??...?##. 1,1,3", 4, 1)

	// The first has 4 permutations. Unknown between it and the second becomes a new wildcard in the
	// second and thus a candidate for the first 1, thus increasing permutations in the second to 8.
	// The same logic applies to all parts except the first.
	expected := 4 * pow(8, 4)
	run(".??..??...?##. 1,1,3", expected, 5)

	run("?###???????? 3,2,1", 10, 1)
	run("????.#...#... 4,1,1", 16, 5)
	run("????.######..#####. 1,6,5", 2_500, 5)
	run("?###???????? 3,2,1", 506_250, 5)
	run("#???#?????.?#?. 2,1,2,1", 87_489, 5)
}

func pow(base, exp int) int {
	if exp == 0 {
		return 1
	}
	result := base
	for range exp - 1 {
		result *= base
	}
	return result
}

func BenchmarkCountPermutations(b *testing.B) {
	shared.InitNullLogging()
	for range b.N {
		countPermutations(parseLine("???.### 1,1,3"), 1)
		countPermutations(parseLine(".??..??...?##. 1,1,3"), 1)
		countPermutations(parseLine("?###???????? 3,2,1"), 1)
	}
}

func BenchmarkCountPermutationsWithMultiplier(b *testing.B) {
	shared.InitNullLogging()
	mul := 5
	for range b.N {
		countPermutations(parseLine("?#?.??##?????#.???? 1,1,4,1,1,3"), mul)
		countPermutations(parseLine("???.### 1,1,3"), mul)
		countPermutations(parseLine(".??..??...?##. 1,1,3"), mul)
		countPermutations(parseLine("?###???????? 3,2,1"), mul)
	}
}

func BenchmarkCaching(b *testing.B) {
	shared.InitNullLogging()
	for range b.N {
		countPermutations(parseLine("?..?#?????.. 2,1"), 5)
	}
}

func TestHypothesize(t *testing.T) {
	run := func(s string, expected int) {
		t.Run(s, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			parsed := parseLine(s)
			// EXERCISE
			count := hypothesize(parsed.springs, parsed.groups)

			// VERIFY
			req.Equal(expected, count)
		})
	}
	run("?##??.?#.????? 5,1,4", 2)
	run("???.### 1,1,3", 1)
	run("?????#??????????# 6,1,6", 6)
}

func TestMultiplySpring(t *testing.T) {
	req := require.New(t)
	req.Equal(
		parseSpring("???.###????.###"),
		multiplySpring(parseSpring("???.###"), 2))
	req.Equal(
		parseSpring("???.###"),
		multiplySpring(parseSpring("???.###"), 1))
}

func TestMultiplyGroups(t *testing.T) {
	req := require.New(t)
	req.Equal([]int{1, 2, 3, 1, 2, 3}, multiplyGroups([]int{1, 2, 3}, 2))
	req.Equal([]int{1, 2, 3}, multiplyGroups([]int{1, 2, 3}, 1))
}

func TestCreateCondCounter(t *testing.T) {
	req := require.New(t)
	counter := createCondCounter(parseSpring(".##?#.??."), []int{2, 1, 1})
	req.Equal(
		condCounter{
			target: condPair{damaged: 4, operational: 5},
			status: condPair{damaged: 3, operational: 3},
		},
		*counter)
}
