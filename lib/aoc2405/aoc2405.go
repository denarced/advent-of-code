package aoc2405

import (
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumCorrectMiddlePageNumbers(lines []string) int {
	return sumMiddlePageNumbers(lines, true)
}

func SumIncorrectMiddlePageNumbers(lines []string) int {
	return sumMiddlePageNumbers(lines, false)
}

func sumMiddlePageNumbers(lines []string, correct bool) int {
	rules, pages := toRulesAndPages(lines)
	shared.Logger.Info(
		"Sum middle page numbers.",
		"rule count",
		len(rules),
		"page list count",
		len(pages),
		"correct",
		correct,
	)
	sum := 0
	filtered := gent.Filter(
		pages,
		func(s []int) bool {
			return isSortedAccordingToRules(rules, s) == correct
		})
	shared.Logger.Info(
		"Page lists filtered.",
		"page list count",
		len(filtered))
	for _, each := range filtered {
		if !correct {
			each = sortWithRules(rules, each)
		}
		middle := each[len(each)/2]
		shared.Logger.Debug("Add middle page to the sum.", "pages", each, "middle", middle)
		sum += middle
	}
	shared.Logger.Info("Middle page numbers summed.", "sum", sum)
	return sum
}

func toRulesAndPages(lines []string) ([][]int, [][]int) {
	toInts := func(s []string) []int {
		var ints []int
		for _, each := range s {
			i, err := strconv.Atoi(each)
			if err != nil {
				panic(err)
			}
			ints = append(ints, i)
		}
		return ints
	}
	contains := func(sub string) func(s string) bool {
		return func(s string) bool {
			return strings.Contains(s, sub)
		}
	}
	split := func(sep string) func(s string) []string {
		return func(s string) []string {
			return strings.Split(s, sep)
		}
	}
	filterAndSplit := func(sep string) [][]int {
		return gent.Map(
			gent.Map(
				gent.Filter(lines, contains(sep)),
				split(sep)),
			toInts)
	}
	return filterAndSplit("|"), filterAndSplit(",")
}

func isSortedAccordingToRules(rules [][]int, pages []int) bool {
	sorted := sortWithRules(rules, pages)
	return slices.Equal(pages, sorted)
}

func sortWithRules(rules [][]int, pages []int) []int {
	dup := append([]int{}, pages...)
	slices.SortStableFunc(
		dup,
		func(a, b int) int {
			rule := findRelevantRule(rules, a, b)
			if rule == nil {
				return 0
			}
			if rule[0] == a {
				return -1
			}
			return 1
		})
	return dup
}

func findRelevantRule(rules [][]int, a, b int) []int {
	var relevant []int
	for _, each := range rules {
		if a == each[0] && b == each[1] || a == each[1] && b == each[0] {
			return each
		}
	}
	return relevant
}
