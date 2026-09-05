package aoc2213

import (
	"encoding/json"
	"log/slog"
	"slices"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	correctOrder orderType = iota
	incorrectOrder
	unknownOrder
)

func SumSortedIndexes(lines []string) int {
	shared.Logger.Info("Derive sum of sorted indexes.", "line count", len(lines))
	pairSets := toPairs(parseLinesToPairs(lines))
	var sum int
	for i, pair := range pairSets {
		inOrder := isInOrder(pair[0], pair[1])
		if inOrder == unknownOrder {
			panic("impossible: here order should be defined")
		}
		if inOrder == correctOrder {
			sum += i + 1
		}
	}
	return sum
}

func SumValues(values []any) int {
	var sum int
	for _, each := range values {
		switch typed := each.(type) {
		case int:
			sum += typed
		case []any:
			sum += SumValues(typed)
		default:
			panic("unknown type")
		}
	}
	return sum
}

func parseLinesToPairs(lines []string) [][]string {
	var linePairs [][]string
	var current []string
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			if len(current) != 2 {
				shared.Logger.Error("Current should have 2 lines.", "current", current)
				panic("current should have 2 lines")
			}
			linePairs = append(linePairs, current)
			current = nil
		} else {
			current = append(current, trimmed)
		}
	}
	if len(current) > 0 {
		if len(current) != 2 {
			shared.Logger.Error("Final current should have 2 lines.", "current", current)
			panic("final current should have 2 lines")
		}
		linePairs = append(linePairs, current)
		current = nil
	}
	return linePairs
}

func toPairs(lineSets [][]string) [][2][]any {
	var annies [][2][]any
	for _, each := range lineSets {
		first := toAny(each[0])
		second := toAny(each[1])
		pair := [2][]any{first, second}
		annies = append(annies, pair)
	}
	return annies
}

func toAny(line string) []any {
	var result []any
	err := json.Unmarshal([]byte(line), &result)
	if err != nil {
		shared.Logger.Error("Failed to unmarshal a JSON line.", "err", err, "line", line)
		panic(err)
	}
	return result
}

type orderType int

func (v orderType) String() string {
	switch v {
	case correctOrder:
		return "order"
	case incorrectOrder:
		return "inorder"
	case unknownOrder:
		return "unknown"
	default:
		panic("unknown order type")
	}
}

func isInOrder(first, second []any) (sub orderType) {
	// 0 and 1: correct order.
	// 1 and 0: wrong order.
	// 1 and 1: continue.
	// [1] and [1,2]: correct order.
	// [1,2] and [1]: wrong order.
	// [1,2] and [1,2]: continue.
	// 1 and [1,2] -> [1] and [1,2].
	maxLength := max(len(first), len(second))
	for i := 0; i < maxLength; i++ {
		if i >= len(first) && i < len(second) {
			return correctOrder
		}
		if i < len(first) && i >= len(second) {
			return incorrectOrder
		}
		firstInt, firstIsFloat := first[i].(float64)
		secondInt, secondIsFloat := second[i].(float64)
		var debugLogger *slog.Logger
		if shared.IsDebugEnabled() {
			debugLogger = shared.Logger.With("first", first[i], "second", second[i])
		}
		if firstIsFloat && secondIsFloat {
			if firstInt < secondInt {
				return correctOrder
			}
			if firstInt > secondInt {
				return incorrectOrder
			}
			continue
		}
		if firstIsFloat && !secondIsFloat {
			sub = isInOrder([]any{firstInt}, second[i].([]any))
			if sub != unknownOrder {
				if debugLogger != nil {
					debugLogger.Debug("Order result.", "result", sub)
				}
				return sub
			}
			continue
		}
		if !firstIsFloat && secondIsFloat {
			sub = isInOrder(first[i].([]any), []any{secondInt})
			if sub != unknownOrder {
				if debugLogger != nil {
					debugLogger.Debug("Order result.", "result", sub)
				}
				return sub
			}
			continue
		}
		sub = isInOrder(first[i].([]any), second[i].([]any))
		if sub != unknownOrder {
			if debugLogger != nil {
				debugLogger.Debug("Order result.", "result", sub)
			}
			return sub
		}
	}

	return unknownOrder
}

func SortAll(lines []string) int {
	shared.Logger.Info("Sort to find divider packets.", "line count", len(lines))
	parsed := parseLines(lines)
	dividerValues := []float64{2, 6}
	shared.Logger.Info(
		"Lines parsed, about to add dividers.",
		"count", len(parsed),
		"dividers", dividerValues)
	for _, each := range dividerValues {
		line := []any{[]any{each}}
		parsed = append(parsed, line)
	}
	slices.SortFunc(parsed, func(a, b []any) int {
		order := isInOrder(a, b)
		switch order {
		case correctOrder:
			return -1
		case incorrectOrder:
			return 1
		default:
			return 0
		}
	})
	var indexes []int
	for i, each := range parsed {
		if shared.IsDebugEnabled() {
			s := gent.OrPanic2(json.Marshal(each))("failed to marshal to JSON")
			shared.Logger.Debug("Sorted.", "i", i, "line", string(s))
		}
		if len(each) == 1 {
			arr, ok := each[0].([]any)
			if ok {
				if len(arr) == 1 {
					f, fOk := arr[0].(float64)
					if fOk && (int(f) == 2 || int(f) == 6) {
						indexes = append(indexes, i+1)
					}
				}
			}
		}
	}
	if len(indexes) != 2 {
		shared.Logger.Error("Broken code, should have 2 indexes.", "len", len(indexes))
	}
	product := indexes[0] * indexes[1]
	shared.Logger.Info("Sort result.", "product", product, "indexes", indexes)
	return product
}

func parseLines(lines []string) [][]any {
	var lists [][]any
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		lists = append(lists, toAny(trimmed))
	}
	return lists
}
