package aoc2309

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumExtrapolatedValues(lines []string) int {
	table := parseLines(lines)
	shared.Logger.Info("Extrapolate values.", "table size", len(table))
	var sum int
	for _, each := range table {
		result := extrapolate(each)
		sum += result
		shared.Logger.Info("Sum extrapolated.", "sum", result, "values", each, "total", sum)
	}
	shared.Logger.Info("Total sum extrapolated.", "sum", sum)
	return sum
}

func parseLines(lines []string) [][]int {
	result := make([][]int, len(lines))
	for i, each := range lines {
		result[i] = gent.Map(strings.Fields(each), func(s string) int {
			return gent.OrPanic2(strconv.Atoi(s))("failed to convert number")
		})
	}
	return result
}

func isAllZeroes(ints []int) bool {
	for _, each := range ints {
		if each != 0 {
			return false
		}
	}
	return true
}

func extrapolate(values []int) int {
	if len(values) < 3 {
		panic(fmt.Sprintf("too few values: %v", values))
	}

	table := [][]int{values}
	for !isAllZeroes(getLast(table)) {
		last := getLast(table)
		added := make([]int, len(last)-1)
		for i := 0; i < len(last)-1; i++ {
			current := last[i]
			next := last[i+1]
			added[i] = next - current
		}
		table = append(table, added)
	}
	table[len(table)-1] = append(getLast(table), 0)
	for i := len(table) - 2; i >= 0; i-- {
		top := table[i]
		bottom := table[i+1]
		added := getLast(top) + getLast(bottom)
		top = append(top, added)
		table[i] = top
	}
	return getLast(table[0])
}

func getLast[T any](values []T) (t T) {
	if len(values) == 0 {
		return
	}
	return values[len(values)-1]
}
