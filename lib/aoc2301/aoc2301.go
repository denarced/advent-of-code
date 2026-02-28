package aoc2301

import "github.com/denarced/advent-of-code/shared"

const (
	prefixTarget seekTarget = iota
	suffixTarget
)

type seekTarget int

var (
	digits = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
	numbers = []string{
		"zero",
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
	}
)

func SumCalibrationValues(lines []string, justDigits bool) int {
	var sum int
	for _, each := range lines {
		sum += parseDigit(each, justDigits)
	}
	return sum
}

func parseDigit(s string, justDigits bool) int {
	shared.Logger.Debug("Parse digit.", "value", s, "just digits", justDigits)
	var first, last int
	for i := range len(s) {
		sub := s[i:]
		d, found := parseDigitIn(sub, prefixTarget, justDigits)
		if found {
			first = d
			break
		}
	}
	if first <= 0 {
		shared.Logger.Error("First digit not found.", "value", s, "just digits", justDigits)
		panic("first not found")
	}
	for i := len(s); i > 0; i-- {
		sub := s[0:i]
		d, found := parseDigitIn(sub, suffixTarget, justDigits)
		if found {
			last = d
			break
		}
	}
	if last <= 0 {
		shared.Logger.Error("Last digit not found.", "value", s, "just digits", justDigits)
		panic("last digit not found")
	}
	return first*10 + last
}

func parseDigitIn(s string, prefix seekTarget, justDigits bool) (int, bool) {
	sub := s[len(s)-1:]
	if prefix == prefixTarget {
		sub = s[0:1]
	}
	for i, each := range digits {
		if i == 0 {
			continue
		}
		if each == sub {
			return i, true
		}
	}
	if justDigits {
		return 0, false
	}
	for i, each := range numbers {
		if len(each) > len(s) {
			continue
		}
		compared := s[:len(each)]
		if prefix == suffixTarget {
			compared = s[len(s)-len(each):]
		}
		if each == compared {
			return i, true
		}
	}
	return 0, false
}
