package aoc2505

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func CountFreshAvailableIngredients(lines []string) int {
	freshIDRanges, availableIDs := parseLines(lines)
	var count int
	for _, each := range availableIDs {
		for _, aRange := range freshIDRanges {
			if aRange.contains(each) {
				count++
				break
			}
		}
	}
	return count
}

type intRange struct {
	from int
	to   int
}

func (v intRange) contains(value int) bool {
	return v.from <= value && value <= v.to
}

func parseLines(lines []string) ([]intRange, []int) {
	toInt := func(s string) int {
		value, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return value
	}

	var ranges []intRange
	var availableIDs []int
	var emptyLineSeen bool
	for _, each := range lines {
		if each == "" {
			emptyLineSeen = true
			continue
		}
		if !emptyLineSeen {
			pieces := strings.Split(each, "-")
			if len(pieces) != 2 {
				panic(fmt.Sprintf("invalid line: %s", each))
			}
			ranges = append(ranges, intRange{
				from: toInt(pieces[0]),
				to:   toInt(pieces[1]),
			})
			continue
		}
		value := toInt(each)
		availableIDs = append(availableIDs, value)
	}
	return ranges, availableIDs
}

func CountFreshIngredients(lines []string) int {
	freshIDRanges, _ := parseLines(lines)
	merged := mergeIntRanges(freshIDRanges)
	var count int
	for _, each := range merged {
		count += each.to - each.from + 1
	}
	return count
}

func mergeIntRanges(ranges []intRange) []intRange {
	if len(ranges) == 0 {
		return ranges
	}

	// Use pointers so we can set latter pair to nil after merge.
	working := make([]*intRange, len(ranges))
	for i, each := range ranges {
		copied := each
		working[i] = &copied
	}

	merged := true
	// We always need to have a limit of some kind to prevent eternal loops.
	safety := 1000
	for merged && safety > 0 {
		merged = runMergeRound(working)
		safety--
	}
	if safety <= 0 {
		panic("safety reached")
	}

	finalCount := func() int {
		var count int
		for _, each := range working {
			if each != nil {
				count++
			}
		}
		return count
	}()
	result := make([]intRange, 0, finalCount)
	for _, each := range working {
		if each == nil {
			continue
		}
		result = append(result, *each)
	}
	// Just to make results deterministic for tests.
	slices.SortFunc(result, func(a, b intRange) int {
		if a.from < b.from {
			return -1
		}
		if b.from < a.from {
			return 1
		}
		if a.to < b.to {
			return -1
		}
		if b.to < a.to {
			return 1
		}
		return 0
	})
	return result
}

func runMergeRound(ranges []*intRange) bool {
	length := len(ranges)
	merged := false
	for i := 0; i < length-1; i++ {
		if ranges[i] == nil {
			continue
		}
		for j := i + 1; j < length; j++ {
			if ranges[j] == nil {
				continue
			}
			joined := joinIntRanges(ranges[i], ranges[j])
			if joined != nil {
				ranges[i] = joined
				ranges[j] = nil
				merged = true
			}
		}
	}
	return merged
}

func joinIntRanges(first, second *intRange) *intRange {
	if first.to < second.from || second.to < first.from {
		return nil
	}
	start := min(first.from, second.from)
	end := max(first.to, second.to)
	return &intRange{from: start, to: end}
}
