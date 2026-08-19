package aoc2201

import (
	"strconv"
	"strings"

	"github.com/denarced/gent"
)

func DeriveHeaviest(lines []string) int {
	weights := []int{0}
	maximus := -1
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			maximus = max(maximus, weights[len(weights)-1])
			weights = append(weights, 0)
			continue
		}
		value := gent.OrPanic2(strconv.Atoi(trimmed))("str-to-int failed")
		weights[len(weights)-1] += value
	}
	return maximus
}
