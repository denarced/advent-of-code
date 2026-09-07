package aoc2215

import (
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func FindNoGoForBeacons(lines []string, y int) int {
	shared.Logger.Info("Find no-go locations for beacon.", "y", y)
	rSet := new(rangeSet)
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
		shared.Logger.Info(
			"Add range.",
			"alpha", alpha,
			"omega", omega,
			"omega-alpha", omega-alpha+1)
		rSet.add([2]int{alpha, omega})
	})
	reservedX.ForEachAll(func(x int) {
		rSet.remove(x)
	})
	var result int
	for _, each := range rSet.ranges {
		result += each[1] - each[0] + 1
	}
	shared.Logger.Info("No-go position count found.", "count", result)
	return result
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

type rangeSet struct {
	ranges [][2]int
}

func (v *rangeSet) add(r [2]int) {
	if v.ranges == nil {
		v.ranges = [][2]int{r}
		return
	}
	v.ranges = append(v.ranges, r)
	slices.SortFunc(v.ranges, func(a, b [2]int) int {
		return a[0] - b[0]
	})
	var i int
	for {
		if i >= len(v.ranges)-1 {
			break
		}
		shared.Logger.Debug("Should join?", "i", v.ranges[i], "i+1", v.ranges[i+1])
		if v.ranges[i][1]+1 >= v.ranges[i+1][0] {
			shared.Logger.Debug("Join ranges.")
			var next [][2]int
			if i > 0 {
				next = append([][2]int(nil), v.ranges[:i]...)
			}
			next = append(next, [2]int{v.ranges[i][0], max(v.ranges[i+1][1], v.ranges[i][1])})
			if i+2 < len(v.ranges) {
				next = append(next, v.ranges[i+2:]...)
			}
			v.ranges = next
		} else {
			i++
		}
	}
}

func (v *rangeSet) remove(value int) {
	if len(v.ranges) == 0 {
		return
	}
	if value < v.ranges[0][0] || v.ranges[len(v.ranges)-1][1] < value {
		return
	}
	for i, each := range v.ranges {
		if each[0] == value && each[1] == value {
			if len(v.ranges) == 1 {
				v.ranges = nil
			} else {
				var next [][2]int
				next = append(next, v.ranges[:i]...)
				next = append(next, v.ranges[i+1:]...)
				v.ranges = next
			}
			return
		}
		if each[0] == value {
			v.ranges[i][0] = value + 1
			return
		} else if each[1] == value {
			v.ranges[i][1] = value - 1
			return
		} else if each[0] < value && value < each[1] {
			shared.Logger.Debug("Charged bodyguard.")
			var next [][2]int
			next = append(next, v.ranges[:i]...)
			next = append(next, [2]int{each[0], value - 1})
			next = append(next, [2]int{value + 1, each[1]})
			next = append(next, v.ranges[i+1:]...)
			v.ranges = next
			return
		}
	}
}
