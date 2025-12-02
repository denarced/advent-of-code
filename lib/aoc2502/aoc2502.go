package aoc2502

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumInvalidIDs(line string, twice bool) int64 {
	var sum int64
	shared.Logger.Info("Derive sum of invalid IDs.", "twice", twice)
	for _, each := range splitToRanges(line) {
		for n := each.from; n <= each.to; n++ {
			maxSplit := 2
			if !twice {
				maxSplit = len(strconv.FormatInt(n, 10))
			}
			if breaks(n, 2, maxSplit) {
				shared.Logger.Info("Invalid ID found.", "ID", n)
				sum += n
			}
		}
	}
	shared.Logger.Info("Sum calculated.", "sum", sum)
	return sum
}

func breaks(n int64, minSplit, maxSplit int) bool {
	if n < 10 {
		return false
	}
	s := strconv.FormatInt(n, 10)
	for i := minSplit; i <= maxSplit; i++ {
		length := len(s)
		if i > length {
			break
		}
		if length%i != 0 {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Skip split.", "piece count", i)
			}
			continue
		}
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Test split.", "piece count", i)
		}

		if allEqual(splitString(s, i)) {
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Breaks.", "n", n, "split count", i)
			}
			return true
		}
	}
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Doesn't break.", "n", n, "min", minSplit, "max", maxSplit)
	}
	return false
}

func allEqual(s []string) bool {
	if len(s) < 2 {
		return false
	}
	first := s[0]
	for _, each := range s[1:] {
		if first != each {
			return false
		}
	}
	return true
}

func splitString(s string, n int) []string {
	var pieces []string
	inc := len(s) / n
	for i := 0; i < len(s); i += inc {
		pieces = append(pieces, s[i:i+inc])
	}
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Result of string split.", "string", s, "n", n, "pieces", pieces)
	}
	return pieces
}

type intRange struct {
	from int64
	to   int64
}

func splitToRanges(line string) []intRange {
	var ranges []intRange
	for _, each := range strings.Split(line, ",") {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		pieces := strings.Split(trimmed, "-")
		if len(pieces) != 2 {
			panic(fmt.Sprintf("invalid range: %s", trimmed))
		}
		var r intRange
		r.from = gent.OrPanic2(strconv.ParseInt(pieces[0], 10, 64))("convert start of range")
		r.to = gent.OrPanic2(strconv.ParseInt(pieces[1], 10, 64))("convert end of range")
		ranges = append(ranges, r)
	}
	return ranges
}
