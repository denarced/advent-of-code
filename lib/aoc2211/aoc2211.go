package aoc2211

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func DeriveMonkeyBusiness(lines []string) int {
	monkeys := parseLines(lines)
	for range 20 {
		for i, each := range monkeys {
			for _, value := range each.startItems {
				newValue := each.op(value)
				rounded := newValue / 3

				monkeyID := each.monkeyIDFalse
				if rounded%each.testMod == 0 {
					monkeyID = each.monkeyIDTrue
				}
				if each.ID == monkeyID {
					panic("self send is not monkey-like")
				}

				if shared.IsDebugEnabled() {
					shared.Logger.Debug(
						"New value in the monkey's bag.",
						"monkey", each.ID,
						"value", value,
						"after op", newValue,
						"/3", rounded,
						"send to", monkeyID)
				}
				monkeys[monkeyID].startItems = append(monkeys[monkeyID].startItems, rounded)
			}
			monkeys[i].effortCount += len(each.startItems)
			monkeys[i].startItems = monkeys[i].startItems[:0]
		}
	}
	effortLevels := make([]int, len(monkeys))
	for i, each := range monkeys {
		effortLevels[i] = each.effortCount
	}
	slices.Sort(effortLevels)
	businessLevel := effortLevels[len(effortLevels)-2] * effortLevels[len(effortLevels)-1]
	shared.Logger.Info(
		"Business is done.",
		"effort counts", gent.Map(monkeys, func(each monkey) int { return each.effortCount }),
		"business level", businessLevel)
	return businessLevel
}

func parseLines(lines []string) []monkey {
	groups := groupLines(lines)
	monkeys := make([]monkey, len(groups))
	for i, each := range groups {
		monk := parseMonkeyLines(each)
		monkeys[i] = monk
	}
	shared.Logger.Info("Lines parsed, monkeys shaved.", "monkey count", len(monkeys))
	return monkeys
}

func groupLines(lines []string) [][]string {
	lineSlices := [][]string{}
	var current []string
	for _, each := range lines {
		if strings.HasPrefix(each, "Monkey") {
			if len(current) > 0 {
				lineSlices = append(lineSlices, current)
				if len(current) != 6 {
					panic("should have 6 lines")
				}
			}
			current = nil
		}
		current = append(current, each)
	}
	if len(current) > 0 {
		lineSlices = append(lineSlices, current)
	}
	return lineSlices
}

type operation func(old int) int

type monkey struct {
	ID            int
	startItems    []int
	op            operation
	testMod       int
	monkeyIDTrue  int
	monkeyIDFalse int
	effortCount   int
}

func parseMonkeyLines(lines []string) monkey {
	assert := func(line string) func(count int, err error) func(expectedCount int) {
		return func(count int, err error) func(expectedCount int) {
			logger := shared.Logger.With("line", line, "count", count, "err", err)
			if err != nil {
				logger.Error("Failed to parse a monkey, non-nil error.")
				panic("failed to parse a monkey")
			}
			return func(expectedCount int) {
				if count != expectedCount {
					logger.Error(
						"Failed to parse a monkey, count != expected.",
						"expected count",
						expectedCount,
					)
				}
			}
		}
	}
	var monk monkey
	for i, each := range lines {
		switch i {
		case 0:
			assert(each)(fmt.Sscanf(each, "Monkey %d:", &monk.ID))(1)
		case 1:
			withoutSpaces := strings.ReplaceAll(each, " ", "")
			monk.startItems = gent.Map(strings.Split(strings.Split(withoutSpaces, ":")[1], ","),
				func(s string) int {
					value, err := strconv.Atoi(s)
					if err != nil {
						shared.Logger.Error(
							"Failed to parse starting item integer.",
							"str", s,
							"err", err)
						panic(err)
					}
					return value
				})
		case 2:
			monk.op = parseOperation(each)
		case 3:
			assert(each)(fmt.Sscanf(each, "Test: divisible by %d", &monk.testMod))(1)
		case 4:
			assert(each)(fmt.Sscanf(each, "If true: throw to monkey %d", &monk.monkeyIDTrue))(1)
		case 5:
			assert(each)(fmt.Sscanf(each, "If false: throw to monkey %d", &monk.monkeyIDFalse))(1)
		default:
			panic("unknown monkey line: " + each)
		}
	}
	return monk
}

func parseOperation(line string) operation {
	mainPieces := strings.Split(line, "=")
	if len(mainPieces) != 2 {
		shared.Logger.Error("Invalid operation line.", "line", line, "piece count", len(mainPieces))
		panic("invalid operation line")
	}
	pieces := strings.Fields(strings.TrimSpace(mainPieces[1]))
	if len(pieces) != 3 {
		shared.Logger.Error(
			"Invalid operation line formula.",
			"formula", mainPieces[1],
			"pieces", pieces)
		panic("invalid operation formula")
	}
	return func(old int) int {
		var other int
		if pieces[2] == "old" {
			other = old
		} else {
			i, err := strconv.Atoi(pieces[2])
			if err != nil {
				shared.Logger.Error(
					"Failed parse fixed formula value.",
					"err", err,
					"pieces", pieces)
				panic(err)
			}
			other = i
		}
		switch pieces[1] {
		case "*":
			return old * other
		case "+":
			return old + other
		default:
			panic("unknown arithmetic operation: " + pieces[1])
		}
	}
}
