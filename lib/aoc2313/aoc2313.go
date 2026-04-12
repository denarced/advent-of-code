package aoc2313

import (
	"github.com/denarced/advent-of-code/shared"
)

func SumReflections(lines []string, fixSmudge bool) int {
	blocks := shared.SplitToBlocks(lines)
	var sum int
	for _, each := range blocks {
		sum += sumReflection(each, fixSmudge)
	}
	return sum
}

func sumReflection(lines []string, fixSmudge bool) int {
	if len(lines) < 2 {
		return 0
	}
	shared.Logger.Info(
		"Sum block's reflections.",
		"line count", len(lines),
		"start", shared.SubString(lines[0], 10),
		"fix smudges", fixSmudge,
	)

	sum := findHorizontalReflection(lines, fixSmudge)
	if sum > 0 {
		shared.Logger.Info("Reflection sum found.", "sum", sum, "vertical", false)
		return sum
	}
	sum = findVerticalReflection(lines, fixSmudge)
	if sum > 0 {
		shared.Logger.Info("Reflection sum found.", "sum", sum, "vertical", true)
		return sum
	}
	panic("no reflection found")
}

func findHorizontalReflection(lines []string, fixSmudge bool) (sum int) {
	shared.Logger.Debug("Find horizontal reflection.", "first", shared.SubString(lines[0], 10))
	return findReflection(
		lines,
		fixSmudge,
		len(lines)-1,
		false,
		func(i int) int { return 100 * (i + 1) },
	)
}

func findVerticalReflection(lines []string, fixSmudge bool) (sum int) {
	shared.Logger.Debug("Find vertical reflection.", "first", shared.SubString(lines[0], 10))
	return findReflection(
		lines,
		fixSmudge,
		len(lines[0])-1,
		true,
		func(i int) int { return i + 1 },
	)
}

func findReflection(
	lines []string,
	fixSmudge bool,
	maximum int,
	vertical bool,
	deriveSum func(int) int,
) (sum int) {
	for i := range maximum {
		ok := checkReflection(
			lines,
			i,
			1,
			deriveMaxDistance(i, maximum+1),
			vertical,
			fixSmudge,
		)
		if ok {
			shared.Logger.Info("Reflection found.", "i", i, "vertical", vertical)
			return deriveSum(i)
		}
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Failed to find reflection.",
				"start", shared.SubString(lines[0], 10),
				"i", i,
				"vertical", vertical)
		}
	}
	return
}

func deriveMaxDistance(index, count int) int {
	below := count - index - 2
	return min(index, below) + 1
}

func countDifference(first, second string, maximum int) int {
	var count int
	for i := range first {
		if first[i] != second[i] {
			count++
			if count >= maximum {
				return count
			}
		}
	}
	return count
}

func checkReflection(
	lines []string,
	start, distance, maxDistance int,
	vertical, fixSmudge bool,
) (ok bool) {
	var diffDistance int
	for d := distance; d <= maxDistance; d++ {
		firstIndex := start - (d - 1)
		secondIndex := start + d
		var first, second string
		if vertical {
			first = getColumn(lines, firstIndex)
			second = getColumn(lines, secondIndex)
		} else {
			first = lines[firstIndex]
			second = lines[secondIndex]
		}
		switch difference := countDifference(first, second, 2); difference {
		case 0:
			continue
		case 1:
			if !fixSmudge {
				return
			}
			if diffDistance != 0 {
				return
			}
			diffDistance = difference
		default:
			return
		}
	}
	if fixSmudge {
		return diffDistance == 1
	}
	return diffDistance == 0
}

func getColumn(lines []string, col int) string {
	chars := make([]rune, len(lines))
	for i := range lines {
		c := ([]rune(lines[i]))[col]
		chars[i] = c
	}
	return string(chars)
}
