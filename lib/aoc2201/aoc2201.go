package aoc2201

import (
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/gent"
)

func DeriveHeaviest(lines []string, limit int) int {
	weights := []int{0}
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			weights = append(weights, 0)
			continue
		}
		value := gent.OrPanic2(strconv.Atoi(trimmed))("str-to-int failed")
		weights[len(weights)-1] += value
	}
	slices.Sort(weights)
	lastIndex := len(weights) - 1
	var total int
	for i := lastIndex; i > lastIndex-limit; i-- {
		if i < 0 {
			break
		}
		total += weights[i]
	}
	return total
}
