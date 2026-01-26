package aoc2510

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestDeriveFewestClicks(t *testing.T) {
	readLines := func(req *require.Assertions) []string {
		lines, err := inr.ReadPath("testdata/in.txt")
		req.NoError(err, "failed to read test data")
		return lines
	}
	t.Run("indicator", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		req.Equal(7, DeriveFewestClicks(readLines(req), true))
	})
	t.Run("joltage", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)
		req.Equal(33, DeriveFewestClicks(readLines(req), false))
	})
}

func BenchmarkDeriveFewestClicks(b *testing.B) {
	shared.InitNullLogging()
	req := require.New(b)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	for range b.N {
		DeriveFewestClicks(lines, false)
	}
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
