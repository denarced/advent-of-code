package aoc2317

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveLeastHeatLoss(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err)

	req.Equal(102, DeriveLeastHeatLoss(lines))
}

func BenchmarkDeriveLeastHeatLoss(b *testing.B) {
	shared.InitNullLogging()
	lines, _ := inr.ReadPath("testdata/in.txt")

	for range b.N {
		DeriveLeastHeatLoss(lines)
	}
}

type testCase struct {
	name     string
	lines    []string
	expected int
}

func TestDeriveLeastHeatLossCases(t *testing.T) {
	run := func(name string, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			least := DeriveLeastHeatLoss(lines)
			if expected > 0 {
				req.Equal(expected, least)
			} else {
				req.Less(least, 1_000)
			}
		})
	}

	for _, each := range readTestCases(require.New(t)) {
		run(each.name, each.lines, each.expected)
	}
}

func readTestCases(req *require.Assertions) (cases []testCase) {
	dirp := "testdata/gen"
	entries, err := os.ReadDir(dirp)
	req.NoError(err, "failed to read gen dir")
	for _, each := range entries {
		if each.IsDir() {
			continue
		}
		if !strings.HasSuffix(each.Name(), ".txt") {
			continue
		}
		filep := filepath.Join(dirp, each.Name())
		b, err := os.ReadFile(filep)
		req.NoErrorf(err, "failed to read file %s", filep)
		lines := strings.Split(strings.TrimSpace(string(b)), "\n")
		sum, err := strconv.Atoi(strings.Split(lines[0], ":")[1])
		req.NoErrorf(err, "failed to convert test case sum to int")
		cases = append(cases, testCase{
			name:     each.Name(),
			lines:    lines[1:],
			expected: sum,
		})
	}
	return
}

func TestBetween(t *testing.T) {
	run := func(from, to shared.Loc, expected []shared.Loc) {
		name := fmt.Sprintf("%v -> %v", from, to)
		t.Run(name, func(t *testing.T) {
			req := require.New(t)
			steps := between(from, to)
			req.Equal(expected, steps)
		})
	}

	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 1, Y: 3}, []shared.Loc{{X: 2, Y: 3}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 2, Y: 3}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 3}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 4, Y: 3}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 5, Y: 3}, []shared.Loc{{X: 4, Y: 3}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 5, Y: 3}, []shared.Loc{{X: 4, Y: 3}})

	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 1}, []shared.Loc{{X: 3, Y: 2}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 2}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 4}, nil)
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 5}, []shared.Loc{{X: 3, Y: 4}})
	run(shared.Loc{X: 3, Y: 3}, shared.Loc{X: 3, Y: 6}, []shared.Loc{{X: 3, Y: 4}, {X: 3, Y: 5}})
}

func TestDeriveNextHops(t *testing.T) {
	run := func(name string, latest hop, arrow shared.Loc, expected []hop) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			brd := shared.NewBoard([]string{
				"012345",
				"012345",
				"012345",
			})
			aRunner := &runner{latestHop: latest, arrow: arrow}
			// EXERCISE
			hops := deriveNextHops(brd, aRunner)

			// VERIFY
			req.Equal(expected, hops)
		})
	}

	run(
		"happy path X",
		hop{loc: shared.Loc{X: 4, Y: 1}, dir: shared.RealEast},
		shared.Loc{X: 2},
		[]hop{
			{
				loc: shared.Loc{X: 5, Y: 1},
				dir: shared.Direction{X: 1},
			},
			{
				loc: shared.Loc{X: 4, Y: 2},
				dir: shared.Direction{Y: 1},
			},
			{
				loc: shared.Loc{X: 4, Y: 0},
				dir: shared.Direction{Y: -1},
			},
		})
	run(
		"happy path Y",
		hop{loc: shared.Loc{X: 3, Y: 1}, dir: shared.RealNorth},
		shared.Loc{Y: 1},
		[]hop{
			{
				loc: shared.Loc{X: 3, Y: 2},
				dir: shared.RealNorth,
			},
			{
				loc: shared.Loc{X: 2, Y: 1},
				dir: shared.RealWest,
			},
			{
				loc: shared.Loc{X: 4, Y: 1},
				dir: shared.RealEast,
			},
		})
	run(
		"dead end X",
		hop{loc: shared.Loc{}, dir: shared.RealWest},
		shared.Loc{X: 1},
		[]hop{{loc: shared.Loc{Y: 1}, dir: shared.RealNorth}})
	run(
		"dead end Y",
		hop{loc: shared.Loc{X: 5, Y: 2}, dir: shared.RealEast},
		shared.Loc{Y: 2},
		[]hop{{loc: shared.Loc{X: 5, Y: 1}, dir: shared.RealSouth}})
	run(
		"blocked straight X",
		hop{loc: shared.Loc{X: 2}, dir: shared.RealEast},
		shared.Loc{X: 3},
		[]hop{{loc: shared.Loc{X: 2, Y: 1}, dir: shared.RealNorth}})
	run(
		"blocked straight Y",
		hop{loc: shared.Loc{X: 1, Y: 1}, dir: shared.RealSouth},
		shared.Loc{Y: 3},
		[]hop{
			{loc: shared.Loc{X: 2, Y: 1}, dir: shared.RealEast},
			{loc: shared.Loc{X: 0, Y: 1}, dir: shared.RealWest},
		})
}
