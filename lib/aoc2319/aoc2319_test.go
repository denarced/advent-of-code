package aoc2319

import (
	"slices"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSumRatings(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err, "failed to read test data")

	// EXERCISE
	sum := SumRatings(lines)

	// VERIFY
	req.Equal(19_114, sum)
}

func TestParseWorkflows(t *testing.T) {
	req := require.New(t)

	// EXERCISE
	workflows := parseWorkflows([]string{
		"px{a<2006:qkq,m>2090:A,rfg}",
	})

	// VERIFY
	req.Equal(
		map[string]workflow{
			"px": {
				name: "px",
				specs: []spec{
					{attr: "a", dest: "qkq", value: 2006, less: true},
					{attr: "m", dest: "A", value: 2090},
					{dest: "rfg", endComplete: true},
				},
			},
		},
		workflows)
}

func TestParseParts(t *testing.T) {
	req := require.New(t)

	// EXERCISE
	parts := parseParts([]string{
		"{x=787,m=2655,a=1222,s=2876}",
	})

	// VERIFY
	req.Equal(
		[]part{
			map[string]int{"x": 787, "m": 2655, "a": 1222, "s": 2876},
		},
		parts)
}

func TestNegotiate(t *testing.T) {
	run := func(name string, lines []string, genPolicy func() policy, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// Just for in.txt, to get rid of parts.
			if index := slices.Index(lines, ""); index > 0 {
				lines = lines[:index]
			}

			// EXERCISE
			result := Negotiate(
				append(
					lines,
					"",
					"{x=787,m=2655,a=1222,s=2876}"),
				genPolicy)

			// VERIFY
			req.Equal(expected, result)
		})
	}

	run(
		"two dimensions",
		[]string{
			"in{a<3:A,sec}",
			"sec{b>1:A}",
		},
		func() policy {
			return policy{
				"a": restriction{low: 1, high: 3},
				"b": restriction{low: 1, high: 3},
			}
		},
		8)

	run(
		"three dimensions #1",
		[]string{
			"in{a<2:jinn,b<3:kill,R}",
			"jinn{b<2:kill,R}",
			"kill{c<2:A,R}",
		},
		func() policy {
			return policy{
				"a": restriction{low: 1, high: 2},
				"b": restriction{low: 1, high: 2},
				"c": restriction{low: 1, high: 2},
			}
		},
		3)

	run(
		"three dimensions #2",
		[]string{
			"in{a<4:j,a>6:k,R}",
			"j{b<7:k,A}",
			"k{c<2:l,R}",
			"l{A}",
		},
		func() policy {
			return policy{
				"a": restriction{low: 1, high: 10},
				"b": restriction{low: 1, high: 10},
				"c": restriction{low: 1, high: 10},
			}
		},
		178)

	run(
		"reproduce in.txt error",
		[]string{
			"in{qqz}",
			"qqz{m<4:hdj,A}",
			"hdj{m>1:A,pv}",
			"pv{a>2:R,A}",
		},
		func() policy {
			return policy{
				"s": restriction{low: 1, high: 6},
				"m": restriction{low: 1, high: 6},
				"a": restriction{low: 1, high: 6},
			}
		},
		192)

	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	require.NoError(t, err, "failed to read test data")
	run("in.txt", lines, nil, 167_409_079_868_000)
}

func TestPolicyIntersection(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	first := policy{
		"a": restriction{low: 2, high: 12},
		"b": restriction{low: 8, high: 14},
		"c": restriction{low: 10, high: 12},
	}
	second := policy{
		"a": restriction{low: 1, high: 13},
		"b": restriction{low: 9, high: 13},
		"c": restriction{low: 12, high: 14},
	}
	// EXERCISE
	result := first.intersection(second)

	// VERIFY
	req.Equal(
		policy{
			"a": restriction{low: 2, high: 12},
			"b": restriction{low: 9, high: 13},
			"c": restriction{low: 12, high: 12},
		},
		result)
}

func TestPierce(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	ass := assert.New(t)

	workflows := map[string]workflow{
		"in": {
			name: "in",
			specs: []spec{
				{attr: "a", value: 10, dest: "pty"},
				{dest: "rab", endComplete: true},
			},
		},
		"pty": {
			name: "pty",
			specs: []spec{
				{attr: "b", dest: "A", value: 20, less: true},
				{dest: "R", endComplete: true},
			},
		},
		"rab": {
			name: "rab",
			specs: []spec{
				{attr: "b", dest: "A", value: 15},
				{dest: "R", endComplete: true},
			},
		},
	}
	var actual [][]comparison
	cb := func(comparisons *shared.Link[comparison]) {
		var comps []comparison
		for comparisons != nil {
			comps = append(comps, comparisons.Item)
			comparisons = comparisons.Parent
		}
		actual = append(actual, comps)
	}

	// EXERCISE
	pierce(workflows, "in", nil, cb)

	// VERIFY
	// #1
	//     in - a>10
	//     pty - b<20
	//     A
	// #2
	//     in - a>10
	//     pty - end
	//     R
	// #3  in - a<11
	//     rab - b>15
	//     A
	// #4  in - a<11
	//     rab - b<16
	//     R
	// expected := [][]comparison{
	// 	{
	// 		{attr: "b", val: 20, less: true},
	// 		{attr: "a", val: 10},
	// 	},
	// 	{
	// 		{attr: "b", val: 15},
	// 	},
	// }
	expected := [][]comparison{
		{
			{attr: "b", val: 20, less: true},
			{attr: "a", val: 10},
		},
		{
			{attr: "b", val: 15},
			{attr: "a", val: 11, less: true},
		},
	}
	req.Equal(expected, actual)
	if !ass.Equal(expected, actual) {
		t.Log("Expected")
		for i, each := range expected {
			t.Logf("%d %+v", i, each)
		}
		t.Log("Actual")
		for i, each := range actual {
			t.Logf("%d %+v", i, each)
		}
	}
}

func TestPierceDefect(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	var workflows map[string]workflow = map[string]workflow{
		"in": {
			name:  "in",
			specs: []spec{{dest: "qqz", endComplete: true}},
		},
		"qqz": {
			name: "qqz",
			specs: []spec{
				{attr: "m", dest: "hdj", value: 4, less: true},
				{dest: "A", endComplete: true},
			},
		},
		"hdj": {
			name: "hdj",
			specs: []spec{
				{attr: "m", dest: "A", value: 1},
				{dest: "pv", endComplete: true},
			},
		},
		"pv": {
			name: "pv",
			specs: []spec{
				{attr: "a", dest: "R", value: 2},
				{dest: "A", endComplete: true},
			},
		},
	}
	var result [][]comparison
	// EXERCISE
	pierce(workflows, "in", nil, func(link *shared.Link[comparison]) {
		var comparisons []comparison
		for link != nil {
			c := link.Item
			comparisons = append(comparisons, c)
			link = link.Parent
		}
		result = append(result, comparisons)
	})

	// VERIFY
	req.Equal(
		stringifyComparisonSlices([][]comparison{
			{
				{attr: "m", val: 1},
				{attr: "m", val: 4, less: true},
			},
			{
				{attr: "a", val: 3, less: true},
				{attr: "m", val: 2, less: true},
				{attr: "m", val: 4, less: true},
			},
			{
				{attr: "m", val: 3, less: false},
			},
		}),
		stringifyComparisonSlices(result))
}

func stringifyComparisonSlices(comps [][]comparison) []string {
	result := make([]string, len(comps))
	for i := range comps {
		result[i] = stringifyComparisons(comps[i])
	}
	return result
}

func stringifyComparisons(comps []comparison) string {
	result := make([]string, len(comps))
	for i := range comps {
		result[i] = comps[i].String()
	}
	return strings.Join(result, ",")
}

func TestCutRestriction(t *testing.T) {
	run := func(name string, comp comparison, expectedRest restriction, expectedSuccess bool) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			result, success := cutRestriction(restriction{low: 10, high: 19}, comp)

			// VERIFY
			req.Equal(expectedSuccess, success)
			if expectedSuccess {
				req.Equal(expectedRest, result)
			}
		})
	}

	run("low-1 - no effect", comparison{val: 8, less: false}, restriction{low: 10, high: 19}, true)
	run("low - no effect", comparison{val: 9, less: false}, restriction{low: 10, high: 19}, true)
	run("low+1", comparison{val: 10, less: false}, restriction{low: 11, high: 19}, true)
	run("high-1", comparison{val: 19, less: true}, restriction{low: 10, high: 18}, true)
	run("high - no effect", comparison{val: 20, less: true}, restriction{low: 10, high: 19}, true)
	run("high+1 - no effect", comparison{val: 21, less: true}, restriction{low: 10, high: 19}, true)

	run("high < low", comparison{val: 10, less: true}, restriction{}, false)
	run("low > high", comparison{val: 19}, restriction{}, false)
}

func TestDerivePolicy(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	link := shared.AddLink(nil, comparison{attr: "a", val: 5, less: false})
	link = shared.AddLink(link, comparison{attr: "a", val: 15, less: true})
	// EXERCISE
	pol := derivePolicy(link, func() policy {
		return policy{
			"a": restriction{low: 1, high: 20},
			"b": restriction{low: 1, high: 20},
		}
	})

	// VERIFY
	req.NotNil(pol)
	req.Equal(
		policy{
			"a": restriction{low: 6, high: 14},
			"b": restriction{low: 1, high: 20},
		},
		pol)
}

func sumNumbers(lines []string) int {
	var sum int
	var on bool
	for _, each := range lines {
		var lineTotal, current int
		for _, c := range each {
			digit := '0' <= c && c <= '9'
			if digit {
				current = current*10 + int(c-'0')
				on = true
			} else {
				if on {
					lineTotal += current
					current = 0
					on = false
				}
			}
		}
		sum += lineTotal + current
		if shared.IsDebugEnabled() {
			shared.Logger.Debug(
				"Add line total to the sum.",
				"line", each,
				"sum now", sum,
				"added", lineTotal)
		}
	}
	return sum
}

func TestSumNumbers(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	sum := sumNumbers([]string{
		"",
		"1",
		"a2",
		"4b",
		"c8d16",
	})

	req.Equal(
		1+2+4+8+16,
		sum)
}

func TestCountLand(t *testing.T) {
	run := func(name string, policies []policy, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			land := countLand(policies)

			// VERIFY
			req.Equal(expected, land)
		})
	}

	rest := func(low, high int) restriction {
		return restriction{low: low, high: high}
	}
	run("nil", nil, 0)
	run(
		"one policy",
		[]policy{{"a": rest(3, 5), "b": rest(10, 12)}},
		9)
	run(
		"full overlap",
		[]policy{
			{"a": rest(3, 5), "b": rest(10, 12)},
			{"a": rest(3, 5), "b": rest(10, 12)},
		},
		9)
	run(
		"shift by 1 in 1 dimension",
		[]policy{
			{"a": rest(3, 5), "b": rest(10, 12)},
			{"a": rest(4, 6), "b": rest(10, 12)},
		},
		12)
	run(
		"shift by 1 in 2 dimensions",
		[]policy{
			{"a": rest(3, 5), "b": rest(10, 12)},
			{"a": rest(4, 6), "b": rest(11, 13)},
		},
		14)
	run(
		"overlap of one cell",
		[]policy{
			{"a": rest(1, 3), "b": rest(101, 103)},
			{"a": rest(3, 5), "b": rest(103, 105)},
		},
		17)
}
