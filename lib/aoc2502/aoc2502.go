package aoc2502

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumInvalidIDs(line string, twice bool) int64 {
	var sum int64
	shared.Logger.Info("Derive sum of invalid IDs.", "twice", twice)
	var wg sync.WaitGroup
	for _, each := range splitToRanges(line) {
		wg.Add(1)
		go func(intr intRange) {
			defer wg.Done()
			for n := intr.from; n <= intr.to; n++ {
				maxSplit := 2
				if !twice {
					maxSplit = deriveIntLength(n)
				}
				if breaks(n, 2, maxSplit) {
					shared.Logger.Info("Invalid ID found.", "ID", n)
					sum += n
				}
			}
		}(each)
	}
	wg.Wait()
	shared.Logger.Info("Sum calculated.", "sum", sum)
	return sum
}

func breaks(n int64, minSplit, maxSplit int) bool {
	if n < 10 {
		return false
	}
	pieces := splitInt(n)
	for i := minSplit; i <= maxSplit; i++ {
		length := len(pieces)
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

		if allIntsEqual(splitInts(pieces, i)) {
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

func deriveIntLength(n int64) int {
	if n == 0 {
		return 1
	}
	length := 0
	for n > 0 {
		n /= 10
		length++
	}
	return length
}

func splitInt(n int64) []int {
	if n < 0 {
		panic("negative number")
	}
	if n == 0 {
		return []int{0}
	}
	i := deriveIntLength(n)
	digits := make([]int, i)
	i--
	for n > 0 {
		digits[i] = int(n % 10)
		i--
		n /= 10
	}
	return digits
}

func splitInts(pieces []int, n int) [][]int {
	inc := len(pieces) / n
	result := make([][]int, 0, len(pieces)/inc)
	for i := 0; i < len(pieces); i += inc {
		result = append(result, pieces[i:i+inc])
	}
	return result
}

func allIntsEqual(s [][]int) bool {
	if len(s) < 2 {
		return false
	}
	first := s[0]
	for _, each := range s[1:] {
		if !slices.Equal(first, each) {
			return false
		}
	}
	return true
}
