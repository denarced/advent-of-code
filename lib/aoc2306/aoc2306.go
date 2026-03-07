package aoc2306

import (
	"math"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func MultiplyCounts(lines []string) int {
	races := parseLines(lines)
	shared.Logger.Info("Multiple counts to win.", "race count", len(races))
	product := 1
	for _, each := range races {
		count := countWaysToWin(each)
		shared.Logger.Info("Race solved.", "race", each, "count", count)
		product *= count
	}
	shared.Logger.Info("Counts multiplied.", "result", product)
	return product
}

func countWaysToWin(aRace race) int {
	shared.Logger.Debug("Find ways to win.", "race", aRace)
	neg, pos := findRoots(aRace)
	shared.Logger.Debug("Roots found.", "neg", neg, "pos", pos)
	low := toRoot(neg, true)
	high := toRoot(pos, false)
	shared.Logger.Debug("Limits found.", "low", low, "high", high)
	return high - low + 1
}

func toRoot(num float64, low bool) int {
	rounded := math.Round(num)
	zero := math.Abs(rounded-num) < 0.001
	if zero {
		return int(rounded)
	}
	if low {
		return int(math.Ceil(num))
	}
	return int(math.Floor(num))
}

func findRoots(aRace race) (neg, pos float64) {
	// x = time - (distance / x) ||*x
	// x^2 - time*x + distance
	// In equation distance = race.distance + 1.
	aRace.distance++
	squareFirst := -aRace.time * -aRace.time
	squareSec := 4 * aRace.distance
	squared := math.Sqrt(float64(squareFirst - squareSec))
	neg = (float64(aRace.time) - squared) / 2
	pos = (float64(aRace.time) + squared) / 2
	return
}

type race struct {
	time     int
	distance int
}

func parseLines(lines []string) []race {
	var times []int
	var distances []int
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		pieces := strings.SplitN(strings.ToLower(trimmed), ":", 2)
		if len(pieces) != 2 {
			shared.Logger.Error("Invalid piece count for line.", "line", each)
			panic("invalid piece count")
		}
		switch pieces[0] {
		case "time":
			times = toInts(pieces[1])
		case "distance":
			distances = toInts(pieces[1])
		default:
			shared.Logger.Error("Unknown line type.", "prefix", pieces[0])
			panic("unknown kind of line")
		}
	}
	if len(times) != len(distances) {
		shared.Logger.Error(
			"Mismatch between times and distances.",
			"time count",
			len(times),
			"distance count",
			len(distances),
		)
		panic("mismatch between times and distances")
	}
	races := make([]race, len(times))
	for i := range len(times) {
		races[i] = race{
			time:     times[i],
			distance: distances[i],
		}
	}
	return races
}

func toInts(s string) []int {
	fields := strings.Fields(s)
	ints := make([]int, len(fields))
	for i, each := range fields {
		value, err := strconv.Atoi(each)
		if err != nil {
			shared.Logger.Error("Failed to convert to int.", "s", each, "err", err)
			panic("failed to convert to int")
		}
		ints[i] = value
	}
	return ints
}
