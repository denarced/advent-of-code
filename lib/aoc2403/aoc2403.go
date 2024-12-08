package aoc2403

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	mulPattern = regexp.MustCompile(`mul\(\d+,\d+\)|don't\(\)|do\(\)`)
)

func Multiply(text string, logic bool) int {
	total := 0
	skipping := false
	for {
		pair := mulPattern.FindStringIndex(text)
		if pair == nil {
			break
		}
		piece := text[pair[0]:pair[1]]
		text = text[pair[1]:]

		if strings.HasPrefix(piece, "don't") {
			skipping = true
			continue
		}
		if strings.HasPrefix(piece, "do") {
			skipping = false
			continue
		}
		if logic && skipping {
			continue
		}
		if strings.HasPrefix(piece, "mul") {
			a, b := splitMul(piece)
			total += a * b
		}
	}
	return total
}

//revive:disable-next-line:confusing-results
func splitMul(s string) (int, int) {
	separated := s[4 : len(s)-1]
	broken := strings.Split(separated, ",")
	return toInt(broken[0]), toInt(broken[1])
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Failed to convert to int: " + s)
	}
	return i
}
