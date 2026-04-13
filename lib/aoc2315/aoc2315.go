package aoc2315

import (
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func SumHashes(lines []string) int {
	pieces := parseLines(lines)
	shared.Logger.Info("Sum hashes.", "count", len(pieces))
	var sum int
	for _, each := range pieces {
		sum += hash(each)
	}
	shared.Logger.Info("Hashes summed.", "sum", sum)
	return sum
}

func hash(s string) (h int) {
	for _, r := range s {
		h = ((h + int(r)) * 17) % 256
	}
	return
}

func parseLines(lines []string) []string {
	var result []string
	for _, line := range lines {
		for _, each := range strings.Split(line, ",") {
			if each != "" {
				result = append(result, each)
			}
		}
	}
	return result
}
