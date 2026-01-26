package aoc2510

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const maxVersion = 9

func TestDeriveFewestClicks(t *testing.T) {
	readLines := func(req *require.Assertions) []string {
		lines, err := inr.ReadPath("testdata/in.txt")
		req.NoError(err, "failed to read test data")
		return lines
	}
	t.Run("indicator", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		req.Equal(7, DeriveFewestClicks(readLines(req), true, 0))
	})
	for i := 1; i <= maxVersion; i++ {
		t.Run(fmt.Sprintf("joltage-%d", i), func(t *testing.T) {
			// 4: concurrent map writes for globalButtonIndexes.
			if i == 4 || i == 5 || i == 8 {
				return
			}
			shared.InitTestLogging(t)
			req := require.New(t)
			req.Equal(33, DeriveFewestClicks(readLines(req), false, i))
		})
	}
}

func BenchmarkDeriveFewestClicks(b *testing.B) {
	shared.InitNullLogging()
	req := require.New(b)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	for i := 1; i <= maxVersion; i++ {
		b.Run(fmt.Sprint(i), func(b *testing.B) {
			for range b.N {
				DeriveFewestClicks(lines, false, i)
			}
		})
	}
}

func TestDeriveFewestJoltageClicks(t *testing.T) {
	shared.InitTestLogging(t)
	line := strings.Join(
		[]string{
			"[..#.##] ",
			"(0,1,3,4,5) (3) (0,1,3,5) (3,5) (1,5) (0,2,3,5) (0,1,2,3) (0,2,4) ",
			"{25,12,13,57,14,38}",
		},
		"")
	mach := parseMachine(line)
	require.Equal(t, 59, deriveFewestJoltageClicks3(mach))
}

func TestParseInts(t *testing.T) {
	s := func(i ...int) []int {
		var ints []int
		return append(ints, i...)
	}
	var tests = []struct {
		in  string
		out []int
	}{
		{"", nil},
		{"0", s(0)},
		{"10", s(10)},
		{"0100 10 00", s(100, 10, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			require.Equal(t, tt.out, parseInts(tt.in))
		})
	}
}

func TestParseState(t *testing.T) {
	parsed := parseState("[.##..#.]")
	require.Equal(t, []bool{false, true, true, false, false, true, false}, parsed)
}

func TestToNumericState(t *testing.T) {
	require.Equal(t, int64(5), toNumericState([]bool{true, false, true}))
}

func TestDeriveClickLimits(t *testing.T) {
	mach := Machine{
		Buttons: [][]int{
			{1},
			{1, 2, 3},
		},
		Joltages: []int{20, 30, 60},
	}
	minClicks, maxClicks := deriveClickLimits(mach)
	ass := assert.New(t)
	ass.Equal(7, minClicks, "min")
	ass.Equal(60, maxClicks, "max")
}

func TestIter(t *testing.T) {
	container := make([]int, 2)
	result := [][]int{}
	iter(container, 2, 2, 0, 999, func() bool {
		added := make([]int, len(container))
		copy(added, container)
		result = append(result, added)
		return false
	})
	require.Equal(
		t,
		[][]int{
			{2, 2}, {2, 1}, {2, 0},
			{1, 2}, {1, 1}, {1, 0},
			{0, 2}, {0, 1}, {0, 0},
		},
		result)
}

func TestIt(t *testing.T) {
	shared.InitTestLogging(t)
	mach := Machine{
		Buttons: [][]int{
			{0},
			{1},
		},
		Joltages: []int{6, 8},
	}
	require.Equal(t, 14, deriveFewestJoltageClicks5(mach))
}

func TestDeriveButtonCombinations(t *testing.T) {
	for _, each := range []struct {
		name           string
		mach           Machine
		expectedCombos [][]int
	}{
		{
			name: "simplest",
			mach: Machine{
				Buttons:  [][]int{{0}, {1}},
				Joltages: []int{2, 3},
			},
			expectedCombos: [][]int{{0, 1}},
		},
		{
			name: "2-1",
			mach: Machine{
				Buttons:  [][]int{{0, 1}, {2}},
				Joltages: []int{2, 2, 4},
			},
			expectedCombos: [][]int{{0, 1}},
		},
		{
			name: "one overlap",
			mach: Machine{
				Buttons:  [][]int{{0, 1}, {1, 2}},
				Joltages: []int{4, 5, 1},
			},
			expectedCombos: [][]int{{0, 1}},
		},
		{
			name: "3+2+1",
			mach: Machine{
				Buttons:  [][]int{{0, 1, 2}, {1, 2}, {0}},
				Joltages: []int{10, 10, 10},
			},
			expectedCombos: [][]int{
				{0},
				{0, 1},
				{0, 1, 2},
				{0, 2},
				{1, 2},
			},
		},
		// Verify that once the first button has been clicked even once, the second is guaranteed to
		// be useless because joltage becomes all zeroes.
		{
			name: "useless second",
			mach: Machine{
				Buttons:  [][]int{{0, 1}, {1}},
				Joltages: []int{1, 1},
			},
			expectedCombos: [][]int{
				{0},
			},
		},
	} {
		t.Run(each.name, func(t *testing.T) {
			shared.InitTestLogging(t)
			var results [][]int
			initDeriveButtonCombinations(each.mach, func(combination []int) {
				results = append(results, combination)
			})
			require.Equal(t, each.expectedCombos, results)
		})
	}
}

func TestDeriveFewestJoltageClicks6(t *testing.T) {
	shared.InitTestLogging(t)
	buttons := [][]int{
		{1, 2, 3, 4},
		{1, 3},
		{0, 2, 3, 4},
		{1, 2},
	}
	expectedTotal := 39
	joltages := []int{8, 31, 21, 26, 8}
	mach := Machine{
		Buttons:  buttons,
		Joltages: joltages,
	}
	require.Equal(t, expectedTotal, deriveFewestJoltageClicks6(mach))
}

func TestSetupOverlord(t *testing.T) {
	mach := Machine{
		Buttons:  [][]int{{0, 1}, {2}},
		Joltages: []int{2, 1, 9},
	}
	ctx := context.TODO()
	// EXERCISE
	lord, angels := setupOverlord(ctx, mach, -1)
	ass := assert.New(t)
	ass.Equal(
		overlord{
			buttons:  mach.Buttons,
			joltages: mach.Joltages,
			minCount: math.MaxInt,
		},
		*lord)
	// Impossible to verify these so nullify before "ass.Equal".
	for _, each := range angels {
		each.reduceJoltage = nil
		each.increaseJoltage = nil
	}
	ass.Equal(
		[]*joltaAngel{
			{
				ctx:              ctx,
				joltageID:        0,
				value:            1,
				buttons:          []int{0},
				buttonJoltageIDs: [][]int{{0, 1}},
				lord:             lord,
				maxMultiplier:    -1,
			},
			{
				ctx:              ctx,
				joltageID:        1,
				value:            1,
				buttons:          []int{0},
				buttonJoltageIDs: [][]int{{0, 1}},
				lord:             lord,
				maxMultiplier:    -1,
			},
			{
				ctx:              ctx,
				joltageID:        2,
				value:            9,
				buttons:          []int{1},
				buttonJoltageIDs: [][]int{{2}},
				lord:             lord,
				maxMultiplier:    -1,
			},
		},
		angels)
}

func TestGetNextSeed(t *testing.T) {
	angel := newJoltaAngel(context.TODO(), 0, 3, []int{4, 6}, nil, nil, nil, nil, -1)
	var steps [][]step
	for {
		// EXERCISE
		each := angel.getNextSeed()
		if each == nil {
			break
		}
		steps = append(steps, each)
	}

	// VERIFY
	req := require.New(t)
	req.Nil(angel.ch, "channel should be nil")
	req.True(angel.done, "angel should be done")
	makeStep := func(id, count int) step {
		return step{
			buttonID: id,
			count:    count,
		}
	}
	req.Equal(
		[][]step{
			{
				makeStep(4, 3),
			},
			{
				makeStep(4, 2),
				makeStep(6, 1),
			},
			{
				makeStep(4, 1),
				makeStep(6, 2),
			},
			{
				makeStep(6, 3),
			},
		},
		steps)
}

func Test7(t *testing.T) {
	shared.InitTestLogging(t)
	var cases = []struct {
		buttons  [][]int
		joltages []int
		expected int
	}{
		{[][]int{{0, 1}}, []int{3, 3}, 3},
		{[][]int{{0, 1}, {1}}, []int{3, 3}, 3},
		{[][]int{{0}, {1}}, []int{3, 3}, 6},
		{[][]int{{0, 1}, {1, 2}}, []int{2, 4, 2}, 4},
		{[][]int{{3}, {1, 2, 3}, {0, 1, 2}}, []int{2, 2, 2, 1}, 3},
		{[][]int{{0}, {1}, {2}, {0, 1}, {2}, {0, 1, 2}}, []int{5, 5, 5}, 5},
	}
	for i, each := range cases {
		name := fmt.Sprintf("#%d %v - %v", i, each.joltages, each.buttons)
		t.Run(name, func(t *testing.T) {
			mach := Machine{
				Buttons:  each.buttons,
				Joltages: each.joltages,
			}
			require.Equal(t, each.expected, deriveFewestJoltageClicks7(mach))
		})
	}
}

func TestReduceJoltage(t *testing.T) {
	buttons := [][]int{{0, 1}, {1, 3}, {2}}
	joltages := []int{10, 10, 10, 10}
	expected := make([]int, len(joltages))
	copy(expected, joltages)

	var steps []step
	steps = append(steps, step{buttonID: 0, count: 2})
	expected[0] -= 2
	expected[1] -= 2

	steps = append(steps, step{buttonID: 0, count: 1})
	expected[0]--
	expected[1]--

	steps = append(steps, step{buttonID: 1, count: 4})
	expected[1] -= 4
	expected[3] -= 4

	// EXERCISE
	reduceJoltage(buttons, joltages, steps)

	// VERIFY
	require.Equal(t, expected, joltages)
}

func TestDeriveAngelMaximums(t *testing.T) {
	mach := Machine{
		Buttons: [][]int{
			{0, 1, 3},
			{1, 2},
		},
		Joltages: []int{9, 4, 6, 1},
	}
	maximums := deriveAngelMaximums(mach)
	req := require.New(t)
	req.Equal(
		[]int{1, 1, 4, 1},
		maximums)
}

func TestTable(t *testing.T) {
	req := require.New(t)

	originalJoltages := []int{8, 5, 22, 19}
	mach := Machine{
		Buttons:  [][]int{{2, 3}, {0, 1}, {0, 2}, {3}},
		Joltages: originalJoltages,
	}
	tbl := newTable(mach)
	assertTable := func(clicks, joltages []int, msg string) {
		req.Equal(clicks, tbl.clicks, "clicks "+msg)
		req.Equal(joltages, tbl.joltages, "joltages "+msg)
	}
	originalClicks := []int{19, 5, 8, 19}
	assertTable(originalClicks, originalJoltages, "before any clicks")

	clicked := make([]int, 0, 30)
	states := [][][]int{
		// Click button 0: 2,3.
		{{18, 5, 8, 18}, {8, 5, 21, 18}},
		{{17, 5, 8, 17}, {8, 5, 20, 17}},
		{{16, 5, 8, 16}, {8, 5, 19, 16}},
		{{15, 5, 8, 15}, {8, 5, 18, 15}},
		{{14, 5, 8, 14}, {8, 5, 17, 14}},
		{{13, 5, 8, 13}, {8, 5, 16, 13}},
		{{12, 5, 8, 12}, {8, 5, 15, 12}},
		{{11, 5, 8, 11}, {8, 5, 14, 11}},
		{{10, 5, 8, 10}, {8, 5, 13, 10}},
		{{9, 5, 8, 9}, {8, 5, 12, 9}},
		{{8, 5, 8, 8}, {8, 5, 11, 8}},
		{{7, 5, 8, 7}, {8, 5, 10, 7}},
		{{6, 5, 8, 6}, {8, 5, 9, 6}},
		{{5, 5, 8, 5}, {8, 5, 8, 5}},
		{{4, 5, 7, 4}, {8, 5, 7, 4}},
		{{3, 5, 6, 3}, {8, 5, 6, 3}},
		{{2, 5, 5, 2}, {8, 5, 5, 2}},
		{{1, 5, 4, 1}, {8, 5, 4, 1}},
		{{0, 5, 3, 0}, {8, 5, 3, 0}},
		// Click button 1: 0, 1.
		{{0, 4, 3, 0}, {7, 4, 3, 0}},
		{{0, 3, 3, 0}, {6, 3, 3, 0}},
		{{0, 2, 3, 0}, {5, 2, 3, 0}},
		{{0, 1, 3, 0}, {4, 1, 3, 0}},
		{{0, 0, 3, 0}, {3, 0, 3, 0}},
		// Click button: 2: 0, 2.
		{{0, 0, 2, 0}, {2, 0, 2, 0}},
		{{0, 0, 1, 0}, {1, 0, 1, 0}},
		{{0, 0, 0, 0}, {0, 0, 0, 0}},
	}
	for i, each := range states {
		index, found := tbl.getFirstClickable()
		req.True(found, fmt.Sprint("round index ", i))
		tbl.click(index)
		assertTable(each[0], each[1], fmt.Sprint("after round", i))
		if i < len(states)-1 {
			req.False(tbl.done(), fmt.Sprint("shouldn't be done until the end ", i))
		}
		clicked = append(clicked, index)
	}

	lastIndex, lastFound := tbl.getFirstClickable()
	req.False(lastFound, "should have nothing more to click")
	req.Equal(lastIndex, -1)
	req.True(tbl.done(), "should be done since joltages are all zero")

	for i := len(clicked) - 1; i >= 0; i-- {
		tbl.undo(clicked[i])
	}
	assertTable(originalClicks, originalJoltages, "after undoing everything")
}

func BenchmarkTableClick(b *testing.B) {
	for range b.N {
		joltages := []int{8, 5, 22, 19}
		mach := Machine{
			Buttons:  [][]int{{2, 3}, {0, 1}, {0, 2}, {3}},
			Joltages: joltages,
		}
		tbl := newTable(mach)
		for range 27 {
			i, _ := tbl.getFirstClickable()
			tbl.click(i)
		}
	}
}

func TestShiftButtons(t *testing.T) {
	shared.InitTestLogging(t)
	buttons := [][]int{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
	}
	shifted := shiftButtons(buttons, 2)
	require.Equal(
		t,
		[][]int{
			{2, 3},
			{3, 4},
			{0, 1},
			{1, 2},
		},
		shifted)
}

func TestSetupDemons(t *testing.T) {
	buttonLists := [][]*button{
		{
			{id: 0, joltageIndexes: []int{2, 3}, clicks: 19},
			{id: 1, joltageIndexes: []int{0, 1}, clicks: 5},
			{id: 2, joltageIndexes: []int{0, 2}, clicks: 3},
			{id: 3, joltageIndexes: []int{3}, clicks: -1},
		},
		{
			{id: 0, joltageIndexes: []int{0, 1, 2, 3, 4}, clicks: -1},
			{id: 1, joltageIndexes: []int{0, 3, 4}, clicks: -1},
			{id: 2, joltageIndexes: []int{0, 1, 2, 4, 5}, clicks: 5},
			{id: 3, joltageIndexes: []int{1, 2}, clicks: -1},
		},
		{
			{id: 0, joltageIndexes: []int{0, 3, 4}, clicks: 1},
			{id: 1, joltageIndexes: []int{2, 3, 4}, clicks: 16},
			{id: 2, joltageIndexes: []int{1, 5}, clicks: 6},
			{id: 3, joltageIndexes: []int{3, 4}, clicks: 12},
			{id: 4, joltageIndexes: []int{0, 2, 4, 5}, clicks: 2},
		},
	}
	getButtons := func(topIndex int) func(ids ...int) []*button {
		parentButton := buttonLists[topIndex]
		return func(ids ...int) []*button {
			buttons := make([]*button, 0, len(ids))
			for _, i := range ids {
				buttons = append(buttons, parentButton[i])
			}
			return buttons
		}
	}
	demonLists := [][]*demon{
		{
			{joltageID: 1, joltage: 5, buttons: getButtons(0)(1), lockedToButton: 1},
			{joltageID: 0, joltage: 8, buttons: getButtons(0)(1, 2), lockedToButton: 2},
			{joltageID: 2, joltage: 22, buttons: getButtons(0)(0, 2), lockedToButton: 0},
			{joltageID: 3, joltage: 19, buttons: getButtons(0)(0, 3), lockedToButton: -1},
		},
		{
			{joltageID: 5, joltage: 5, buttons: getButtons(1)(2), lockedToButton: 2},
			{joltageID: 3, joltage: 5, buttons: getButtons(1)(0, 1), lockedToButton: -1},
			{joltageID: 0, joltage: 10, buttons: getButtons(1)(0, 1, 2), lockedToButton: -1},
			{joltageID: 4, joltage: 10, buttons: getButtons(1)(0, 1, 2), lockedToButton: -1},
			{joltageID: 1, joltage: 11, buttons: getButtons(1)(0, 2, 3), lockedToButton: -1},
			{joltageID: 2, joltage: 11, buttons: getButtons(1)(0, 2, 3), lockedToButton: -1},
		},
		{
			{joltageID: 1, joltage: 6, buttons: getButtons(2)(2), lockedToButton: 2},
			{joltageID: 0, joltage: 3, buttons: getButtons(2)(0, 4), lockedToButton: 0},
			{joltageID: 5, joltage: 8, buttons: getButtons(2)(2, 4), lockedToButton: 4},
			{joltageID: 2, joltage: 18, buttons: getButtons(2)(1, 4), lockedToButton: 1},
			{joltageID: 3, joltage: 29, buttons: getButtons(2)(0, 1, 3), lockedToButton: 3},
			{joltageID: 4, joltage: 31, buttons: getButtons(2)(0, 1, 3, 4), lockedToButton: -1},
		},
	}
	for i, machine := range []string{
		"[...#] (2,3) (0,1) (0,2) (3) {8,5,22,19}",
		"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
		"[#...##] (0,3,4) (2,3,4) (1,5) (3,4) (0,2,4,5) {3,6,18,29,31,8}",
	} {
		t.Run(fmt.Sprint(machine), func(t *testing.T) {
			shared.InitTestLogging(t)
			most := min(999, len(demonLists[i]))
			require.Equal(t, demonLists[i][:most], setupDemons(parseMachine(machine))[:most])
		})
	}
}

func BenchmarkSetupDemons(b *testing.B) {
	shared.InitNullLogging()
	mach := parseMachine("[#...##] (0,3,4) (2,3,4) (1,5) (3,4) (0,2,4,5) {3,6,18,29,31,8}")
	for range b.N {
		setupDemons(mach)
	}
}

func TestDeriveJoltageComplexity(t *testing.T) {
	// 3,0,0
	// 2,0,1 ; 2,1,0
	// 1,2,0 ; 1,0,2 ; 1,1,1
	// 0,3,0 ; 0,2,1 ; 0,1,2
	// 0,0,3
	require.Equal(t, 10, deriveJoltageComplexity(3, 3))
	require.Equal(t, 1, deriveJoltageComplexity(1, 1))
}

func TestNk(t *testing.T) {
	for _, each := range [][3]int{
		{3, 3, 2},
		{6, 4, 2},
		{1, 4, 4},
		{210, 10, 4},
	} {
		t.Run(fmt.Sprint(each), func(t *testing.T) {
			require.Equal(t, each[0], nk(each[1], each[2]))
		})
	}
}

func TestDeriveDemons(t *testing.T) {
	mach := Machine{
		Buttons: [][]int{
			{1, 2},
			{0},
			{2},
			{2, 3},
		},
		Joltages: []int{1, 3, 4, 5},
	}
	// EXERCISE
	demons := deriveDemons(mach)
	for _, aDemon := range demons {
		for _, btn := range aDemon.buttons {
			btn.clicks = 0
		}
		aDemon.lockedToButton = 0
	}
	// VERIFY
	require.Equal(
		t,
		[]*demon{
			{
				joltageID: 0,
				joltage:   1,
				buttons: []*button{
					{id: 1, joltageIndexes: []int{0}},
				},
			},
			{
				joltageID: 1,
				joltage:   3,
				buttons: []*button{
					{id: 0, joltageIndexes: []int{1, 2}},
				},
			},
			{
				joltageID: 2,
				joltage:   4,
				buttons: []*button{
					{id: 0, joltageIndexes: []int{1, 2}},
					{id: 2, joltageIndexes: []int{2}},
					{id: 3, joltageIndexes: []int{2, 3}},
				},
			},
			{
				joltageID: 3,
				joltage:   5,
				buttons: []*button{
					{id: 3, joltageIndexes: []int{2, 3}},
				},
			},
		},
		demons)
}

func TestGenerateButtonCombinations(t *testing.T) {
	var tests = []struct {
		name         string
		buttons      [][]int
		expected     [][]int
		joltageCount int
	}{
		{
			name:         "0,1+1,2",
			buttons:      [][]int{{0, 1}, {1, 2}},
			expected:     [][]int{{0, 1}},
			joltageCount: 3,
		},
		{
			name: "8,5,22,19",
			buttons: [][]int{
				{2, 3},
				{0, 1},
				{0, 2},
				{3},
			},
			expected: [][]int{
				{0, 1},
				{0, 1, 2},
				{0, 1, 3},
				{1, 2, 3},
				{0, 1, 2, 3},
			},
			joltageCount: 4,
		},
		{
			name: "sorted",
			buttons: [][]int{
				{0},
				{1, 2, 3},
				{1, 3},
			},
			expected: [][]int{
				{1, 0},
				{1, 2, 0},
			},
			joltageCount: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shared.InitTestLogging(t)
			combos := [][]int{}
			cb := func(combo []int) {
				combos = append(combos, combo)
			}
			// EXERCISE
			generateButtonCombinations(tt.buttons, tt.joltageCount, cb)
			require.Equal(t, tt.expected, combos)
		})
	}
}

func TestPermute(t *testing.T) {
	var perms [][]int
	cb := func(perm []int) bool {
		perms = append(perms, perm)
		return true
	}
	// EXERCISE
	permute([]int{0, 1, 2}, cb)
	// VERIFY
	require.Equal(
		t,
		[][]int{
			{0, 1, 2},
			{0, 2, 1},
			{1, 0, 2},
			{1, 2, 0},
			{2, 1, 0},
			{2, 0, 1},
		},
		perms)
}

func TestMe(t *testing.T) {
	hand := []int{3, 5, 7}
	pack := [][]int{append([]int(nil), hand...)}
	hand[0]++
	require.Equal(t, [][]int{{3, 5, 7}}, pack)
}

func TestIsSuperset(t *testing.T) {
	var tests = []struct {
		name      string
		container []int
		nested    []int
		common    []int
		superset  bool
	}{
		{
			name:      "no overlap",
			container: []int{3, 4, 5},
			nested:    []int{2, 6},
			common:    nil,
			superset:  false,
		},
		{
			name:      "partial overlap",
			container: []int{10, 11, 12},
			nested:    []int{10, 11},
			common:    []int{10, 11},
			superset:  true,
		},
		{
			name:      "total overlap",
			container: []int{20, 21, 22, 23},
			nested:    []int{20, 21, 22, 23},
			common:    []int{20, 21, 22, 23},
			superset:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ass := assert.New(t)
			ass.Equal(tt.superset, isSuperset(tt.container, tt.nested))
			ass.Equal(tt.common, union(tt.container, tt.nested))
		})
	}
}

func TestExtractUnblockingButtonIDs(t *testing.T) {
	buttons := [][]int{
		{0, 1},
		{1, 2},
		{2, 3},
	}
	counter := &clickCounter{
		buttonIDToCount: map[int]int{1: 1, 2: 1},
		count:           2,
	}
	blocked := blockedButton{
		buttonID:              1,
		joltageIndexes:        buttons[1],
		blockedJoltageIndexes: []int{2},
	}
	ids := extractUnblockingButtonIDs(buttons, counter, blocked)
	require.Equal(t, []int{2}, ids)
}

func TestFilterBlockedButtons(t *testing.T) {
	joltages := []int{1, 2, 3, 0, 5}
	buttons := [][]int{
		// Not blocked.
		{0, 1, 2, 4},
		// Partially blocked.
		{2, 3},
		// Partially blocked.
		{3, 4},
		// Completely blocked.
		{3},
	}
	blockedButtons := filterBlockedButtons(joltages, buttons)
	require.Equal(
		t,
		[]blockedButton{
			{
				buttonID:              1,
				joltageIndexes:        buttons[1],
				blockedJoltageIndexes: []int{3},
			},
			{
				buttonID:              2,
				joltageIndexes:        buttons[2],
				blockedJoltageIndexes: []int{3},
			},
			{
				buttonID:              3,
				joltageIndexes:        buttons[3],
				blockedJoltageIndexes: []int{3},
			},
		},
		blockedButtons,
	)
}

func TestSortHops(t *testing.T) {
	source := []hop{
		{total: 10},
		{total: 5},
		{total: 15, clickIndex: 9},
		{total: 15, clickIndex: 64},
	}
	hops := append([]hop(nil), source...)
	sortHops(hops)
	require.Equal(
		t,
		[]hop{source[2], source[3], source[0], source[1]},
		hops)
}

func TestDeriveBestDoubleClicks(t *testing.T) {
	joltages := []int{2, 0, 2}
	buttons := [][]int{
		{0, 1, 2},
		{1},
	}
	counter := &clickCounter{
		buttonIDToCount: map[int]int{1: 1, 0: 1},
		count:           2,
	}
	// EXERCISE
	hops := deriveBestDoubleClicks(joltages, buttons, counter)

	// VERIFY
	// Notice that here we're verifying two things. First is obvious: that unclicking 1 enables
	// clicking 0. The second is that despite the fact that button 1 is blocked, unclicking 0 isn't
	// considered because net benefit would be negative.
	require.Equal(
		t,
		[]hop{
			{
				total:        2,
				clickIndex:   0,
				unclickIndex: 1,
			},
		},
		hops)
}

func TestCreateHive(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	mach := parseMachine("[...#] (2,3) (0,1) (0,2) (3) {8,5,22,19}")
	// EXERCISE
	hive := createHiveMind(mach)
	// VERIFY
	req.Equal(
		[]valueRestrictor{
			{maximum: 19, multiplier: 2, minimum: 19},
			{maximum: 5, multiplier: 2, minimum: 5},
			{maximum: 3, multiplier: 2, minimum: 3},
			{multiplier: 1},
		},
		hive.restrictors,
	)

	hive.derive()
}

func TestIsValid(t *testing.T) {
	run := func(
		name string,
		expected bool,
		sum int,
		buttonIndexes []int,
		clicks []*int) {
		t.Run(name, func(t *testing.T) {
			expr := &expressionRestrictor{
				sum:           sum,
				buttonIndexes: buttonIndexes,
			}
			req := require.New(t)
			shelf := &clickShelf{clicks: clicks}
			req.Equal(expected, expr.isValid(shelf))
		})
	}
	run("just nils", true, 7, []int{1, 3}, []*int{nil, nil, nil, nil})
	run("less no nils", false, 3, []int{0}, toIntPointers(2))
	run("more no nils", false, 3, []int{0}, toIntPointers(4))
	run("less with nil", true, 3, []int{0, 1}, toIntPointers(-1, 2))
	run("equal with nil", true, 3, []int{0, 1, 2}, toIntPointers(1, -1, 2))
	run("equal", true, 2, []int{0, 1}, toIntPointers(1, 1))
}

func toIntPointers(s ...int) []*int {
	result := make([]*int, len(s))
	for i, each := range s {
		if each < 0 {
			result[i] = nil
		} else {
			v := each
			result[i] = &v
		}
	}
	return result
}

func TestExpressionRestrictorDerive(t *testing.T) {
	createDeriveClick := func(
		req *require.Assertions,
		restrictors []expressionRestrictor,
		shelf *clickShelf) func() []derivedValue {
		return func() []derivedValue {
			for i, each := range restrictors {
				req.True(each.isValid(shelf), fmt.Sprint("restrictor isn't valid: ", i))
				derivedValues := each.derive(shelf)
				if len(derivedValues) > 0 {
					return derivedValues
				}
			}
			return nil
		}
	}

	t.Run("only one possible click combination", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)

		mach := parseMachine("[...#] (2,3) (0,1) (0,2) (3) {8,5,22,19}")
		restrictors := createHiveMind(mach).expressionRestrictors
		shelf := createClickShelf(len(mach.Buttons), mach.Joltages, mach.Buttons)
		deriveClick := createDeriveClick(req, restrictors, shelf)

		desiredValues := []derivedValue{
			{i: 2, clicks: 3},
			{i: 1, clicks: 5},
			{i: 0, clicks: 19},
			{i: 3, clicks: 0},
		}
		for i, expected := range desiredValues {
			// EXERCISE
			firstValues := deriveClick()
			// VERIFY
			req.NotEmpty(firstValues, fmt.Sprint("empty: ", i))
			req.Equal([]derivedValue{expected}, firstValues, fmt.Sprint("not equal: ", i))
			for j := range firstValues {
				shelf.setp(firstValues[j].i, &firstValues[j].clicks)
			}
		}
		// VERIFY
		req.Equal(toIntPointers(19, 5, 3, 0), shelf.clicks, "final clicks")
		req.Nil(deriveClick(), "last expected nil")
	})

	t.Run("non-deterministic", func(t *testing.T) {
		req := require.New(t)
		shared.InitTestLogging(t)

		mach := parseMachine("[#..#] (1,3) (2,3) (0,2) (0,3) (0,1,3) (0) {40,22,15,34}")
		//                                                                18, 0,15,12
		//                                                                18, 0, 3, 0
		//                                                                15, 0, 0, 0
		//                                                                 0, 0, 0, 0
		restrictors := createHiveMind(mach).expressionRestrictors
		shelf := createClickShelf(len(mach.Buttons), mach.Joltages, mach.Buttons)
		deriveClick := createDeriveClick(req, restrictors, shelf)

		for i, each := range restrictors {
			req.True(each.isValid(shelf), fmt.Sprint("initial isValid should be true: ", i))
			req.Nil(each.derive(shelf), fmt.Sprint("derive result should be nil: ", i))
		}
		req.Nil(
			deriveClick(),
			"first derive should be nil "+
				"because there are no buttons that have a deterministic click count",
		)

		shelf.setp(4, intp(22))
		derivedValues := deriveClick()
		req.NotEmpty(derivedValues)
		req.Equal(1, len(derivedValues))
		req.Equal(
			derivedValue{},
			derivedValues[0])
		shelf.setp(derivedValues[0].i, &derivedValues[0].clicks)

		derivedValues = deriveClick()
		req.Empty(derivedValues)
		req.Nil(derivedValues)

		shelf.setp(1, intp(12))
		shelf.setp(2, intp(3))
		shelf.setp(3, intp(0))
		shelf.setp(5, intp(15))
		derivedValues = deriveClick()
		req.Nil(derivedValues)
	})
}

func intp(i int) *int {
	return &i
}

func TestCopyIntp(t *testing.T) {
	original := 13
	copied := copyIntp(&original)
	*copied = 17
	require.Equal(t, 13, original)
}

func TestHiveMindSortButtons(t *testing.T) {
	hive := &hiveMind{buttonRates: []int{10, 30, 20, 40}}
	actual := hive.sortButtons()
	require.Equal(t, []int{0, 2, 1, 3}, actual)
}

func TestHiveMindDive(t *testing.T) {
	shared.InitTestLogging(t)
	tests := []struct {
		spec     string
		expected int
	}{
		{spec: "[#..#] (1,3) (2,3) (0,2) (0,3) (0,1,3) (0) {40,22,15,34}", expected: 45},
		{spec: "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}", expected: 10},
		{spec: "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}", expected: 12},
		{spec: "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}", expected: 11},
		{spec: "[..###] (0,1,4) (0,3,4) (1,2,3) (0,1,2,3) {14,6,2,10,12}", expected: 14},
		{
			spec: "[#.....##] " +
				"(1,6,7) (0,2,4,5,6) (0,3,4) (3,4,6) (0,1,2,4,5,6,7) (0,1,7) (0,6,7) (1,4,7) " +
				"{42,39,24,16,52,24,36,42}",
			expected: 63,
		},
	}
	for _, tt := range tests {
		t.Run(tt.spec, func(t *testing.T) {
			req := require.New(t)
			req.Equal(tt.expected, createHiveMind(parseMachine(tt.spec)).derive())
		})
	}
}

func TestClickShelfSum(t *testing.T) {
	req := require.New(t)
	shelf := createClickShelf(2, []int{1, 3, 2}, [][]int{{0, 1}, {1, 2}})
	assertState := func(clickSum, currentSum int, correctSum bool) {
		req.Equal(clickSum, shelf.clickSum)
		req.Equal(currentSum, shelf.currentSum)
		req.Equal(correctSum, shelf.isSumCorrect())
	}
	assertState(0, 0, false)

	shelf.set(0, 1)
	assertState(1, 2, false)

	shelf.set(1, 1)
	assertState(2, 4, false)

	shelf.set(1, 2)
	assertState(3, 6, true)

	shelf.setp(1, nil)
	assertState(1, 2, false)

	shelf.setp(0, nil)
	assertState(0, 0, false)
}

func BenchmarkSubSetHiveMindDive(b *testing.B) {
	shared.InitNullLogging()
	spec := "[#.....##] " +
		"(1,6,7) (0,2,4,5,6) (0,3,4) (3,4,6) (0,1,2,4,5,6,7) (0,1,7) (0,6,7) (1,4,7) " +
		"{42,39,24,16,52,24,36,42}"
	for range b.N {
		if createHiveMind(parseMachine(spec)).derive() != 63 {
			b.FailNow()
		}
	}
}
