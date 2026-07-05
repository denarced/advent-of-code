package aoc2321

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestCountRangeFromLines(t *testing.T) {
	run := func(stepCount int, infinite bool, expected int) {
		kind := "restricted"
		if infinite {
			kind = "infinite"
		}
		name := fmt.Sprintf("%s-%d", kind, stepCount)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "failed to read test data")

			// EXERCISE & VERIFY
			req.Equal(expected, CountRangeFromLines(lines, stepCount, infinite))
		})
	}

	run(1, false, 2)
	run(1, true, 2)
	run(2, false, 4)
	run(2, true, 4)
	run(6, false, 16)
	run(6, true, 16)
	run(50, false, 42)
	run(10, true, 50)
	run(50, true, 1594)
	run(100, true, 6536)
}

func TestDeriveSection(t *testing.T) {
	deriveSection := createDeriveSection(3)
	run := func(loc, expected shared.Loc) {
		t.Run(loc.ToString(), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			result := deriveSection(loc)

			// VERIFY
			req.Equal(expected, result)
		})
	}

	l := func(s string) shared.Loc {
		return shared.ParseLoc(s)
	}
	run(l("0x0"), l("0x0"))
	run(l("-2x-2"), l("-1x-1"))
	run(l("3x3"), l("1x1"))
	run(l("-3x-3"), l("-1x-1"))
	run(l("-4x-4"), l("-2x-2"))
	run(l("5x5"), l("1x1"))
}

func TestRepeatMonitor(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	// EXERCISE
	monitor := createRepeatMonitor()
	// Fill up.
	for i := range 10 {
		sizes := monitor(i + 1)
		req.Nilf(sizes, "expected nil on index %d", i)
	}
	// Start zig-zag.
	for i := range 4 {
		sizes := monitor(41)
		req.Nilf(sizes, "expected nil during initial zig-zag, first call, index: %d", i)
		sizes = monitor(42)
		req.Nilf(sizes, "expected nil during initial zig-zag, second call, index: %d", i)
	}
	req.Nil(monitor(41), "nil still expected on second to last call")
	req.Equal([]int{41, 42}, monitor(42), "10th zig-zag call should result in returned pattern")
}

func TestDiamond(t *testing.T) {
	run := func(stepCount, expectedRadius int) {
		t.Run(strconv.Itoa(stepCount), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			size := 5
			// EXERCISE
			diam := newDiamond(shared.Loc{X: size / 2, Y: size / 2}, stepCount, size)

			// VERIFY
			req.Equal(expectedRadius, diam.radius)
		})
	}

	run(1, 0)
	run(2, 0)
	run(5+2, 0)
	run(2*5+2, 0)
	run(2*5+3, 1)
	run(4*5+3, 3)
	run(4*5+3, 3)
}

func TestDiamondCountTotal(t *testing.T) {
	run := func(stepCount, expectedTotal int) {
		t.Run(strconv.Itoa(stepCount), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			size := 5
			diam := newDiamond(shared.Loc{X: size / 2, Y: size / 2}, stepCount, size)
			// EXERCISE
			total := diam.countTotal(11, 13)

			// VERIFY
			req.Equal(expectedTotal, total)
		})
	}

	run(2*5+2, 0)
	run(2*5+3, 1*11+4*13)
	run(3*5+3, 9*11+4*13)
	run(4*5+3, 9*11+16*13)
}

func TestDiamondForbiddenFilter(t *testing.T) {
	// Step from S to x.
	// .....|.....|.....<.....|.....|.....
	// .....|.....|.....|.....|.....|.....
	// .....|.....|.....|.....|.....|.....
	// .....|.....|.....|.....|.....|.....
	// .....|.....|.....|.....|.....|.....
	// -----+-----+-----+-----+-----+-----
	// .....|.....|.....|.....<.....|.....
	// .....|.....|.....|.....|.....|.....
	// ..S..|.....|.....|.....|.....|x....
	// .....|.....|.....|.....|.....|.....
	// .....|.....|.....|.....|.....|.....
	run := func(loc shared.Loc, expected bool) {
		t.Run(loc.ToString(), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			size := 5
			stepCount := 4*size + size/2 + 1
			diam := newDiamond(shared.Loc{X: size / 2, Y: size / 2}, stepCount, size)
			// EXERCISE
			forbidden := diam.isForbidden(loc)

			// VERIFY
			req.Equal(expected, forbidden)
		})
	}

	l := func(x, y int) shared.Loc {
		return shared.Loc{X: x, Y: y}
	}

	run(l(0, 0), true)

	// Left edge of Y:0.
	run(l(0-3*5, 4), true)
	run(l(0-3*5-1, 4), false)
	run(l(0-3*5-1, 5), false)
	run(l(0-3*5, 5), false)

	// Left edge of Y:1.
	run(l(0-2*5, 5), true)
	run(l(0-2*5-1, 5), false)

	// Left edge of Y:2.
	run(l(0-1*5, 10), true)
	run(l(0-1*5-1, 10), false)

	// Left edge of Y:3.
	run(l(0-0*5, 15), true)
	run(l(0-0*5-1, 15), false)

	// Edge of top.
	run(l(0, 19), true)
	run(l(-1, 19), false)
	run(l(0, 20), false)

	// Right edge of Y:0.
	run(l(4+3*5, 0), true)
	run(l(4+3*5+1, 0), false)

	// Right edge of Y:-1.
	run(l(4+2*5, -1), true)
	run(l(4+2*5+1, -1), false)

	// Right edge of Y:-2.
	run(l(4+1*5, -10), true)
	run(l(4+1*5+1, -10), false)

	// Right edge of Y:-3.
	run(l(4+0*5, -11), true)
	run(l(4+0*5+1, -11), false)

	// Edge of bottom.
	run(l(4, -3*5), true)
	run(l(4, -3*5-1), false)
	run(l(5, -3*5), false)
}

func TestDiamondStepOutside(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	size := 5
	diam := newDiamond(shared.Loc{X: size / 2, Y: size / 2}, 3+3*size, size)
	// EXERCISE
	starts, steps := diam.stepOutside()

	// VERIFY
	req.Equal(
		[]shared.Loc{
			// Top
			{X: 2, Y: 15},
			// Descend on the right side of triangle
			{X: 5, Y: 12},
			{X: 10, Y: 7},
			{X: 15, Y: 2},
			{X: 10, Y: -3},
			{X: 5, Y: -8},
			// Botttom
			{X: 2, Y: -11},
			// Ascend on the left side of triangle
			{X: -1, Y: -8},
			{X: -6, Y: -3},
			{X: -11, Y: 2},
			{X: -6, Y: 7},
			{X: -1, Y: 12},
		},
		starts)
	req.Equal(13, steps)
}

func TestCountInfiniteRangeWithKnownFailing(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines := []string{
		"...........",
		"....#....#.",
		"..#.#......",
		"......#....",
		"...#.......",
		".....S.....",
		"..#.....#..",
		"...........",
		".#.......#.",
		"..#......#.",
		"...........",
	}
	stepCount := 151
	// EXERCISE
	efficientCount := CountInfiniteRange(lines, stepCount)
	naiveCount := CountRangeFromLines(lines, stepCount, true)
	brd := shared.NewBoard(lines)
	center := brd.FindOrDie('S')
	brd.Set(center, '.')
	discoveredCount := discoverPlots(
		brd,
		[]shared.Loc{center},
		stepCount,
		func(shared.Loc) bool { return false },
	).Count()

	// VERIFY
	expected := 20808
	req.Equal(expected, naiveCount, "naive")
	req.Equal(expected, discoveredCount, "discovered")
	req.Equal(expected, efficientCount, "efficient")
}

func TestQuarterFeed(t *testing.T) {
	run := func(prefix string, size int, dtos ...[][2]int) {
		t.Run(prefix+"-"+strconv.Itoa(size), func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			verifyPair := func(f func() (int, int), start, end int, message string) {
				a, b := f()
				req.Equalf(start, a, "start: %s", message)
				req.Equalf(end, b, "end: %s", message)
			}

			feed := newQuarterFeed(size)
			// EXERCISE & VERIFY
			for i, each := range dtos {
				f := feed.create()
				if each == nil {
					req.Nilf(f, "nil expected %d", i)
					continue
				}
				for j, dto := range each {
					verifyPair(f, dto[0], dto[1], fmt.Sprintf("pair %d,%d", i, j))
				}
			}
		})
	}

	convertInt := func(s string) int {
		return gent.OrPanic2(strconv.Atoi(s))("failed to convert int " + s)
	}
	createSingleDtos := func(spec string) (result [][][2]int) {
		outer := strings.Split(spec, ";")
		for _, each := range outer {
			inner := strings.Split(each, ",")
			var nested [][2]int
			for _, anIn := range inner {
				value := convertInt(anIn)
				nested = append(nested, [2]int{value, -1})
			}
			result = append(result, append(nested, [2]int{-1, -1}))
		}
		result = append(result, nil)
		return
	}
	createProperDtos := func(spec string) (result [][][2]int) {
		outer := strings.Split(spec, "|")
		for _, each := range outer {
			inner := strings.Split(each, ";")
			var nested [][2]int
			for _, anIn := range inner {
				pieces := strings.Split(anIn, ",")
				pair := [2]int{convertInt(pieces[0]), convertInt(pieces[1])}
				nested = append(nested, pair)
			}
			result = append(result, append(nested, [2]int{-1, -1}))
		}
		result = append(result, nil)
		return
	}

	run("single", 5, createSingleDtos("0,1;2,3;4")...)
	run("single", 8, createSingleDtos("0,1;2,3;4,5;6,7")...)
	run("single", 11, createSingleDtos("0,1,2;3,4,5;6,7,8;9,10")...)
	run("proper", 48, createProperDtos(strings.Join([]string{
		"0,11;1,10;2,9;3,8;4,7;5,6",
		"12,23;13,22;14,21;15,20;16,19;17,18",
		"24,35;25,34;26,33;27,32;28,31;29,30",
		"36,47;37,46;38,45;39,44;40,43;41,42",
	}, "|"))...)
	run("proper", 49, createProperDtos(strings.Join([]string{
		"0,12;1,11;2,10;3,9;4,8;5,7;6,-1",
		"13,25;14,24;15,23;16,22;17,21;18,20;19,-1",
		"26,38;27,37;28,36;29,35;30,34;31,33;32,-1",
		"39,48;40,47;41,46;42,45;43,44",
	}, "|"))...)
	run("proper", 95, createProperDtos(strings.Join([]string{
		"0,23;1,22;2,21;3,20;4,19;5,18;6,17;7,16;8,15;9,14;10,13;11,12",
		"24,47;25,46;26,45;27,44;28,43;29,42;30,41;31,40;32,39;33,38;34,37;35,36",
		"48,71;49,70;50,69;51,68;52,67;53,66;54,65;55,64;56,63;57,62;58,61;59,60",
		"72,94;73,93;74,92;75,91;76,90;77,89;78,88;79,87;80,86;81,85;82,84;83,-1",
	}, "|"))...)
}
