package aoc2322

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

func TestCountBricks(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	// EXERCISE & VERIFY
	req.Equal(5, CountBricksFromLines(lines))
}

func TestFindSlackers(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	bricks := parseLines(lines)
	byLowZ, byHighZ := createSearchIndexes(bricks)
	descend(bricks, byLowZ, byHighZ)
	// EXERCISE
	slackers := findSlackers(bricks, byLowZ, byHighZ).ToSlice()
	slices.SortFunc(slackers, func(a, b brick) int {
		z := a.start.z - b.start.z
		if z != 0 {
			return z
		}
		x := a.start.x - b.start.x
		if x != 0 {
			return x
		}
		return a.start.y - b.start.y
	})

	// VERIFY
	expected := []brick{
		// B
		parseBrick("0-2,0,2"),
		// C
		parseBrick("0-2,2,2"),
		// D
		parseBrick("0,0-2,3"),
		// E
		parseBrick("2,0-2,3"),
		// g
		parseBrick("1,1,5-6"),
	}
	// Much easier to read as strings but something could be left out of them.
	req.Equal(stringify(expected), stringify(slackers))
	req.Equal(expected, slackers)
}

func TestCountBricksSelectCases(t *testing.T) {
	run := func(name string, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			slackerCount := CountBricksFromLines(lines)

			// VERIFY
			req.Equal(expected, slackerCount)
		})
	}

	run(
		"2 * 1x1 on top of each other",
		[]string{
			"1,1,1~1,1,1",
			"1,1,2~1,1,2",
		},
		1)
	run(
		"2 * 1x1 side by side",
		[]string{
			"0,0,1~0,0,1",
			"0,1,1~0,1,1",
		},
		2)
	run(
		"2 * 1x1 side by side, one in air",
		[]string{
			"0,0,1~0,0,1",
			"0,1,2~0,1,2",
		},
		2)
	run(
		"2 * 2x1 overlapping, one in air",
		[]string{
			"0,0,1~1,0,1",
			"1,0,3~1,1,3",
		},
		1)
	run(
		"2 * 2x1 overlapping, one in air",
		[]string{
			"0,0,1~1,0,1",
			"1,0,3~1,1,3",
		},
		1)
	run(
		"supported by three pillars, all dropping",
		[]string{
			"2,0,10~2,0,20",
			"4,0,110~4,0,120",
			"6,0,210~6,0,220",
			"0,0,310~8,0,310",
		},
		4)
	run(
		"three on top of each other",
		[]string{
			"0,0,1~4,0,1",
			"1,0,2~1,4,2",
			"0,3,3~4,3,3",
		},
		1)
	// z
	// 4        a a
	// 3      b b c c
	// 2    d d e e f f
	// 1  g g h h i i j j
	//  0 1 2 3 4 5 6 7 8 x
	run(
		"pyramid - all but top share support with another",
		[]string{
			/* a */ "4,0,4~5,0,4",
			/* b */ "3,0,3~4,0,3",
			/* c */ "5,0,3~6,0,3",
			/* d */ "2,0,2~3,0,2",
			/* d */ "4,0,2~5,0,2",
			/* d */ "6,0,2~7,0,2",
			/* e */ "1,0,1~2,0,1",
			/* f */ "3,0,1~4,0,1",
			/* g */ "5,0,1~6,0,1",
			/* h */ "7,0,1~8,0,1",
		},
		10)
	// z
	// 2|bdd
	// 1|aac
	//  +---
	//   012x
	run(
		"brick supports two bricks, one only supported by it",
		[]string{
			"0,0,1~1,0,1", // a
			"0,0,2~0,0,2", // b
			"2,0,1~2,0,1", // c
			"1,0,2~2,0,2", // d
		},
		3)
}

func stringify[T any, S ~[]T](s S) []string {
	result := make([]string, len(s))
	for i, each := range s {
		result[i] = fmt.Sprint(each)
	}
	return result
}

func TestCreateSearchIndexes(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	b := func(start, end int) brick {
		return brick{
			start: coordinate{z: start},
			end:   coordinate{z: end},
		}
	}
	bricks := []brick{
		b(0, 0),
		b(1, 2),
		b(0, 2),
		b(3, 3),
	}
	// EXERCISE
	low, high := createSearchIndexes(bricks)

	// VERIFY
	expectedLow := map[int][]int{
		0: {0, 2},
		1: {1},
		3: {3},
	}
	req.Equal(expectedLow, low)
	expectedHigh := map[int][]int{
		0: {0},
		2: {1, 2},
		3: {3},
	}
	req.Equal(expectedHigh, high)
}

func TestFindOverlaps(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	//   +---+---+---+---+
	// 3 |   |E  |G  |   |
	//   +---/////////---+
	// 2 |   /E  /F  /   |
	//   +---/---/---/---+
	// 1 |   /A  /B  /BC  |
	//   +---/////////---+
	// 0 |D  |   |   |   |
	//   +---+---+---+---+
	//     0   1   2   3
	bricks := []brick{
		// A: contained
		parseBrick("1,1,0"),
		// B: overlapping
		parseBrick("2-3,1,0"),
		// C: outside
		parseBrick("3,1,0"),
		// D: outside and not included in indexes
		{},
		// E: overlapping
		parseBrick("1,2-3,0"),
		// F: inside but not included in indexes
		parseBrick("2,2,0"),
		// G: outside
		parseBrick("2,3,0"),
	}
	indexes := []int{0, 1, 2, 4, 6}
	base := parseBrick("1-2,1-2,0")
	// EXERCISE
	overlaps := findOverlaps(bricks, indexes, base)

	// VERIFY
	req.Equal([]int{0, 1, 4}, overlaps)
}

func parseInt(s string) int {
	return gent.OrPanic2(strconv.Atoi(s))("parse int")
}

func TestDescend(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	bricks := []brick{
		parseBrick("0-1,1,1"),   // 0
		parseBrick("1,1-2,3-4"), // 1
		parseBrick("1-2,1,6-8"), // 2
		parseBrick("1,0-1,10"),  // 3
		parseBrick("0,1,2-100"), // 4, over half of 0
		parseBrick("1,2,1"),     // 5, under 1
	}
	originalBricks := make([]brick, len(bricks))
	copy(originalBricks, bricks)
	byLow := map[int][]int{
		1:  {0, 5},
		2:  {4},
		3:  {1},
		6:  {2},
		10: {3},
	}
	byHigh := map[int][]int{
		1:   {0, 5},
		4:   {1},
		8:   {2},
		10:  {3},
		100: {4},
	}
	// EXERCISE
	descend(bricks, byLow, byHigh)

	changeZ := func(b brick, delta int) brick {
		b.start.z += delta
		b.end.z += delta
		return b
	}

	// VERIFY
	expected := []brick{
		// Not changed
		originalBricks[0],
		changeZ(originalBricks[1], -1),
		changeZ(originalBricks[2], -2),
		changeZ(originalBricks[3], -3),
		originalBricks[4],
		originalBricks[5],
	}
	req.Equal(stringify(expected), stringify(bricks), "stringified bricks")
	req.Equal(expected, bricks, "bricks")

	req.Equal(
		map[int][]int{
			1: {0, 5},
			2: {1, 4},
			4: {2},
			7: {3},
		},
		byLow,
		"by low")
	req.Equal(
		map[int][]int{
			1:   {0, 5},
			3:   {1},
			6:   {2},
			7:   {3},
			100: {4},
		},
		byHigh,
		"by high")
}

func claimExactly(message string, value, expected int) {
	if value == expected {
		return
	}
	panic(fmt.Sprintf("expected %d, got %d, message: %s", expected, value, message))
}

func parseBrick(s string) brick {
	mainPieces := strings.Split(s, ",")
	claimExactly("main pieces", len(mainPieces), 3)
	r := func(coord string) (int, int) {
		cPieces := strings.Split(coord, "-")
		first := parseInt(cPieces[0])
		if len(cPieces) == 1 {
			return first, first
		}
		return first, parseInt(cPieces[1])
	}
	result := brick{}
	from, to := r(mainPieces[0])
	result.start.x = from
	result.end.x = to
	from, to = r(mainPieces[1])
	result.start.y = from
	result.end.y = to
	from, to = r(mainPieces[2])
	result.start.z = from
	result.end.z = to
	return result
}

func TestIsOverlap(t *testing.T) {
	run := func(a, b brick, expected bool) {
		name := fmt.Sprintf("%v and %v", a, b)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			result := isOverlap(a, b)

			// VERIFY
			req.Equal(expected, result)
		})
	}
	run(parseBrick("0-1,3-4,1"), parseBrick("2-3,1-2,1"), false)
	run(parseBrick("0-1,3-4,1"), parseBrick("1-2,1-2,1"), false)
	run(parseBrick("0-1,3-4,1"), parseBrick("1-2,2-3,1"), true)
	run(parseBrick("0-1,3-4,1"), parseBrick("1-2,3-4,1"), true)
	run(parseBrick("0-1,3-4,1"), parseBrick("1-2,4-5,1"), true)
	run(parseBrick("0-1,3-4,1"), parseBrick("1-2,5-6,1"), false)

	run(parseBrick("1-2,3-4,1"), parseBrick("1-2,5-6,1"), false)
	run(parseBrick("1-2,4-5,1"), parseBrick("1-2,5-6,1"), true)
	run(parseBrick("1-2,5-6,1"), parseBrick("1-2,5-6,1"), true)
	run(parseBrick("2-3,5-6,1"), parseBrick("1-2,5-6,1"), true)
	run(parseBrick("3-4,5-6,1"), parseBrick("1-2,5-6,1"), false)
}
