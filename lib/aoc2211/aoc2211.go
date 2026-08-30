package aoc2211

import (
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func DeriveMonkeyBusiness(lines []string, roundCount int, worries bool) int {
	monkeys := parseLines(lines)

	var commonDiv *big.Int
	for _, each := range monkeys {
		if commonDiv == nil {
			commonDiv = new(big.Int).Set(each.testMod)
		} else {
			commonDiv.Mul(commonDiv, each.testMod)
		}
	}

	bigZero := big.NewInt(0)
	bigOne := big.NewInt(1)
	bigThree := big.NewInt(3)
	bigTen := big.NewInt(10)
	bigMod := new(big.Int)
	for round := 1; round <= roundCount; round++ {
		for i, each := range monkeys {
			for _, value := range each.startItems {
				each.op(value)
				if !worries {
					value.Div(value, bigThree)
				}

				monkeyID := each.monkeyIDFalse
				if bigMod.Mod(value, each.testMod).Cmp(bigZero) == 0 {
					monkeyID = each.monkeyIDTrue
				}
				if each.ID == monkeyID {
					panic("self send is not monkey-like")
				}
				monkeys[monkeyID].startItems = append(monkeys[monkeyID].startItems, value)
			}
			monkeys[i].effortCount += len(each.startItems)
			monkeys[i].startItems = monkeys[i].startItems[:0]
		}
		if !worries {
			continue
		}
		for _, monk := range monkeys {
			for _, startingValue := range monk.startItems {
				if bigMod.Div(startingValue, commonDiv).Cmp(bigTen) > 0 {
					bigMod.Sub(bigMod, bigOne)
					bigMod.Mul(bigMod, commonDiv)
					startingValue.Sub(startingValue, bigMod)
				}
			}
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

type operation func(old *big.Int)

type monkey struct {
	ID            int
	startItems    []*big.Int
	op            operation
	testMod       *big.Int
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
				func(s string) *big.Int {
					value, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						shared.Logger.Error(
							"Failed to parse starting item integer.",
							"str", s,
							"err", err)
						panic(err)
					}
					return big.NewInt(value)
				})
		case 2:
			monk.op = parseOperation(each)
		case 3:
			var testMod int64
			assert(each)(fmt.Sscanf(each, "Test: divisible by %d", &testMod))(1)
			monk.testMod = big.NewInt(testMod)
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
	var fixedValue *big.Int
	if pieces[2] != "old" {
		i, success := new(big.Int).SetString(pieces[2], 10)
		if !success {
			shared.Logger.Error(
				"Failed parse fixed formula value.",
				"pieces", pieces)
			panic("failed to parse fixed formula value")
		}
		fixedValue = i
	}

	return func(old *big.Int) {
		var other *big.Int
		if fixedValue == nil {
			other = old
		} else {
			other = fixedValue
		}
		switch pieces[1] {
		case "*":
			old.Mul(old, other)
		case "+":
			old.Add(old, other)
		default:
			panic("unknown arithmetic operation: " + pieces[1])
		}
	}
}
