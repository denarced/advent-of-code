package aoc2215

import (
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func FindNoGoForBeacons(lines []string, y int) int {
	xSet := gent.NewSet[int]()
	reservedX := gent.NewSet[int]()
	parseLines(lines, func(sens sensor) {
		if sens.loc.Y == y {
			reservedX.Add(sens.loc.X)
		}
		if sens.closest.Y == y {
			reservedX.Add(sens.closest.X)
		}
		distance := calculateManhattanDistance(sens)
		centerToY := shared.Abs(sens.loc.Y - y)
		logger := shared.Logger.With("sensor", sens, "distance", distance)
		if distance < centerToY {
			logger.Debug("Skipped.")
			return
		}
		alpha, omega := deriveRange(sens.loc.X, distance-centerToY)
		logger.Debug("Add points.", "alpha", alpha, "omega", omega)
		for i := alpha; i <= omega; i++ {
			xSet.Add(i)
		}
	})
	reservedX.ForEachAll(func(x int) {
		xSet.Remove(x)
	})
	xCoords := xSet.ToSlice()
	slices.Sort(xCoords)
	result := xSet.Count()
	shared.Logger.Info("No-go position count found.", "count", result, "x-coords", xCoords)
	return xSet.Count()
}

type sensor struct {
	loc     shared.Loc
	closest shared.Loc
}

func parseLines(lines []string, cb func(sensor)) {
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		parseLine(trimmed, cb)
	}
}

func parseLine(line string, cb func(sensor)) {
	var sens sensor
	indexes := findIndexes(line, '=')
	var values []int
	for _, each := range indexes {
		values = append(values, grabInt(line, each+1))
	}
	sens.loc.X = values[0]
	sens.loc.Y = values[1]
	sens.closest.X = values[2]
	sens.closest.Y = values[3]
	cb(sens)
}

func findIndexes(line string, char rune) []int {
	var indexes []int
	for i, c := range line {
		if c == char {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func grabInt(line string, start int) int {
	var chars []rune
	full := []rune(line)
	for i := 0; ; i++ {
		lineIndex := start + i
		if lineIndex >= len(line) {
			break
		}
		c := full[lineIndex]
		if c != '-' && c < '0' || '9' < c {
			break
		}
		chars = append(chars, c)
	}
	i, err := strconv.Atoi(string(chars))
	if err != nil {
		shared.Logger.Error("Failed grab int.", "line", line, "start", start, "chars", chars)
		panic(err)
	}
	return i
}

func calculateManhattanDistance(sens sensor) int {
	x := sens.loc.X - sens.closest.X
	y := sens.loc.Y - sens.closest.Y
	return shared.Abs(x) + shared.Abs(y)
}

func deriveRange(centerX, overlap int) (from, to int) {
	return centerX - overlap, centerX + overlap
}
