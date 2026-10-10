package aoc2217

import (
	"fmt"
	"slices"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestMeasureTowerHeight(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	req.Equal(3068, MeasureTowerHeight(lines, 2022))
}

func TestBuryHidden(t *testing.T) {
	run := func(name string, pieces, expected [][]shared.Loc) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			board := new(pieceBoard)
			board.xToWhys = convertLocTableToExToWhys(pieces)
			board.wanderBook = createGuestbook()
			board.visibleBook = createGuestbook()
			// EXERCISE
			board.buryHidden()

			// VERIFY
			req.Equal(convertLocTableToExToWhys(expected), board.xToWhys)
		})
	}

	floor := createPieces([]string{"aaaaaaa"})
	run("just floor", floor, floor)

	run(
		"pyramid",
		createPieces(
			[]string{
				"   d   ",
				"  ccc  ",
				" bbbbb ",
				"aaaaaaa",
			},
		),
		createPieces(
			[]string{
				"   d   ",
				"  c c  ",
				" b   b ",
				"a     a",
			},
		))
	run(
		"hidden corner",
		createPieces([]string{
			"cc     ",
			"bc     ",
			"aaaaaaa",
		}),
		createPieces([]string{
			"cc     ",
			" c     ",
			"  aaaaa",
		}))
	run(
		"complex",
		createPieces([]string{
			"  c    ",
			"  c    ",
			"cccb   ",
			"  bbb  ",
			"   b   ",
			"  aaaa ",
			"zzzzzzz",
		}),
		createPieces([]string{
			"  c    ",
			"  c    ",
			"cc b   ",
			"    b  ",
			"   b   ",
			"    aa ",
			"      z",
		}))
}

func createPieces(lines []string) [][]shared.Loc {
	idToLocs := map[rune][]shared.Loc{}
	for i := range lines {
		for j, r := range lines[i] {
			if r == ' ' {
				continue
			}
			trueIndex := reverseLineIndex(len(lines), i)
			idToLocs[r] = append(idToLocs[r], shared.Loc{X: j, Y: trueIndex})
		}
	}
	sortedIDs := make([]rune, 0, len(idToLocs))
	for r := range idToLocs {
		sortedIDs = append(sortedIDs, r)
	}
	slices.Sort(sortedIDs)

	pieces := make([][]shared.Loc, 0, len(idToLocs))
	for _, id := range sortedIDs {
		pieces = append(pieces, idToLocs[id])
	}
	return pieces
}

func reverseLineIndex(height int, index int) int {
	return shared.Abs(height-index) - 2
}

func TestReverseLineIndex(t *testing.T) {
	run := func(index, expected int) {
		t.Run(fmt.Sprintf("%d-%d", index, expected), func(t *testing.T) {
			require.Equal(t, expected, reverseLineIndex(9, index))
		})
	}

	run(0, 7)
	run(1, 6)
	run(6, 1)
	run(8, -1)
}

func BenchmarkMeasureTowerHeight(b *testing.B) {
	shared.InitNullLogging()
	req := require.New(b)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")

	for range b.N {
		MeasureTowerHeight(lines, 10_000)
	}
}

func convertLocTableToExToWhys(table [][]shared.Loc) []map[int]struct{} {
	xToWhys := make([]map[int]struct{}, boardWidth)
	for _, arr := range table {
		for _, each := range arr {
			exes := xToWhys[each.X]
			if exes == nil {
				exes = map[int]struct{}{}
				xToWhys[each.X] = exes
			}
			exes[each.Y] = struct{}{}
		}
	}
	return xToWhys
}
