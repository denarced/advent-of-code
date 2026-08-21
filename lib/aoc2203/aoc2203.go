package aoc2203

import (
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumPriorities(lines []string, groupSearch bool) int {
	var table [][]string
	if groupSearch {
		var current []string
		for _, each := range lines {
			current = append(current, each)
			if len(current) >= 3 {
				table = append(table, current)
				current = nil
			}
		}
		if len(current) != 0 {
			shared.Logger.Error(
				"Current slice has content.",
				"size", len(current),
				"table size", len(table))
			panic("current has content")
		}
	} else {
		for _, line := range lines {
			contents := splitLine(line)
			table = append(table, contents)
		}
	}
	var sum int
	for _, each := range table {
		commonItem := deriveCommon(each)
		var added int
		if commonItem <= 'Z' {
			added = int(commonItem-'A') + 27
		} else {
			added = int(commonItem-'a') + 1
		}
		sum += added
		shared.Logger.Info("Add new value.", "common", commonItem, "value", added)
	}
	shared.Logger.Info("Sum priorities summed.", "sum", sum)
	return sum
}

func splitLine(line string) []string {
	mid := len(line) / 2
	return []string{line[:mid], line[mid:]}
}

func deriveCommon(contents []string) rune {
	uniques := deriveUniques(contents)
	var found rune
	uniques[0].ForEach(func(each rune, stop func()) {
		for i := 1; i < len(uniques); i++ {
			if !uniques[i].Has(each) {
				return
			}
		}
		found = each
		stop()
	})
	return found
}

func deriveUniques(contents []string) []*gent.Set[rune] {
	sets := make([]*gent.Set[rune], len(contents))
	for i, each := range contents {
		current := gent.NewSet[rune]()
		sets[i] = current
		for _, c := range each {
			current.Add(c)
		}
	}
	return sets
}
