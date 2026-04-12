package aoc2313

import (
	"github.com/denarced/advent-of-code/shared"
)

func SumReflections(lines []string) int {
	blocks := shared.SplitToBlocks(lines)
	var sum int
	for _, each := range blocks {
		sum += sumReflection(each)
	}
	return sum
}

func sumReflection(lines []string) int {
	if len(lines) < 2 {
		return 0
	}
	shared.Logger.Info(
		"Sum block's reflections.",
		"line count", len(lines),
		"start", shared.SubString(lines[0], 10),
	)

	var sum int
	for i := 0; i < len(lines)-1; i++ {
		if checkHorizontalReflection(lines, i, 1, deriveMaxDistance(i, len(lines))) {
			shared.Logger.Info("Reflection found.", "i", i, "kind", "horizontal")
			sum += 100 * (i + 1)
		}
	}
	for i := 0; i < len(lines[0])-1; i++ {
		if checkVerticalReflection(lines, i, 1, deriveMaxDistance(i, len(lines[0]))) {
			shared.Logger.Info("Reflection found.", "i", i, "kind", "vertical")
			sum += (i + 1)
		}
	}
	return sum
}

func deriveMaxDistance(index, count int) int {
	below := count - index - 2
	return min(index, below) + 1
}

func checkHorizontalReflection(lines []string, row, distance, maxDistance int) bool {
	for d := distance; d <= maxDistance; d++ {
		above := row - (d - 1)
		below := row + d
		if lines[above] != lines[below] {
			return false
		}
	}
	return true
}

func checkVerticalReflection(lines []string, col, distance, maxDistance int) bool {
	for d := distance; d <= maxDistance; d++ {
		left := col - (d - 1)
		right := col + d
		if getColumn(lines, left) != getColumn(lines, right) {
			return false
		}
	}
	return true
}

func getColumn(lines []string, col int) string {
	chars := make([]rune, len(lines))
	for i := range lines {
		c := ([]rune(lines[i]))[col]
		chars[i] = c
	}
	return string(chars)
}
