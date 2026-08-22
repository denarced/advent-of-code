package aoc2204

import (
	"fmt"

	"github.com/denarced/advent-of-code/shared"
)

func CountOverlaps(lines []string, total bool) int {
	var count int
	for _, line := range lines {
		var af, at, bf, bt int
		parsedCount, err := fmt.Sscanf(line, "%d-%d,%d-%d", &af, &at, &bf, &bt)
		if err != nil {
			shared.Logger.Error("Failed to parse line.", "err", err)
			panic(err)
		}
		if parsedCount != 4 {
			shared.Logger.Error("Expected count of parsed numbers.", "count", parsedCount)
			panic("count should be 4")
		}
		if total {
			if isEnclosed(af, at, bf, bt) {
				count++
			}
		} else if overlaps(af, at, bf, bt) {
			count++
		}
	}
	return count
}

func isEnclosed(af, at, bf, bt int) bool {
	is := func(inFrom, inTo, outFrom, outTo int) bool {
		return outFrom <= inFrom && inTo <= outTo
	}
	return is(af, at, bf, bt) || is(bf, bt, af, at)
}

func overlaps(af, at, bf, bt int) bool {
	if af < bf {
		return at >= bf
	}
	return bt >= af
}
